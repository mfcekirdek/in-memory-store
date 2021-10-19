package utils

import (
	"encoding/json"
	"log"
	"net/http"
)

func HandleError(w http.ResponseWriter, r *http.Request, status int) {
	w.WriteHeader(status)
	var resp []byte
	if status == http.StatusNotFound {
		resp, _ = json.Marshal("Not found")
	} else if status == http.StatusBadRequest {
		resp, _ = json.Marshal("custom 400")
	}

	if _, err := w.Write(resp); err != nil {
		log.Printf("Err: %s\n", err.Error())
		return
	}
}

func ReturnJSONResponse(w http.ResponseWriter, r *http.Request, result interface{}) {
	resp, err := json.Marshal(result)
	if err != nil {
		log.Printf("Error happened in JSON marshal. Err: %s\n", err.Error())
		HandleError(w, r, http.StatusInternalServerError)
		return
	}
	if _, err := w.Write(resp); err != nil {
		log.Printf("Err: %s\n", err.Error())
		return
	}
}
