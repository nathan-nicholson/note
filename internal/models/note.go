package models

import "time"

type Note struct {
	ID          int
	Content     string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	IsImportant bool
	Tags        []string
}
