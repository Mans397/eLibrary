package serverConnection

import (
	"log"
	"net/http"
)

func Connect() {
	http.HandleFunc("/", FirstHandler)

	log.Println("Server starting on port 8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}

}

func FirstHandler(w http.ResponseWriter, r *http.Request) {
	http.FileServer(http.Dir("./FrontEnd"))
	switch r.Method {
	case http.MethodGet:
		log.Println("GET")
		GetHandler(w, r)
	case http.MethodPost:
		log.Println("POST")
		PostHandler(w, r)
	default:
		http.Error(w, "Wrong type of http method", http.StatusMethodNotAllowed)
	}
}
