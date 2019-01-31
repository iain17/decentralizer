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

package actions

import (
	"cirello.io/bookmarkd/pkg/errors"
	"cirello.io/bookmarkd/pkg/models"
	"github.com/jmoiron/sqlx"
)

// DeleteBookmark deletes one bookmark from the database.
func DeleteBookmark(db *sqlx.DB, b *models.Bookmark, broadcast func(interface{})) error {
	err := models.NewBookmarkDAO(db).Delete(b)
	if err != nil {
		return errors.E(errors.Internal, err)
	}
	broadcast(&struct {
		WSType string `json:"type"`
		ID     int64  `json:"id"`
	}{
		"BOOKMARK_DELETED",
		b.ID,
	})
	return nil
}
