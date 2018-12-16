package user

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"net/http"
)

type userController struct {
	service *userService
}

func (uc *userController) getAllUsersHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		users, err := uc.service.find(); if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		if err = json.NewEncoder(w).Encode(users); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func (uc *userController) createUserHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		user, err := uc.service.create(decoder); if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		if err = json.NewEncoder(w).Encode(user); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func (uc *userController) getUserByIdHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_id := mux.Vars(r)["id"]
		objectId, _ := primitive.ObjectIDFromHex(_id)
		user, err := uc.service.get(objectId); if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		if err = json.NewEncoder(w).Encode(user); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func (uc *userController) updateUserByIdHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_id := mux.Vars(r)["id"]
		objectId, _ := primitive.ObjectIDFromHex(_id)
		decoder := json.NewDecoder(r.Body)
		user, err := uc.service.update(objectId, decoder); if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		if err := json.NewEncoder(w).Encode(user); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func (uc *userController) deleteUserByIdHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_id := mux.Vars(r)["id"]
		objectId, _ := primitive.ObjectIDFromHex(_id)
		err := uc.service.delete(objectId); if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}