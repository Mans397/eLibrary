package serverConnection

import (
	"encoding/json"
	"log"
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
		SendResponse(w, Response{Status: "Fail", Message: "Not string type of 'message' key"})
		log.Println(err)
		return
	}

	if request.Message == "" {
		w.WriteHeader(http.StatusBadRequest)
		SendResponse(w, Response{Status: "Fail", Message: `"message" field is required`})
		return
	}

	RequestHistory = append(RequestHistory, request)

	SendResponse(w, Response{Status: "Success", Message: request.Message})

}
