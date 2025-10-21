// Package promptdb centralizes helpers for working with the PromptSentinel
// PostgreSQL database. The functions in this file intentionally focus on small
// and easily testable behaviors so that students can understand each step of
// the data flow when reviewing the repository.
package promptdb

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

const defaultPostgresPort = 5432

// Config describes how to connect to the database. Only the fields that are
// commonly customized are exposed to keep configuration approachable for new
// contributors.
type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
	SSLMode  string
}

// ConnString builds a lib/pq compatible PostgreSQL connection string. The
// method validates that the required fields are present before joining them into
// the final connection string.
func (c Config) ConnString() (string, error) {
	if c.Host == "" {
		return "", errors.New("host is required")
	}
	if c.User == "" {
		return "", errors.New("user is required")
	}
	if c.Database == "" {
		return "", errors.New("database is required")
	}

	port := c.Port
	if port == 0 {
		port = defaultPostgresPort
	}

	sslMode := c.SSLMode
	if sslMode == "" {
		sslMode = "require"
	}

	parts := []string{
		fmt.Sprintf("host=%s", c.Host),
		fmt.Sprintf("port=%d", port),
		fmt.Sprintf("user=%s", c.User),
		fmt.Sprintf("dbname=%s", c.Database),
		fmt.Sprintf("sslmode=%s", sslMode),
	}

	if c.Password != "" {
		parts = append(parts, fmt.Sprintf("password=%s", c.Password))
	}

	return strings.Join(parts, " "), nil
}

// Open creates a *sql.DB and ensures that the database is reachable with the
// provided configuration. Callers must import a PostgreSQL driver (for example
// github.com/lib/pq) in their main package before invoking this helper. Open is
// intentionally optional for tests because most unit tests should rely on
// lightweight stubs instead of a real database connection.
func Open(ctx context.Context, cfg Config) (*sql.DB, error) {
	connString, err := cfg.ConnString()
	if err != nil {
		return nil, err
	}

	db, err := sql.Open("postgres", connString)
	if err != nil {
		return nil, fmt.Errorf("open postgres: %w", err)
	}

	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("ping postgres: %w", err)
	}

	return db, nil
}

// APIKeyRecord describes the information stored in the api_keys table. The
// struct definition doubles as documentation for new contributors.
type APIKeyRecord struct {
	Prefix  string
	Hash    string
	OwnerID string
}

func (r APIKeyRecord) validate() error {
	if strings.TrimSpace(r.Prefix) == "" {
		return errors.New("prefix is required")
	}
	if strings.TrimSpace(r.Hash) == "" {
		return errors.New("hash is required")
	}
	if strings.TrimSpace(r.OwnerID) == "" {
		return errors.New("owner id is required")
	}

	return nil
}

// execContext is satisfied by *sql.DB, *sql.Tx, and lightweight stubs used in
// tests. It keeps InsertAPIKey flexible for different call sites without
// importing additional packages.
type execContext interface {
	ExecContext(context.Context, string, ...any) (sql.Result, error)
}

const insertAPIKeyQuery = `INSERT INTO api_keys (key_prefix, key_hash, owner_id) VALUES ($1, $2, $3)`

// InsertAPIKey writes an APIKeyRecord to the database. The function validates
// the record before executing the INSERT statement so that students understand
// where validation should occur in a Go codebase.
func InsertAPIKey(ctx context.Context, db execContext, record APIKeyRecord) error {
	if err := record.validate(); err != nil {
		return err
	}

	if _, err := db.ExecContext(ctx, insertAPIKeyQuery, record.Prefix, record.Hash, record.OwnerID); err != nil {
		return fmt.Errorf("insert api key: %w", err)
	}

	return nil
}

// queryContext mirrors execContext but for SELECT style queries.
const listAPIKeyOwnersQuery = `SELECT owner_id FROM api_keys ORDER BY owner_id`

// RowIterator is the minimal interface required to loop through query results.
// sql.Rows satisfies this contract, allowing callers to pass database query
// results directly while keeping tests simple to write.
type RowIterator interface {
	Close() error
	Err() error
	Next() bool
	Scan(dest ...any) error
}

// RowQueryFunc performs a query and returns a RowIterator. Callers can pass
// `db.QueryContext` from *sql.DB or provide a lightweight stub in tests.
type RowQueryFunc func(ctx context.Context, query string, args ...any) (RowIterator, error)

// ListAPIKeyOwners reads all API key owners from the database. Returning a
// string slice keeps the example approachable while still demonstrating how to
// loop through RowIterator values.
func ListAPIKeyOwners(ctx context.Context, query RowQueryFunc) ([]string, error) {
	rows, err := query(ctx, listAPIKeyOwnersQuery)
	if err != nil {
		return nil, fmt.Errorf("list api key owners: %w", err)
	}
	defer rows.Close()

	var owners []string
	for rows.Next() {
		var owner string
		if err := rows.Scan(&owner); err != nil {
			return nil, fmt.Errorf("scan owner: %w", err)
		}
		owners = append(owners, owner)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate owners: %w", err)
	}

	return owners, nil
}
