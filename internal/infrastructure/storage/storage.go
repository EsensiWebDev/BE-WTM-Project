package storage

import (
	"context"
	"fmt"
	"wtm-backend/config"
	"wtm-backend/internal/domain"
	"wtm-backend/pkg/logger"
)

type MultiStorageClient struct {
	Minio         domain.StorageClient
	S3            domain.StorageClient
	ActiveStorage domain.StorageClient
}

func NewMultiStorageClient(config *config.Config) (*MultiStorageClient, error) {
	ctx := context.Background()
	minioClient, err := NewMinioClient(config)
	if err != nil {
		logger.Error(ctx, "Failed to initialize minio client", err.Error())
	}
	s3Client, err := NewS3Client(config)
	if err != nil {
		logger.Error(ctx, "failed to initialize s3 client", err.Error())
	}
	if minioClient == nil && s3Client == nil {
		return nil, fmt.Errorf("failed to initialize minio client and s3 client")
	}
	if s3Client != nil {
		logger.Info(ctx, "Using S3 as the active storage client")
		return &MultiStorageClient{Minio: minioClient, S3: s3Client, ActiveStorage: s3Client}, nil
	}

	logger.Info(ctx, "Using Minio as the active storage client")
	return &MultiStorageClient{Minio: minioClient, S3: s3Client, ActiveStorage: minioClient}, nil
}
