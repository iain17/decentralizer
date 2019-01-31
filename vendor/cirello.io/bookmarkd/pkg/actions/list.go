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

package actions // import "cirello.io/bookmarkd/pkg/actions"

import (
	"cirello.io/bookmarkd/pkg/errors"
	"cirello.io/bookmarkd/pkg/models"
	"github.com/jmoiron/sqlx"
)

// ListBookmarks list all bookmarks.
func ListBookmarks(db *sqlx.DB) ([]*models.Bookmark, error) {
	bookmarkDAO := models.NewBookmarkDAO(db)
	bookmarks, err := bookmarkDAO.All()
	if err != nil {
		return nil, errors.E(errors.Internal, err, "cannot load all bookmarks")
	}
	return bookmarks, nil
}
