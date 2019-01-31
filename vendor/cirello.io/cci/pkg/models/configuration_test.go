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
	"reflect"
	"strings"
	"testing"

	"cirello.io/cci/pkg/grpc/api"
)

func TestConfigurationParser(t *testing.T) {
	yaml := strings.NewReader(yamlSample)
	got, err := LoadConfiguration(yaml)
	if err != nil {
		t.Fatalf("cannot parse configuration: %v", err)
	}
	expected := Configuration{
		"org/account": api.Recipe{
			Concurrency:  2,
			Clone:        "git@github.com:org/account.git",
			SlackWebhook: "https://hooks.slack.com/services/AAAA/BBB/CCC",
			GithubSecret: "ghsecret",
			Environment:  "ENV1=1\nENV2=2\n",
			Commands:     "vgo test ./errors/... ./supervisor/...\necho OK\n",
		},
	}

	if !reflect.DeepEqual(got, expected) {
		t.Errorf("wrong parsing. got: %#v\nexpected:%#v\n", got, expected)
	}
}

const yamlSample = `---
org/account:
  concurrency: 2
  clone: git@github.com:org/account.git
  slack_webhook: https://hooks.slack.com/services/AAAA/BBB/CCC
  github_secret: ghsecret
  environment: |
    ENV1=1
    ENV2=2
  commands: |
    vgo test ./errors/... ./supervisor/...
    echo OK
`
