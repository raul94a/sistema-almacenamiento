package repository

import (
	"os"
	"testing"

	"github.com/storage-system/database"
	"github.com/storage-system/server/models"
	"gorm.io/gorm"
)



func useSqlite(t *testing.T) *gorm.DB{

	dbPath := "../test.db"
	os.Remove(dbPath)

	env := &database.DatabaseConfig{
		DatabaseType: "sqlite",
	}
	initializer := database.InitializeDatabaseDriver(env)
	initializer.DatabaseConfig = &database.DatabaseConfig{
		DatabaseName: dbPath,
	}
	db := initializer.InitDatabase()
	createObjectsTable(db,t)
	if db == nil {
		t.Errorf("Error initializing database")
	}
	sqliteDb, _ := db.DB()

	if err := sqliteDb.Ping(); err != nil {
		t.Errorf("TestSqlite3_Init error: %v", err)
	}
	return db

	
}

func createObjectsTable(db *gorm.DB, t *testing.T ){
	bytes ,err := os.ReadFile("../../database/schema/schema.sql")
	if err != nil {
		t.Fatal(err.Error())
		} 
	t.Logf("SCHEMA leido correctamente\n%s ",string(bytes))
	if tx := db.Exec(string(bytes));tx.Error != nil {
		t.Fatal(tx.Error.Error())
	}
	t.Log("SCHEMA se ha ejecutado. Accediendo")
	var count int64
	db.Exec("SELECT COUNT(*) FROM objects",&count)
	t.Logf("NUMERO DE OBJECTOS ENCONTRADOS %v",count)
	

}

func createRepository(t *testing.T)*ObjectRepository{
	db := useSqlite(t)
	return &ObjectRepository{
		Db: db,
	}
}

func TestUploadFile(t *testing.T){
	repository := createRepository(t)
	object := &models.Object{
		Disk: "C",
		Location: "pwd/Desktop",
		Filename: "testfile.html",
		Description: "mydec",
		Hash: "hash",
		Extension: ".html",

	}
	err := repository.UploadFile(object)
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Log("SCHEMA se ha ejecutado. Accediendo")
	var count int64
	repository.Db.Model(&models.Object{}).Count(&count)
	t.Logf("NUMERO DE OBJETOS ENCONTRADOS %v",count)
	const NUMBER_OBJECTS_MUST_BE_FOUND = 1
	if count != NUMBER_OBJECTS_MUST_BE_FOUND {
		t.Fatalf("Error. Objects found (%d) is different from %d",count,NUMBER_OBJECTS_MUST_BE_FOUND)
	}

	os.Remove("../test.db")

}