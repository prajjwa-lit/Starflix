package api

import (
	"net/http"

	"DevMaan707/streamer/services"
)

// RegisterRoutes sets up all API routes
func RegisterRoutes(mux *http.ServeMux, videoSvc *services.VideoService, uploadSvc *services.UploadService) {
	// Video routes
	mux.HandleFunc("/api/videos", videoListHandler(videoSvc))
	mux.HandleFunc("/api/videos/genre/", videoListByGenreHandler(videoSvc))
	mux.HandleFunc("/api/genres", genreListHandler(videoSvc))
	mux.HandleFunc("/videos/", videoStreamHandler(videoSvc))
	mux.HandleFunc("/covers/", coverImageHandler(videoSvc))

	// Upload routes
	mux.HandleFunc("/api/upload", uploadHandler(uploadSvc))
}
