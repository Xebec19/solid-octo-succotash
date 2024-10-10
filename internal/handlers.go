package internal

import (
	"bytes"
	"encoding/json"
	"io"
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
	// Copy the response body to a new buffer
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, res.Body); err != nil {
		slog.Error("Failed to read response body", "error", err)
		http.Error(w, "Failed to read response", http.StatusInternalServerError)
		return
	}

	// Reset the response body with the copied data
	res.Body = io.NopCloser(&buf)

	// Set the content type to JSON
	w.Header().Set("Content-Type", "application/json")

	// Write the complete response body to the client
	if _, err := io.Copy(w, &buf); err != nil {
		slog.Error("Failed to write response to client", "error", err)
		http.Error(w, "Failed to send response", http.StatusInternalServerError)
		return
	}

}
