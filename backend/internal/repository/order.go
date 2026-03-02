package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"time"
)

const (
	OrderStatusDraft            = "DRAFT"
	OrderStatusSubmitted        = "SUBMITTED"
	OrderStatusConfirmed        = "CONFIRMED"
	OrderStatusInPreparation    = "IN_PREPARATION"
	OrderStatusReady            = "READY"
	OrderStatusServed           = "SERVED"
	OrderStatusPaymentRequested = "PAYMENT_REQUESTED"
	OrderStatusPaid             = "PAID"
	OrderStatusCancelled        = "CANCELLED"
	OrderStatusRejected         = "REJECTED"
)

type Order struct {
	ID             string
	TenantID       string
	TableID        string
	SessionID      string
	Status         string
	CustomerID     *string
	Subtotal       int64
	DiscountAmount int64
	Total          int64
	Notes          *string
	PlacedAt       *time.Time
	ConfirmedAt    *time.Time
	ReadyAt        *time.Time
	ServedAt       *time.Time
	PaidAt         *time.Time
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type OrderItem struct {
	ID        string
	OrderID   string
	TenantID  string
	ItemID    string
	ItemName  string
	ItemPrice int64
	Quantity  int
	Modifiers json.RawMessage
	Notes     *string
	Status    string
	RoutedTo  string
}

type OrderEvent struct {
	ID         string
	OrderID    string
	TenantID   string
	FromStatus *string
	ToStatus   string
	ActorID    *string
	ActorType  *string
	Metadata   json.RawMessage
	CreatedAt  time.Time
}

type OrderRepository struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) Create(ctx context.Context, o *Order) error {
	return r.db.QueryRowContext(ctx,
		`INSERT INTO orders (tenant_id, table_id, session_id, status, customer_id, notes)
		 VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, created_at, updated_at`,
		o.TenantID, o.TableID, o.SessionID, o.Status, o.CustomerID, o.Notes,
	).Scan(&o.ID, &o.CreatedAt, &o.UpdatedAt)
}

func (r *OrderRepository) GetByID(ctx context.Context, id string) (*Order, error) {
	o := &Order{}
	err := r.db.QueryRowContext(ctx,
		`SELECT id, tenant_id, table_id, session_id, status, customer_id, subtotal, discount_amount, total, notes,
		        placed_at, confirmed_at, ready_at, served_at, paid_at, created_at, updated_at
		 FROM orders WHERE id=$1`, id,
	).Scan(&o.ID, &o.TenantID, &o.TableID, &o.SessionID, &o.Status, &o.CustomerID,
		&o.Subtotal, &o.DiscountAmount, &o.Total, &o.Notes,
		&o.PlacedAt, &o.ConfirmedAt, &o.ReadyAt, &o.ServedAt, &o.PaidAt,
		&o.CreatedAt, &o.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return o, err
}

func (r *OrderRepository) UpdateStatus(ctx context.Context, id, tenantID, status string) error {
	res, err := r.db.ExecContext(ctx,
		`UPDATE orders SET status=$1, updated_at=NOW() WHERE id=$2 AND tenant_id=$3`,
		status, id, tenantID)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *OrderRepository) ListByTenant(ctx context.Context, tenantID string) ([]*Order, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, tenant_id, table_id, session_id, status, customer_id, subtotal, discount_amount, total, notes,
		        placed_at, confirmed_at, ready_at, served_at, paid_at, created_at, updated_at
		 FROM orders WHERE tenant_id=$1 ORDER BY created_at DESC`, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var orders []*Order
	for rows.Next() {
		o := &Order{}
		if err := rows.Scan(&o.ID, &o.TenantID, &o.TableID, &o.SessionID, &o.Status, &o.CustomerID,
			&o.Subtotal, &o.DiscountAmount, &o.Total, &o.Notes,
			&o.PlacedAt, &o.ConfirmedAt, &o.ReadyAt, &o.ServedAt, &o.PaidAt,
			&o.CreatedAt, &o.UpdatedAt); err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}
	return orders, rows.Err()
}

func (r *OrderRepository) GetItems(ctx context.Context, orderID string) ([]*OrderItem, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, order_id, tenant_id, item_id, item_name, item_price, quantity, modifiers, notes, status, routed_to
		 FROM order_items WHERE order_id=$1`, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*OrderItem
	for rows.Next() {
		item := &OrderItem{}
		if err := rows.Scan(&item.ID, &item.OrderID, &item.TenantID, &item.ItemID, &item.ItemName,
			&item.ItemPrice, &item.Quantity, &item.Modifiers, &item.Notes, &item.Status, &item.RoutedTo); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *OrderRepository) UpdateItems(ctx context.Context, orderID, tenantID string, items []*OrderItem) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback() //nolint:errcheck

	if _, err := tx.ExecContext(ctx, `DELETE FROM order_items WHERE order_id=$1`, orderID); err != nil {
		return err
	}
	for _, item := range items {
		mods := item.Modifiers
		if mods == nil {
			mods = json.RawMessage("[]")
		}
		if _, err := tx.ExecContext(ctx,
			`INSERT INTO order_items (order_id, tenant_id, item_id, item_name, item_price, quantity, modifiers, notes, status, routed_to)
			 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
			orderID, tenantID, item.ItemID, item.ItemName, item.ItemPrice, item.Quantity, mods, item.Notes, item.Status, item.RoutedTo,
		); err != nil {
			return err
		}
	}

	if _, err := tx.ExecContext(ctx, `UPDATE orders SET updated_at=NOW() WHERE id=$1`, orderID); err != nil {
		return err
	}
	return tx.Commit()
}

func (r *OrderRepository) UpdateItem(ctx context.Context, item *OrderItem) error {
	res, err := r.db.ExecContext(ctx,
		`UPDATE order_items SET quantity=$1, notes=$2, status=$3 WHERE id=$4 AND order_id=$5`,
		item.Quantity, item.Notes, item.Status, item.ID, item.OrderID)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *OrderRepository) CreateEvent(ctx context.Context, e *OrderEvent) error {
	meta := e.Metadata
	if meta == nil {
		meta = json.RawMessage("{}")
	}
	return r.db.QueryRowContext(ctx,
		`INSERT INTO order_events (order_id, tenant_id, from_status, to_status, actor_id, actor_type, metadata)
		 VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id, created_at`,
		e.OrderID, e.TenantID, e.FromStatus, e.ToStatus, e.ActorID, e.ActorType, meta,
	).Scan(&e.ID, &e.CreatedAt)
}
