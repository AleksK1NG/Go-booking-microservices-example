package image

import (
	"context"

	uuid "github.com/satori/go.uuid"

	"github.com/AleksK1NG/hotels-mocroservices/images-microservice/internal/models"
)

type PgRepository interface {
	Create(ctx context.Context, msg *models.Image) (*models.Image, error)
	GetImageByID(ctx context.Context, imageID uuid.UUID) (*models.Image, error)
}
