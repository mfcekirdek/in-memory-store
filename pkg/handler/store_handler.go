// Package Handler handles the HTTP requests.
package handler

import (
	"encoding/json"
	"gitlab.com/mfcekirdek/in-memory-store/pkg/service"
	"gitlab.com/mfcekirdek/in-memory-store/pkg/utils"
	"io"
	"net/http"
)

// StoreHandler interface has ServeHTTP and Flush functions.
type StoreHandler interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
	Flush(w http.ResponseWriter, r *http.Request)
}

// storeHandler implements StoreHandler Interface.
// Contains StoreService.
type storeHandler struct {
	service service.StoreService
}

// NewStoreHandler creates a new storeHandler instance.
func NewStoreHandler(svc service.StoreService) StoreHandler {
	handler := &storeHandler{service: svc}
	return handler
}

// ServeHTTP handles HTTP requests to "/api/v1/store/<keyParam>"
// If HTTP method is GET or PUT and parameters are valid,
// ServeHTTP responds to the request using StoreService to do GET/SET operations.
func (s storeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Path[len("/api/v1/store/"):]

	switch r.Method {
	case http.MethodGet:
		if key == "" {
			utils.HandleError(w, r, http.StatusBadRequest)
			return
		}
		if result, err := s.service.Get(key); err != nil {
			utils.HandleError(w, r, http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusOK)
			utils.ReturnJSONResponse(w, r, result, "item fetched")
		}
	case http.MethodPut:
		if key == "" {
			utils.HandleError(w, r, http.StatusBadRequest)
			return
		}
		store := map[string]string{}
		body, _ := io.ReadAll(r.Body)
		err := json.Unmarshal(body, &store)
		if err != nil {
			utils.HandleError(w, r, http.StatusBadRequest)
			return
		}

		if store["value"] == "" {
			utils.HandleError(w, r, http.StatusBadRequest)
			return
		}

		result, keyAlreadyExist := s.service.Set(key, store["value"])
		if keyAlreadyExist {
			w.WriteHeader(http.StatusOK)
			utils.ReturnJSONResponse(w, r, result, "item updated")
		} else {
			w.WriteHeader(http.StatusCreated)
			utils.ReturnJSONResponse(w, r, result, "item created")
		}
	default:
		utils.HandleError(w, r, http.StatusMethodNotAllowed)
	}
}

// Flush handles DELETE requests to "/api/v1/store"
// Responds with empty object after deleting all store data via StoreService.
func (s storeHandler) Flush(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodDelete:
		w.WriteHeader(http.StatusOK)
		result := s.service.Flush()
		utils.ReturnJSONResponse(w, r, result, "all items deleted")
		return
	default:
		utils.HandleError(w, r, http.StatusMethodNotAllowed)
	}
}
