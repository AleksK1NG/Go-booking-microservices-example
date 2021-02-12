package usecase

import (
	"context"

	uuid "github.com/satori/go.uuid"

	"github.com/AleksK1NG/hotels-mocroservices/api-gateway/internal/comments"
	"github.com/AleksK1NG/hotels-mocroservices/api-gateway/internal/models"
	"github.com/AleksK1NG/hotels-mocroservices/api-gateway/pkg/logger"
	commentsService "github.com/AleksK1NG/hotels-mocroservices/api-gateway/proto/comments"
)

// CommentUseCase
type commentUseCase struct {
	logger      logger.Logger
	commService commentsService.CommentsServiceClient
	commRepo    comments.RedisRepository
}

// NewCommentUseCase
func NewCommentUseCase(logger logger.Logger, commService commentsService.CommentsServiceClient, commRepo comments.RedisRepository) *commentUseCase {
	return &commentUseCase{logger: logger, commService: commService, commRepo: commRepo}
}

func (c *commentUseCase) CreateComment(ctx context.Context, comment *models.Comment) (*models.Comment, error) {
	panic("implement me")
}

func (c *commentUseCase) GetCommByID(ctx context.Context, commentID uuid.UUID) (*models.Comment, error) {
	panic("implement me")
}

func (c *commentUseCase) UpdateComment(ctx context.Context, comment *models.Comment) (*models.Comment, error) {
	panic("implement me")
}

func (c *commentUseCase) GetByHotelID(ctx context.Context, hotelID uuid.UUID) (*models.CommentsList, error) {
	panic("implement me")
}
