package usecase

import (
	"context"

	"github.com/opentracing/opentracing-go"

	"github.com/AleksK1NG/hotels-mocroservices/comments/internal/comment"
	"github.com/AleksK1NG/hotels-mocroservices/comments/internal/models"
	"github.com/AleksK1NG/hotels-mocroservices/comments/pkg/logger"
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
