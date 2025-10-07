package models

import "time"

type DefectStatus string

const (
	DefectStatusNew        DefectStatus = "new"
	DefectStatusInProgress DefectStatus = "in_progress"
	DefectStatusOnReview   DefectStatus = "on_review"
	DefectStatusClosed     DefectStatus = "closed"
)

type DefectPriority string

const (
	DefectPriorityLow    DefectPriority = "low"
	DefectPriorityMedium DefectPriority = "medium"
	DefectPriorityHigh   DefectPriority = "high"
)

type Defect struct {
	ID          int64          `json:"id"`
	Title       string         `json:"title"`
	Description string         `json:"description"`
	ProjectID   int64          `json:"project_id"`
	AssignedTo  *int64         `json:"assigned_to,omitempty"`
	Status      DefectStatus   `json:"status"`
	Priority    DefectPriority `json:"priority"`
	DueDate     *time.Time     `json:"due_date,omitempty"`
	CreatedBy   int64          `json:"created_by"`
	CreatedAt   time.Time      `json:"created_at"`
}
