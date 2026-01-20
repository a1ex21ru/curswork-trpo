package postgres

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Client struct {
	pool *pgxpool.Pool
}

func NewClient(ctx context.Context) (*Client, error) {
	const fn = "NewClient"

	host := os.Getenv("POSTGRES_HOST")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbName := os.Getenv("POSTGRES_DB")

	connString := fmt.Sprintf("postgres://%s:%s@%s:5432/%s?sslmode=disable", user, password, host, dbName)

	conn, err := pgx.Connect(ctx, connString)
	if err != nil {
		return nil, fmt.Errorf("%s : %w", fn, err)
	}

	err = conn.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s : %w", fn, err)
	}

	poolCfg, err := newConfig(ctx, connString)

	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		return nil, fmt.Errorf("%s : %w", fn, err)
	}

	return &Client{
		pool: pool,
	}, nil
}

// newConfig creates a new pgxpool.Config object based on the provided configuration.
func newConfig(ctx context.Context, dataBaseURL string) (*pgxpool.Config, error) {
	const fn = "newConfig"

	dbConfig, err := pgxpool.ParseConfig(dataBaseURL)
	if err != nil {
		return nil, err
	}

	setWithConfig(dbConfig)

	dbConfig.BeforeAcquire = func(_ context.Context, _ *pgx.Conn) bool {
		return true
	}

	dbConfig.AfterRelease = func(_ *pgx.Conn) bool {
		return true
	}

	dbConfig.BeforeClose = func(_ *pgx.Conn) {
		log.Println(ctx, "Closed the connection pool to the database!")
	}

	return dbConfig, nil
}

func setWithConfig(dbConfig *pgxpool.Config) {
	dbConfig.MaxConns = 50

	dbConfig.MinConns = 1

	dbConfig.MaxConnLifetime = time.Hour

	dbConfig.MaxConnIdleTime = time.Minute * 30

	dbConfig.HealthCheckPeriod = time.Minute

	dbConfig.ConnConfig.ConnectTimeout = time.Second * 5
}

func (c *Client) Close(ctx context.Context) error {
	c.pool.Close()
	return nil
}

// Query executes a query that returns rows
func (c *Client) Query(ctx context.Context, query string, args ...interface{}) (pgx.Rows, error) {
	return c.pool.Query(ctx, query, args...)
}

// QueryRow executes a query that returns a single row
func (c *Client) QueryRow(ctx context.Context, query string, args ...interface{}) pgx.Row {
	return c.pool.QueryRow(ctx, query, args...)
}

// Exec executes a query that doesn't return rows
func (c *Client) Exec(ctx context.Context, query string, args ...interface{}) (pgconn.CommandTag, error) {
	return c.pool.Exec(ctx, query, args...)
}

// Ping checks the database connection
func (c *Client) Ping(ctx context.Context) error {
	return c.pool.Ping(ctx)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
