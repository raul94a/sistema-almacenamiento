package database

import
(
	"testing"
	"os"
)


func TestSqlite3_Init(t *testing.T) {
	dbPath := "../test.db"
	os.Remove(dbPath)

	initializer := InitializeDatabaseDriver("sqlite3",dbPath)
	db := initializer.InitDatabase()

	if db == nil {
		t.Errorf("Error initializing database")
	}

	if err := db.Ping(); err != nil {
		t.Errorf("TestSqlite3_Init error: %v",err)
	}

	t.Cleanup(func() {
		os.Remove(dbPath)
	})
	
}