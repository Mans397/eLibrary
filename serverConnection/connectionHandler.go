package serverConnection

import (
	"encoding/json"
	"fmt"
	"github.com/Mans397/eLibrary/Database"
	es "github.com/Mans397/eLibrary/emailSender"
	"github.com/golang-jwt/jwt/v4"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

const port = ":8080"

func ConnectToServer() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	// Общие маршруты
	http.HandleFunc("/", MainHandler)
	http.HandleFunc("/auth/userLogin", UserLoginHandler)   // Логин
	http.HandleFunc("/auth/register", UserRegisterHandler) // Регистрация
	http.HandleFunc("/logout", LogoutHandler)              // Выход из системы

	// Админские маршруты (с проверкой доступа)
	http.HandleFunc("/admin", AdminMiddleware(AdminPageHandler))           // Админская страница
	http.HandleFunc("/admin/sendEmail", AdminMiddleware(SendEmailHandler)) // Страница отправки email

	// Маршрут для книг (для всех авторизованных)
	http.HandleFunc("/books", AuthMiddleware(BooksHandler))

	// База данных (используется для админов)
	http.HandleFunc("/db/createUser", AdminMiddleware(CreateUserHandler))
	http.HandleFunc("/db/readUser", AdminMiddleware(ReadUserHandler))
	http.HandleFunc("/db/updateUser", AdminMiddleware(UpdateUserHandler))
	http.HandleFunc("/db/deleteUser", AdminMiddleware(DeleteUserHandler))

	fmt.Println("Server starting on port", port)
	fmt.Printf("http://localhost%s\n", port)
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
	case "/userLogin":
		FilePath = "./FrontEnd/loginuser.html"
	case "/admin":
		FilePath = "./FrontEnd/admin.html"
	case "/register":
		FilePath = "./FrontEnd/register.html"
	default:
		FilePath = "./FrontEnd/error.html"
	}

	log.Println("Request Path:", r.URL.Path)

	http.ServeFile(w, r, FilePath)
}

var jwtSecret = []byte("wqdLXn3w8SgFWYd7kwg1ZZ21YAPYIjqe52hjz++0lvo=")

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var user Database.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := Database.CreateUser(user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": user.ID,
		"exp":    time.Now().Add(24 * time.Hour).Unix(),
	})
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	emailConfirmation := Database.EmailConfirmation{
		UserID: user.ID,
		Token:  tokenString,
		Expiry: time.Now().Add(24 * time.Hour),
	}
	if err := Database.DB.Create(&emailConfirmation).Error; err != nil {
		http.Error(w, "Failed to save email confirmation", http.StatusInternalServerError)
		return
	}

	go func() {
		message := "Please confirm your email using the link: http://localhost:8080/auth/confirmEmail?token=" + tokenString
		es.SendEmailLogin(user.Email, &message, "")
	}()

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User registered. Please confirm your email."})
}

func ConfirmEmailHandler(w http.ResponseWriter, r *http.Request) {
	tokenString := r.URL.Query().Get("token")
	if tokenString == "" {
		http.Error(w, "Token is required", http.StatusBadRequest)
		return
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil || !token.Valid {
		http.Error(w, "Invalid token", http.StatusBadRequest)
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		http.Error(w, "Invalid token claims", http.StatusBadRequest)
		return
	}

	userID := uint(claims["userID"].(float64))
	var confirmation Database.EmailConfirmation
	if err := Database.DB.Where("user_id = ? AND token = ?", userID, tokenString).First(&confirmation).Error; err != nil {
		http.Error(w, "Token not found or expired", http.StatusNotFound)
		return
	}

	Database.DB.Delete(&confirmation)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Email confirmed successfully."})
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
func UserLoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var creds struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Проверка логина админа
	if creds.Email == "admin@gmail.com" && creds.Password == "admin2005" {
		http.SetCookie(w, &http.Cookie{
			Name:  "session_token",
			Value: "admin_session",
			Path:  "/",
		})
		json.NewEncoder(w).Encode(map[string]string{
			"status":   "success",
			"role":     "admin",
			"redirect": "/admin",
		})
		return
	}

	// Проверка логина пользователя
	var user Database.User
	if err := Database.DB.Where("email = ?", creds.Email).First(&user).Error; err != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	if user.Password != creds.Password {
		http.Error(w, "Invalid password", http.StatusUnauthorized)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:  "session_token",
		Value: "user_session",
		Path:  "/",
	})
	json.NewEncoder(w).Encode(map[string]string{
		"status":   "success",
		"role":     "user",
		"redirect": "/books",
	})
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Удаляем cookie
	http.SetCookie(w, &http.Cookie{
		Name:   "session_token",
		Value:  "",
		Path:   "/",
		MaxAge: -1, // Устанавливаем отрицательный возраст, чтобы удалить cookie
	})

	// Возвращаем успешный ответ
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "success",
		"message": "Logged out successfully",
	})
}

func UserRegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var creds struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var existingUser Database.User
	if err := Database.DB.Where("email = ?", creds.Email).First(&existingUser).Error; err == nil {
		http.Error(w, "User already exists", http.StatusConflict)
		return
	}

	user := Database.User{
		Name:     creds.Name,
		Email:    creds.Email,
		Password: creds.Password,
	}

	if err := Database.DB.Create(&user).Error; err != nil {
		http.Error(w, "Failed to register user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "success",
		"message": "User registered successfully",
	})
}
func AdminPageHandler(w http.ResponseWriter, r *http.Request) {
	// Проверяем, что запрос сделан с методом GET
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Проверяем cookie, чтобы удостовериться, что пользователь — админ
	cookie, err := r.Cookie("session_token")
	if err != nil || cookie.Value != "admin_session" {
		// Если пользователь не авторизован как админ, перенаправляем на страницу логина
		http.Redirect(w, r, "/auth/userLogin", http.StatusFound)
		return
	}

	// Отображаем страницу админа
	http.ServeFile(w, r, "./FrontEnd/admin.html")
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
func AdminMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Проверяем cookie
		cookie, err := r.Cookie("session_token")
		if err != nil || cookie.Value != "admin_session" {
			// Если нет прав доступа, перенаправляем на логин
			http.Redirect(w, r, "/auth/userLogin", http.StatusFound)
			return
		}

		// Если права доступа есть, вызываем следующий обработчик
		next(w, r)
	}
}
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Проверяем cookie
		cookie, err := r.Cookie("session_token")
		if err != nil || (cookie.Value != "admin_session" && cookie.Value != "user_session") {
			// Если нет прав доступа, перенаправляем на логин
			http.Redirect(w, r, "/auth/userLogin", http.StatusFound)
			return
		}

		// Если пользователь авторизован, вызываем следующий обработчик
		next(w, r)
	}
}

func BooksHandler(w http.ResponseWriter, r *http.Request) {
	// Проверяем наличие cookie
	cookie, err := r.Cookie("session_token")
	if err != nil || cookie.Value != "user_session" {
		// Если cookie нет или она не соответствует пользователю, перенаправляем на страницу авторизации для пользователей
		http.Redirect(w, r, "/userLogin", http.StatusFound)
		return
	}

	// Логика отображения книг
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
