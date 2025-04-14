package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"DevMaan707/streamer/config"
	"DevMaan707/streamer/server"
)

func expandPath(path string) string {
	if len(path) == 0 || path[0] != '~' {
		return path
	}

	home, err := os.UserHomeDir()
	if err != nil {
		log.Printf("Warning: could not expand home directory: %v", err)
		return path
	}

	return filepath.Join(home, path[1:])
}
func main() {
	// Parse command line flags
	cfg := config.NewConfig()
	flag.IntVar(&cfg.Port, "port", 8080, "Port to serve on")
	flag.StringVar(&cfg.VideoDir, "videos", "./videos", "Directory containing video files")
	flag.IntVar(&cfg.MaxUploadSize, "max-upload", 1024, "Maximum upload size in MB")
	flag.Parse()

	// Ensure video directory exists, create if not
	cfg.VideoDir = expandPath(cfg.VideoDir)

	// Ensure video directory exists, create if not
	absPath, err := filepath.Abs(cfg.VideoDir)
	if err != nil {
		log.Fatalf("Error resolving video directory path: %v", err)
	}

	cfg.VideoDir = absPath // Store the absolute path

	// Check if directory exists, create if not
	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		log.Printf("Video directory does not exist, creating: %s", absPath)
		if err := os.MkdirAll(absPath, 0755); err != nil {
			log.Fatalf("Failed to create video directory: %v", err)
		}
	} else if err != nil {
		log.Fatalf("Error accessing video directory: %v", err)
	}
	// Create and start server
	srv, err := server.NewServer(cfg)
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	fmt.Printf("Starting video streaming server on http://localhost:%d\n", cfg.Port)
	fmt.Printf("Serving videos from: %s\n", absPath)
	fmt.Printf("Maximum upload size: %d MB\n", cfg.MaxUploadSize)

	if err := srv.Start(); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
