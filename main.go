package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"DevMaan707/streamer/config"
	"DevMaan707/streamer/db"
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
	cfg := config.NewConfig()
	flag.IntVar(&cfg.Port, "port", 5101, "Port to serve on")
	flag.StringVar(&cfg.VideoDir, "videos", "./videos", "Directory containing video files")
	flag.StringVar(&cfg.CoverImageDir, "covers", "./covers", "Directory for video cover images")
	flag.IntVar(&cfg.MaxUploadSize, "max-upload", 1024, "Maximum upload size in MB")
	flag.Parse()
	cfg.VideoDir = expandPath(cfg.VideoDir)
	cfg.CoverImageDir = expandPath(cfg.CoverImageDir)
	absPath, err := filepath.Abs(cfg.VideoDir)
	if err != nil {
		log.Fatalf("Error resolving video directory path: %v", err)
	}

	cfg.VideoDir = absPath
	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		log.Printf("Video directory does not exist, creating: %s", absPath)
		if err := os.MkdirAll(absPath, 0755); err != nil {
			log.Fatalf("Failed to create video directory: %v", err)
		}
	} else if err != nil {
		log.Fatalf("Error accessing video directory: %v", err)
	}
	absPathCovers, err := filepath.Abs(cfg.CoverImageDir)
	if err != nil {
		log.Fatalf("Error resolving cover images directory path: %v", err)
	}

	cfg.CoverImageDir = absPathCovers
	if _, err := os.Stat(absPathCovers); os.IsNotExist(err) {
		log.Printf("Cover images directory does not exist, creating: %s", absPathCovers)
		if err := os.MkdirAll(absPathCovers, 0755); err != nil {
			log.Fatalf("Failed to create cover images directory: %v", err)
		}
	}
	if err := db.Initialize(); err != nil {
		log.Fatalf("Database initialization failed: %v", err)
	}
	if err := db.TestConnection(); err != nil {
		log.Printf("WARNING: Database connection test failed: %v", err)
		log.Println("The application will continue, but some features may not work properly")
	}
	if err := db.EnsureTablesExist(); err != nil {
		log.Printf("WARNING: Failed to verify database tables: %v", err)
	}

	srv, err := server.NewServer(cfg)
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	fmt.Printf("Starting StreamFlix video platform on http://localhost:%d\n", cfg.Port)
	fmt.Printf("Serving videos from: %s\n", absPath)
	fmt.Printf("Storing cover images in: %s\n", absPathCovers)
	fmt.Printf("Maximum upload size: %d MB\n", cfg.MaxUploadSize)

	if err := srv.Start(); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
