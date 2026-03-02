# Sprint Ceremonies — IOrder

> All ceremonies run asynchronously via GitHub Issues. Formats are standardized so AI agents can process them automatically.

## 1. Sprint Planning (Monday)

### Trigger
PO Agent creates a Sprint Planning issue on Monday at 09:00.

### Template

```markdown
---
name: Sprint Planning
about: Sprint planning for Sprint N
title: "Sprint [N] Planning — [Date Range]"
labels: ceremony:planning
---

## Sprint Goal
> [1-2 sentences describing the sprint objective]

## Capacity
| Agent | Available days | Notes |
|---|---|---|
| Dev Backend | 5 | |
| Dev Frontend | 5 | |
| Dev Mobile | 5 | |
| QA | 5 | |
| DevOps | 5 | |

## Sprint Backlog
| # | Issue | Area | Size | Assignee | Priority |
|---|---|---|---|---|---|
| 1 | #XX — [title] | backend | S | Dev Backend | P1 |
| 2 | #YY — [title] | frontend | M | Dev Frontend | P1 |
| ... | | | | | |

## Dependencies
- #XX blocks #YY (API contract needed first)

## Risks
- [List any known risks]

## Checklist
- [ ] Sprint Goal agreed
- [ ] All issues meet DoR
- [ ] Dependencies identified
- [ ] Capacity confirmed
- [ ] Board updated (issues in Ready)
```

### Participants
- **PO:** Facilitates, presents Sprint Goal
- **Tech Lead:** Validates technical feasibility
- **All Dev agents:** Confirm capacity, raise concerns
- **QA:** Confirm test approach

---

## 2. Daily Standup (Async — Daily by 10:00)

### Trigger
Auto-created issue or comment thread daily.

### Template

```markdown
## Daily Standup — Sprint [N] — Day [X] — [Date]

### Dev Backend
- **Yesterday:** [completed]
- **Today:** [planned]
- **Blockers:** [none / details]

### Dev Frontend
- **Yesterday:** [completed]
- **Today:** [planned]
- **Blockers:** [none / details]

### Dev Mobile
- **Yesterday:** [completed]
- **Today:** [planned]
- **Blockers:** [none / details]

### QA
- **Yesterday:** [completed]
- **Today:** [planned]
- **Blockers:** [none / details]

### DevOps
- **Yesterday:** [completed]
- **Today:** [planned]
- **Blockers:** [none / details]
```

### Rules
- Each agent posts their update before 10:00.
- If there is a blocker → tag Tech Lead or PO.
- Keep it brief — max 3 bullets per section.

---

## 3. Backlog Refinement (Wednesday)

### Trigger
PO Agent creates a Refinement issue mid-sprint.

### Template

```markdown
## Backlog Refinement — Sprint [N] — [Date]

### Issues Refined
| # | Issue | Status | Notes |
|---|---|---|---|
| 1 | #XX | Ready ✅ | AC clarified |
| 2 | #YY | Needs work 🔧 | Missing AC for error case |
| 3 | #ZZ | Split → #AA + #BB | Too large |

### New Issues Created
- #AA — [split from #ZZ]
- #BB — [split from #ZZ]

### Removed from Backlog
- #WW — [reason: out of scope]

### Action Items
- [ ] PO: Clarify AC for #YY
- [ ] Tech Lead: Write API contract for #XX
```

### Rules
- Refine the top 10 backlog items.
- Issues > 1 day → split.
- Ensure next sprint has ≥ 8 Ready issues.

---

## 4. Sprint Review (Friday 14:00)

### Trigger
PO Agent creates a Sprint Review issue on Friday.

### Template

```markdown
## Sprint Review — Sprint [N] — [Date]

### Sprint Goal
> [Repeat sprint goal]

### Goal Achieved: ✅ Yes / ❌ No / 🟡 Partial

### Completed Items
| # | Issue | Area | Demo |
|---|---|---|---|
| 1 | #XX — [title] | backend | [screenshot/link] |
| 2 | #YY — [title] | frontend | [screenshot/link] |

### Not Completed
| # | Issue | Reason | Action |
|---|---|---|---|
| 1 | #ZZ — [title] | Blocked by API | Carry to Sprint N+1 |

### Metrics
- Velocity: X story points / Y issues completed
- Bugs found: Z
- CI pass rate: XX%

### Feedback & Insights
- [Key learnings]
```

---

## 5. Sprint Retrospective (Friday 15:00)

### Trigger
Automatically after Sprint Review.

### Template

```markdown
## Sprint Retrospective — Sprint [N] — [Date]

### 🟢 What Went Well
- [item 1]
- [item 2]

### 🔴 What Didn't Go Well
- [item 1]
- [item 2]

### 💡 Action Items for Next Sprint
| # | Action | Owner | Due |
|---|---|---|---|
| 1 | [action] | [agent] | Sprint N+1 |
| 2 | [action] | [agent] | Sprint N+1 |

### Process Improvements
- [ ] [Improvement to try next sprint]
```

### Rules
- Each agent must post at least 1 item per section.
- Action items must have an owner and deadline.
- Max 3 action items per sprint (focus).
