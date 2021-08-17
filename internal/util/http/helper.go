package http

import (
	errUtil "dns/internal/util/err"
	"encoding/json"
	"fmt"
	"net/http"
)

type ErrorResponse struct {
	Errors error `json:"errors"`
}

func WriteJson(w http.ResponseWriter, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(data)
}

func HandleValidationError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	e := json.NewEncoder(w).Encode(ErrorResponse{
		Errors: err,
	})
	if e != nil {
		fmt.Printf(errUtil.HttpResponseError+"\n", e)
	}
}
