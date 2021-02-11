package v1

import (
	"github.com/labstack/echo/v4"
)

// MapRoutes
func (h *hotelsHandlers) MapRoutes() {
	h.group.GET("", func(c echo.Context) error {
		return c.JSON(200, "Ok")
	})
}
