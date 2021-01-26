package server

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	userHandlers "github.com/AleksK1NG/hotels-mocroservices/user/internal/user/delivery/http"
	"github.com/AleksK1NG/hotels-mocroservices/user/internal/user/repository"
	"github.com/AleksK1NG/hotels-mocroservices/user/internal/user/usecase"
)

const (
	gzipLevel       = 5
	stackSize       = 1 << 10 // 1 KB
	csrfTokenHeader = "X-CSRF-Token"
	bodyLimit       = "2M"
)

func (s *Server) MapRoutes() {
	v1 := s.echo.Group("/api/v1")
	usersGroup := v1.Group("/users")

	userPGRepository := repository.NewUserPGRepository(s.pgxPool)
	userUseCase := usecase.NewUserUseCase(userPGRepository)
	uh := userHandlers.NewUserHandlers(usersGroup, userUseCase, s.logger)
	uh.MapUserRoutes()

	s.echo.Pre(middleware.HTTPSRedirect())

	s.echo.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderXRequestID, csrfTokenHeader},
	}))
	s.echo.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		StackSize:         stackSize,
		DisablePrintStack: true,
		DisableStackAll:   true,
	}))
	s.echo.Use(middleware.RequestID())

	s.echo.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: gzipLevel,
		Skipper: func(c echo.Context) bool {
			return strings.Contains(c.Request().URL.Path, "swagger")
		},
	}))

	s.echo.Use(middleware.Secure())
	s.echo.Use(middleware.BodyLimit(bodyLimit))
	// if s.cfg.Server.Debug {
	// 	s.echo.Use(mw.DebugMiddleware)
	// }

	v1.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "Ok")
	})
}
