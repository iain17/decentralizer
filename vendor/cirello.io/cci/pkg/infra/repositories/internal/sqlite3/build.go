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

package sqlite3 // import "cirello.io/cci/pkg/infra/repositories/internal/sqlite3"

import (
	"time"

	"cirello.io/cci/pkg/models"
	"cirello.io/errors"
	"github.com/jmoiron/sqlx"
)

// BuildDAO provides access to the collection of Builds.
type BuildDAO struct {
	db *sqlx.DB
}

// NewBuildDAO creates a new Build data access object.
func NewBuildDAO(db *sqlx.DB) *BuildDAO {
	return &BuildDAO{
		db: db,
	}
}

// Bootstrap creates the necessary table to operate builds.
func (b *BuildDAO) Bootstrap() error {
	ops := []string{
		`CREATE TABLE IF NOT EXISTS builds (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			repo_full_name text,
			commit_hash text,
			commit_message text,
			environment text,
			commands text,
			started_at datetime default (datetime('now')) not null,
			success bool default false not null,
			log bigtext default '' not null,
			completed_at datetime default '' not null
		);`,
		`CREATE INDEX builds_repo_full_nam ON builds (repo_full_name)`,
		`CREATE INDEX builds_started_at ON builds (started_at)`,
	}
	for _, op := range ops {
		_, err := b.db.Exec(op)
		if err != nil {
			errors.E(err, "cannot bootstrap database")
		}
	}
	return nil
}

// Register stores an new build in the database.
func (b *BuildDAO) Register(build *models.Build) (*models.Build, error) {
	res, err := b.db.NamedExec(`
		INSERT INTO builds
		(repo_full_name, commit_hash, commit_message, environment, commands)
		VALUES
		(:repo_full_name, :commit_hash, :commit_message, :environment, :commands)
	`, build)
	if err != nil {
		return build, errors.E(err, "cannot add job to database")
	}
	id, err := res.LastInsertId()
	if err != nil {
		return build, errors.E(err, "cannot load ID from the added job")
	}
	build.ID = id
	return build, nil
}

// MarkInProgress determines a build has started and update its build
// information in the database.
func (b *BuildDAO) MarkInProgress(build *models.Build) error {
	now := time.Now()
	build.StartedAt = &now
	_, err := b.db.NamedExec(`
		UPDATE builds
		SET started_at = :started_at
		WHERE id = :id
	`, build)
	return errors.E(err, "cannot mark build to in-progress")
}

// MarkComplete determines a build has completed and update its build
// information in the database.
func (b *BuildDAO) MarkComplete(build *models.Build) error {
	now := time.Now()
	build.CompletedAt = &now
	_, err := b.db.NamedExec(`
		UPDATE builds
		SET
			success = :success,
			log = :log,
			completed_at = :completed_at
		WHERE id = :id AND completed_at = ''
	`, build)
	return errors.E(err, "cannot mark build to complete")
}

// GetLastBuild loads last known build for a repository.
func (b *BuildDAO) GetLastBuild(repoFullName string) (*models.Build, error) {
	var build models.Build
	err := b.db.Get(&build, `
		SELECT
			id,
			repo_full_name,
			commit_hash,
			commit_message,
			environment,
			commands,
			started_at,
			success,
			log,
			completed_at
		FROM builds
		WHERE repo_full_name = :repoFullName
		ORDER BY started_at DESC
		LIMIT 1
	`, repoFullName)
	return &build, errors.E(err, "cannot load last known build")
}

// ListByRepoFullName all builds for a repository
func (b *BuildDAO) ListByRepoFullName(repoFullName string) ([]*models.Build, error) {
	var builds []*models.Build
	err := b.db.Select(&builds, `
		SELECT
			id,
			repo_full_name,
			commit_hash,
			commit_message,
			environment,
			commands,
			started_at,
			success,
			log,
			completed_at
		FROM builds
		WHERE repo_full_name = :repoFullName
		ORDER BY started_at DESC
	`, repoFullName)
	return builds, errors.E(err, "cannot load builds for repository")
}

// SweepExpired mark expired builds as failed.
func (b *BuildDAO) SweepExpired(timeout time.Duration) (int64, error) {
	now := time.Now()
	resp, err := b.db.Exec(`
		UPDATE builds
		SET
			success = 0,
			completed_at = $1,
			log = 'timeout'
		WHERE
			started_at < $2 AND
			completed_at = '' AND
			success = 0
	`, now, now.Add(-timeout))
	if err != nil {
		return 0, errors.E(err, "cannot mark build to complete")
	}
	rows, err := resp.RowsAffected()
	if err != nil {
		return 0, errors.E(err, "cannot count rows affected")
	}
	return rows, nil
}
