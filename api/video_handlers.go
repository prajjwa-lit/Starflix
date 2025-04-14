package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"DevMaan707/streamer/services"
	utils "DevMaan707/streamer/utils"
)

// videoListHandler returns a handler for listing videos
func videoListHandler(svc *services.VideoService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Check if refresh requested
		refresh := r.URL.Query().Get("refresh") == "true"

		// Get videos
		videos, err := svc.ListVideos(refresh)
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
