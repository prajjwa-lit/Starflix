package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"DevMaan707/streamer/services"
)

// uploadResponse defines the upload response format
type uploadResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	File    string `json:"file,omitempty"`
}

// uploadHandler returns a handler for file uploads
func uploadHandler(svc *services.UploadService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Only allow POST method
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Handle the file upload
		filename, err := svc.HandleUpload(r)

		// Set content type for response
		w.Header().Set("Content-Type", "application/json")

		// Create response
		resp := uploadResponse{}

		if err != nil {
			resp.Success = false
			resp.Message = fmt.Sprintf("Upload failed: %v", err)
			w.WriteHeader(http.StatusBadRequest)
		} else {
			resp.Success = true
			resp.Message = "Upload successful"
			resp.File = filename
			w.WriteHeader(http.StatusOK)
		}

		// Send response
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}
	}
}
