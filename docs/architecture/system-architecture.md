# IOrder — System Architecture

> **Author:** Tech Lead
> **Date:** 2026-03-02
> **Version:** 1.0
> **Status:** Accepted

---

## 1. Overview

IOrder là hệ thống multi-tenant SaaS cho phép nhiều nhà hàng vận hành trên cùng nền tảng. Kiến trúc ưu tiên AWS-native services, hỗ trợ horizontal scaling, và đảm bảo availability ≥ 99.9%.

---

## 2. Architecture Principles

1. **Cloud-native first** — Ưu tiên managed AWS services để giảm operational overhead
2. **API-first** — Tất cả tính năng được expose qua REST API; FE/Mobile là client thuần
3. **Event-driven for real-time** — Order state changes được push qua WebSocket/SSE
4. **Multi-tenant isolation** — Dữ liệu nhà hàng được cô lập ở tầng application (row-level tenant ID)
5. **Security by default** — TLS everywhere, least-privilege IAM, secrets in AWS Secrets Manager
6. **Observability** — Logs, metrics, traces từ ngày đầu

---

## 3. High-Level Architecture

```
┌─────────────────────────────────────────────────────────────────────────┐
│                          AWS Cloud (ap-southeast-1)                      │
│                                                                           │
│  ┌──────────────┐    ┌──────────────────────────────────────────────┐   │
│  │   CloudFront │    │              VPC (10.0.0.0/16)               │   │
│  │   + WAF      │    │                                              │   │
│  │              │    │  ┌─────────────────────────────────────┐    │   │
│  │  Static      │    │  │   Public Subnets (2 AZs)            │    │   │
│  │  Assets      │    │  │  ┌───────────┐  ┌───────────────┐  │    │   │
│  │  (S3)        │    │  │  │    ALB    │  │  NAT Gateway  │  │    │   │
│  └──────────────┘    │  │  └─────┬─────┘  └───────────────┘  │    │   │
│                       │  └────────┼────────────────────────────┘    │   │
│  ┌──────────────┐    │           │                                  │   │
│  │   Route 53   │    │  ┌────────┼────────────────────────────┐    │   │
│  │   (DNS)      │    │  │ Private Subnets (2 AZs)             │    │   │
│  └──────────────┘    │  │        │                             │    │   │
│                       │  │  ┌─────▼─────────────────────┐      │    │   │
│  ┌──────────────┐    │  │  │    ECS Fargate Cluster     │      │    │   │
│  │   Cognito    │    │  │  │  ┌──────┐  ┌──────────┐   │      │    │   │
│  │   (Auth)     │    │  │  │  │ API  │  │  Worker  │   │      │    │   │
│  └──────────────┘    │  │  │  │ svc  │  │  (SQS)   │   │      │    │   │
│                       │  │  │  └──────┘  └──────────┘   │      │    │   │
│                       │  │  └────────────────────────────┘      │    │   │
│                       │  │                                        │    │   │
│                       │  │  ┌──────────────────────────────────┐ │    │   │
│                       │  │  │  Data Layer                      │ │    │   │
│                       │  │  │  ┌──────────┐  ┌──────────────┐  │ │    │   │
│                       │  │  │  │  Aurora  │  │ ElastiCache  │  │ │    │   │
│                       │  │  │  │ Postgres │  │ Redis        │  │ │    │   │
│                       │  │  │  └──────────┘  └──────────────┘  │ │    │   │
│                       │  │  └──────────────────────────────────┘ │    │   │
│                       │  └────────────────────────────────────────┘    │   │
│                       └──────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────────────────┘
```

---

## 4. AWS Services Selection

### 4.1 Compute

| Service | Use Case | Justification |
|---|---|---|
| **ECS Fargate** | API server, background workers | Serverless containers, auto-scaling, no EC2 management |
| **Lambda** | Webhook processors, image resize, scheduled jobs | Event-driven, pay-per-use, natural fit for spiky workloads |
| **ECR** | Container registry | Native ECS integration, image scanning |

### 4.2 Networking

| Service | Use Case |
|---|---|
| **VPC** | Network isolation với public/private subnets |
| **ALB** (Application Load Balancer) | L7 load balancing, SSL termination, WebSocket support |
| **CloudFront** | CDN cho static assets (menu images, frontend) |
| **Route 53** | DNS management, health checks, failover |
| **WAF** | Web Application Firewall trước CloudFront và ALB |
| **NAT Gateway** | Outbound internet cho private subnets |

### 4.3 Database & Storage

| Service | Use Case | Config |
|---|---|---|
| **Aurora PostgreSQL Serverless v2** | Primary database (multi-tenant) | Auto-scales 0.5–128 ACU, Multi-AZ |
| **ElastiCache Redis** | Session cache, real-time order pub/sub, rate limiting | Cluster mode, 2 nodes, Multi-AZ |
| **S3** | Menu images, QR code PDFs, e-receipts, backups | Versioning enabled, lifecycle rules |
| **DynamoDB** | QR token store, session store (high-read) | On-demand capacity |

### 4.4 Messaging & Real-time

| Service | Use Case |
|---|---|
| **SQS** | Order event queue, payment webhooks, async notifications |
| **SNS** | Fan-out notifications (SMS, email, push triggers) |
| **API Gateway (WebSocket)** | Real-time order status push to customer browser/app |
| **EventBridge** | Domain events routing (OrderPlaced, PaymentCompleted,...) |

### 4.5 Auth & Security

| Service | Use Case |
|---|---|
| **Cognito** | Customer auth (social login, phone OTP) |
| **Cognito User Pools** | Staff/owner auth với MFA |
| **Secrets Manager** | Database credentials, API keys, payment secrets |
| **KMS** | Encryption keys cho S3, RDS, Secrets Manager |
| **IAM** | Fine-grained access control |

### 4.6 Observability

| Service | Use Case |
|---|---|
| **CloudWatch Logs** | Application logs, API access logs |
| **CloudWatch Metrics** | Custom business metrics (orders/min, payment success rate) |
| **CloudWatch Alarms** | Alerting trên error rate, latency, DB connections |
| **X-Ray** | Distributed tracing |
| **OpenSearch** | Log analytics, full-text search (menu search) |

### 4.7 CI/CD & DevOps

| Service | Use Case |
|---|---|
| **GitHub Actions** | CI/CD pipeline |
| **ECR** | Docker image registry |
| **CodeDeploy** (via ECS rolling update) | Zero-downtime deployment |
| **CloudFormation / CDK** | Infrastructure as Code |

### 4.8 Notifications

| Service | Use Case |
|---|---|
| **SNS** | SMS gateway |
| **SES** | Transactional emails (e-receipt, booking confirmation) |
| **Pinpoint** | Push notifications, marketing campaigns cho membership |

---

## 5. Application Architecture

### 5.1 Backend Services (Go + Gin)

```
backend/
├── cmd/
│   ├── server/          # HTTP API server
│   └── worker/          # SQS consumer (background jobs)
├── internal/
│   ├── config/          # Env-based config
│   ├── db/              # Aurora PostgreSQL connection pool
│   ├── cache/           # Redis client
│   ├── handler/         # HTTP handlers (one file per domain)
│   │   ├── auth.go      # POST /api/v1/auth/*
│   │   ├── restaurant.go # CRUD /api/v1/restaurants
│   │   ├── menu.go      # CRUD /api/v1/restaurants/:id/menu
│   │   ├── table.go     # /api/v1/restaurants/:id/tables
│   │   ├── order.go     # /api/v1/orders
│   │   ├── payment.go   # /api/v1/payments
│   │   └── membership.go # /api/v1/membership
│   ├── middleware/
│   │   ├── auth.go      # JWT validation (Cognito)
│   │   ├── tenant.go    # Tenant resolution from subdomain/header
│   │   ├── logger.go    # slog request logger
│   │   ├── ratelimit.go # Redis-based rate limiting
│   │   └── recovery.go  # JSON panic recovery
│   ├── domain/          # Domain models & business logic
│   │   ├── order/       # Order state machine
│   │   ├── payment/     # Payment processing
│   │   └── membership/  # Points & tiers
│   ├── repository/      # Database access layer
│   ├── event/           # EventBridge publisher
│   └── router/          # Gin route registration
└── migrations/          # SQL migrations
```

### 5.2 Frontend (Next.js 15)

```
frontend/src/
├── app/
│   ├── (customer)/          # Customer-facing PWA
│   │   ├── [restaurant]/    # Restaurant landing
│   │   │   └── table/[id]/  # Table ordering flow
│   │   └── membership/      # Member profile
│   ├── (dashboard)/         # Restaurant owner/staff dashboard
│   │   ├── layout.tsx       # Auth guard
│   │   ├── menu/            # Menu management
│   │   ├── tables/          # Floor plan editor
│   │   ├── orders/          # Live order management
│   │   ├── kds/             # Kitchen Display System
│   │   ├── payments/        # Payment history
│   │   └── analytics/       # Reports & charts
│   └── (admin)/             # Platform admin
├── components/
│   ├── ui/                  # Base UI (shadcn/ui)
│   ├── floor-map/           # Interactive floor plan editor
│   ├── order-tracker/       # Real-time order status
│   └── payment/             # Payment method selector
├── lib/
│   ├── api/client.ts        # Typed API client
│   ├── ws/client.ts         # WebSocket client (order updates)
│   └── auth/                # Cognito auth helpers
└── types/                   # Shared TypeScript types
```

### 5.3 Mobile (React Native / Expo)

```
mobile/src/
├── app/
│   ├── (tabs)/              # Main navigation tabs
│   ├── scan/                # QR scanner
│   ├── menu/                # Menu browsing
│   ├── order/               # Order flow
│   ├── payment/             # Payment flow
│   └── profile/             # Membership profile
├── components/              # Reusable components
├── lib/
│   ├── api/client.ts        # API client (matches frontend pattern)
│   └── notifications/       # Push notification setup
└── types/                   # Shared types
```

---

## 6. Database Schema Design

### 6.1 Multi-tenancy Strategy

**Strategy: Shared Database, Shared Schema với `tenant_id` column**

Tất cả bảng có `tenant_id UUID NOT NULL` và Row-Level Security (RLS) trên PostgreSQL.

```sql
-- Enable RLS on orders table
ALTER TABLE orders ENABLE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON orders
  USING (tenant_id = current_setting('app.current_tenant')::UUID);
```

Lý do chọn shared schema thay vì separate schema per tenant:
- Đơn giản hơn trong việc migration schema
- Chi phí thấp hơn khi số tenant tăng
- Dễ dàng chạy cross-tenant analytics (platform level)
- Aurora Serverless auto-scale xử lý tốt mixed workloads

### 6.2 Core Tables

```sql
-- Tenants (restaurants)
CREATE TABLE tenants (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    slug        VARCHAR(100) UNIQUE NOT NULL,  -- used in QR URL
    name        VARCHAR(255) NOT NULL,
    logo_url    TEXT,
    timezone    VARCHAR(50) NOT NULL DEFAULT 'Asia/Ho_Chi_Minh',
    currency    CHAR(3) NOT NULL DEFAULT 'VND',
    is_active   BOOLEAN NOT NULL DEFAULT true,
    settings    JSONB NOT NULL DEFAULT '{}',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Floor plans
CREATE TABLE floor_plans (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id   UUID NOT NULL REFERENCES tenants(id),
    name        VARCHAR(100) NOT NULL,  -- "Tầng 1", "Sân thượng"
    floor_level INT NOT NULL DEFAULT 1,
    is_active   BOOLEAN NOT NULL DEFAULT true,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Tables
CREATE TABLE restaurant_tables (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id   UUID NOT NULL REFERENCES tenants(id),
    floor_plan_id UUID NOT NULL REFERENCES floor_plans(id),
    name        VARCHAR(50) NOT NULL,  -- "A1", "VIP 01"
    capacity    INT NOT NULL DEFAULT 4,
    pos_x       FLOAT NOT NULL DEFAULT 0,  -- position on floor map
    pos_y       FLOAT NOT NULL DEFAULT 0,
    shape       VARCHAR(20) NOT NULL DEFAULT 'rectangle',  -- circle, rectangle
    status      VARCHAR(20) NOT NULL DEFAULT 'available',  -- available, occupied, reserved
    qr_token    TEXT,          -- current active QR JWT
    qr_expires_at TIMESTAMPTZ,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Menu categories
CREATE TABLE menu_categories (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id   UUID NOT NULL REFERENCES tenants(id),
    name        VARCHAR(100) NOT NULL,
    description TEXT,
    sort_order  INT NOT NULL DEFAULT 0,
    is_active   BOOLEAN NOT NULL DEFAULT true,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Menu items
CREATE TABLE menu_items (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id       UUID NOT NULL REFERENCES tenants(id),
    category_id     UUID NOT NULL REFERENCES menu_categories(id),
    name            VARCHAR(255) NOT NULL,
    description     TEXT,
    price           BIGINT NOT NULL,  -- stored in smallest unit (xu/cent)
    image_url       TEXT,
    tags            TEXT[] NOT NULL DEFAULT '{}',  -- vegetarian, vegan, ...
    is_available    BOOLEAN NOT NULL DEFAULT true,
    sort_order      INT NOT NULL DEFAULT 0,
    serve_from      TIME,   -- null = always available
    serve_until     TIME,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Item modifiers (e.g., size, spice level)
CREATE TABLE item_modifiers (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id   UUID NOT NULL REFERENCES tenants(id),
    item_id     UUID NOT NULL REFERENCES menu_items(id),
    name        VARCHAR(100) NOT NULL,  -- "Kích thước", "Mức độ cay"
    required    BOOLEAN NOT NULL DEFAULT false,
    options     JSONB NOT NULL  -- [{"name":"Nhỏ","price_delta":0},{"name":"Lớn","price_delta":5000}]
);

-- Orders
CREATE TABLE orders (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id       UUID NOT NULL REFERENCES tenants(id),
    table_id        UUID NOT NULL REFERENCES restaurant_tables(id),
    session_id      UUID NOT NULL,  -- groups multiple order rounds at same table
    status          VARCHAR(30) NOT NULL DEFAULT 'DRAFT',
    customer_id     UUID REFERENCES customers(id),  -- null if not logged in
    subtotal        BIGINT NOT NULL DEFAULT 0,
    discount_amount BIGINT NOT NULL DEFAULT 0,
    total           BIGINT NOT NULL DEFAULT 0,
    notes           TEXT,
    placed_at       TIMESTAMPTZ,
    confirmed_at    TIMESTAMPTZ,
    ready_at        TIMESTAMPTZ,
    served_at       TIMESTAMPTZ,
    paid_at         TIMESTAMPTZ,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Order items
CREATE TABLE order_items (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id        UUID NOT NULL REFERENCES orders(id),
    tenant_id       UUID NOT NULL REFERENCES tenants(id),
    item_id         UUID NOT NULL REFERENCES menu_items(id),
    item_name       VARCHAR(255) NOT NULL,  -- snapshot at order time
    item_price      BIGINT NOT NULL,        -- snapshot at order time
    quantity        INT NOT NULL DEFAULT 1,
    modifiers       JSONB NOT NULL DEFAULT '[]',
    notes           TEXT,
    status          VARCHAR(20) NOT NULL DEFAULT 'pending'  -- pending, preparing, ready, served
);

-- Order state transitions (audit log)
CREATE TABLE order_events (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id    UUID NOT NULL REFERENCES orders(id),
    tenant_id   UUID NOT NULL REFERENCES tenants(id),
    from_status VARCHAR(30),
    to_status   VARCHAR(30) NOT NULL,
    actor_id    UUID,        -- staff/customer who triggered transition
    actor_type  VARCHAR(20), -- staff, customer, system
    metadata    JSONB NOT NULL DEFAULT '{}',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Payments
CREATE TABLE payments (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id       UUID NOT NULL REFERENCES tenants(id),
    order_id        UUID NOT NULL REFERENCES orders(id),
    method          VARCHAR(30) NOT NULL,  -- cash, vietqr, momo, zalopay, shopee_pay, adyen
    status          VARCHAR(20) NOT NULL DEFAULT 'pending',  -- pending, processing, completed, failed, refunded
    amount          BIGINT NOT NULL,
    currency        CHAR(3) NOT NULL DEFAULT 'VND',
    provider_ref    TEXT,   -- external payment reference
    provider_data   JSONB NOT NULL DEFAULT '{}',
    expires_at      TIMESTAMPTZ,
    completed_at    TIMESTAMPTZ,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Customers (membership)
CREATE TABLE customers (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id       UUID NOT NULL REFERENCES tenants(id),
    cognito_sub     TEXT,       -- Cognito user ID
    phone           VARCHAR(20),
    email           VARCHAR(255),
    full_name       VARCHAR(255),
    tier            VARCHAR(20) NOT NULL DEFAULT 'bronze',
    total_spend     BIGINT NOT NULL DEFAULT 0,
    points_balance  INT NOT NULL DEFAULT 0,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(tenant_id, phone),
    UNIQUE(tenant_id, email)
);

-- Membership transactions (points)
CREATE TABLE membership_transactions (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id       UUID NOT NULL REFERENCES tenants(id),
    customer_id     UUID NOT NULL REFERENCES customers(id),
    order_id        UUID REFERENCES orders(id),
    type            VARCHAR(20) NOT NULL,  -- earn, redeem, expire, adjust
    points          INT NOT NULL,          -- positive=earn, negative=redeem
    balance_after   INT NOT NULL,
    description     TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

---

## 7. API Design

### 7.1 Base URL & Versioning

```
https://api.iorder.vn/api/v1/
```

### 7.2 Authentication

Tất cả API (trừ public customer endpoints) yêu cầu:
```
Authorization: Bearer <cognito-jwt>
X-Tenant-ID: <tenant-uuid>   (hoặc resolve từ subdomain)
```

### 7.3 Key Endpoints

#### Tenant / Restaurant
```
POST   /api/v1/tenants              # Platform admin: create restaurant
GET    /api/v1/tenants/:id          # Get restaurant details
PATCH  /api/v1/tenants/:id          # Update restaurant settings
```

#### Menu
```
GET    /api/v1/menu                           # Public: get full menu (for customer QR flow)
GET    /api/v1/menu/categories                # List categories
POST   /api/v1/menu/categories                # Create category
PUT    /api/v1/menu/categories/:id            # Update category
DELETE /api/v1/menu/categories/:id            # Delete category
GET    /api/v1/menu/items                     # List items
POST   /api/v1/menu/items                     # Create item
PUT    /api/v1/menu/items/:id                 # Update item
PATCH  /api/v1/menu/items/:id/availability    # Toggle availability
DELETE /api/v1/menu/items/:id                 # Delete item
```

#### Tables & Floor Plans
```
GET    /api/v1/floor-plans                    # List floor plans
POST   /api/v1/floor-plans                    # Create floor plan
PUT    /api/v1/floor-plans/:id                # Update floor plan
GET    /api/v1/floor-plans/:id/tables         # List tables on floor plan
POST   /api/v1/tables                         # Create table
PUT    /api/v1/tables/:id                     # Update table (position, name, capacity)
DELETE /api/v1/tables/:id                     # Delete table
POST   /api/v1/tables/:id/qr                  # Regenerate QR token
GET    /api/v1/tables/:id/qr/download         # Download QR as PNG/PDF
```

#### Orders
```
POST   /api/v1/orders                         # Customer: create order (draft)
GET    /api/v1/orders/:id                     # Get order details
PATCH  /api/v1/orders/:id/items               # Update cart items
POST   /api/v1/orders/:id/submit              # Customer: submit order
POST   /api/v1/orders/:id/confirm             # Staff: confirm order
POST   /api/v1/orders/:id/reject              # Staff: reject order
POST   /api/v1/orders/:id/start-preparation   # Kitchen: mark IN_PREPARATION
POST   /api/v1/orders/:id/ready               # Kitchen: mark READY
POST   /api/v1/orders/:id/serve               # Staff: mark SERVED
POST   /api/v1/orders/:id/request-payment     # Customer/Staff: request payment
POST   /api/v1/orders/:id/cancel              # Cancel order
GET    /api/v1/orders                         # Staff: list orders (with filters)
```

#### Payments
```
POST   /api/v1/payments                       # Initiate payment
GET    /api/v1/payments/:id                   # Get payment status
POST   /api/v1/payments/webhooks/momo         # MoMo webhook (public, HMAC-verified)
POST   /api/v1/payments/webhooks/zalopay      # ZaloPay webhook
POST   /api/v1/payments/webhooks/adyen        # Adyen webhook
```

#### Membership
```
POST   /api/v1/membership/register            # Register as member
GET    /api/v1/membership/profile             # Get member profile + points
GET    /api/v1/membership/transactions        # Points history
POST   /api/v1/membership/redeem              # Redeem points for voucher
```

#### WebSocket (Real-time)
```
wss://api.iorder.vn/ws/orders/:session_id    # Customer: subscribe to order updates
wss://api.iorder.vn/ws/kds                   # Staff: KDS live order feed
```

### 7.4 Error Response Format

```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Invalid input",
    "details": [
      { "field": "price", "message": "must be greater than 0" }
    ]
  }
}
```

---

## 8. Real-time Architecture

### 8.1 Order Status Updates Flow

```
Customer Browser ──── WebSocket ────── API Gateway (WS) ──── Lambda ──── Redis Pub/Sub
                                                                              │
Kitchen Staff ──────── WebSocket ──────────────────────────────────────────── │
                                                                              │
API Server ──── Publishes OrderStateChanged event ──────────────────────────→─┘
                        │
                        └────── SQS ──── Worker ──── SNS ──── Push/SMS
```

### 8.2 WebSocket Connection Management

- API Gateway WebSocket API manage connections
- Lambda handler stores `connectionId` → `{tenantId, orderId, role}` in DynamoDB
- API server publishes events to Redis channel
- Lambda subscriber pushes to relevant WebSocket connections

---

## 9. Payment Integration Architecture

### 9.1 Payment Flow

```
Customer selects payment method
         │
         ▼
POST /api/v1/payments (method=momo)
         │
         ▼
Payment Service creates payment record (status=pending)
         │
         ▼
Call payment provider API (create QR / deep link)
         │
         ▼
Return payment URL / QR to customer
         │
         ▼
Customer completes payment in bank app
         │
         ▼
Provider sends webhook → POST /api/v1/payments/webhooks/{provider}
         │
         ▼
Verify HMAC signature → Update payment (status=completed)
         │
         ▼
Publish PaymentCompleted event → SQS
         │
         ▼
Worker consumes → Update order status to PAID → Notify customer
```

### 9.2 Idempotency

- Tất cả payment webhook handlers phải idempotent
- Kiểm tra `provider_ref` trước khi xử lý để tránh duplicate

---

## 10. Security Architecture

### 10.1 Network Security

- ALB chỉ accept traffic từ CloudFront (IP restriction)
- ECS tasks trong private subnet, không có public IP
- Security Groups: ALB → ECS (443/8080), ECS → RDS (5432), ECS → Redis (6379)
- WAF rules: SQL injection, XSS, rate limiting, geo-blocking

### 10.2 Data Security

- Aurora PostgreSQL: encrypted at rest (KMS)
- S3: SSE-S3 mặc định, SSE-KMS cho sensitive data
- Redis: in-transit encryption (TLS), at-rest encryption
- Secrets Manager: rotate credentials tự động 90 ngày

### 10.3 Application Security

- JWT validation: verify Cognito public key
- Tenant isolation: middleware inject `tenant_id` vào mọi query
- QR token: JWT HS256 với 6-hour expiry, stored in DynamoDB
- Payment webhook: HMAC-SHA256 signature verification
- Input validation: tất cả request body validated trước khi xử lý

---

## 11. Deployment Architecture

### 11.1 Environments

| Environment | Branch | Auto-deploy | Database |
|---|---|---|---|
| `dev` | `develop` | Yes (on push) | Aurora Dev cluster |
| `staging` | `main` | Yes (on merge) | Aurora Staging cluster |
| `production` | `main` + tag | Manual approval | Aurora Prod Multi-AZ |

### 11.2 CI/CD Pipeline (GitHub Actions)

```
Push to branch
     │
     ▼
1. Lint (golangci-lint, eslint, tsc)
     │
     ▼
2. Unit Tests (go test -race, jest)
     │
     ▼
3. Build Docker images
     │
     ▼
4. Push to ECR (tag: git SHA)
     │
     ▼
5. Deploy to ECS (rolling update)
     │
     ▼
6. Run DB migrations (golang-migrate)
     │
     ▼
7. Smoke tests
     │
     ▼
[Production only] Manual approval gate
```

### 11.3 Infrastructure as Code

Toàn bộ AWS infrastructure được quản lý bằng **AWS CDK (TypeScript)** tại `infra/`.

```
infra/
├── bin/infra.ts          # CDK app entry
├── lib/
│   ├── network-stack.ts  # VPC, subnets, security groups
│   ├── database-stack.ts # Aurora, ElastiCache
│   ├── compute-stack.ts  # ECS cluster, services, ALB
│   ├── auth-stack.ts     # Cognito user pools
│   ├── storage-stack.ts  # S3 buckets
│   └── cdn-stack.ts      # CloudFront, WAF
└── config/
    ├── dev.ts
    ├── staging.ts
    └── prod.ts
```

---

## 12. Scalability & Cost Optimization

### 12.1 Auto-scaling

- **ECS Fargate**: Target tracking scaling — CPU 70%, min 2 tasks, max 20 tasks
- **Aurora Serverless v2**: Auto-scales 0.5 → 128 ACU based on load
- **ElastiCache**: Read replicas scale out manually based on traffic patterns
- **CloudFront**: Inherently scalable (edge caching)

### 12.2 Cost Optimization

- Aurora Serverless v2 pauses at 0 ACU trong dev/staging khi không có traffic
- S3 lifecycle: move old images to S3 Glacier after 90 days
- CloudFront caches menu images (TTL 24h) → giảm S3 GET requests
- Lambda cho webhook processors (pay-per-execution, không idle cost)
- Reserved instances cho production ECS tasks chạy 24/7

---

## 13. Disaster Recovery

| Scenario | RPO | RTO | Strategy |
|---|---|---|---|
| Single AZ failure | 0 | < 1 min | Multi-AZ active-passive (Aurora, ECS) |
| Region failure | 1 hour | 4 hours | S3 cross-region replication, DB snapshot restore |
| Accidental data deletion | 5 min | 30 min | Aurora Point-in-Time Recovery (35 days) |
| Application bug | 0 | 10 min | ECS rolling deploy với rollback |

---

## 14. Related ADRs

- [ADR 0001](adr/0001-record-architecture-decisions.md) — Record Architecture Decisions
- [ADR 0002](adr/0002-aws-infrastructure.md) — AWS Infrastructure Architecture
- [ADR 0003](adr/0003-order-state-machine.md) — Order State Machine Design
- [ADR 0004](adr/0004-payment-integration.md) — Payment Integration Strategy
- [ADR 0005](adr/0005-qr-table-ordering.md) — QR Code Table Ordering
