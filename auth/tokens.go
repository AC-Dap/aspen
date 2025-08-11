package auth

import (
	"fmt"
)

// CreateToken creates an access token for the given user
func CreateToken(username string) error {
	conn, cleanup, err := getConnection()
	if err != nil {
		return err
	}
	defer cleanup()

	// Generate access token
	token, err := GenerateAccessToken()
	if err != nil {
		return fmt.Errorf("failed to generate access token: %w", err)
	}

	stmt := conn.Prep(`
		INSERT INTO access_tokens (user_id, token)
		SELECT user_id
		FROM users
		WHERE username = ?
		VALUES (user_id, ?)
	`)
	stmt.BindText(1, username)
	stmt.BindText(2, token)

	_, err = stmt.Step()
	if err != nil {
		return fmt.Errorf("failed to store access token: %w", err)
	}

	return nil
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
		SELECT 1
		FROM access_tokens at
		JOIN roles r ON at.user_id = r.user_id
		WHERE at.token = ? AND r.role = ?
	`)
	stmt.BindText(1, token)
	stmt.BindText(2, role)

	if hasRow, err := stmt.Step(); err != nil {
		return false, fmt.Errorf("failed to verify access token: %w", err)
	} else if !hasRow {
		return false, fmt.Errorf("invalid access token")
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

	_, err = stmt.Step()
	if err != nil {
		return fmt.Errorf("failed to revoke access tokens: %w", err)
	}

	return nil
}
