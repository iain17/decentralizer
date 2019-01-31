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

// Package dashboard implements a build web dashboard.
package dashboard // import "cirello.io/cci/pkg/ui/dashboard"

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net"
	"net/http"
	"strings"
	"time"

	"cirello.io/cci/pkg/models"
	"cirello.io/errors"
)

// Server implements the dashboard.
type Server struct {
	buildDAO models.BuildRepository
}

// New creates a new builds dashboard.
func New(buildDAO models.BuildRepository) *Server {
	return &Server{
		buildDAO: buildDAO,
	}
}

// ServeContext exposes the build dashboard.
func (s *Server) ServeContext(ctx context.Context, l net.Listener) error {
	srv := &http.Server{}
	mux := http.NewServeMux()
	mux.HandleFunc("/dashboard/", s.listBuilds)
	go func() {
		select {
		case <-ctx.Done():
			ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
			defer cancel()
			srv.Shutdown(ctx)
		}
	}()
	srv.Handler = mux
	go func() {
		select {
		case <-ctx.Done():
			ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
			defer cancel()
			srv.Shutdown(ctx)
		}
	}()
	return errors.E(srv.Serve(l), "error when serving dashboard interface")
}

func (s *Server) listBuilds(w http.ResponseWriter, r *http.Request) {
	repoFullName := strings.TrimPrefix(r.RequestURI, "/dashboard/")
	builds, err := s.buildDAO.ListByRepoFullName(repoFullName)
	if err != nil && errors.RootCause(err) != sql.ErrNoRows {
		log.Println("cannot load builds:", err)
	}
	if err := json.NewEncoder(w).Encode(builds); err != nil {
		log.Println("cannot encode builds:", err)
	}
}
