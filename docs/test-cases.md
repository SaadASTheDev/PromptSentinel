# Test Cases

This document summarizes the automated checks that accompany the repository. It
links each Go unit test to the behavior it validates so that reviewers can map
results to the peer evaluation form requirement: "Test Case(s) are provided with
descriptions."

## API Key Helpers (`internal/auth`)

| Test Name | Description | Expected Result |
|-----------|-------------|-----------------|
| `TestNewAPIKeyTrimAndValidate` | Ensures that whitespace is trimmed and that keys shorter than `MinimumAPIKeyLength` are rejected. | The test passes when valid keys are preserved and short keys return an error. |
| `TestNewAPIKeyTooShort` | Verifies that clearly invalid keys fail fast with an explanatory message. | The test passes when an error is returned. |
| `TestAPIKeyMatches` | Confirms that constant-time comparisons succeed for identical keys and fail for mismatched values. | The test passes when identical keys return `true` and different keys return `false`. |
| `TestAPIKeyPrefix` | Demonstrates how prefixes hide the majority of an API key while remaining configurable. | The test passes when prefixes of varying lengths match expectations. |

## Database Utilities (`internal/promptdb`)

| Test Name | Description | Expected Result |
|-----------|-------------|-----------------|
| `TestConfigConnString` | Builds a connection string from minimal configuration and enforces default values. | The resulting string matches the documented format. |
| `TestInsertAPIKey` | Executes the parameterized insert using a lightweight stub executor. | The insert query runs with the provided arguments. |
| `TestInsertAPIKeyValidation` | Ensures that obviously incomplete records are rejected before reaching the database. | The function returns an error and no SQL statements are executed. |
| `TestListAPIKeyOwners` | Streams rows from a stubbed result set to demonstrate safe iteration. | The function returns the owner IDs in order without errors. |

To rerun all cases locally, execute `go test ./...` from the project root.
