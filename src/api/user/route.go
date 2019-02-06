package user

import (
	"github.com/gorilla/mux"
	"github.com/mongodb/mongo-go-driver/mongo"
	"kartiz/auth"
)

func ApplyRoutes(router *mux.Router, db *mongo.Database, auth *auth.Auth) {
	c := db.Collection("user")
	h := handler{service: &service{c}, auth: auth}
	authSub := router.PathPrefix("/auth").Subrouter()
	authSub.HandleFunc("/login", h.login()).Methods("POST")
	authSub.Handle("/logout",auth.AuthMiddleWare(h.logout())).Methods("GET")
	sub := router.PathPrefix("/users").Subrouter()
	sub.Use(h.auth.AuthMiddleWare)
	sub.HandleFunc("", h.find()).Methods("GET")
	sub.HandleFunc("", h.create()).Methods("POST")
	sub.HandleFunc("/{id}", h.getById()).Methods("GET")
	sub.HandleFunc("/{id}", h.updateById()).Methods("PUT")
	sub.HandleFunc("/{id}", h.deleteById()).Methods("DELETE")
}