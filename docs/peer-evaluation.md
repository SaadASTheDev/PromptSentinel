# CISC 4900 Peer Evaluation Checklist

This checklist documents how the repository satisfies each item of the CISC
4900 Peer Evaluation Form. Use it as a quick reference when reviewing the
project.

| Evaluation Item | Evidence |
|-----------------|----------|
| Project Description | [`README.md`](../README.md) – "Project Description" section explains the goal of PromptSentinel. |
| Feature Description | [`README.md`](../README.md#features) and [`CHANGELOG.md`](../CHANGELOG.md) summarize available features. |
| Installation / Execution Instructions | [`README.md`](../README.md#installation--setup) provides setup steps. |
| Technologies Used | [`README.md`](../README.md#technologies-used) lists the primary stack. |
| Sensible filenames and structure | Directory tree documented in [`README.md`](../README.md#repository-structure); Go packages follow idiomatic lowercase names. |
| Source code includes comments | See exported symbols in [`internal/auth/auth.go`](../internal/auth/auth.go) and [`internal/promptdb/db.go`](../internal/promptdb/db.go). |
| Readable, well-formatted code | All Go files are formatted with `gofmt`, use descriptive identifiers, and include focused helper functions. |
| Test cases provided with expected results | Automated tests live in `*_test.go` files; descriptions and expectations are in [`docs/test-cases.md`](test-cases.md). |
| Each test case described | [`docs/test-cases.md`](test-cases.md) maps each test to its intent and outcome. |
| `.gitignore` present | Repository root contains [`.gitignore`](../.gitignore). |
| Issue tracker activity | Manage via GitHub Issues; update or create items when submitting work. |
| Branching strategy | Use topic branches for new features or fixes; merge via pull requests. |
| Submodules | Not required for this project; confirm "No" on the evaluation form. |
| Pull requests | Submit changes through GitHub pull requests to document review history. |

For items tracked outside the repository (issues, branches, pull requests), the
notes above describe the expected workflow.
