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
	"text/template"
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
	s.logger.Info().Msgf("Method: [%s] URL: [%s] Code: [%d %s]",
		r.Method, r.URL, code, http.StatusText(code))
}

func (s *server) initRouter() {
	s.router.Post("/signup", s.handleSignUp())
	s.router.Get("/signup", s.getSignup())
	s.router.Post("/login", s.handleLogIn())
	s.router.Get("/login", s.getLogin())

	s.router.Get("/tasks", s.handleTasks())

	s.router.Get("/statuses", s.handleGetStatuses())
	s.router.Get("/statusId", s.handlestatusById())
}

func (s *server) getSignup() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tpl, err := template.ParseFiles("templates/signup.html")
		err = tpl.Execute(w, nil)
		if err != nil {
			s.error(w, r, 500, err)
		}
	}
}
func (s *server) handleSignUp() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		err, code := handlers.SignUp(s.store, w, r)
		if err != nil {
			s.error(w, r, code, err)
			return
		} else {
			s.success(w, r, code, nil)
		}
	}
}

func (s *server) getLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tpl, err := template.ParseFiles("templates/login.html")
		err = tpl.Execute(w, nil)
		if err != nil {
			s.error(w, r, 500, err)
		}
	}
}
func (s *server) handleLogIn() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		err, code := handlers.LogIn(s.store, w, r)
		if err != nil {
			s.error(w, r, code, err)
			return
		} else {
			s.success(w, r, code, nil)
		}
	}
}

func (s *server) handleGetStatuses() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		statuses, err := s.store.Status().GetAll()
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		} else {
			s.success(w, r, http.StatusOK, statuses)
		}
	}
}

func (s *server) handlestatusById() http.HandlerFunc {
	var name string
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewDecoder(r.Body).Decode(&name); err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		id, err := s.store.Status().GetIdByName(name)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.success(w, r, http.StatusOK, id)
	}
}

func (s *server) handleTasks() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		tasks, err := s.store.Task().GetAll()
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
		s.success(w, r, http.StatusOK, tasks)
	}
}
