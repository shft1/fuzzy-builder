package repositories

import (
	"context"

	"github.com/alexm/fuzzy-builder/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProjectRepository struct {
	pool *pgxpool.Pool
}

func NewProjectRepository(pool *pgxpool.Pool) *ProjectRepository {
	return &ProjectRepository{pool: pool}
}

func (r *ProjectRepository) Create(ctx context.Context, p *models.Project) error {
	row := r.pool.QueryRow(ctx,
		`INSERT INTO projects (name, description, created_by) VALUES ($1,$2,$3) RETURNING id, created_at`,
		p.Name, p.Description, p.CreatedBy,
	)
	return row.Scan(&p.ID, &p.CreatedAt)
}

func (r *ProjectRepository) Update(ctx context.Context, p *models.Project) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE projects SET name=$1, description=$2 WHERE id=$3`,
		p.Name, p.Description, p.ID,
	)
	return err
}

func (r *ProjectRepository) Delete(ctx context.Context, id int64) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM projects WHERE id=$1`, id)
	return err
}

func (r *ProjectRepository) List(ctx context.Context, limit, offset int32) ([]models.Project, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, name, description, created_by, created_at FROM projects ORDER BY id DESC LIMIT $1 OFFSET $2`,
		limit, offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []models.Project
	for rows.Next() {
		var p models.Project
		if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.CreatedBy, &p.CreatedAt); err != nil {
			return nil, err
		}
		items = append(items, p)
	}
	return items, rows.Err()
}
