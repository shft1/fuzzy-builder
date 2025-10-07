package models

import "time"

type Attachment struct {
	ID         int64     `json:"id"`
	DefectID   int64     `json:"defect_id"`
	Filename   string    `json:"filename"`
	Filepath   string    `json:"filepath"`
	UploadedBy int64     `json:"uploaded_by"`
	CreatedAt  time.Time `json:"created_at"`
}
