package errs

import "errors"

var (
	NotFound      = errors.New("Resource not found")
	Validation    = errors.New("Validation error")
	BadRequest    = errors.New("Bad Request")
	Database      = errors.New("Database error")
	Internal      = errors.New("Internal error")
	AlreadyExists = errors.New("Already exists")
	Unauthorized  = errors.New("Unauthorized")

	ClientResponseNoOK = errors.New("External api response not OK")
	ClientNotRightType = errors.New("External api response not correct type")
	ImageCropError     = errors.New("Error will croping image")
)
