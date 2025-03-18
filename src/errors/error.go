package errs

import "errors"

var (
	ErrNotFound   = errors.New("Resource not found")
	ErrValidation = errors.New("Validation error")
	ErrBadRequest = errors.New("Bad Request")
	ErrDatabase   = errors.New("Database error")
	ErrInternal   = errors.New("Internal error")
)
