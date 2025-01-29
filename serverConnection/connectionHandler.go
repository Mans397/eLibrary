package serverConnection

import (
	"encoding/json"
	"fmt"
	"github.com/Mans397/eLibrary/Database"
	es "github.com/Mans397/eLibrary/emailSender"
	"golang.org/x/crypto/nacl/auth"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
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
	http.HandleFunc("/books", BooksHandler)

	// Новые маршруты для регистрации и аутентификации
	http.HandleFunc("/register", RegisterHandler)
	http.HandleFunc("/confirm", ConfirmEmailHandler)
	http.HandleFunc("/login", LoginHandler)

	// Защищенные маршруты (JWT Middleware)
	http.HandleFunc("/books", auth.AuthMiddleware(BooksHandler))
	http.HandleFunc("/admin/sendEmail", auth.AuthMiddleware(SendEmailHandler))

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
			log.Println("Read User Name:", name)
			errName := user.ReadUserName(name)
			if errName != nil {
				SendResponse(w, Response{Status: "fail", Message: "Error: " + errName.Error()})
				return
			}
			SendResponse(w, user)
		} else if email != "" {
			log.Println("Read User Email:", email)
			err := user.ReadUserEmail(email)
			if err != nil {
				SendResponse(w, Response{Status: "fail", Message: "Error: " + err.Error()})
				return
			}
			SendResponse(w, user)
		} else {
			log.Println("Read all users")
			users := make([]Database.User, 10)
			var err error
			users, err = Database.ReadUserAll()
			fmt.Println(users)
			if err != nil {
				SendResponse(w, Response{Status: "fail", Message: "Error: " + err.Error()})
			}
			SendResponse(w, users)
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

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "Failed to parse form: "+err.Error(), http.StatusBadRequest)
		log.Println(err)
		return
	}

	message := r.FormValue("message")
	if message == "" {
		http.Error(w, "Message is required", http.StatusBadRequest)
		return
	}

	file, _, err := r.FormFile("attachment")
	if err != nil && err.Error() != "http: no such file" {
		http.Error(w, "Failed to read image: "+err.Error(), http.StatusBadRequest)
		log.Println(err)
		return
	}

	var imagePath string
	if file != nil {
		defer file.Close()

		tempFile, err := os.CreateTemp("", "upload-*.jpg")
		if err != nil {
			http.Error(w, "Failed to save image: "+err.Error(), http.StatusInternalServerError)
			return
		}
		defer tempFile.Close()

		_, err = io.Copy(tempFile, file)
		if err != nil {
			http.Error(w, "Failed to save image: "+err.Error(), http.StatusInternalServerError)
			return
		}

		imagePath = tempFile.Name()
	}

	err = es.SendEmailAll(&message, imagePath)
	if err != nil {
		SendResponse(w, Response{Status: "Fail", Message: err.Error()})
		log.Println(err)
		return
	}

	SendResponse(w, Response{Status: "Success", Message: "Emails sent successfully"})
}

// 📌 Регистрация пользователя
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user := Database.User{Name: req.Name, Email: req.Email}
	err := Database.CreateUser(user, req.Password)
	if err != nil {
		http.Error(w, "Error creating user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Registration successful. Check your email for confirmation code."})
}

// 📌 Подтверждение email
func ConfirmEmailHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Email string `json:"email"`
		Code  string `json:"code"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := Database.ConfirmUser(req.Email, req.Code); err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Email confirmed successfully. You can now log in."})
}

// 📌 Логин с JWT
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user, err := Database.AuthenticateUser(req.Email, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	token, err := auth.GenerateJWT(user.Email, user.Role)
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func BooksHandler(w http.ResponseWriter, r *http.Request) {

	filter := r.URL.Query().Get("filter")
	sort := r.URL.Query().Get("sort")
	page := r.URL.Query().Get("page")

	if page == "" {
		page = "1"
	}
	if sort == "" {
		sort = "title"
	}

	limit := 10
	offset := 0

	pageInt, err := strconv.Atoi(page)
	if err != nil || pageInt < 1 {
		pageInt = 1
	}

	offset = (pageInt - 1) * limit

	query := "SELECT * FROM books"
	if filter != "" {
		query += " WHERE title LIKE '%" + filter + "%'"
	}
	validSortFields := []string{"title", "price", "date"}

	if sort != "" {
		found := false
		for _, field := range validSortFields {
			if sort == field {
				found = true
				query += " ORDER BY " + field
				break
			}
		}

		if !found {
			query += " ORDER BY title"
		}
	}
	query += fmt.Sprintf(" LIMIT %d OFFSET %d", limit, offset)

	var books []Database.Book
	if err := Database.DB.Raw(query).Scan(&books).Error; err != nil {
		http.Error(w, "Ошибка при получении данных", http.StatusInternalServerError)
		return
	}

	countQuery := "SELECT COUNT(*) FROM books"
	if filter != "" {
		countQuery += " WHERE title LIKE '%" + filter + "%'"
	}

	var totalCount int
	if err := Database.DB.Raw(countQuery).Scan(&totalCount).Error; err != nil {
		http.Error(w, "Ошибка при подсчете количества данных", http.StatusInternalServerError)
		return
	}

	totalPages := (totalCount + limit - 1) / limit
	paginationPages := []int{}
	for i := 1; i <= totalPages; i++ {
		paginationPages = append(paginationPages, i)
	}

	data := struct {
		Books           []Database.Book
		Filter          string
		Sort            string
		Page            int
		TotalPages      int
		PaginationPages []int
	}{
		Books:           books,
		Filter:          filter,
		Sort:            sort,
		Page:            pageInt,
		TotalPages:      totalPages,
		PaginationPages: paginationPages,
	}

	tmpl, err := template.ParseFiles("FrontEnd/book.html")
	if err != nil {
		http.Error(w, "Ошибка при рендеринге шаблона", http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, data)
}
