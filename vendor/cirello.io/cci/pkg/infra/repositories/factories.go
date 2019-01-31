// Copyright 2018 github.com/ucirello
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

// Package repositories implements the models repositories of the application,
// it derives the underlying implementation according to the database driver
// injected in.
package repositories // import "cirello.io/cci/pkg/infra/repositories"

import (
	sqlite3 "cirello.io/cci/pkg/infra/repositories/internal/sqlite3"
	"cirello.io/cci/pkg/models"
	"github.com/jmoiron/sqlx"
	sqlite3Driver "github.com/mattn/go-sqlite3"
)

// Builds creates a Builds repository.
func Builds(db *sqlx.DB) models.BuildRepository {
	switch db.Driver().(type) {
	case *sqlite3Driver.SQLiteDriver:
		return sqlite3.NewBuildDAO(db)
	default:
		panic("invalid DB driver")
	}
}
