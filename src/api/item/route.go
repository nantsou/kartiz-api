package item

import (
	"github.com/gorilla/mux"
	"github.com/mongodb/mongo-go-driver/mongo"
	"kartiz/auth"
)

func ApplyRoutes(router *mux.Router, db *mongo.Database, auth *auth.Auth) {
	c := db.Collection("item")
	h := handler{service: &service{c}}
	sub := router.PathPrefix("/items").Subrouter()
	sub.Use(auth.AuthMiddleWare)
	sub.HandleFunc("", h.find()).Methods("GET")
	sub.HandleFunc("", h.create()).Methods("POST")
	sub.HandleFunc("/{id}", h.getById()).Methods("GET")
	sub.HandleFunc("/{id}", h.updateById()).Methods("PUT")
	sub.HandleFunc("/{id}", h.deleteById()).Methods("DELETE")
}