package user

import (
	"github.com/gorilla/mux"
	"github.com/mongodb/mongo-go-driver/mongo"
)

func ApplyUserRoutes(router *mux.Router, db *mongo.Database) *mux.Router {
	c := db.Collection("user")
	uc := userController{&userService{c}}
	sub := router.PathPrefix("/users").Subrouter()
	sub.HandleFunc("", uc.getAllUsersHandler()).Methods("GET")
	sub.HandleFunc("", uc.createUserHandler()).Methods("POST")
	sub.HandleFunc("/{id}", uc.getUserByIdHandler()).Methods("GET")
	sub.HandleFunc("/{id}", uc.updateUserByIdHandler()).Methods("PUT")
	sub.HandleFunc("/{id}", uc.deleteUserByIdHandler()).Methods("DELETE")
	return router
}