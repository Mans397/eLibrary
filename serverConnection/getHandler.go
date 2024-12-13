package serverConnection

import (
	"net/http"
)

func GetHandlerDataJson(w http.ResponseWriter) {
	SendResponse(w, RequestHistory)
}
