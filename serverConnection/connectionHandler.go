package serverConnection

import (
	"encoding/json"
	"github.com/Mans397/eLibrary/Database"
	"log"
	"net/http"
)

const port = ":8080"

func ConnectToServer() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	http.HandleFunc("/", MainHandler)
	http.HandleFunc("/data/json", DataJsonHandler)
	http.HandleFunc("/post/json", SendJsonHandler)
	http.HandleFunc("/db/createUser", CreateUserHandler)

	log.Println("Server starting on port", port)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		panic(err)
	}

}

func MainHandler(w http.ResponseWriter, r *http.Request) {
	var FilePath string

	switch r.URL.Path {
	case "/":
		FilePath = "./FrontEnd/index.html"
	case "/login":
		FilePath = "./FrontEnd/login.html"
	default:
		FilePath = "./FrontEnd/error.html"
	}

	log.Println("Request Path:", r.URL.Path)

	http.ServeFile(w, r, FilePath)
}

func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var user Database.User
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			SendResponse(w, Response{Status: "fail", Message: "Error: " + err.Error()})
			return
		}
		err = Database.CreateUser(user)
		if err != nil {
			SendResponse(w, Response{Status: "fail", Message: "Error: " + err.Error()})
			return
		}
		SendResponse(w, Response{Status: "success", Message: "User created successfully"})
	}

}

func DataJsonHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		log.Println("GET/data/json")
		GetHandlerDataJson(w)
	default:
		http.Error(w, "Wrong type of http method", http.StatusMethodNotAllowed)
	}
}

func SendJsonHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Wrong type of http method", http.StatusMethodNotAllowed)
		return
	}
	log.Println("GET/post/json")
	PostHandler(w, r)
}
