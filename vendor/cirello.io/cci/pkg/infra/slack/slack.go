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

// Package slack is a service to operate Slack integrations.
package slack // import "cirello.io/cci/pkg/infra/slack"

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	"cirello.io/errors"
)

// Sends a message to a given Slack webhook.
func Send(webhookURL string, msg string) error {
	var payload bytes.Buffer
	err := json.NewEncoder(&payload).Encode(struct {
		Text string `json:"text"`
	}{Text: msg})
	if err != nil {
		return errors.E(err, "cannot encode slack message")
	}
	response, err := http.Post(webhookURL, "application/json", &payload)
	if err != nil {
		return errors.E(err, "cannot send slack message")
	}
	if _, err := io.Copy(ioutil.Discard, response.Body); err != nil {
		return errors.E(err, "cannot drain response body")
	}
	return nil
}
