package utils

import (
	"errors"
	"mime"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	ErrNotFound    = errors.New("file not found")
	ErrInvalidPath = errors.New("invalid path")
)

func IsVideoFile(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".mp4", ".webm", ".ogg", ".mov", ".avi", ".mkv", ".flv", ".ts":
		return true
	default:
		return false
	}
}
func IsImageFile(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".jpg", ".jpeg", ".png", ".gif", ".webp":
		return true
	default:
		return false
	}
}

func GetContentType(path string) string {
	ext := strings.ToLower(filepath.Ext(path))
	if ext == ".ts" {
		return "video/mp2t"
	}

	contentType := mime.TypeByExtension(ext)
	if contentType == "" {
		contentType = "application/octet-stream"
	}
	return contentType
}
func SafeFilename(filename string) string {
	ext := filepath.Ext(filename)
	name := filename[:len(filename)-len(ext)]
	re := regexp.MustCompile(`[^\w\-.]`)
	name = re.ReplaceAllString(name, "_")
	if len(name) > 200 {
		name = name[:200]
	}

	return name + ext
}
