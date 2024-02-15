package database

import (
	"database/sql"
	"fmt"
	"github.com/julioisaac/users/logger"
	"time"

	_ "github.com/lib/pq"
)

type DBConfig struct {
	ServiceName string
	Host        string
	Port        string
	User        string
	Password    string
	DBName      string
	PoolSize    int
	ConnMaxTTL  time.Duration
}

func (c DBConfig) Dsn() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", c.Host, c.Port, c.User, c.Password, c.DBName)
}

// OpenConnection given config
func OpenConnection(c *DBConfig) (*sql.DB, error) {

	db, err := sql.Open("postgres", c.Dsn())
	if err != nil {
		logger.Logger.Fatalf("error opening potgres connection, err: %v", err)
	}

	db.SetMaxIdleConns(c.PoolSize)
	db.SetMaxOpenConns(c.PoolSize)
	db.SetConnMaxLifetime(c.ConnMaxTTL)

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
