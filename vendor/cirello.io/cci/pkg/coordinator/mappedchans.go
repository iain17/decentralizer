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

package coordinator

import (
	"sync"

	"cirello.io/cci/pkg/models"
)

// TODO: remove when using proper queues.
type mappedChans struct {
	m sync.Map // map of models.Build.RepoFullName to chan *models.Build
}

func (mc *mappedChans) ch(repoFullName string) chan *models.Build {
	ch := make(chan *models.Build, 10)
	foundCh, _ := mc.m.LoadOrStore(repoFullName, ch)
	return foundCh.(chan *models.Build)
}
