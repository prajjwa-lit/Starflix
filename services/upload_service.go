package services

import (
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"DevMaan707/streamer/db"
	"DevMaan707/streamer/utils"
)

// UploadService handles file uploading functionality
type UploadService struct {
	uploadDir     string
	coverDir      string
	maxUploadSize int64
}

// UploadMetadata contains metadata for an uploaded video
type UploadMetadata struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Genre       string `json:"genre"`
	ReleaseYear int    `json:"release_year"`
}

// NewUploadService creates a new upload service
func NewUploadService(uploadDir string, coverDir string, maxUploadSize int64) *UploadService {
	return &UploadService{
		uploadDir:     uploadDir,
		coverDir:      coverDir,
		maxUploadSize: maxUploadSize,
	}
}

// HandleUpload processes a file upload from HTTP request
func (s *UploadService) HandleUpload(r *http.Request) (string, error) {
	// Parse multipart form
	err := r.ParseMultipartForm(s.maxUploadSize)
	if err != nil {
		return "", fmt.Errorf("failed to parse form: %w", err)
	}

	// Get the video file from form data
	file, header, err := r.FormFile("file")
	if err != nil {
		return "", fmt.Errorf("failed to get video file: %w", err)
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

	// Now handle the metadata and cover image
	title := r.FormValue("title")
	if title == "" {
		// Use filename without extension as default title
		title = strings.TrimSuffix(filepath.Base(filename), filepath.Ext(filename))
	}

	description := r.FormValue("description")
	genre := r.FormValue("genre")

	// Try to parse release year
	releaseYear := 0
	if yearStr := r.FormValue("release_year"); yearStr != "" {
		fmt.Sscanf(yearStr, "%d", &releaseYear)
	}

	// Handle cover image if provided
	var coverPath string
	coverFile, coverHeader, err := r.FormFile("cover_image")
	if err == nil && coverFile != nil {
		defer coverFile.Close()

		// Verify it's an image file
		if utils.IsImageFile(coverHeader.Filename) {
			// Create a safe filename for the cover image
			coverFilename := "cover_" + utils.SafeFilename(filename) + filepath.Ext(coverHeader.Filename)
			coverFullPath := filepath.Join(s.coverDir, coverFilename)

			coverDst, err := os.Create(coverFullPath)
			if err == nil {
				defer coverDst.Close()
				_, err = io.Copy(coverDst, coverFile)
				if err == nil {
					coverPath = coverFilename
				} else {
					log.Printf("Failed to save cover image: %v", err)
				}
			} else {
				log.Printf("Failed to create cover image file: %v", err)
			}
		}
	}

	// Store video metadata in the database
	video := &db.Video{
		Filename:    filepath.Base(safeName),
		Title:       title,
		Description: description,
		Genre:       genre,
		ReleaseYear: releaseYear,
		CoverImage:  coverPath,
		FilePath:    safeName,
		FileSize:    header.Size,
	}

	err = db.InsertVideo(video)
	if err != nil {
		log.Printf("Warning: Failed to store video metadata: %v", err)
		// We'll still return success since the file was uploaded
	}

	return safeName, nil
}

// SaveCoverImage saves a cover image file
func (s *UploadService) SaveCoverImage(file multipart.File, header *multipart.FileHeader, baseName string) (string, error) {
	// Verify file is an image
	if !utils.IsImageFile(header.Filename) {
		return "", errors.New("only image files are allowed for covers")
	}

	// Create safe filename for the cover image
	ext := filepath.Ext(header.Filename)
	safeName := "cover_" + utils.SafeFilename(baseName) + ext
	filePath := filepath.Join(s.coverDir, safeName)

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
