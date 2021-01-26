package server

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func (s *Server) MapRoutes() {
	s.echo.Pre(middleware.HTTPSRedirect())

	s.echo.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderXRequestID, "X-CSRF-Token"},
	}))
	s.echo.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		StackSize:         1 << 10, // 1 KB
		DisablePrintStack: true,
		DisableStackAll:   true,
	}))
	s.echo.Use(middleware.RequestID())

	s.echo.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 5,
		Skipper: func(c echo.Context) bool {
			return strings.Contains(c.Request().URL.Path, "swagger")
		},
	}))

	s.echo.Use(middleware.Secure())
	s.echo.Use(middleware.BodyLimit("2M"))
	// if s.cfg.Server.Debug {
	// 	s.echo.Use(mw.DebugMiddleware)
	// }

	v1 := s.echo.Group("/api/v1")

	// usersRoutes := v1.Group("/users")

	v1.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "Ok")
	})
}
