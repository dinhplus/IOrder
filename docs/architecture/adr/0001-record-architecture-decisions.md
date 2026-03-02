# ADR 0001: Record Architecture Decisions

## Status

Accepted

## Date

2026-03-02

## Context

We need to record significant architecture decisions for the IOrder project (mono-repo: Backend Go+Gin, Frontend Next.js, Mobile React Native) so that:

1. AI agents understand the context when implementing features.
2. Decisions have traceability and can be revisited.
3. Onboarding new agents/developers is faster.

## Decision

Use **Architecture Decision Records (ADR)** following Michael Nygard's format, stored at `docs/architecture/adr/`.

### Format

Each ADR is a single markdown file:

```
docs/architecture/adr/NNNN-<short-title>.md
```

Contents include:
- **Status:** Proposed / Accepted / Deprecated / Superseded by NNNN
- **Date:** YYYY-MM-DD
- **Context:** The problem to be solved
- **Decision:** The chosen decision
- **Consequences:** Outcomes (positive + negative)

### Numbering

Sequential 4-digit: `0001`, `0002`, ...

### When to Write an ADR

- Choosing a new framework/library
- Changing project structure
- API design pattern decisions
- Database schema decisions
- CI/CD pipeline changes
- Cross-platform conventions

### Who Creates ADRs

- **Tech Lead Agent** is the primary author.
- **Dev agents** may propose ADRs via PR.
- **PO** approves ADRs related to product direction.

## Consequences

### Positive

- All decisions are documented and available.
- AI agents have context to implement correctly.
- Easy to onboard, easy to revisit decisions.

### Negative

- Small overhead when writing ADRs.
- Requires discipline to keep them updated.

## Template for New ADRs

```markdown
# ADR NNNN: [Title]

## Status
Proposed

## Date
YYYY-MM-DD

## Context
[Describe the issue/force that requires a decision]

## Decision
[Describe the decision and rationale]

## Consequences
### Positive
- [benefit 1]

### Negative
- [trade-off 1]
```
