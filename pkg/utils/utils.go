package utils

import (
	"encoding/json"
	model2 "gitlab.com/mfcekirdek/in-memory-store/pkg/model"
	"log"
	"net/http"
)

func HandleError(w http.ResponseWriter, r *http.Request, status int) {
	var response *model2.BaseResponse
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

func GenerateResponse(data interface{}, description string) *model2.BaseResponse {
	response := model2.BaseResponse{
		Data:        data,
		Description: description,
	}
	return &response
}
