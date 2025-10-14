package services

import (
	"testing"

	"github.com/alexm/fuzzy-builder/internal/models"
)

func TestDefectServiceTransitions(t *testing.T) {
	svc := NewDefectService(nil)
	if !svc.CanTransition(models.DefectStatusNew, models.DefectStatusInProgress) {
		t.Fatalf("expected transition new->in_progress")
	}
	if svc.CanTransition(models.DefectStatusClosed, models.DefectStatusOnReview) {
		t.Fatalf("closed should not transition")
	}
}
