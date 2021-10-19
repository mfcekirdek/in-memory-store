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
	w.Write(resp)
}

func ReturnJSONResponse(w http.ResponseWriter, r *http.Request, result interface{}) {
	resp, err := json.Marshal(result)
	if err != nil {
		log.Printf("Error happened in JSON marshal. Err: %s\n", err.Error())
		HandleError(w, r, http.StatusInternalServerError)
		return
	}
	w.Write(resp)
}
