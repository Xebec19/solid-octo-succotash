package internal

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

func (s *Server) CreateIndex(w http.ResponseWriter, r *http.Request) {
	res, err := s.OpensearchAPI.CreateIndex("fake-data")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	slog.Info(res.String())

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]bool{"ok": true})
}

func (s *Server) AddFakeData(w http.ResponseWriter, r *http.Request) {
	res, err := s.OpensearchAPI.AddFakeDocuments("fake-data", 100)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	slog.Info(res.String())

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]bool{"ok": true})
}
