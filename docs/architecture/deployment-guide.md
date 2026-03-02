# IOrder — Developer Deployment Guide

> **Author:** Tech Lead
> **Date:** 2026-03-02
> **Version:** 1.0

---

## Overview

Hướng dẫn này mô tả cách thiết lập môi trường local, deploy lên AWS, và vận hành hệ thống IOrder.

---

## 1. Prerequisites

### 1.1 Local Development Tools

```bash
# Required versions
go 1.24+
node 22+
pnpm 9+
docker 24+
docker-compose 2.24+
aws-cli 2.15+
make
```

Cài đặt:
```bash
# Go (Linux/macOS)
wget https://go.dev/dl/go1.24.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.24.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin

# Node via nvm
curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.7/install.sh | bash
nvm install 22
nvm use 22

# pnpm
npm install -g pnpm@9

# AWS CLI
curl "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" -o "awscliv2.zip"
unzip awscliv2.zip && sudo ./aws/install

# golangci-lint
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh \
  | sh -s -- -b $(go env GOPATH)/bin v1.61.0

# golang-migrate CLI
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

### 1.2 AWS Setup

```bash
# Configure AWS credentials
aws configure
# AWS Access Key ID: <your-key>
# AWS Secret Access Key: <your-secret>
# Default region: ap-southeast-1
# Default output format: json

# Verify
aws sts get-caller-identity
```

### 1.3 Infrastructure Tools (for deploying infra)

```bash
# AWS CDK
npm install -g aws-cdk@2

# Verify
cdk --version
```

---

## 2. Local Development Setup

### 2.1 Clone & Setup

```bash
git clone https://github.com/dinhplus/IOrder.git
cd IOrder
```

### 2.2 Backend Setup

```bash
cd backend

# Copy environment config
cp .env.example .env
# Edit .env with your local PostgreSQL credentials

# Download dependencies
go mod download

# Start PostgreSQL and Redis via Docker
docker-compose up -d postgres redis

# Run database migrations
make migrate-up

# Start development server
make run
# Server starts at http://localhost:8080

# Verify health check
curl http://localhost:8080/health
```

**`backend/.env.example` variables:**

| Variable | Example | Description |
|---|---|---|
| `PORT` | `8080` | HTTP server port |
| `DATABASE_URL` | `postgres://iorder:password@localhost:5432/iorder?sslmode=disable` | PostgreSQL connection string |
| `REDIS_URL` | `redis://localhost:6379` | Redis connection string |
| `LOG_LEVEL` | `debug` | Logging level (debug/info/warn/error) |
| `ENV` | `development` | Environment name |
| `JWT_SECRET` | `local-dev-secret-min-32-chars` | JWT signing secret (dev only) |
| `CORS_ORIGINS` | `http://localhost:3000,http://localhost:8081` | Allowed CORS origins |
| `TABLE_QR_SECRET` | `qr-secret-min-32-chars` | QR token signing secret |
| `AWS_REGION` | `ap-southeast-1` | AWS region (for S3, SQS, etc.) |
| `S3_BUCKET_ASSETS` | `iorder-assets-dev` | S3 bucket for images |
| `SQS_ORDER_EVENTS_URL` | `http://localhost:4566/...` | SQS queue URL (LocalStack for dev) |

### 2.3 Frontend Setup

```bash
cd frontend

# Install dependencies
pnpm install

# Copy environment config
cp .env.example .env.local
# Edit NEXT_PUBLIC_API_URL=http://localhost:8080

# Start development server
pnpm dev
# App starts at http://localhost:3000
```

### 2.4 Mobile Setup

```bash
cd mobile

# Install dependencies
pnpm install

# Copy environment config
# Edit app.json: extra.apiUrl = "http://localhost:8080"

# Start Expo development server
pnpm start

# Run on iOS simulator (macOS only)
pnpm ios

# Run on Android emulator
pnpm android
```

### 2.5 Docker Compose (All Services)

```bash
# Start all local services
docker-compose up -d

# Services started:
# - postgres:5432
# - redis:6379
# - localstack:4566 (AWS services mock: S3, SQS, SNS)
# - backend:8080
# - frontend:3000

# View logs
docker-compose logs -f backend

# Stop all
docker-compose down
```

---

## 3. Database Migrations

```bash
cd backend

# Run all pending migrations
make migrate-up

# Rollback last migration
make migrate-down

# Create new migration
make migrate-create NAME=add_reservations_table
# Creates:
# migrations/000002_add_reservations_table.up.sql
# migrations/000002_add_reservations_table.down.sql

# Check migration status
make migrate-status

# Force migration to specific version (use with caution)
migrate -database "$DATABASE_URL" -path migrations force 1
```

**Migration naming convention:**
```
NNNNNN_<description>.up.sql     # Apply migration
NNNNNN_<description>.down.sql   # Rollback migration
```

---

## 4. Running Tests

### 4.1 Backend Tests

```bash
cd backend

# Run all tests
go test ./...

# Run with race detector (required before PR)
go test -race ./...

# Run tests for specific package
go test -v ./internal/handler/...

# Run with coverage
go test -cover ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

### 4.2 Frontend Tests

```bash
cd frontend

# Type check
pnpm typecheck

# Lint
pnpm lint

# Unit tests
pnpm test

# Build (validates no build errors)
pnpm build
```

### 4.3 Mobile Tests

```bash
cd mobile

# Type check
pnpm typecheck

# Lint
pnpm lint

# Unit tests
pnpm test
```

---

## 5. Code Quality Gates

Before submitting a PR, run all quality checks:

```bash
# Backend
cd backend
go build ./...           # Must pass
go test -race ./...      # Must pass
golangci-lint run        # 0 warnings
govulncheck ./...        # Must pass

# Frontend
cd frontend
pnpm typecheck           # 0 errors
pnpm lint                # 0 warnings
pnpm build               # Must succeed

# Mobile
cd mobile
pnpm typecheck           # 0 errors
pnpm lint                # 0 warnings
pnpm test                # All tests pass
```

---

## 6. AWS Infrastructure Deployment

### 6.1 Infrastructure Overview

```
infra/
├── bin/infra.ts          # CDK app entry
└── lib/
    ├── network-stack.ts  # VPC, subnets, NAT, security groups
    ├── database-stack.ts # Aurora Serverless v2, ElastiCache Redis
    ├── compute-stack.ts  # ECS Fargate cluster, services, ALB
    ├── auth-stack.ts     # Cognito User Pools
    ├── storage-stack.ts  # S3 buckets, CloudFront
    └── messaging-stack.ts # SQS, SNS, EventBridge
```

### 6.2 First-Time Bootstrap

```bash
cd infra

# Install CDK dependencies
pnpm install

# Bootstrap CDK in your AWS account (run once per account/region)
cdk bootstrap aws://ACCOUNT_ID/ap-southeast-1

# Synthesize CloudFormation templates (dry run)
cdk synth

# Deploy all stacks (first time)
cdk deploy --all
```

### 6.3 Per-Environment Deployment

```bash
# Deploy to dev
cdk deploy --all --context env=dev

# Deploy to staging
cdk deploy --all --context env=staging

# Deploy to production (requires MFA/approval)
cdk deploy --all --context env=prod
```

### 6.4 Deploy Individual Stack

```bash
# Deploy only compute stack (faster iteration)
cdk deploy IOrder-prod-ComputeStack

# Deploy with approval prompt for each change
cdk deploy IOrder-prod-ComputeStack --require-approval never
```

### 6.5 Destroy Environment (dev only)

```bash
# WARNING: Destroys all resources including database!
cdk destroy --all --context env=dev
```

---

## 7. Application Deployment (CI/CD)

### 7.1 Automated Pipeline (GitHub Actions)

Deployment is automatic:
- **Push to `develop`** → deploys to `dev` environment
- **Merge to `main`** → deploys to `staging` environment
- **Create tag `v*`** → triggers production deployment (requires manual approval)

### 7.2 Manual Deployment

```bash
# Build and push Docker image
cd backend
IMAGE_TAG=$(git rev-parse --short HEAD)
ECR_URI="123456789.dkr.ecr.ap-southeast-1.amazonaws.com/iorder-backend"

# Login to ECR
aws ecr get-login-password --region ap-southeast-1 \
  | docker login --username AWS --password-stdin $ECR_URI

# Build
docker build -t $ECR_URI:$IMAGE_TAG -t $ECR_URI:latest .

# Push
docker push $ECR_URI:$IMAGE_TAG
docker push $ECR_URI:latest

# Force ECS to redeploy with new image
aws ecs update-service \
  --cluster iorder-prod \
  --service iorder-backend \
  --force-new-deployment \
  --region ap-southeast-1
```

### 7.3 Run Migrations in Production

```bash
# Migrations run automatically in CI/CD pipeline after deployment
# For manual run:
aws ecs run-task \
  --cluster iorder-prod \
  --task-definition iorder-migrate \
  --launch-type FARGATE \
  --network-configuration "awsvpcConfiguration={subnets=[subnet-xxx],securityGroups=[sg-xxx]}"
```

---

## 8. Secrets Management

### 8.1 AWS Secrets Manager Structure

```
/iorder/prod/database-url          # Aurora connection string
/iorder/prod/redis-url             # ElastiCache connection string
/iorder/prod/jwt-secret            # JWT signing secret
/iorder/prod/table-qr-secret       # QR token signing secret
/iorder/prod/momo/secret-key       # MoMo payment secret
/iorder/prod/zalopay/key2          # ZaloPay HMAC key
/iorder/prod/adyen/api-key         # Adyen API key
/iorder/prod/adyen/hmac-key        # Adyen webhook HMAC key
/iorder/prod/ses/smtp-password     # AWS SES SMTP password
```

### 8.2 Creating/Updating Secrets

```bash
# Create secret
aws secretsmanager create-secret \
  --name /iorder/prod/jwt-secret \
  --secret-string "your-secret-value" \
  --region ap-southeast-1

# Update secret
aws secretsmanager put-secret-value \
  --secret-id /iorder/prod/jwt-secret \
  --secret-string "new-secret-value" \
  --region ap-southeast-1

# Retrieve secret (for debugging)
aws secretsmanager get-secret-value \
  --secret-id /iorder/prod/jwt-secret \
  --query SecretString \
  --output text
```

### 8.3 Local Development (No AWS Secrets Manager)

Trong local dev, sử dụng `.env` file thay vì Secrets Manager:
```bash
cp backend/.env.example backend/.env
# Edit .env với giá trị local
```

---

## 9. Monitoring & Operations

### 9.1 CloudWatch Dashboards

**Business Metrics Dashboard** (`IOrder-Business-Metrics`):
- Orders per minute (by tenant, by status)
- Payment success rate
- Average order processing time
- Active table sessions

**Technical Metrics Dashboard** (`IOrder-Technical`):
- ECS CPU/Memory utilization
- Aurora ACU usage
- Redis cache hit rate
- ALB request count, 5xx rate, latency p50/p95/p99
- SQS queue depth

### 9.2 Key Alarms

| Alarm | Threshold | Action |
|---|---|---|
| API 5xx rate | > 1% per 5 min | PagerDuty + Slack |
| API latency p95 | > 500ms | Slack warning |
| Payment success rate | < 95% per 5 min | PagerDuty urgent |
| ECS task health | < 50% healthy | PagerDuty |
| Aurora CPU | > 80% per 10 min | Slack warning |
| SQS dead-letter queue | > 0 messages | Slack + investigate |

### 9.3 Log Queries (CloudWatch Insights)

```sql
-- Find all errors in last 1 hour
fields @timestamp, level, message, error
| filter level = "ERROR"
| sort @timestamp desc
| limit 100

-- Orders created per minute
fields @timestamp
| filter message = "order created"
| stats count() by bin(1m)

-- Slow API requests (> 500ms)
fields @timestamp, path, duration_ms, tenant_id
| filter duration_ms > 500
| sort duration_ms desc
| limit 50
```

### 9.4 Accessing Logs

```bash
# View backend logs (latest 100 lines)
aws logs tail /ecs/iorder-backend --follow

# Query logs with CloudWatch Insights
aws logs start-query \
  --log-group-name /ecs/iorder-backend \
  --start-time $(date -d '1 hour ago' +%s) \
  --end-time $(date +%s) \
  --query-string 'fields @timestamp, message | filter level = "ERROR" | limit 20'
```

### 9.5 Connecting to Database (Bastion)

```bash
# Connect to Aurora via SSM Session Manager (no SSH key needed)
aws ssm start-session \
  --target i-BASTION_INSTANCE_ID \
  --document-name AWS-StartPortForwardingSessionToRemoteHost \
  --parameters host="aurora-endpoint.cluster-xxx.ap-southeast-1.rds.amazonaws.com",portNumber="5432",localPortNumber="5432"

# In another terminal
psql "host=localhost port=5432 dbname=iorder user=iorder_admin"
```

---

## 10. Troubleshooting

### 10.1 Common Issues

**ECS task keeps restarting:**
```bash
# Check task stopped reason
aws ecs describe-tasks \
  --cluster iorder-prod \
  --tasks $(aws ecs list-tasks --cluster iorder-prod --service-name iorder-backend --query 'taskArns[0]' --output text)

# Check CloudWatch logs for crash
aws logs tail /ecs/iorder-backend --since 5m
```

**Database connection errors:**
```bash
# Check Aurora status
aws rds describe-db-clusters --db-cluster-identifier iorder-prod

# Check security group allows ECS → Aurora
aws ec2 describe-security-groups --group-ids sg-xxx
```

**Payment webhooks not received:**
```bash
# Check SQS dead-letter queue
aws sqs get-queue-attributes \
  --queue-url https://sqs.ap-southeast-1.amazonaws.com/xxx/iorder-payment-dlq \
  --attribute-names ApproximateNumberOfMessages
```

**QR code not loading:**
- Verify CloudFront distribution is deployed
- Check S3 bucket CORS settings
- Verify QR token expiry (6 hour default)

### 10.2 Rollback Procedure

```bash
# Rollback ECS service to previous task definition
PREVIOUS_TASK_DEF=$(aws ecs describe-services \
  --cluster iorder-prod \
  --services iorder-backend \
  --query 'services[0].deployments[-1].taskDefinition' \
  --output text)

aws ecs update-service \
  --cluster iorder-prod \
  --service iorder-backend \
  --task-definition $PREVIOUS_TASK_DEF

# Rollback database migration
cd backend
make migrate-down  # Rolls back 1 version
```

---

## 11. Onboarding a New Restaurant

### 11.1 Via Platform Admin API

```bash
# Create new tenant
curl -X POST https://api.iorder.vn/api/v1/tenants \
  -H "Authorization: Bearer $ADMIN_JWT" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Phở Hà Nội",
    "slug": "pho-ha-noi",
    "timezone": "Asia/Ho_Chi_Minh",
    "currency": "VND",
    "owner_email": "owner@phonhanoi.vn",
    "owner_name": "Nguyễn Văn A"
  }'
```

### 11.2 Restaurant Setup Checklist

After creating tenant:
- [ ] Restaurant owner receives welcome email with credentials
- [ ] Owner logs in to dashboard, updates profile (logo, description, hours)
- [ ] Create menu categories
- [ ] Add menu items with images
- [ ] Create floor plan(s)
- [ ] Add tables to floor plan
- [ ] Download/print QR codes for all tables
- [ ] Configure accepted payment methods
- [ ] Test with a sample order end-to-end
- [ ] Go live!

---

## 12. Development Best Practices

### 12.1 Branch Workflow

```bash
# Create feature branch
git checkout -b feat/42-add-qr-regeneration

# Make changes, commit with conventional commits
git commit -m "feat(backend): add QR token regeneration endpoint"

# Push and create PR
git push origin feat/42-add-qr-regeneration
gh pr create --title "feat(backend): add QR token regeneration endpoint" --body "..."
```

### 12.2 Adding a New API Endpoint

1. Add migration if new table needed (`make migrate-create NAME=add_table`)
2. Define types in `internal/domain/`
3. Add repository method in `internal/repository/`
4. Implement handler in `internal/handler/`
5. Register route in `internal/router/router.go`
6. Write handler test in `internal/handler/*_test.go`
7. Update API docs in `docs/architecture/api-contracts/`

### 12.3 Adding a New Payment Provider

1. Create `internal/infrastructure/payment/{provider}/provider.go`
2. Implement `PaymentProvider` interface
3. Add webhook handler in `internal/handler/payment.go`
4. Add provider secret to AWS Secrets Manager
5. Register provider in payment service factory
6. Add sandbox tests in `internal/infrastructure/payment/{provider}/provider_test.go`
7. Update ADR 0004 with new provider

---

## 13. Contact & Support

| Role | Contact |
|---|---|
| Tech Lead | Create GitHub issue with `area:infra` label |
| Backend issues | Create GitHub issue with `area:backend` label |
| Frontend/Mobile issues | Create GitHub issue with `area:frontend`/`area:mobile` label |
| Production incident | Create P0 issue immediately, tag Tech Lead |
