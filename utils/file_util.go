package utils

import (
	"errors"
	"mime"
	"path/filepath"
	"regexp"
	"strings"
)

// Common errors
var (
	ErrNotFound    = errors.New("file not found")
	ErrInvalidPath = errors.New("invalid path")
)

// IsVideoFile checks if a file is a video based on extension
func IsVideoFile(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".mp4", ".webm", ".ogg", ".mov", ".avi", ".mkv", ".flv", ".ts":
		return true
	default:
		return false
	}
}

// GetContentType returns the MIME type based on file extension
func GetContentType(path string) string {
	ext := strings.ToLower(filepath.Ext(path))

	// Special case for .ts files which might not be in the mime package
	if ext == ".ts" {
		return "video/mp2t"
	}

	contentType := mime.TypeByExtension(ext)
	if contentType == "" {
		// Default to octet-stream if MIME type can't be determined
		contentType = "application/octet-stream"
	}
	return contentType
}

// SafeFilename creates a safe filename by removing potentially dangerous characters
func SafeFilename(filename string) string {
	// Extract the file extension
	ext := filepath.Ext(filename)
	name := filename[:len(filename)-len(ext)]

	// Replace unsafe characters with underscores
	re := regexp.MustCompile(`[^\w\-.]`)
	name = re.ReplaceAllString(name, "_")

	// Ensure the name isn't excessively long
	if len(name) > 200 {
		name = name[:200]
	}

	return name + ext
}
