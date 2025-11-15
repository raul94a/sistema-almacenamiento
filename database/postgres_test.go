package database

import
(
	"testing"
)

func TestPostgresInitialization(t *testing.T) {
	env := &Env{
		DatabaseType: "postgres",
	}
	initializer := InitializeDatabaseDriver(env)
	

	loader := initializer.Loader
	if l, ok := loader.(PostgresLoader); ok {
		t.Log("TestMySqlInitialization PASSED")
	} else {
    	t.Errorf("TestMySqlInitialization failed. loader type is %T", l) 
	}
}

