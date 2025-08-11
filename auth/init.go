package auth

import (
	"zombiezen.com/go/sqlite"
	"zombiezen.com/go/sqlite/sqlitex"
)

var dbpool *sqlitex.Pool

// openDatabase opens a connection to the SQLite database at the specified path,
// storing the connection in `dbpool` if successful.
func openDatabase(db_path string) error {
	pool, err := sqlitex.NewPool(db_path, sqlitex.PoolOptions{
		Flags: sqlite.OpenCreate | sqlite.OpenReadWrite,
	})
	if err != nil {
		return err
	}

	dbpool = pool
	return nil
}

// Initializes the sqlite database at the specified path.
// Creates the necessary tables and indexes if they do not exist.
func Initialize(db_path string) error {
	err := openDatabase(db_path)
	if err != nil {
		return err
	}

	// Get a connection to the SQLite database.
	conn, cleanup, err := getConnection()
	if err != nil {
		return err
	}
	defer cleanup()

	// Create users table
	err = sqlitex.ExecScript(conn, `
		CREATE TABLE IF NOT EXISTS users (
			user_id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT UNIQUE NOT NULL,
			password_hash TEXT NOT NULL,
			salt TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);
	`)
	if err != nil {
		return err
	}

	// Create user_roles junction table
	err = sqlitex.ExecScript(conn, `
		CREATE TABLE IF NOT EXISTS user_roles (
			user_id INTEGER NOT NULL,
			role TEXT NOT NULL,
			PRIMARY KEY (user_id, role),
			FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE
		);
	`)
	if err != nil {
		return err
	}

	// Create access_tokens table
	err = sqlitex.ExecScript(conn, `
		CREATE TABLE IF NOT EXISTS access_tokens (
			token_id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			token TEXT UNIQUE NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE
		);
	`)
	if err != nil {
		return err
	}

	// Create OAuth2Tokens table
	err = sqlitex.ExecScript(conn, `
		CREATE TABLE IF NOT EXISTS oauth2_tokens (
			token_id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			token TEXT UNIQUE NOT NULL,
			role TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE
		);
	`)
	if err != nil {
		return err
	}

	// Create AccessTokenNonces table
	err = sqlitex.ExecScript(conn, `
		CREATE TABLE IF NOT EXISTS access_token_nonces (
			token_id INTEGER NOT NULL,
			nonce INTEGER NOT NULL,
			PRIMARY KEY (token_id, nonce),
			FOREIGN KEY (token_id) REFERENCES access_tokens(token_id) ON DELETE CASCADE
		);
	`)
	if err != nil {
		return err
	}

	// Create OAuth2TokenNonces table
	err = sqlitex.ExecScript(conn, `
		CREATE TABLE IF NOT EXISTS oauth2_token_nonces (
			token_id INTEGER NOT NULL,
			nonce INTEGER NOT NULL,
			PRIMARY KEY (token_id, nonce),
			FOREIGN KEY (token_id) REFERENCES oauth2_tokens(token_id) ON DELETE CASCADE
		);
	`)
	if err != nil {
		return err
	}

	// Create indexes for better query performance
	err = sqlitex.ExecScript(conn, `
		CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);
		CREATE INDEX IF NOT EXISTS idx_user_roles_user_id ON user_roles(user_id);
		CREATE INDEX IF NOT EXISTS idx_user_roles_role ON user_roles(role);
		CREATE INDEX IF NOT EXISTS idx_access_tokens_token ON access_tokens(token);
		CREATE INDEX IF NOT EXISTS idx_oauth2_tokens_user_id ON oauth2_tokens(user_id);
		CREATE INDEX IF NOT EXISTS idx_oauth2_tokens_token ON oauth2_tokens(token);
		CREATE INDEX IF NOT EXISTS idx_oauth2_tokens_role ON oauth2_tokens(role);
		CREATE INDEX IF NOT EXISTS idx_access_token_nonces_token_id ON access_token_nonces(token_id);
		CREATE INDEX IF NOT EXISTS idx_access_token_nonces_nonce ON access_token_nonces(nonce);
		CREATE INDEX IF NOT EXISTS idx_oauth2_token_nonces_token_id ON oauth2_token_nonces(token_id);
		CREATE INDEX IF NOT EXISTS idx_oauth2_token_nonces_nonce ON oauth2_token_nonces(nonce);
	`)
	if err != nil {
		return err
	}

	return nil
}

func Close() error {
	if dbpool != nil {
		return dbpool.Close()
	}
	return nil
}
