package hotels

import (
	"context"

	"github.com/AleksK1NG/hotels-mocroservices/api-gateway/internal/models"
)

// UseCase
type UseCase interface {
	GetHotels(ctx context.Context, page, size int64) (*models.HotelsListRes, error)
}
