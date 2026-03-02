# ADR 0005: QR Code Table Ordering

## Status

Accepted

## Date

2026-03-02

## Context

Hệ thống đặt món qua QR code là tính năng cốt lõi của IOrder. Cần thiết kế QR flow đảm bảo:

1. **Không cần cài app** — khách chỉ cần camera điện thoại
2. **Bảo mật** — QR không bị giả mạo, không chia sẻ được giữa các bàn
3. **UX mượt** — thời gian từ scan → menu < 2 giây
4. **Multi-session** — nhiều người tại cùng bàn có thể xem và gọi cùng lúc
5. **QR expiry** — QR hết hạn có thể được regenerate bởi staff

## Decision

### QR URL Structure

```
https://order.iorder.vn/{restaurant_slug}/table/{table_id}?token={jwt}
```

Ví dụ:
```
https://order.iorder.vn/pho-ha-noi/table/tbl_abc123?token=eyJhbGciOiJIUzI1...
```

### JWT Token Design

QR token là **JWT HS256** signed với `TABLE_QR_SECRET` (per-restaurant secret, lưu trong AWS Secrets Manager):

```json
// Header
{
  "alg": "HS256",
  "typ": "JWT"
}

// Payload
{
  "sub": "table:{table_id}",
  "tid": "{tenant_id}",           // tenant ID
  "tbl": "{table_id}",            // table ID
  "iat": 1709340000,              // issued at
  "exp": 1709361600,              // expires in 6 hours
  "jti": "{uuid}"                 // unique ID (prevent replay)
}
```

**Thời gian sống: 6 giờ** — đủ cho một bữa tối dài nhất.

### Token Validation Flow

```
Customer scans QR → Browser opens URL
         │
         ▼
Frontend sends GET /api/v1/qr/validate?token={jwt}
         │
         ▼
Backend validates:
  1. JWT signature (HMAC-SHA256 with restaurant secret)
  2. Token not expired (exp > now)
  3. Table exists and is active
  4. Restaurant is currently open (operating hours)
         │
    ┌────┴────┐
   VALID    INVALID
    │           │
    ▼           ▼
Return      Return 401 with
session     "QR expired, please
info        ask staff to refresh"
    │
    ▼
Create/Resume table session
Return: { sessionId, tableId, tenantId, menuUrl }
```

### Session Management

Một **table session** đại diện cho một lượt phục vụ tại bàn (từ khi khách ngồi đến khi thanh toán):

```go
type TableSession struct {
    ID        string    // session UUID
    TableID   string
    TenantID  string
    Status    string    // active, closed
    OpenedAt  time.Time
    ClosedAt  *time.Time
}
```

- Khi QR hợp lệ → tìm hoặc tạo active session cho bàn đó
- Nhiều thiết bị scan cùng 1 QR → cùng session → thấy đơn của nhau (real-time sync)
- Staff reset bàn → session cũ closed, session mới tạo khi QR được scan tiếp theo

### Multi-Device Ordering

Trong cùng một bữa ăn, nhiều người tại cùng bàn có thể gọi mon:
- Tất cả thiết bị join cùng `sessionId` WebSocket channel
- Giỏ hàng được sync real-time giữa các thiết bị
- Khi 1 người thêm món → tất cả người còn lại thấy cập nhật
- Có tùy chọn: "Shared cart" (cùng 1 giỏ) hoặc "Individual cart" (mỗi người 1 giỏ, tách bill)

### QR Regeneration

Staff có thể regenerate QR khi:
- QR cũ bị hết hạn (> 6 giờ)
- Bàn được reset sau khi khách rời
- QR bị phát tán ra ngoài (security concern)

```
POST /api/v1/tables/{tableId}/qr
→ Generate new JWT, update table.qr_token
→ Return new QR image (base64 PNG) and download URL
→ Invalidate old token (by checking jti against DynamoDB blocklist)
```

### QR Generation & Rendering

- QR được render server-side bằng Go library `github.com/skip2/go-qrcode`
- Output formats: PNG (300x300), PDF A4 (với branding nhà hàng, 6 QR per page)
- Error correction level: **M** (15% data recovery) — đủ cho nhà hàng có nhãn sticker
- Lưu QR image vào S3 với URL có signed expiry

### Security Considerations

1. **Token không reusable sau reset**: Khi staff reset bàn, generate `jti` mới → token cũ invalid
2. **Rate limiting**: Validate endpoint giới hạn 10 req/min per IP
3. **HTTPS only**: QR URL phải HTTPS, không HTTP redirect
4. **Table slug không predictable**: `table_id` là UUID, không phải sequential integer
5. **Tenant isolation**: Token chứa `tenant_id`, server verify token đúng tenant

### PWA vs Native App

Chọn **Progressive Web App (PWA)** cho customer ordering flow vì:
- Không cần cài app → giảm friction tối đa
- Camera scan QR → mở browser trực tiếp
- PWA có thể "Add to Home Screen" cho customer loyalty
- Cập nhật instantly (không qua App Store review)

Native app (React Native) vẫn được phát triển cho:
- Staff/KDS interface (cần push notifications reliable)
- Restaurant owner dashboard trên mobile
- Advanced features: offline mode, biometric auth

## Consequences

### Positive

- Zero friction for customers: chỉ cần camera
- Secure: JWT signed, expiry, per-restaurant secret
- Multi-device ordering trong cùng session
- Staff có control rõ ràng (reset, regenerate)
- PWA load nhanh, không cần cài đặt

### Negative

- QR token 6 giờ: nếu server clock drift → token invalid (cần NTP sync)
- Phụ thuộc internet: khách cần kết nối để xem menu (offline mode là v2.0)
- QR in trên sticker có thể bị bẩn/hỏng → cần backup QR dạng laminate

### Mitigation

- NTP sync tự động trên ECS Fargate (AWS managed)
- CDN caching menu data để giảm latency
- Cung cấp QR backup PDF cho nhà hàng tự in thêm
