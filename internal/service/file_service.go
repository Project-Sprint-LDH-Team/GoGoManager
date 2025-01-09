package service

import (
	"context"
	"github.com/Project-Sprint-LDH-Team/GoGoManager/internal/models"
	"mime/multipart"
)

type FileService interface {
	UploadFile(ctx context.Context, file *multipart.FileHeader, userID uint, userEmail string) (*models.FileUploadResponse, error)
}

type fileService struct {
	storageService StorageService
}

func NewFileService(storageService StorageService) FileService {
	return &fileService{
		storageService: storageService,
	}
}

func (s *fileService) UploadFile(ctx context.Context, file *multipart.FileHeader, userID uint, userEmail string) (*models.FileUploadResponse, error) {
	// Validate file
	validator := &models.FileValidator{
		File:      file,
		MaxSize:   102400, // 100KiB in bytes
		MimeTypes: []string{"image/jpeg", "image/jpg", "image/png"},
	}

	if err := validator.Validate(); err != nil {
		return nil, err // Langsung return error validasi
	}

	// Upload to storage
	uri, err := s.storageService.UploadFile(ctx, file, userEmail)
	if err != nil {
		return nil, err
	}

	return &models.FileUploadResponse{
		URI: uri,
	}, nil
}
