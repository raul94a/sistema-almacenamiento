package database

import (
	"context"
	"testing"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

var testConnStr string = "test@connectionstring:9999/test-db"
var testDriver string = "TEST"
var mysqlDriver string = "mysql"
var sqliteDriver string = "sqlite"

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

}

func TestMySql_Env_Connection_String(t *testing.T) {

	env := LoadDatabaseConfigFromEnv("./test/.mysql")
	initializer := InitializeDatabaseDriver(env)
	loader := initializer.Loader

	err, dsn := loader.BuildDsnFromEnv("./test/.mysql")

	if err != nil {
		t.Error(err.Error())
	} else {
		t.Logf("Correct connection String %s", dsn)
	}

}

func Test_Sqlite3_Initialization(t *testing.T) {
	env := &DatabaseConfig{
		DatabaseType: sqliteDriver,
		DatabaseName: "db",
	}
	initializer := InitializeDatabaseDriver(env)
	loader := initializer.Loader
	if l, ok := loader.(SqliteLoader); ok {
		t.Log("TestMySqlInitialization PASSED")
	} else {
		t.Errorf("TestMySqlInitialization failed. loader type is %T", l)
	}
}

func Test_Integration_MySQL(t *testing.T) {
	ctx := context.Background()

    req := testcontainers.ContainerRequest{
        Image:        "mysql:8.0.36",
        ExposedPorts: []string{"3333/tcp"},
        Env: map[string]string{
            "MYSQL_ROOT_PASSWORD": "root",
            "MYSQL_DATABASE":      "foo",
        },
        Cmd: []string{"mysqld", "--port=3333"}, // AQUÍ ESTÁ LA CLAVE
        WaitingFor: wait.ForListeningPort("3333/tcp"),
        HostConfigModifier: func(hc *container.HostConfig) {
            hc.PortBindings = nat.PortMap{
                "3333/tcp": []nat.PortBinding{{HostIP: "127.0.0.1", HostPort: "3333"}},
            }
        },
    }

    // 2. Inicia el contenedor genérico
    mysqlC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
        ContainerRequest: req,
        Started:          true,
    })
    if err != nil {
        t.Fatalf("failed to start container: %v", err)
    }


    // ---------- 2. Auto-remove (stop + rm + rm volumes) ----------
    t.Cleanup(func() {
        // TerminateRemoveVolumes stops + removes container + removes any volumes
        if err := testcontainers.TerminateContainer(mysqlC); err != nil {
            t.Logf("cleanup error: %v", err)
        }
    })

    

    // ---------- 5. Build your config that points at the container ----------
    cfg := &DatabaseConfig{
        DatabaseType:     "mysql",
        DatabaseUrl:      "localhost", // host visible from the test process
        DatabasePort:     "3333",      // mapped host port
        DatabaseUser:     "root",
        DatabasePassword: "root",
        DatabaseName:     "foo",
    }

	

    // ---------- 6. Call your driver ----------
    init := InitializeDatabaseDriver(cfg)
    _, err = init.load()
    if err != nil {
        t.Fatalf("driver load failed: %v", err)
    }

    // ---------- 7. (optional) sanity-check that the DB really exists ----------
    
    t.Logf("Database `%s` is ready – driver loaded successfully", "foo")
}
