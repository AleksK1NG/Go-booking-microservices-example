package middlewares

import (
	"context"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/opentracing/opentracing-go"

	"github.com/AleksK1NG/hotels-mocroservices/user/config"
	"github.com/AleksK1NG/hotels-mocroservices/user/internal/user"
	httpErrors "github.com/AleksK1NG/hotels-mocroservices/user/pkg/http_errors"
	"github.com/AleksK1NG/hotels-mocroservices/user/pkg/logger"
)

// MiddlewareManager
type MiddlewareManager struct {
	logger logger.Logger
	cfg    *config.Config
	userUC user.UseCase
}

// NewMiddlewareManager
func NewMiddlewareManager(logger logger.Logger, cfg *config.Config, userUC user.UseCase) *MiddlewareManager {
	return &MiddlewareManager{logger: logger, cfg: cfg, userUC: userUC}
}

// Request Ctx User key
type RequestCtxUser struct{}

// Request Ctx Session key
type RequestCtxSession struct{}

// SessionMiddleware
func (m *MiddlewareManager) SessionMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		span, ctx := opentracing.StartSpanFromContext(c.Request().Context(), "user.SessionMiddleware")
		defer span.Finish()

		cookie, err := c.Cookie(m.cfg.HttpServer.SessionCookieName)
		if err != nil {
			if errors.Is(err, http.ErrNoCookie) {
				m.logger.Errorf("SessionMiddleware.ErrNoCookie: %v", err)
				return httpErrors.ErrorCtxResponse(c, err)
			}
			m.logger.Errorf("SessionMiddleware.c.Cookie: %v", err)
			return httpErrors.ErrorCtxResponse(c, err)
		}

		sessionByID, err := m.userUC.GetSessionByID(ctx, cookie.Value)
		if err != nil {
			m.logger.Errorf("SessionMiddleware.GetSessionByID: %v", err)
			return httpErrors.ErrorCtxResponse(c, err)
		}

		userResponse, err := m.userUC.GetByID(ctx, sessionByID.UserID)
		if err != nil {
			m.logger.Errorf("SessionMiddleware.userUC.GetByID: %v", err)
			return httpErrors.ErrorCtxResponse(c, err)
		}

		ctx = context.WithValue(c.Request().Context(), RequestCtxUser{}, userResponse)
		ctx = context.WithValue(ctx, RequestCtxSession{}, sessionByID)
		c.SetRequest(c.Request().WithContext(ctx))

		return next(c)
	}
}
