package auth

import (
	"fmt"
	"time"
)

const TokenLifespan = (30 * 24 * time.Hour)

// LoginUser creates an access token for the given user
func LoginUser(username string) (string, error) {
	conn, cleanup, err := getConnection()
	if err != nil {
		return "", err
	}
	defer cleanup()

	// Generate access token
	token, err := generateAccessToken()
	if err != nil {
		return "", fmt.Errorf("failed to generate access token: %w", err)
	}

	stmt := conn.Prep(`
		INSERT INTO access_tokens (user_id, token)
		SELECT user_id, ?
		FROM users
		WHERE username = ?
	`)
	stmt.BindText(1, token)
	stmt.BindText(2, username)
	defer stmt.Finalize()

	_, err = stmt.Step()
	if err != nil {
		return "", fmt.Errorf("failed to store access token: %w", err)
	}

	return token, nil
}

// VerifyAccessToken checks if an access token is valid and its user has the required role
func VerifyAccessToken(token, role string) (bool, error) {
	if token == "" {
		return false, fmt.Errorf("access token is required")
	}

	conn, cleanup, err := getConnection()
	if err != nil {
		return false, err
	}
	defer cleanup()

	stmt := conn.Prep(`
		SELECT created_at
		FROM access_tokens at
		JOIN user_roles ur ON at.user_id = ur.user_id
		WHERE at.token = ? AND ur.role = ?
	`)
	stmt.BindText(1, token)
	stmt.BindText(2, role)
	defer stmt.Finalize()

	if hasRow, err := stmt.Step(); err != nil {
		return false, fmt.Errorf("failed to verify access token: %w", err)
	} else if !hasRow {
		return false, fmt.Errorf("invalid access token")
	}

	// Check that token is not too old
	createdAt, err := parseTimestamp(stmt.ColumnText(0))
	if err != nil {
		return false, fmt.Errorf("failed to parse timestamp: %w", err)
	}
	if createdAt.Add(TokenLifespan).Before(time.Now()) {
		return false, LogoutUser(token)
	}

	return true, nil
}

// LogoutUser removes the access token from the database
func LogoutUser(token string) error {
	if token == "" {
		return fmt.Errorf("access token is required")
	}

	conn, cleanup, err := getConnection()
	if err != nil {
		return err
	}
	defer cleanup()

	stmt := conn.Prep(`DELETE FROM access_tokens WHERE token = ?`)
	stmt.BindText(1, token)
	defer stmt.Finalize()

	_, err = stmt.Step()
	if err != nil {
		return fmt.Errorf("failed to logout user: %w", err)
	}

	return nil
}

// RevokeAllAccessTokens removes all access tokens for a specific user
func RevokeAllAccessTokens(username string) error {
	conn, cleanup, err := getConnection()
	if err != nil {
		return err
	}
	defer cleanup()

	stmt := conn.Prep(`
		DELETE FROM access_tokens
		WHERE user_id = (SELECT user_id FROM users WHERE username = ?)
	`)
	stmt.BindText(1, username)
	defer stmt.Finalize()

	_, err = stmt.Step()
	if err != nil {
		return fmt.Errorf("failed to revoke access tokens: %w", err)
	}

	return nil
}
