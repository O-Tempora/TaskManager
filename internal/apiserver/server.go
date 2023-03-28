package apiserver

import (
	"dip/internal/handlers"
	"dip/internal/store"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
)

type server struct {
	router *chi.Mux
	logger zerolog.Logger
	store  store.Store
}

func newServer(store store.Store) *server {
	s := &server{
		router: chi.NewRouter(),
		logger: zerolog.New(os.Stdout),
		store:  store,
	}

	s.initLogger(os.Stdout)
	s.initRouter()
	return s
}

func (s *server) initLogger(wr io.Writer) {
	output := zerolog.ConsoleWriter{
		Out:        wr,
		NoColor:    false,
		TimeFormat: time.ANSIC,
		FormatLevel: func(i interface{}) string {
			return strings.ToUpper(fmt.Sprintf("[%s]", i))
		},
		FormatTimestamp: func(i interface{}) string {
			t, _ := time.Parse(time.RFC3339, fmt.Sprintf("%s", i))
			return fmt.Sprintf("%s", t.Format(time.RFC1123))
		},
	}
	s.logger = zerolog.New(output).With().Timestamp().Logger().Level(zerolog.DebugLevel)
}

// To implement interface http.Handler
func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}
func (s *server) error(w http.ResponseWriter, r *http.Request, code int, err error) {
	s.success(w, r, code, map[string]string{"error": err.Error()})
}
func (s *server) success(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.WriteHeader(code)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}

func (s *server) initRouter() {
	s.router.Post("/signup", s.handleSignUp())
	s.router.Post("/login", s.handleLogIn())
}

func (s *server) handleSignUp() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		err, code := handlers.SignUp(s.store, w, r)
		if err != nil {
			s.error(w, r, code, err)
		} else {
			s.success(w, r, code, nil)
		}
	}
}
func (s *server) handleLogIn() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		err, code := handlers.LogIn(s.store, w, r)
		if err != nil {
			s.error(w, r, code, err)
		} else {
			s.success(w, r, code, nil)
		}
	}
}
