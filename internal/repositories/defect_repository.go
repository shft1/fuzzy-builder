package repositories

import (
	"context"
	"strings"

	"github.com/alexm/fuzzy-builder/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DefectFilter struct {
	ProjectID  *int64
	AssignedTo *int64
	Status     *models.DefectStatus
	Priority   *models.DefectPriority
	Limit      int32
	Offset     int32
}

type DefectRepository struct {
	pool *pgxpool.Pool
}

func NewDefectRepository(pool *pgxpool.Pool) *DefectRepository {
	return &DefectRepository{pool: pool}
}

func (r *DefectRepository) Create(ctx context.Context, d *models.Defect) error {
	row := r.pool.QueryRow(ctx,
		`INSERT INTO defects (title, description, project_id, assigned_to, status, priority, due_date, created_by)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
		 RETURNING id, created_at`,
		d.Title, d.Description, d.ProjectID, d.AssignedTo, d.Status, d.Priority, d.DueDate, d.CreatedBy,
	)
	return row.Scan(&d.ID, &d.CreatedAt)
}

func (r *DefectRepository) UpdateStatus(ctx context.Context, id int64, status models.DefectStatus) error {
	_, err := r.pool.Exec(ctx, `UPDATE defects SET status=$1 WHERE id=$2`, status, id)
	return err
}

func (r *DefectRepository) List(ctx context.Context, f DefectFilter) ([]models.Defect, error) {
	where := []string{"1=1"}
	args := []any{}
	idx := 1
	if f.ProjectID != nil {
		where = append(where, "project_id=$"+itoa(idx))
		args = append(args, *f.ProjectID)
		idx++
	}
	if f.AssignedTo != nil {
		where = append(where, "assigned_to=$"+itoa(idx))
		args = append(args, *f.AssignedTo)
		idx++
	}
	if f.Status != nil {
		where = append(where, "status=$"+itoa(idx))
		args = append(args, *f.Status)
		idx++
	}
	if f.Priority != nil {
		where = append(where, "priority=$"+itoa(idx))
		args = append(args, *f.Priority)
		idx++
	}
	query := "SELECT id, title, description, project_id, assigned_to, status, priority, due_date, created_by, created_at FROM defects WHERE " + strings.Join(where, " AND ") + " ORDER BY id DESC LIMIT $" + itoa(idx) + " OFFSET $" + itoa(idx+1)
	args = append(args, f.Limit, f.Offset)
	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []models.Defect
	for rows.Next() {
		var d models.Defect
		if err := rows.Scan(&d.ID, &d.Title, &d.Description, &d.ProjectID, &d.AssignedTo, &d.Status, &d.Priority, &d.DueDate, &d.CreatedBy, &d.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, d)
	}
	return out, rows.Err()
}

func itoa(n int) string {
	// fast small int to string
	const digits = "0123456789"
	if n == 0 {
		return "0"
	}
	buf := make([]byte, 0, 4)
	for n > 0 {
		d := n % 10
		buf = append([]byte{digits[d]}, buf...)
		n /= 10
	}
	return string(buf)
}
