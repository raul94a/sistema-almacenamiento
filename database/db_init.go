// Copyright 2025 Raul Albín Alba
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package database

import (
	"strings"

	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

/* The types of db drivers the app can handle */
type DatabaseDriver string

const (
	MySql     DatabaseDriver = "mysql"
	Postgres  DatabaseDriver = "postgres"
	Oracle    DatabaseDriver = "oracle"
	SqlServer DatabaseDriver = "sqlserver"
	MariaDb   DatabaseDriver = "mariadb"
	Sqlite3   DatabaseDriver = "sqlite"
)

// Will handle the connection to the database
// see databaseInitializer.InitDatabase for more info
type databaseInitializer struct {
	Driver           DatabaseDriver
	Loader           DatabaseLoader
	DatabaseConfig              *DatabaseConfig
}


func (d *databaseInitializer) load() (*gorm.DB,error){
	cnf := d.DatabaseConfig
	return d.Loader.LoadDatabase(cnf)
}

// Depending on the database config, usually read from a .env file,
// the correct DatabaseLoader will be allocated.
func InitializeDatabaseDriver(databaseConfig *DatabaseConfig) databaseInitializer {
	driverStr := strings.ToLower(databaseConfig.DatabaseType)
	var systemDriver DatabaseDriver
	var loader DatabaseLoader = NotImplementedLoader{}
	switch driverStr {
	case "mysql":
		systemDriver = MySql
		loader = MySqlLoader{}
		
	case "mariadb":
		systemDriver = MariaDb
		loader = MySqlLoader{}
		
	case "oracle":
		systemDriver = Oracle
		
	case "postgres":
		systemDriver = Postgres
		loader = PostgresLoader{}
		
	case "sqlserver":
		systemDriver = SqlServer
	case "sqlite": 
		systemDriver = Sqlite3
		loader = SqliteLoader{}
		
	}

	return databaseInitializer{
		Driver:           systemDriver,
		Loader:           loader,
		DatabaseConfig:              databaseConfig,
	}
}

// This method from databaseInitializer will gracefully handle the driver
// that will connect the app to the database. 
// It needs a DatabaseConfig object to work
func (initializer *databaseInitializer) InitDatabase() *gorm.DB {
	var db *gorm.DB
	db, err := initializer.load()
	if err != nil {
		panic(err)
	}
	return db
}

// Entrypoint for creating a Database connection
// It uses the required .env file located in the root directory.
// 
// The magic of loading different kind of databases resides in the
// InitializeDatabaseDriver() method, which loads the correct driver
//
func GetDatabase() *gorm.DB {
	godotenv.Load(".env")
	env := LoadDatabaseConfigFromEnv(".env")
	initializer := InitializeDatabaseDriver(env)
	return initializer.InitDatabase()
}
