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

package client

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"cirello.io/cci/pkg/grpc/api"
	"cirello.io/errors"
)

const execScript = `#!/bin/bash

set -e

%s
`

func run(ctx context.Context, recipe *api.Recipe, repoDir, baseDir string) (string, error) {
	tmpfile, err := ioutil.TempFile(repoDir, "agent")
	if err != nil {
		return "", errors.E(errors.FailedPrecondition, err,
			"agent cannot create temporary file")
	}
	defer os.Remove(tmpfile.Name())
	defer tmpfile.Close()
	fmt.Fprintf(tmpfile, execScript, recipe.Commands)
	tmpfile.Close()
	if recipe.Timeout != nil && *recipe.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, *recipe.Timeout)
		defer cancel()
	}
	cmd := exec.CommandContext(ctx, "/bin/sh", tmpfile.Name())
	cmd.Dir = repoDir
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, fmt.Sprintf("CCI_BUILD_BASE_DIRECTORY=%s", baseDir))
	recipeEnvVars := strings.Split(recipe.Environment, "\n")
	for _, v := range recipeEnvVars {
		cmd.Env = append(cmd.Env, os.Expand(v, expandVar(cmd.Env)))
	}
	var buf crbuffer
	cmd.Stdout = io.MultiWriter(&buf, os.Stdout)
	cmd.Stderr = io.MultiWriter(&buf, os.Stderr)
	err = cmd.Run()
	return buf.String(), errors.E(err, "failed when running builder")
}

func expandVar(currentEnv []string) func(string) string {
	return func(s string) string {
		for _, e := range currentEnv {
			if strings.HasPrefix(e, s+"=") {
				return strings.TrimPrefix(e, s+"=")
			}
		}
		return ""
	}
}
