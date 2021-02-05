package repository

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"

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
func (i *ImagePGRepository) GetImageByID(ctx context.Context, imageID uuid.UUID) (*models.Image, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ImagePGRepository.GetImageByID")
	defer span.Finish()

	var img models.Image
	if err := i.pgxPool.QueryRow(ctx, getImageByIDQuery, imageID).Scan(
		&img.ImageID,
		&img.ImageURL,
		&img.IsUploaded,
		&img.CreatedAt,
		&img.UpdatedAt,
	); err != nil {
		return nil, errors.Wrap(err, "ImagePGRepository.Scan")
	}

	return &img, nil
}
