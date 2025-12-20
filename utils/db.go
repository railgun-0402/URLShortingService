package utils

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"
)

func OpenDB() (*sql.DB, error) {
	dsn := os.Getenv("DATABASE_URL") //
	if dsn == "" {
		// local
		dsn = "postgres://urlshort:urlshort@localhost:5432/urlshort?sslmode=disable"
	} // postgres://user:pass@host:5432/db?sslmode=require
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("failed to open db: %v", err)
	}

	// この辺の数字はお試し実装なので、任意で設定している
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(30 * time.Second)

	// Test Connection
	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("failed to ping db: %w", err)
	}

	log.Printf("connected: %s\n", maskDSN(dsn))
	return db, nil
}

func maskDSN(dsn string) string {
	// 表示を隠す
	if len(dsn) > 10 {
		return dsn[:20] + "..." + dsn[len(dsn)-10:]
	}
	return dsn
}
