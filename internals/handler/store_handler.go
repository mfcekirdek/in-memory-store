package handler

import (
	"gitlab.com/mfcekirdek/in-memory-store/internals/utils"
	"net/http"
)

type StoreHandler interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

type storeHandler struct {
}

func NewStoreHandler() StoreHandler {
	handler := &storeHandler{}
	return handler
}

func (s storeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/api/v1/store/"):]

	if id == "" {
		utils.HandleError(w, r, http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		return
	case http.MethodPut:
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		utils.ReturnJSONResponse(w, r, "Method Not Allowed")
		return
	}
}
