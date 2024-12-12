package serverConnection

import (
	"encoding/json"
	"net/http"
)

func GetHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./FrontEnd/index.html")
}

func GetHandlerJson(w http.ResponseWriter) {
	err := json.NewEncoder(w).Encode(RequestHistory)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
