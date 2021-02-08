package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"

	"github.com/AleksK1NG/hotels-mocroservices/images-microservice/config"
)

// Init new AWS S3 session
func NewS3Session(cfg *config.Config) *s3.S3 {
	// s3Config := &aws.Config{
	// 	Credentials:      credentials.NewStaticCredentials("minio", "minio", ""),
	// 	Endpoint:         aws.String("http://localhost:9000"),
	// 	Region:           aws.String("us-east-1"),
	// 	DisableSSL:       aws.Bool(true),
	// 	S3ForcePathStyle: aws.Bool(true),
	// }
	//
	// newSession, err := session.NewSession(s3Config)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	//
	// s3Client := s3.New(newSession)
	return s3.New(session.Must(session.NewSession(&aws.Config{
		Credentials:      credentials.NewStaticCredentials("minio123", "minio123", ""),
		Region:           aws.String(cfg.AWS.S3Region),
		Endpoint:         aws.String(cfg.AWS.S3EndPoint),
		DisableSSL:       aws.Bool(cfg.AWS.DisableSSL),
		S3ForcePathStyle: aws.Bool(cfg.AWS.S3ForcePathStyle),
	})))

}
