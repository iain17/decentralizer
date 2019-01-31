// Copyright 2018 github.com/ucirello
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package repositories // import "cirello.io/snippetsd/pkg/infra/repositories"
import (
	sqlite3Snippets "cirello.io/snippetsd/pkg/infra/repositories/internal/sqlite3/snippets"
	sqlite3Users "cirello.io/snippetsd/pkg/infra/repositories/internal/sqlite3/users"
	"cirello.io/snippetsd/pkg/models/snippet"
	"cirello.io/snippetsd/pkg/models/user"
	"github.com/jmoiron/sqlx"
	sqlite3Driver "github.com/mattn/go-sqlite3"
)

// Snippets creates a Snippet repository.
func Snippets(db *sqlx.DB) snippet.Repository {
	switch db.Driver().(type) {
	case *sqlite3Driver.SQLiteDriver:
		return sqlite3Snippets.NewRepository(db)
	default:
		panic("invalid DB driver")
	}
}

// Users creates a User repository.
func Users(db *sqlx.DB) user.Repository {
	switch db.Driver().(type) {
	case *sqlite3Driver.SQLiteDriver:
		return sqlite3Users.NewRepository(db)
	default:
		panic("invalid DB driver")
	}
}
