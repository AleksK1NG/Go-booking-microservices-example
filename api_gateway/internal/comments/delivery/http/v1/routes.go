package v1

import "github.com/labstack/echo/v4"

// MapRoutes
func (c *commentsHandlers) MapRoutes() {
	c.group.GET("", func(c echo.Context) error {
		return c.JSON(200, "Ok")
	})
}
