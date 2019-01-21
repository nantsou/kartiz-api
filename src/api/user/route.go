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
	authSub.HandleFunc("/login", h.getLoginHandler()).Methods("POST")
	authSub.Handle("/logout",auth.AuthMiddleWare(h.getLogoutHandler())).Methods("GET")
	sub := router.PathPrefix("/users").Subrouter()
	sub.Use(h.auth.AuthMiddleWare)
	sub.HandleFunc("", h.getAllUsersHandler()).Methods("GET")
	sub.HandleFunc("", h.createUserHandler()).Methods("POST")
	sub.HandleFunc("/{id}", h.getUserByIdHandler()).Methods("GET")
	sub.HandleFunc("/{id}", h.updateUserByIdHandler()).Methods("PUT")
	sub.HandleFunc("/{id}", h.deleteUserByIdHandler()).Methods("DELETE")
}