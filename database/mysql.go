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
	"fmt"
	"github.com/storage-system/database/patterns"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"regexp"
)

type MySqlLoader struct {
}

func (m MySqlLoader) LoadDatabase(connectionString string) (*gorm.DB, error) {
	return gorm.Open(mysql.Open(connectionString), &gorm.Config{})
}

func (m MySqlLoader) buildDsnFromEnvStruct(env *Env) string {
	return env.DatabaseUser + ":" + env.DatabasePassword + "@" + "tcp(" + env.DatabaseUrl + ":" + env.DatabasePort + ")" + "/" + env.DatabaseName

}

func (m MySqlLoader) BuildDsn(env *Env) (error, string) {

	mysqlPattern := patterns.CreateDSNPatterns().MySQL

	dsn := m.buildDsnFromEnvStruct(env)
	re := regexp.MustCompile(mysqlPattern)

	if re.MatchString(dsn) {
		return nil, dsn

	}
	return fmt.Errorf("bad connection string %s", dsn), ""

}

// Testing
func (m MySqlLoader) BuildDsnFromEnv(path string) (error, string) {
	fmt.Printf("MysqlLoader BuildDsnFromEnv %s\n", path)
	env := LoadDatabaseEnvVariables(path)
	if len(env.DatabaseUrl) == 0 || len(env.DatabaseUrl) == 0 {
		panic(fmt.Errorf("Env file not loaded"))
	}
	dsn := m.buildDsnFromEnvStruct(env)
	mysqlPattern := patterns.CreateDSNPatterns().MySQL

	re := regexp.MustCompile(mysqlPattern)

	if re.MatchString(dsn) {
		return nil, dsn
	}

	return fmt.Errorf("bad connection string %s", dsn), ""

}
