package auth

import (
	"fmt"
	"time"
)

type User struct {
	Username  string
	CreatedAt time.Time
}

// CreateUser creates a new user with the given username and password
func CreateUser(username, password string) error {
	conn, cleanup, err := getConnection()
	if err != nil {
		return err
	}
	defer cleanup()

	// Generate salt and hash password
	salt, err := GenerateSalt()
	if err != nil {
		return fmt.Errorf("failed to generate salt: %w", err)
	}

	passwordHash, err := HashPassword(password, salt)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// Insert user into database
	stmt := conn.Prep(`
		INSERT INTO users (username, password_hash, salt)
		VALUES (?, ?, ?)
	`)
	stmt.BindText(1, username)
	stmt.BindText(2, passwordHash)
	stmt.BindText(3, salt)

	_, err = stmt.Step()
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

// VerifyUserCredentials checks if the username and password are correct
func VerifyUserCredentials(username, password string) error {
	conn, cleanup, err := getConnection()
	if err != nil {
		return err
	}
	defer cleanup()

	stmt := conn.Prep(`
		SELECT password_hash, salt
		FROM users
		WHERE username = ?
	`)
	stmt.BindText(1, username)

	if hasRow, err := stmt.Step(); err != nil {
		return fmt.Errorf("failed to query user: %w", err)
	} else if !hasRow {
		return fmt.Errorf("invalid credentials")
	}

	passwordHash := stmt.ColumnText(0)
	salt := stmt.ColumnText(1)

	// Verify password
	isValid, err := VerifyPassword(password, passwordHash, salt)
	if err != nil {
		return fmt.Errorf("failed to verify password: %w", err)
	}
	if !isValid {
		return fmt.Errorf("invalid credentials")
	}

	return nil
}

// DeleteUser removes a user from the database
func DeleteUser(username string) error {
	conn, cleanup, err := getConnection()
	if err != nil {
		return err
	}
	defer cleanup()

	stmt := conn.Prep(`DELETE FROM users WHERE username = ?`)
	stmt.BindText(1, username)

	_, err = stmt.Step()
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

// ListUsers returns all users in the system
func ListUsers(dbPath string) ([]*User, error) {
	conn, cleanup, err := getConnection()
	if err != nil {
		return nil, err
	}
	defer cleanup()

	stmt := conn.Prep(`
		SELECT username, created_at
		FROM users
		ORDER BY username
	`)

	var users []*User
	for {
		if hasRow, err := stmt.Step(); err != nil {
			return nil, fmt.Errorf("failed to query users: %w", err)
		} else if !hasRow {
			break
		}

		createdAt, err := parseTimestamp(stmt.ColumnText(1))
		if err != nil {
			return nil, fmt.Errorf("failed to parse timestamp: %w", err)
		}

		users = append(users, &User{
			Username:  stmt.ColumnText(0),
			CreatedAt: createdAt,
		})
	}

	return users, nil
}
