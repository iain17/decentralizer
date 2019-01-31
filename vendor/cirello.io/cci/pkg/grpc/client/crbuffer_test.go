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
	"testing"
)

func Test_crbuffer_Write(t *testing.T) {
	type args struct {
		p []byte
	}
	tests := []struct {
		name string
		args args
		buf  []byte
	}{
		{"simple", args{[]byte("a")}, []byte("a")},
		{"ln", args{[]byte("a\nb")}, []byte("a\nb")},
		{"cr", args{[]byte("a\nb\r1")}, []byte("a\n1")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &crbuffer{}
			_, err := c.Write(tt.args.p)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if !bytes.Equal(c.buf, tt.buf) {
				t.Errorf("error processing buffer. got: %q\nexpected: %q", c.buf, tt.buf)
			}
		})
	}
}
