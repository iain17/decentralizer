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

package cli

import (
	"fmt"

	"cirello.io/errors"
	"cirello.io/snippetsd/pkg/infra/repositories"
	"cirello.io/snippetsd/pkg/models/user"
	"gopkg.in/urfave/cli.v1"
)

func (c *commands) addUser() cli.Command {
	return cli.Command{
		Name:  "add",
		Usage: "add a user",
		Action: func(ctx *cli.Context) error {
			u, err := user.NewFromEmail(ctx.Args().Get(0), ctx.Args().Get(1), ctx.Args().Get(2))
			if err != nil {
				return errors.E(ctx, err, "cannot create user from email")
			}

			if _, err := repositories.Users(c.db).Insert(u); err != nil {
				return errors.E(ctx, err, "cannot store the new user")
			}

			fmt.Println(u, "added")
			return nil
		},
	}
}
