# GitHub Project Setup — IOrder Scrum Board

> Step-by-step guide to create and configure the GitHub Project board for the team.

## 1. Create Project

1. Go to https://github.com/dinhplus/IOrder
2. Click the **Projects** tab → **New project**
3. Select **Board** (Kanban view)
4. Name: **IOrder Sprint Board**
5. Click **Create project**

## 2. Configure Columns (Status)

GitHub creates default columns: Todo, In Progress, Done. Modify as follows:

| Column | Description |
|---|---|
| **Backlog** | Not yet refined, not ready |
| **Ready** | Meets DoR, available for pickup |
| **In Progress** | Currently being implemented |
| **In Review** | PR opened, under review |
| **QA** | QA agent is verifying |
| **Done** | Merged + verified |

### Steps:
1. Click column header → **Rename** → rename according to the table above
2. Click **+ Add column** to add missing columns
3. Arrange order: Backlog → Ready → In Progress → In Review → QA → Done

## 3. Add Custom Fields

### Priority Field
1. Click **+** (Add field) → **Single select**
2. Name: **Priority**
3. Options: `P0` (🔴), `P1` (🟡), `P2` (🟢)

### Area Field
1. Click **+** → **Single select**
2. Name: **Area**
3. Options: `backend`, `frontend`, `mobile`, `infra`

### Size Field
1. Click **+** → **Single select**
2. Name: **Size**
3. Options: `XS` (< 2h), `S` (2-4h), `M` (4-8h)

### Sprint Field
1. Click **+** → **Iteration**
2. Name: **Sprint**
3. Duration: 1 week
4. Start date: Monday

## 4. Views

### Board View (default)
- Group by: **Status** (columns)
- Sort by: **Priority** desc

### Table View
1. Click **+ New view** → **Table**
2. Name: "Sprint Backlog"
3. Show fields: Title, Status, Priority, Area, Size, Sprint, Assignees
4. Filter: Sprint = current

### Area View
1. Click **+ New view** → **Board**
2. Name: "By Area"
3. Group by: **Area**

## 5. Automation Rules

> GitHub Projects (v2) has built-in automations. Configure as follows:

### Auto-add issues
1. Project Settings → **Workflows**
2. Enable: **Item added to project** → Set Status = **Backlog**

### Auto-move on PR
1. Enable: **Pull request merged** → Set Status = **Done**
2. Enable: **Pull request opened** → Set Status = **In Review**

### Custom automation (for advanced use)
Use GitHub Actions + GraphQL API. Example:

```yaml
# .github/workflows/project-automation.yml
name: Project Automation
on:
  issues:
    types: [labeled]
  pull_request:
    types: [opened, closed]

jobs:
  automate:
    runs-on: ubuntu-latest
    steps:
      - name: Move issue to Ready when labeled
        if: github.event.label.name == 'status:ready'
        uses: actions/github-script@v7
        with:
          script: |
            // Use GitHub Projects GraphQL API
            // to move item to Ready column
            console.log('Would move to Ready column');
```

## 6. Labels Setup

Run the following script using GitHub CLI (`gh`):

```bash
#!/bin/bash
REPO="dinhplus/IOrder"

# Type labels
gh label create "type:feature" --color "0E8A16" --description "New feature" -R $REPO
gh label create "type:bug" --color "D93F0B" --description "Bug fix" -R $REPO
gh label create "type:chore" --color "FEF2C0" --description "Maintenance task" -R $REPO

# Priority labels
gh label create "prio:P0" --color "B60205" --description "Critical - fix immediately" -R $REPO
gh label create "prio:P1" --color "FF9500" --description "High priority" -R $REPO
gh label create "prio:P2" --color "0E8A16" --description "Normal priority" -R $REPO

# Area labels
gh label create "area:backend" --color "5319E7" --description "Backend (Go/Gin)" -R $REPO
gh label create "area:frontend" --color "1D76DB" --description "Frontend (Next.js)" -R $REPO
gh label create "area:mobile" --color "D4C5F9" --description "Mobile (React Native)" -R $REPO
gh label create "area:infra" --color "C2E0C6" --description "Infrastructure/CI/CD" -R $REPO

# Status labels
gh label create "status:blocked" --color "E4E669" --description "Blocked by dependency" -R $REPO
gh label create "status:needs-review" --color "FBCA04" --description "Needs review" -R $REPO
gh label create "status:ready" --color "0E8A16" --description "Ready for pickup" -R $REPO

# Size labels
gh label create "size:XS" --color "C5DEF5" --description "Extra small (< 2h)" -R $REPO
gh label create "size:S" --color "BFD4F2" --description "Small (2-4h)" -R $REPO
gh label create "size:M" --color "A2C4E0" --description "Medium (4-8h)" -R $REPO

# Ceremony labels
gh label create "ceremony:planning" --color "D4C5F9" --description "Sprint planning" -R $REPO
gh label create "ceremony:review" --color "D4C5F9" --description "Sprint review" -R $REPO
gh label create "ceremony:retro" --color "D4C5F9" --description "Retrospective" -R $REPO
gh label create "ceremony:daily" --color "D4C5F9" --description "Daily standup" -R $REPO
```

## 7. Link Repo to Project

1. Project Settings → **Manage access**
2. Add repository: `dinhplus/IOrder`
3. Role: **Write** (so workflows can update the project)

## 8. Single Board for All 3 Platforms

**Decision:** 1 shared board.

**Rationale:**
- Mono-repo → a single project board is natural.
- Cross-platform dependencies are easier to track on the same board.
- Filter by `Area` field instead of switching boards.
- PO + Tech Lead see the entire sprint progress in one place.

**Usage:**
- Use the **Area** field to filter: Backend / Frontend / Mobile.
- Default board view: all areas.
- Table view: filter by area when focused work is needed.
