package handler

import (
	"encoding/json"
	"io"
	"net/http"

	"gitlab.com/mfcekirdek/in-memory-store/internals/service"

	"gitlab.com/mfcekirdek/in-memory-store/internals/utils"
)

type StoreHandler interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
	Flush(w http.ResponseWriter, r *http.Request)
}

type storeHandler struct {
	service service.StoreService
}

func NewStoreHandler(svc service.StoreService) StoreHandler {
	handler := &storeHandler{service: svc}
	return handler
}

func (s *storeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
			utils.ReturnJSONResponse(w, r, result)
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

		result := s.service.Set(key, store["value"])
		w.WriteHeader(http.StatusOK)
		utils.ReturnJSONResponse(w, r, result)
	default:
		utils.HandleError(w, r, http.StatusMethodNotAllowed)
	}
}

func (s *storeHandler) Flush(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodDelete:
		w.WriteHeader(http.StatusOK)
		result := s.service.Flush()
		utils.ReturnJSONResponse(w, r, result)
		return
	default:
		utils.HandleError(w, r, http.StatusMethodNotAllowed)
	}
}
