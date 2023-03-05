package api

import (
	"database/sql"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
	"os"
	"time"

	// nolint: goimports
	_ "github.com/lib/pq"

	apiHandlers "github.com/pavbis/repositories-api/api/handlers"
	"github.com/pavbis/repositories-api/application/storage"
)

var (
	userName = os.Getenv("AUTH_USER")
	password = os.Getenv("AUTH_PASS")
)

// Server represents server
type Server struct {
	router *chi.Mux
	logger *log.Logger
	db     *sql.DB
}

// Initialize initializes the server with necessary deps
func (s *Server) Initialize() {
	s.router = chi.NewRouter()
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
	srv := &http.Server{
		Handler:           s.router,
		Addr:              addr,
		WriteTimeout:      15 * time.Second,
		ReadTimeout:       15 * time.Second,
		IdleTimeout:       120 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
	}

	s.logger.Fatal(srv.ListenAndServe())
}

func (s *Server) initializeRoutes() {
	// Health check
	s.router.Get("/health", apiHandlers.HealthRequestHandler)

	// Language
	s.router.Post("/languages/{languageName}", s.handleRequestWithDBInstance(apiHandlers.ReceiveRepositoriesRequestHandler))
	s.GetWithBasicAuth("/languages/{languageName}", s.handleRequestWithDBInstance(apiHandlers.ReadRepositoriesRequestHandler))
	s.GetWithBasicAuth("/languages", s.handleRequestWithDBInstance(apiHandlers.ListLanguagesAndRepositoriesRequestHandler))
	s.GetWithBasicAuth("/stats/count-repositories", s.handleRequestWithDBInstance(apiHandlers.CountRepositoriesStarsForLanguagesRequestHandler))

	// Repositories
	s.router.Post("/repositories/{repositoryId}", s.handleRequestWithDBInstance(apiHandlers.RemoveRepositoryRequestHandler))
	s.GetWithBasicAuth("/stats/top-list", s.handleRequestWithDBInstance(apiHandlers.TopRepositoryForLanguageRequestHandler))
}

func (s *Server) GetWithBasicAuth(path string, handler http.HandlerFunc) {
	s.router.Get(path, apiHandlers.BasicAuthMiddleware(userName, password, handler))
}

// RequestHandlerFunction is the function which represents any handler
type RequestHandlerFunction func(db storage.Executor, w http.ResponseWriter, r *http.Request)

func (s *Server) handleRequestWithDBInstance(handler RequestHandlerFunction) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handler(s.db, w, r)
	}
}
