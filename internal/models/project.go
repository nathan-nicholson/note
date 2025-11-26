package models

import (
	"database/sql"
	"regexp"
	"time"
)

type Project struct {
	ID               int
	Name             string
	CreatedAt        time.Time
	FirstActivatedAt sql.NullTime
	LastActivityAt   sql.NullTime
	ClosedAt         sql.NullTime
	IsClosed         bool
	Tags             []string
}

var reservedProjectNames = map[string]bool{
	"create": true,
	"close":  true,
	"reopen": true,
	"list":   true,
	"status": true,
	"show":   true,
	"edit":   true,
	"delete": true,
}

var kebabCaseRegex = regexp.MustCompile(`^[a-z0-9]+(-[a-z0-9]+)*$`)

func ValidateProjectName(name string) error {
	if reservedProjectNames[name] {
		return &ProjectNameReservedError{Name: name}
	}

	if !kebabCaseRegex.MatchString(name) {
		return &ProjectNameInvalidError{Name: name}
	}

	return nil
}

type ProjectNameReservedError struct {
	Name string
}

func (e *ProjectNameReservedError) Error() string {
	return "project name '" + e.Name + "' is reserved"
}

type ProjectNameInvalidError struct {
	Name string
}

func (e *ProjectNameInvalidError) Error() string {
	return "project name must be kebab-case (lowercase letters and hyphens only)"
}
