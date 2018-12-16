package main

import (
	"encoding/json"
	"net/http"
)

func (s *server) getHealthCheckHandlerFunc() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewEncoder(w).Encode(map[string]string{"status": "alive"}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func (s *server) getNotFoundHandleFunc() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		if err := json.NewEncoder(w).Encode(map[string]string{"status": "not found"}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}