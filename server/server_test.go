package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"os"
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
	go initHttpServer(port)
	time.Sleep(1 * time.Second)
	f, _ := os.Create("DATA.txt")
	content := "TEST FILE"
	contentBytes := []byte(content)
	f.Write(contentBytes)
	f.Close()
	method := "POST"
	url := "http://localhost:8888/api/v1/PutObject"
	httpClient := &http.Client{
		Timeout: time.Second * 30,
	}
	const UPLOAD_FILEPATH = "Users/Raul9/Documents/Desarrollo/Go/sistema-almacenamiento/files/DATA.txt"
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

	os.Remove("C:/Users/Raul9/Documents/Desarrollo/Go/sistema-almacenamiento/server/DATA.txt")
	os.Remove(fmt.Sprintf("C:/%s",UPLOAD_FILEPATH))
}

func initHttpServer(port string) {
	serverHandler := createTestServer()
	serverHandler.HttpServer(&port)
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
	os.Remove(dbPath)

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
