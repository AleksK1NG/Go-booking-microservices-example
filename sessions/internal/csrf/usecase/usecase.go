package usecase

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"io"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"

	"github.com/AleksK1NG/hotels-mocroservices/sessions/internal/csrf"
	"github.com/AleksK1NG/hotels-mocroservices/sessions/internal/models"
)

// CSRF usecase
type CsrfUC struct {
	csrfRepo       csrf.RedisRepository
	secretTokenKey string
	csrfExpire     int
}

// NewCsrfUC
func NewCsrfUC(csrfRepo csrf.RedisRepository, secretTokenKey string, csrfExpire int) *CsrfUC {
	return &CsrfUC{csrfRepo: csrfRepo, secretTokenKey: secretTokenKey, csrfExpire: csrfExpire}
}

// CreateToken
func (c *CsrfUC) CreateToken(ctx context.Context, sesID string, timeStamp int64) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CsrfUC.CreateToken")
	defer span.Finish()

	block, err := aes.NewCipher([]byte(c.secretTokenKey))
	if err != nil {
		return "", errors.Wrap(err, "CsrfUC.CreateToken.aes.NewCipher")
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", errors.Wrap(err, "CsrfUC.CreateToken.cipher.NewGCM")
	}

	nonce := make([]byte, aesgcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", errors.Wrap(err, "CsrfUC.CreateToken.io.ReadFull")
	}

	td := &models.CsrfToken{SessionID: sesID, Timestamp: timeStamp}
	data, err := json.Marshal(td)
	if err != nil {
		return "", errors.Wrap(err, "CsrfUC.CreateToken.json.Marshal")
	}
	ciphertext := aesgcm.Seal(nil, nonce, data, nil)

	res := append([]byte(nil), nonce...)
	res = append(res, ciphertext...)

	token := base64.StdEncoding.EncodeToString(res)

	return token, nil
}

// CheckToken
func (c *CsrfUC) CheckToken(ctx context.Context, sesID string, token string) (bool, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CsrfUC.CheckToken")
	defer span.Finish()

	ciphertext, _ := base64.StdEncoding.DecodeString(token)

	block, err := aes.NewCipher([]byte(c.secretTokenKey))
	if err != nil {
		return false, errors.Wrap(err, "CsrfUC.CheckToken.aes.NewCipher")
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return false, err
	}

	nonceSize := aesgcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return false, errors.New("ciphertext < nonceSize")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

	plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return false, errors.Wrap(err, "CsrfUC.CheckToken.aesgcm.Open")
	}

	CsrfTok := models.CsrfToken{}
	if err := json.Unmarshal(plaintext, &CsrfTok); err != nil {
		return false, errors.Wrap(err, "CsrfUC.CheckToken.json.Unmarshal")
	}

	if time.Now().Unix()-CsrfTok.Timestamp >
		int64(time.Duration(c.csrfExpire)*time.Minute) {
		return false, errors.New("token expired")
	}

	expected := models.CsrfToken{SessionID: sesID, Timestamp: CsrfTok.Timestamp}

	err = c.csrfRepo.Check(ctx, token)

	if CsrfTok != expected || err != nil {
		return false, errors.New("token expired")
	}

	if err := c.csrfRepo.Create(ctx, token); err != nil {
		return false, errors.Wrap(err, "CsrfUC.CheckToken.csrfRepo.Create")
	}

	return true, nil
}
