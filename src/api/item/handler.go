package item

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
    "net/http"

	"kartiz/utils"
)

type handler struct {
	service *service
}

func (h *handler) find() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var statusCode int
		var output map[string]interface{}
		queries := r.URL.Query()
		filter := utils.BuildFilter(queries)
		items, err := h.service.find(filter)

		if err != nil {
			statusCode = http.StatusBadRequest
			output = utils.BuildOutput(nil, err, statusCode)
		} else {
			statusCode = http.StatusOK
			output = utils.BuildOutput(items, nil, statusCode)
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
		item, err := h.service.create(decoder)

		if err != nil {
			statusCode = http.StatusBadRequest
			output = utils.BuildOutput(nil, err, statusCode)
		} else {
			statusCode = http.StatusOK
			output = utils.BuildOutput(item, nil, http.StatusCreated)
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
		item, err := h.service.get(objectId)

		if err != nil {
			statusCode = http.StatusBadRequest
			output = utils.BuildOutput(nil, err, statusCode)
		} else {
			statusCode = http.StatusOK
			output = utils.BuildOutput(item, nil, http.StatusCreated)
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
		item, err := h.service.update(objectId, decoder)

		if err != nil {
			statusCode = http.StatusBadRequest
			output = utils.BuildOutput(nil, err, statusCode)
		} else {
			statusCode = http.StatusOK
			output = utils.BuildOutput(item, nil, http.StatusCreated)
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
