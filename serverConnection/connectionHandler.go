package serverConnection

import (
	"encoding/json"
	"github.com/Mans397/eLibrary/Database"
	es "github.com/Mans397/eLibrary/emailSender"
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
	http.HandleFunc("/db/readUser", ReadUserHandler)
	http.HandleFunc("/db/updateUser", UpdateUserHandler)
	http.HandleFunc("/db/deleteUser", DeleteUserHandler)
	http.HandleFunc("/admin/sendEmail", SendEmailHandler)

	log.Println("Server starting on port", port)
	log.Printf("http://localhost%s\n", port)
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
	case "/admin":
		FilePath = "./FrontEnd/admin.html"
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

func ReadUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		email := r.URL.Query().Get("email")
		name := r.URL.Query().Get("name")
		user := Database.User{}

		if name != "" {
			errName := user.ReadUserName(name)
			if errName != nil {
				SendResponse(w, Response{Status: "fail", Message: "Error: " + errName.Error()})
				return
			}
			SendResponse(w, user)
		} else if email != "" {
			err := user.ReadUserEmail(email)
			if err != nil {
				SendResponse(w, Response{Status: "fail", Message: "Error: " + err.Error()})
				return
			}
			SendResponse(w, user)
		} else {
			w.WriteHeader(http.StatusBadRequest)
			SendResponse(w, Response{Status: "Fail", Message: "Email and name is empty"})
		}

	}
}

func UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPut {
		var user Database.User
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			SendResponse(w, Response{Status: "fail", Message: "Error: " + err.Error()})
			return
		}

		err = Database.UpdateUser(user.Email, user.Name)
		if err != nil {
			SendResponse(w, Response{Status: "fail", Message: "Error: " + err.Error()})
			return
		}
		SendResponse(w, Response{Status: "success", Message: "User updated successfully"})
	}
}

func DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodDelete {
		var user Database.User
		err := json.NewDecoder(r.Body).Decode(&user)

		if user.Email == "" {
			SendResponse(w, Response{Status: "fail", Message: "Email is empty"})
			return
		}

		err = Database.DeleteUser(user.Email)
		if err != nil {
			SendResponse(w, Response{Status: "fail", Message: "Error: " + err.Error()})
			return
		}
		SendResponse(w, Response{Status: "success", Message: "User deleted successfully"})
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

func SendEmailHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Wrong type of http method", http.StatusMethodNotAllowed)
		return
	}
	log.Println("GET/post/email")
	request := Request{}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		SendResponse(w, Response{Status: "Fail", Message: "error in json decode: " + err.Error()})
		log.Println(err)
		return
	}

	err = es.SendEmailAll(request.Message)
	if err != nil {
		SendResponse(w, Response{Status: "Fail", Message: err.Error()})
		log.Println(err)
	}
	SendResponse(w, Response{Status: "Success", Message: "Emails sent successfully"})
}
