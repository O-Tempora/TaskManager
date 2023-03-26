package apiserver

import (
	"dip/internal/models"
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

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *server) initRouter() {
	s.router.Post("/persons", s.handleCreatePerson())
}
func (s *server) handleCreatePerson() http.HandlerFunc {
	type request struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Phone    string `json:"phone"`
		Password string `json:"password"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		p := &models.Person{
			Email:    req.Email,
			Password: req.Password,
			Name:     req.Name,
			Phone:    req.Phone,
			Settings: "",
		}

		if err := s.store.Person().Create(p); err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		s.respond(w, r, http.StatusCreated, p)
	}
}

func (s *server) error(w http.ResponseWriter, r *http.Request, code int, err error) {
	s.respond(w, r, code, map[string]string{"error": err.Error()})
}
func (s *server) respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.WriteHeader(code)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}
