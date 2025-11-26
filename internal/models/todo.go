package models

import (
	"database/sql"
	"time"
)

type Todo struct {
	ID          int
	Content     string
	IsComplete  bool
	DueDate     sql.NullTime
	CreatedAt   time.Time
	UpdatedAt   time.Time
	CompletedAt sql.NullTime
	Tags        []string
}
