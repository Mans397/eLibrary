package serverConnection

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Request struct {
	Message string `json:"message"`
}

var RequestHistory []Request

func PostHandler(w http.ResponseWriter, r *http.Request) {
	request := Request{}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	RequestHistory = append(RequestHistory, request)
	_, err = fmt.Fprintln(w, RequestHistory)
	if err != nil {
		return
	}
}
