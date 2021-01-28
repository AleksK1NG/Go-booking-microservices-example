package user

import "github.com/labstack/echo/v4"

// Delivery
type Delivery interface {
	Register() echo.HandlerFunc
	Login() echo.HandlerFunc
	Logout() echo.HandlerFunc
	Update() echo.HandlerFunc
	Delete() echo.HandlerFunc
	GetUserByID() echo.HandlerFunc
	GetMe() echo.HandlerFunc
	GetCSRFToken() echo.HandlerFunc
}
