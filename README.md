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

## Business Requirements

See [docs/business/requirements.md](docs/business/requirements.md) for full PO business analysis including user stories, acceptance criteria, and non-functional requirements.

## Architecture

See [docs/architecture/system-architecture.md](docs/architecture/system-architecture.md) for the full system architecture (AWS infrastructure, service design, database schema, API design).

### Architecture Decisions

| ADR | Title |
|---|---|
| [ADR 0001](docs/architecture/adr/0001-record-architecture-decisions.md) | Record Architecture Decisions |
| [ADR 0002](docs/architecture/adr/0002-aws-infrastructure.md) | AWS Infrastructure Architecture |
| [ADR 0003](docs/architecture/adr/0003-order-state-machine.md) | Order State Machine Design |
| [ADR 0004](docs/architecture/adr/0004-payment-integration.md) | Payment Integration Strategy |
| [ADR 0005](docs/architecture/adr/0005-qr-table-ordering.md) | QR Code Table Ordering |

## Deployment

See [docs/architecture/deployment-guide.md](docs/architecture/deployment-guide.md) for instructions on local setup, AWS deployment, and operations.

## License

Private
