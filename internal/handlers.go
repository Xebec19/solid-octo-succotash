package internal

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

func (s *Server) CreateIndex(w http.ResponseWriter, r *http.Request) {
	res, err := s.OpensearchAPI.CreateIndex("go-test-index2")

	if err != nil {
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	slog.Info(res.String())

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]bool{"ok": true})
}

func (s *Server) AddFakeData(w http.ResponseWriter, r *http.Request) {
	res, err := s.OpensearchAPI.AddFakeDocuments("go-test-index2", 20000)

	if err != nil {
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	slog.Info(res.String())

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func (s *Server) FetchData(w http.ResponseWriter, r *http.Request) {

	res, err := s.OpensearchAPI.SearchData("go-test-index2")

	if err != nil {
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	slog.Info(res.String())

	defer res.Body.Close()

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res.String())
}
