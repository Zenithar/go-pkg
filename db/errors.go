package db

import "golang.org/x/xerrors"

var (
	// ErrNoResult is raised when data query returns no result
	ErrNoResult = xerrors.New("No result")
	// ErrTooManyResults is raised when data query returns too many results
	ErrTooManyResults = xerrors.New("too many results returned")
	// ErrNoModification is raised when updating an entity without any changes
	ErrNoModification = xerrors.New("No changes made")
)
