package repositories

import (
	"context"

	"github.com/alexm/fuzzy-builder/internal/models"
)

type Users interface {
	Create(ctx context.Context, user *models.User) error
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	GetByID(ctx context.Context, id int64) (*models.User, error)
}

type Projects interface {
	Create(ctx context.Context, p *models.Project) error
	Update(ctx context.Context, p *models.Project) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context, limit, offset int32) ([]models.Project, error)
}

type Defects interface {
	Create(ctx context.Context, d *models.Defect) error
	GetByID(ctx context.Context, id int64) (*models.Defect, error)
	UpdateStatus(ctx context.Context, id int64, status models.DefectStatus) error
	List(ctx context.Context, f DefectFilter) ([]models.Defect, error)
	CountByStatus(ctx context.Context) (map[string]int64, error)
	CountByProject(ctx context.Context) (map[int64]int64, error)
}

type Attachments interface {
	Create(ctx context.Context, a *models.Attachment) error
	ListByDefect(ctx context.Context, defectID int64, limit, offset int32) ([]models.Attachment, error)
	GetByID(ctx context.Context, id int64) (*models.Attachment, error)
}
