package repository

import (
	"context"
	"database/sql"
	"time"
)

type FloorPlan struct {
	ID         string
	TenantID   string
	Name       string
	FloorLevel int
	IsActive   bool
	CreatedAt  time.Time
}

type RestaurantTable struct {
	ID          string
	TenantID    string
	FloorPlanID string
	Name        string
	Capacity    int
	PosX        float64
	PosY        float64
	Shape       string
	Status      string
	QRToken     *string
	QRExpiresAt *time.Time
	CreatedAt   time.Time
}

type TableRepository struct {
	db *sql.DB
}

func NewTableRepository(db *sql.DB) *TableRepository {
	return &TableRepository{db: db}
}

func (r *TableRepository) CreateFloorPlan(ctx context.Context, fp *FloorPlan) error {
	return r.db.QueryRowContext(ctx,
		`INSERT INTO floor_plans (tenant_id, name, floor_level, is_active)
		 VALUES ($1, $2, $3, $4) RETURNING id, created_at`,
		fp.TenantID, fp.Name, fp.FloorLevel, fp.IsActive,
	).Scan(&fp.ID, &fp.CreatedAt)
}

func (r *TableRepository) GetFloorPlan(ctx context.Context, id, tenantID string) (*FloorPlan, error) {
	fp := &FloorPlan{}
	err := r.db.QueryRowContext(ctx,
		`SELECT id, tenant_id, name, floor_level, is_active, created_at FROM floor_plans WHERE id=$1 AND tenant_id=$2`,
		id, tenantID,
	).Scan(&fp.ID, &fp.TenantID, &fp.Name, &fp.FloorLevel, &fp.IsActive, &fp.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return fp, err
}

func (r *TableRepository) UpdateFloorPlan(ctx context.Context, fp *FloorPlan) error {
	res, err := r.db.ExecContext(ctx,
		`UPDATE floor_plans SET name=$1, floor_level=$2, is_active=$3 WHERE id=$4 AND tenant_id=$5`,
		fp.Name, fp.FloorLevel, fp.IsActive, fp.ID, fp.TenantID)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *TableRepository) ListFloorPlans(ctx context.Context, tenantID string) ([]*FloorPlan, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, tenant_id, name, floor_level, is_active, created_at FROM floor_plans WHERE tenant_id=$1 ORDER BY floor_level, created_at`,
		tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var fps []*FloorPlan
	for rows.Next() {
		fp := &FloorPlan{}
		if err := rows.Scan(&fp.ID, &fp.TenantID, &fp.Name, &fp.FloorLevel, &fp.IsActive, &fp.CreatedAt); err != nil {
			return nil, err
		}
		fps = append(fps, fp)
	}
	return fps, rows.Err()
}

func (r *TableRepository) ListTablesByFloorPlan(ctx context.Context, floorPlanID, tenantID string) ([]*RestaurantTable, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, tenant_id, floor_plan_id, name, capacity, pos_x, pos_y, shape, status, qr_token, qr_expires_at, created_at
		 FROM restaurant_tables WHERE floor_plan_id=$1 AND tenant_id=$2 ORDER BY name`,
		floorPlanID, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var tables []*RestaurantTable
	for rows.Next() {
		t := &RestaurantTable{}
		if err := rows.Scan(&t.ID, &t.TenantID, &t.FloorPlanID, &t.Name, &t.Capacity, &t.PosX, &t.PosY, &t.Shape, &t.Status, &t.QRToken, &t.QRExpiresAt, &t.CreatedAt); err != nil {
			return nil, err
		}
		tables = append(tables, t)
	}
	return tables, rows.Err()
}

func (r *TableRepository) CreateTable(ctx context.Context, t *RestaurantTable) error {
	return r.db.QueryRowContext(ctx,
		`INSERT INTO restaurant_tables (tenant_id, floor_plan_id, name, capacity, pos_x, pos_y, shape, status)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id, created_at`,
		t.TenantID, t.FloorPlanID, t.Name, t.Capacity, t.PosX, t.PosY, t.Shape, t.Status,
	).Scan(&t.ID, &t.CreatedAt)
}

func (r *TableRepository) GetTable(ctx context.Context, id, tenantID string) (*RestaurantTable, error) {
	t := &RestaurantTable{}
	err := r.db.QueryRowContext(ctx,
		`SELECT id, tenant_id, floor_plan_id, name, capacity, pos_x, pos_y, shape, status, qr_token, qr_expires_at, created_at
		 FROM restaurant_tables WHERE id=$1 AND tenant_id=$2`, id, tenantID,
	).Scan(&t.ID, &t.TenantID, &t.FloorPlanID, &t.Name, &t.Capacity, &t.PosX, &t.PosY, &t.Shape, &t.Status, &t.QRToken, &t.QRExpiresAt, &t.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return t, err
}

func (r *TableRepository) UpdateTable(ctx context.Context, t *RestaurantTable) error {
	res, err := r.db.ExecContext(ctx,
		`UPDATE restaurant_tables SET floor_plan_id=$1, name=$2, capacity=$3, pos_x=$4, pos_y=$5, shape=$6, status=$7
		 WHERE id=$8 AND tenant_id=$9`,
		t.FloorPlanID, t.Name, t.Capacity, t.PosX, t.PosY, t.Shape, t.Status, t.ID, t.TenantID)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *TableRepository) DeleteTable(ctx context.Context, id, tenantID string) error {
	res, err := r.db.ExecContext(ctx,
		`DELETE FROM restaurant_tables WHERE id=$1 AND tenant_id=$2`, id, tenantID)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return sql.ErrNoRows
	}
	return nil
}
