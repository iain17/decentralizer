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
	"compress/gzip"
	"net/http"
	"strings"
	"time"

	"cirello.io/bookmarkd/pkg/models"
	"github.com/PuerkitoBio/goquery"
)

var now = time.Now

// CheckLink dials bookmark URL and updates its state with the errors if any.
func CheckLink(bookmark *models.Bookmark) *models.Bookmark {
	res, err := http.Get(bookmark.URL)
	if err != nil {
		bookmark.LastStatusCheck = now().Unix()
		bookmark.LastStatusCode = 0
		bookmark.LastStatusReason = err.Error()
		return bookmark
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		bookmark.LastStatusCheck = now().Unix()
		bookmark.LastStatusCode = int64(res.StatusCode)
		bookmark.LastStatusReason = http.StatusText(res.StatusCode)
		return bookmark
	}

	bookmark.LastStatusCode = int64(res.StatusCode)
	bookmark.LastStatusReason = http.StatusText(res.StatusCode)
	bookmark.LastStatusCheck = now().Unix()

	isHTML := strings.Contains(res.Header.Get("Content-Type"), "text/html")
	if bookmark.Title != "" || !isHTML {
		return bookmark
	}

	isGzipped := strings.Contains(res.Header.Get("Content-Encoding"), "gzip")
	if isGzipped {
		res.Body, err = gzip.NewReader(res.Body)
		if err != nil {
			bookmark.LastStatusReason = "not gzip content"
			return bookmark
		}
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		bookmark.LastStatusReason = "cannot parse body"
		return bookmark
	}

	doc.Find("HEAD>TITLE").Each(func(i int, s *goquery.Selection) {
		bookmark.Title = s.Text()
	})

	return bookmark
}
