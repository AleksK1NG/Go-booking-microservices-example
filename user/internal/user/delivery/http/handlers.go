package http

import (
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/opentracing/opentracing-go"

	"github.com/AleksK1NG/hotels-mocroservices/user/config"
	"github.com/AleksK1NG/hotels-mocroservices/user/internal/models"
	"github.com/AleksK1NG/hotels-mocroservices/user/internal/user"
	httpErrors "github.com/AleksK1NG/hotels-mocroservices/user/pkg/http_errors"
	"github.com/AleksK1NG/hotels-mocroservices/user/pkg/logger"
)

// UserHandlers
type UserHandlers struct {
	cfg      *config.Config
	group    *echo.Group
	userUC   user.UseCase
	logger   logger.Logger
	validate *validator.Validate
}

// NewUserHandlers
func NewUserHandlers(group *echo.Group, userUC user.UseCase, logger logger.Logger, validate *validator.Validate, cfg *config.Config) *UserHandlers {
	return &UserHandlers{group: group, userUC: userUC, logger: logger, validate: validate, cfg: cfg}
}

// Register godoc
// @Summary Register new user
// @Description register new user account, returns user data and session
// @Accept json
// @Produce json
// @Param data body models.User true "user data"
// @Success 201 {object} models.UserResponse
// @Router /auth/register [post]
func (h *UserHandlers) Register() echo.HandlerFunc {
	return func(c echo.Context) error {
		span, ctx := opentracing.StartSpanFromContext(c.Request().Context(), "auth.Register")
		defer span.Finish()

		var u models.User
		if err := c.Bind(&u); err != nil {
			h.logger.Errorf("c.Bind: %v", err)
			return httpErrors.ErrorCtxResponse(c, err)
		}

		if err := h.validate.StructCtx(ctx, &u); err != nil {
			h.logger.Errorf("validate.StructCtx: %v", err)
			return httpErrors.ErrorCtxResponse(c, err)
		}

		regUser, err := h.userUC.Register(ctx, &u)
		if err != nil {
			h.logger.Errorf("UserHandlers.Register.userUC.Register: %v", err)
			return httpErrors.ErrorCtxResponse(c, err)
		}

		sessionID, err := h.userUC.CreateSession(ctx, regUser.UserID)
		if err != nil {
			h.logger.Errorf("UserHandlers.userUC.CreateSession: %v", err)
			return httpErrors.ErrorCtxResponse(c, err)
		}

		c.SetCookie(&http.Cookie{
			Name:     "session_token",
			Value:    sessionID,
			Path:     "/",
			HttpOnly: true,
			Expires:  time.Now().Add(time.Duration(h.cfg.HttpServer.CookieLifeTime) * time.Minute),
		})

		return c.JSON(http.StatusCreated, regUser)
	}
}

// Register godoc
// @Summary Login user
// @Description login user, returns user data and session
// @Accept json
// @Produce json
// @Param data body models.Login true "email and password"
// @Success 200 {object} models.UserResponse
// @Router /auth/login [post]
func (h *UserHandlers) Login() echo.HandlerFunc {
	return func(c echo.Context) error {
		span, ctx := opentracing.StartSpanFromContext(c.Request().Context(), "auth.Login")
		defer span.Finish()

		var login models.Login
		if err := c.Bind(&login); err != nil {
			h.logger.Errorf("c.Bind: %v", err)
			return httpErrors.ErrorCtxResponse(c, err)
		}

		if err := h.validate.StructCtx(ctx, &login); err != nil {
			h.logger.Errorf("validate.StructCtx: %v", err)
			return httpErrors.ErrorCtxResponse(c, err)
		}

		userResponse, err := h.userUC.Login(ctx, login)
		if err != nil {
			h.logger.Errorf("UserHandlers.userUC.Login: %v", err)
			return httpErrors.ErrorCtxResponse(c, err)
		}

		sessionID, err := h.userUC.CreateSession(ctx, userResponse.UserID)
		if err != nil {
			h.logger.Errorf("UserHandlers.Login.CreateSession: %v", err)
			return httpErrors.ErrorCtxResponse(c, err)
		}

		c.SetCookie(&http.Cookie{
			Name:     "session_token",
			Value:    sessionID,
			Path:     "/",
			HttpOnly: true,
			Expires:  time.Now().Add(time.Duration(h.cfg.HttpServer.CookieLifeTime) * time.Minute),
		})

		return c.JSON(http.StatusOK, userResponse)
	}
}

func (h *UserHandlers) Logout() echo.HandlerFunc {
	panic("implement me")
}

func (h *UserHandlers) Update() echo.HandlerFunc {
	panic("implement me")
}

func (h *UserHandlers) Delete() echo.HandlerFunc {
	panic("implement me")
}

func (h *UserHandlers) GetUserByID() echo.HandlerFunc {
	panic("implement me")
}
