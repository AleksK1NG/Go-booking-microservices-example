package v1

import (
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/opentracing/opentracing-go"
	uuid "github.com/satori/go.uuid"

	"github.com/AleksK1NG/hotels-mocroservices/api-gateway/config"
	"github.com/AleksK1NG/hotels-mocroservices/api-gateway/internal/comments"
	"github.com/AleksK1NG/hotels-mocroservices/api-gateway/internal/models"
	httpErrors "github.com/AleksK1NG/hotels-mocroservices/api-gateway/pkg/http_errors"
	"github.com/AleksK1NG/hotels-mocroservices/api-gateway/pkg/logger"
)

// CommentsHandlers
type commentsHandlers struct {
	cfg      *config.Config
	group    *echo.Group
	logger   logger.Logger
	validate *validator.Validate
	commUC   comments.UseCase
}

// NewCommentsHandlers
func NewCommentsHandlers(
	cfg *config.Config,
	group *echo.Group,
	logger logger.Logger,
	validate *validator.Validate,
	commUC comments.UseCase,
) *commentsHandlers {
	return &commentsHandlers{cfg: cfg, group: group, logger: logger, validate: validate, commUC: commUC}
}

// Register CreateComment
// @Tags Comments
// @Summary Create new comment
// @Description Create new single comment
// @Accept json
// @Produce json
// @Success 201 {object} models.Comment
// @Router /comments [post]
func (h *commentsHandlers) CreateComment() echo.HandlerFunc {
	return func(c echo.Context) error {
		span, ctx := opentracing.StartSpanFromContext(c.Request().Context(), "commentsHandlers.CreateComment")
		defer span.Finish()

		var comm models.Comment
		if err := c.Bind(&comm); err != nil {
			h.logger.Error("c.Bind")
			return httpErrors.ErrorCtxResponse(c, err)
		}

		if err := h.validate.StructCtx(ctx, &comm); err != nil {
			h.logger.Error("validate.StructCtx")
			return httpErrors.ErrorCtxResponse(c, err)
		}

		comment, err := h.commUC.CreateComment(ctx, &comm)
		if err != nil {
			h.logger.Error("commUC.CreateComment")
			return httpErrors.ErrorCtxResponse(c, err)
		}

		return c.JSON(http.StatusCreated, comment)
	}
}

// Register GetCommByID
// @Tags Comments
// @Summary Get comment by id
// @Description Get comment by uuid
// @Accept json
// @Produce json
// @Param comment_id query string false "comment uuid"
// @Success 200 {object} models.Comment
// @Router /comments/{comment_id} [get]
func (h *commentsHandlers) GetCommByID() echo.HandlerFunc {
	return func(c echo.Context) error {
		span, ctx := opentracing.StartSpanFromContext(c.Request().Context(), "commentsHandlers.GetCommByID")
		defer span.Finish()

		commUUID, err := uuid.FromString(c.QueryParam("comment_id"))
		if err != nil {
			h.logger.Error("uuid.FromString")
			return httpErrors.ErrorCtxResponse(c, err)
		}

		commByID, err := h.commUC.GetCommByID(ctx, commUUID)
		if err != nil {
			h.logger.Error("commUC.GetCommByID")
			return httpErrors.ErrorCtxResponse(c, err)
		}

		return c.JSON(http.StatusOK, commByID)
	}
}

// Register UpdateComment
// @Tags Comments
// @Summary Update comment by id
// @Description Update comment by uuid
// @Accept json
// @Produce json
// @Param comment_id query string false "comment uuid"
// @Success 200 {object} models.Comment
// @Router /comments/{comment_id} [put]
func (h *commentsHandlers) UpdateComment() echo.HandlerFunc {
	return func(c echo.Context) error {
		span, ctx := opentracing.StartSpanFromContext(c.Request().Context(), "commentsHandlers.UpdateComment")
		defer span.Finish()

		commUUID, err := uuid.FromString(c.QueryParam("comment_id"))
		if err != nil {
			h.logger.Error("uuid.FromString")
			return httpErrors.ErrorCtxResponse(c, err)
		}

		var comm models.Comment
		if err := c.Bind(&comm); err != nil {
			h.logger.Error("c.Bind")
			return httpErrors.ErrorCtxResponse(c, err)
		}
		comm.CommentID = commUUID

		if err := h.validate.StructCtx(ctx, &comm); err != nil {
			h.logger.Error("validate.StructCtx")
			return httpErrors.ErrorCtxResponse(c, err)
		}

		comment, err := h.commUC.UpdateComment(ctx, &comm)
		if err != nil {
			h.logger.Error("commUC.UpdateComment")
			return httpErrors.ErrorCtxResponse(c, err)
		}

		return c.JSON(http.StatusOK, comment)
	}
}

// Register GetByHotelID
// @Tags Comments
// @Summary Get comments by hotel id
// @Description Get comments list by hotel uuid
// @Accept json
// @Produce json
// @Param hotel_id query string false "hotel uuid"
// @Success 200 {object} models.CommentsList
// @Router /comments/hotel/{hotel_id} [get]
func (h *commentsHandlers) GetByHotelID() echo.HandlerFunc {
	return func(c echo.Context) error {
		span, ctx := opentracing.StartSpanFromContext(c.Request().Context(), "commentsHandlers.GetByHotelID")
		defer span.Finish()

		hotelUUID, err := uuid.FromString(c.QueryParam("hotel_id"))
		if err != nil {
			h.logger.Error("uuid.FromString")
			return httpErrors.ErrorCtxResponse(c, err)
		}

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

		commentsList, err := h.commUC.GetByHotelID(ctx, hotelUUID, int64(page), int64(size))
		if err != nil {
			h.logger.Error("commUC.GetByHotelID")
			return httpErrors.ErrorCtxResponse(c, err)
		}

		return c.JSON(http.StatusOK, commentsList)
	}
}
