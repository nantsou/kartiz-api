package user

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
    "net/http"

	"kartiz/auth"
	"kartiz/utils"
)

type handler struct {
	service *service
	auth *auth.Auth
}

func (h *handler) getAllUsersHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		users, err := h.service.find(); if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		var output map[string]interface{}
		if err != nil {
			output = utils.BuildOutput(nil, err, http.StatusBadRequest)
		} else {
			var ups []userProfile
			for _, user := range users {
				ups = append(ups, user.toProfile())
			}
			output = utils.BuildOutput(ups, nil, http.StatusOK)
		}
		if err = json.NewEncoder(w).Encode(output); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func (h *handler) createUserHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		user, err := h.service.create(decoder)
		var output map[string]interface{}
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			output = utils.BuildOutput(nil, err, http.StatusBadRequest)
		} else {
			w.WriteHeader(http.StatusCreated)
			output = utils.BuildOutput(user.toProfile(), nil, http.StatusCreated)
		}
		if err = json.NewEncoder(w).Encode(output); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func (h *handler) getUserByIdHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_id := mux.Vars(r)["id"]
		objectId, _ := primitive.ObjectIDFromHex(_id)
		user, err := h.service.get(objectId)
		var output map[string]interface{}
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			output = utils.BuildOutput(nil, err, http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusOK)
			output = utils.BuildOutput(user.toProfile(), nil, http.StatusOK)
		}
		if err = json.NewEncoder(w).Encode(output); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func (h *handler) updateUserByIdHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_id := mux.Vars(r)["id"]
		objectId, _ := primitive.ObjectIDFromHex(_id)
		decoder := json.NewDecoder(r.Body)
		user, err := h.service.update(objectId, decoder)

		var output map[string]interface{}
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			output = utils.BuildOutput(nil, err, http.StatusBadRequest)
		} else {
			w.WriteHeader(http.StatusOK)
			output = utils.BuildOutput(user.toProfile(), nil, http.StatusOK)
		}
		if err := json.NewEncoder(w).Encode(output); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func (h *handler) deleteUserByIdHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_id := mux.Vars(r)["id"]
		objectId, _ := primitive.ObjectIDFromHex(_id)
		err := h.service.delete(objectId)
		var output map[string]interface{}
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			output = utils.BuildOutput(nil, err, http.StatusBadRequest)
		} else {
			w.WriteHeader(http.StatusNoContent)
			output = utils.BuildOutput(nil, nil, http.StatusNoContent)
		}

		if err = json.NewEncoder(w).Encode(output); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func (h *handler) getLoginHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		user, err := h.service.login(decoder)
		var output map[string]interface{}
		var statusCode int
		if err != nil {
			statusCode = http.StatusUnauthorized
			output = utils.BuildOutput(nil, err, statusCode)
		} else {
		    payload := make(map[string]interface{})
		    payload["userId"] = user.Id
		    payload["isAdmin"] = user.IsAdmin
		    tokenString, err := h.auth.GetToken(payload)
		    if err != nil {
				statusCode = http.StatusUnauthorized
				output = utils.BuildOutput(nil, err, statusCode)
            } else {
                payload["token"] = tokenString
				statusCode = http.StatusOK
				output = utils.BuildOutput(payload, nil, statusCode)
            }
		}
		w.WriteHeader(statusCode)
		_ = json.NewEncoder(w).Encode(output)
	}
}

func (h *handler) getLogoutHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("token")
		var output map[string]interface{}
		var statusCode int
		err := h.auth.Invalidate(tokenString)
		statusCode = http.StatusUnauthorized
		output = utils.BuildOutput(nil, err, statusCode)
		w.WriteHeader(statusCode)
		_ = json.NewEncoder(w).Encode(output)
	}
}