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

	"github.com/AleksK1NG/hotels-mocroservices/sessions/internal/models"
)

// sessionRedisRepo
type sessionRedisRepo struct {
	redis      *redis.Client
	prefix     string
	expiration time.Duration
}

func NewSessionRedisRepo(redis *redis.Client, prefix string, expiration time.Duration) *sessionRedisRepo {
	return &sessionRedisRepo{redis: redis, prefix: prefix, expiration: expiration}
}

func (s *sessionRedisRepo) CreateSession(ctx context.Context, userID uuid.UUID) (*models.Session, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "sessionRedisRepo.CreateSession")
	defer span.Finish()

	sess := &models.Session{
		SessionID: uuid.NewV4().String(),
		UserID:    userID,
	}

	sessBytes, err := json.Marshal(&sess)
	if err != nil {
		return nil, errors.Wrap(err, "sessionRepo.CreateSession.json.Marshal")
	}

	if err := s.redis.SetEX(ctx, s.createKey(sess.SessionID), string(sessBytes), s.expiration).Err(); err != nil {
		return nil, errors.Wrap(err, "sessionRepo.CreateSession.redis.SetEX")
	}

	return sess, nil
}

func (s *sessionRedisRepo) GetSessionByID(ctx context.Context, sessID string) (*models.Session, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "sessionRedisRepo.GetSessionByID")
	defer span.Finish()

	result, err := s.redis.Get(ctx, s.createKey(sessID)).Result()
	if err != nil {
		return nil, errors.Wrap(err, "sessionRepo.GetSessionByID.redis.Get")
	}

	var sess models.Session
	if err := json.Unmarshal([]byte(result), &sess); err != nil {
		return nil, errors.Wrap(err, "sessionRepo.GetSessionByID.json.Unmarshal")
	}
	return &sess, nil
}

func (s *sessionRedisRepo) DeleteSession(ctx context.Context, sessID string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "sessionRedisRepo.DeleteSession")
	defer span.Finish()

	if err := s.redis.Del(ctx, s.createKey(sessID)).Err(); err != nil {
		return errors.Wrap(err, "sessionRepo.DeleteSession.redis.Del")
	}
	return nil
}

func (s *sessionRedisRepo) createKey(sessionID string) string {
	return fmt.Sprintf("%s: %s", s.prefix, sessionID)
}
