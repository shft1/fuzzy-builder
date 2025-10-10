package repositories

import (
	"context"

	"github.com/alexm/fuzzy-builder/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AttachmentRepository struct{ pool *pgxpool.Pool }

func NewAttachmentRepository(pool *pgxpool.Pool) *AttachmentRepository {
	return &AttachmentRepository{pool: pool}
}

func (r *AttachmentRepository) Create(ctx context.Context, a *models.Attachment) error {
	row := r.pool.QueryRow(ctx,
		`INSERT INTO attachments (defect_id, filename, filepath, uploaded_by) VALUES ($1,$2,$3,$4) RETURNING id, created_at`,
		a.DefectID, a.Filename, a.Filepath, a.UploadedBy,
	)
	return row.Scan(&a.ID, &a.CreatedAt)
}

func (r *AttachmentRepository) ListByDefect(ctx context.Context, defectID int64, limit, offset int32) ([]models.Attachment, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, defect_id, filename, filepath, uploaded_by, created_at FROM attachments WHERE defect_id=$1 ORDER BY id DESC LIMIT $2 OFFSET $3`,
		defectID, limit, offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []models.Attachment
	for rows.Next() {
		var a models.Attachment
		if err := rows.Scan(&a.ID, &a.DefectID, &a.Filename, &a.Filepath, &a.UploadedBy, &a.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, a)
	}
	return out, rows.Err()
}
