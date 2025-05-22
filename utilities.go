package main

import (
	"encoding/json"
	"log"
	"net/http"
)

// response types

type errorResponse struct {
	ErrorMessage string `json:"errorMessage"`
}

// === Utility Response Handlers ===

func respondWithError(error error, code int, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	if error != nil {
		errorString := error.Error()
		errorResponse := errorResponse{ErrorMessage: errorString}
		errorData, err := json.Marshal(errorResponse)
		if err != nil {
			log.Printf("Error occured marshaling error response: %q", err)
			return
		}

		w.Write(errorData)
		return
	}

	w.Write([]byte(`{"error":"internal error"}`))
}
