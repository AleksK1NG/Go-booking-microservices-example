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
	userExchange     = "user"
	resizeRoutingKey = "resize"
)

// UserUseCase
type UserUseCase struct {
	userPGRepo    user.PGRepository
	sessClient    sessionService.AuthorizationServiceClient
	redisRepo     user.RedisRepository
	log           logger.Logger
	amqpPublisher rabbitmq.Publisher
}

// NewUserUseCase
func NewUserUseCase(
	userPGRepo user.PGRepository,
	sessClient sessionService.AuthorizationServiceClient,
	redisRepo user.RedisRepository,
	log logger.Logger,
	amqpPublisher rabbitmq.Publisher,
) *UserUseCase {
	return &UserUseCase{
		userPGRepo:    userPGRepo,
		sessClient:    sessClient,
		redisRepo:     redisRepo,
		log:           log,
		amqpPublisher: amqpPublisher,
	}
}

// GetByID
func (u *UserUseCase) GetByID(ctx context.Context, userID uuid.UUID) (*models.UserResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UserUseCase.GetByID")
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
		return nil, errors.Wrap(err, "UserUseCase.userPGRepo.GetByID")
	}

	if err := u.redisRepo.SaveUser(ctx, userResponse); err != nil {
		u.log.Errorf("redisRepo.SaveUser: %v", err)
	}

	return userResponse, nil
}

// Register
func (u *UserUseCase) Register(ctx context.Context, user *models.User) (*models.UserResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UserUseCase.Register")
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

// Login
func (u *UserUseCase) Login(ctx context.Context, login models.Login) (*models.User, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UserUseCase.Login")
	defer span.Finish()

	userByEmail, err := u.userPGRepo.GetByEmail(ctx, login.Email)
	if err != nil {
		return nil, errors.Wrap(err, "userPGRepo.GetByEmail")
	}

	if err := userByEmail.ComparePasswords(login.Password); err != nil {
		return nil, errors.Wrap(err, "UserUseCase.ComparePasswords")
	}

	userByEmail.SanitizePassword()

	return userByEmail, nil
}

// CreateSession
func (u *UserUseCase) CreateSession(ctx context.Context, userID uuid.UUID) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UserUseCase.CreateSession")
	defer span.Finish()

	session, err := u.sessClient.CreateSession(ctx, &sessionService.CreateSessionRequest{UserID: userID.String()})
	if err != nil {
		return "", errors.Wrap(err, "sessClient.CreateSession")
	}

	return session.GetSession().GetSessionID(), err
}

// DeleteSession
func (u *UserUseCase) DeleteSession(ctx context.Context, sessionID string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UserUseCase.DeleteSession")
	defer span.Finish()

	_, err := u.sessClient.DeleteSession(ctx, &sessionService.DeleteSessionRequest{SessionID: sessionID})
	if err != nil {
		return errors.Wrap(err, "sessClient.DeleteSession")
	}

	return nil
}

// GetSessionByID
func (u *UserUseCase) GetSessionByID(ctx context.Context, sessionID string) (*models.Session, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UserUseCase.GetSessionByID")
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

// GetCSRFToken
func (u *UserUseCase) GetCSRFToken(ctx context.Context, sessionID string) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UserUseCase.GetCSRFToken")
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

// Update
func (u *UserUseCase) Update(ctx context.Context, user *models.UserUpdate) (*models.UserResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UserUseCase.Update")
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
		return nil, errors.Wrap(err, "UserUseCase.Update.userPGRepo.Update")
	}

	if err := u.redisRepo.SaveUser(ctx, userResponse); err != nil {
		u.log.Errorf("redisRepo.SaveUser: %v", err)
	}

	return userResponse, nil
}

func (u *UserUseCase) UpdateUploadedAvatar(ctx context.Context, delivery amqp.Delivery) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UserUseCase.UpdateUploadedAvatar")
	defer span.Finish()

	u.log.Infof("USER UPDATE UPLOADED AVATAR **************** : %v", delivery.Body)

	var img models.Image
	if err := json.Unmarshal(delivery.Body, &img); err != nil {
		return errors.Wrap(err, "UserUseCase.UpdateUploadedAvatar.json.Unmarshal")
	}

	userUUID, ok := delivery.Headers["user_uuid"]
	if !ok {
		return errors.Wrap(errors.New("not ok"), "UserUseCase.UpdateUploadedAvatar.json.Unmarshal")
	}

	uuidDromStr, err := uuid.FromString(userUUID.(string))
	if err != nil {
		return errors.Wrap(err, "UserUseCase.UpdateUploadedAvatar.uuid.FromString")
	}

	u.log.Infof("USER UNMARSHAL **************** : %v", img)
	created, err := u.userPGRepo.UpdateAvatar(ctx, models.UploadedImageMsg{
		ImageID:    img.ImageID,
		UserID:     uuidDromStr,
		ImageURL:   img.ImageURL,
		IsUploaded: img.IsUploaded,
	})
	if err != nil {
		return err
	}

	u.log.Infof("USER CREATED AVATAR WOWOWOWOWO **************** : %v", created)

	return nil
}

const (
	imagesExchange = "images"
	resizeKey      = "resize_image_key"
)

func (u *UserUseCase) UpdateAvatar(ctx context.Context, data *models.UpdateAvatarMsg) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UserUseCase.UpdateAvatar")
	defer span.Finish()

	headers := make(amqp.Table)
	headers["user_uuid"] = data.UserID.String()

	u.log.Infof("AMQP headers: %v", headers)

	u.log.Infof("PUBLISH UpdateAvatar USER ***************** %-v", headers)
	if err := u.amqpPublisher.Publish(ctx, imagesExchange, resizeKey, data.ContentType, headers, data.Body); err != nil {
		return errors.Wrap(err, "UserUseCase.UpdateUploadedAvatar.amqpPublisher.Publish")
	}

	return nil
}
