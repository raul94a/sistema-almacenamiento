package database

import
(
	"testing"
)

var testConnStr string = "test@connectionstring:9999/test-db"
var testDriver string = "TEST"
var mysqlDriver string = "mysql"
var sqliteDriver string = "sqlite3"
func TestNotImplementedInitialization(t *testing.T) {
	initializer := InitializeDatabaseDriver(testDriver,testConnStr)

	loader := initializer.Loader
	if l, ok := loader.(NotImplementedLoader); ok {
		t.Log("TestNotImplementedInitialization PASSED")
	} else {
    	t.Errorf("TestNotImplementedInitialization failed. loader type is %T", l) 
	}
}

func TestMySqlInitialization(t *testing.T) {
	initializer := InitializeDatabaseDriver(mysqlDriver,testConnStr)

	loader := initializer.Loader
	if l, ok := loader.(MySqlLoader); ok {
		t.Log("TestMySqlInitialization PASSED")
	} else {
    	t.Errorf("TestMySqlInitialization failed. loader type is %T", l) 
	}
}


func Test_Sqlite3_Initialization(t *testing.T) {
	initializer := InitializeDatabaseDriver(sqliteDriver,testConnStr)

	loader := initializer.Loader
	if l, ok := loader.(SqliteLoader); ok {
		t.Log("TestMySqlInitialization PASSED")
	} else {
    	t.Errorf("TestMySqlInitialization failed. loader type is %T", l) 
	}
}