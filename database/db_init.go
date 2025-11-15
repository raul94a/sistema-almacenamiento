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
	Sqlite3  DatabaseDriver = "sqlite3"
)

// Will handle the connection to the database
// see databaseInitializer.InitDatabase for more info
type databaseInitializer struct {
	Driver           DatabaseDriver
	Loader           DatabaseLoader
	ConnectionString string
	Env				 *Env
}


// TODO: Delete driverStr from the method signature.
//		 This field will be read from .env file
func InitializeDatabaseDriver(env *Env) databaseInitializer {
	driverStr := strings.ToLower(env.DatabaseType)
	var systemDriver DatabaseDriver
	var loader DatabaseLoader = NotImplementedLoader{}
	var connectionString string
	var err error
	switch driverStr {
	case "mysql":
		systemDriver = MySql
		loader = MySqlLoader{}
		err, connectionString = loader.BuildDsn(env)
	case "mariadb":
		systemDriver = MariaDb
		loader = MySqlLoader{}
		err, connectionString = loader.BuildDsn(env)
	case "oracle":
		systemDriver = Oracle
		err, connectionString = loader.BuildDsn(env)
	case "postgres":
		systemDriver = Postgres
		loader = PostgresLoader{}
		err, connectionString = loader.BuildDsn(env)
	case "sqlserver":
		systemDriver = SqlServer
	case "sqlite3":
		systemDriver = Sqlite3
		loader = SqliteLoader{}
		err, connectionString = loader.BuildDsn(env)
	}

	if err != nil {
		panic(err)
	}
	
	
	return databaseInitializer{
		Driver:           systemDriver,
		ConnectionString: connectionString,
		Loader:           loader,
		Env: 			  env,
	}
}

// This method from databaseInitializer will gracefully handle the driver
// that will connect the app to the database.
func (initializer *databaseInitializer) InitDatabase() *gorm.DB {
	var db *gorm.DB
	// Read env variables: driver + connection
	
	db, err := initializer.Loader.LoadDatabase(initializer.ConnectionString)
	if err != nil {
		panic(err)
	}
	return db
}

func GetDatabase() (*gorm.DB) {
	godotenv.Load(".env")
	env := LoadDatabaseEnvVariables(".env")
	initializer := InitializeDatabaseDriver(env)
	return initializer.InitDatabase()
}
/*
	Reglas:
	1. Utilizar solo un método que cargue la base de datos configurada con el entorno -InitDatabase-

	Pasos:
	1. Carga del .env
	2. Connection to database -LoaderDatabase
	
*/
