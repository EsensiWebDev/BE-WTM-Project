package domain

import (
	"context"
	"io"
	"mime/multipart"
)

type StorageClient interface {
	UploadFile(ctx context.Context, file multipart.File, fileHeader *multipart.FileHeader, bucketName string, objectName string) (string, error)
	GetFile(ctx context.Context, bucketName, objectName string) (string, error)
	GetFileObject(ctx context.Context, bucketName, objectName string) (StreamableObject, error)
}

type StreamableObject interface {
	io.ReadCloser
	GetContentType() string
	GetContentLength() int64
	GetFilename() string
}
