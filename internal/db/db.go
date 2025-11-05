package db

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
)

var Conn *pgx.Conn

func ConnectDB() error {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return fmt.Errorf("DATABASE_URL not set")
	}

	// Parse connection config
	config, err := pgx.ParseConfig(dbURL)
	if err != nil {
		return fmt.Errorf("❌ failed to parse config: %w", err)
	}

	// Use a dialer that supports both IPv4 & IPv6 with timeout
	dialer := &net.Dialer{
		Timeout:   5 * time.Second,
		DualStack: true, // ✅ allows both IPv4 & IPv6
	}

	config.DialFunc = func(ctx context.Context, network, addr string) (net.Conn, error) {
		return dialer.DialContext(ctx, "tcp", addr)
	}

	// Try connecting
	Conn, err = pgx.ConnectConfig(context.Background(), config)
	if err != nil {
		return fmt.Errorf("❌ failed to connect: %w", err)
	}

	// Test connection
	if err := Conn.Ping(context.Background()); err != nil {
		return fmt.Errorf("❌ ping failed: %w", err)
	}

	log.Println("✅ Connected to Supabase database successfully!")
	return nil
}
