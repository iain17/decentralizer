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

// Package git is wrapper to $PATH/bin.
package git // import "cirello.io/cci/pkg/infra/git"

import (
	"context"
	"log"
	"os"
	"os/exec"

	"cirello.io/errors"
)

// Checkout clones and reset build directory to a given commit.
func Checkout(ctx context.Context, cloneURL, repoDir, commit string) error {
	if _, err := os.Stat(repoDir); os.IsNotExist(err) {
		os.MkdirAll(repoDir, os.ModePerm&0755)
		cmd := exec.CommandContext(ctx, "git", "clone", cloneURL, ".")
		cmd.Dir = repoDir
		out, err := cmd.CombinedOutput()
		log.Println("cloning:", string(out))
		if err != nil {
			return errors.E(err, "cannot clone repository")
		}
	}
	cmd := exec.CommandContext(ctx, "git", "fetch", "--all")
	cmd.Dir = repoDir
	out, err := cmd.CombinedOutput()
	log.Println("fetching objects", string(out))
	if err != nil {
		return errors.E(err, "cannot fetch objects")
	}
	cmd = exec.CommandContext(ctx, "git", "reset", "--hard", commit)
	cmd.Dir = repoDir
	out, err = cmd.CombinedOutput()
	log.Println("reset to", string(out))
	return errors.E(err, "cannot reconfigure repository")
}
