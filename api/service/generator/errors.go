package generator

import "errors"

var (
	ErrPlateNotFound    = errors.New("plate not found")
	ErrPlateNotApproved = errors.New("plate is not approved")
	ErrNoRepoURL        = errors.New("plate has no repository URL")
	ErrFetchFailed      = errors.New("failed to fetch plate.yaml from repository")
	ErrMissingYAML      = errors.New("plate.yaml not found in repository")
	ErrInvalidInput     = errors.New("invalid generation input")
	ErrTemplateFailed   = errors.New("template rendering failed")
)
