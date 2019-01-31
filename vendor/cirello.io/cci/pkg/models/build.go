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

package models

import (
	"time"

	"cirello.io/cci/pkg/grpc/api"
)

// Build defines the necessary data to run a build successfully.
type Build struct {
	*api.Build
	*api.Recipe
}

// Status define the current build status of a Build
type Status int

// Possible build status
const (
	Unknown Status = iota
	Success
	Failed
	InProgress
)

// Status returns the current status of the build.
func (b *Build) Status() Status {
	switch {
	case b.StartedAt.IsZero():
		return Unknown
	case b.CompletedAt.IsZero():
		return InProgress
	case b.Success:
		return Success
	default:
		return Failed
	}
}

// BuildRepository manipulate a collection of Build requests.
type BuildRepository interface {
	Bootstrap() error
	Register(build *Build) (*Build, error)
	MarkInProgress(build *Build) error
	MarkComplete(build *Build) error
	GetLastBuild(repoFullName string) (*Build, error)
	ListByRepoFullName(repoFullName string) ([]*Build, error)
	SweepExpired(timeout time.Duration) (int64, error)
}
