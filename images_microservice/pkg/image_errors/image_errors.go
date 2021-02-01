package image_errors

import "github.com/pkg/errors"

var (
	ErrInvalidUUID            = errors.New("Invalid uuid")
	ErrInvalidDeliveryHeaders = errors.New("Invalid uuid")
	ErrInternalServerError    = errors.New("Internal server error")
	ErrInvalidImageFormat     = errors.New("Invalid image format")
)
