package serverConnection

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func GetHandler(w http.ResponseWriter, r *http.Request) {
	SendResponse(w, Response{Status: "OK", Message: "Hello User"})
}

func SendResponse(w http.ResponseWriter, response Response) {
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		panic(err)
	}
}
