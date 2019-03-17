package user

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"

	"kartiz/auth"
	"kartiz/utils"
)

type handler struct {
	service *service
	auth *auth.Auth
}

func (h *handler) find() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var statusCode int
		var output map[string]interface{}
		queries := r.URL.Query()
		filter := utils.BuildFilter(queries)
		users, err := h.service.find(filter)

		if err != nil {
			statusCode = http.StatusBadRequest
			output = utils.BuildOutput(nil, err, statusCode)
		} else {
			var ups []userProfile
			for _, user := range users {
				ups = append(ups, user.toProfile())
			}
			statusCode = http.StatusOK
			output = utils.BuildOutput(ups, nil, statusCode)
		}
		w.WriteHeader(statusCode)
		_ = json.NewEncoder(w).Encode(output)
	}
}

func (h *handler) create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var statusCode int
		var output map[string]interface{}

		decoder := json.NewDecoder(r.Body)
		user, err := h.service.create(decoder)

		if err != nil {
			statusCode = http.StatusBadRequest

			output = utils.BuildOutput(nil, err, http.StatusBadRequest)
		} else {
			statusCode = http.StatusCreated
			output = utils.BuildOutput(user.toProfile(), nil, statusCode)
		}
		w.WriteHeader(statusCode)
		_ = json.NewEncoder(w).Encode(output)
	}
}

func (h *handler) getById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var statusCode int
		var output map[string]interface{}

		_id := mux.Vars(r)["id"]
		objectId, _ := primitive.ObjectIDFromHex(_id)
		user, err := h.service.get(objectId)

		if err != nil {
			statusCode = http.StatusNotFound
			output = utils.BuildOutput(nil, err, statusCode)
		} else {
			statusCode = http.StatusOK
			output = utils.BuildOutput(user.toProfile(), nil, statusCode)
		}

		w.WriteHeader(statusCode)
		_ = json.NewEncoder(w).Encode(output)
	}
}

func (h *handler) updateById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var statusCode int
		var output map[string]interface{}

		_id := mux.Vars(r)["id"]
		objectId, _ := primitive.ObjectIDFromHex(_id)
		decoder := json.NewDecoder(r.Body)
		user, err := h.service.update(objectId, decoder)

		if err != nil {
			statusCode = http.StatusBadRequest
			output = utils.BuildOutput(nil, err, statusCode)
		} else {
			statusCode = http.StatusOK
			output = utils.BuildOutput(user.toProfile(), nil, statusCode)
		}

		w.WriteHeader(statusCode)
		_ = json.NewEncoder(w).Encode(output)
	}
}

func (h *handler) deleteById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var statusCode int
		var output map[string]interface{}

		_id := mux.Vars(r)["id"]
		objectId, _ := primitive.ObjectIDFromHex(_id)
		err := h.service.delete(objectId)

		if err != nil {
			statusCode = http.StatusBadRequest
			output = utils.BuildOutput(nil, err, statusCode)
		} else {
			statusCode = http.StatusNoContent
			output = utils.BuildOutput(nil, nil, statusCode)
		}

		w.WriteHeader(statusCode)
		_ = json.NewEncoder(w).Encode(output)
	}
}

func (h *handler) login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var statusCode int
		var output map[string]interface{}

		decoder := json.NewDecoder(r.Body)
		user, err := h.service.login(decoder)

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

func (h *handler) logout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var output map[string]interface{}
		var statusCode int

		tokenString := r.Header.Get("token")
		err := h.auth.Invalidate(tokenString)
		statusCode = http.StatusUnauthorized
		output = utils.BuildOutput(nil, err, statusCode)
		w.WriteHeader(statusCode)
		_ = json.NewEncoder(w).Encode(output)
	}
}