package api

import (
	"database/sql"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	apiHandlers "github.com/pavbis/zal-case-study/api/handlers"
	"github.com/pavbis/zal-case-study/application/storage"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	// Register postgres driver
	_ "github.com/lib/pq"
)

const apiPathPrefix = "/api"

var (
	userName = os.Getenv("AUTH_USER")
	password = os.Getenv("AUTH_PASS")
)

// Server represents server
type Server struct {
	router *mux.Router
	logger *log.Logger
	db     *sql.DB
}

// Initialize initializes the server with necessary deps
func (s *Server) Initialize() {
	s.router = mux.NewRouter().PathPrefix(apiPathPrefix).Subrouter()
	s.logger = log.New(os.Stdout, "", log.LstdFlags)

	var err error
	s.db, err = sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		s.logger.Fatal(err)
	}

	err = s.db.Ping()
	if err != nil {
		s.logger.Fatal(err)
	}

	s.initializeRoutes()
}

// Run starts the server on the provided port
func (s *Server) Run(addr string) {
	loggedRouter := s.createLoggingRouter(s.logger.Writer())

	srv := &http.Server{
		Handler:      loggedRouter,
		Addr:         addr,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	s.logger.Fatal(srv.ListenAndServe())
}

// Post wraps the router for POST method
func (s *Server) Post(path string, f func(w http.ResponseWriter, r *http.Request)) {
	s.router.HandleFunc(path, f).Methods(http.MethodPost)
}

// Get wraps the router for Get method
func (s *Server) Get(path string, f func(w http.ResponseWriter, r *http.Request)) {
	s.router.HandleFunc(path, apiHandlers.BasicAuthMiddleware(userName, password, f)).Methods(http.MethodGet)
}

func (s *Server) createLoggingRouter(out io.Writer) http.Handler {
	return handlers.LoggingHandler(out, s.router)
}

func (s *Server) initializeRoutes() {
	// Health check
	s.router.HandleFunc("/health", apiHandlers.HealthRequestHandler).Methods(http.MethodGet)

	// Language
	s.Post("/languages/{languageName}", s.handleRequestWithDBInstance(apiHandlers.ReceiveRepositoriesRequestHandler))
	s.Get("/languages/{languageName}", s.handleRequestWithDBInstance(apiHandlers.ReadRepositoriesRequestHandler))
	s.Get("/languages", s.handleRequestWithDBInstance(apiHandlers.ListLanguagesAndRepositoriesRequestHandler))
	s.Get("/stats/count-repositories", s.handleRequestWithDBInstance(apiHandlers.CountRepositoriesStarsForLanguagesRequestHandler))

	// Repositories
	s.Post("/repositories/{repositoryId}", s.handleRequestWithDBInstance(apiHandlers.RemoveRepositoryRequestHandler))
	s.Get("/stats/top-list", s.handleRequestWithDBInstance(apiHandlers.TopRepositoryForLanguageRequestHandler))
}

// RequestHandlerFunction is the function which represents any handler
type RequestHandlerFunction func(db storage.Executor, w http.ResponseWriter, r *http.Request)

func (s *Server) handleRequestWithDBInstance(handler RequestHandlerFunction) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handler(s.db, w, r)
	}
}
