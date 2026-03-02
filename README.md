# IOrder

> Multi-platform ordering system — Backend (Go/Gin) + Frontend (Next.js) + Mobile (React Native)

## Mono-repo Structure

```
IOrder/
├── backend/          # Go 1.23 + Gin + PostgreSQL
├── frontend/         # Next.js + TypeScript (Node 22)
├── mobile/           # React Native + TypeScript (Node 22)
├── docs/             # Process docs, ADRs, agent instructions
├── scripts/          # Setup & utility scripts
└── .github/          # CI/CD workflows, templates, CODEOWNERS
```

## Quick Start

### Prerequisites

- Go 1.23+
- Node.js 22+
- PostgreSQL 16+
- Docker (for CD)
- GitHub CLI (`gh`) — for label setup

### Setup Labels

```bash
chmod +x scripts/setup-labels.sh
./scripts/setup-labels.sh
```

### Backend

```bash
cd backend
cp .env.example .env
go mod download
go run ./cmd/server
```

### Frontend

```bash
cd frontend
pnpm install  # or npm install
pnpm dev
```

### Mobile

```bash
cd mobile
pnpm install  # or npm install
pnpm start
```

## Process

This project uses an **Agentic Scrum** process with AI agents.

- [Working Agreement](docs/process/working-agreement.md) — DoR/DoD, conventions
- [Sprint Ceremonies](docs/process/sprint-ceremonies.md) — planning, daily, review, retro
- [GitHub Project Setup](docs/process/github-project-setup.md) — board configuration
- [Branch Protection](docs/process/branch-protection.md) — guardrails

### Agent Roles

| Role | Doc |
|---|---|
| Product Owner | [docs/process/agents/po.md](docs/process/agents/po.md) |
| Tech Lead | [docs/process/agents/tech-lead.md](docs/process/agents/tech-lead.md) |
| Dev Backend | [docs/process/agents/dev-backend.md](docs/process/agents/dev-backend.md) |
| Dev Frontend | [docs/process/agents/dev-frontend.md](docs/process/agents/dev-frontend.md) |
| Dev Mobile | [docs/process/agents/dev-mobile.md](docs/process/agents/dev-mobile.md) |
| QA / CI | [docs/process/agents/qa-ci.md](docs/process/agents/qa-ci.md) |
| DevOps / Release | [docs/process/agents/release-devops.md](docs/process/agents/release-devops.md) |

## Architecture Decisions

See [docs/architecture/adr/](docs/architecture/adr/) for all ADRs.

## License

Private
