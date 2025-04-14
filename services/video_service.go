package services

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"DevMaan707/streamer/models"
	"DevMaan707/streamer/utils"
)

// VideoService handles video-related operations
type VideoService struct {
	videoDir    string
	videos      []models.Video
	videosMutex sync.RWMutex
	lastUpdate  time.Time
}

// NewVideoService creates a new video service
func NewVideoService(videoDir string) (*VideoService, error) {
	svc := &VideoService{
		videoDir: videoDir,
	}

	// Initialize video list
	_, err := svc.ListVideos(true)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize video list: %w", err)
	}

	// Start periodic updates
	go svc.periodicUpdate()

	return svc, nil
}

// ListVideos returns the list of videos
func (s *VideoService) ListVideos(refresh bool) ([]models.Video, error) {
	// If refresh requested or first load, update the list
	if refresh || time.Since(s.lastUpdate) > 5*time.Minute {
		if err := s.updateVideoList(); err != nil {
			return nil, err
		}
	}

	// Return the cached list
	s.videosMutex.RLock()
	defer s.videosMutex.RUnlock()

	// Return a copy to prevent external modification
	videos := make([]models.Video, len(s.videos))
	copy(videos, s.videos)

	return videos, nil
}

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

// updateVideoList refreshes the list of available videos
func (s *VideoService) updateVideoList() error {
	var videos []models.Video

	err := filepath.Walk(s.videoDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Check if it's a video file based on extension
		if utils.IsVideoFile(path) {
			relPath, err := filepath.Rel(s.videoDir, path)
			if err != nil {
				return err
			}

			videos = append(videos, models.Video{
				Name: info.Name(),
				Path: relPath,
				Size: info.Size(),
			})
		}

		return nil
	})

	if err != nil {
		return err
	}

	// Update the videos list in a thread-safe way
	s.videosMutex.Lock()
	s.videos = videos
	s.lastUpdate = time.Now()
	s.videosMutex.Unlock()

	return nil
}

// periodicUpdate refreshes the video list periodically
func (s *VideoService) periodicUpdate() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		if err := s.updateVideoList(); err != nil {
			fmt.Printf("Error updating video list: %v\n", err)
		}
	}
}
