# Machine State: Dev Frontend Agent

## States

```
IDLE → READING → PLANNING → CODING → TESTING → REVIEWING → DONE
         ↑                                ↓
         └──────── BLOCKED ───────────────┘
```

| State | Description | Entry Condition | Exit Condition |
|---|---|---|---|
| **IDLE** | No active task | PR merged / task blocked / sprint start | Issue assigned with `status:ready` |
| **READING** | Reading issue, API contracts, existing UI | Issue picked | Requirements clear |
| **PLANNING** | Designing component tree, data flow | Reading done | Component plan ready |
| **CODING** | Writing Next.js/React/TypeScript code | Plan ready | TypeScript compiles, UI renders |
| **TESTING** | Running typecheck, lint, build | Code complete | `tsc --noEmit` + `eslint` + `next build` green |
| **REVIEWING** | PR open, addressing review feedback | PR opened | PR approved, CI green |
| **BLOCKED** | Waiting on API endpoint or design spec | Blocker identified | Blocker resolved |
| **DONE** | PR merged, issue closed | PR squash-merged to main | — |

## Transitions

```
IDLE       → READING    : pick issue with area:frontend + status:ready
READING    → PLANNING   : UI spec and API contract understood
PLANNING   → CODING     : branch created (feat/<issue>-<desc>)
CODING     → TESTING    : TypeScript compiles with no errors
TESTING    → CODING     : lint errors or build failures found
TESTING    → REVIEWING  : tsc + eslint + next build all green; PR opened
REVIEWING  → CODING     : review requests changes
REVIEWING  → DONE       : PR approved + CI green → squash merged
ANY        → BLOCKED    : API endpoint not ready or design unclear
BLOCKED    → READING    : blocker resolved (API endpoint merged / design provided)
```

## Context Variables

| Variable | Type | Description |
|---|---|---|
| `current_issue` | string | Active GitHub issue number |
| `branch` | string | Active git branch name |
| `api_base_url` | string | Backend API URL in use |
| `components_added` | list | New components created |
| `blockers` | list | Active blockers (missing API, design, etc.) |

## Output per State

| State | Output |
|---|---|
| READING | Component/page list and data requirements |
| PLANNING | Component tree sketch, API calls identified |
| CODING | `.tsx`/`.ts` files in `frontend/src/` |
| TESTING | TypeScript clean, ESLint clean, build output |
| REVIEWING | PR description, screenshots, review responses |
| DONE | Merged PR, closed issue |
