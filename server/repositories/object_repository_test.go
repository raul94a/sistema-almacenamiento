package repository

import (
	"os"
	"testing"

	"github.com/storage-system/database"
	"github.com/storage-system/server/models"
	"gorm.io/gorm"
)

func useSqlite(t *testing.T) *gorm.DB {

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
	createObjectsTable(db, t)
	if db == nil {
		t.Errorf("Error initializing database")
	}
	sqliteDb, _ := db.DB()

	if err := sqliteDb.Ping(); err != nil {
		t.Errorf("TestSqlite3_Init error: %v", err)
	}
	return db

}

func createObjectsTable(db *gorm.DB, t *testing.T) {
	bytes, err := os.ReadFile("../../database/schema/schema.sql")
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Logf("SCHEMA leido correctamente\n%s ", string(bytes))
	if tx := db.Exec(string(bytes)); tx.Error != nil {
		t.Fatal(tx.Error.Error())
	}
	t.Log("SCHEMA se ha ejecutado. Accediendo")
	var count int64
	db.Exec("SELECT COUNT(*) FROM objects", &count)
	t.Logf("NUMERO DE OBJECTOS ENCONTRADOS %v", count)

}

func createRepository(t *testing.T) *ObjectRepository {
	db := useSqlite(t)
	return &ObjectRepository{
		Db: db,
	}
}

func tearDown() {
	os.Remove("../test.db")
}

func countObjectsHook(t *testing.T,repository ObjectRepository) int64{
	var count int64
	repository.Db.Model(&models.Object{}).Count(&count)
	t.Logf("NUMERO DE OBJETOS ENCONTRADOS %v", count)
	return count
}

func uploadObjectHook(t *testing.T, repository ObjectRepository, id *string) *models.Object {
	var ownerid string
	ownerid = "testuser"
	object := &models.Object{
		ID:          "___TEST_ID___",
		Disk:        "",
		Location:    "../../",
		Filename:    ".env-copy",
		Description: "mydec",
		Hash:        "hash",
		Extension:   "",
		UserOwner:   &ownerid,
	}
	if id != nil {
		object.ID = *id;
	}
	err := repository.UploadFile(object)
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Log("SCHEMA se ha ejecutado. Accediendo")
	
	return object

}

func TestUploadFile(t *testing.T) {
	defer tearDown()
	repository := createRepository(t)
	objectId := "myfile"
	uploadObjectHook(t,*repository,&objectId)
	count := countObjectsHook(t,*repository)
	const NUMBER_OBJECTS_MUST_BE_FOUND = 1
	if count != NUMBER_OBJECTS_MUST_BE_FOUND {
		t.Fatalf("Error. Objects found (%d) is different from %d", count, NUMBER_OBJECTS_MUST_BE_FOUND)
	}

}

// The object is logically deleted
func TestSoftDeleteObject(t *testing.T) {
	defer tearDown()

	repository := createRepository(t)
	objectId := "myfile2"

	object := uploadObjectHook(t,*repository,&objectId)
	if err := repository.DeleteFile(object); err != nil{
		t.Fatal(err.Error())
	}
	
	count := countObjectsHook(t,*repository)
	const ZERO = 0
	if count != ZERO{
		t.Fatalf("Error deleting object")
	}
}

func TestFetchFiles(t *testing.T){
	tearDown()
	defer tearDown()

	repository := createRepository(t)
	var id1 string
	id1 = "TEST_ID_1___2" 
 	_ = uploadObjectHook(t,*repository, &id1)

	id2 := "TEST_ID_2"
	_ = uploadObjectHook(t,*repository, &id2)

	
	objects,err := repository.GetUserFiles("testuser",1,10)
	if err != nil {
		t.Fatal(err.Error())
	}
	const TWO = 2

	if len(objects.Items) != TWO {
		t.Fatalf("Number of objects retrieved (%d) failed. %d where expected",len(objects.Items), TWO)
	}
}

func TestGetFile(t *testing.T){
	tearDown()
	defer tearDown()

	repository := createRepository(t)
	var id1 string
	id1 = "TEST_ID_1___2" 
 	_ = uploadObjectHook(t,*repository, &id1)

	id2 := "TEST_ID_2"
	_ = uploadObjectHook(t,*repository, &id2)

	
	objects,err := repository.GetUserFiles("testuser",1,10)
	if err != nil {
		t.Fatal(err.Error())
	}
	const TWO = 2

	if len(objects.Items) != TWO {
		t.Fatalf("Number of objects retrieved (%d) failed. %d where expected",len(objects.Items), TWO)
	}

	for _, object := range objects.Items {
		if obj,err := repository.GetUserFile(object.ID);err != nil{
			t.Fatal(err.Error())
		} else {
			t.Logf("Object %v",string(obj))
		}
	}
}


