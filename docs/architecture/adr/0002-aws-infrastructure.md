# ADR 0002: AWS Infrastructure Architecture

## Status

Accepted

## Date

2026-03-02

## Context

IOrder requires a cloud architecture that can:
- Support multi-tenancy (multiple independent restaurants)
- Auto-scale according to real traffic patterns (lunch and dinner peak hours)
- Guarantee uptime ≥ 99.9% (SLA)
- Optimize cost for the early stage (pay-as-you-go)
- Provide a good developer experience (fast CI/CD, easy debugging)

AWS was chosen as the primary cloud provider because:
- The team already has experience with AWS
- A strong managed services ecosystem (Aurora, Fargate, Cognito)
- The `ap-southeast-1` region (Singapore) is close to the Vietnamese market, giving low latency
- AWS has partnerships with Vietnamese payment gateways (VNPay, MoMo)

## Decision

### Primary Region: `ap-southeast-1` (Singapore)

### Compute: ECS Fargate

**ECS Fargate** was chosen over EC2 or EKS because:
- No server/node management — reduces operational overhead
- Container-native, well-suited to Go microservices
- Built-in auto-scaling
- Reasonable cost for small-to-medium scale
- Simpler than Kubernetes for current team size

**Rejected alternatives:**
- EC2: Requires OS management, patching, and more complex security groups
- EKS: Over-engineered at this stage; fixed control plane cost of $0.10/h
- Lambda (for API): Cold-start latency unsuitable for an API server; no WebSocket support
- App Runner: Less flexible than Fargate (no custom networking, no GPU)

### Database: Aurora PostgreSQL Serverless v2

**Aurora PostgreSQL Serverless v2** was chosen over RDS or DynamoDB because:
- Auto-scales from 0 ACU (pauses when idle) — saves cost in dev/staging
- Full PostgreSQL compatibility — works with golang-migrate and pgx
- Built-in Multi-AZ, failover in < 30 seconds
- Row-Level Security (RLS) for multi-tenant isolation
- JSONB support for flexible schema (modifiers, settings)

**Rejected alternatives:**
- RDS PostgreSQL: Cannot pause; higher cost for non-production environments
- DynamoDB: Not suitable for complex relational queries (orders + items + payments)
- PlanetScale/Neon: Vendor risk; not AWS-native

### Cache/Pub-Sub: ElastiCache Redis

**ElastiCache Redis** was chosen because:
- Session caching, rate limiting, and real-time pub/sub in a single service
- Redis Pub/Sub suits WebSocket order-status broadcasting
- Managed service with automatic failover
- Widely used in the Go ecosystem (go-redis)

### CDN: CloudFront

Required for:
- Static assets (menu images, frontend bundles): TTL 24h
- API edge caching for public endpoints (menu GET): TTL 5 minutes
- Integrated WAF
- DDoS protection (AWS Shield Standard, included free)

### Auth: Cognito

**Cognito** was chosen because:
- Fully managed; no need to build authentication from scratch
- Social login (Google, Facebook) for customers
- MFA for staff and admin
- JWT tokens compatible with Go middleware
- Native integration with API Gateway and ALB

**Rejected alternatives:**
- Auth0: High cost at scale; external vendor dependency
- Keycloak: Requires self-hosting and ongoing maintenance
- Custom JWT: Security risk; not recommended

### Real-time: API Gateway WebSocket + Redis Pub/Sub

Architecture:
1. Client connects to API Gateway WebSocket endpoint
2. `$connect` Lambda stores `connectionId` → `{tenantId, orderId, role}` in DynamoDB
3. API server publishes events to Redis channel when order state changes
4. Lambda subscriber receives from Redis and sends via API Gateway Management API to relevant connections

Reason for not using WebSocket directly on ECS:
- ALB WebSocket connections have a limited idle timeout (max 4000s)
- API Gateway WS has no such limit, auto-scales, and manages connection state automatically

### Infrastructure as Code: AWS CDK (TypeScript)

**CDK** was chosen over Terraform or raw CloudFormation because:
- TypeScript — familiar to the frontend team
- Higher-level abstractions (L2/L3 constructs) reduce boilerplate
- Type-safe with full IDE autocomplete
- Native AWS support

**Rejected alternatives:**
- Terraform: HCL syntax less familiar to the team; provider sometimes lags behind the latest AWS releases
- CloudFormation YAML: Verbose; no type safety
- Pulumi: Smaller community; fewer examples

## Consequences

### Positive

- Zero-management infrastructure (Fargate, Aurora Serverless, Cognito)
- Auto-scaling handles peak hours automatically
- Strong security baseline from AWS managed services
- Developers can focus on business logic instead of infrastructure
- Cost-optimized: dev/staging Aurora Serverless pauses when idle

### Negative

- AWS vendor lock-in: migrating to another cloud would be expensive
- CDK learning curve for team members not familiar with it
- Aurora Serverless v2 has a ~1-second cold start after a pause — not ideal for developers working at night (needs a warm-up strategy)
- ElastiCache Redis requires a minimum of 2 nodes (~$50/month) even under low traffic

### Mitigation

- Vendor lock-in: abstraction layer in code (repository pattern); most infrastructure is commodity
- Aurora cold start: disable pause on staging, or schedule a warm-up Lambda
