package images

import (
	"context"

	"github.com/AleksK1NG/hotels-mocroservices/images-microservice/internal/models"
)

type PgRepository interface {
	Create(ctx context.Context, msg *models.Image) (*models.Image, error)
}
