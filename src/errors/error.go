package errs

import "errors"

var (
	NotFound      = errors.New("Resource not found")
	Validation    = errors.New("Validation error")
	BadRequest    = errors.New("Bad Request")
	Database      = errors.New("Database error")
	Internal      = errors.New("Internal error")
	AlreadyExists = errors.New("Already exists")
)
