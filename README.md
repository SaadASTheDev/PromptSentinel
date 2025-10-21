# PromptSentinel

PromptSentinel is a teaching repository that showcases how to secure API keys
and interact with a PostgreSQL database. The project focuses on clarity and
small, well-tested functions so that students can quickly evaluate the structure
of a Go codebase.

## Features

- **API key helpers** – normalize keys, expose safe prefixes, and perform
  constant-time comparisons to reduce the risk of leaking secrets.
- **Database utilities** – build validated PostgreSQL connection strings,
  execute parameterized inserts, and stream query results using familiar Go
  patterns.
- **Thorough unit tests** – lightweight stubs demonstrate how to validate
  database behavior without requiring a live server.

Additional release highlights are documented in [`CHANGELOG.md`](CHANGELOG.md).

## Installation & Setup

1. Ensure Go 1.22 or newer is installed.
2. Clone the repository and install dependencies:

   ```bash
   git clone https://github.com/example/PromptSentinel.git
   cd PromptSentinel
   go mod tidy
   ```

3. (Optional) Provide PostgreSQL credentials via environment variables or a
   configuration file when using the helpers in your own application.

## Usage

The library exposes packages under `internal/` for instructional purposes. The
snippet below shows how to validate a key and prepare a database connection:

```go
ctx := context.Background()
key, err := auth.NewAPIKey(os.Getenv("PROMPTSENTINEL_API_KEY"))
if err != nil {
    log.Fatalf("invalid key: %v", err)
}

cfg := promptdb.Config{
    Host:     "localhost",
    User:     "prompt",
    Password: os.Getenv("PROMPTSENTINEL_DB_PASSWORD"),
    Database: "promptdb",
}

connString, err := cfg.ConnString()
if err != nil {
    log.Fatalf("bad database config: %v", err)
}
log.Printf("connecting with key prefix %s", key.Prefix(4))
log.Printf("connection string: %s", connString)
```

> **Note:** Import a PostgreSQL driver (for example `_ "github.com/lib/pq"`) in
> the main package when executing database operations.

See the inline documentation within each package for more detail.

## Tests

Run all automated checks with:

```bash
go test ./...
```

Each test case is explained in [`docs/test-cases.md`](docs/test-cases.md), which
links expected behavior to the evaluation form.

## Technologies Used

- [Go 1.22](https://go.dev/doc/)
- PostgreSQL (tested with the [`github.com/lib/pq`](https://pkg.go.dev/github.com/lib/pq) driver when running against a live database)
- Standard library testing and context packages for deterministic unit tests

## Repository Structure

```
.
├── CHANGELOG.md          # Feature highlights and release history
├── docs                  # Supplementary documentation (test cases, evaluation mapping)
├── internal/auth         # API key helpers
└── internal/promptdb     # PostgreSQL utilities and accompanying tests
```

## Further Reading

- [`docs/test-cases.md`](docs/test-cases.md) – descriptions of the automated
  tests.
- [`docs/peer-evaluation.md`](docs/peer-evaluation.md) – checklist mapping the
  repository to the CISC 4900 Peer Evaluation Form.
