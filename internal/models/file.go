package models

import (
	"encoding/json"
	"fmt"
	"mime/multipart"
	"time"
)

type File struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"not null" json:"user_id"`
	Filename  string    `gorm:"size:255;not null" json:"filename"`
	URI       string    `gorm:"size:255;not null" json:"uri"`
	FileType  string    `gorm:"size:50;not null" json:"file_type"`
	FileSize  int64     `gorm:"not null" json:"file_size"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// FileUploadResponse sesuai dengan contract API
type FileUploadResponse struct {
	URI string `json:"uri"`
}

// FileValidationError dengan implementasi interface error
type FileValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// Implementasi interface error
func (e FileValidationError) Error() string {
	errBytes, _ := json.Marshal(e)
	return string(errBytes)
}

// Custom error type untuk multiple validation errors
type FileValidationErrors []FileValidationError

func (e FileValidationErrors) Error() string {
	errBytes, _ := json.Marshal(e)
	return string(errBytes)
}

// FileValidator untuk validasi file upload
type FileValidator struct {
	File      *multipart.FileHeader
	MaxSize   int64
	MimeTypes []string
}

func (v *FileValidator) Validate() error {
	var errors FileValidationErrors

	// Validate file size (max 100KiB)
	if v.File.Size > v.MaxSize {
		errors = append(errors, FileValidationError{
			Field:   "file",
			Message: fmt.Sprintf("File size exceeds maximum limit of %d KiB", v.MaxSize/1024),
		})
	}

	// Validate mime type
	contentType := v.File.Header.Get("Content-Type")
	isValidType := false
	for _, mimeType := range v.MimeTypes {
		if contentType == mimeType {
			isValidType = true
			break
		}
	}
	if !isValidType {
		errors = append(errors, FileValidationError{
			Field:   "file",
			Message: "File type must be jpeg, jpg, or png",
		})
	}

	if len(errors) > 0 {
		return errors
	}

	return nil
}
