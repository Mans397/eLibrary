package serverConnection

import (
	"net/http"
)

func GetHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./FrontEnd/index.html")
}

func GetHandlerJson(w http.ResponseWriter) {
	SendResponse(w, RequestHistory)
}
