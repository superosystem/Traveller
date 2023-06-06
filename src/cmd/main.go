package main

import (
	"cloud.google.com/go/storage"
	"context"
	"github.com/labstack/echo/v4"
	"github.com/superosystem/trainingsystem-backend/src/common/config"
	"github.com/superosystem/trainingsystem-backend/src/middlewares"
	"github.com/superosystem/trainingsystem-backend/src/routes"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type operation func(context.Context) error

func main() {
	// init mysql config
	configMySQL := config.MySQLConfig{
		MYSQL_USERNAME: config.GetEnv("MYSQL_USERNAME"),
		MYSQL_PASSWORD: config.GetEnv("MYSQL_PASSWORD"),
		MYSQL_HOST:     config.GetEnv("MYSQL_HOST"),
		MYSQL_PORT:     config.GetEnv("MYSQL_PORT"),
		MYSQL_NAME:     config.GetEnv("MYSQL_NAME"),
	}

	mysqlDB := configMySQL.InitMySQLDatabase()

	config.MySQLMigrate(mysqlDB)

	// init redis config
	configRedis := config.RedisConfig{
		REDIS_HOST:     config.GetEnv("REDIS_HOST"),
		REDIS_PORT:     config.GetEnv("REDIS_PORT"),
		REDIS_PASSWORD: config.GetEnv("REDIS_PASSWORD"),
		REDIS_DB:       config.GetEnv("REDIS_DB"),
	}

	redisDB := configRedis.InitRedisDatabase()

	// init jwt config
	jwtConfig := config.NewJWTConfig(config.GetEnv("JWT_SECRET"))

	// init mailer config
	mailerConfig := config.NewMailer(
		config.GetEnv("SMTP_HOST"),
		config.GetEnv("SMTP_PORT"),
		config.GetEnv("EMAIL_SENDER_NAME"),
		config.GetEnv("AUTH_EMAIL"),
		config.GetEnv("AUTH_PASSWORD_EMAIL"),
	)

	ctx := context.Background()

	// init cloud storage config
	storageClient, _ := storage.NewClient(ctx)

	storageConfig := config.NewCloudStorage(storageClient, config.GetEnv("BUCKET_NAME"))

	e := echo.New()

	// CORS
	e.Use(middlewares.CORS())

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
		if err := e.Start(config.GetEnv("APP_PORT")); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	wait := gracefulShutdown(ctx, 5*time.Second, map[string]operation{
		"mysql": func(ctx context.Context) error {
			return config.MySQLClose(mysqlDB)
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
