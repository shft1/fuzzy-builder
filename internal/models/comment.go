package models

import "time"

type Comment struct {
	ID        int64     `json:"id"`
	DefectID  int64     `json:"defect_id"`
	UserID    int64     `json:"user_id"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"created_at"`
}
