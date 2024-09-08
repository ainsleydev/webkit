package storage

import (
	"context"
	"errors"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/ainsleydev/webkit/pkg/util/ptr"
)

// S3 implements the Storage interface for AWS S3.
type S3 struct {
	client s3ClientAPI
	bucket string
}

// S3Config holds the configuration for S3.
type S3Config struct {
	Bucket          string
	Region          string
	AccessKeyID     string
	SecretAccessKey string
	Endpoint        string
}

// s3ClientAPI defines the interface for S3 client operations we use
type s3ClientAPI interface {
	PutObject(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error)
	DeleteObject(ctx context.Context, params *s3.DeleteObjectInput, optFns ...func(*s3.Options)) (*s3.DeleteObjectOutput, error)
	ListObjectsV2(ctx context.Context, params *s3.ListObjectsV2Input, optFns ...func(*s3.Options)) (*s3.ListObjectsV2Output, error)
	GetObject(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error)
	HeadObject(ctx context.Context, params *s3.HeadObjectInput, optFns ...func(*s3.Options)) (*s3.HeadObjectOutput, error)
}

// NewS3Storage creates a new S3 instance with the provided configuration.
func NewS3Storage(ctx context.Context, cfg S3Config) (*S3, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	options := []func(*config.LoadOptions) error{
		config.WithRegion(cfg.Region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			cfg.AccessKeyID,
			cfg.SecretAccessKey,
			"",
		)),
	}

	awsCfg, err := config.LoadDefaultConfig(ctx, options...)
	if err != nil {
		return nil, err
	}

	var s3Options []func(*s3.Options)
	if cfg.Endpoint != "" {
		s3Options = append(s3Options, func(o *s3.Options) {
			o.BaseEndpoint = aws.String(cfg.Endpoint)
		})
	}

	return &S3{
		client: s3.NewFromConfig(awsCfg, s3Options...),
		bucket: cfg.Bucket,
	}, nil
}

// Validate validates the S3Config to ensure that all required fields are set.
func (cfg *S3Config) Validate() error {
	if cfg.Bucket == "" {
		return errors.New("bucket is required")
	}
	if cfg.Region == "" {
		return errors.New("region is required")
	}
	if cfg.AccessKeyID == "" {
		return errors.New("access key ID is required")
	}
	if cfg.SecretAccessKey == "" {
		return errors.New("secret access key is required")
	}
	return nil
}

func (s *S3) Upload(ctx context.Context, path string, content io.Reader) error {
	_, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(path),
		Body:   content,
	})
	return err
}

func (s *S3) Delete(ctx context.Context, path string) error {
	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(path),
	})
	return err
}

func (s *S3) List(ctx context.Context, prefix string) ([]string, error) {
	result, err := s.client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
		Bucket: aws.String(s.bucket),
		Prefix: aws.String(prefix),
	})
	if err != nil {
		return nil, err
	}

	var keys []string
	for _, obj := range result.Contents {
		keys = append(keys, *obj.Key)
	}
	return keys, nil
}

func (s *S3) Download(ctx context.Context, path string) (io.ReadCloser, error) {
	result, err := s.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(path),
	})
	if err != nil {
		return nil, err
	}
	return result.Body, nil
}

func (s *S3) Exists(ctx context.Context, path string) (bool, error) {
	_, err := s.client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(path),
	})
	if err != nil {
		return false, nil
	}
	return true, nil
}

func (s *S3) Stat(ctx context.Context, path string) (*FileInfo, error) {
	result, err := s.client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(path),
	})
	if err != nil {
		return nil, err
	}
	return &FileInfo{
		Size:         ptr.Int64(result.ContentLength),
		LastModified: ptr.Time(result.LastModified),
		IsDir:        false,
		ContentType:  ptr.String(result.ContentType),
	}, nil
}
