package s3

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/smithy-go"

	logger_pkg "boilerplate/internal/pkg/logger"
)

type Client interface {
	CreateBucket(ctx context.Context, bucket string) (bool, error)
	UploadFile(ctx context.Context, path string, content *bytes.Reader) error
	// DownloadFile загружает файл по указанному пути
	// ВАЖНО: Вызывающий должен закрыть возвращаемый io.ReadCloser
	DownloadFile(ctx context.Context, path string) (io.ReadCloser, error)
	DeleteFile(ctx context.Context, path string) error
}

type client struct {
	logger   logger_pkg.Logger
	s3client *s3.Client
	bucket   string
}

func NewClient(ctx context.Context, host, port, accessKey, secretKey, bucket string, opts ...option) (Client, error) {
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")),
		config.WithRegion("us-east-1"),
	)
	if err != nil {
		return nil, fmt.Errorf("s3 load config: %w", err)
	}

	s3client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(fmt.Sprintf("http://%s:%s", host, port))
		o.UsePathStyle = true
	})

	c := &client{
		s3client: s3client,
		bucket:   bucket,
	}

	for _, opt := range opts {
		opt(c)
	}

	return c, nil
}

func (c *client) CreateBucket(ctx context.Context, bucket string) (bool, error) {
	if c.logger != nil {
		c.logger.DebugKV(ctx, "s3 create bucket", "bucket", bucket)
	}

	// Пробуем создать bucket. Если он уже существует - игнорируем ошибку
	_, err := c.s3client.CreateBucket(ctx, &s3.CreateBucketInput{
		Bucket: aws.String(bucket),
	})
	if err != nil {
		var aerr smithy.APIError
		if errors.As(err, &aerr) && (aerr.ErrorCode() == "BucketAlreadyOwnedByYou" || aerr.ErrorCode() == "BucketAlreadyExists") {
			return false, nil
		}
		return false, fmt.Errorf("s3 create bucket: %w", err)
	}

	return true, nil
}

func (c *client) UploadFile(ctx context.Context, path string, content *bytes.Reader) error {
	if c.logger != nil {
		c.logger.DebugKV(ctx, "s3 upload", "path", path)
	}

	_, err := c.s3client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(c.bucket),
		Key:    aws.String(path),
		Body:   content,
	})
	if err != nil {
		return fmt.Errorf("s3 put object: %w", err)
	}

	return nil
}

func (c *client) DownloadFile(ctx context.Context, path string) (io.ReadCloser, error) {
	if c.logger != nil {
		c.logger.DebugKV(ctx, "s3 download", "path", path)
	}

	result, err := c.s3client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(c.bucket),
		Key:    aws.String(path),
	})
	if err != nil {
		return nil, fmt.Errorf("s3 get object: %w", err)
	}

	return result.Body, nil
}

func (c *client) DeleteFile(ctx context.Context, path string) error {
	if c.logger != nil {
		c.logger.DebugKV(ctx, "s3 delete", "path", path)
	}

	_, err := c.s3client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(c.bucket),
		Key:    aws.String(path),
	})
	if err != nil {
		return fmt.Errorf("s3 delete object: %w", err)
	}

	return nil
}
