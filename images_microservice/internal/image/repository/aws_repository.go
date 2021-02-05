package repository

import (
	"bytes"
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"

	"github.com/AleksK1NG/hotels-mocroservices/images-microservice/config"
)

const (
	imagesBucket = "images"
)

type ImageAWSRepository struct {
	cfg *config.Config
	s3  *s3.S3
}

func NewImageAWSRepository(cfg *config.Config, s3 *s3.S3) *ImageAWSRepository {
	return &ImageAWSRepository{cfg: cfg, s3: s3}
}

func (i *ImageAWSRepository) PutObject(ctx context.Context, data []byte, fileType string) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ImageAWSRepository.PutObject")
	defer span.Finish()

	newFilename := uuid.NewV4().String()
	key := i.getFileKey(newFilename, fileType)

	object, err := i.s3.PutObjectWithContext(ctx, &s3.PutObjectInput{
		Body:   bytes.NewReader(data),
		Bucket: aws.String(imagesBucket),
		Key:    aws.String(key),
		ACL:    aws.String(s3.BucketCannedACLPublicRead),
	})
	if err != nil {
		return "", errors.Wrap(err, "s3.PutObjectWithContext")
	}

	log.Printf("object : %-v", object)

	return i.getFilePublicURL(key), err
}

func (i *ImageAWSRepository) GetObject(ctx context.Context, key string) (*s3.GetObjectOutput, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ImageAWSRepository.GetObject")
	defer span.Finish()

	obj, err := i.s3.GetObjectWithContext(ctx, &s3.GetObjectInput{
		Bucket: aws.String(imagesBucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, errors.Wrap(err, "s3.GetObjectWithContext")
	}

	return obj, nil
}

func (i *ImageAWSRepository) DeleteObject(ctx context.Context, key string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ImageAWSRepository.DeleteObject")
	defer span.Finish()

	_, err := i.s3.DeleteObjectWithContext(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(imagesBucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return errors.Wrap(err, "s3.DeleteObjectWithContext")
	}

	return nil
}

func (i *ImageAWSRepository) getFileKey(fileID string, fileType string) string {
	return fmt.Sprintf("%s.%s", fileID, fileType)
}

func (i *ImageAWSRepository) getFilePublicURL(key string) string {
	return i.cfg.AWS.S3EndPointMinio + "/" + imagesBucket + "/" + key
}
