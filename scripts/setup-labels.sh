#!/bin/bash
# Setup labels for IOrder repo
# Usage: ./scripts/setup-labels.sh

set -euo pipefail

REPO="dinhplus/IOrder"

echo "Setting up labels for $REPO..."

# Delete default labels that we don't use
for label in "bug" "documentation" "duplicate" "enhancement" "good first issue" "help wanted" "invalid" "question" "wontfix"; do
  gh label delete "$label" -R "$REPO" --yes 2>/dev/null || true
done

# Type labels
gh label create "type:feature" --color "0E8A16" --description "New feature" -R "$REPO" --force
gh label create "type:bug" --color "D93F0B" --description "Bug fix" -R "$REPO" --force
gh label create "type:chore" --color "FEF2C0" --description "Maintenance task" -R "$REPO" --force

# Priority labels
gh label create "prio:P0" --color "B60205" --description "Critical — fix immediately" -R "$REPO" --force
gh label create "prio:P1" --color "FF9500" --description "High priority" -R "$REPO" --force
gh label create "prio:P2" --color "0E8A16" --description "Normal priority" -R "$REPO" --force

# Area labels
gh label create "area:backend" --color "5319E7" --description "Backend (Go/Gin)" -R "$REPO" --force
gh label create "area:frontend" --color "1D76DB" --description "Frontend (Next.js)" -R "$REPO" --force
gh label create "area:mobile" --color "D4C5F9" --description "Mobile (React Native)" -R "$REPO" --force
gh label create "area:infra" --color "C2E0C6" --description "Infrastructure/CI/CD" -R "$REPO" --force

# Status labels
gh label create "status:blocked" --color "E4E669" --description "Blocked by dependency" -R "$REPO" --force
gh label create "status:needs-review" --color "FBCA04" --description "Needs code review" -R "$REPO" --force
gh label create "status:ready" --color "0E8A16" --description "Ready for pickup (meets DoR)" -R "$REPO" --force

# Size labels
gh label create "size:XS" --color "C5DEF5" --description "Extra small (< 2h)" -R "$REPO" --force
gh label create "size:S" --color "BFD4F2" --description "Small (2-4h)" -R "$REPO" --force
gh label create "size:M" --color "A2C4E0" --description "Medium (4-8h)" -R "$REPO" --force

# Ceremony labels
gh label create "ceremony:planning" --color "D4C5F9" --description "Sprint planning" -R "$REPO" --force
gh label create "ceremony:review" --color "D4C5F9" --description "Sprint review" -R "$REPO" --force
gh label create "ceremony:retro" --color "D4C5F9" --description "Sprint retrospective" -R "$REPO" --force
gh label create "ceremony:daily" --color "D4C5F9" --description "Daily standup" -R "$REPO" --force

echo "✅ Labels setup complete!"
