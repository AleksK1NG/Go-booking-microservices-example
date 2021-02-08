package usecase

import (
	"context"

	"github.com/opentracing/opentracing-go"
	uuid "github.com/satori/go.uuid"

	"github.com/AleksK1NG/hotels-mocroservices/comments/internal/comment"
	"github.com/AleksK1NG/hotels-mocroservices/comments/internal/models"
	"github.com/AleksK1NG/hotels-mocroservices/comments/pkg/logger"
	"github.com/AleksK1NG/hotels-mocroservices/comments/pkg/utils"
)

// CommUseCase
type commUseCase struct {
	commRepo comment.PGRepository
	logger   logger.Logger
}

// NewCommUseCase
func NewCommUseCase(commRepo comment.PGRepository, logger logger.Logger) *commUseCase {
	return &commUseCase{commRepo: commRepo, logger: logger}
}

// Create
func (c *commUseCase) Create(ctx context.Context, comment *models.Comment) (*models.Comment, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "commUseCase.Create")
	defer span.Finish()
	return c.commRepo.Create(ctx, comment)
}

// GetByID
func (c *commUseCase) GetByID(ctx context.Context, commentID uuid.UUID) (*models.Comment, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "commUseCase.GetByID")
	defer span.Finish()
	return c.commRepo.GetByID(ctx, commentID)
}

// Update
func (c *commUseCase) Update(ctx context.Context, comment *models.Comment) (*models.Comment, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "commUseCase.Update")
	defer span.Finish()
	return c.commRepo.Update(ctx, comment)
}

// GetByHotelID
func (c *commUseCase) GetByHotelID(ctx context.Context, hotelID uuid.UUID, query *utils.Pagination) (*models.CommentsList, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "commUseCase.GetByHotelID")
	defer span.Finish()
	return c.commRepo.GetByHotelID(ctx, hotelID, query)
}
