package ioc

import (
	"database/sql"
	"github.com/julioisaac/users/config"
	"github.com/julioisaac/users/logger"
	"github.com/julioisaac/users/providers/cache"
	"github.com/julioisaac/users/providers/database"
	rabbitMQ "github.com/julioisaac/users/providers/rabbitmq"
	"sync"
)

var (
	dbInstanceOnce sync.Once
	dbInstance     *sql.DB

	redisCacheOnce sync.Once
	redisCache     cache.Cache

	rabbitmqOnce sync.Once
	rabbitmq     rabbitMQ.RabbitMQProvider
)

// DB instance
func DB() *sql.DB {
	dbInstanceOnce.Do(func() {
		cfg := database.DBConfig{
			ServiceName: "users-api",
			Host:        config.GetString("DATABASE_HOST"),
			Port:        config.GetString("POSTGRES_PORT"),
			User:        config.GetString("POSTGRES_USER"),
			Password:    config.GetString("POSTGRES_PASSWORD"),
			DBName:      config.GetString("POSTGRES_DB"),
			PoolSize:    config.GetInt("POSTGRES_POOL_SIZE"),
			ConnMaxTTL:  config.GetDuration("POSTGRES_CONN_MAX_TTL_MILLIS"),
		}

		var err error
		dbInstance, err = database.OpenConnection(&cfg)
		if err != nil {
			logger.Logger.Errorf("Error creating postgres connection, err: %v", err)
		}
	})
	return dbInstance
}

// RedisCache with dependencies
func RedisCache() cache.Cache {
	redisCacheOnce.Do(func() {
		redisCache = cache.NewRedisCache()
	})

	return redisCache
}

// RabbitMQ with dependencies
func RabbitMQ() rabbitMQ.RabbitMQProvider {
	rabbitmqOnce.Do(func() {
		rabbitmq = rabbitMQ.NewRabbit()
	})

	return rabbitmq
}
