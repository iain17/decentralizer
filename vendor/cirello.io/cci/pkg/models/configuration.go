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
	"io"

	"cirello.io/cci/pkg/grpc/api"
	"cirello.io/errors"
	yaml "gopkg.in/yaml.v2"
)

// Configuration defines the internal parameters for the application.
type Configuration map[string]api.Recipe

// LoadConfiguration loads a given fd with YAML content into Configuration.
func LoadConfiguration(r io.Reader) (Configuration, error) {
	var c Configuration
	err := yaml.NewDecoder(r).Decode(&c)
	for k, v := range c {
		if v.Concurrency == 0 {
			v.Concurrency = 1
			c[k] = v
		}
	}
	return c, errors.E(err, "cannot parse configuration")
}
