package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

// Staff represents a restaurant staff member.
type Staff struct {
	ID         string
	TenantID   string
	CognitoSub *string
	FullName   string
	Email      *string
	Role       string
	IsActive   bool
	CreatedAt  time.Time
}

// StaffRepository implements staff persistence using PostgreSQL.
type StaffRepository struct {
	db *sql.DB
}

// NewStaffRepository creates a StaffRepository.
func NewStaffRepository(db *sql.DB) *StaffRepository {
	return &StaffRepository{db: db}
}

// Create inserts a new staff member and fills in the generated ID and CreatedAt.
func (r *StaffRepository) Create(ctx context.Context, s *Staff) error {
	return r.db.QueryRowContext(ctx,
		`INSERT INTO staff (tenant_id, cognito_sub, full_name, email, role, is_active)
		 VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, created_at`,
		s.TenantID, s.CognitoSub, s.FullName, s.Email, s.Role, s.IsActive,
	).Scan(&s.ID, &s.CreatedAt)
}

// GetByID returns a staff member by ID, or nil if not found.
func (r *StaffRepository) GetByID(ctx context.Context, id, tenantID string) (*Staff, error) {
	s := &Staff{}
	err := r.db.QueryRowContext(ctx,
		`SELECT id, tenant_id, cognito_sub, full_name, email, role, is_active, created_at
		 FROM staff WHERE id=$1 AND tenant_id=$2`, id, tenantID,
	).Scan(&s.ID, &s.TenantID, &s.CognitoSub, &s.FullName, &s.Email, &s.Role, &s.IsActive, &s.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return s, err
}

// ListByTenant returns all staff for a given tenant.
func (r *StaffRepository) ListByTenant(ctx context.Context, tenantID string) ([]*Staff, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, tenant_id, cognito_sub, full_name, email, role, is_active, created_at
		 FROM staff WHERE tenant_id=$1 ORDER BY full_name`, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var staff []*Staff
	for rows.Next() {
		s := &Staff{}
		if err := rows.Scan(&s.ID, &s.TenantID, &s.CognitoSub, &s.FullName, &s.Email, &s.Role, &s.IsActive, &s.CreatedAt); err != nil {
			return nil, err
		}
		staff = append(staff, s)
	}
	return staff, rows.Err()
}

// Update updates mutable staff fields (full_name, email, role, is_active).
// Returns sql.ErrNoRows if the record does not exist.
func (r *StaffRepository) Update(ctx context.Context, s *Staff) error {
	res, err := r.db.ExecContext(ctx,
		`UPDATE staff SET full_name=$1, email=$2, role=$3, is_active=$4
		 WHERE id=$5 AND tenant_id=$6`,
		s.FullName, s.Email, s.Role, s.IsActive, s.ID, s.TenantID)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return sql.ErrNoRows
	}
	return nil
}

// Delete removes a staff member. Returns sql.ErrNoRows if not found.
func (r *StaffRepository) Delete(ctx context.Context, id, tenantID string) error {
	res, err := r.db.ExecContext(ctx,
		`DELETE FROM staff WHERE id=$1 AND tenant_id=$2`, id, tenantID)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return sql.ErrNoRows
	}
	return nil
}
