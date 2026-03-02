# ADR 0002: AWS Infrastructure Architecture

## Status

Accepted

## Date

2026-03-02

## Context

IOrder cần một kiến trúc cloud có khả năng:
- Hỗ trợ multi-tenant (nhiều nhà hàng độc lập)
- Scale tự động theo traffic thực tế (giờ cao điểm bữa trưa/tối)
- Đảm bảo uptime ≥ 99.9% (SLA)
- Chi phí tối ưu cho giai đoạn đầu (pay-as-you-go)
- Developer experience tốt (CI/CD nhanh, dễ debug)

AWS được chọn là cloud provider chính vì:
- Team đã có kinh nghiệm với AWS
- Strong managed services ecosystem (Aurora, Fargate, Cognito)
- Region `ap-southeast-1` (Singapore) gần thị trường Việt Nam, latency thấp
- AWS có partnership với các payment gateway Việt Nam (VNPay, MoMo)

## Decision

### Primary Region: `ap-southeast-1` (Singapore)

### Compute: ECS Fargate

Chọn **ECS Fargate** thay vì EC2 hoặc EKS vì:
- Không quản lý server/node — giảm operational overhead
- Container-native, phù hợp với Go microservices
- Auto-scaling tích hợp sẵn
- Chi phí hợp lý cho scale nhỏ → vừa
- Simpler than Kubernetes cho team size hiện tại

**Rejected alternatives:**
- EC2: Phải quản lý OS, patching, security groups phức tạp hơn
- EKS: Over-engineered cho giai đoạn đầu, chi phí control plane cố định $0.10/h
- Lambda (API): Cold start latency không phù hợp cho API server (WebSocket không support)
- App Runner: Kém linh hoạt hơn Fargate (no custom networking, no GPU support)

### Database: Aurora PostgreSQL Serverless v2

Chọn **Aurora PostgreSQL Serverless v2** thay vì RDS hoặc DynamoDB vì:
- Auto-scales từ 0 ACU (pause khi không dùng) — tiết kiệm chi phí dev/staging
- PostgreSQL compatibility — dùng được với golang-migrate và pgx
- Multi-AZ tích hợp, failover < 30 giây
- Row-Level Security (RLS) cho multi-tenant isolation
- JSONB support cho flexible schema (modifiers, settings)

**Rejected alternatives:**
- RDS PostgreSQL: Không pause được, chi phí cao hơn cho non-prod
- DynamoDB: Không phù hợp cho relational queries phức tạp (orders + items + payments)
- PlanetScale/Neon: Vendor risk, không AWS-native

### Cache/Pub-Sub: ElastiCache Redis

Chọn **ElastiCache Redis** vì:
- Session caching, rate limiting, real-time pub/sub trong một service
- Redis Pub/Sub phù hợp cho WebSocket order status broadcasting
- Managed service với automatic failover
- Phổ biến trong Go ecosystem (go-redis)

### CDN: CloudFront

Bắt buộc cho:
- Static assets (menu images, frontend bundle): cache TTL 24h
- API edge caching cho public endpoints (menu GET): TTL 5 phút
- WAF tích hợp sẵn
- DDoS protection (AWS Shield Standard free tier)

### Auth: Cognito

Chọn **Cognito** vì:
- Managed, không cần tự build auth
- Hỗ trợ social login (Google, Facebook) cho customer
- MFA cho staff/admin
- JWT tokens compatible với Go middleware
- Tích hợp với API Gateway và ALB

**Rejected alternatives:**
- Auth0: Chi phí cao khi scale, external vendor
- Keycloak: Cần tự host và maintain
- Custom JWT: Security risk, không recommend

### Real-time: API Gateway WebSocket + Redis Pub/Sub

Architecture:
1. Client kết nối WebSocket tới API Gateway WS endpoint
2. `$connect` Lambda lưu `connectionId` vào DynamoDB
3. API server publish event vào Redis channel khi order state thay đổi
4. Lambda subscriber nhận từ Redis → gửi qua API Gateway Management API tới các connection liên quan

Lý do không dùng WebSocket trực tiếp trên ECS:
- ALB WebSocket connections bị giới hạn idle timeout (4000s max)
- API Gateway WS không bị giới hạn, tự scale, không manage connection state

### Infrastructure as Code: AWS CDK (TypeScript)

Chọn **CDK** thay vì Terraform hoặc CloudFormation raw vì:
- TypeScript — familiar với frontend team
- Higher-level abstractions (L2/L3 constructs) giảm boilerplate
- Type-safe, IDE autocomplete
- Native AWS support

**Rejected alternatives:**
- Terraform: HCL syntax ít familiar, provider lag với AWS mới nhất
- CloudFormation YAML: Verbose, không có type safety
- Pulumi: Smaller community, ít examples

## Consequences

### Positive

- Zero-management infrastructure (Fargate, Aurora Serverless, Cognito)
- Auto-scaling xử lý giờ cao điểm tự động
- Strong security baseline từ AWS managed services
- Developer có thể focus vào business logic thay vì infra
- Cost-optimized: staging/dev environments chạy Aurora Serverless pause khi idle

### Negative

- AWS vendor lock-in: Migration sang cloud khác tốn kém
- Learning curve cho CDK nếu team không familiar
- Aurora Serverless v2 có cold start ~1 giây sau pause — không phù hợp cho dev làm việc ban đêm (cần warm-up strategy)
- ElastiCache Redis tối thiểu 2 nodes (~$50/month) ngay cả khi traffic thấp

### Mitigation

- Vendor lock-in: Abstraction layer trong code (repository pattern), infrastructure phần lớn là commodity
- Aurora cold start: Disable pause trên staging, hoặc schedule warm-up Lambda
