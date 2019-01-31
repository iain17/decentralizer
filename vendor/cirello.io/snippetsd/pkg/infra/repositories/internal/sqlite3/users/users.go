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

package users // import "cirello.io/snippetsd/pkg/infra/repositories/internal/sqlite3/users"
import (
	"cirello.io/errors"
	"cirello.io/snippetsd/pkg/models/user"
	"github.com/jmoiron/sqlx"
)

// Repository provides DB persistence to snippets users.
type Repository struct {
	db *sqlx.DB
}

// NewRepository instanties a Repository
func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

// Bootstrap creates table if missing.
func (b *Repository) Bootstrap() error {
	cmds := []string{
		`create table if not exists users (
			id integer primary key autoincrement,
			email varchar(255),
			password varchar(255),
			team varchar(255),
			constraint users_email_unique unique (email)
		);
		`,
		`create index if not exists users_email on users (email)`,
	}

	for _, cmd := range cmds {
		_, err := b.db.Exec(cmd)
		if err != nil {
			return errors.E(err)
		}
	}

	return nil
}

// GetByID loads a user by ID.
func (b *Repository) GetByID(id int64) (*user.User, error) {
	var u user.User
	err := b.db.Get(&u, "SELECT * FROM users WHERE id = $1", id)
	return &u, errors.E(err)
}

// GetByEmail loads a user by email.
func (b *Repository) GetByEmail(email string) (*user.User, error) {
	var u user.User
	err := b.db.Get(&u, "SELECT * FROM users WHERE email = $1", email)
	return &u, errors.E(err)
}

// Insert a user.
func (b *Repository) Insert(user *user.User) (*user.User, error) {
	if err := user.EncryptPassword(); err != nil {
		return nil, errors.E(err, "cannot encrypt password")
	}
	_, err := b.db.NamedExec(`
		INSERT INTO users
		(email, password, team)
		VALUES (:email, :password, :team)
	`, user)
	if err != nil {
		return nil, errors.E(err)
	}
	user.Password = "" // avoid leaking

	err = b.db.Get(user, `
		SELECT
			*
		FROM
			users
		WHERE
			id = last_insert_rowid()
	`)
	if err != nil {
		return nil, errors.E(err)
	}
	return user, errors.E(err)
}
