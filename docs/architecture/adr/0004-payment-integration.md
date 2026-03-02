# ADR 0004: Payment Integration Strategy

## Status

Accepted

## Date

2026-03-02

## Context

IOrder cần hỗ trợ nhiều phương thức thanh toán phổ biến tại thị trường Việt Nam và quốc tế. Yêu cầu:

1. Hỗ trợ VietQR (chuẩn ngân hàng Việt Nam), MoMo, ZaloPay, ShopeePay
2. Hỗ trợ thẻ quốc tế (Visa/Mastercard) cho nhà hàng cao cấp / khách quốc tế
3. Webhook-based confirmation (async) — không poll API liên tục
4. Idempotent webhook handling
5. PCI DSS compliance (không lưu thông tin thẻ)
6. Dễ thêm provider mới trong tương lai

## Decision

### Provider Selection

| Provider | Market | Use Case |
|---|---|---|
| **VietQR** (via VietQR open standard) | Việt Nam | QR code liên ngân hàng, miễn phí |
| **MoMo** | Việt Nam | Ví điện tử phổ biến nhất VN |
| **ZaloPay** | Việt Nam | Ví điện tử tích hợp Zalo ecosystem |
| **ShopeePay** (Airpay) | SEA | Ví điện tử Shopee, phổ biến giới trẻ |
| **Adyen** | International | Thẻ Visa/MC/JCB, Apple Pay, Google Pay |

**Rejected:**
- Stripe: Chưa hỗ trợ Vietnam merchant natively (2026)
- VNPay: Phức tạp onboarding, phí transaction cao hơn
- PayOS: Phù hợp cho cá nhân/startup nhỏ, ít tính năng enterprise

### Architecture Pattern: Provider Abstraction

```go
// domain/payment/provider.go

type PaymentProvider interface {
    // CreatePayment initiates a payment and returns checkout URL or QR data
    CreatePayment(ctx context.Context, req CreatePaymentRequest) (*PaymentResult, error)
    // VerifyWebhook validates the webhook signature and returns parsed event
    VerifyWebhook(ctx context.Context, payload []byte, headers map[string]string) (*WebhookEvent, error)
    // GetPaymentStatus queries current payment status (for polling fallback)
    GetPaymentStatus(ctx context.Context, providerRef string) (*PaymentStatus, error)
}

type CreatePaymentRequest struct {
    OrderID     string
    Amount      int64    // in smallest currency unit (xu/satoshi)
    Currency    string
    Description string
    ReturnURL   string
    WebhookURL  string
    ExtraData   map[string]string
}

type PaymentResult struct {
    ProviderRef    string  // external transaction ID
    CheckoutURL    string  // redirect URL (for MoMo/ZaloPay deep link)
    QRData         string  // VietQR QR string or MoMo QR
    QRImageURL     string  // rendered QR image URL (optional)
    ExpiresAt      time.Time
}
```

### VietQR Implementation

VietQR là chuẩn mở, không cần tích hợp qua một nhà cung cấp cụ thể:

```go
// Tạo VietQR string theo chuẩn NAPAS/VIETQR
func BuildVietQRString(bankBIN, accountNo, amount int64, description string) string {
    // Format theo chuẩn EMV QR Code (ISO 20022)
    // Không cần API key, tạo trực tiếp
}
```

Để xác nhận thanh toán VietQR:
- Integrate với **Casso** hoặc **Payos** để nhận webhook khi có giao dịch vào tài khoản
- Hoặc bank API trực tiếp (MB Bank, VCB, ACB có open banking)

### MoMo Integration

```go
// infrastructure/payment/momo/provider.go
// Sử dụng MoMo Payment Gateway API v2
// Doc: https://developers.momo.vn/
```

Flow:
1. POST `https://payment.momo.vn/v2/gateway/api/create`
2. Nhận `payUrl` hoặc `qrCodeUrl`
3. Customer thanh toán
4. MoMo POST webhook tới `POST /api/v1/payments/webhooks/momo`
5. Verify `signature` = HMAC-SHA256(rawBody, secretKey)
6. Update payment status

### ZaloPay Integration

Flow tương tự MoMo, sử dụng ZaloPay Merchant API:
- Endpoint: `https://sb-openapi.zalopay.vn/v2/create`
- Webhook verification: HMAC-SHA256 với `key2` (ZaloPay specific)

### Adyen Integration

Sử dụng **Adyen Drop-in** (hosted payment UI) cho thẻ quốc tế:
- Adyen handle PCI DSS scope, không cần lưu card data
- Webhook: Adyen Notification Service với HMAC
- Hỗ trợ Apple Pay, Google Pay automatically

### Webhook Security

Tất cả webhook endpoint phải:
1. Verify HMAC signature trước khi xử lý
2. Return 200 ngay (processing async qua SQS)
3. Idempotent: check `provider_ref` đã xử lý chưa
4. Retry: nếu SQS consumer fail, message được retry tự động (max 3 lần)

```go
func (h *PaymentHandler) HandleMoMoWebhook(c *gin.Context) {
    body, err := io.ReadAll(c.Request.Body)
    if err != nil {
        handler.RespondError(c, http.StatusBadRequest, handler.ErrValidation, "invalid body")
        return
    }

    // 1. Verify HMAC before any processing
    if !h.momoProvider.VerifyHMAC(body, c.GetHeader("X-Signature")) {
        handler.RespondError(c, http.StatusUnauthorized, handler.ErrUnauthorized, "invalid signature")
        return
    }

    // 2. Return 200 immediately
    c.JSON(http.StatusOK, gin.H{"status": "received"})

    // 3. Enqueue for async processing
    h.queue.Enqueue(ctx, "payment.webhook.momo", body)
}
```

### Payment Status Polling Fallback

Trong trường hợp webhook không đến (network issue):
- Background Lambda chạy mỗi 2 phút query `payments` table cho orders `pending` > 5 phút
- Gọi `GetPaymentStatus()` từ provider
- Cập nhật trạng thái nếu đã paid

### Refund Policy

- Refund không nằm trong scope v1.0
- Staff có thể đánh dấu "manual refund" trong dashboard
- Tracking trong `payment_refunds` table (v2.0)

## Consequences

### Positive

- Provider abstraction cho phép thêm provider mới không cần thay đổi business logic
- Webhook + async processing không block API response
- Idempotent handling ngăn duplicate processing
- PCI DSS compliant vì không lưu card data (delegated to Adyen)

### Negative

- Cần maintain nhiều webhook endpoint và HMAC keys
- Test integration với các sandbox environment (MoMo, ZaloPay đều có sandbox)
- VietQR confirmation phụ thuộc vào Casso/bank API — thêm external dependency

### Mitigation

- Viết integration tests với sandbox của từng provider
- Mock `PaymentProvider` interface trong unit tests
- Monitoring: alert nếu payment completion rate < 95% trong 5 phút
