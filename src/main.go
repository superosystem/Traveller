package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"cloud.google.com/go/storage"

	"github.com/labstack/echo/v4"

	"github.com/superosystem/trainingsystem-backend/src/config"
	"github.com/superosystem/trainingsystem-backend/src/middleware"
	"github.com/superosystem/trainingsystem-backend/src/routes"
)

type operation func(context.Context) error

func main() {
	// init mysql config
	configMySQL := config.ConfigMySQL{
		MYSQL_USERNAME: config.GetConfig("MYSQL_USERNAME"),
		MYSQL_PASSWORD: config.GetConfig("MYSQL_PASSWORD"),
		MYSQL_HOST:     config.GetConfig("MYSQL_HOST"),
		MYSQL_PORT:     config.GetConfig("MYSQL_PORT"),
		MYSQL_NAME:     config.GetConfig("MYSQL_NAME"),
	}

	mysqlDB := configMySQL.InitMySQLDatabase()

	config.MySQLMigrate(mysqlDB)

	// init redis config
	configRedis := config.ConfigRedis{
		REDIS_HOST:     config.GetConfig("REDIS_HOST"),
		REDIS_PORT:     config.GetConfig("REDIS_PORT"),
		REDIS_PASSWORD: config.GetConfig("REDIS_PASSWORD"),
		REDIS_DB:       config.GetConfig("REDIS_DB"),
	}

	redisDB := configRedis.InitRedisDatabase()

	// init jwt config
	jwtConfig := config.NewJWTConfig(config.GetConfig("JWT_SECRET"))

	// init mailer config
	mailerConfig := config.NewMailer(
		config.GetConfig("SMTP_HOST"),
		config.GetConfig("SMTP_PORT"),
		config.GetConfig("EMAIL_SENDER_NAME"),
		config.GetConfig("AUTH_EMAIL"),
		config.GetConfig("AUTH_PASSWORD_EMAIL"),
	)

	ctx := context.Background()

	// init cloud storage config
	storageClient, _ := storage.NewClient(ctx)

	storageConfig := config.NewCloudStorage(storageClient, config.GetConfig("BUCKET_NAME"))

	e := echo.New()

	// CORS
	e.Use(middleware.CORS())

	// init routes config
	route := routes.RouteConfig{
		Echo:          e,
		MySQLDB:       mysqlDB,
		RedisDB:       redisDB,
		JWTConfig:     jwtConfig,
		Mailer:        mailerConfig,
		StorageConfig: storageConfig,
	}

	route.New()

	go func() {
		if err := e.Start(config.GetConfig("APP_PORT")); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	wait := gracefulShutdown(ctx, 5*time.Second, map[string]operation{
		"mysql": func(ctx context.Context) error {
			return config.CloseDB(mysqlDB)
		},
		"http-server": func(ctx context.Context) error {
			return e.Shutdown(ctx)
		},
	})

	<-wait
}

// graceful shutdown perform application shutdown gracefully
func gracefulShutdown(ctx context.Context, timeout time.Duration, ops map[string]operation) <-chan struct{} {
	wait := make(chan struct{})

	go func() {
		s := make(chan os.Signal, 1)

		// add any other syscall that you want to be notified with
		signal.Notify(s, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP)
		<-s

		log.Println("shutting down")

		// set timeout for the ops to be done to prevent system hang
		timeoutFunc := time.AfterFunc(timeout, func() {
			log.Printf("timeout %d ms has been elapsed, force exit", timeout.Milliseconds())
			os.Exit(0)
		})

		defer timeoutFunc.Stop()

		var wg sync.WaitGroup

		// do the operation asynchronously to save time
		for key, op := range ops {
			wg.Add(1)
			innerOp := op
			innerKey := key
			go func() {
				defer wg.Done()

				log.Printf("cleaning up: %v", innerKey)

				if err := innerOp(ctx); err != nil {
					log.Printf("%s: clean up failed: %s", innerKey, err.Error())
					return
				}

				log.Printf("%s was shutdown gracefully", innerKey)
			}()
		}

		wg.Wait()

		close(wait)
	}()

	return wait
}
