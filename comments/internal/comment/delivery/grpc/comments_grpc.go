package grpc

import (
	"context"

	"github.com/go-playground/validator/v10"

	"github.com/AleksK1NG/hotels-mocroservices/comments/config"
	"github.com/AleksK1NG/hotels-mocroservices/comments/internal/comment"
	"github.com/AleksK1NG/hotels-mocroservices/comments/pkg/logger"
	"github.com/AleksK1NG/hotels-mocroservices/comments/proto"
)

// CommentsService
type CommentsService struct {
	commUC   comment.UseCase
	logger   logger.Logger
	cfg      *config.Config
	validate *validator.Validate
}

// NewCommentsService
func NewCommentsService(commUC comment.UseCase, logger logger.Logger, cfg *config.Config, validate *validator.Validate) *CommentsService {
	return &CommentsService{commUC: commUC, logger: logger, cfg: cfg, validate: validate}
}

// CreateComment
func (c *CommentsService) CreateComment(ctx context.Context, req *commentsService.CreateCommentReq) (*commentsService.CreateCommentRes, error) {
	panic("implement me")
}
