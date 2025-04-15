package server

import (
	"embed"
	"fmt"
	"io/fs"
	"net/http"
	"time"

	"DevMaan707/streamer/api"
	"DevMaan707/streamer/config"
	"DevMaan707/streamer/services"
)

//go:embed static/*
var staticFiles embed.FS

type Server struct {
	cfg       *config.Config
	videoSvc  *services.VideoService
	uploadSvc *services.UploadService
}

func NewServer(cfg *config.Config) (*Server, error) {
	videoSvc, err := services.NewVideoService(cfg.VideoDir, cfg.CoverImageDir)
	if err != nil {
		return nil, fmt.Errorf("failed to create video service: %w", err)
	}

	uploadSvc := services.NewUploadService(cfg.VideoDir, cfg.CoverImageDir, cfg.MaxUploadSizeBytes())

	return &Server{
		cfg:       cfg,
		videoSvc:  videoSvc,
		uploadSvc: uploadSvc,
	}, nil
}

func (s *Server) Start() error {
	mux := http.NewServeMux()
	api.RegisterRoutes(mux, s.videoSvc, s.uploadSvc)
	staticFS, err := fs.Sub(staticFiles, "static")
	if err != nil {
		return fmt.Errorf("failed to load static files: %w", err)
	}
	mux.Handle("/", http.FileServer(http.FS(staticFS)))
	handler := CloudflareMiddleware(mux)
	handler = ErrorLoggingMiddleware(handler)
	handler = LoggingMiddleware(handler)
	handler = CORSMiddleware(handler)
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", s.cfg.Port),
		Handler:      handler,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 30 * time.Minute,
		IdleTimeout:  120 * time.Second,
	}

	return server.ListenAndServe()
}
