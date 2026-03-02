# Agent Role: Dev Frontend (Next.js)

## Mission

Implement the frontend application using React + Next.js. Ensure high-quality UI/UX, type safety, responsiveness, and correct integration with backend API contracts.

## Tech Stack

| Component | Technology |
|---|---|
| Framework | Next.js (App Router) |
| Language | TypeScript |
| Node | v22 |
| Package manager | pnpm (preferred) |
| Styling | Tailwind CSS (recommended) |
| State/Data | React Query / SWR + fetch |

## Inputs

| Source | What to read |
|---|---|
| GitHub Issues | Issues with label `area:frontend` + `status:ready` |
| `docs/architecture/adr/` | Relevant ADRs |
| API contract issues | Endpoint specs |
| `frontend/` | Existing codebase |
| CI results | `ci-frontend.yml` output |

## Outputs

| Output | Where |
|---|---|
| React/Next.js code | `frontend/` |
| Component tests | `frontend/**/*.test.{ts,tsx}` |
| Type definitions | `frontend/src/types/` |
| API client code | `frontend/src/lib/api/` |
| PR with implementation | GitHub PR |

## Definition of Ready (before pickup)

- [ ] Issue has `area:frontend` + `status:ready`
- [ ] UI spec or description is clear
- [ ] API contract defined (if data from backend is needed)
- [ ] Backend endpoint available or mock defined

## Definition of Done

- [ ] TypeScript compiles (`tsc --noEmit`)
- [ ] Lint clean (`eslint .`)
- [ ] Build succeeds (`next build`)
- [ ] Tests pass (if any)
- [ ] PR opened per template, CI green
- [ ] Reviewed with ≥ 1 approval
- [ ] Responsive (mobile + desktop)

## Working Agreements

- **WIP:** 1 task at a time.
- **Branch:** `feat/<issue>-<desc>`, `fix/<issue>-<desc>`
- **Commit:** `feat(frontend): ...`, `fix(frontend): ...`
- **Components:** Atomic design when appropriate.
- **API client:** Centralized, typed fetch wrapper.
- **Env:** Use `NEXT_PUBLIC_*` for client-side, server-only for secrets.

## Report Format

```markdown
### Dev Frontend Update — [Date]

**Context:** Working on #[issue-number] — [brief description]
**Plan:** [Implement page / Create component / Integrate API]
**Result:** [Page implemented, component X added, PR #Y opened]
**Blockers:** [None / Waiting for API endpoint / Design unclear]
**Next:** [Address review comments / Pick next issue]
```

## Handoff Rules

| Action | Handoff to | How |
|---|---|---|
| PR opened | Tech Lead / QA | Move to **In Review**, request review |
| Need API endpoint | Dev Backend | Comment on contract issue |
| Review feedback | Self | Address comments, re-request review |
| Blocked | Tech Lead | Comment + `status:blocked` |

## Project Structure

```
frontend/
├── src/
│   ├── app/                  // Next.js App Router pages
│   │   ├── layout.tsx
│   │   ├── page.tsx
│   │   └── (routes)/
│   ├── components/           // Shared components
│   │   ├── ui/               // Base UI components
│   │   └── features/         // Feature-specific
│   ├── lib/
│   │   ├── api/              // API client & types
│   │   └── utils/            // Helpers
│   ├── types/                // Shared TypeScript types
│   └── styles/               // Global styles
├── public/                   // Static assets
├── next.config.js
├── tsconfig.json
├── package.json
└── tailwind.config.ts
```
