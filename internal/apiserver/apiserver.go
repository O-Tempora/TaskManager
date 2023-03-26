package apiserver

import (
	"dip/internal/store"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
)

type APIServer struct {
	config *Config
	logger zerolog.Logger
	router *chi.Mux
	store  *store.Store
}

func New(config *Config) *APIServer {
	return &APIServer{
		config: config,
		logger: zerolog.New(os.Stdout),
		router: chi.NewRouter(),
	}
}

func (s *APIServer) Start() error {
	s.logger = initLogger(os.Stdout)
	s.logger.Info().Msg("Initialized")

	s.initRouter()

	if err := s.initStore(); err != nil {
		return err
	}

	return http.ListenAndServe(s.config.Port, s.router)
}

// Logger initialization
func initLogger(wr io.Writer) zerolog.Logger {
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
	return zerolog.New(output).With().Timestamp().Logger().Level(zerolog.DebugLevel)
}

// Store initialization
func (s *APIServer) initStore() error {
	st := store.New(s.config.Store)
	if err := st.Open(); err != nil {
		return err
	}
	s.store = st
	return nil
}

// Router initialization
func (s *APIServer) initRouter() {
	s.router.HandleFunc("/hello", s.handleHello())
}

func (s *APIServer) handleHello() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "Test Hello")
	}
}
