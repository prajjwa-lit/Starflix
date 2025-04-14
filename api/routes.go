package api

import (
	"net/http"

	"DevMaan707/streamer/services"
)

// RegisterRoutes sets up all API routes
func RegisterRoutes(mux *http.ServeMux, videoSvc *services.VideoService, uploadSvc *services.UploadService) {
	// Video routes
	mux.HandleFunc("/api/videos", videoListHandler(videoSvc))
	mux.HandleFunc("/videos/", videoStreamHandler(videoSvc))

	// Upload routes
	mux.HandleFunc("/api/upload", uploadHandler(uploadSvc))
}
