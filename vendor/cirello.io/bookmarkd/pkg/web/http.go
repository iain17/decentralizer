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

package web // import "cirello.io/bookmarkd/pkg/web"

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"cirello.io/bookmarkd/generated"
	"cirello.io/bookmarkd/pkg/actions"
	"cirello.io/bookmarkd/pkg/errors"
	"cirello.io/bookmarkd/pkg/models"
	"cirello.io/bookmarkd/pkg/net"
	"cirello.io/bookmarkd/pkg/pubsub"
	svcjwt "cirello.io/svc/pkg/jwt"
	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/websocket"
	"github.com/jmoiron/sqlx"
)

// Server implements the web interface.
type Server struct {
	db               *sqlx.DB
	handler          http.Handler
	pubsub           *pubsub.Broker
	jwtSecret        []byte
	acceptableEmails map[string]struct{}
	authMiddleware   *jwtmiddleware.JWTMiddleware
}

// New creates a web interface handler.
func New(db *sqlx.DB, jwtSecret []byte, acceptableEmails []string) (*Server, error) {
	s := &Server{
		db:        db,
		pubsub:    pubsub.New(),
		jwtSecret: jwtSecret,
	}
	s.acceptableEmails = make(map[string]struct{})
	for _, e := range acceptableEmails {
		s.acceptableEmails[e] = struct{}{}
	}

	if len(jwtSecret) > 0 {
		s.authMiddleware = jwtmiddleware.New(jwtmiddleware.Options{
			UserProperty: "user",
			ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
				return s.jwtSecret, nil
			},
			SigningMethod: jwt.SigningMethodHS512,
			ErrorHandler:  s.unauthorized,
		})
	}

	err := s.registerRoutes()
	return s, err
}

func (s *Server) registerRoutes() error {
	rootFS := generated.AssetFS()
	rootFS.Prefix = "frontend/build/"
	index, err := rootFS.Open("index.html")
	if err != nil {
		return err
	}
	rootHandler := http.FileServer(rootFS)

	router := http.NewServeMux()
	router.HandleFunc("/state", s.state)
	router.HandleFunc("/loadBookmark", s.loadBookmark)
	router.HandleFunc("/newBookmark", s.newBookmark)
	router.HandleFunc("/deleteBookmark", s.deleteBookmark)
	router.HandleFunc("/ws", s.websocket)
	router.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		rootHandler.ServeHTTP(&recoverableResponseWriter{
			responseWriter: w,
			request:        req,
			fallback: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "text/html; charset=utf-8")
				http.ServeContent(w, r, "index.html", time.Now(), index)
			},
		}, req)
	})

	s.handler = router
	return nil
}

func (s *Server) unauthorized(w http.ResponseWriter, r *http.Request, err string) {
	log.Println("access denied", err)
	w.WriteHeader(http.StatusUnauthorized)
}

type authClaims string

var trustLevel = authClaims("trust-level")

func (s *Server) authentication(w http.ResponseWriter, r *http.Request) error {
	if s.authMiddleware == nil {
		return nil
	}
	if err := s.authMiddleware.CheckJWT(w, r); err != nil {
		return errors.E("cannot find JWT in the request")
	}
	token, ok := r.Context().Value("user").(*jwt.Token)
	if !ok {
		return errors.E("cannot find token in context")
	}
	claims, err := svcjwt.Claims(token)
	if err != nil {
		return errors.E(err, "unexpected token set")
	}
	if claims.Target != "bookmarkd.cirello.io" {
		return errors.E("invalid target in token")
	}
	if _, ok := s.acceptableEmails[claims.Email]; !ok {
		return errors.E("access for this account")
	}
	*r = *r.WithContext(context.WithValue(r.Context(),
		trustLevel, claims.Trust))
	return nil
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := s.authentication(w, r); err != nil {
		s.unauthorized(w, r, err.Error())
		return
	}
	s.handler.ServeHTTP(w, r)
}

func (s *Server) state(w http.ResponseWriter, r *http.Request) {
	// TODO: handle Access-Control-Allow-Origin correctly
	w.Header().Set("Access-Control-Allow-Origin", "*")
	bookmarks, err := actions.ListBookmarks(s.db)
	if err != nil {
		log.Println("cannot load all bookmarks:", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(bookmarks); err != nil {
		log.Println("cannot marshal bookmarks:", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (s *Server) loadBookmark(w http.ResponseWriter, r *http.Request) {
	// TODO: handle Access-Control-Allow-Origin correctly
	w.Header().Set("Access-Control-Allow-Origin", "*")

	bookmark := &models.Bookmark{}
	if err := json.NewDecoder(r.Body).Decode(bookmark); err != nil {
		log.Println("cannot unmarshal bookmark request:", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}

	bookmark = net.CheckLink(bookmark)

	if err := json.NewEncoder(w).Encode(bookmark); err != nil {
		log.Println("cannot marshal bookmark:", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}
}

func (s *Server) newBookmark(w http.ResponseWriter, r *http.Request) {
	// TODO: handle Access-Control-Allow-Origin correctly
	w.Header().Set("Access-Control-Allow-Origin", "*")

	bookmark := &models.Bookmark{}
	if err := json.NewDecoder(r.Body).Decode(bookmark); err != nil {
		log.Println("cannot unmarshal bookmark request:", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}

	if err := actions.AddBookmark(s.db, bookmark, s.pubsub.Broadcast); err != nil {
		log.Println("cannot save bookmark:", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}

	w.Write([]byte("{}"))
}

func (s *Server) deleteBookmark(w http.ResponseWriter, r *http.Request) {
	if v, ok := r.Context().Value(trustLevel).(string); !ok || v != "high" {
		http.Error(w, http.StatusText(http.StatusUnauthorized),
			http.StatusUnauthorized)
		return
	}

	// TODO: handle Access-Control-Allow-Origin correctly
	w.Header().Set("Access-Control-Allow-Origin", "*")

	bookmark := &models.Bookmark{}
	if err := json.NewDecoder(r.Body).Decode(bookmark); err != nil {
		log.Println("cannot unmarshal bookmark request:", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}

	if err := actions.DeleteBookmark(s.db, bookmark, s.pubsub.Broadcast); err != nil {
		log.Println("cannot save bookmark:", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}

	w.Write([]byte("{}"))
}

func (s *Server) websocket(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		ReadBufferSize:  4096,
		WriteBufferSize: 4096,
	}

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("cannot upgrade to websocket:", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}
	if ws != nil {
		defer ws.Close()
	}

	unsubscribe := s.pubsub.Subscribe(func(msg interface{}) {
		if err := ws.WriteJSON(msg); err != nil {
			log.Println("cannot write websocket message:", err)
			ws.Close()
		}
	})
	defer unsubscribe()
	defer ws.Close()

	log.Println("listening for pings...")
	for {
		msgType, _, err := ws.NextReader()
		if err != nil || msgType == websocket.CloseMessage {
			return
		}
	}
}
