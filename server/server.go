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

// Server represents the HTTP server
type Server struct {
	cfg       *config.Config
	videoSvc  *services.VideoService
	uploadSvc *services.UploadService
}

// NewServer creates a new server instance
func NewServer(cfg *config.Config) (*Server, error) {
	// Create services
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

// Start begins the HTTP server
func (s *Server) Start() error {
	// Set up router with middlewares
	mux := http.NewServeMux()

	// Register API routes
	api.RegisterRoutes(mux, s.videoSvc, s.uploadSvc)

	// Set up static files
	staticFS, err := fs.Sub(staticFiles, "static")
	if err != nil {
		return fmt.Errorf("failed to load static files: %w", err)
	}

	// Serve static files
	mux.Handle("/", http.FileServer(http.FS(staticFS)))

	// Wrap with logging middleware
	handler := LoggingMiddleware(mux)

	// Configure server with reasonable timeouts
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", s.cfg.Port),
		Handler:      handler,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 30 * time.Minute, // Long timeout for video uploads/streaming
		IdleTimeout:  120 * time.Second,
	}

	return server.ListenAndServe()
}
