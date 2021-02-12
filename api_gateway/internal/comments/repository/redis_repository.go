package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"

	"github.com/AleksK1NG/hotels-mocroservices/api-gateway/internal/models"
)

const (
	prefix     = "comments:"
	expiration = time.Second * 3600
)

// CommRedisRepository
type commRedisRepository struct {
	redisConn *redis.Client
}

// NewCommRedisRepository
func NewCommRedisRepository(redisConn *redis.Client) *commRedisRepository {
	return &commRedisRepository{redisConn: redisConn}
}

// GetCommentByID
func (c *commRedisRepository) GetCommentByID(ctx context.Context, commentID uuid.UUID) (*models.Comment, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "commRedisRepository.GetCommentByID")
	defer span.Finish()

	result, err := c.redisConn.Get(ctx, c.createKey(commentID)).Bytes()
	if err != nil {
		return nil, errors.Wrap(err, "commRedisRepository.GetCommentByID.Get")
	}

	var res models.Comment
	if err := json.Unmarshal(result, &res); err != nil {
		return nil, errors.Wrap(err, "json.Unmarshal")
	}
	return &res, nil
}

// SetComment
func (c *commRedisRepository) SetComment(ctx context.Context, comment *models.Comment) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "commRedisRepository.GetCommentByID")
	defer span.Finish()

	commBytes, err := json.Marshal(comment)
	if err != nil {
		return errors.Wrap(err, "commRedisRepository.Marshal")
	}

	if err := c.redisConn.SetEX(ctx, c.createKey(comment.CommentID), string(commBytes), expiration).Err(); err != nil {
		return errors.Wrap(err, "commRedisRepository.GetCommentByID.SetEX")
	}

	return nil
}

// DeleteComment
func (c *commRedisRepository) DeleteComment(ctx context.Context, commentID uuid.UUID) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "commRedisRepository.DeleteComment")
	defer span.Finish()

	if err := c.redisConn.Del(ctx, c.createKey(commentID)).Err(); err != nil {
		return errors.Wrap(err, "commRedisRepository.DeleteComment.Del")
	}

	return nil
}

func (c *commRedisRepository) createKey(commID uuid.UUID) string {
	return fmt.Sprintf("%s: %s", prefix, commID.String())
}
