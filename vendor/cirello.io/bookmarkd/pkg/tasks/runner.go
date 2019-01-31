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

package tasks // import "cirello.io/bookmarkd/pkg/tasks"

import (
	"context"
	"log"
	"sync"
	"time"

	"cirello.io/bookmarkd/pkg/errors"
	"cirello.io/bookmarkd/pkg/models"
	"cirello.io/bookmarkd/pkg/net"
	"github.com/jmoiron/sqlx"
	"golang.org/x/sync/singleflight"
)

// Task represents one periodic task executed by the runner. Key must be unique,
// as it is used as key to lock.
type Task struct {
	Name      string
	Exec      func(db *sqlx.DB) error
	Frequency time.Duration
}

var execGroup singleflight.Group
var tasks = []Task{
	{"check link health", LinkHealth, 6 * time.Hour},
	{"vacuum", Vacuum, 24 * time.Hour},
}

// Run executes background maintenance tasks.
func Run(db *sqlx.DB) {
	run(context.Background(), db, tasks)
}

func run(ctx context.Context, db *sqlx.DB, tasks []Task) {
	for _, t := range tasks {
		t := t
		go func() {
			log.Println("tasks: scheduled", t.Name)
			for {
				go func() {
					_, err, _ := execGroup.Do(t.Name, func() (interface{}, error) {
						log.Println("tasks:", t.Name, "running")
						defer log.Println("tasks:", t.Name, "done")
						err := t.Exec(db)
						return nil, err
					})
					if err != nil {
						log.Println(t.Name, " failed:", err)
					}
				}()
				select {
				case <-time.After(t.Frequency):
				case <-ctx.Done():
				}
			}
		}()
	}
}

// LinkHealth checks if the expired links are still valid.
func LinkHealth(db *sqlx.DB) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.E(errors.Internal, r)
		}
	}()
	bookmarkDAO := models.NewBookmarkDAO(db)
	bookmarks, err := bookmarkDAO.Expired()
	if err != nil {
		return errors.E(errors.Internal, err, "cannot load expired bookmarks")
	}

	bookmarkCh := make(chan *models.Bookmark)
	var wg sync.WaitGroup
	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for bookmark := range bookmarkCh {
				log.Println("linkHealth:", bookmark.ID, bookmark.URL)
				bookmark = net.CheckLink(bookmark)
				if err := bookmarkDAO.Update(bookmark); err != nil {
					log.Println(err, "cannot update link during link health check - status OK")
				}
			}
		}()
	}
	for _, bookmark := range bookmarks {
		bookmarkCh <- bookmark
	}
	close(bookmarkCh)
	wg.Wait()

	return nil
}

// Vacuum executes a SQLite3 vacuum clean up.
func Vacuum(db *sqlx.DB) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.E(errors.Internal, r)
		}
	}()

	_, err = db.Exec("VACUUM")
	return errors.E(err, "cannot run vacuum")
}
