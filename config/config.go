package config

import (
	"os"
	"strconv"
	"strings"
	"time"
)

var config = map[string]string{
	"APP_ENV":           "local",
	"APP_SERVICE_TOKEN": "xpto",
	"APP_NAME":          "users",
	"ENABLE_SWAGGER":    "true",

	// LOG LEVEL
	"LOG_LEVEL": "debug",

	// HTTP CONFIGURATION
	"HTTP_PORT":                           "4040",
	"HTTP_PORT_OPEN_API":                  "5050",
	"HTTP_SERVER_READ_TIMEOUT_SECONDS":    "30",
	"HTTP_SERVER_WRITE_TIMEOUT_SECONDS":   "30",
	"HTTP_SERVER_MAX_IDLE_CONNS":          "3",
	"HTTP_SERVER_MAX_IDLE_CONNS_PER_HOST": "2",

	// HTTP CLIENT CONFIGURATION
	"HTTP_CLIENT_READ_TIMEOUT_SECONDS":    "30",
	"HTTP_CLIENT_CONN_TIMEOUT_SECONDS":    "30",
	"HTTP_CLIENT_MAX_IDLE_CONNS":          "3",
	"HTTP_CLIENT_MAX_IDLE_CONNS_PER_HOST": "2",

	// RABBITMQ
	"RABBITMQ_URI": "amqp://guest:guest@rabbitmq:5672/",
	"BATCH_SIZE":   "10000",

	// POSTGRES DATABASE CONNECTION DEV
	"POSTGRES_TAG":                 "15",
	"POSTGRES_USER":                "user-app",
	"POSTGRES_PASSWORD":            "VmY3JHIydUt3UDlz",
	"POSTGRES_DB":                  "user-dev-db",
	"POSTGRES_PORT":                "5432",
	"DATABASE_HOST":                "localhost",
	"POSTGRES_POOL_SIZE":           "5",
	"POSTGRES_CONN_MAX_TTL_MILLIS": "1800000",

	// REDIS
	"REDIS_ADDR":             "localhost:6379",
	"REDIS_POOL_SIZE":        "10",
	"REDIS_TIMEOUT_MS":       "1000",
	"REDIS_MAX_RETRIES":      "1",
	"REDIS_WRITE_TIMEOUT_MS": "1000",
}

// GetString value of a given env var
func GetString(k string) string {
	v := os.Getenv(k)
	if v == "" {
		return config[k]
	}

	return v
}

// GetInt value of a given env var
func GetInt(k string) int {
	v := GetString(k)
	i, err := strconv.Atoi(v)
	if err != nil {
		panic(err)
	}

	return i
}

// GetDuration value of a given env var
func GetDuration(k string) time.Duration {
	return time.Duration(GetInt(k)) * time.Millisecond
}

// GetBool value of a given env var
func GetBool(k string) bool {
	v := GetString(k)
	return strings.ToLower(v) == "true"
}
