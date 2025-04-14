package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"DevMaan707/streamer/services"
	"DevMaan707/streamer/utils"
)

// videoListHandler returns a handler for listing videos
func videoListHandler(svc *services.VideoService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Get videos
		videos, err := svc.ListVideos()
		if err != nil {
			log.Printf("Error listing videos: %v", err)
			http.Error(w, fmt.Sprintf("Failed to list videos: %v", err), http.StatusInternalServerError)
			return
		}

		// Return as JSON
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Cache-Control", "no-cache")

		if err := json.NewEncoder(w).Encode(videos); err != nil {
			log.Printf("Error encoding video response: %v", err)
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}
	}
}

// videoListByGenreHandler returns videos filtered by genre
func videoListByGenreHandler(svc *services.VideoService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Extract genre from path
		genre := strings.TrimPrefix(r.URL.Path, "/api/videos/genre/")
		if genre == "" {
			http.Error(w, "Genre required", http.StatusBadRequest)
			return
		}

		// Get videos
		videos, err := svc.ListVideosByGenre(genre)
		if err != nil {
			http.Error(w, "Failed to list videos", http.StatusInternalServerError)
			return
		}

		// Return as JSON
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Cache-Control", "no-cache")

		if err := json.NewEncoder(w).Encode(videos); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}
	}
}

// genreListHandler returns all available genres
func genreListHandler(svc *services.VideoService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Get genres
		genres, err := svc.GetGenres()
		if err != nil {
			http.Error(w, "Failed to list genres", http.StatusInternalServerError)
			return
		}

		// Return as JSON
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Cache-Control", "no-cache")

		if err := json.NewEncoder(w).Encode(genres); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}
	}
}

// videoStreamHandler returns a handler for streaming videos
func videoStreamHandler(svc *services.VideoService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Extract video path
		path := strings.TrimPrefix(r.URL.Path, "/videos/")
		if path == "" {
			http.Error(w, "Video path required", http.StatusBadRequest)
			return
		}

		// Stream the video
		err := svc.StreamVideo(w, r, path)
		if err != nil {
			if err == utils.ErrNotFound {
				http.Error(w, "Video not found", http.StatusNotFound)
			} else if err == utils.ErrInvalidPath {
				http.Error(w, "Invalid video path", http.StatusForbidden)
			} else {
				http.Error(w, "Error streaming video", http.StatusInternalServerError)
			}
		}
	}
}

// coverImageHandler serves cover images
func coverImageHandler(svc *services.VideoService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Extract cover image path
		path := strings.TrimPrefix(r.URL.Path, "/covers/")
		if path == "" {
			http.Error(w, "Image path required", http.StatusBadRequest)
			return
		}

		// Serve the cover image
		err := svc.ServeCoverImage(w, r, path)
		if err != nil {
			if err == utils.ErrNotFound {
				http.Error(w, "Image not found", http.StatusNotFound)
			} else if err == utils.ErrInvalidPath {
				http.Error(w, "Invalid image path", http.StatusForbidden)
			} else {
				http.Error(w, "Error serving image", http.StatusInternalServerError)
			}
		}
	}
}
