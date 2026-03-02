# Machine State: Dev Mobile Agent

## States

```
IDLE → READING → PLANNING → CODING → TESTING → REVIEWING → DONE
         ↑                                ↓
         └──────── BLOCKED ───────────────┘
```

| State | Description | Entry Condition | Exit Condition |
|---|---|---|---|
| **IDLE** | No active task | PR merged / task blocked / sprint start | Issue assigned with `status:ready` |
| **READING** | Reading issue, API contracts, existing screens | Issue picked | Requirements clear |
| **PLANNING** | Designing screen flow, navigation, data | Reading done | Screen plan ready |
| **CODING** | Writing React Native / TypeScript code | Plan ready | TypeScript compiles, screen renders |
| **TESTING** | Running typecheck, lint, Jest | Code complete | `tsc --noEmit` + `eslint` + `jest` green |
| **REVIEWING** | PR open, addressing review feedback | PR opened | PR approved, CI green |
| **BLOCKED** | Waiting on API endpoint, native module, or design | Blocker identified | Blocker resolved |
| **DONE** | PR merged, issue closed | PR squash-merged to main | — |

## Transitions

```
IDLE       → READING    : pick issue with area:mobile + status:ready
READING    → PLANNING   : screen spec and API contract understood
PLANNING   → CODING     : branch created (feat/<issue>-<desc>)
CODING     → TESTING    : TypeScript compiles with no errors
TESTING    → CODING     : lint errors or test failures found
TESTING    → REVIEWING  : tsc + eslint + jest all green; PR opened
REVIEWING  → CODING     : review requests changes
REVIEWING  → DONE       : PR approved + CI green → squash merged
ANY        → BLOCKED    : API endpoint not ready, native issue, design unclear
BLOCKED    → READING    : blocker resolved
```

## Context Variables

| Variable | Type | Description |
|---|---|---|
| `current_issue` | string | Active GitHub issue number |
| `branch` | string | Active git branch name |
| `platform_target` | list | `["ios", "android"]` |
| `screens_added` | list | New screens created |
| `blockers` | list | Active blockers |

## Output per State

| State | Output |
|---|---|
| READING | Screen list and navigation flow |
| PLANNING | Screen tree sketch, API calls identified |
| CODING | `.tsx`/`.ts` files in `mobile/` |
| TESTING | TypeScript clean, ESLint clean, Jest output |
| REVIEWING | PR description, simulator screenshots, review responses |
| DONE | Merged PR, closed issue |
