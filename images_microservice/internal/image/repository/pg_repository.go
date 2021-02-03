package repository

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"

	"github.com/AleksK1NG/hotels-mocroservices/images-microservice/internal/models"
)

type ImagePGRepository struct {
	pgxPool *pgxpool.Pool
}

func NewImagePGRepository(pgxPool *pgxpool.Pool) *ImagePGRepository {
	return &ImagePGRepository{pgxPool: pgxPool}
}

func (i *ImagePGRepository) Create(ctx context.Context, msg *models.Image) (*models.Image, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ImagePGRepository.Create")
	defer span.Finish()

	var res models.Image
	if err := i.pgxPool.QueryRow(
		ctx,
		createImageQuery,
		msg.ImageURL,
		msg.IsUploaded,
	).Scan(&res.ImageID, &res.ImageURL, &res.IsUploaded, &res.CreatedAt); err != nil {
		return nil, errors.Wrap(err, "ImagePGRepository.Scan")
	}

	return &res, nil
}
