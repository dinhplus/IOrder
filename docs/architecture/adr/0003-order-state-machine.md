# ADR 0003: Order State Machine Design

## Status

Accepted

## Date

2026-03-02

## Context

Đơn hàng trong nhà hàng có vòng đời phức tạp với nhiều actor (customer, staff, kitchen, payment system) thay đổi trạng thái. Cần thiết kế state machine rõ ràng để:

1. Ngăn chặn transition không hợp lệ (ví dụ: nhảy từ DRAFT sang PAID)
2. Đảm bảo audit trail đầy đủ (ai thay đổi gì, khi nào)
3. Real-time notification đúng đối tượng khi state thay đổi
4. Hỗ trợ nhiều loại đơn (gọi thêm món, tách bill)

## Decision

### State Diagram

```
                   ┌──────────────────────────────────────────────┐
                   │                                              │
                   ▼                                              │
  ┌────────┐   submit   ┌───────────┐  confirm  ┌───────────┐   │
  │ DRAFT  │──────────→ │ SUBMITTED │──────────→ │ CONFIRMED │   │
  └────────┘            └───────────┘            └─────┬─────┘   │
      ▲                       │                        │          │
      │                    reject                  start_prep     │
      │                       │                        │          │
      │                       ▼                        ▼          │
      │                 ┌──────────┐          ┌──────────────┐   │
      │                 │ REJECTED │          │IN_PREPARATION│   │
      │                 └──────────┘          └──────┬───────┘   │
      │                                              │            │
      │                                           mark_ready      │
      │                                              │            │
      │                                              ▼            │
      │                                         ┌─────────┐      │
      │                                         │  READY  │      │
      │                                         └────┬────┘      │
      │                                              │            │
      │                                           serve           │
      │                                              │            │
      │                                              ▼            │
      │                              ┌───────────────────────┐   │
      │   add_items (re-open DRAFT)  │       SERVED          │───┘
      │ ←────────────────────────────┤  (customer can order  │
      │                              │   more items)         │
      │                              └──────────┬────────────┘
      │                                         │
      │                                   request_payment
      │                                         │
      │                                         ▼
      │                              ┌──────────────────────┐
      │                              │  PAYMENT_REQUESTED   │
      │                              └──────────┬───────────┘
      │                                         │
      │                              payment_completed / payment_failed
      │                                         │
      │                              ┌──────────▼───────────┐
      │                              │        PAID          │
      │                              └──────────┬───────────┘
      │                                         │
      │                                      close
      │                                         │
      │                                         ▼
      │                              ┌──────────────────────┐
      │                              │       CLOSED         │
      │                              └──────────────────────┘
      │
  CANCELLED ←── cancel (from DRAFT, SUBMITTED, or CONFIRMED only)
```

### Valid Transitions

| From | To | Actor | Trigger |
|---|---|---|---|
| `DRAFT` | `SUBMITTED` | Customer | Submit order button |
| `SUBMITTED` | `CONFIRMED` | Staff / Auto | Staff confirms, or auto-confirm enabled |
| `SUBMITTED` | `REJECTED` | Staff | Staff rejects (out of stock, etc.) |
| `CONFIRMED` | `IN_PREPARATION` | Kitchen Staff | KDS: start cooking |
| `IN_PREPARATION` | `READY` | Kitchen Staff | KDS: mark done |
| `READY` | `SERVED` | Staff | Mark as delivered to table |
| `SERVED` | `DRAFT` | Customer | Add more items (new order round) |
| `SERVED` | `PAYMENT_REQUESTED` | Customer/Staff | Request bill |
| `PAYMENT_REQUESTED` | `PAID` | Payment System | Webhook success |
| `PAYMENT_REQUESTED` | `SERVED` | System | Payment timeout → back to SERVED |
| `PAID` | `CLOSED` | System | Auto-close after 15 minutes |
| `DRAFT` | `CANCELLED` | Customer/Staff | Cancel before submit |
| `SUBMITTED` | `CANCELLED` | Staff | Cancel after submit |
| `CONFIRMED` | `CANCELLED` | Staff | Cancel before preparation |

### Implementation in Go

```go
// domain/order/state_machine.go

type OrderStatus string

const (
    StatusDraft             OrderStatus = "DRAFT"
    StatusSubmitted         OrderStatus = "SUBMITTED"
    StatusConfirmed         OrderStatus = "CONFIRMED"
    StatusRejected          OrderStatus = "REJECTED"
    StatusInPreparation     OrderStatus = "IN_PREPARATION"
    StatusReady             OrderStatus = "READY"
    StatusServed            OrderStatus = "SERVED"
    StatusPaymentRequested  OrderStatus = "PAYMENT_REQUESTED"
    StatusPaid              OrderStatus = "PAID"
    StatusClosed            OrderStatus = "CLOSED"
    StatusCancelled         OrderStatus = "CANCELLED"
)

// validTransitions defines allowed state changes
var validTransitions = map[OrderStatus][]OrderStatus{
    StatusDraft:            {StatusSubmitted, StatusCancelled},
    StatusSubmitted:        {StatusConfirmed, StatusRejected, StatusCancelled},
    StatusConfirmed:        {StatusInPreparation, StatusCancelled},
    StatusInPreparation:    {StatusReady},
    StatusReady:            {StatusServed},
    StatusServed:           {StatusDraft, StatusPaymentRequested},
    StatusPaymentRequested: {StatusPaid, StatusServed},
    StatusPaid:             {StatusClosed},
}

// CanTransition checks if a transition is allowed
func CanTransition(from, to OrderStatus) bool {
    allowed, ok := validTransitions[from]
    if !ok {
        return false
    }
    for _, s := range allowed {
        if s == to {
            return true
        }
    }
    return false
}
```

### Session-Based Multi-Round Ordering

Khi khách muốn gọi thêm món sau khi đơn đầu đã SERVED:
- Tạo đơn mới với cùng `session_id` và `table_id`
- Đơn mới bắt đầu từ DRAFT
- Khi thanh toán: tổng hợp tất cả đơn cùng `session_id` thành một bill

### Event Publishing

Mỗi state transition publish một event vào EventBridge:

```json
{
  "source": "iorder.orders",
  "detail-type": "OrderStatusChanged",
  "detail": {
    "orderId": "uuid",
    "tenantId": "uuid",
    "tableId": "uuid",
    "sessionId": "uuid",
    "fromStatus": "SUBMITTED",
    "toStatus": "CONFIRMED",
    "actorId": "uuid",
    "actorType": "staff",
    "timestamp": "2026-03-02T10:00:00Z"
  }
}
```

EventBridge rules:
- `OrderStatusChanged` → SQS → Notification Worker (Push/SMS to customer)
- `OrderStatusChanged` where `toStatus=CONFIRMED` → SQS → KDS Worker
- `OrderStatusChanged` where `toStatus=PAID` → SQS → Membership Points Worker

### Audit Log

Tất cả transitions được ghi vào bảng `order_events` (immutable append-only log).

## Consequences

### Positive

- Trạng thái đơn hàng luôn nhất quán, không thể có transition bất hợp lệ
- Audit trail đầy đủ cho mọi thay đổi trạng thái
- Event-driven architecture giúp decoupling giữa các concerns
- Dễ thêm trạng thái mới trong tương lai (ví dụ: REFUNDED)
- Clear actor responsibility cho mỗi transition

### Negative

- Phức tạp hơn simple CRUD approach
- Cần đồng bộ state machine giữa backend logic và frontend UI
- Thêm bảng `order_events` và overhead ghi log

### Mitigation

- Viết unit tests toàn bộ valid/invalid transitions
- Document state diagram rõ ràng (file này)
- Share TypeScript types cho FE/Mobile để tránh inconsistency
