package auth

import (
	"fmt"
)

// AssignRoleToUser assigns a role to a user
func AssignRoleToUser(username, role string) error {
	conn, cleanup, err := getConnection()
	if err != nil {
		return err
	}
	defer cleanup()

	// Get user ID and role ID
	stmt := conn.Prep(`
		INSERT INTO user_roles (user_id, role)
		SELECT u.user_id
		FROM users u
		WHERE u.username = ?
	`)
	stmt.BindText(1, username)
	stmt.BindText(2, role)
	defer stmt.Finalize()

	_, err = stmt.Step()
	if err != nil {
		return fmt.Errorf("failed to assign role to user: %w", err)
	}

	if conn.Changes() == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

// RemoveRoleFromUser removes a role from a user
func RemoveRoleFromUser(username, role string) error {
	conn, cleanup, err := getConnection()
	if err != nil {
		return err
	}
	defer cleanup()

	stmt := conn.Prep(`
		DELETE FROM user_roles
		WHERE user_id = (SELECT user_id FROM users WHERE username = ?)
		  AND role = ?
	`)
	stmt.BindText(1, username)
	stmt.BindText(2, role)
	defer stmt.Finalize()

	_, err = stmt.Step()
	if err != nil {
		return fmt.Errorf("failed to remove role from user: %w", err)
	}

	return nil
}

// CheckUserPermissions checks if a user has a specific role
func CheckUserPermissions(username, role string) (bool, error) {
	conn, cleanup, err := getConnection()
	if err != nil {
		return false, err
	}
	defer cleanup()

	stmt := conn.Prep(`
		SELECT 1
		FROM user_roles ur
		JOIN users u ON ur.user_id = u.user_id
		WHERE u.username = ? AND ur.role = ?
	`)
	stmt.BindText(1, username)
	stmt.BindText(2, role)
	defer stmt.Finalize()

	if hasRow, err := stmt.Step(); err != nil {
		return false, fmt.Errorf("failed to check user permissions: %w", err)
	} else if !hasRow {
		return false, nil
	}

	return true, nil
}

// GetUserRoles returns all roles assigned to a user
func GetUserRoles(dbPath, username string) ([]string, error) {
	conn, cleanup, err := getConnection()
	if err != nil {
		return nil, err
	}
	defer cleanup()

	stmt := conn.Prep(`
		SELECT ur.role
		FROM user_roles ur
		JOIN users u ON ur.user_id = u.user_id
		WHERE u.username = ?
		ORDER BY ur.role
	`)
	stmt.BindText(1, username)
	defer stmt.Finalize()

	var roles []string
	for {
		if hasRow, err := stmt.Step(); err != nil {
			return nil, fmt.Errorf("failed to query user roles: %w", err)
		} else if !hasRow {
			break
		}

		roles = append(roles, stmt.ColumnText(0))
	}

	return roles, nil
}

// GetUsersWithRole returns all users that have a specific role
func GetUsersWithRole(role string) ([]*User, error) {
	conn, cleanup, err := getConnection()
	if err != nil {
		return nil, err
	}
	defer cleanup()

	stmt := conn.Prep(`
		SELECT u.username, u.created_at
		FROM user_roles ur
		JOIN users u ON ur.user_id = u.user_id
		WHERE ur.role = ?
		ORDER BY u.username
	`)
	stmt.BindText(1, role)
	defer stmt.Finalize()

	var users []*User
	for {
		if hasRow, err := stmt.Step(); err != nil {
			return nil, fmt.Errorf("failed to query users with role: %w", err)
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
