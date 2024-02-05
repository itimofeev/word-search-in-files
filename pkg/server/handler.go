package server

import (
	"encoding/json"
	"net/http"
	"strings"
)

func (s *Server) handleSearch(w http.ResponseWriter, r *http.Request) {
	wordToSearch := strings.TrimSpace(r.URL.Query().Get("word"))
	if wordToSearch == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	fileNames, err := s.searcher.Search(wordToSearch)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(fileNames)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
