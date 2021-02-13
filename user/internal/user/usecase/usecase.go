package usecase

import (
	"context"
	"encoding/json"

	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"github.com/streadway/amqp"

	"github.com/AleksK1NG/hotels-mocroservices/user/internal/middlewares"
	"github.com/AleksK1NG/hotels-mocroservices/user/internal/models"
	"github.com/AleksK1NG/hotels-mocroservices/user/internal/user"
	"github.com/AleksK1NG/hotels-mocroservices/user/internal/user/delivery/rabbitmq"
	httpErrors "github.com/AleksK1NG/hotels-mocroservices/user/pkg/http_errors"
	"github.com/AleksK1NG/hotels-mocroservices/user/pkg/logger"
	sessionService "github.com/AleksK1NG/hotels-mocroservices/user/proto/session"
)

const (
	imagesExchange = "images"
	resizeKey      = "resize_image_key"
	userUUIDHeader = "user_uuid"
)

type userUseCase struct {
	userPGRepo    user.PGRepository
	sessClient    sessionService.AuthorizationServiceClient
	redisRepo     user.RedisRepository
	log           logger.Logger
	amqpPublisher rabbitmq.Publisher
}

func NewUserUseCase(
	userPGRepo user.PGRepository,
	sessClient sessionService.AuthorizationServiceClient,
	redisRepo user.RedisRepository,
	log logger.Logger,
	amqpPublisher rabbitmq.Publisher,
) *userUseCase {
	return &userUseCase{
		userPGRepo:    userPGRepo,
		sessClient:    sessClient,
		redisRepo:     redisRepo,
		log:           log,
		amqpPublisher: amqpPublisher,
	}
}

func (u *userUseCase) GetByID(ctx context.Context, userID uuid.UUID) (*models.UserResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "userUseCase.GetByID")
	defer span.Finish()

	cachedUser, err := u.redisRepo.GetUserByID(ctx, userID)
	if err != nil {
		u.log.Errorf("redisRepo.GetUserByID: %v", err)
	}
	if cachedUser != nil {
		return cachedUser, nil
	}

	userResponse, err := u.userPGRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, errors.Wrap(err, "userUseCase.userPGRepo.GetByID")
	}

	if err := u.redisRepo.SaveUser(ctx, userResponse); err != nil {
		u.log.Errorf("redisRepo.SaveUser: %v", err)
	}

	return userResponse, nil
}

func (u *userUseCase) Register(ctx context.Context, user *models.User) (*models.UserResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "userUseCase.Register")
	defer span.Finish()

	if err := user.PrepareCreate(); err != nil {
		return nil, errors.Wrap(err, "user.PrepareCreate")
	}

	created, err := u.userPGRepo.Create(ctx, user)
	if err != nil {
		return nil, errors.Wrap(err, "userPGRepo.Create")
	}

	return created, err
}

func (u *userUseCase) Login(ctx context.Context, login models.Login) (*models.User, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "userUseCase.Login")
	defer span.Finish()

	userByEmail, err := u.userPGRepo.GetByEmail(ctx, login.Email)
	if err != nil {
		return nil, errors.Wrap(err, "userPGRepo.GetByEmail")
	}

	if err := userByEmail.ComparePasswords(login.Password); err != nil {
		return nil, errors.Wrap(err, "userUseCase.ComparePasswords")
	}

	userByEmail.SanitizePassword()

	return userByEmail, nil
}

func (u *userUseCase) CreateSession(ctx context.Context, userID uuid.UUID) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "userUseCase.CreateSession")
	defer span.Finish()

	session, err := u.sessClient.CreateSession(ctx, &sessionService.CreateSessionRequest{UserID: userID.String()})
	if err != nil {
		return "", errors.Wrap(err, "sessClient.CreateSession")
	}

	return session.GetSession().GetSessionID(), err
}

func (u *userUseCase) DeleteSession(ctx context.Context, sessionID string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "userUseCase.DeleteSession")
	defer span.Finish()

	_, err := u.sessClient.DeleteSession(ctx, &sessionService.DeleteSessionRequest{SessionID: sessionID})
	if err != nil {
		return errors.Wrap(err, "sessClient.DeleteSession")
	}

	return nil
}

func (u *userUseCase) GetSessionByID(ctx context.Context, sessionID string) (*models.Session, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "userUseCase.GetSessionByID")
	defer span.Finish()

	sessionByID, err := u.sessClient.GetSessionByID(ctx, &sessionService.GetSessionByIDRequest{SessionID: sessionID})
	if err != nil {
		return nil, errors.Wrap(err, "sessClient.GetSessionByID")
	}

	sess := &models.Session{}
	sess, err = sess.FromProto(sessionByID.GetSession())
	if err != nil {
		return nil, errors.Wrap(err, "sess.FromProto")
	}

	return sess, nil
}

func (u *userUseCase) GetCSRFToken(ctx context.Context, sessionID string) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "userUseCase.GetCSRFToken")
	defer span.Finish()

	csrfToken, err := u.sessClient.CreateCsrfToken(
		ctx,
		&sessionService.CreateCsrfTokenRequest{CsrfTokenInput: &sessionService.CsrfTokenInput{SessionID: sessionID}},
	)
	if err != nil {
		return "", errors.Wrap(err, "sessClient.CreateCsrfToken")
	}

	return csrfToken.GetCsrfToken().GetToken(), nil
}

func (u *userUseCase) Update(ctx context.Context, user *models.UserUpdate) (*models.UserResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "userUseCase.Update")
	defer span.Finish()

	ctxUser, ok := ctx.Value(middlewares.RequestCtxUser{}).(*models.UserResponse)
	if !ok {
		return nil, errors.Wrap(httpErrors.Unauthorized, "ctx.Value user")
	}

	if ctxUser.UserID != user.UserID || *ctxUser.Role != models.RoleAdmin {
		return nil, errors.Wrap(httpErrors.WrongCredentials, "user is not owner or admin")
	}

	userResponse, err := u.userPGRepo.Update(ctx, user)
	if err != nil {
		return nil, errors.Wrap(err, "userUseCase.Update.userPGRepo.Update")
	}

	if err := u.redisRepo.SaveUser(ctx, userResponse); err != nil {
		u.log.Errorf("redisRepo.SaveUser: %v", err)
	}

	return userResponse, nil
}

func (u *userUseCase) UpdateUploadedAvatar(ctx context.Context, delivery amqp.Delivery) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "userUseCase.UpdateUploadedAvatar")
	defer span.Finish()

	var img models.Image
	if err := json.Unmarshal(delivery.Body, &img); err != nil {
		return errors.Wrap(err, "UpdateUploadedAvatar.json.Unmarshal")
	}

	userUUID, ok := delivery.Headers[userUUIDHeader].(string)
	if !ok {
		return errors.Wrap(httpErrors.InvalidUUID, "delivery.Headers")
	}

	uid, err := uuid.FromString(userUUID)
	if err != nil {
		return errors.Wrap(err, "uuid.FromString")
	}

	created, err := u.userPGRepo.UpdateAvatar(ctx, models.UploadedImageMsg{
		ImageID:    img.ImageID,
		UserID:     uid,
		ImageURL:   img.ImageURL,
		IsUploaded: img.IsUploaded,
	})
	if err != nil {
		return err
	}

	u.log.Infof("UpdateUploadedAvatar: %s", created.Avatar)

	return nil
}

func (u *userUseCase) UpdateAvatar(ctx context.Context, data *models.UpdateAvatarMsg) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "userUseCase.UpdateAvatar")
	defer span.Finish()

	headers := make(amqp.Table, 1)
	headers[userUUIDHeader] = data.UserID.String()
	if err := u.amqpPublisher.Publish(
		ctx,
		imagesExchange,
		resizeKey,
		data.ContentType,
		headers,
		data.Body,
	); err != nil {
		return errors.Wrap(err, "UpdateUploadedAvatar.Publish")
	}

	u.log.Infof("Publish UpdateAvatar %-v", headers)
	return nil
}

func (u *userUseCase) GetUsersByIDs(ctx context.Context, userIDs []string) ([]*models.UserResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "userUseCase.GetUsersByIDs")
	defer span.Finish()
	return u.userPGRepo.GetUsersByIDs(ctx, userIDs)
}
