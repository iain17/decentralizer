// Copyright 2018 github.com/ucirello
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package web // import "cirello.io/snippetsd/pkg/ui/web"

import (
	"encoding/json"
	"log"
	"net/http"

	"cirello.io/snippetsd/pkg/infra/repositories"
	"cirello.io/snippetsd/pkg/models/snippet"
	"cirello.io/snippetsd/pkg/models/user"
	"github.com/jmoiron/sqlx"
)

// Server implements the web interface.
type Server struct {
	db  *sqlx.DB
	mux *http.ServeMux
}

// New creates a web interface handler.
func New(db *sqlx.DB) *Server {
	s := &Server{
		db:  db,
		mux: http.NewServeMux(),
	}
	s.registerRoutes()
	return s
}

func (s *Server) registerRoutes() {
	s.mux.HandleFunc("/storeSnippet", s.storeSnippet)
	s.mux.HandleFunc("/state", s.state)
	s.mux.HandleFunc("/", http.NotFound)
}

func (s *Server) unauthorized(w http.ResponseWriter) {
	w.Header().Add("WWW-Authenticate", `Basic realm="snippetsd"`)
	w.WriteHeader(http.StatusUnauthorized)
}

// ServeHTTP process HTTP requests.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// TODO: handle Access-Control-Allow-Origin correctly
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5200")
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	email, password, ok := r.BasicAuth()
	if !ok {
		s.unauthorized(w)
		return
	}

	u, err := repositories.Users(s.db).GetByEmail(email)
	if err != nil {
		s.unauthorized(w)
		return
	}

	if err := user.Authenticate(u, email, password); err != nil {
		s.unauthorized(w)
		return
	}

	r = r.WithContext(user.WithContext(r.Context(), u))
	s.mux.ServeHTTP(w, r)
}

func (s *Server) state(w http.ResponseWriter, r *http.Request) {
	snippets, err := repositories.Snippets(s.db).All()
	if err != nil {
		log.Println("cannot load all snippets:", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "    ")
	if err := enc.Encode(snippets); err != nil {
		log.Println("cannot marshal snippets:", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}
}

func (s *Server) storeSnippet(w http.ResponseWriter, r *http.Request) {
	user := user.WhoAmI(r.Context())

	var req struct {
		Contents string
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Println("cannot parse snippet storage request:", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}

	snippet := snippet.New(user, req.Contents)
	savedSnippet, err := repositories.Snippets(s.db).Save(snippet)
	if err != nil {
		log.Println("cannot save user's saved snippet:", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "    ")
	if err := enc.Encode(savedSnippet); err != nil {
		log.Println("cannot marshal saved snippet:", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}
}
