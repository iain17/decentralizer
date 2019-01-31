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

package snippets // import "cirello.io/snippetsd/pkg/infra/repositories/internal/sqlite3/snippets"
import (
	"cirello.io/errors"
	"cirello.io/snippetsd/pkg/infra/repositories/internal/sqlite3/users"
	"cirello.io/snippetsd/pkg/models/snippet"
	"github.com/jmoiron/sqlx"
)

// Repository provides a repository of Snippets.
type Repository struct {
	db *sqlx.DB
}

// NewRepository instanties a Repository
func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		db: db,
	}
}

// Bootstrap creates table if missing.
func (b *Repository) Bootstrap() error {
	cmds := []string{
		`create table if not exists snippets (
			id integer primary key autoincrement,
			user_id bigint,
			week_start datetime,
			contents bigtext
		);
		`,
		`create index if not exists snippets_user_id on snippets (user_id)`,
		`create index if not exists snippets_week_start on snippets (week_start)`,
		`create unique index if not exists snippets_user_id_week_start on snippets (user_id, week_start)`,
	}

	for _, cmd := range cmds {
		_, err := b.db.Exec(cmd)
		if err != nil {
			return errors.E(err, cmd)
		}
	}

	return nil
}

func (b *Repository) loadUsers(snippets *[]*snippet.Snippet) error {
	repo := users.NewRepository(b.db)
	for i, s := range *snippets {
		u, err := repo.GetByID(s.UserID)
		if err != nil {
			return errors.E(err, "cannot load snippets user")
		}
		s.User = u
		(*snippets)[i] = s
	}
	return nil
}

// All returns all known snippets.
func (b *Repository) All() ([]*snippet.Snippet, error) {
	var snippets []*snippet.Snippet
	err := b.db.Select(&snippets, "SELECT * FROM snippets ORDER BY week_start DESC")
	if err != nil {
		return snippets, errors.E(err, "cannot load snippets")
	}
	if err := b.loadUsers(&snippets); err != nil {
		return snippets, errors.E(err, "cannot load users information")
	}
	return snippets, nil
}

// Save one snippet entry.
func (b *Repository) Save(snippet *snippet.Snippet) (*snippet.Snippet, error) {
	_, err := b.db.NamedExec(`
		REPLACE INTO snippets (user_id, week_start, contents)
		VALUES (:user_id, :week_start, :contents)
	`, snippet)
	if err != nil {
		return nil, errors.E(err, "upsert operation failed")
	}

	return snippet, nil
}
