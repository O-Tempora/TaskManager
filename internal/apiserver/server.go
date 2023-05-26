package apiserver

import (
	"dip/internal/handlers"
	"dip/internal/middleware"
	"dip/internal/models"
	"dip/internal/store"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
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
	cors := cors.AllowAll()
	s.router.Use(cors.Handler)
	s.router.Use(middleware.LogRequest(s.logger))
	s.router.Post("/signup", s.handleSignUp())
	s.router.Post("/login", s.handleLogIn())
	s.router.Group(func(r chi.Router) {
		r.Use(middleware.AuthorizeToken())
		r.Route("/home", func(r chi.Router) {
			r.Get("/", s.handleHome())
			r.Post("/", s.handleCreateWS())
		})
	})
	s.router.Route("/workspace-{ws}", func(r chi.Router) {
		r.Get("/", s.handleWorkspace())
	})
	s.router.Route("/task", func(r chi.Router) {
		r.Post("/", s.handleCreateTask())
		r.Get("/{id}-{ws}", s.handleTask())
		r.Put("/{id}", s.handleUpdateTask())
		r.Delete("/{id}", s.handleDeleteTask())
	})
	s.router.Route("/group", func(r chi.Router) {
		r.Post("/", s.handleCreateGroup())
		r.Put("/{id}", s.handleUpdateGroup())
		r.Delete("/{id}", s.handleDeleteGroup())
	})
	s.router.Route("/person", func(r chi.Router) {
		r.With(middleware.AuthorizeToken()).Get("/isAdmin-{ws}", s.handleIsAdmin())
		r.Get("/ws-{ws}", s.handleAllPersonsInWs())
		r.With(middleware.AuthorizeToken()).Get("/byToken", s.handleGetPerson())
		r.Post("/{name}/assign-{task}", s.handleAssign())
		r.Delete("/{name}/dismiss-{task}", s.handleDismiss())
	})
	s.router.Route("/status", func(r chi.Router) {
		r.Get("/", s.handleStatuses())
	})
}

func (s *server) handleSignUp() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		code, err := handlers.SignUp(s.store, w, r)
		s.respond(w, r, code, nil, err)
	}
}
func (s *server) handleLogIn() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		code, token, err := handlers.LogIn(s.store, w, r)
		if err != nil {
			s.respond(w, r, code, nil, err)
			return
		}
		cookie := http.Cookie{
			Name:  "usr",
			Value: token,
		}
		http.SetCookie(w, &cookie)
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
		if err != nil {
			s.respond(w, r, code, nil, err)
			return
		}
		s.respond(w, r, code, ws, nil)
	}
}
func (s *server) handleTask() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			s.respond(w, r, http.StatusBadRequest, nil, err)
			return
		}
		ws, err := strconv.Atoi(chi.URLParam(r, "ws"))
		if err != nil {
			s.respond(w, r, http.StatusBadRequest, nil, err)
			return
		}
		task, err := s.store.Task().GetById(id)
		if err != nil {
			s.respond(w, r, http.StatusBadRequest, nil, err)
			return
		}
		persons, err := s.store.Person().GetAllAssignedToTask(id, ws)
		if err != nil {
			s.respond(w, r, http.StatusBadRequest, nil, err)
			return
		}

		t := struct {
			Task    *models.Task          `json:"task"`
			Persons []models.PersonInTask `json:"persons"`
		}{
			Task:    task,
			Persons: persons,
		}
		s.respond(w, r, http.StatusOK, t, nil)
	}
}
func (s *server) handleUpdateTask() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		t := &models.Task{}
		if err := json.NewDecoder(r.Body).Decode(t); err != nil {
			s.respond(w, r, http.StatusBadRequest, nil, err)
			return
		}
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			s.respond(w, r, http.StatusBadRequest, nil, err)
			return
		}
		t.Id = id
		err = s.store.Task().Update(t)
		if err != nil {
			s.respond(w, r, http.StatusInternalServerError, nil, err)
			return
		}
		s.respond(w, r, http.StatusOK, nil, nil)
	}
}
func (s *server) handleDeleteTask() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			s.respond(w, r, http.StatusBadRequest, nil, err)
			return
		}
		if err = s.store.Task().Delete(id); err != nil {
			s.respond(w, r, http.StatusInternalServerError, nil, err)
			return
		}
		s.respond(w, r, http.StatusOK, nil, nil)
	}
}
func (s *server) handleCreateTask() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		group_id := -1
		if err := json.NewDecoder(r.Body).Decode(&group_id); err != nil {
			s.respond(w, r, http.StatusBadRequest, nil, err)
			return
		}
		to, err := s.store.Task().Create(group_id)
		if err != nil {
			s.respond(w, r, http.StatusInternalServerError, nil, err)
			return
		}
		s.respond(w, r, http.StatusOK, to, nil)
	}
}
func (s *server) handleCreateGroup() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		tg := &models.TG{}
		if err := json.NewDecoder(r.Body).Decode(&tg); err != nil {
			s.respond(w, r, http.StatusBadRequest, nil, err)
			return
		}
		tg, err := s.store.TaskGroup().Create(tg.WorkspaceId, tg.Name, tg.Color)
		if err != nil {
			s.respond(w, r, http.StatusInternalServerError, nil, err)
			return
		}
		s.respond(w, r, http.StatusOK, tg, nil)
	}
}
func (s *server) handleUpdateGroup() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		tg := &models.TG{}
		if err := json.NewDecoder(r.Body).Decode(&tg); err != nil {
			s.respond(w, r, http.StatusBadRequest, nil, err)
			return
		}
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			s.respond(w, r, http.StatusBadRequest, nil, err)
			return
		}
		tg.Id = id
		err = s.store.TaskGroup().Update(tg)
		if err != nil {
			s.respond(w, r, http.StatusInternalServerError, nil, err)
			return
		}
		s.respond(w, r, http.StatusOK, nil, nil)
	}
}
func (s *server) handleDeleteGroup() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			s.respond(w, r, http.StatusBadRequest, nil, err)
			return
		}
		if err = s.store.TaskGroup().Delete(id); err != nil {
			s.respond(w, r, http.StatusInternalServerError, nil, err)
			return
		}
		s.respond(w, r, http.StatusOK, nil, nil)
	}
}
func (s *server) handleCreateWS() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		type pl struct {
			Name        string `json:"name"`
			Description string `json:"description"`
		}
		payload := &pl{}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			fmt.Println("Error - ", err)
			s.respond(w, r, http.StatusBadRequest, nil, err)
			return
		}
		tp, err := middleware.ParseCredentials(r.Context())
		if err != nil {
			s.respond(w, r, http.StatusUnauthorized, nil, err)
			return
		}

		ws, err := s.store.Workspace().Create(tp.Id, payload.Name, payload.Description)
		if err != nil {
			s.respond(w, r, http.StatusInternalServerError, nil, err)
			return
		}

		s.respond(w, r, http.StatusOK, *ws, nil)
	}
}
func (s *server) handleIsAdmin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		id, err := strconv.Atoi(chi.URLParam(r, "ws"))
		if err != nil {
			s.respond(w, r, http.StatusBadRequest, nil, err)
			return
		}

		tp, err := middleware.ParseCredentials(r.Context())
		if err != nil {
			s.respond(w, r, http.StatusUnauthorized, nil, err)
			return
		}
		isAdmin, err := s.store.Person().IsAdmin(tp.Login, id)
		if err != nil {
			s.respond(w, r, http.StatusUnauthorized, nil, err)
			return
		}
		s.respond(w, r, http.StatusOK, isAdmin, nil)
	}
}
func (s *server) handleStatuses() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		st, err := s.store.Status().GetAll()
		if err != nil {
			s.respond(w, r, http.StatusInternalServerError, nil, err)
			return
		}
		s.respond(w, r, http.StatusOK, st, nil)
	}
}
func (s *server) handleAllPersonsInWs() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		ws, err := strconv.Atoi(chi.URLParam(r, "ws"))
		if err != nil {
			s.respond(w, r, http.StatusBadRequest, nil, err)
			return
		}
		pl, err := s.store.Person().GetAllByWorkspace(ws)
		if err != nil {
			s.respond(w, r, http.StatusBadRequest, nil, err)
			return
		}
		s.respond(w, r, http.StatusOK, pl, nil)
	}
}
func (s *server) handleGetPerson() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		tp, err := middleware.ParseCredentials(r.Context())
		if err != nil {
			s.respond(w, r, http.StatusUnauthorized, nil, err)
			return
		}
		s.respond(w, r, http.StatusOK, tp, nil)
	}
}
func (s *server) handleAssign() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		id, err := strconv.Atoi(chi.URLParam(r, "task"))
		if err != nil {
			s.respond(w, r, http.StatusBadRequest, nil, err)
			return
		}
		err = s.store.Person().Assign(chi.URLParam(r, "name"), id)
		if err != nil {
			s.respond(w, r, http.StatusInternalServerError, nil, err)
			return
		}
		s.respond(w, r, http.StatusOK, nil, nil)
	}
}
func (s *server) handleDismiss() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		id, err := strconv.Atoi(chi.URLParam(r, "task"))
		if err != nil {
			s.respond(w, r, http.StatusBadRequest, nil, err)
			return
		}
		err = s.store.Person().Dismiss(chi.URLParam(r, "name"), id)
		if err != nil {
			s.respond(w, r, http.StatusInternalServerError, nil, err)
			return
		}
		s.respond(w, r, http.StatusOK, nil, nil)
	}
}
