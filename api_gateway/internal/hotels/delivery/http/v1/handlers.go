package v1

import (
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/opentracing/opentracing-go"

	"github.com/AleksK1NG/hotels-mocroservices/api-gateway/config"
	"github.com/AleksK1NG/hotels-mocroservices/api-gateway/internal/hotels"
	"github.com/AleksK1NG/hotels-mocroservices/api-gateway/pkg/logger"
)

const (
	csrfHeader  = "X-CSRF-Token"
	maxFileSize = 1024 * 1024 * 10
)

// HotelsHandlers
type hotelsHandlers struct {
	cfg      *config.Config
	group    *echo.Group
	logger   logger.Logger
	validate *validator.Validate
	hotelsUC hotels.UseCase
	// mw       *middlewares.MiddlewareManager
}

// NewHotelsHandlers
func NewHotelsHandlers(
	cfg *config.Config,
	group *echo.Group,
	logger logger.Logger,
	validate *validator.Validate,
	hotelsUC hotels.UseCase,
) *hotelsHandlers {
	return &hotelsHandlers{cfg: cfg, group: group, logger: logger, validate: validate, hotelsUC: hotelsUC}
}

func (h *hotelsHandlers) CreateHotel() echo.HandlerFunc {
	panic("implement me")
}

func (h *hotelsHandlers) UpdateHotel() echo.HandlerFunc {
	panic("implement me")
}

func (h *hotelsHandlers) GetHotelByID() echo.HandlerFunc {
	panic("implement me")
}

// Register GetHotels
// @Tags Hotels
// @Summary Get hotels list new user
// @Description Get hotels list with pagination using page and size query parameters
// @Accept json
// @Produce json
// @Param page query int false "page number"
// @Param size query int false "number of elements"
// @Success 200 {object} models.HotelsListRes
// @Router /hotels [get]
func (h *hotelsHandlers) GetHotels() echo.HandlerFunc {
	return func(c echo.Context) error {
		span, ctx := opentracing.StartSpanFromContext(c.Request().Context(), "hotelsHandlers.GetHotels")
		defer span.Finish()

		page, err := strconv.Atoi(c.QueryParam("page"))
		if err != nil {
			return err
		}
		size, err := strconv.Atoi(c.QueryParam("size"))
		if err != nil {
			return err
		}

		hotelsList, err := h.hotelsUC.GetHotels(ctx, int64(page), int64(size))
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, hotelsList)
	}
}

func (h *hotelsHandlers) UploadImage() echo.HandlerFunc {
	panic("implement me")
}
