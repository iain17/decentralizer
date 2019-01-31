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

package cli

import (
	"cirello.io/bookmarkd/pkg/actions"
	"cirello.io/bookmarkd/pkg/errors"
	"gopkg.in/urfave/cli.v1"
)

func (c *commands) addBookmark() cli.Command {
	return cli.Command{
		Name:        "add",
		Aliases:     []string{"post"},
		Usage:       "add a bookmark",
		Description: "add a bookmarks",
		Action: func(ctx *cli.Context) error {
			if err := actions.AddBookmarkByURL(c.db, ctx.Args().First()); err != nil {
				return errors.E(ctx, err, "cannot add the bookmark")
			}
			return nil
		},
	}
}
