# Unified Machine States — IOrder Agentic Scrum

All Dev agents (Backend, Frontend, Mobile) share the same top-level state machine.
Area-specific details are in the individual state files.

## Canonical States

```
┌─────────────────────────────────────────────────────────────────────┐
│                      UNIFIED AGENT STATE MACHINE                    │
│                                                                     │
│   ┌──────┐   issue    ┌─────────┐  understood  ┌──────────┐        │
│   │ IDLE │──────────►│ READING │─────────────►│ PLANNING │        │
│   └──────┘  assigned  └─────────┘              └────┬─────┘        │
│      ▲                                              │ plan ready    │
│      │                                              ▼               │
│      │                                         ┌────────┐           │
│      │                                         │ CODING │           │
│      │                                         └───┬────┘           │
│      │                                             │ compiles       │
│      │                                             ▼               │
│      │                                        ┌─────────┐          │
│      │           ┌───────────────────────────│ TESTING │          │
│      │           │ failures found             └────┬────┘          │
│      │           └──────────────────────────────► │ all green      │
│      │                                             ▼               │
│      │                                       ┌──────────┐          │
│      │           ┌──────────────────────────│ REVIEWING│          │
│      │           │ review requests changes   └────┬─────┘          │
│      │           └──────────────────────────────► │ approved       │
│      │                                             ▼               │
│      └────────────────────────────────────── ┌──────┐             │
│                       PR merged               │ DONE │             │
│                                               └──────┘             │
│                                                                     │
│   ┌─────────┐  ◄── any state can transition to BLOCKED             │
│   │ BLOCKED │  ──► READING when blocker is resolved                │
│   └─────────┘                                                      │
└─────────────────────────────────────────────────────────────────────┘
```

## State Definitions (Common)

| State | All Agents |
|---|---|
| **IDLE** | Waiting for issue. Read sprint board. |
| **READING** | Read issue AC, related ADRs, existing code in the area. |
| **PLANNING** | Design the solution. Identify files, dependencies, test cases. |
| **CODING** | Implement the solution. Follow conventions for the area. |
| **TESTING** | Run all checks for the area (see table below). |
| **REVIEWING** | Open PR, respond to feedback, keep CI green. |
| **BLOCKED** | Document blocker on issue, add `status:blocked`, wait. |
| **DONE** | PR merged to main, issue closed. |

## Area-Specific Checks in TESTING State

| Area | Commands |
|---|---|
| **Backend** | `go build ./...` · `go test ./...` · `golangci-lint run` · `govulncheck ./...` |
| **Frontend** | `pnpm run typecheck` · `pnpm run lint` · `pnpm run build` |
| **Mobile** | `pnpm run typecheck` · `pnpm run lint` · `pnpm test` |

## Cross-Agent Handoffs

| From | To | Trigger |
|---|---|---|
| Backend → Frontend | Backend merges API endpoint | Comment "endpoint live" on API contract issue |
| Backend → Mobile | Backend merges API endpoint | Comment "endpoint live" on API contract issue |
| Frontend → QA | PR opened | Move issue to **In Review** |
| Mobile → QA | PR opened | Move issue to **In Review** |
| QA → Dev | Review requests changes | Dev moves to REVIEWING→CODING |
| QA → Release | PR approved + CI green | Merge + move to **Done** |

## Issue Labels per State

| State | Label to Apply |
|---|---|
| IDLE → READING | (remove `status:ready`, add to Sprint milestone) |
| CODING / TESTING | `status:in-progress` (move column on board) |
| REVIEWING | `status:needs-review` |
| BLOCKED | `status:blocked` |
| DONE | (close issue) |

## References

- Per-agent details: `docs/process/agents/machine-states/`
  - [`backend-agent.md`](backend-agent.md)
  - [`frontend-agent.md`](frontend-agent.md)
  - [`mobile-agent.md`](mobile-agent.md)
- Working agreement: `docs/process/working-agreement.md`
- Sprint backlog: `docs/process/sprint-1-backlog.md`
