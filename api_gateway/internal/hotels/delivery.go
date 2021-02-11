package hotels

import "github.com/labstack/echo/v4"

// Delivery
type Delivery interface {
	CreateHotel() echo.HandlerFunc
	UpdateHotel() echo.HandlerFunc
	GetHotelByID() echo.HandlerFunc
	GetHotels() echo.HandlerFunc
	UploadImage() echo.HandlerFunc
}
