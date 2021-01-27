package http

import (
	"github.com/labstack/echo/v4"
)

// MapUserRoutes
func (h *UserHandlers) MapUserRoutes() {
	h.group.GET("/me", func(c echo.Context) error {
		return c.String(200, "PRO")
	})
	h.group.POST("/register", h.Register())
}
