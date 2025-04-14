package services

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"DevMaan707/streamer/db"
	"DevMaan707/streamer/utils"
)

// VideoService handles video-related operations
type VideoService struct {
	videoDir   string
	coverDir   string
	lastUpdate time.Time
}

// NewVideoService creates a new video service
func NewVideoService(videoDir string, coverDir string) (*VideoService, error) {
	svc := &VideoService{
		videoDir: videoDir,
		coverDir: coverDir,
	}

	return svc, nil
}

// ListVideos returns the list of videos
func (s *VideoService) ListVideos() ([]db.Video, error) {
	// Get videos from database
	videos, err := db.GetAllVideos()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve videos: %w", err)
	}

	return videos, nil
}

// ListVideosByGenre returns videos filtered by genre
func (s *VideoService) ListVideosByGenre(genre string) ([]db.Video, error) {
	videos, err := db.GetVideosByGenre(genre)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve videos: %w", err)
	}

	return videos, nil
}

// GetGenres returns all available genres
func (s *VideoService) GetGenres() ([]db.Genre, error) {
	genres, err := db.GetAllGenres()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve genres: %w", err)
	}

	return genres, nil
}

// StreamVideo streams a video file
func (s *VideoService) StreamVideo(w http.ResponseWriter, r *http.Request, path string) error {
	// Validate and clean the path
	fullPath := filepath.Join(s.videoDir, filepath.Clean(path))

	// Ensure the path doesn't try to navigate outside the videos directory
	if !strings.HasPrefix(fullPath, s.videoDir) {
		return utils.ErrInvalidPath
	}

	// Open the file
	file, err := os.Open(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			return utils.ErrNotFound
		}
		return err
	}
	defer file.Close()

	// Get file info
	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}

	// Set content disposition and type
	contentType := utils.GetContentType(fullPath)
	w.Header().Set("Content-Disposition", fmt.Sprintf("inline; filename=\"%s\"", filepath.Base(fullPath)))
	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Accept-Ranges", "bytes")

	// For .ts files, set additional headers for better streaming
	if strings.HasSuffix(strings.ToLower(fullPath), ".ts") {
		// These headers can help with MPEG-TS streaming
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")
	}

	fileSize := fileInfo.Size()

	// Handle range requests
	rangeHeader := r.Header.Get("Range")
	if rangeHeader != "" {
		// Process range header
		ranges, err := utils.ParseRangeHeader(rangeHeader, fileSize)
		if err != nil {
			return err
		}

		if len(ranges) > 0 {
			start, end := ranges[0].Start, ranges[0].End

			// Set headers for partial content
			w.Header().Set("Content-Length", strconv.FormatInt(end-start+1, 10))
			w.Header().Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, end, fileSize))
			w.WriteHeader(http.StatusPartialContent)

			// Seek to the start position
			if _, err := file.Seek(start, 0); err != nil {
				return err
			}

			// Stream the range
			_, err = utils.CopyN(w, file, end-start+1)
			return err
		}
	}

	// For regular requests (non-range), serve the entire file
	w.Header().Set("Content-Length", strconv.FormatInt(fileSize, 10))
	_, err = io.Copy(w, file)
	return err
}

// GetCoverImage returns the path to a cover image
func (s *VideoService) GetCoverImagePath(coverFilename string) string {
	if coverFilename == "" {
		return ""
	}
	return filepath.Join(s.coverDir, coverFilename)
}

// ServeCoverImage serves a cover image file
func (s *VideoService) ServeCoverImage(w http.ResponseWriter, r *http.Request, filename string) error {
	if filename == "" {
		return utils.ErrNotFound
	}

	// Validate and clean the path
	fullPath := filepath.Join(s.coverDir, filepath.Clean(filename))

	// Ensure the path doesn't try to navigate outside the covers directory
	if !strings.HasPrefix(fullPath, s.coverDir) {
		return utils.ErrInvalidPath
	}

	// Check if file exists
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return utils.ErrNotFound
	}

	// Set proper content type
	contentType := "image/jpeg" // default
	if strings.HasSuffix(filename, ".png") {
		contentType = "image/png"
	} else if strings.HasSuffix(filename, ".gif") {
		contentType = "image/gif"
	} else if strings.HasSuffix(filename, ".webp") {
		contentType = "image/webp"
	}

	w.Header().Set("Content-Type", contentType)

	// Serve the file
	http.ServeFile(w, r, fullPath)
	return nil
}
