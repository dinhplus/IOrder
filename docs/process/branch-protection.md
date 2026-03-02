# Branch Protection Rules — IOrder

> Guide to configure branch protection for the `main` branch.
> Navigate to: **Settings → Branches → Branch protection rules → Add rule**

## Rule for `main`

### Branch name pattern
```
main
```

### Settings

| Setting | Value | Notes |
|---|---|---|
| **Require a pull request before merging** | ✅ Enabled | |
| — Required approvals | **1** | Minimum 1 reviewer |
| — Dismiss stale reviews | ✅ | When new commits are pushed → stale approvals are dismissed |
| — Require review from CODEOWNERS | ✅ | CODEOWNERS must approve |
| **Require status checks to pass** | ✅ Enabled | |
| — Require branches to be up to date | ✅ | Branch must merge latest main |
| — Status checks required: | See list below | |
| **Require conversation resolution** | ✅ | All comments must be resolved |
| **Require signed commits** | ❌ Optional | |
| **Require linear history** | ✅ | Squash merge enforced |
| **Include administrators** | ✅ | Admins must also follow the rules |
| **Allow force pushes** | ❌ Never | |
| **Allow deletions** | ❌ Never | |

### Required Status Checks

Add the following status checks (names match workflow job names):

**Backend:**
- `Lint` (from ci-backend.yml)
- `Test` (from ci-backend.yml)
- `Vulnerability Check` (from ci-backend.yml)
- `Build` (from ci-backend.yml)

**Frontend:**
- `Lint & Typecheck` (from ci-frontend.yml)
- `Build` (from ci-frontend.yml)

**Mobile:**
- `Lint & Typecheck` (from ci-mobile.yml)

> **Note:** Status checks are only required when files in the corresponding path change (due to `paths:` filters in workflows). GitHub will skip checks for unrelated paths.

## Merge Strategy

**Recommended: Squash and merge**

- Settings → General → Pull Requests:
  - ✅ Allow squash merging (default)
  - ❌ Allow merge commits
  - ❌ Allow rebase merging
  - ✅ Default to PR title for squash commit message

## Setup via GitHub CLI (Optional)

```bash
# Requires gh CLI with admin access
gh api repos/dinhplus/IOrder/branches/main/protection \
  --method PUT \
  --field required_status_checks='{"strict":true,"contexts":["Lint","Test","Build","Vulnerability Check","Lint & Typecheck"]}' \
  --field enforce_admins=true \
  --field required_pull_request_reviews='{"required_approving_review_count":1,"dismiss_stale_reviews":true,"require_code_owner_reviews":true}' \
  --field restrictions=null \
  --field allow_force_pushes=false \
  --field allow_deletions=false
```
