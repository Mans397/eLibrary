package serverConnection

import (
	"net/http"
)

func GetHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./FrontEnd/index.html")
}