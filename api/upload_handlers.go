package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"DevMaan707/streamer/services"
)

type uploadResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	File    string `json:"file,omitempty"`
}

func uploadHandler(svc *services.UploadService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Upload request received")

		if r.Method != http.MethodPost {
			log.Println("Method not allowed:", r.Method)
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		log.Println("Attempting to handle upload")
		filename, err := svc.HandleUpload(r)
		w.Header().Set("Content-Type", "application/json")
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
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Printf("Failed to encode response: %v", err)
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}
	}
}
