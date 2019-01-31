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

package net

import (
	"reflect"
	"testing"
	"time"

	"cirello.io/bookmarkd/pkg/models"
)

func TestCheckLink(t *testing.T) {
	oldNow := now
	now = func() time.Time {
		return time.Unix(0, 0)
	}
	type args struct {
		bookmark *models.Bookmark
	}
	tests := []struct {
		name string
		args args
		want *models.Bookmark
	}{
		{
			"404",
			args{
				&models.Bookmark{URL: "http://example.com/404"},
			},
			&models.Bookmark{
				URL:              "http://example.com/404",
				LastStatusCode:   404,
				LastStatusReason: "Not Found",
				LastStatusCheck:  now().Unix(),
			},
		},
		{
			"200",
			args{
				&models.Bookmark{URL: "http://example.com/"},
			},
			&models.Bookmark{
				URL:              "http://example.com/",
				LastStatusCode:   200,
				LastStatusReason: "OK",
				LastStatusCheck:  now().Unix(),
				Title:            "Example Domain",
			},
		},
		{
			"invalid URL",
			args{
				&models.Bookmark{URL: "invalid-url"},
			},
			&models.Bookmark{
				URL:              "invalid-url",
				LastStatusCode:   0,
				LastStatusReason: "Get invalid-url: unsupported protocol scheme \"\"",
				LastStatusCheck:  now().Unix(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CheckLink(tt.args.bookmark); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("%s CheckLink() = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
	now = oldNow
}

func TestContentExtraction(t *testing.T) {
	b := &models.Bookmark{
		URL: "https://www.example.org",
	}
	b = CheckLink(b)
	if b.Title != "Example Domain" {
		t.Fatal("cannot extract HTML title")
	}
}
