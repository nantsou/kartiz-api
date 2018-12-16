package main

import (
	"kartiz/api/user"
)

func (s *server) routes() {
	s.router.HandleFunc("/", s.getHealthCheckHandlerFunc()).Methods("GET")

	// apply services
	sub := s.router.PathPrefix("/api").Subrouter()
	user.ApplyUserRoutes(sub, s.db)

	// set middleware
	s.router.Use(s.setContentType)
	s.router.Use(s.accessLogging)

	// set default handler to not found with access logging
	s.router.NotFoundHandler = s.accessLogging(s.setContentType(s.getNotFoundHandleFunc()))
}
