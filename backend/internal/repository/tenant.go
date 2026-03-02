package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"time"
)

type Tenant struct {
	ID        string
	Slug      string
	Name      string
	LogoURL   *string
	Timezone  string
	Currency  string
	IsActive  bool
	Settings  json.RawMessage
	CreatedAt time.Time
	UpdatedAt time.Time
}

type TenantRepository struct {
	db *sql.DB
}

func NewTenantRepository(db *sql.DB) *TenantRepository {
	return &TenantRepository{db: db}
}

func (r *TenantRepository) Create(ctx context.Context, t *Tenant) error {
	settings := t.Settings
	if settings == nil {
		settings = json.RawMessage("{}")
	}
	return r.db.QueryRowContext(ctx,
		`INSERT INTO tenants (slug, name, logo_url, timezone, currency, is_active, settings)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)
		 RETURNING id, created_at, updated_at`,
		t.Slug, t.Name, t.LogoURL, t.Timezone, t.Currency, t.IsActive, settings,
	).Scan(&t.ID, &t.CreatedAt, &t.UpdatedAt)
}

func (r *TenantRepository) GetByID(ctx context.Context, id string) (*Tenant, error) {
	t := &Tenant{}
	err := r.db.QueryRowContext(ctx,
		`SELECT id, slug, name, logo_url, timezone, currency, is_active, settings, created_at, updated_at
		 FROM tenants WHERE id = $1`, id,
	).Scan(&t.ID, &t.Slug, &t.Name, &t.LogoURL, &t.Timezone, &t.Currency, &t.IsActive, &t.Settings, &t.CreatedAt, &t.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return t, err
}

func (r *TenantRepository) GetBySlug(ctx context.Context, slug string) (*Tenant, error) {
	t := &Tenant{}
	err := r.db.QueryRowContext(ctx,
		`SELECT id, slug, name, logo_url, timezone, currency, is_active, settings, created_at, updated_at
		 FROM tenants WHERE slug = $1`, slug,
	).Scan(&t.ID, &t.Slug, &t.Name, &t.LogoURL, &t.Timezone, &t.Currency, &t.IsActive, &t.Settings, &t.CreatedAt, &t.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return t, err
}

func (r *TenantRepository) Update(ctx context.Context, t *Tenant) error {
	return r.db.QueryRowContext(ctx,
		`UPDATE tenants SET slug=$1, name=$2, logo_url=$3, timezone=$4, currency=$5, is_active=$6, settings=$7, updated_at=NOW()
		 WHERE id=$8 RETURNING updated_at`,
		t.Slug, t.Name, t.LogoURL, t.Timezone, t.Currency, t.IsActive, t.Settings, t.ID,
	).Scan(&t.UpdatedAt)
}

func (r *TenantRepository) List(ctx context.Context) ([]*Tenant, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, slug, name, logo_url, timezone, currency, is_active, settings, created_at, updated_at FROM tenants ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var tenants []*Tenant
	for rows.Next() {
		t := &Tenant{}
		if err := rows.Scan(&t.ID, &t.Slug, &t.Name, &t.LogoURL, &t.Timezone, &t.Currency, &t.IsActive, &t.Settings, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, err
		}
		tenants = append(tenants, t)
	}
	return tenants, rows.Err()
}
