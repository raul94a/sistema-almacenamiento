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
 "gorm.io/driver/sqlite" // Sqlite driver based on CGO
  // "github.com/glebarez/sqlite" // Pure go SQLite driver, checkout https://github.com/glebarez/sqlite for details
  "gorm.io/gorm"

)

type SqliteLoader struct{}

func (m SqliteLoader) LoadDatabase(connectionString string) (*gorm.DB, error) {
	// EXAMPLE file:test.db?cache=shared&mode=memory
	return  gorm.Open(sqlite.Open(connectionString), &gorm.Config{})

}
