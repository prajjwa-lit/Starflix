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

type UploadService struct {
	uploadDir     string
	coverDir      string
	maxUploadSize int64
}

type UploadMetadata struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Genre       string `json:"genre"`
	ReleaseYear int    `json:"release_year"`
}

func NewUploadService(uploadDir string, coverDir string, maxUploadSize int64) *UploadService {
	return &UploadService{
		uploadDir:     uploadDir,
		coverDir:      coverDir,
		maxUploadSize: maxUploadSize,
	}
}
func checkDirPermissions(dir string) error {
	info, err := os.Stat(dir)
	if err != nil {
		return fmt.Errorf("failed to stat directory: %w", err)
	}
	if !info.IsDir() {
		return fmt.Errorf("%s is not a directory", dir)
	}
	testFile := filepath.Join(dir, ".test_write_permission")
	f, err := os.Create(testFile)
	if err != nil {
		return fmt.Errorf("failed to write to directory: %w", err)
	}
	f.Close()
	os.Remove(testFile)

	return nil
}

func (s *UploadService) HandleUpload(r *http.Request) (string, error) {
	log.Println("Starting file upload handling")
	log.Printf("Content-Length: %d", r.ContentLength)
	log.Printf("Transfer-Encoding: %v", r.TransferEncoding)
	log.Printf("X-Forwarded-For: %v", r.Header.Get("X-Forwarded-For"))
	maxMemory := int64(32 << 20)
	if err := r.ParseMultipartForm(maxMemory); err != nil {
		log.Printf("Failed to parse form: %v", err)
		return "", fmt.Errorf("failed to parse form: %w", err)
	}
	file, header, err := r.FormFile("file")
	if err != nil {
		log.Printf("Failed to get file: %v", err)
		return "", fmt.Errorf("failed to get file: %w", err)
	}
	defer file.Close()
	log.Printf("Received file: %s, size: %d bytes", header.Filename, header.Size)
	if header.Size > s.maxUploadSize {
		return "", fmt.Errorf("file too large (max %d bytes)", s.maxUploadSize)
	}
	filename := header.Filename
	if !utils.IsVideoFile(filename) {
		return "", errors.New("only video files are allowed")
	}

	safeName := utils.SafeFilename(header.Filename)
	filePath := filepath.Join(s.uploadDir, safeName)
	dst, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Printf("Failed to create file: %v", err)
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer dst.Close()
	written, err := io.Copy(dst, file)
	if err != nil {
		log.Printf("Failed to save file: %v", err)
		os.Remove(filePath)
		return "", fmt.Errorf("failed to save file: %w", err)
	}

	log.Printf("Successfully wrote %d bytes to %s", written, filePath)
	title := r.FormValue("title")
	if title == "" {
		title = strings.TrimSuffix(filepath.Base(filename), filepath.Ext(filename))
	}

	description := r.FormValue("description")
	genre := r.FormValue("genre")
	releaseYear := 0
	if yearStr := r.FormValue("release_year"); yearStr != "" {
		fmt.Sscanf(yearStr, "%d", &releaseYear)
	}
	var coverPath string
	coverFile, coverHeader, err := r.FormFile("cover_image")
	if err == nil && coverFile != nil {
		defer coverFile.Close()
		if utils.IsImageFile(coverHeader.Filename) {
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
	}

	return safeName, nil
}
func (s *UploadService) SaveCoverImage(file multipart.File, header *multipart.FileHeader, baseName string) (string, error) {
	if !utils.IsImageFile(header.Filename) {
		return "", errors.New("only image files are allowed for covers")
	}
	ext := filepath.Ext(header.Filename)
	safeName := "cover_" + utils.SafeFilename(baseName) + ext
	filePath := filepath.Join(s.coverDir, safeName)
	dst, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create destination file: %w", err)
	}
	defer dst.Close()
	_, err = io.Copy(dst, file)
	if err != nil {
		return "", fmt.Errorf("failed to save file: %w", err)
	}

	return safeName, nil
}
