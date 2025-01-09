package service

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
	"mime/multipart"
	"path/filepath"
)

type StorageService interface {
	UploadFile(ctx context.Context, file *multipart.FileHeader, userEmail string) (string, error)
}

type s3StorageService struct {
	s3Client   *s3.Client
	bucketName string
}

func NewS3StorageService(s3Client *s3.Client, bucketName string) StorageService {
	return &s3StorageService{
		s3Client:   s3Client,
		bucketName: bucketName,
	}
}

func (s *s3StorageService) UploadFile(ctx context.Context, file *multipart.FileHeader, userEmail string) (string, error) {
	// Open file
	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer src.Close()

	// Generate unique filename
	fileExt := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("%s/%s%s", userEmail, uuid.New().String(), fileExt)

	// Upload to S3
	_, err = s.s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.bucketName),
		Key:         aws.String(filename),
		Body:        src,
		ContentType: aws.String(file.Header.Get("Content-Type")),
	})

	if err != nil {
		return "", fmt.Errorf("failed to upload file: %w", err)
	}

	// Return S3 URI
	return fmt.Sprintf("%s/%s", s.bucketName, filename), nil
}
