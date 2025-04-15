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

type VideoService struct {
	videoDir   string
	coverDir   string
	lastUpdate time.Time
}

func NewVideoService(videoDir string, coverDir string) (*VideoService, error) {
	svc := &VideoService{
		videoDir: videoDir,
		coverDir: coverDir,
	}

	return svc, nil
}

func (s *VideoService) ListVideos() ([]db.Video, error) {
	videos, err := db.GetAllVideos()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve videos: %w", err)
	}

	return videos, nil
}
func (s *VideoService) ListVideosByGenre(genre string) ([]db.Video, error) {
	videos, err := db.GetVideosByGenre(genre)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve videos: %w", err)
	}

	return videos, nil
}

func (s *VideoService) GetGenres() ([]db.Genre, error) {
	genres, err := db.GetAllGenres()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve genres: %w", err)
	}

	return genres, nil
}
func (s *VideoService) StreamVideo(w http.ResponseWriter, r *http.Request, path string) error {
	fullPath := filepath.Join(s.videoDir, filepath.Clean(path))
	if !strings.HasPrefix(fullPath, s.videoDir) {
		return utils.ErrInvalidPath
	}
	file, err := os.Open(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			return utils.ErrNotFound
		}
		return err
	}
	defer file.Close()
	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}
	contentType := utils.GetContentType(fullPath)
	w.Header().Set("Content-Disposition", fmt.Sprintf("inline; filename=\"%s\"", filepath.Base(fullPath)))
	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Accept-Ranges", "bytes")
	if strings.HasSuffix(strings.ToLower(fullPath), ".ts") {
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")
	}

	fileSize := fileInfo.Size()
	rangeHeader := r.Header.Get("Range")
	if rangeHeader != "" {
		ranges, err := utils.ParseRangeHeader(rangeHeader, fileSize)
		if err != nil {
			return err
		}

		if len(ranges) > 0 {
			start, end := ranges[0].Start, ranges[0].End
			w.Header().Set("Content-Length", strconv.FormatInt(end-start+1, 10))
			w.Header().Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, end, fileSize))
			w.WriteHeader(http.StatusPartialContent)
			if _, err := file.Seek(start, 0); err != nil {
				return err
			}
			_, err = utils.CopyN(w, file, end-start+1)
			return err
		}
	}
	w.Header().Set("Content-Length", strconv.FormatInt(fileSize, 10))
	_, err = io.Copy(w, file)
	return err
}

func (s *VideoService) GetCoverImagePath(coverFilename string) string {
	if coverFilename == "" {
		return ""
	}
	return filepath.Join(s.coverDir, coverFilename)
}
func (s *VideoService) ServeCoverImage(w http.ResponseWriter, r *http.Request, filename string) error {
	if filename == "" {
		return utils.ErrNotFound
	}
	fullPath := filepath.Join(s.coverDir, filepath.Clean(filename))
	if !strings.HasPrefix(fullPath, s.coverDir) {
		return utils.ErrInvalidPath
	}
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return utils.ErrNotFound
	}
	contentType := "image/jpeg"
	if strings.HasSuffix(filename, ".png") {
		contentType = "image/png"
	} else if strings.HasSuffix(filename, ".gif") {
		contentType = "image/gif"
	} else if strings.HasSuffix(filename, ".webp") {
		contentType = "image/webp"
	}

	w.Header().Set("Content-Type", contentType)
	http.ServeFile(w, r, fullPath)
	return nil
}
