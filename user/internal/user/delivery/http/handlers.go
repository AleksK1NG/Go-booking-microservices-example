package http

import (
	"github.com/labstack/echo/v4"

	"github.com/AleksK1NG/hotels-mocroservices/user/internal/user"
	"github.com/AleksK1NG/hotels-mocroservices/user/pkg/logger"
)

// UserHandlers
type UserHandlers struct {
	group  *echo.Group
	userUC user.UseCase
	logger logger.Logger
}

// NewUserHandlers
func NewUserHandlers(group *echo.Group, userUC user.UseCase, logger logger.Logger) *UserHandlers {
	return &UserHandlers{group: group, userUC: userUC, logger: logger}
}
