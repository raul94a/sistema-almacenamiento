package database

import (
	"os"
	"github.com/joho/godotenv"
)

type DatabaseConfig struct {
	DatabaseType     string
	DatabaseUrl      string
	DatabaseUser     string
	DatabasePassword string
	DatabasePort     string
	DatabaseName     string
}

func LoadDatabaseEnvVariables(path string) *DatabaseConfig {
	godotenv.Load(path)
	dbType := os.Getenv("DATABASE_TYPE")
	url := os.Getenv("DATABASE_URL")
	port := os.Getenv("DATABASE_PORT")
	user := os.Getenv("DATABASE_USER")
	pwd := os.Getenv("DATABASE_PASSWORD")
	dbName := os.Getenv("DATABASE_NAME")
	return &DatabaseConfig{
		DatabaseType:     dbType,
		DatabaseUrl:      url,
		DatabaseUser:     user,
		DatabasePassword: pwd,
		DatabasePort:     port,
		DatabaseName:     dbName,
	}

}
