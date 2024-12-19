package serverConnection

import (
	"encoding/json"
	"errors"
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

	err = CheckRequest(w, request)
	if err != nil {
		SendResponse(w, Response{Status: "Fail", Message: err.Error()})
	}
}

func CheckRequest(w http.ResponseWriter, request Request) error {
	switch request.Message {
	case "Hello!":
		SendResponse(w, Response{Status: "success", Message: "Data successfully received"})
		return nil
	default:
		return errors.New("invalid JSON message: " + request.Message)
	}

}
