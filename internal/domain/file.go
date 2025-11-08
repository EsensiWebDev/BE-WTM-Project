package domain

import "context"

type FileUsecase interface {
	GetFiles(ctx context.Context, bucketName, objectName string) (StreamableObject, error)
}
