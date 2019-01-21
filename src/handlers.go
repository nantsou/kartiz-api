package main

import (
	"encoding/json"
	"errors"
	"kartiz/utils"
	"net/http"
)

func (s *server) getHealthCheckHandlerFunc() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		output := utils.BuildOutput(nil, nil, http.StatusOK)
		_ = json.NewEncoder(w).Encode(output)
	}
}

func (s *server) getNotFoundHandleFunc() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		output := utils.BuildOutput(nil, errors.New("resource not found"), http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(output)
	}
}