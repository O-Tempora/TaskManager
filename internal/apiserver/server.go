package apiserver

import (
	"dip/internal/middleware"
	"dip/internal/models"
	"dip/internal/service"
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
	router  *chi.Mux
	logger  zerolog.Logger
	store   store.Store
	service service.IService
}

func newServer(store store.Store, service service.IService) *server {
	s := &server{
		router:  chi.NewRouter(),
		logger:  zerolog.New(os.Stdout),
		store:   store,
		service: service,
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

	//s.router.Get("/test", s.TEST())
	s.router.Post("/signup", s.handleSignUp())
	s.router.Post("/login", s.handleLogIn())
	s.router.Route("/workspace", func(r chi.Router) {
		r.With(middleware.AuthorizeToken()).Get("/", s.handleHome())
		r.With(middleware.AuthorizeToken()).Post("/", s.handleCreateWS())
		r.Get("/{id}", s.handleWorkspace())
		r.Put("/{id}", s.handleUpdateWS())
		r.Delete("/{id}", s.handleDeleteWS())
	})
	s.router.Route("/task", func(r chi.Router) {
		r.Post("/", s.handleCreateTask())
		r.Get("/{id}-{ws}", s.handleTask())
		r.Put("/{id}", s.handleUpdateTask())
		r.Put("/{id}-to-{gr}", s.handleMoveTask())
		r.Delete("/{id}", s.handleDeleteTask())
		r.With(middleware.AuthorizeToken()).Get("/", s.handleGetPersonalTasks())

		r.Route("/{task}/comment", func(r chi.Router) {
			r.Get("/", s.handleGetComments())
			r.Post("/", s.handleCreateComment())
			r.Delete("/{id}", s.handleDeleteComment())
		})
	})
	s.router.Route("/group", func(r chi.Router) {
		r.Post("/", s.handleCreateGroup())
		r.Put("/{id}", s.handleUpdateGroup())
		r.Delete("/{id}", s.handleDeleteGroup())
	})
	s.router.Route("/person", func(r chi.Router) {
		r.Get("/ws-{ws}", s.handleAllPersonsInWs())
		r.With(middleware.AuthorizeToken()).Get("/isAdmin-{ws}", s.handleIsAdmin())
		r.With(middleware.AuthorizeToken()).Get("/byToken", s.handleGetPerson())
		r.Post("/{name}/assign-{task}", s.handleAssign())
		r.Put("/{id}", s.handleUpdatePerson())
		r.Delete("/{name}/dismiss-{task}", s.handleDismiss())
		r.Delete("/{id}/{ws}", s.handleLeaveWS())
	})
	s.router.Route("/status", func(r chi.Router) {
		r.Get("/", s.handleStatuses())
	})
	s.router.Route("/maintainer", func(r chi.Router) {
		r.Use(middleware.AuthorizeToken(), middleware.VerifyMaintainer)
		r.Route("/person", func(r chi.Router) {
			r.Get("/", s.handleGetAllPersons())
			r.Delete("/{id}", s.handleDeletePerson())
		})
		r.Route("/workspace", func(r chi.Router) {
			r.Get("/", s.handleGetAllWS())
		})
	})
	s.router.Route("/invite", func(r chi.Router) {
		r.With(middleware.AuthorizeToken()).Get("/", s.handleGetInvites())
		r.With(middleware.AuthorizeToken()).Post("/", s.handleSendInvite())
		r.With(middleware.AuthorizeToken()).Post("/{id}/{ws}", s.handleAcceptInvite())
		r.Delete("/{id}", s.handleDeclineInvite())
	})
}

//	func (s *server) TEST() http.HandlerFunc {
//		return func(w http.ResponseWriter, r *http.Request) {
//			w.Header().Set("Content-Type", "application/json")
//			res, err := s.store.Task().GetAllByUser(1)
//			s.respond(w, r, 322, res, err)
//		}
//	}
func (s *server) handleSignUp() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		type request struct {
			Name     string `json:"name"`
			Email    string `json:"email"`
			Phone    string `json:"phone"`
			Password string `json:"password"`
		}
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.respond(w, r, http.StatusBadRequest, nil, err)
			return
		}
		if err := s.service.Logic().SignUp(req.Name, req.Email, req.Phone, req.Password); err != nil {
			s.respond(w, r, http.StatusInternalServerError, nil, err)
			return
		}
		s.respond(w, r, http.StatusOK, nil, nil)
	}
}
func (s *server) handleLogIn() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		type request struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.respond(w, r, http.StatusBadRequest, nil, err)
			return
		}
		token, err := s.service.Logic().LogIn(req.Email, req.Password)
		if err != nil {
			s.respond(w, r, http.StatusInternalServerError, nil, err)
			return
		}
		s.respond(w, r, http.StatusOK, map[string]string{
			"accessToken": token,
		}, nil)
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
		ws, err := s.service.Logic().GetWsByUser(tp.Id)
		if err != nil {
			s.respond(w, r, http.StatusUnauthorized, nil, err)
			return
		}
		s.respond(w, r, http.StatusOK, ws, nil)
	}
}
func (s *server) handleWorkspace() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			s.respond(w, r, http.StatusBadRequest, nil, err)
			return
		}
		ws, err := s.service.Logic().GetFullWorkspace(id)
		if err != nil {
			s.respond(w, r, http.StatusInternalServerError, nil, err)
			return
		}
		s.respond(w, r, http.StatusBadRequest, ws, nil)
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
		persons, err := s.service.Logic().GetAllAssignedToTask(id, ws)
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
		isAdmin, err := s.service.Logic().IsAdmin(tp.Login, id)
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
		err = s.service.Logic().Assign(chi.URLParam(r, "name"), id)
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
		err = s.service.Logic().Dismiss(chi.URLParam(r, "name"), id)
		if err != nil {
			s.respond(w, r, http.StatusInternalServerError, nil, err)
			return
		}
		s.respond(w, r, http.StatusOK, nil, nil)
	}
}
func (s *server) handleUpdateWS() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			s.respond(w, r, http.StatusBadRequest, nil, err)
			return
		}
		var ws *models.Workspace
		if err := json.NewDecoder(r.Body).Decode(&ws); err != nil {
			s.respond(w, r, http.StatusBadRequest, nil, err)
			return
		}
		if err := s.store.Workspace().Update(ws, id); err != nil {
			s.respond(w, r, http.StatusInternalServerError, nil, err)
			return
		}
		s.respond(w, r, http.StatusOK, nil, nil)
	}
}
func (s *server) handleDeleteWS() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			s.respond(w, r, http.StatusBadRequest, nil, err)
			return
		}
		if err = s.store.Workspace().Delete(id); err != nil {
			s.respond(w, r, http.StatusInternalServerError, nil, err)
			return
		}
		s.respond(w, r, http.StatusOK, nil, nil)
	}
}
func (s *server) handleGetPersonalTasks() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		tp, err := middleware.ParseCredentials(r.Context())
		if err != nil {
			s.respond(w, r, http.StatusUnauthorized, nil, err)
			return
		}
		res, err := s.service.Logic().GetAllTasksByUser(tp.Id)
		if err != nil {
			s.respond(w, r, http.StatusNotFound, nil, err)
			return
		}
		s.respond(w, r, http.StatusOK, res, nil)
	}
}
func (s *server) handleDeletePerson() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			s.respond(w, r, http.StatusBadRequest, nil, err)
			return
		}
		if err = s.store.Person().Delete(id); err != nil {
			s.respond(w, r, http.StatusInternalServerError, nil, err)
			return
		}
		s.respond(w, r, http.StatusOK, nil, nil)
	}
}
func (s *server) handleUpdatePerson() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			s.respond(w, r, http.StatusBadRequest, nil, err)
			return
		}
		var p = models.Person{}
		if err = json.NewDecoder(r.Body).Decode(&p); err != nil {
			s.respond(w, r, http.StatusBadRequest, nil, err)
			return
		}
		if err := p.Validate(); err != nil {
			s.respond(w, r, http.StatusBadRequest, nil, err)
			return
		}
		if err = s.store.Person().Update(id, p); err != nil {
			s.respond(w, r, http.StatusInternalServerError, nil, err)
			return
		}
		s.respond(w, r, http.StatusOK, nil, nil)
	}
}
func (s *server) handleGetComments() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		task, err := strconv.Atoi(chi.URLParam(r, "task"))
		if err != nil {
			s.respond(w, r, http.StatusBadRequest, nil, err)
			return
		}
		comments, err := s.store.Comment().GetByTask(task)
		if err != nil {
			s.respond(w, r, http.StatusInternalServerError, nil, err)
			return
		}
		s.respond(w, r, http.StatusOK, comments, nil)
	}
}
func (s *server) handleCreateComment() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		cm := models.Comment{}
		err := json.NewDecoder(r.Body).Decode(&cm)
		if err != nil {
			s.respond(w, r, http.StatusBadRequest, nil, err)
			return
		}
		c, err := s.store.Comment().Create(cm)
		if err != nil {
			s.respond(w, r, http.StatusInternalServerError, nil, err)
			return
		}
		s.respond(w, r, http.StatusOK, c, nil)
	}
}
func (s *server) handleDeleteComment() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			s.respond(w, r, http.StatusBadRequest, nil, err)
			return
		}
		err = s.store.Comment().Delete(id)
		if err != nil {
			s.respond(w, r, http.StatusInternalServerError, nil, err)
			return
		}
		s.respond(w, r, http.StatusOK, nil, nil)
	}
}
func (s *server) handleMoveTask() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			s.respond(w, r, http.StatusBadRequest, nil, err)
			return
		}
		gr, err := strconv.Atoi(chi.URLParam(r, "gr"))
		if err != nil {
			s.respond(w, r, http.StatusBadRequest, nil, err)
			return
		}
		if err = s.service.Logic().MoveTask(id, gr); err != nil {
			s.respond(w, r, http.StatusBadRequest, nil, err)
			return
		}
		s.respond(w, r, http.StatusOK, nil, nil)
	}
}
func (s *server) handleLeaveWS() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		var nextAdmin int
		na := r.URL.Query().Get("next")

		//lacks "next" flag in request -> sent by user
		if na == "" {
			nextAdmin = -1
		} else {
			//else sent by admin -> make new admin
			next_id, err := strconv.Atoi(na)
			if err != nil {
				s.respond(w, r, http.StatusBadRequest, nil, err)
				return
			}
			nextAdmin = next_id
		}

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
		if err = s.service.Logic().LeaveWs(id, ws, nextAdmin); err != nil {
			s.respond(w, r, http.StatusInternalServerError, nil, err)
			return
		}
		s.respond(w, r, http.StatusOK, nil, nil)
	}
}
func (s *server) handleGetAllPersons() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		params := r.URL.Query()
		page, err := strconv.Atoi(params.Get("page"))
		if err != nil {
			s.respond(w, r, http.StatusBadRequest, nil, err)
			return
		}
		take, err := strconv.Atoi(params.Get("take"))
		if err != nil {
			s.respond(w, r, http.StatusBadRequest, nil, err)
			return
		}
		p, err := s.store.Person().GetAll(page, take)
		if err != nil {
			s.respond(w, r, http.StatusInternalServerError, nil, err)
			return
		}
		payload := struct {
			Persons []models.PersonShow
			Total   int
		}{
			Persons: p,
			Total:   len(p),
		}
		s.respond(w, r, http.StatusOK, payload, nil)
	}
}
func (s *server) handleGetAllWS() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		params := r.URL.Query()
		page, err := strconv.Atoi(params.Get("page"))
		if err != nil {
			s.respond(w, r, http.StatusBadRequest, nil, err)
			return
		}
		take, err := strconv.Atoi(params.Get("take"))
		if err != nil {
			s.respond(w, r, http.StatusBadRequest, nil, err)
			return
		}
		ws, err := s.store.Workspace().GetAll(page, take)
		if err != nil {
			s.respond(w, r, http.StatusInternalServerError, nil, err)
			return
		}
		payload := struct {
			WS    []models.Workspace
			Total int
		}{
			WS:    ws,
			Total: len(ws),
		}
		s.respond(w, r, http.StatusOK, payload, nil)
	}
}
func (s *server) handleGetInvites() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		tp, err := middleware.ParseCredentials(r.Context())
		if err != nil {
			s.respond(w, r, http.StatusUnauthorized, nil, err)
			return
		}
		inv, err := s.store.Invite().GetAll(tp.Id)
		if err != nil {
			s.respond(w, r, http.StatusInternalServerError, nil, err)
			return
		}
		s.respond(w, r, http.StatusOK, inv, nil)
	}
}
func (s *server) handleSendInvite() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		tp, err := middleware.ParseCredentials(r.Context())
		if err != nil {
			s.respond(w, r, http.StatusUnauthorized, nil, err)
			return
		}
		payload := struct {
			Email string `json:"email"`
			WsId  int    `json:"ws_id"`
		}{}
		if err = json.NewDecoder(r.Body).Decode(&payload); err != nil {
			s.respond(w, r, http.StatusBadRequest, nil, err)
			return
		}
		if err := s.service.Logic().SendInvite(payload.Email, payload.WsId, tp.Id); err != nil {
			s.respond(w, r, http.StatusInternalServerError, nil, err)
			return
		}
		s.respond(w, r, http.StatusOK, nil, nil)
	}
}
func (s *server) handleDeclineInvite() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			s.respond(w, r, http.StatusBadRequest, nil, err)
			return
		}
		if err := s.service.Logic().DeclineInvite(id); err != nil {
			s.respond(w, r, http.StatusInternalServerError, nil, err)
			return
		}
		s.respond(w, r, http.StatusOK, nil, nil)
	}
}
func (s *server) handleAcceptInvite() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		tp, err := middleware.ParseCredentials(r.Context())
		if err != nil {
			s.respond(w, r, http.StatusUnauthorized, nil, err)
			return
		}
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
		newWs, err := s.service.Logic().AcceptInvite(id, ws, tp.Id)
		if err != nil {
			s.respond(w, r, http.StatusInternalServerError, nil, err)
			return
		}
		s.respond(w, r, http.StatusOK, newWs, nil)
	}
}
