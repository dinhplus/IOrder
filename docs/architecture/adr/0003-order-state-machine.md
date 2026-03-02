# ADR 0003: Order State Machine Design

## Status

Accepted

## Date

2026-03-02

## Context

Restaurant orders have a complex lifecycle with multiple actors (customer, waiter, chef, bartender, server, cashier, payment system) changing the order status. A well-defined state machine is necessary to:

1. Prevent invalid transitions (e.g., jumping directly from DRAFT to PAID)
2. Guarantee a complete audit trail (who changed what, when)
3. Send real-time notifications to the correct actor when state changes
4. Support multiple order rounds per table session and split-bill scenarios

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
| `DRAFT` | `SUBMITTED` | Customer / Waiter | Submit order button |
| `SUBMITTED` | `CONFIRMED` | Waiter / Auto | Waiter confirms, or auto-confirm is enabled |
| `SUBMITTED` | `REJECTED` | Waiter / Manager | Item out of stock or other issue |
| `CONFIRMED` | `IN_PREPARATION` | Chef / Bartender | KDS/BDS: start cooking/preparing |
| `IN_PREPARATION` | `READY` | Chef / Bartender | KDS/BDS: mark done |
| `READY` | `SERVED` | Server / Waiter | Delivered to table |
| `SERVED` | `DRAFT` | Customer / Waiter | Add more items (new order round) |
| `SERVED` | `PAYMENT_REQUESTED` | Customer / Waiter | Request bill |
| `PAYMENT_REQUESTED` | `PAID` | Payment System | Webhook success |
| `PAYMENT_REQUESTED` | `SERVED` | System | Payment timeout → back to SERVED |
| `PAID` | `CLOSED` | System | Auto-close after 15 minutes |
| `DRAFT` | `CANCELLED` | Customer / Waiter | Cancel before submit |
| `SUBMITTED` | `CANCELLED` | Waiter / Manager | Cancel after submit |
| `CONFIRMED` | `CANCELLED` | Waiter / Manager | Cancel before preparation |

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

### KDS / BDS Item Routing

When an order is CONFIRMED, its items are routed based on the category type:
- `category.type = "food"` → routed to KDS (Kitchen Display System)
- `category.type = "beverage"` → routed to BDS (Bar Display System)

Each display only shows the relevant items, allowing kitchen and bar to work independently and in parallel.

### Session-Based Multi-Round Ordering

When a customer wants to order additional items after their first order has been SERVED:
- A new order is created with the same `session_id` and `table_id`
- The new order starts from DRAFT
- At checkout, all orders sharing the same `session_id` are consolidated into a single bill

### Event Publishing

Every state transition publishes an event to EventBridge:

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
- `OrderStatusChanged` where `toStatus=CONFIRMED` → SQS → KDS/BDS Worker
- `OrderStatusChanged` where `toStatus=PAID` → SQS → Membership Points Worker + Inventory Deduction Worker

### Audit Log

All transitions are written to the `order_events` table (immutable append-only log).

## Consequences

### Positive

- Order status is always consistent; invalid transitions are impossible
- Complete audit trail for every status change
- Event-driven architecture decouples notifications, KDS, inventory, and membership concerns
- Easy to add new states in the future (e.g., REFUNDED)
- Clear actor responsibility for every transition

### Negative

- More complex than a simple CRUD approach
- State machine must be kept in sync between backend logic and frontend/mobile UI
- Additional `order_events` table adds write overhead

### Mitigation

- Write unit tests covering all valid and invalid transitions
- Document state diagram clearly (this file)
- Share TypeScript enum types between FE and mobile to avoid inconsistency
