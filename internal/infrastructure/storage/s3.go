package storage

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"io"
	"mime/multipart"
	"strings"
	"time"
	"wtm-backend/config"
	"wtm-backend/internal/domain"
	"wtm-backend/pkg/logger"
)

type S3Client struct {
	s3            *s3.Client
	bucketURL     string
	duration      time.Duration
	presignClient *s3.PresignClient
}

type S3Object struct {
	body        io.ReadCloser
	contentType string
	contentLen  int64
	filename    string
}

func NewS3Client(config *config.Config) (*S3Client, error) {
	ctx := context.Background()
	awsCfg, err := setConfigS3(&config.AWSConfig)
	if err != nil {
		logger.Error(ctx, "failed to set AWS S3 config", err.Error())
		return nil, err
	}
	s3Client := s3.NewFromConfig(awsCfg)
	presignClient := s3.NewPresignClient(s3Client)

	// Test koneksi ke bucket tertentu
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// HeadBucket untuk mengecek apakah bucket ada dan accessible
	_, err = s3Client.HeadBucket(ctx, &s3.HeadBucketInput{
		Bucket: aws.String(extractBucketName(config.AWSConfig.S3BucketURL)),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to access S3 bucket: %s", err.Error())
	}

	return &S3Client{
		s3:            s3Client,
		bucketURL:     config.AWSConfig.S3BucketURL,
		duration:      config.AWSConfig.S3PresignDuration,
		presignClient: presignClient,
	}, nil
}

func extractBucketName(bucketURL string) string {
	// Implementasi sederhana extract bucket name dari URL
	// Contoh: "https://my-bucket.s3.region.amazonaws.com" -> "my-bucket"
	if !strings.Contains(bucketURL, "//") {
		return bucketURL
	}
	return strings.Split(bucketURL, "//")[1] // sederhana, bisa diperbaiki
}

func setConfigS3(cfg *config.AWSConfig) (aws.Config, error) {
	ctx := context.Background()
	awsCfg, err := awsConfig.LoadDefaultConfig(context.TODO(),
		awsConfig.WithRegion(cfg.Region),
		awsConfig.WithCredentialsProvider(aws.CredentialsProviderFunc(
			func(ctx context.Context) (aws.Credentials, error) {
				return aws.Credentials{
					AccessKeyID:     cfg.AccessKeyID,
					SecretAccessKey: cfg.SecretAccessKey,
				}, nil
			},
		)),
	)
	awsCfg.BaseEndpoint = aws.String(cfg.Endpoint)
	if err != nil {
		logger.Error(ctx, "failed to load AWS config", err.Error())
		return aws.Config{}, err
	}

	return awsCfg, nil
}

func (s *S3Client) UploadFile(ctx context.Context, file multipart.File, fileHeader *multipart.FileHeader, bucketName string, objectName string) (string, error) {
	// Reset file pointer ke awal
	if seeker, ok := file.(io.Seeker); ok {
		seek, err := seeker.Seek(0, io.SeekStart)
		if err != nil {
			logger.Error(ctx, "failed to seek file", err.Error())
			return "", err
		}
		if seek != 0 {
			return "", errors.New("failed to reset file pointer to start")
		}
	}

	if file == nil || fileHeader == nil {
		return "", errors.New("file and fileHeader cannot be nil")
	}
	if bucketName == "" || objectName == "" {
		return "", errors.New("bucketName and objectName cannot be empty")
	}

	// Convert file to byte stream
	buf := new(bytes.Buffer)
	_, err := io.Copy(buf, file)
	if err != nil {
		return "", err
	}

	_, err = s.s3.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(bucketName),
		Key:         aws.String(objectName),
		Body:        bytes.NewReader(buf.Bytes()),
		ContentType: aws.String(fileHeader.Header.Get("Content-Type")),
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload file to S3: %s", err.Error())
	}

	return objectName, nil
}

func (s *S3Client) GetFile(ctx context.Context, bucketName, objectName string) (string, error) {
	presignedReq, err := s.presignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectName),
	}, s3.WithPresignExpires(s.duration))
	if err != nil {
		return "", err
	}
	return presignedReq.URL, nil
}

func (s *S3Client) GetFileObject(ctx context.Context, bucketName, objectName string) (domain.StreamableObject, error) {
	output, err := s.s3.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectName),
	})
	if err != nil {
		logger.Error(ctx, "failed to get file object from S3", err.Error())
		return nil, err
	}

	if output.ContentLength == nil || *output.ContentLength == 0 {
		logger.Error(ctx, "file not found or is empty", "objectName", objectName)
		return nil, errors.New("file not found or is empty")
	}

	return &S3Object{
		body:        output.Body,
		contentType: *output.ContentType,
		contentLen:  *output.ContentLength,
		filename:    objectName,
	}, nil
}

func (so *S3Object) Read(p []byte) (int, error) { return so.body.Read(p) }
func (so *S3Object) Close() error               { return so.body.Close() }
func (so *S3Object) GetContentType() string     { return so.contentType }
func (so *S3Object) GetContentLength() int64    { return so.contentLen }
func (so *S3Object) GetFilename() string        { return so.filename }
