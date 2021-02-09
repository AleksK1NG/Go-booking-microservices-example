package usecase

import (
	"context"

	"github.com/opentracing/opentracing-go"
	uuid "github.com/satori/go.uuid"

	"github.com/AleksK1NG/hotels-mocroservices/comments/internal/comment"
	"github.com/AleksK1NG/hotels-mocroservices/comments/internal/models"
	"github.com/AleksK1NG/hotels-mocroservices/comments/pkg/logger"
	"github.com/AleksK1NG/hotels-mocroservices/comments/pkg/utils"
	userService "github.com/AleksK1NG/hotels-mocroservices/comments/proto/user"
)

// CommUseCase
type commUseCase struct {
	commRepo   comment.PGRepository
	logger     logger.Logger
	userClient userService.UserServiceClient
}

// NewCommUseCase
func NewCommUseCase(commRepo comment.PGRepository, logger logger.Logger, userClient userService.UserServiceClient) *commUseCase {
	return &commUseCase{commRepo: commRepo, logger: logger, userClient: userClient}
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

	commentsList, err := c.commRepo.GetByHotelID(ctx, hotelID, query)
	if err != nil {
		return nil, err
	}

	uniqUserIDsMap := make(map[string]struct{}, len(commentsList.Comments))
	for _, comm := range commentsList.Comments {
		uniqUserIDsMap[comm.UserID.String()] = struct{}{}
	}

	userIDS := make([]string, 0, len(commentsList.Comments))
	for key, _ := range uniqUserIDsMap {
		userIDS = append(userIDS, key)
	}

	usersByIDs, err := c.userClient.GetUsersByIDs(ctx, &userService.GetByIDsReq{UsersIDs: userIDS})
	if err != nil {
		return nil, err
	}

	c.logger.Infof("USER CLIENT RESPONSE: %-v", usersByIDs)

	return commentsList, nil
}
