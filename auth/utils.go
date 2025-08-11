package auth

import (
	"context"
	"fmt"
	"time"

	"zombiezen.com/go/sqlite"
)

// getConnection returns a new connection to the database.
// This connection times out after 5 seconds, to avoid blocking the server.
// It also returns a cleanup function to close the connection after use.
func getConnection() (*sqlite.Conn, func(), error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	conn, err := dbpool.Take(ctx)
	if err != nil {
		cancel()
		return nil, nil, err
	}

	cleanup := func() {
		dbpool.Put(conn)
		cancel()
	}
	return conn, cleanup, nil
}

// parseTimestamp parses a timestamp string from the database.
func parseTimestamp(ts string) (time.Time, error) {
	t, err := time.Parse("2006-01-02 15:04:05", ts)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to parse timestamp: %w", err)
	}
	return t, nil
}
