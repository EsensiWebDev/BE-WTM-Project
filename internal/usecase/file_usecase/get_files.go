package file_usecase

import (
	"context"
	"wtm-backend/internal/domain"
	"wtm-backend/pkg/logger"
)

func (fu *FileUsecase) GetFiles(ctx context.Context, bucketName, objectName string) (domain.StreamableObject, error) {
	file, err := fu.fileStorage.GetFileObject(ctx, bucketName, objectName)
	if err != nil {
		logger.Error(ctx, "Error getting file", err.Error())
		return nil, err
	}

	return file, nil
}
