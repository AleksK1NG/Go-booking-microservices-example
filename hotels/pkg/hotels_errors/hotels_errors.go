package hotels_errors

import "github.com/pkg/errors"

var (
	ErrInvalidUUID            = errors.New("Invalid uuid")
	ErrInvalidDeliveryHeaders = errors.New("Invalid delivery headers")
	ErrInternalServerError    = errors.New("Internal server error")
	ErrInvalidImageFormat     = errors.New("Invalid file format")
	ErrHotelNotFound          = errors.New("Hotel not found")
)
