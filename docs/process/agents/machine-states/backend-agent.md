# Machine State: Dev Backend Agent

## States

```
IDLE → READING → PLANNING → CODING → TESTING → REVIEWING → DONE
         ↑                                ↓
         └──────── BLOCKED ───────────────┘
```

| State | Description | Entry Condition | Exit Condition |
|---|---|---|---|
| **IDLE** | No active task | PR merged / task blocked / sprint start | Issue assigned with `status:ready` |
| **READING** | Reading issue, ADRs, existing code | Issue picked | Understanding complete |
| **PLANNING** | Designing approach, identifying files | Reading done | Implementation plan clear |
| **CODING** | Writing Go code, migrations | Plan approved (or self-approved) | Code compiles, logic complete |
| **TESTING** | Writing/running tests, lint | Code complete | `go test ./...` + `golangci-lint run` green |
| **REVIEWING** | PR open, addressing review feedback | PR opened | PR approved, CI green |
| **BLOCKED** | Waiting on dependency or clarification | Blocker identified | Blocker resolved |
| **DONE** | PR merged, issue closed | PR squash-merged to main | — |

## Transitions

```
IDLE       → READING    : pick issue with area:backend + status:ready
READING    → PLANNING   : issue fully understood
PLANNING   → CODING     : branch created (feat/<issue>-<desc>)
CODING     → TESTING    : go build ./... passes
TESTING    → CODING     : test failures found
TESTING    → REVIEWING  : go test ./... + golangci-lint clean; PR opened
REVIEWING  → CODING     : review requests changes
REVIEWING  → DONE       : PR approved + CI green → squash merged
ANY        → BLOCKED    : dependency missing or question unanswered
BLOCKED    → READING    : blocker resolved
```

## Context Variables

| Variable | Type | Description |
|---|---|---|
| `current_issue` | string | Active GitHub issue number |
| `branch` | string | Active git branch name |
| `files_modified` | list | Files changed in current task |
| `test_coverage` | float | Current coverage percentage |
| `blockers` | list | Active blockers |

## Output per State

| State | Output |
|---|---|
| READING | Notes on AC and technical approach |
| PLANNING | File list + implementation sketch |
| CODING | `.go` files in `backend/` |
| TESTING | `*_test.go` files, coverage report |
| REVIEWING | PR description, review responses |
| DONE | Merged PR, closed issue |
