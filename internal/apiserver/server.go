package apiserver

import (
	"dip/internal/handlers"
	"dip/internal/models"
	"dip/internal/store"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
)

var tmp = template.Must(template.ParseGlob("templates/*"))

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
	s.router.Get("/delete", s.deleteTask())

	s.router.Get("/new", s.addNew())
	s.router.Post("/create", s.createTask())

	s.router.Get("/edit", s.editTask())
	s.router.Post("/update", s.updateTask())

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
			//http.Redirect(w, r, "/tasks", code)
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
		//w.Header().Set("Content-Type", "application/json")
		tasks, _ := s.store.Task().GetAll()
		// tpl, _ := template.ParseFiles("templates/tasks.html")
		// tpl.Execute(w, tasks)
		tmp.ExecuteTemplate(w, "Index", tasks)
		// if err != nil {
		// 	s.error(w, r, http.StatusInternalServerError, err)
		// 	return
		// }
		// s.success(w, r, http.StatusOK, tasks)
	}
}

func (s *server) deleteTask() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		value, err := strconv.Atoi(id)
		if err != nil {
			s.error(w, r, 500, err)
			return
		}

		s.store.Task().DeleteTask(value)
		http.Redirect(w, r, "/tasks", 301)
	}
}

func (s *server) createTask() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		t := &models.JoinedTask{}
		// if err := json.NewDecoder(r.Body).Decode(t); err != nil {
		// 	s.error(w, r, 400, err)
		// 	return
		// }
		t.Description = r.FormValue("descr")
		t.StartAt = r.FormValue("date1")
		t.FinishAt = r.FormValue("date2")
		t.Status = r.FormValue("status")
		status, err := s.store.Status().GetIdByName(t.Status)
		if err != nil {
			s.error(w, r, 400, err)
			return
		}

		s.store.Task().Create(t, status)
		http.Redirect(w, r, "/tasks", 301)
	}
}

func (s *server) updateTask() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		t := &models.JoinedTask{}
		t.Id, _ = strconv.Atoi(r.FormValue("id"))
		t.Description = r.FormValue("description")
		t.StartAt = r.FormValue("date1")
		t.FinishAt = r.FormValue("date2")
		t.Status = r.FormValue("status")
		status, err := s.store.Status().GetIdByName(t.Status)
		if err != nil {
			s.error(w, r, 400, err)
			return
		}

		err = s.store.Task().Update(t, status)
		if err != nil {
			s.error(w, r, 400, err)
		}
		http.Redirect(w, r, "/tasks", 301)
	}
}

func (s *server) addNew() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		statuses, _ := s.store.Status().GetAll()
		names := make([]string, 0)
		for _, v := range statuses {
			names = append(names, v.Name)
		}
		tmp.ExecuteTemplate(w, "New", names)
	}
}

func (s *server) editTask() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type respone struct {
			Task     models.JoinedTask
			Statuses []string
		}
		buf := &respone{models.JoinedTask{}, make([]string, 0)}
		id := r.URL.Query().Get("id")
		i, _ := strconv.Atoi(id)
		task, err := s.store.Task().Get(i)
		if err != nil {
			s.error(w, r, 500, err)
			return
		}
		statuses, _ := s.store.Status().GetAll()
		//names := make([]string, 0)
		for _, v := range statuses {
			buf.Statuses = append(buf.Statuses, v.Name)
		}

		buf.Task = task

		tmp.ExecuteTemplate(w, "Edit", buf)
	}
}
