package main

import (
	// "database/sql"
	// "github.com/gin-gonic/gin"
	// "go.uber.org/zap"
	// "github.com/joho/godotenv"
	"fmt"
	"os"

	"github.com/storage-system/database"
	"github.com/storage-system/server"
	"github.com/storage-system/server/controller"
	repository "github.com/storage-system/server/repositories"
	"github.com/storage-system/server/service"
	"gorm.io/gorm"
)

func main()  {
	fmt.Println("🇪🇸 🇪🇸 🇪🇸 🇪🇸  Hello World! 🇪🇸 🇪🇸 🇪🇸 🇪🇸")
	db := database.GetDatabase()
	// testing
	createObjectsTable(db)
	server.HttpServer(CreateUploadObjectController(db))
}

func CreateUploadObjectController(db *gorm.DB) controller.UploadObjectController {
	
	uoc := controller.UploadObjectController {
		UploadService: &service.UploadObjectService{
			Repository: &repository.ObjectRepository{Db: db},
		},
	}
	return uoc;
}

//TODO: Solo develop
func createObjectsTable(db *gorm.DB) {
	bytes, _ := os.ReadFile("database/schema/schema.sql")
	
	if tx := db.Exec(string(bytes)); tx.Error != nil {
		
	}
	
	

}