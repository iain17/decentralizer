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

package tasks

import (
	"context"
	"testing"
	"time"

	"cirello.io/bookmarkd/pkg/errors"
	"github.com/jmoiron/sqlx"
)

func Test_run(t *testing.T) {
	count := 0
	tasks := []Task{
		{
			Name: "overlapping tasks",
			Exec: func(*sqlx.DB) error {
				count++
				time.Sleep(5 * time.Second)
				return errors.E("fake error")
			},
			Frequency: 1 * time.Second,
		},
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	run(ctx, nil, tasks)
	<-ctx.Done()
	if count != 1 {
		t.Error("overlaping protection is not working. count:", count)
	}
}
