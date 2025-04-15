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

func videoListHandler(svc *services.VideoService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		videos, err := svc.ListVideos()
		if err != nil {
			log.Printf("Error listing videos: %v", err)
			http.Error(w, fmt.Sprintf("Failed to list videos: %v", err), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Cache-Control", "no-cache")

		if err := json.NewEncoder(w).Encode(videos); err != nil {
			log.Printf("Error encoding video response: %v", err)
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}
	}
}
func videoListByGenreHandler(svc *services.VideoService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		genre := strings.TrimPrefix(r.URL.Path, "/api/videos/genre/")
		if genre == "" {
			http.Error(w, "Genre required", http.StatusBadRequest)
			return
		}

		videos, err := svc.ListVideosByGenre(genre)
		if err != nil {
			http.Error(w, "Failed to list videos", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Cache-Control", "no-cache")

		if err := json.NewEncoder(w).Encode(videos); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}
	}
}

func genreListHandler(svc *services.VideoService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		genres, err := svc.GetGenres()
		if err != nil {
			http.Error(w, "Failed to list genres", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Cache-Control", "no-cache")

		if err := json.NewEncoder(w).Encode(genres); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}
	}
}

func videoStreamHandler(svc *services.VideoService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		path := strings.TrimPrefix(r.URL.Path, "/videos/")
		if path == "" {
			http.Error(w, "Video path required", http.StatusBadRequest)
			return
		}
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

func coverImageHandler(svc *services.VideoService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		path := strings.TrimPrefix(r.URL.Path, "/covers/")
		if path == "" {
			http.Error(w, "Image path required", http.StatusBadRequest)
			return
		}
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
