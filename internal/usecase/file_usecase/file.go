package file_usecase

import "wtm-backend/internal/domain"

type FileUsecase struct {
	fileStorage domain.StorageClient
}

func NewFileUsecase(fileStorage domain.StorageClient) *FileUsecase {
	return &FileUsecase{
		fileStorage: fileStorage,
	}
}
