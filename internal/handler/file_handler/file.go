package file_handler

import "wtm-backend/internal/domain"

type FileHandler struct {
	fileUsecase domain.FileUsecase
}

func NewFileHandler(fileUsecase domain.FileUsecase) *FileHandler {
	return &FileHandler{
		fileUsecase: fileUsecase,
	}
}
