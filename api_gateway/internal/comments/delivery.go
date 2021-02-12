package comments

import "github.com/labstack/echo/v4"

// Delivery
type Delivery interface {
	CreateComment() echo.HandlerFunc
	GetCommByID() echo.HandlerFunc
	UpdateComment() echo.HandlerFunc
	GetByHotelID() echo.HandlerFunc
}
