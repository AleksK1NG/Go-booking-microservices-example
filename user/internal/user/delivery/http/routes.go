package http

import (
	"github.com/labstack/echo/v4"
)

func (h *UserHandlers) MapUserRoutes() {
	h.group.GET("/me", func(c echo.Context) error {
		return c.String(200, "PRO")
	})
}
