package grpc

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/opentracing/opentracing-go"
	uuid "github.com/satori/go.uuid"

	"github.com/AleksK1NG/hotels-mocroservices/comments/config"
	"github.com/AleksK1NG/hotels-mocroservices/comments/internal/comment"
	"github.com/AleksK1NG/hotels-mocroservices/comments/internal/models"
	grpcErrors "github.com/AleksK1NG/hotels-mocroservices/comments/pkg/grpc_errors"
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
	span, ctx := opentracing.StartSpanFromContext(ctx, "CommentsService.CreateComment")
	defer span.Finish()

	comm, err := c.protoToModel(req)
	if err != nil {
		c.logger.Errorf("validate.StructCtx: %v", err)
		return nil, grpcErrors.ErrorResponse(err, err.Error())
	}

	if err := c.validate.StructCtx(ctx, comm); err != nil {
		c.logger.Errorf("validate.StructCtx: %v", err)
		return nil, grpcErrors.ErrorResponse(err, err.Error())
	}

	createdComm, err := c.commUC.Create(ctx, comm)
	if err != nil {
		c.logger.Errorf("commUC.Create: %v", err)
		return nil, grpcErrors.ErrorResponse(err, err.Error())
	}

	c.logger.Infof("CREATED: %-v", createdComm)

	return &commentsService.CreateCommentRes{Comment: createdComm.ToProto()}, nil
}

func (c *CommentsService) protoToModel(req *commentsService.CreateCommentReq) (*models.Comment, error) {
	hotelUUID, err := uuid.FromString(req.GetHotelID())
	if err != nil {
		return nil, err
	}
	userUUID, err := uuid.FromString(req.GetUserID())
	if err != nil {
		return nil, err
	}

	return &models.Comment{
		HotelID: hotelUUID,
		UserID:  userUUID,
		Message: req.GetMessage(),
		Photos:  req.GetPhotos(),
		Rating:  req.GetRating(),
	}, nil
}
