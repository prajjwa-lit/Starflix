package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"DevMaan707/streamer/services"
)

// uploadResponse defines the upload response format
type uploadResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	File    string `json:"file,omitempty"`
}

// uploadHandler returns a handler for file  uploads
func uploadHandler(svc *services.UploadService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Upload request received")

		// Only allow POST method
		if r.Method != http.MethodPost {
			log.Println("Method not allowed:", r.Method)
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Handle the file upload
		log.Println("Attempting to handle upload")
		filename, err := svc.HandleUpload(r)

		// Set content type for response
		w.Header().Set("Content-Type", "application/json")

		// Create response
		resp := uploadResponse{}

		if err != nil {
			log.Printf("Upload failed: %v", err)
			resp.Success = false
			resp.Message = fmt.Sprintf("Upload failed: %v", err)
			w.WriteHeader(http.StatusBadRequest)
		} else {
			log.Printf("Upload successful: %s", filename)
			resp.Success = true
			resp.Message = "Upload successful"
			resp.File = filename
			w.WriteHeader(http.StatusOK)
		}

		// Send response
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Printf("Failed to encode response: %v", err)
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}
	}
}
