package db

import "github.com/pkg/errors"

var (
	// ErrNoResult is raised when data query returns no result
	ErrNoResult = errors.New("No result")
	// ErrTooManyResults is raised when data query returns too many results
	ErrTooManyResults = errors.New("too many results returned")
	// ErrNoModification is raised when updating an entity without any changes
	ErrNoModification = errors.New("No changes made")
)
