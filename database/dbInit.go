/*
 This file will manage the database initialization driver
*/

package database

import (
	"strings"
	"database/sql"
)

type DatabaseDriver string

const 
(
	MySql DatabaseDriver = "mysql"
	Postgres DatabaseDriver = "postgres"
	Oracle DatabaseDriver = "oracle"
	SqlServer DatabaseDriver = "sqlserver"
	MariaDb DatabaseDriver = "mariadb"
)



type databaseInitializer struct {
	Driver DatabaseDriver
	ConnectionString string

}

func InitializeDatabaseDriver(driverStr, connectionString string) databaseInitializer {
	driverStr = strings.ToLower(driverStr)
	var systemDriver DatabaseDriver
	switch driverStr {
		case "mysql":
			systemDriver = MySql
		case "mariadb":
			systemDriver = MariaDb
		case "oracle":
			systemDriver = Oracle
		
		case "postgres":
			systemDriver = Postgres
	
		case "sqlserver":
			systemDriver = SqlServer
			
		
	}
	return databaseInitializer{
		Driver: systemDriver,
		ConnectionString: connectionString,
	}
}

func (initializer *databaseInitializer) InitDatabase() *sql.DB{
	var db *sql.DB
	/*
		Esto no va a funcionar
		En primer lugar, hay que usar la funcion sql.OpenDB
		que devuelve un connector.Driver.
		Debemos crear un archivo por driver que vayamos a utilizar y utilizar
		un Facade para instanciar el driver correcto en función del driverString que se le pase
	*/

	db, err := sql.Open(string(initializer.Driver), initializer.ConnectionString )
	if err != nil {
		panic(err)
	}
	return db
}



