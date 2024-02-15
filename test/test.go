package test

import (
	"database/sql"
	"github.com/google/go-cmp/cmp"
	"github.com/julioisaac/users/config"
	"github.com/julioisaac/users/providers/database"
	"testing"
)

var dbInstance *sql.DB

func GetDBInstance() (*sql.DB, error) {

	if dbInstance != nil {
		return dbInstance, nil
	}
	cfg := GetDBConfig()
	dbInstance, _ = database.OpenConnection(&cfg)

	return dbInstance, nil
}

func GetDBConfig() database.DBConfig {
	return database.DBConfig{
		Host:       config.GetString("DATABASE_HOST"),
		Port:       config.GetString("POSTGRES_PORT"),
		User:       config.GetString("POSTGRES_USER"),
		Password:   config.GetString("POSTGRES_PASSWORD"),
		DBName:     config.GetString("POSTGRES_DB"),
		PoolSize:   config.GetInt("POSTGRES_POOL_SIZE"),
		ConnMaxTTL: config.GetDuration("POSTGRES_CONN_MAX_TTL_MILLIS"),
	}
}

// AssertDomainError message
func AssertDomainError(t *testing.T, err error, expectedErr error) {
	t.Helper()
	if diff := cmp.Diff(err, expectedErr); diff != "" {
		t.Errorf("unexpected domain error %s", diff)
	}
}
