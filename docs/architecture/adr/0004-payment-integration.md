# ADR 0004: Payment Integration Strategy

## Status

Accepted

## Date

2026-03-02

## Context

IOrder needs to support the payment methods most popular in the Vietnamese and international markets. Requirements:

1. Support VietQR (Vietnamese inter-bank standard), MoMo, ZaloPay, and ShopeePay
2. Support international cards (Visa/Mastercard) for upscale restaurants and international guests
3. Webhook-based confirmation (async) — avoid continuous polling
4. Idempotent webhook handling
5. PCI DSS compliance (no card data stored on our servers)
6. Easy to add new providers in the future

## Decision

### Provider Selection

| Provider | Market | Use Case |
|---|---|---|
| **VietQR** (open standard) | Vietnam | Inter-bank QR code; free to generate |
| **MoMo** | Vietnam | Most popular digital wallet in Vietnam |
| **ZaloPay** | Vietnam | Digital wallet integrated into Zalo ecosystem |
| **ShopeePay** (Airpay) | SEA | Digital wallet popular among younger users |
| **Adyen** | International | Visa/MC/JCB; Apple Pay, Google Pay |

**Rejected:**
- Stripe: No native Vietnam merchant support as of 2026
- VNPay: Complex onboarding; higher transaction fees
- PayOS: Suitable for individuals/small startups; insufficient enterprise features

### Architecture Pattern: Provider Abstraction

```go
// domain/payment/provider.go

type PaymentProvider interface {
    // CreatePayment initiates a payment and returns a checkout URL or QR data
    CreatePayment(ctx context.Context, req CreatePaymentRequest) (*PaymentResult, error)
    // VerifyWebhook validates the webhook signature and returns the parsed event
    VerifyWebhook(ctx context.Context, payload []byte, headers map[string]string) (*WebhookEvent, error)
    // GetPaymentStatus queries the current payment status (polling fallback)
    GetPaymentStatus(ctx context.Context, providerRef string) (*PaymentStatus, error)
}

type CreatePaymentRequest struct {
    OrderID     string
    Amount      int64    // in smallest currency unit (xu/cent)
    Currency    string
    Description string
    ReturnURL   string
    WebhookURL  string
    ExtraData   map[string]string
}

type PaymentResult struct {
    ProviderRef    string    // external transaction ID
    CheckoutURL    string    // redirect URL (for MoMo/ZaloPay deep link)
    QRData         string    // VietQR string or MoMo QR
    QRImageURL     string    // rendered QR image URL (optional)
    ExpiresAt      time.Time
}
```

### VietQR Implementation

VietQR is an open standard; no API key is required to generate the QR string:

```go
// Build a VietQR string following the NAPAS/VIETQR EMV QR Code standard (ISO 20022)
func BuildVietQRString(bankBIN, accountNo string, amount int64, description string) string {
    // Format per the EMV QR Code standard
    // No API key needed; generated directly in the backend
}
```

To confirm VietQR payments:
- Integrate with **Casso** or **PayOS** to receive a webhook when a transaction arrives in the bank account
- Or use a bank's open banking API directly (MB Bank, VCB, ACB all offer this)

### MoMo Integration

```go
// infrastructure/payment/momo/provider.go
// Uses MoMo Payment Gateway API v2
// Docs: https://developers.momo.vn/
```

Flow:
1. POST `https://payment.momo.vn/v2/gateway/api/create`
2. Receive `payUrl` or `qrCodeUrl`
3. Customer pays in MoMo app
4. MoMo POSTs a webhook to `POST /api/v1/payments/webhooks/momo`
5. Verify `signature` = HMAC-SHA256(rawBody, secretKey)
6. Update payment status

### ZaloPay Integration

Similar flow to MoMo, using the ZaloPay Merchant API:
- Endpoint: `https://sb-openapi.zalopay.vn/v2/create`
- Webhook verification: HMAC-SHA256 with `key2` (ZaloPay-specific)

### Adyen Integration

Uses **Adyen Drop-in** (hosted payment UI) for international cards:
- Adyen handles PCI DSS scope; no card data is stored on IOrder servers
- Webhook: Adyen Notification Service with HMAC
- Apple Pay and Google Pay are supported automatically

### Webhook Security

All webhook endpoints must:
1. Verify the HMAC signature before any processing
2. Return HTTP 200 immediately (processing happens asynchronously via SQS)
3. Be idempotent: check `provider_ref` to avoid duplicate processing
4. Retry: if the SQS consumer fails, the message is retried automatically (max 3 attempts)

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

In case a webhook does not arrive (network issue):
- A background Lambda runs every 2 minutes and queries the `payments` table for orders with status `pending` for more than 5 minutes
- Calls `GetPaymentStatus()` from the provider
- Updates status if already paid

### Refund Policy

- Refunds are out of scope for v1.0
- Staff can mark a payment as "manual refund" in the dashboard
- Full refund tracking will be in `payment_refunds` table (v2.0)

## Consequences

### Positive

- Provider abstraction allows adding new providers without changing business logic
- Webhook + async processing does not block the API response
- Idempotent handling prevents duplicate processing
- PCI DSS compliant because card data is never stored (delegated to Adyen)

### Negative

- Multiple webhook endpoints and HMAC keys to maintain
- Each provider requires integration testing against sandbox environments (MoMo and ZaloPay both have sandboxes)
- VietQR confirmation depends on Casso/bank API — an additional external dependency

### Mitigation

- Write integration tests against each provider's sandbox
- Mock `PaymentProvider` interface in unit tests
- Monitoring: alert if payment completion rate drops below 95% over 5 minutes
