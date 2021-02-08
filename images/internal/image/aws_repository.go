package image

import (
	"context"

	"github.com/aws/aws-sdk-go/service/s3"
)

type AWSRepository interface {
	PutObject(ctx context.Context, data []byte, fileType string) (string, error)
	GetObject(ctx context.Context, key string) (*s3.GetObjectOutput, error)
	DeleteObject(ctx context.Context, key string) error
}
