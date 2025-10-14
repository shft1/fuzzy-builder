package services

import (
	"context"
	"errors"

	"github.com/alexm/fuzzy-builder/internal/models"
	"github.com/alexm/fuzzy-builder/internal/repositories"
)

var ErrInvalidStatusTransition = errors.New("invalid status transition")
var ErrPermissionDenied = errors.New("permission denied")

type DefectService struct {
	repo *repositories.DefectRepository
}

func NewDefectService(repo *repositories.DefectRepository) *DefectService {
	return &DefectService{repo: repo}
}

// CanTransition enforces allowed transitions in the defect lifecycle.
func (s *DefectService) CanTransition(from, to models.DefectStatus) bool {
	switch from {
	case models.DefectStatusNew:
		return to == models.DefectStatusInProgress
	case models.DefectStatusInProgress:
		return to == models.DefectStatusOnReview
	case models.DefectStatusOnReview:
		return to == models.DefectStatusClosed || to == models.DefectStatusInProgress
	case models.DefectStatusClosed:
		return false
	default:
		return false
	}
}

// UpdateStatus validates permissions and transition rules.
func (s *DefectService) UpdateStatus(ctx context.Context, currentUserID int64, defectID int64, newStatus models.DefectStatus) error {
	defect, err := s.repo.GetByID(ctx, defectID)
	if err != nil {
		return err
	}
	if !s.CanTransition(defect.Status, newStatus) {
		return ErrInvalidStatusTransition
	}
	// Permission rule example: only assigned engineer can move from in_progress -> on_review or back
	if defect.AssignedTo != nil && *defect.AssignedTo != currentUserID {
		return ErrPermissionDenied
	}
	return s.repo.UpdateStatus(ctx, defectID, newStatus)
}
