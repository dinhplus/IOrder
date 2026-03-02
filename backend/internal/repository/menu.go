package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/lib/pq"
)

type MenuCategory struct {
	ID          string
	TenantID    string
	Name        string
	Type        string
	Description *string
	SortOrder   int
	IsActive    bool
	CreatedAt   time.Time
}

type MenuItem struct {
	ID          string
	TenantID    string
	CategoryID  string
	Name        string
	Description *string
	Price       int64
	ImageURL    *string
	Tags        []string
	IsAvailable bool
	SortOrder   int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type MenuRepository struct {
	db *sql.DB
}

func NewMenuRepository(db *sql.DB) *MenuRepository {
	return &MenuRepository{db: db}
}

func (r *MenuRepository) CreateCategory(ctx context.Context, c *MenuCategory) error {
	return r.db.QueryRowContext(ctx,
		`INSERT INTO menu_categories (tenant_id, name, type, description, sort_order, is_active)
		 VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, created_at`,
		c.TenantID, c.Name, c.Type, c.Description, c.SortOrder, c.IsActive,
	).Scan(&c.ID, &c.CreatedAt)
}

func (r *MenuRepository) GetCategory(ctx context.Context, id, tenantID string) (*MenuCategory, error) {
	c := &MenuCategory{}
	err := r.db.QueryRowContext(ctx,
		`SELECT id, tenant_id, name, type, description, sort_order, is_active, created_at
		 FROM menu_categories WHERE id=$1 AND tenant_id=$2`, id, tenantID,
	).Scan(&c.ID, &c.TenantID, &c.Name, &c.Type, &c.Description, &c.SortOrder, &c.IsActive, &c.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return c, err
}

func (r *MenuRepository) UpdateCategory(ctx context.Context, c *MenuCategory) error {
	res, err := r.db.ExecContext(ctx,
		`UPDATE menu_categories SET name=$1, type=$2, description=$3, sort_order=$4, is_active=$5
		 WHERE id=$6 AND tenant_id=$7`,
		c.Name, c.Type, c.Description, c.SortOrder, c.IsActive, c.ID, c.TenantID)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *MenuRepository) DeleteCategory(ctx context.Context, id, tenantID string) error {
	res, err := r.db.ExecContext(ctx,
		`DELETE FROM menu_categories WHERE id=$1 AND tenant_id=$2`, id, tenantID)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *MenuRepository) ListCategories(ctx context.Context, tenantID string) ([]*MenuCategory, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, tenant_id, name, type, description, sort_order, is_active, created_at
		 FROM menu_categories WHERE tenant_id=$1 ORDER BY sort_order, created_at`, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var cats []*MenuCategory
	for rows.Next() {
		c := &MenuCategory{}
		if err := rows.Scan(&c.ID, &c.TenantID, &c.Name, &c.Type, &c.Description, &c.SortOrder, &c.IsActive, &c.CreatedAt); err != nil {
			return nil, err
		}
		cats = append(cats, c)
	}
	return cats, rows.Err()
}

func (r *MenuRepository) CreateItem(ctx context.Context, item *MenuItem) error {
	tags := item.Tags
	if tags == nil {
		tags = []string{}
	}
	return r.db.QueryRowContext(ctx,
		`INSERT INTO menu_items (tenant_id, category_id, name, description, price, image_url, tags, is_available, sort_order)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id, created_at, updated_at`,
		item.TenantID, item.CategoryID, item.Name, item.Description, item.Price, item.ImageURL, pq.Array(tags), item.IsAvailable, item.SortOrder,
	).Scan(&item.ID, &item.CreatedAt, &item.UpdatedAt)
}

func (r *MenuRepository) GetItem(ctx context.Context, id, tenantID string) (*MenuItem, error) {
	item := &MenuItem{}
	err := r.db.QueryRowContext(ctx,
		`SELECT id, tenant_id, category_id, name, description, price, image_url, tags, is_available, sort_order, created_at, updated_at
		 FROM menu_items WHERE id=$1 AND tenant_id=$2`, id, tenantID,
	).Scan(&item.ID, &item.TenantID, &item.CategoryID, &item.Name, &item.Description, &item.Price, &item.ImageURL, pq.Array(&item.Tags), &item.IsAvailable, &item.SortOrder, &item.CreatedAt, &item.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return item, err
}

func (r *MenuRepository) UpdateItem(ctx context.Context, item *MenuItem) error {
	tags := item.Tags
	if tags == nil {
		tags = []string{}
	}
	res, err := r.db.ExecContext(ctx,
		`UPDATE menu_items SET category_id=$1, name=$2, description=$3, price=$4, image_url=$5, tags=$6, is_available=$7, sort_order=$8, updated_at=NOW()
		 WHERE id=$9 AND tenant_id=$10`,
		item.CategoryID, item.Name, item.Description, item.Price, item.ImageURL, pq.Array(tags), item.IsAvailable, item.SortOrder, item.ID, item.TenantID)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *MenuRepository) DeleteItem(ctx context.Context, id, tenantID string) error {
	res, err := r.db.ExecContext(ctx,
		`DELETE FROM menu_items WHERE id=$1 AND tenant_id=$2`, id, tenantID)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *MenuRepository) ListItems(ctx context.Context, tenantID string) ([]*MenuItem, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, tenant_id, category_id, name, description, price, image_url, tags, is_available, sort_order, created_at, updated_at
		 FROM menu_items WHERE tenant_id=$1 ORDER BY sort_order, created_at`, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*MenuItem
	for rows.Next() {
		item := &MenuItem{}
		if err := rows.Scan(&item.ID, &item.TenantID, &item.CategoryID, &item.Name, &item.Description, &item.Price, &item.ImageURL, pq.Array(&item.Tags), &item.IsAvailable, &item.SortOrder, &item.CreatedAt, &item.UpdatedAt); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *MenuRepository) ToggleItemAvailability(ctx context.Context, id, tenantID string, available bool) error {
	res, err := r.db.ExecContext(ctx,
		`UPDATE menu_items SET is_available=$1, updated_at=NOW() WHERE id=$2 AND tenant_id=$3`,
		available, id, tenantID)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return sql.ErrNoRows
	}
	return nil
}
