package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"mime/multipart"
	"strings"
	"time"
	"wtm-backend/config"
	"wtm-backend/internal/domain"
	"wtm-backend/pkg/logger"
)

const publicString = "public"
const bucketNameCheck = "check-bucket"

type MinioClient struct {
	client       *minio.Client
	baseURL      string
	durationFile time.Duration
}

type MinioObject struct {
	obj  *minio.Object
	stat minio.ObjectInfo
}

func NewMinioClient(config *config.Config) (*MinioClient, error) {
	ctx := context.Background()
	endpoint := fmt.Sprintf("%s:%s", config.MinioHost, config.MinioPort)
	accessKey := config.MinioRootUser
	secretKey := config.MinioRootPassword

	// Inisialisasi MinIO client
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: false, // true jika pakai HTTPS
	})
	if err != nil {
		logger.Fatal("Failed to initialize Minio client", err.Error())
		return nil, err
	}

	// Cek koneksi dengan bucket khusus
	if err := client.MakeBucket(context.Background(), bucketNameCheck, minio.MakeBucketOptions{}); err != nil {
		if minio.ToErrorResponse(err).Code != "BucketAlreadyOwnedByYou" {
			logger.Error(ctx, "Failed to connect to Minio bucket", err.Error())
			return nil, err
		}
	}
	// Hapus bucket khusus yang digunakan untuk pengecekan
	if err := client.RemoveBucket(context.Background(), bucketNameCheck); err != nil {
		logger.Error(ctx, "Failed to remove check bucket", err.Error())
		return nil, err
	}

	logger.Info(ctx, "Minio client initialized successfully")

	baseURL := config.Host + "/api/files"

	return &MinioClient{client: client, baseURL: baseURL, durationFile: config.DurationAccessFileMinio}, nil
}

func (m *MinioClient) UploadFile(ctx context.Context, file multipart.File, fileHeader *multipart.FileHeader, bucketName string, objectName string) (string, error) {
	if err := m.ensureBucket(ctx, bucketName); err != nil {
		logger.Error(ctx, "Error to ensure bucket", err.Error())
		return "", err
	}

	contentType := fileHeader.Header.Get("Content-Type")
	_, err := m.client.PutObject(ctx, bucketName, objectName, file, fileHeader.Size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		logger.Error(ctx, "Error to upload file", err.Error())
		return "", err
	}

	return objectName, nil
}

func (m *MinioClient) GetFile(ctx context.Context, bucketName, objectName string) (string, error) {
	// Cek apakah file ada
	_, err := m.client.StatObject(ctx, bucketName, objectName, minio.StatObjectOptions{})
	if err != nil {
		logger.Error(ctx, "File not found", err.Error())
		return "", err
	}

	// Jika bucket publik
	if strings.Contains(bucketName, publicString) {
		return fmt.Sprintf("%s/%s/%s", m.baseURL, bucketName, objectName), nil
	}

	// Jika private, generate presigned URL
	url, err := m.client.PresignedGetObject(ctx, bucketName, objectName, m.durationFile, nil)
	if err != nil {
		logger.Error(ctx, "Error to generate accessible URL", err.Error())
		return "", err
	}
	return url.String(), nil
}

func (m *MinioClient) GetFileObject(ctx context.Context, bucketName, objectName string) (domain.StreamableObject, error) {
	obj, err := m.client.GetObject(ctx, bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		logger.Error(ctx, "Error to get file object", err.Error())
		return nil, err
	}

	stat, err := obj.Stat()
	if err != nil {
		logger.Error(ctx, "Error to get file object stat", err.Error())
		return nil, err
	}

	if stat.Size == 0 {
		logger.Error(ctx, "File object is empty", "objectName", objectName)
		return nil, fmt.Errorf("file object is empty")
	}

	return &MinioObject{obj: obj, stat: stat}, nil
}

func (mo *MinioObject) Read(p []byte) (int, error) { return mo.obj.Read(p) }
func (mo *MinioObject) Close() error               { return mo.obj.Close() }
func (mo *MinioObject) GetContentType() string     { return mo.stat.ContentType }
func (mo *MinioObject) GetContentLength() int64    { return mo.stat.Size }
func (mo *MinioObject) GetFilename() string        { return mo.stat.Key }

func (m *MinioClient) ensureBucket(ctx context.Context, bucketName string) error {
	exists, err := m.client.BucketExists(ctx, bucketName)
	if err != nil {
		return err
	}
	if !exists {
		err = m.client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
		if err != nil {
			logger.Error(ctx, "Error to create bucket", err.Error())
			return err
		}

		if strings.Contains(bucketName, publicString) {
			if err := m.setPublicReadPolicy(ctx, bucketName); err != nil {
				logger.Error(ctx, "Error to set read policy", err.Error())
				return err
			}
		}
	}
	return nil
}

func (m *MinioClient) setPublicReadPolicy(ctx context.Context, bucketName string) error {
	policy := map[string]interface{}{
		"Version": "2012-10-17",
		"Statement": []map[string]interface{}{
			{
				"Effect":    "Allow",
				"Principal": map[string]interface{}{"AWS": "*"},
				"Action":    []string{"s3:GetObject"},
				"Resource":  fmt.Sprintf("arn:aws:s3:::%s/*", bucketName),
			},
		},
	}

	policyBytes, err := json.Marshal(policy)
	if err != nil {
		logger.Error(ctx, "Error to marshal policy", err.Error())
		return err
	}

	return m.client.SetBucketPolicy(ctx, bucketName, string(policyBytes))
}
