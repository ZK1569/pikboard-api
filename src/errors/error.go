package errs

import "errors"

var (
	NotFound     = errors.New("Resource not found")
	Validation   = errors.New("Validation error")
	BadRequest   = errors.New("Bad Request")
	Database     = errors.New("Database error")
	Internal     = errors.New("Internal error")
	Unauthorized = errors.New("Unauthorized")

	AlreadyExists      = errors.New("Already exists")
	UserAlreadyExists  = errors.New("User already exists")
	EmailAlreadyExists = errors.New("Email already exists")

	ClientResponseNoOK = errors.New("External api response not OK")
	ClientNotRightType = errors.New("External api response not correct type")
	ImageCropError     = errors.New("Error will croping image")

	S3Error = errors.New("Error with the S3")
)
