package internal

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

type Server struct {
	Port          string
	OpensearchAPI *OpensearchAPI
}

func (s *Server) NewServer() *http.Server {

	handler := slog.NewJSONHandler(os.Stdout, nil)

	logger := slog.NewLogLogger(handler, slog.LevelError)

	r := mux.NewRouter()

	api := r.PathPrefix("/api").Subrouter()

	api.HandleFunc("/health", s.HealthCheck).Methods("GET")

	return &http.Server{
		Handler:  api,
		Addr:     fmt.Sprintf(":%s", s.Port),
		ErrorLog: logger,
	}
}
