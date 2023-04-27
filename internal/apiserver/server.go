package apiserver

import (
	"dip/internal/handlers"
	"dip/internal/middleware"
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
			return t.Format(time.RFC1123)
		},
	}
	s.logger = zerolog.New(output).With().Timestamp().Logger().Level(zerolog.DebugLevel)
}

// To implement interface http.Handler
func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}
func (s *server) respond(w http.ResponseWriter, r *http.Request, code int, data interface{}, err error) {
	w.WriteHeader(code)
	if err != nil {
		response := map[string]string{"error": err.Error()}
		json.NewEncoder(w).Encode(response)
		s.logger.Error().Msgf("Resonse: method  %s, URL  %s, code  %d %s, error  %s",
			r.Method, r.URL, code, http.StatusText(code), err.Error())
		return
	}

	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
	s.logger.Info().Msgf("Response: method  %s, URL  %s, Code  %d %s",
		r.Method, r.URL, code, http.StatusText(code))
}

func (s *server) initRouter() {
	s.router.Use(middleware.LogRequest(s.logger))

	s.router.Post("/signup", s.handleSignUp())
	s.router.Post("/login", s.handleLogIn())

	s.router.Group(func(r chi.Router) {
		r.Use(middleware.AuthorizeToken())
		r.Get("/home", s.handleHome())
		r.Get("/home/{id}/{ws}", s.handleWorkspace())
	})
}

func (s *server) handleSignUp() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		err, code := handlers.SignUp(s.store, w, r)
		s.respond(w, r, code, nil, err)
	}
}
func (s *server) handleLogIn() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		code, token, err := handlers.LogIn(s.store, w, r)
		s.respond(w, r, code, map[string]string{
			"accessToken": token,
		}, err)
	}
}

func (s *server) handleHome() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		tp, err := middleware.ParseCredentials(r.Context())
		if err != nil {
			s.respond(w, r, http.StatusUnauthorized, nil, err)
			return
		}
		ws, code, err := handlers.GetHome(s.store, tp.Id)
		s.respond(w, r, code, ws, err)
	}
}

func (s *server) handleWorkspace() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		ws, code, err := handlers.GetFullWorkspace(s.store, chi.URLParam(r, "ws"))
		s.respond(w, r, code, ws, err)
	}
}
