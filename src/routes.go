package main

import (
	"kartiz/api/user"
)

func (s *server) routes() {
	// set middleware
	s.router.Use(s.setContentType)
	s.router.Use(s.accessLogging)

	// set default handler to not found with access logging
	s.router.NotFoundHandler = s.accessLogging(s.setContentType(s.getNotFoundHandleFunc()))

	// setup routes
	s.router.HandleFunc("/", s.getHealthCheckHandlerFunc()).Methods("GET")
	s.router.Handle("/needAuth", s.auth.AuthMiddleWare(s.getHealthCheckHandlerFunc())).Methods("GET")

	// api
	sub := s.router.PathPrefix("/api").Subrouter()
	user.ApplyRoutes(sub, s.db, s.auth)
}
