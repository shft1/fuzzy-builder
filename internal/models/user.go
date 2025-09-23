// Package models contains the data models for the application.
package models

import "time"

type Role string

const (
	RoleEngineer Role = "engineer"
	RoleManager  Role = "manager"
	RoleObserver Role = "observer"
)

type User struct {
	ID           int64     `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	Role         Role      `json:"role"`
	FullName     string    `json:"full_name"`
	CreatedAt    time.Time `json:"created_at"`
}
