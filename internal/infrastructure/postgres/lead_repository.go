package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nomad-pixel/imperial/internal/domain/entities"
	"github.com/nomad-pixel/imperial/internal/domain/ports"
)

type leadRepository struct {
	db *pgxpool.Pool
}

func NewLeadRepository(db *pgxpool.Pool) ports.LeadRepository {
	return &leadRepository{db: db}
}

func (r *leadRepository) CreateLead(ctx context.Context, lead *entities.Lead) error {
	query := `
		INSERT INTO leads (full_name, phone, period, created_at)
		VALUES ($1, $2, tstzrange($3, $4, '[]'), $5)
		RETURNING id
	`
	return r.db.QueryRow(ctx, query, lead.FullName, lead.Phone, lead.StartDate, lead.EndDate, lead.CreatedAt).Scan(&lead.ID)
}

func (r *leadRepository) GetLeadByID(ctx context.Context, id int64) (*entities.Lead, error) {
	query := `
		SELECT id, full_name, phone, lower(period), upper(period), created_at
		FROM leads
		WHERE id = $1
	`
	lead := &entities.Lead{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&lead.ID,
		&lead.FullName,
		&lead.Phone,
		&lead.StartDate,
		&lead.EndDate,
		&lead.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return lead, nil
}

func (r *leadRepository) ListLeads(ctx context.Context, offset, limit int64) (int64, []*entities.Lead, error) {
	countQuery := `SELECT COUNT(*) FROM leads`
	var total int64
	if err := r.db.QueryRow(ctx, countQuery).Scan(&total); err != nil {
		return 0, nil, err
	}

	query := `
		SELECT id, full_name, phone, lower(period), upper(period), created_at
		FROM leads
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`
	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return 0, nil, err
	}
	defer rows.Close()

	leads := make([]*entities.Lead, 0)
	for rows.Next() {
		lead := &entities.Lead{}
		if err := rows.Scan(
			&lead.ID,
			&lead.FullName,
			&lead.Phone,
			&lead.StartDate,
			&lead.EndDate,
			&lead.CreatedAt,
		); err != nil {
			return 0, nil, err
		}
		leads = append(leads, lead)
	}

	if err := rows.Err(); err != nil {
		return 0, nil, err
	}

	return total, leads, nil
}

func (r *leadRepository) DeleteLead(ctx context.Context, id int64) error {
	query := `DELETE FROM leads WHERE id = $1`
	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("lead not found")
	}

	return nil
}
