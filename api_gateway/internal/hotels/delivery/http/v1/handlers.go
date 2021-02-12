package v1

import (
	"bytes"
	"io"
	"net/http"
	"strconv"
	"sync"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/opentracing/opentracing-go"
	uuid "github.com/satori/go.uuid"

	"github.com/AleksK1NG/hotels-mocroservices/api-gateway/config"
	"github.com/AleksK1NG/hotels-mocroservices/api-gateway/internal/hotels"
	"github.com/AleksK1NG/hotels-mocroservices/api-gateway/internal/models"
	httpErrors "github.com/AleksK1NG/hotels-mocroservices/api-gateway/pkg/http_errors"
	"github.com/AleksK1NG/hotels-mocroservices/api-gateway/pkg/logger"
	"github.com/AleksK1NG/hotels-mocroservices/api-gateway/pkg/utils"
)

const (
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

// Register CreateHotel
// @Tags Hotels
// @Summary Create new hotel
// @Description Create new hotel instance
// @Accept json
// @Produce json
// @Success 201 {object} models.Hotel
// @Router /hotels [post]
func (h *hotelsHandlers) CreateHotel() echo.HandlerFunc {
	return func(c echo.Context) error {
		span, ctx := opentracing.StartSpanFromContext(c.Request().Context(), "hotelsHandlers.CreateHotel")
		defer span.Finish()

		var hotelReq models.Hotel
		if err := c.Bind(&hotelReq); err != nil {
			h.logger.Error("c.Bind")
			return httpErrors.ErrorCtxResponse(c, err)
		}

		if err := h.validate.StructCtx(ctx, &hotelReq); err != nil {
			h.logger.Error("validate.StructCtx")
			return httpErrors.ErrorCtxResponse(c, err)
		}

		hotel, err := h.hotelsUC.CreateHotel(ctx, &hotelReq)
		if err != nil {
			h.logger.Error("hotelsUC.CreateHotel")
			return httpErrors.ErrorCtxResponse(c, err)
		}

		return c.JSON(http.StatusCreated, hotel)
	}
}

// Register UpdateHotel
// @Tags Hotels
// @Summary Update hotel data
// @Description Update single hotel data
// @Accept json
// @Produce json
// @Param hotel_id path int true "Hotel UUID"
// @Success 200 {object} models.Hotel
// @Router /hotels/{hotel_id} [put]
func (h *hotelsHandlers) UpdateHotel() echo.HandlerFunc {
	return func(c echo.Context) error {
		span, ctx := opentracing.StartSpanFromContext(c.Request().Context(), "hotelsHandlers.UpdateHotel")
		defer span.Finish()

		hotelUUID, err := uuid.FromString(c.QueryParam("hotel_id"))
		if err != nil {
			h.logger.Error("uuid.FromString")
			return httpErrors.ErrorCtxResponse(c, err)
		}

		var hotelReq models.Hotel
		if err := c.Bind(&hotelReq); err != nil {
			h.logger.Error("c.Bind")
			return httpErrors.ErrorCtxResponse(c, err)
		}
		hotelReq.HotelID = hotelUUID

		if err := h.validate.StructCtx(ctx, &hotelReq); err != nil {
			h.logger.Error("validate.StructCtx")
			return httpErrors.ErrorCtxResponse(c, err)
		}

		hotel, err := h.hotelsUC.UpdateHotel(ctx, &hotelReq)
		if err != nil {
			h.logger.Error("hotelsUC.UpdateHotel")
			return httpErrors.ErrorCtxResponse(c, err)
		}

		return c.JSON(http.StatusOK, hotel)
	}
}

// Register GetHotelByID
// @Tags Hotels
// @Summary Get hotel by id
// @Description Get single hotel by uuid
// @Accept json
// @Produce json
// @Param hotel_id query string false "hotel uuid"
// @Success 200 {object} models.Hotel
// @Router /hotels/{hotel_id} [get]
func (h *hotelsHandlers) GetHotelByID() echo.HandlerFunc {
	return func(c echo.Context) error {
		span, ctx := opentracing.StartSpanFromContext(c.Request().Context(), "hotelsHandlers.GetHotelByID")
		defer span.Finish()

		hotelUUID, err := uuid.FromString(c.QueryParam("hotel_id"))
		if err != nil {
			h.logger.Error("uuid.FromString")
			return httpErrors.ErrorCtxResponse(c, err)
		}

		hotelByID, err := h.hotelsUC.GetHotelByID(ctx, hotelUUID)
		if err != nil {
			h.logger.Error("hotelsUC.GetHotelByID")
			return httpErrors.ErrorCtxResponse(c, err)
		}

		return c.JSON(http.StatusOK, hotelByID)
	}
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
			h.logger.Error("strconv.Atoi")
			return httpErrors.ErrorCtxResponse(c, err)
		}
		size, err := strconv.Atoi(c.QueryParam("size"))
		if err != nil {
			h.logger.Error("strconv.Atoi")
			return httpErrors.ErrorCtxResponse(c, err)
		}

		hotelsList, err := h.hotelsUC.GetHotels(ctx, int64(page), int64(size))
		if err != nil {
			h.logger.Error("hotelsUC.GetHotels")
			return httpErrors.ErrorCtxResponse(c, err)
		}

		return c.JSON(http.StatusOK, hotelsList)
	}
}

// UploadImage godoc
// @Summary Upload hotel image
// @Tags Hotels
// @Description Upload hotel logo image
// @Accept mpfd
// @Produce json
// @Param hotel_id query string false "hotel uuid"
// @Success 200 {string} ""
// @Router /hotels/{id}/image [put]
func (h *hotelsHandlers) UploadImage() echo.HandlerFunc {
	bufferPool := &sync.Pool{New: func() interface{} {
		return &bytes.Buffer{}
	}}
	return func(c echo.Context) error {
		span, ctx := opentracing.StartSpanFromContext(c.Request().Context(), "hotelsHandlers.UploadImage")
		defer span.Finish()

		hotelUUID, err := uuid.FromString(c.QueryParam("hotel_id"))
		if err != nil {
			return err
		}

		if err := c.Request().ParseMultipartForm(maxFileSize); err != nil {
			h.logger.Error("c.ParseMultipartForm")
			return httpErrors.ErrorCtxResponse(c, err)
		}

		c.Request().Body = http.MaxBytesReader(c.Response(), c.Request().Body, maxFileSize)
		defer c.Request().Body.Close()

		formFile, _, err := c.Request().FormFile("avatar")
		if err != nil {
			h.logger.Error("c.FormFile")
			return httpErrors.ErrorCtxResponse(c, err)
		}

		fileType, err := utils.CheckImageUpload(formFile)
		if err != nil {
			h.logger.Error("h.checkAvatar")
			return httpErrors.ErrorCtxResponse(c, err)
		}

		buf, ok := bufferPool.Get().(*bytes.Buffer)
		if !ok {
			h.logger.Error("bufferPool.Get")
			return httpErrors.ErrorCtxResponse(c, httpErrors.InternalServerError)
		}
		defer bufferPool.Put(buf)
		buf.Reset()

		if _, err := io.Copy(buf, formFile); err != nil {
			h.logger.Error("io.Copy")
			return httpErrors.ErrorCtxResponse(c, err)
		}

		if err := h.hotelsUC.UploadImage(ctx, buf.Bytes(), fileType, hotelUUID.String()); err != nil {
			h.logger.Error("hotelsUC.UploadImage")
			return httpErrors.ErrorCtxResponse(c, err)
		}

		return c.NoContent(http.StatusOK)
	}
}
