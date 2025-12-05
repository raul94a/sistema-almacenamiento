package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Configuración básica
var (
	UploadDir = "./uploads"
)

// Object simula tu modelo de datos
type Object struct {
	ID        string    `json:"id"`
	Filename  string    `json:"filename"`
	Bucket    string    `json:"bucket"`
	Size      int64     `json:"size"`
	CreatedAt time.Time `json:"created_at"`
}

func main() {
	// 1. Asegurar directorio de subidas
	if err := os.MkdirAll(UploadDir, os.ModePerm); err != nil {
		log.Fatalf("Error creando directorio de uploads: %v", err)
	}

	// 2. Definición de Rutas
	mux := http.NewServeMux()

	// --- RUTAS PÚBLICAS (Sin Auth) ---

	// Ruta: /objects (GET para listar, POST para subir)
	mux.HandleFunc("/objects", func(w http.ResponseWriter, r *http.Request) {
		// Log para ver qué llega
		log.Printf("[%s] %s", r.Method, r.URL.Path)
		
		switch r.Method {
		case http.MethodGet:
			handleListObjects(w, r)
		case http.MethodPost:
			PutFile(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Ruta: /objects/ (para operaciones por ID)
	mux.HandleFunc("/objects/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[%s] %s", r.Method, r.URL.Path)

		pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
		
		// Esperamos /objects/{id} o /objects/{id}/download
		if len(pathParts) < 2 {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}
		objectID := pathParts[1]

		// DELETE /objects/{id}
		if len(pathParts) == 2 && r.Method == http.MethodDelete {
			handleDeleteObject(w, r, objectID)
			return
		}
		
		// GET /objects/{id}/download
		if len(pathParts) == 3 && pathParts[2] == "download" && r.Method == http.MethodGet {
			handleDownloadObject(w, r, objectID)
			return
		}

		http.Error(w, "Not found or method not allowed", http.StatusNotFound)
	})

	// 3. Arrancar servidor directamente (sin envoltorios de seguridad)
	log.Println("🚀 Servidor 'Inseguro' escuchando en :4444")
	log.Fatal(http.ListenAndServe(":4444", mux))
}

// --- HANDLERS (Lógica de negocio) ---

func PutFile(w http.ResponseWriter, r *http.Request) {
	// 1. Parsear Multipart (Max 10MB de buffer en RAM)
	// Esto procesa el boundary y separa los archivos de los campos de texto
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		log.Printf("Error parsing multipart: %v", err)
		http.Error(w, "Error parsing multipart form", http.StatusBadRequest)
		return
	}

	// 2. Obtener el archivo (Clave "file" coincidente con el SDK)
	file, handler, err := r.FormFile("file")
	if err != nil {
		log.Printf("Error retrieving file key: %v", err)
		http.Error(w, "Error retrieving file (check 'file' key)", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// 3. Obtener metadatos (Clave "bucket" coincidente con el SDK)
	bucket := r.FormValue("bucket")
	
	// 4. Lógica de guardado en disco
	objID := uuid.New().String()
	// Prevenimos colisiones de nombre añadiendo el UUID
	filename := fmt.Sprintf("%s_%s", objID, handler.Filename)
	dstPath := filepath.Join(UploadDir, filename)
	
	dst, err := os.Create(dstPath)
	if err != nil {
		log.Printf("Error creating file on disk: %v", err)
		http.Error(w, "Error saving file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		log.Printf("Error copying bytes: %v", err)
		http.Error(w, "Error copying file", http.StatusInternalServerError)
		return
	}

	log.Printf("✅ Archivo recibido: %s (Bucket: %s) -> Guardado como: %s", handler.Filename, bucket, filename)

	// 5. Respuesta JSON (El cliente espera un JSON con "id")
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"id":      objID,
		"message": "File uploaded successfully",
		"bucket":  bucket,
	})
}

func handleListObjects(w http.ResponseWriter, r *http.Request) {
	// Dummy response
	objects := []Object{
		{ID: uuid.New().String(), Filename: "test.txt", Bucket: "default", Size: 123, CreatedAt: time.Now()},
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(objects)
}

func handleDownloadObject(w http.ResponseWriter, r *http.Request, id string) {
	// Buscar archivo que empiece por el ID
	files, _ := os.ReadDir(UploadDir)
	var foundPath string
	for _, f := range files {
		if strings.HasPrefix(f.Name(), id) {
			foundPath = filepath.Join(UploadDir, f.Name())
			break
		}
	}

	if foundPath == "" {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	log.Printf("Sirviendo archivo: %s", foundPath)
	http.ServeFile(w, r, foundPath)
}

func handleDeleteObject(w http.ResponseWriter, r *http.Request, id string) {
	// Buscar y borrar
	files, _ := os.ReadDir(UploadDir)
	for _, f := range files {
		if strings.HasPrefix(f.Name(), id) {
			err := os.Remove(filepath.Join(UploadDir, f.Name()))
			if err != nil {
				http.Error(w, "Error deleting file", http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}
	http.Error(w, "File not found", http.StatusNotFound)
}