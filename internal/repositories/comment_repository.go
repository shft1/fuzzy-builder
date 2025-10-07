package repositories

import (
	"context"

	"github.com/alexm/fuzzy-builder/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CommentRepository struct{ pool *pgxpool.Pool }

func NewCommentRepository(pool *pgxpool.Pool) *CommentRepository {
	return &CommentRepository{pool: pool}
}

func (r *CommentRepository) Create(ctx context.Context, c *models.Comment) error {
	row := r.pool.QueryRow(ctx,
		`INSERT INTO comments (defect_id, user_id, text) VALUES ($1,$2,$3) RETURNING id, created_at`,
		c.DefectID, c.UserID, c.Text,
	)
	return row.Scan(&c.ID, &c.CreatedAt)
}

func (r *CommentRepository) ListByDefect(ctx context.Context, defectID int64, limit, offset int32) ([]models.Comment, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, defect_id, user_id, text, created_at FROM comments WHERE defect_id=$1 ORDER BY id DESC LIMIT $2 OFFSET $3`,
		defectID, limit, offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []models.Comment
	for rows.Next() {
		var c models.Comment
		if err := rows.Scan(&c.ID, &c.DefectID, &c.UserID, &c.Text, &c.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	return out, rows.Err()
}
