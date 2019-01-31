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
	"bytes"
)

type crbuffer struct {
	buf []byte
}

func (c *crbuffer) String() string {
	return string(c.buf)
}

func (c *crbuffer) Write(p []byte) (int, error) {
	for _, b := range p {
		switch b {
		case '\r':
			c.buf = c.buf[:bytes.LastIndexByte(c.buf, '\n')+1]
		default:
			c.buf = append(c.buf, b)
		}
	}
	return len(p), nil
}
