package http

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/opentracing/opentracing-go"

	"github.com/AleksK1NG/hotels-mocroservices/user/internal/models"
	"github.com/AleksK1NG/hotels-mocroservices/user/internal/user"
	httpErrors "github.com/AleksK1NG/hotels-mocroservices/user/pkg/http_errors"
	"github.com/AleksK1NG/hotels-mocroservices/user/pkg/logger"
)

// UserHandlers
type UserHandlers struct {
	group    *echo.Group
	userUC   user.UseCase
	logger   logger.Logger
	validate *validator.Validate
}

// NewUserHandlers
func NewUserHandlers(group *echo.Group, userUC user.UseCase, logger logger.Logger, validate *validator.Validate) *UserHandlers {
	return &UserHandlers{group: group, userUC: userUC, logger: logger, validate: validate}
}

// Register new user
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
			h.logger.Errorf("userUC.Register: %v", err)
			return httpErrors.ErrorCtxResponse(c, err)
		}

		return c.JSON(http.StatusOK, regUser)
	}
}

func (h *UserHandlers) Login() echo.HandlerFunc {
	panic("implement me")
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
