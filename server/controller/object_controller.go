package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/storage-system/server/models"
	"github.com/storage-system/server/service"
)

type ObjectController struct {
	UploadService *service.UploadObjectService
}

func (u *ObjectController) UploadObject(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// 1. MUST BE FIRST: Set the max bytes reader on the raw body
	maxBytes := int64(33 * 1024 * 1024)
	r.Body = http.MaxBytesReader(w, r.Body, maxBytes)

	// 2. Parse the multipart form
	// This reads the body and separates the "metadata" and "file" parts
	if err := r.ParseMultipartForm(maxBytes); err != nil {
		fmt.Printf("Parse Error: %v\n", err)
		http.Error(w, "Error parsing form or file too large", http.StatusBadRequest)
		return
	}

	// 3. Extract Metadata
	metaData := r.FormValue("metadata")
	var uploadObject models.UploadObject
	if err := json.Unmarshal([]byte(metaData), &uploadObject); err != nil {
		fmt.Printf("JSON Error: %v | Data: %s\n", err, metaData)
		http.Error(w, "Invalid JSON metadata", http.StatusBadRequest)
		return
	}

	// 4. Extract File
	file, header, err := r.FormFile("file")
	if err != nil {
		fmt.Printf("File Error: %v\n", err)
		http.Error(w, "Could not get file part", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Debugging info (Now it will have values)
	fmt.Printf("Metadata received: %+v\n", uploadObject)
	fmt.Printf("Received file: %s, Size: %d bytes\n", header.Filename, header.Size)

	// 5. Upload via Service
	id, err := u.UploadService.Upload(uploadObject, file, nil)
	if err != nil {
		fmt.Printf("Service Error: %v\n", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// 6. Respond with JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"id": id})
}
