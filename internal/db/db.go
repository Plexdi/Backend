package db

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var Pool *pgxpool.Pool

func ConnectDB() error {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return fmt.Errorf("DATABASE_URL not set")
	}

	// Parse pool config
	config, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		return fmt.Errorf("❌ failed to parse config: %w", err)
	}

	// Same dialer as before (force IPv4, timeout)
	dialer := &net.Dialer{
		Timeout:   5 * time.Second,
		DualStack: true,
	}

	config.ConnConfig.DialFunc = func(ctx context.Context, network, addr string) (net.Conn, error) {
		// use tcp4 so Neon/Heroku don’t get weird with IPv6
		return dialer.DialContext(ctx, "tcp4", addr)
	}

	// Create a *pool*, not a single connection
	Pool, err = pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return fmt.Errorf("❌ failed to connect: %w", err)
	}

	// Test connection with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := Pool.Ping(ctx); err != nil {
		return fmt.Errorf("❌ ping failed: %w", err)
	}

	log.Println("✅ Connected to Neon Console database successfully!")
	return nil
}
