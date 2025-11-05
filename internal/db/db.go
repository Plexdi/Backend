package db

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/jackc/pgx/v5"
)

var Conn *pgx.Conn

func ConnectDB() error {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return fmt.Errorf("DATABASE_URL not set")
	}

	// Force IPv4 resolution (fixes "no such host" on Windows/Supabase)
	dialer := &net.Dialer{DualStack: false}
	customDialer := func(ctx context.Context, network, addr string) (net.Conn, error) {
		return dialer.DialContext(ctx, "tcp4", addr)
	}

	config, err := pgx.ParseConfig(dbURL)
	if err != nil {
		return fmt.Errorf("❌ failed to parse config: %w", err)
	}
	config.DialFunc = customDialer

	Conn, err = pgx.ConnectConfig(context.Background(), config)
	if err != nil {
		return fmt.Errorf("❌ failed to connect: %w", err)
	}

	if err := Conn.Ping(context.Background()); err != nil {
		return fmt.Errorf("❌ ping failed: %w", err)
	}

	log.Println("✅ Connected to Supabase database successfully!")
	return nil
}
