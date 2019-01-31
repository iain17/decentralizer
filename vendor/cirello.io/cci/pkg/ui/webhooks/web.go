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

// Package webhooks implements the public interface used to communicate with
// Github.
package webhooks // import "cirello.io/cci/pkg/ui/webhooks"

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"strings"
	"time"

	"cirello.io/cci/pkg/coordinator"
	"cirello.io/cci/pkg/models"
	"cirello.io/errors"
)

// Server implements the webhooks part of the CI service. For now, compatible
// only with Github Webhooks.
type Server struct {
	coordinator *coordinator.Coordinator
}

// New creates a new web-facing server.
func New(c *coordinator.Coordinator) *Server {
	return &Server{
		coordinator: c,
	}
}

// ServeContext handles the HTTP requests.
func (s *Server) ServeContext(ctx context.Context, l net.Listener) error {
	srv := &http.Server{}
	mux := http.NewServeMux()
	mux.HandleFunc("/github-webhook/", func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println("cannot read payload body:", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
			return
		}
		var payload githubHookPayload
		if err := json.Unmarshal(body, &payload); err != nil {
			log.Println("cannot parse payload:", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
			return
		}
		sig := r.Header.Get("X-Hub-Signature")
		if err := s.coordinator.Enqueue(payload.Repository.FullName,
			payload.CommitHash, payload.HeadCommit.Message,
			sig, body); err != nil {
			log.Println("cannot enqueue payload:", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
			return
		}
	})
	mux.HandleFunc("/badge/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/svg+xml;charset=utf-8")
		repoFullName := strings.TrimPrefix(r.RequestURI, "/badge/")
		status := s.coordinator.GetLastBuildStatus(repoFullName)
		badge := badgeUnknown
		switch status {
		case models.Success:
			badge = badgePassing
		case models.Failed:
			badge = badgeFailing
		case models.InProgress:
			badge = badgeRunning
		}
		fmt.Fprint(w, badge)
	})
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		out, err := httputil.DumpRequest(r, true)
		fmt.Println(string(out), err)
	})
	srv.Handler = mux
	go func() {
		select {
		case <-ctx.Done():
			ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
			defer cancel()
			srv.Shutdown(ctx)
		}
	}()
	return errors.E(srv.Serve(l), "error when serving web interface")
}

type githubHookPayload struct {
	CommitHash string `json:"after"`
	HeadCommit struct {
		Message string `json:"message"`
	} `json:"head_commit"`
	Repository struct {
		FullName string `json:"full_name"`
	} `json:"repository"`
	Sender struct {
		Login     string
		AvatarURL string `json:"avatar_url"`
	}
}

const badgePassing = `<svg xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" width="88" height="20">
<linearGradient id="b" x2="0" y2="100%"><stop offset="0" stop-color="#bbb" stop-opacity=".1"/><stop offset="1" stop-opacity=".1"/></linearGradient>
<clipPath id="a"><rect width="88" height="20" rx="3" fill="#fff"/></clipPath>
<g clip-path="url(#a)"><path fill="#555" d="M0 0h37v20H0z"/><path fill="#4c1" d="M37 0h51v20H37z"/><path fill="url(#b)" d="M0 0h88v20H0z"/></g>
<g fill="#fff" text-anchor="middle" font-family="DejaVu Sans,Verdana,Geneva,sans-serif" font-size="110">
<text x="195" y="150" fill="#010101" fill-opacity=".3" transform="scale(.1)" textLength="270">build</text>
<text x="195" y="140" transform="scale(.1)" textLength="270">build</text>
<text x="615" y="150" fill="#010101" fill-opacity=".3" transform="scale(.1)" textLength="410">passing</text>
<text x="615" y="140" transform="scale(.1)" textLength="410">passing</text>
</g>
</svg>`

const badgeFailing = `<svg xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" width="80" height="20">
<linearGradient id="b" x2="0" y2="100%"><stop offset="0" stop-color="#bbb" stop-opacity=".1"/><stop offset="1" stop-opacity=".1"/></linearGradient>
<clipPath id="a"><rect width="80" height="20" rx="3" fill="#fff"/></clipPath>
<g clip-path="url(#a)"><path fill="#555" d="M0 0h37v20H0z"/><path fill="#e05d44" d="M37 0h43v20H37z"/><path fill="url(#b)" d="M0 0h80v20H0z"/></g>
<g fill="#fff" text-anchor="middle" font-family="DejaVu Sans,Verdana,Geneva,sans-serif" font-size="110">
<text x="195" y="150" fill="#010101" fill-opacity=".3" transform="scale(.1)" textLength="270">build</text>
<text x="195" y="140" transform="scale(.1)" textLength="270">build</text>
<text x="575" y="150" fill="#010101" fill-opacity=".3" transform="scale(.1)" textLength="330">failing</text>
<text x="575" y="140" transform="scale(.1)" textLength="330">failing</text>
</g>
</svg>`

const badgeRunning = `<svg xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" width="90" height="20">
<linearGradient id="b" x2="0" y2="100%"><stop offset="0" stop-color="#bbb" stop-opacity=".1"/><stop offset="1" stop-opacity=".1"/></linearGradient>
<clipPath id="a"><rect width="90" height="20" rx="3" fill="#fff"/></clipPath>
<g clip-path="url(#a)"><path fill="#555" d="M0 0h37v20H0z"/><path fill="#9f9f9f" d="M37 0h53v20H37z"/><path fill="url(#b)" d="M0 0h90v20H0z"/></g>
<g fill="#fff" text-anchor="middle" font-family="DejaVu Sans,Verdana,Geneva,sans-serif" font-size="110">
<text x="195" y="150" fill="#010101" fill-opacity=".3" transform="scale(.1)" textLength="270">build</text>
<text x="195" y="140" transform="scale(.1)" textLength="270">build</text>
<text x="625" y="150" fill="#010101" fill-opacity=".3" transform="scale(.1)" textLength="430">running</text>
<text x="625" y="140" transform="scale(.1)" textLength="430">running</text>
</g>
</svg>`

const badgeUnknown = `<svg xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" width="90" height="20">
<linearGradient id="b" x2="0" y2="100%"><stop offset="0" stop-color="#bbb" stop-opacity=".1"/><stop offset="1" stop-opacity=".1"/></linearGradient>
<clipPath id="a"><rect width="90" height="20" rx="3" fill="#fff"/></clipPath>
<g clip-path="url(#a)"><path fill="#555" d="M0 0h37v20H0z"/><path fill="#9f9f9f" d="M37 0h53v20H37z"/><path fill="url(#b)" d="M0 0h90v20H0z"/></g>
<g fill="#fff" text-anchor="middle" font-family="DejaVu Sans,Verdana,Geneva,sans-serif" font-size="110">
<text x="195" y="150" fill="#010101" fill-opacity=".3" transform="scale(.1)" textLength="270">build</text>
<text x="195" y="140" transform="scale(.1)" textLength="270">build</text>
<text x="625" y="150" fill="#010101" fill-opacity=".3" transform="scale(.1)" textLength="430">unknown</text>
<text x="625" y="140" transform="scale(.1)" textLength="430">unknown</text>
</g>
</svg>`
