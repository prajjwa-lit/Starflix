package api

import (
	"net/http"

	"DevMaan707/streamer/services"
)

func RegisterRoutes(mux *http.ServeMux, videoSvc *services.VideoService, uploadSvc *services.UploadService) {

	mux.HandleFunc("/api/videos", videoListHandler(videoSvc))
	mux.HandleFunc("/api/videos/genre/", videoListByGenreHandler(videoSvc))
	mux.HandleFunc("/api/genres", genreListHandler(videoSvc))
	mux.HandleFunc("/videos/", videoStreamHandler(videoSvc))
	mux.HandleFunc("/covers/", coverImageHandler(videoSvc))

	mux.HandleFunc("/api/upload", uploadHandler(uploadSvc))
}
