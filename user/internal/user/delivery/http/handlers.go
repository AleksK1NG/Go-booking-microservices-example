package http

import (
	"bytes"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"sync"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/opentracing/opentracing-go"
	uuid "github.com/satori/go.uuid"

	"github.com/AleksK1NG/hotels-mocroservices/user/config"
	"github.com/AleksK1NG/hotels-mocroservices/user/internal/middlewares"
	"github.com/AleksK1NG/hotels-mocroservices/user/internal/models"
	"github.com/AleksK1NG/hotels-mocroservices/user/internal/user"
	httpErrors "github.com/AleksK1NG/hotels-mocroservices/user/pkg/http_errors"
	"github.com/AleksK1NG/hotels-mocroservices/user/pkg/logger"
)

const (
	csrfHeader  = "X-CSRF-Token"
	maxFileSize = 1024 * 1024 * 10
)

// UserHandlers
type UserHandlers struct {
	cfg      *config.Config
	group    *echo.Group
	userUC   user.UseCase
	logger   logger.Logger
	validate *validator.Validate
	mw       *middlewares.MiddlewareManager
}

// NewUserHandlers
func NewUserHandlers(
	group *echo.Group,
	userUC user.UseCase,
	logger logger.Logger,
	validate *validator.Validate,
	cfg *config.Config,
	mw *middlewares.MiddlewareManager,
) *UserHandlers {
	return &UserHandlers{group: group, userUC: userUC, logger: logger, validate: validate, cfg: cfg, mw: mw}
}

// Register godoc
// @Summary Register new user
// @Description register new user account, returns user data and session
// @Accept json
// @Produce json
// @Param data body models.User true "user data"
// @Success 201 {object} models.UserResponse
// @Router /user/register [post]
func (h *UserHandlers) Register() echo.HandlerFunc {
	return func(c echo.Context) error {
		span, ctx := opentracing.StartSpanFromContext(c.Request().Context(), "auth.Register")
		defer span.Finish()

		var u models.User
		if err := c.Bind(&u); err != nil {
			h.logger.Errorf("c.Bind: %v", err)
			return httpErrors.ErrorCtxResponse(c, err)
		}

		if err := h.validate.StructCtx(ctx, &u); err != nil {
			h.logger.Errorf("validate.StructCtx: %v", err)
			return httpErrors.ErrorCtxResponse(c, err)
		}

		regUser, err := h.userUC.Register(ctx, &u)
		if err != nil {
			h.logger.Errorf("UserHandlers.Register.userUC.Register: %v", err)
			return httpErrors.ErrorCtxResponse(c, err)
		}

		sessionID, err := h.userUC.CreateSession(ctx, regUser.UserID)
		if err != nil {
			h.logger.Errorf("UserHandlers.userUC.CreateSession: %v", err)
			return httpErrors.ErrorCtxResponse(c, err)
		}

		c.SetCookie(&http.Cookie{
			Name:     h.cfg.HttpServer.SessionCookieName,
			Value:    sessionID,
			Path:     "/",
			HttpOnly: true,
			Expires:  time.Now().Add(time.Duration(h.cfg.HttpServer.CookieLifeTime) * time.Minute),
		})

		return c.JSON(http.StatusCreated, regUser)
	}
}

// Login godoc
// @Summary Login user
// @Description login user, returns user data and session
// @Accept json
// @Produce json
// @Param data body models.Login true "email and password"
// @Success 200 {object} models.UserResponse
// @Router /user/login [post]
func (h *UserHandlers) Login() echo.HandlerFunc {
	return func(c echo.Context) error {
		span, ctx := opentracing.StartSpanFromContext(c.Request().Context(), "user.Login")
		defer span.Finish()

		var login models.Login
		if err := c.Bind(&login); err != nil {
			h.logger.Errorf("c.Bind: %v", err)
			return httpErrors.ErrorCtxResponse(c, err)
		}

		if err := h.validate.StructCtx(ctx, &login); err != nil {
			h.logger.Errorf("validate.StructCtx: %v", err)
			return httpErrors.ErrorCtxResponse(c, err)
		}

		userResponse, err := h.userUC.Login(ctx, login)
		if err != nil {
			h.logger.Errorf("UserHandlers.userUC.Login: %v", err)
			return httpErrors.ErrorCtxResponse(c, err)
		}

		sessionID, err := h.userUC.CreateSession(ctx, userResponse.UserID)
		if err != nil {
			h.logger.Errorf("UserHandlers.Login.CreateSession: %v", err)
			return httpErrors.ErrorCtxResponse(c, err)
		}

		c.SetCookie(&http.Cookie{
			Name:     h.cfg.HttpServer.SessionCookieName,
			Value:    sessionID,
			Path:     "/",
			HttpOnly: true,
			Expires:  time.Now().Add(time.Duration(h.cfg.HttpServer.CookieLifeTime) * time.Minute),
		})

		return c.JSON(http.StatusOK, userResponse)
	}
}

// Logout godoc
// @Summary Logout user
// @Description Logout user, return no content
// @Accept json
// @Produce json
// @Success 204 ""
// @Router /user/logout [post]
func (h *UserHandlers) Logout() echo.HandlerFunc {
	return func(c echo.Context) error {
		span, ctx := opentracing.StartSpanFromContext(c.Request().Context(), "user.Logout")
		defer span.Finish()

		cookie, err := c.Cookie(h.cfg.HttpServer.SessionCookieName)
		if err != nil {
			if errors.Is(err, http.ErrNoCookie) {
				h.logger.Errorf("UserHandlers.Logout.http.ErrNoCookie: %v", err)
				return httpErrors.ErrorCtxResponse(c, err)
			}
			h.logger.Errorf("UserHandlers.Logout.c.Cookie: %v", err)
			return httpErrors.ErrorCtxResponse(c, err)
		}

		if err := h.userUC.DeleteSession(ctx, cookie.Value); err != nil {
			h.logger.Errorf("UserHandlers.userUC.DeleteSession: %v", err)
			return httpErrors.ErrorCtxResponse(c, err)
		}

		c.SetCookie(&http.Cookie{
			Name:   h.cfg.HttpServer.SessionCookieName,
			Value:  "",
			Path:   "/",
			MaxAge: -1,
		})

		return c.NoContent(http.StatusNoContent)
	}
}

// GetMe godoc
// @Summary Get current user data
// @Description Get current user data, required session cookie
// @Accept json
// @Produce json
// @Success 200 {object} models.UserResponse
// @Router /user/me [get]
func (h *UserHandlers) GetMe() echo.HandlerFunc {
	return func(c echo.Context) error {
		span, ctx := opentracing.StartSpanFromContext(c.Request().Context(), "user.GetMe")
		defer span.Finish()

		userResponse, ok := ctx.Value(middlewares.RequestCtxUser{}).(*models.UserResponse)
		if !ok {
			h.logger.Error("invalid middleware user ctx")
			return httpErrors.ErrorCtxResponse(c, httpErrors.WrongCredentials)
		}

		return c.JSON(http.StatusOK, userResponse)
	}
}

// GetCSRFToken godoc
// @Summary Get csrf token
// @Description Get csrf token, required session
// @Accept json
// @Produce json
// @Success 204 ""
// @Router /user/csrf [get]
func (h *UserHandlers) GetCSRFToken() echo.HandlerFunc {
	return func(c echo.Context) error {
		span, ctx := opentracing.StartSpanFromContext(c.Request().Context(), "user.GetCSRFToken")
		defer span.Finish()

		userSession, ok := ctx.Value(middlewares.RequestCtxSession{}).(*models.Session)
		if !ok {
			h.logger.Error("invalid middleware session ctx")
			return httpErrors.ErrorCtxResponse(c, httpErrors.WrongCredentials)
		}

		csrfToken, err := h.userUC.GetCSRFToken(ctx, userSession.SessionID)
		if err != nil {
			h.logger.Error("userUC.GetCSRFToken")
			return httpErrors.ErrorCtxResponse(c, err)
		}

		c.Response().Header().Set(csrfHeader, csrfToken)

		return c.NoContent(http.StatusOK)
	}
}

func (h *UserHandlers) Update() echo.HandlerFunc {
	return func(c echo.Context) error {
		span, ctx := opentracing.StartSpanFromContext(c.Request().Context(), "user.Update")
		defer span.Finish()

		userID := c.Param("id")
		if userID == "" {
			h.logger.Error("invalid user id param")
			return httpErrors.ErrorCtxResponse(c, httpErrors.BadRequest)
		}

		userUUID, err := uuid.FromString(userID)
		if err != nil {
			h.logger.Error("invalid user uuid")
			return httpErrors.ErrorCtxResponse(c, err)
		}

		var updUser models.UserUpdate
		if err := c.Bind(&updUser); err != nil {
			h.logger.Errorf("c.Bind: %v", err)
			return httpErrors.ErrorCtxResponse(c, err)
		}
		updUser.UserID = userUUID

		if err := h.validate.StructCtx(ctx, updUser); err != nil {
			h.logger.Errorf("UserHandlers.validate.StructCtx: %v", err)
			return httpErrors.ErrorCtxResponse(c, err)
		}

		userResponse, err := h.userUC.Update(ctx, &updUser)
		if err != nil {
			h.logger.Errorf("UserHandlers.userUC.Update: %v", err)
			return httpErrors.ErrorCtxResponse(c, err)
		}

		return c.JSON(http.StatusOK, userResponse)
	}
}

func (h *UserHandlers) Delete() echo.HandlerFunc {
	panic("implement me")
}

// GetUserByID godoc
// @Summary Get user by id
// @Description Get user data by id
// @Accept json
// @Produce json
// @Param id path int false "user uuid"
// @Success 200 {object} models.UserResponse
// @Router /user/{id} [get]
func (h *UserHandlers) GetUserByID() echo.HandlerFunc {
	return func(c echo.Context) error {
		span, ctx := opentracing.StartSpanFromContext(c.Request().Context(), "user.GetUserByID")
		defer span.Finish()

		userID := c.Param("id")
		if userID == "" {
			h.logger.Error("invalid user id param")
			return httpErrors.ErrorCtxResponse(c, httpErrors.BadRequest)
		}

		userUUID, err := uuid.FromString(userID)
		if err != nil {
			h.logger.Error("uuid.FromString")
			return httpErrors.ErrorCtxResponse(c, err)
		}

		userResponse, err := h.userUC.GetByID(ctx, userUUID)
		if err != nil {
			h.logger.Error("userUC.GetByID")
			return httpErrors.ErrorCtxResponse(c, err)
		}

		return c.JSON(http.StatusOK, userResponse)
	}
}

// UpdateAvatar godoc
// @Summary Update user avatar
// @Description Upload user avatar image
// @Accept mpfd
// @Produce json
// @Param id path int false "user uuid"
// @Success 200 {object} models.UserResponse
// @Router /user/{id}/avatar [put]
func (h *UserHandlers) UpdateAvatar() echo.HandlerFunc {
	bufferPool := &sync.Pool{New: func() interface{} {
		return &bytes.Buffer{}
	}}
	return func(c echo.Context) error {
		span, ctx := opentracing.StartSpanFromContext(c.Request().Context(), "user.UpdateAvatar")
		defer span.Finish()

		userResponse, ok := ctx.Value(middlewares.RequestCtxUser{}).(*models.UserResponse)
		if !ok {
			h.logger.Error("invalid middleware user ctx")
			return httpErrors.ErrorCtxResponse(c, httpErrors.WrongCredentials)
		}

		if err := c.Request().ParseMultipartForm(maxFileSize); err != nil {
			h.logger.Error("c.ParseMultipartForm")
			return httpErrors.ErrorCtxResponse(c, err)
		}

		c.Request().Body = http.MaxBytesReader(c.Response(), c.Request().Body, maxFileSize)
		defer c.Request().Body.Close()

		formFile, _, err := c.Request().FormFile("avatar")
		if err != nil {
			h.logger.Error("c.FormFile")
			return httpErrors.ErrorCtxResponse(c, err)
		}

		fileType, err := h.checkAvatar(formFile)
		if err != nil {
			h.logger.Error("h.checkAvatar")
			return httpErrors.ErrorCtxResponse(c, err)
		}

		buf, ok := bufferPool.Get().(*bytes.Buffer)
		if !ok {
			h.logger.Error("bufferPool.Get")
			return httpErrors.ErrorCtxResponse(c, httpErrors.InternalServerError)
		}
		defer bufferPool.Put(buf)
		buf.Reset()

		if _, err := io.Copy(buf, formFile); err != nil {
			h.logger.Error("io.Copy")
			return httpErrors.ErrorCtxResponse(c, err)
		}

		data := &models.UpdateAvatarMsg{
			UserID:      userResponse.UserID,
			ContentType: fileType,
			Body:        buf.Bytes(),
		}

		if err := h.userUC.UpdateAvatar(ctx, data); err != nil {
			h.logger.Error("h.userUC.UpdateAvatar")
			return httpErrors.ErrorCtxResponse(c, err)
		}

		return c.NoContent(http.StatusOK)
	}
}

func (h *UserHandlers) checkAvatar(file multipart.File) (string, error) {
	fileHeader := make([]byte, maxFileSize)
	ContentType := ""
	if _, err := file.Read(fileHeader); err != nil {
		return ContentType, err
	}

	if _, err := file.Seek(0, 0); err != nil {
		return ContentType, err
	}

	count, err := file.Seek(0, 2)
	if err != nil {
		return ContentType, err
	}
	if count > maxFileSize {
		return ContentType, err
	}
	if _, err := file.Seek(0, 0); err != nil {
		return ContentType, err
	}
	ContentType = http.DetectContentType(fileHeader)

	if ContentType != "image/jpg" && ContentType != "image/png" && ContentType != "image/jpeg" {
		return ContentType, err
	}

	return ContentType, nil
}
