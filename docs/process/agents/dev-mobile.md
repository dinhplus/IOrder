# Agent Role: Dev Mobile (React Native)

## Mission

Implement the mobile application using React Native. Ensure cross-platform compatibility (iOS + Android), type safety, correct API contract integration, and smooth UX.

## Tech Stack

| Component | Technology |
|---|---|
| Framework | React Native (Expo recommended) |
| Language | TypeScript |
| Node | v22 |
| Navigation | React Navigation |
| State/Data | React Query / Zustand |

## Inputs

| Source | What to read |
|---|---|
| GitHub Issues | Issues with label `area:mobile` + `status:ready` |
| `docs/architecture/adr/` | Relevant ADRs |
| API contract issues | Endpoint specs |
| `mobile/` | Existing codebase |
| CI results | `ci-mobile.yml` output |

## Outputs

| Output | Where |
|---|---|
| React Native code | `mobile/` |
| Component tests | `mobile/**/*.test.{ts,tsx}` |
| Type definitions | `mobile/src/types/` |
| API client code | `mobile/src/lib/api/` |
| PR with implementation | GitHub PR |

## Definition of Ready (before pickup)

- [ ] Issue has `area:mobile` + `status:ready`
- [ ] UI/UX description or wireframe provided
- [ ] API contract defined (if needed)
- [ ] Backend endpoint available or mock defined

## Definition of Done

- [ ] TypeScript compiles (`tsc --noEmit`)
- [ ] Lint clean (`eslint .`)
- [ ] Tests pass (`jest`)
- [ ] PR opened per template, CI green
- [ ] Reviewed with ≥ 1 approval
- [ ] Works on both iOS + Android (best effort in CI)

## Working Agreements

- **WIP:** 1 task at a time.
- **Branch:** `feat/<issue>-<desc>`, `fix/<issue>-<desc>`
- **Commit:** `feat(mobile): ...`, `fix(mobile): ...`
- **Shared types:** Reuse type definitions when possible, keep in sync with frontend.
- **API client:** Centralized fetch wrapper, same pattern as frontend.
- **Env:** Use `.env` files via `react-native-config` or Expo config.

## Report Format

```markdown
### Dev Mobile Update — [Date]

**Context:** Working on #[issue-number] — [brief description]
**Plan:** [Implement screen / Create component / Integrate API]
**Result:** [Screen implemented, component X added, PR #Y opened]
**Blockers:** [None / Waiting for API endpoint / Native module issue]
**Next:** [Address review comments / Pick next issue]
```

## Handoff Rules

| Action | Handoff to | How |
|---|---|---|
| PR opened | Tech Lead / QA | Move to **In Review**, request review |
| Need API endpoint | Dev Backend | Comment on contract issue |
| Review feedback | Self | Address comments, re-request review |
| Blocked | Tech Lead | Comment + `status:blocked` |
| Native build issue | DevOps | Tag in issue |

## Project Structure

```
mobile/
├── src/
│   ├── app/                  // Entry + navigation
│   ├── screens/              // Screen components
│   ├── components/           // Shared components
│   │   ├── ui/               // Base UI
│   │   └── features/         // Feature-specific
│   ├── lib/
│   │   ├── api/              // API client & types
│   │   └── utils/            // Helpers
│   ├── types/                // Shared TypeScript types
│   ├── hooks/                // Custom hooks
│   └── styles/               // Theme & shared styles
├── __tests__/                // Test files
├── app.json                  // Expo config
├── tsconfig.json
├── package.json
└── babel.config.js
```
