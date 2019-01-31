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
	"os"

	"cirello.io/bookmarkd/pkg/actions"
	"cirello.io/bookmarkd/pkg/errors"
	"gopkg.in/urfave/cli.v1"
)

func (c *commands) importBookmarks() cli.Command {
	return cli.Command{
		Name:        "import",
		Usage:       "import bookmarks",
		Description: "import all bookmarks",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name: "filename,f",
			},
		},
		Action: func(ctx *cli.Context) error {
			fd, err := os.Open(ctx.String("filename"))
			if err != nil {
				return errors.E(ctx, err, "cannot open file")
			}

			if err := actions.ImportBookmarks(c.db, fd); err != nil {
				return errors.E(ctx, err, "cannot import bookmarks")
			}

			return nil
		},
	}
}
