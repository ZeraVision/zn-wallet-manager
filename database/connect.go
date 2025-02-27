package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

// Database defines the interface for database operations
type Database interface {
	Exec(ctx context.Context, query string, args ...interface{}) error
}

// PgxPoolWrapper wraps a pgxpool.Pool to implement the Database interface
type PgxPoolWrapper struct {
	Pool *pgxpool.Pool
}

// Exec executes a query
func (p *PgxPoolWrapper) Exec(ctx context.Context, query string, args ...interface{}) error {
	_, err := p.Pool.Exec(ctx, query, args...)
	return err
}

// Global pool variable
var pool *PgxPoolWrapper

// Connect initializes the connection to the database
func Connect() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	dbName := os.Getenv("DB_NAME")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")

	connStr := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", user, password, host, port, dbName)
	pgxPool, err := pgxpool.Connect(ctx, connStr)
	if err != nil {
		log.Printf("Failed to connect to database %s: %v\n", dbName, err)
		return err
	}

	pool = &PgxPoolWrapper{Pool: pgxPool}
	log.Printf("Connected to database %s successfully\n", dbName)
	return nil
}

// Close shuts down the database connection
func Close() {
	if pool != nil {
		pool.Pool.Close()
		log.Println("Database connection closed.")
	}
}

// Get returns the Database object with the real connection
func Get() Database {
	return pool
}
