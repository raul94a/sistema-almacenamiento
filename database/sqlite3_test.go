package database

import (
	"os"
	"testing"
)

func TestSqlite3_Init(t *testing.T) {
	dbPath := "../test.db"
	os.Remove(dbPath)

	env := &DatabaseConfig{
		DatabaseType: "sqlite3",
	}
	initializer := InitializeDatabaseDriver(env)
	db := initializer.InitDatabase()

	if db == nil {
		t.Errorf("Error initializing database")
	}
	sqliteDb, _ := db.DB()

	if err := sqliteDb.Ping(); err != nil {
		t.Errorf("TestSqlite3_Init error: %v", err)
	}

	t.Cleanup(func() {
		os.Remove(dbPath)
	})

}
