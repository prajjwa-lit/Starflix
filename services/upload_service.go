package services

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"DevMaan707/streamer/utils"
)

// UploadService handles file uploading functionality
type UploadService struct {
	uploadDir     string
	maxUploadSize int64
}

// NewUploadService creates a new upload service
func NewUploadService(uploadDir string, maxUploadSize int64) *UploadService {
	return &UploadService{
		uploadDir:     uploadDir,
		maxUploadSize: maxUploadSize,
	}
}

// HandleUpload processes a file upload from HTTP request
func (s *UploadService) HandleUpload(r *http.Request) (string, error) {
	// Parse our multipart form, 10 << 20 specifies a maximum
	// upload of 10 MB files
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		return "", fmt.Errorf("failed to parse form: %w", err)
	}

	// Get the file from form data
	file, header, err := r.FormFile("file")
	if err != nil {
		return "", fmt.Errorf("failed to get file: %w", err)
	}
	defer file.Close()

	// Limit file size
	if header.Size > s.maxUploadSize {
		return "", fmt.Errorf("file too large (max %d bytes)", s.maxUploadSize)
	}

	// Verify file is a video
	filename := header.Filename
	if !utils.IsVideoFile(filename) {
		return "", errors.New("only video files are allowed")
	}

	// Create safe filename
	safeName := utils.SafeFilename(filename)
	filePath := filepath.Join(s.uploadDir, safeName)

	// Create destination file
	dst, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create destination file: %w", err)
	}
	defer dst.Close()

	// Copy the uploaded file to the destination file
	_, err = io.Copy(dst, file)
	if err != nil {
		return "", fmt.Errorf("failed to save file: %w", err)
	}

	return safeName, nil
}
