package handler

import (
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
	id := r.URL.Path[len("/api/v1/store/"):]

	if id == "" {
		utils.HandleError(w, r, http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		if result, err := s.service.Get("a"); err != nil {
			utils.HandleError(w, r, http.StatusNotFound)
		} else {
			utils.ReturnJSONResponse(w, r, result)
		}
	case http.MethodPut:
		result := s.service.Set("a", "x")
		utils.ReturnJSONResponse(w, r, result)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		utils.ReturnJSONResponse(w, r, "Method Not Allowed")
		return
	}
}

func (s *storeHandler) Flush(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodDelete:
		w.WriteHeader(http.StatusOK)
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		utils.ReturnJSONResponse(w, r, "Method Not Allowed")
		return
	}
}
