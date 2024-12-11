package serverConnection

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func SendResponse(w http.ResponseWriter, i interface{}) {
	err := json.NewEncoder(w).Encode(i)
	if err != nil {
		panic(err)
	}
}
