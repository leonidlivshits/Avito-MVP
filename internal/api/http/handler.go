package httpapi

import (
	"encoding/json"
	"net/http"
)

func WriteJSON(w http.ResponseWriter, code int, v interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		return err
	}
	return nil
}

func WriteError(w http.ResponseWriter, code int, errBody interface{}) error {
	return WriteJSON(w, code, errBody)
}
