// Package utils provides common functions and variables to support other packages.
package utils

import (
	"encoding/json"
	"gitlab.com/mfcekirdek/in-memory-store/pkg/model"
	"log"
	"net/http"
)

// HandleError responds to the requests in error situations according to the status code.
func HandleError(w http.ResponseWriter, r *http.Request, status int) {
	var response *model.BaseResponse
	if status == http.StatusNotFound {
		response = GenerateResponse(nil, "not found")
	} else if status == http.StatusBadRequest {
		response = GenerateResponse(nil, "bad input parameter/body")
	} else if status == http.StatusMethodNotAllowed {
		response = GenerateResponse(nil, "method not allowed")
	}

	resp, _ := json.Marshal(response)
	w.WriteHeader(status)
	_, _ = w.Write(resp)
}

// ReturnJSONResponse responds to the requests using http.ResponseWriter
func ReturnJSONResponse(w http.ResponseWriter, r *http.Request, result interface{}, description string) {
	response := GenerateResponse(result, description)
	resp, err := json.Marshal(response)
	if err != nil {
		log.Printf("Error happened in JSON marshal. Err: %s\n", err.Error())
		HandleError(w, r, http.StatusInternalServerError)
		return
	}

	_, _ = w.Write(resp)
}

// GenerateResponse creates a BaseResponse instance that wraps given data and description parameters and returns it.
func GenerateResponse(data interface{}, description string) *model.BaseResponse {
	response := model.BaseResponse{
		Data:        data,
		Description: description,
	}
	return &response
}
