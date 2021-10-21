package utils

import (
	"encoding/json"
	"log"
	"net/http"

	"gitlab.com/mfcekirdek/in-memory-store/internals/model"
)

func HandleError(w http.ResponseWriter, r *http.Request, status int) {
	var response *model.BaseResponse
	if status == http.StatusNotFound {
		response = GenerateResponse(nil, "not found")
	} else if status == http.StatusBadRequest {
		response = GenerateResponse(nil, "bad request")
	} else if status == http.StatusMethodNotAllowed {
		response = GenerateResponse(nil, "method not allowed")
	}

	resp, err := json.Marshal(response)
	if err != nil {
		log.Printf("Err: %s\n", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(status)
	if _, err := w.Write(resp); err != nil {
		log.Printf("Err: %s\n", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func ReturnJSONResponse(w http.ResponseWriter, r *http.Request, result interface{}, description string) {
	response := GenerateResponse(result, description)
	resp, err := json.Marshal(response)
	if err != nil {
		log.Printf("Error happened in JSON marshal. Err: %s\n", err.Error())
		HandleError(w, r, http.StatusInternalServerError)
		return
	}
	if _, err := w.Write(resp); err != nil {
		log.Printf("Err: %s\n", err.Error())
		HandleError(w, r, http.StatusInternalServerError)
	}
}

func GenerateResponse(data interface{}, description string) *model.BaseResponse {
	response := model.BaseResponse{
		Data:        data,
		Description: description,
	}
	return &response
}
