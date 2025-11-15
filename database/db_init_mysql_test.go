package database

import (
	"testing"
)

var testConnStr string = "test@connectionstring:9999/test-db"
var testDriver string = "TEST"
var mysqlDriver string = "mysql"
var sqliteDriver string = "sqlite3"

func TestNotImplementedInitialization(t *testing.T) {
	env := &DatabaseConfig{
		DatabaseType: testDriver,
	}
	initializer := InitializeDatabaseDriver(env)

	loader := initializer.Loader
	if l, ok := loader.(NotImplementedLoader); ok {
		t.Log("TestNotImplementedInitialization PASSED")
	} else {
		t.Errorf("TestNotImplementedInitialization failed. loader type is %T", l)
	}
}

func TestMySqlInitialization(t *testing.T) {
	env := &DatabaseConfig{
		DatabaseType:     mysqlDriver,
		DatabaseUrl:      "19.9.9.22",
		DatabasePort:     "9888",
		DatabaseUser:     "juser",
		DatabasePassword: "pwd",
		DatabaseName:     "db",
	}
	initializer := InitializeDatabaseDriver(env)

	loader := initializer.Loader
	if l, ok := loader.(MySqlLoader); ok {
		t.Log("TestMySqlInitialization PASSED")
	} else {
		t.Errorf("TestMySqlInitialization failed. loader type is %T", l)
	}

	t.Log(initializer.ConnectionString)

}

func TestMySql_Env_Connection_String(t *testing.T) {
	env := LoadDatabaseConfigFromEnv("./.mysql")
	initializer := InitializeDatabaseDriver(env)
	loader := initializer.Loader

	err, dsn := loader.BuildDsnFromEnv("./.mysql")

	if err != nil {
		t.Error(err.Error())
	} else {
		t.Logf("Correct connection String %s", dsn)
	}

}

func Test_Sqlite3_Initialization(t *testing.T) {
	env := &DatabaseConfig{
		DatabaseType: sqliteDriver,
	}
	initializer := InitializeDatabaseDriver(env)

	loader := initializer.Loader
	if l, ok := loader.(SqliteLoader); ok {
		t.Log("TestMySqlInitialization PASSED")
	} else {
		t.Errorf("TestMySqlInitialization failed. loader type is %T", l)
	}
}
