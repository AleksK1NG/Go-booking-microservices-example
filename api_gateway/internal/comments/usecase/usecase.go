package usecase

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"

	"github.com/AleksK1NG/hotels-mocroservices/api-gateway/internal/comments"
	"github.com/AleksK1NG/hotels-mocroservices/api-gateway/internal/middlewares"
	"github.com/AleksK1NG/hotels-mocroservices/api-gateway/internal/models"
	httpErrors "github.com/AleksK1NG/hotels-mocroservices/api-gateway/pkg/http_errors"
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

// CreateComment
func (c *commentUseCase) CreateComment(ctx context.Context, comment *models.Comment) (*models.Comment, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "commentUseCase.CreateComment")
	defer span.Finish()

	ctxUser, ok := ctx.Value(middlewares.RequestCtxUser{}).(*models.UserResponse)
	if !ok || ctxUser == nil {
		return nil, errors.Wrap(httpErrors.Unauthorized, "ctx.Value user")
	}

	commentRes, err := c.commService.CreateComment(ctx, &commentsService.CreateCommentReq{
		HotelID: comment.HotelID.String(),
		UserID:  ctxUser.UserID.String(),
		Message: comment.Message,
		Photos:  comment.Photos,
		Rating:  comment.Rating,
	})
	if err != nil {
		return nil, errors.Wrap(err, "hotelsService.CreateHotel")
	}

	comm, err := models.CommentFromProto(commentRes.GetComment())
	if err != nil {
		return nil, errors.Wrap(err, "CommentFromProto")
	}

	return comm, nil
}

// GetCommByID
func (c *commentUseCase) GetCommByID(ctx context.Context, commentID uuid.UUID) (*models.Comment, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "commentUseCase.GetCommByID")
	defer span.Finish()

	cacheComm, err := c.commRepo.GetCommentByID(ctx, commentID)
	if err != nil {
		if err != redis.Nil {
			c.logger.Errorf("GetCommentByID: %v", err)
		}
	}
	if cacheComm != nil {
		return cacheComm, nil
	}

	commByID, err := c.commService.GetCommByID(ctx, &commentsService.GetCommByIDReq{CommentID: commentID.String()})
	if err != nil {
		return nil, errors.Wrap(err, "commService.GetCommByID")
	}

	comm, err := models.CommentFromProto(commByID.GetComment())
	if err != nil {
		return nil, errors.Wrap(err, "CommentFromProto")
	}

	if err := c.commRepo.SetComment(ctx, comm); err != nil {
		c.logger.Errorf("SetComment: %v", err)
	}

	return comm, nil
}

// UpdateComment
func (c *commentUseCase) UpdateComment(ctx context.Context, comment *models.Comment) (*models.Comment, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "commentUseCase.CreateComment")
	defer span.Finish()

	ctxUser, ok := ctx.Value(middlewares.RequestCtxUser{}).(*models.UserResponse)
	if !ok || ctxUser == nil {
		return nil, errors.Wrap(httpErrors.Unauthorized, "ctx.Value user")
	}

	if ctxUser.UserID != comment.UserID {
		return nil, errors.Wrap(httpErrors.WrongCredentials, "user is not owner")
	}

	commRes, err := c.commService.UpdateComment(ctx, &commentsService.UpdateCommReq{
		CommentID: comment.CommentID.String(),
		Message:   comment.Message,
		Photos:    comment.Photos,
		Rating:    comment.Rating,
	})
	if err != nil {
		return nil, errors.Wrap(err, "CommentFromProto")
	}

	comm, err := models.CommentFromProto(commRes.GetComment())
	if err != nil {
		return nil, errors.Wrap(err, "CommentFromProto")
	}

	if err := c.commRepo.SetComment(ctx, comm); err != nil {
		c.logger.Errorf("SetComment: %v", err)
	}

	return comm, nil
}

// GetByHotelID
func (c *commentUseCase) GetByHotelID(ctx context.Context, hotelID uuid.UUID, page, size int64) (*models.CommentsList, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "commentUseCase.GetByHotelID")
	defer span.Finish()

	res, err := c.commService.GetByHotelID(ctx, &commentsService.GetByHotelReq{
		HotelID: hotelID.String(),
		Page:    page,
		Size:    size,
	})
	if err != nil {
		return nil, errors.Wrap(err, "CommentFromProto")
	}

	commList := make([]*models.CommentFull, 0, len(res.Comments))
	for _, comment := range res.Comments {
		comm, err := models.CommentFullFromProto(comment)
		if err != nil {
			return nil, errors.Wrap(err, "CommentFullFromProto")
		}
		commList = append(commList, comm)
	}

	return &models.CommentsList{
		TotalCount: res.GetTotalCount(),
		TotalPages: res.GetTotalPages(),
		Page:       res.GetPage(),
		Size:       res.GetSize(),
		HasMore:    res.GetHasMore(),
		Comments:   commList,
	}, nil
}
