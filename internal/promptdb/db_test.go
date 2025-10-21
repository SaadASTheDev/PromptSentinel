package promptdb

import (
	"context"
	"database/sql"
	"errors"
	"testing"
)

func TestConfigConnString(t *testing.T) {
	cfg := Config{
		Host:     "localhost",
		User:     "prompt",
		Database: "promptdb",
		Password: "secret",
	}

	conn, err := cfg.ConnString()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "host=localhost port=5432 user=prompt dbname=promptdb sslmode=require password=secret"
	if conn != expected {
		t.Fatalf("expected connection string %q, got %q", expected, conn)
	}
}

func TestInsertAPIKey(t *testing.T) {
	stub := &stubExec{}
	record := APIKeyRecord{Prefix: "abc123", Hash: "hash", OwnerID: "owner-1"}

	if err := InsertAPIKey(context.Background(), stub, record); err != nil {
		t.Fatalf("unexpected error inserting api key: %v", err)
	}

	if stub.query != insertAPIKeyQuery {
		t.Fatalf("expected query %q, got %q", insertAPIKeyQuery, stub.query)
	}
	if len(stub.args) != 3 || stub.args[0] != "abc123" || stub.args[1] != "hash" || stub.args[2] != "owner-1" {
		t.Fatalf("unexpected arguments: %#v", stub.args)
	}
}

func TestInsertAPIKeyValidation(t *testing.T) {
	stub := &stubExec{}
	if err := InsertAPIKey(context.Background(), stub, APIKeyRecord{}); err == nil {
		t.Fatal("expected validation error for empty record")
	}
	if stub.query != "" {
		t.Fatalf("expected no query to run, got %q", stub.query)
	}
}

func TestListAPIKeyOwners(t *testing.T) {
	owners, err := ListAPIKeyOwners(context.Background(), func(ctx context.Context, query string, args ...any) (RowIterator, error) {
		if query != listAPIKeyOwnersQuery {
			t.Fatalf("expected query %q, got %q", listAPIKeyOwnersQuery, query)
		}
		return &fakeRows{values: []string{"owner-1", "owner-2"}}, nil
	})
	if err != nil {
		t.Fatalf("unexpected error listing owners: %v", err)
	}

	if len(owners) != 2 || owners[0] != "owner-1" || owners[1] != "owner-2" {
		t.Fatalf("unexpected owners slice: %#v", owners)
	}
}

func TestListAPIKeyOwnersQueryError(t *testing.T) {
	queryErr := errors.New("boom")
	_, err := ListAPIKeyOwners(context.Background(), func(ctx context.Context, query string, args ...any) (RowIterator, error) {
		return nil, queryErr
	})
	if !errors.Is(err, queryErr) {
		t.Fatalf("expected wrapped error, got %v", err)
	}
}

type stubExec struct {
	query string
	args  []any
	err   error
}

type stubResult struct{}

func (stubResult) LastInsertId() (int64, error) { return 0, errors.New("not supported") }
func (stubResult) RowsAffected() (int64, error) { return 1, nil }

func (s *stubExec) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	s.query = query
	s.args = append([]any(nil), args...)
	if s.err != nil {
		return stubResult{}, s.err
	}
	return stubResult{}, nil
}

type fakeRows struct {
	values []string
	index  int
}

func (r *fakeRows) Close() error { return nil }
func (r *fakeRows) Err() error   { return nil }

func (r *fakeRows) Next() bool {
	if r.index >= len(r.values) {
		return false
	}
	r.index++
	return true
}

func (r *fakeRows) Scan(dest ...any) error {
	if len(dest) != 1 {
		return errors.New("expected single destination")
	}
	ptr, ok := dest[0].(*string)
	if !ok {
		return errors.New("destination must be *string")
	}
	if r.index == 0 || r.index > len(r.values) {
		return errors.New("iterator not advanced")
	}
	*ptr = r.values[r.index-1]
	return nil
}
