package db

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

func InitializeDB(dsn string) (*sql.DB, error) {
	return initializeWithRetry(dsn, 5, 2*time.Second)
}

func InitializeTestDB(dsn string) (*sql.DB, error) {
	return initializeWithRetry(dsn, 3, 1*time.Second)
}

func initializeWithRetry(dsn string, maxRetries int, delay time.Duration) (*sql.DB, error) {
	var db *sql.DB
	var err error

	if dsn == "" {
		return nil, fmt.Errorf("no postgres DSN provided")
	}

	for i := 0; i < maxRetries; i++ {
		db, err = sql.Open("postgres", dsn)
		if err == nil {
			err = db.Ping()
			if err == nil {
				return db, nil
			}
		}
		time.Sleep(delay)
	}

	return nil, fmt.Errorf("failed to connect to DB after %d retries: %w", maxRetries, err)
}
