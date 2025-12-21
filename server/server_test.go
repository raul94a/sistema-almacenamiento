package server

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/storage-system/database"
	"github.com/storage-system/server/controller"
	"github.com/storage-system/server/models"
	repository "github.com/storage-system/server/repositories"
	"github.com/storage-system/server/service"
	"gorm.io/gorm"
)

func Test_Landing_Page(t *testing.T) {
	t.Skip()
	uoc := controller.ObjectController{
		UploadService: &service.UploadObjectService{
			Repository: &repository.ObjectRepository{},
		},
	}
	serverHandler := HttpServerHandler{
		ObjectController: &uoc,
	}
	go serverHandler.HttpServer(nil)
	time.Sleep(5 * time.Second)
	if r, e := http.Get("http://localhost:4444/index.html"); e != nil {
		t.Fatal(e.Error())
	} else {
		var bytes []byte
		_, e = r.Body.Read(bytes)
		if e != nil {
			t.Fatal(e.Error())
		}
		if r.StatusCode != 200 {
			t.Fatalf("StatusCode %d for body %s", r.StatusCode, string(bytes))
		}
	}
}

func TestUploadFILE(t *testing.T) {
	port := "8888"
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() 
	go initHttpServer(port,ctx)
	time.Sleep(1 * time.Second)
	newFilename := "DATA.txt"
	f, _ := os.Create(newFilename)
	content := "TEST FILE"
	contentBytes := []byte(content)
	f.Write(contentBytes)
	f.Close()
	method := "POST"
	url := "http://localhost:8888/api/v1/PutObject"
	httpClient := &http.Client{
		Timeout: time.Second * 30,
	}
	
	currentFilepath,_ := os.Getwd()
	currentFilePathWithoutVolume := currentFilepath
	volume := ""
	if runtime.GOOS == "windows"{
		volume = filepath.VolumeName(currentFilePathWithoutVolume)	
		t.Log(volume)
		size := len(volume)
		currentFilePathWithoutVolume = currentFilePathWithoutVolume[size + 1:]
	}
	t.Log(currentFilePathWithoutVolume)

	const UPLOAD_FILENAME = "DATA_COPY.txt"
	UPLOAD_FILEPATH := filepath.Join(currentFilePathWithoutVolume,UPLOAD_FILENAME)
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	creator := &models.UploadObject{
		Description:      "Desc",
		Parent:           nil,
		EncryptionMethod: nil,
		Public:           true,
		Filename:         UPLOAD_FILEPATH,
		Extension:        ".txt",
		Hash:             "HASH",
		Size:             22,
		Unit:             "Bytes",
		UserOwner:        nil,
	}
	metaDataBytes, _ := json.Marshal(creator)
	_ = writer.WriteField("metadata", string(metaDataBytes))

	part, err := writer.CreateFormFile("file", creator.Filename)
	if err != nil {
		t.Fatal(err.Error())
	}
	_, err = part.Write(contentBytes)
	if err != nil {
		t.Fatal(err.Error())
	}
	writer.Close()

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		t.Fatal(err.Error())
	}

	if err != nil {
		t.Fatal(err.Error())
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	resp, err := httpClient.Do(req)
	if err != nil {
		t.Fatal(err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		t.Fatal(err.Error())
	}

	var result struct {
		ID uuid.UUID `json:"id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatal(err.Error())
	}

	t.Cleanup(func() {
		os.Remove(filepath.Join(currentFilepath,newFilename))
		os.Remove(filepath.Join(volume,string(filepath.Separator),UPLOAD_FILEPATH))
		time.Sleep(1 *time.Second)

		if err := os.Remove(filepath.Join(currentFilepath,"test.db")); err != nil {
			t.Log(err.Error())
		}
	})

}


func initHttpServer(port string,ctx context.Context) {
	serverHandler := createTestServer()
	
	// Start server in a separate goroutine
    go func() {
        serverHandler.HttpServer(&port)
    }()

    // Wait for the context to be cancelled
    <-ctx.Done()
    
    // Shut down gracefully (this releases the port and file handles)
    _, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	db , _ := serverHandler.ObjectController.UploadService.Repository.Db.DB()
	defer db.Close()
    defer cancel()
}

func createTestServer() *HttpServerHandler {
	db := createSQLiteTestDb()
	serverHandler := HttpServerHandler{
		ObjectController: &controller.ObjectController{
			UploadService: &service.UploadObjectService{
				Repository: &repository.ObjectRepository{
					Db: db,
				},
			},
		},
	}
	
	return &serverHandler
}

func createSQLiteTestDb() *gorm.DB {

	dbPath := "test.db"
	if err := os.Remove(dbPath); err != nil {
		fmt.Println(err.Error())

	}

	env := &database.DatabaseConfig{
		DatabaseType: "sqlite",
	}
	initializer := database.InitializeDatabaseDriver(env)
	initializer.DatabaseConfig = &database.DatabaseConfig{
		DatabaseName: dbPath,
	}
	db := initializer.InitDatabase()
	bytes, _ := os.ReadFile("../database/schema/schema.sql")

	if tx := db.Exec(string(bytes)); tx.Error != nil {

	}
	
	return db

}

