package config

import (
	"context"

	"fmt"
	"log"

	"github.com/go-redis/redis/v8"
)

type ConfigRedis struct {
	REDIS_HOST     string
	REDIS_PORT     string
	REDIS_USER     string
	REDIS_PASSWORD string
	REDIS_DB       string
}

func (config *ConfigRedis) InitRedisDatabase() *redis.Client {
	host := config.REDIS_HOST
	port := config.REDIS_PORT
	// db := config.REDIS_DB
	user := config.REDIS_USER
	password := config.REDIS_PASSWORD

	address := fmt.Sprintf("%s:%s", host, port)

	client := redis.NewClient(&redis.Options{
		Addr:               address,
		Username:           user,
		Password:           password,
	})

	if _, err := client.Ping(context.Background()).Result(); err != nil {
		log.Fatalf("error connecting to redis: %v", err)
	}

	log.Print("connected with redis")

	return client
}

// for PRODUCTION
// func (config *ConfigRedis) InitRedisDatabase() *redis.Client {
// 	host := config.REDIS_HOST
// 	port := config.REDIS_PORT
// 	db := config.REDIS_DB
// 	user := config.REDIS_USER
// 	password := config.REDIS_PASSWORD

// 	address := fmt.Sprintf("%s:%s", host, port)

// 	redisDB, _ := strconv.Atoi(db)

// 	client := redis.NewClient(&redis.Options{
// 		Addr:               address,
// 		Username:           user,
// 		Password:           password,
// 		DB:                 redisDB,
// 		TLSConfig:          &tls.Config{},
// 	})

// 	if _, err := client.Ping(context.Background()).Result(); err != nil {
// 		log.Fatalf("error connecting to redis: %v", err)
// 	}

// 	log.Print("connected with redis")

// 	return client
// }