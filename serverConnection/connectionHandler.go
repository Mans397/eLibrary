package serverConnection

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Mans397/eLibrary/Database"
	"github.com/Mans397/eLibrary/chat"
	es "github.com/Mans397/eLibrary/emailSender"
	"github.com/golang-jwt/jwt/v4"
	"github.com/jung-kurt/gofpdf"
	"html/template"
	"io"
	"log"
	"math/rand"
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
	http.HandleFunc("/auth/userLogin", UserLoginHandler) // Логин
	http.HandleFunc("/auth/confirmEmail", ConfirmEmailHandler)
	http.HandleFunc("/auth/register", UserRegisterHandler) // Регистрация
	http.HandleFunc("/auth/verifyOTP", VerifyOTPHandler)
	http.HandleFunc("/logout", LogoutHandler) // Выход из системы

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

	http.HandleFunc("/bookDetail", AuthMiddleware(BookDetailHandler))
	http.HandleFunc("/processPayment", AuthMiddleware(ProcessPaymentHandler))
	http.HandleFunc("/cart", AuthMiddleware(CartHandler))
	http.HandleFunc("/addToCart", AuthMiddleware(AddItemToCartHandler))
	http.HandleFunc("/removeItem", AuthMiddleware(RemoveItemHandler))
	http.HandleFunc("/updateQuantity", AuthMiddleware(UpdateQuantityHandler))
	http.HandleFunc("/checkout", CheckoutHandler)
	http.HandleFunc("/checkoutCart", CheckoutCartHandler)

	http.HandleFunc("/delete_chat", chat.HandleDeleteChat)
	http.HandleFunc("/ws", chat.HandleConnections)
	http.HandleFunc("/chats", chat.GetActiveChats)
	go chat.HandleMessages()

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
		//замена admin на admmin. Может все сломать!
	case "/admmin":
		FilePath = "./FrontEnd/admin.html"
	case "/register":
		FilePath = "./FrontEnd/register.html"
	case "/chat":
		FilePath = "./FrontEnd/chat.html"
	case "/adminChat":
		FilePath = "./FrontEnd/adminChat.html"
	default:
		FilePath = "./FrontEnd/error.html"
	}

	log.Println("Request Path:", r.URL.Path)

	http.ServeFile(w, r, FilePath)
}

var jwtSecret = []byte("supersecretkey")

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

// 📌 Подтверждение email
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
	if !ok {
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

// 📌 Логин с отправкой OTP
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

	// Проверяем, является ли пользователь администратором
	if creds.Email == "chatgpt15292005@gmail.com" && creds.Password == "admin2005" {
		// Устанавливаем админскую сессию
		http.SetCookie(w, &http.Cookie{
			Name:     "session_token",
			Value:    "admin_session",
			Path:     "/",
			HttpOnly: true,
			Expires:  time.Now().Add(24 * time.Hour),
		})

		// Возвращаем JSON с редиректом на админку
		json.NewEncoder(w).Encode(map[string]string{
			"status":   "success",
			"message":  "Admin logged in successfully",
			"redirect": "/admin",
		})

		return
	}

	// Обычный пользователь (проверка в БД)
	var user Database.User
	if err := Database.DB.Where("email = ?", creds.Email).First(&user).Error; err != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	if user.Password != creds.Password {
		http.Error(w, "Invalid password", http.StatusUnauthorized)
		return
	}

	// Генерация OTP для обычных пользователей
	otpCode := strconv.Itoa(100000 + rand.Intn(900000)) // 6-значный код
	fmt.Println("otpCode:", otpCode)
	err := Database.CreateOTP(user.ID, otpCode, 5*time.Minute)
	if err != nil {
		http.Error(w, "Failed to generate OTP", http.StatusInternalServerError)
		return
	}

	// Отправка OTP на email
	go es.SendOTPEmail(user.Email, otpCode)

	json.NewEncoder(w).Encode(map[string]string{
		"status":  "success",
		"message": "OTP sent to your email",
	})
}

// 📌 Проверка OTP и выдача JWT-токена
func VerifyOTPHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var data struct {
		Email string `json:"email"`
		OTP   string `json:"otp"`
	}

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var user Database.User
	if err := Database.DB.Where("email = ?", data.Email).First(&user).Error; err != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	if !Database.VerifyOTP(user.ID, data.OTP) {
		http.Error(w, "Invalid or expired OTP", http.StatusUnauthorized)
		return
	}

	// Генерация JWT-токена
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": user.ID,
		"exp":    time.Now().Add(24 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	// Установка cookie с токеном
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    "user_session",
		Path:     "/",
		HttpOnly: true,
		Expires:  time.Now().Add(24 * time.Hour),
	})

	// JSON-ответ на случай использования fetch
	json.NewEncoder(w).Encode(map[string]string{
		"status": "success",
		"token":  tokenString,
	})

	// Если вы хотите перенаправлять напрямую
	http.Redirect(w, r, "/books", http.StatusSeeOther)
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

	// Генерация токена подтверждения email
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

	// Отправка email пользователю с подтверждением
	go func() {
		message := "Please confirm your email using the link: http://localhost:8080/auth/confirmEmail?token=" + tokenString
		es.SendEmailLogin(user.Email, &message, "")
	}()

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User registered. Please confirm your email."})
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

// BookDetailHandler обрабатывает /bookDetail?title=...
func BookDetailHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	title := r.URL.Query().Get("title")
	if title == "" {
		http.Error(w, "Title is required", http.StatusBadRequest)
		return
	}

	// Ищем книгу в БД (таблица books)
	var book Database.Book
	err := Database.DB.Where("title = ?", title).First(&book).Error
	if err != nil {
		log.Println("Книга не найдена:", err)
		http.Error(w, "Book not found", http.StatusNotFound)
		return
	}

	// Готовим данные для шаблона
	data := struct {
		Book Database.Book
	}{
		Book: book,
	}

	tmpl, err := template.ParseFiles("FrontEnd/singleBookDetail.html")
	if err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}

	// Рендерим шаблон
	if err := tmpl.Execute(w, data); err != nil {
		log.Println("Ошибка при выполнении шаблона:", err)
	}
}

func ProcessPaymentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Читаем поля формы
	fullName := r.FormValue("fullName")
	email := r.FormValue("email")
	cardNumber := r.FormValue("cardNumber")
	cardExp := r.FormValue("cardExp")
	cardCVV := r.FormValue("cardCVV")

	// Простейшая проверка
	if cardExp == "" || cardCVV == "" {
		http.Error(w, "Invalid card data", http.StatusBadRequest)
		return
	}

	// Допустим, пользователь — userID=1
	userID := uint(1)

	// Получаем или создаём корзину
	cart, err := GetOrCreateCartForUser(userID)
	if err != nil {
		http.Error(w, "Failed to get cart", http.StatusInternalServerError)
		return
	}

	// Создаём транзакцию
	tx := Database.Transaction{
		CartID: cart.ID,
		Status: "pending",
	}
	if err := Database.DB.Create(&tx).Error; err != nil {
		log.Println("Ошибка сохранения транзакции:", err)
		http.Error(w, "Failed to create transaction", http.StatusInternalServerError)
		return
	}

	// Имитируем успех оплаты
	tx.Status = "paid"
	Database.DB.Save(&tx)

	// Меняем статус корзины
	cart.Status = "paid"
	Database.DB.Save(&cart)

	// Генерируем PDF
	pdfPath, err := GenerateCartPDFReceipt(
		"eLibrary Project", // Название проекта/компании
		tx.ID,              // Номер транзакции
		cart,               // Передаём САМО указатель (сейчас cart — *Database.Cart)
		fullName,
		email,
		cardNumber,
	)
	if err != nil {
		log.Println("Ошибка генерации PDF:", err)
		http.Error(w, "PDF generation failed", http.StatusInternalServerError)
		return
	}

	// Отправка PDF (примерно)
	go func() {
		// Если используете свой emailSender:
		// errSend := emailSender.SendEmail(email, &("Спасибо за покупку!"), pdfPath, nil)
		errSend := sendPDFByEmail(pdfPath, email)
		if errSend != nil {
			log.Println("Не удалось отправить PDF:", errSend)
		} else {
			log.Println("PDF отправлен на:", email)
		}
		os.Remove(pdfPath) // Удаляем файл
	}()

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Оплата прошла успешно! Чек отправлен на почту."))
}

func GenerateCartPDFReceipt(
	projectName string,
	transactionID uint,
	cart *Database.Cart, // <-- указатель!
	fullName, email, cardNumber string,
) (string, error) {

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, "Fiscal Receipt")
	pdf.Ln(12)

	pdf.SetFont("Arial", "", 12)
	pdf.Cell(0, 10, fmt.Sprintf("Project: %s", projectName))
	pdf.Ln(8)
	pdf.Cell(0, 10, fmt.Sprintf("Transaction ID: %d", transactionID))
	pdf.Ln(8)
	pdf.Cell(0, 10, fmt.Sprintf("Customer: %s", fullName))
	pdf.Ln(8)
	pdf.Cell(0, 10, fmt.Sprintf("Email: %s", email))
	pdf.Ln(8)

	// Маскируем последние 4 цифры карты
	maskedCard := "****"
	if len(cardNumber) >= 4 {
		maskedCard = "****" + cardNumber[len(cardNumber)-4:]
	}
	pdf.Cell(0, 10, fmt.Sprintf("Payment method: %s", maskedCard))
	pdf.Ln(12)

	// Выводим товары из cart.Items
	var total float64
	if cart.Items != nil {
		for _, item := range cart.Items {
			line := fmt.Sprintf("%s x %d = %.2f", item.ProductName, item.Quantity, item.Price*float64(item.Quantity))
			pdf.Cell(0, 10, line)
			pdf.Ln(6)
			total += item.Price * float64(item.Quantity)
		}
	}
	pdf.Ln(8)
	pdf.Cell(0, 10, fmt.Sprintf("Total sum: %.2f", total))

	filename := fmt.Sprintf("receipt_cart_%d.pdf", transactionID)
	err := pdf.OutputFileAndClose(filename)
	if err != nil {
		return "", err
	}
	return filename, nil
}

func sendPDFByEmail(pdfPath string, toEmail string) error {
	message := "Спасибо за покупку! Во вложении ваш фискальный чек."
	return es.SendEmail(toEmail, &message, pdfPath, nil)
}

// GetOrCreateCartForUser ищет корзину пользователя со статусом "open".
// Если такой нет — создаёт новую.
func GetOrCreateCartForUser(userID uint) (*Database.Cart, error) {
	var cart Database.Cart
	err := Database.DB.Where("user_id = ? AND status = ?", userID, "open").
		Preload("Items").
		First(&cart).Error
	if err != nil {
		// Если корзина не найдена — создаём
		cart = Database.Cart{
			UserID: userID,
			Status: "open",
		}
		if createErr := Database.DB.Create(&cart).Error; createErr != nil {
			return nil, createErr
		}
		return &cart, nil
	}
	return &cart, nil
}

func CartHandler(w http.ResponseWriter, r *http.Request) {
	userID := uint(1)
	cart, err := GetOrCreateCartForUser(userID)
	if err != nil {
		http.Error(w, "Failed to get cart", http.StatusInternalServerError)
		return
	}

	// Регистрируем funcMap, чтобы использовать {{add x y}} или {{sub x y}}:
	funcMap := template.FuncMap{
		"add": func(a, b int) int { return a + b },
		"sub": func(a, b int) int { return a - b },
	}

	tmpl := template.New("cart.html").Funcs(funcMap)
	tmpl, err = tmpl.ParseFiles("FrontEnd/cart.html")
	if err != nil {
		http.Error(w, "Template parse error", http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, cart); err != nil {
		http.Error(w, "Template execute error", http.StatusInternalServerError)
	}
}

func RemoveItemHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	itemIDStr := r.URL.Query().Get("itemID")
	if itemIDStr == "" {
		http.Error(w, "itemID is required", http.StatusBadRequest)
		return
	}

	itemID, err := strconv.Atoi(itemIDStr)
	if err != nil {
		http.Error(w, "invalid itemID", http.StatusBadRequest)
		return
	}

	var cartItem Database.CartItem
	if err := Database.DB.First(&cartItem, itemID).Error; err != nil {
		http.Error(w, "item not found", http.StatusNotFound)
		return
	}

	if err := Database.DB.Delete(&cartItem).Error; err != nil {
		http.Error(w, "failed to remove item", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Товар удалён из корзины"))
}

func CheckoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	tmpl, err := template.ParseFiles("FrontEnd/checkout.html")
	if err != nil {
		log.Println("Ошибка при загрузке шаблона checkout.html:", err)
		http.Error(w, "Template parse error", http.StatusInternalServerError)
		return
	}

	// Выполняем шаблон, передавая (например) nil или структуру, если нужно
	err = tmpl.Execute(w, nil)
	if err != nil {
		log.Println("Ошибка при выполнении шаблона checkout.html:", err)
		http.Error(w, "Template execution error", http.StatusInternalServerError)
		return
	}
}

func CheckoutCartHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Допустим userID=1
	userID := uint(1)
	// 1. Получаем корзину (Cart) со статусом "open"
	cart, err := GetOrCreateCartForUser(userID)
	if err != nil {
		http.Error(w, "Failed to get/create cart", http.StatusInternalServerError)
		return
	}

	// 2. Создаём транзакцию (pending) в локальной БД
	tx := Database.Transaction{
		CartID: cart.ID,
		Status: "pending",
	}
	if err := Database.DB.Create(&tx).Error; err != nil {
		http.Error(w, "Failed to create transaction", http.StatusInternalServerError)
		return
	}

	// 3. Собираем JSON с данными корзины и клиента
	//    (какие поля именно передавать — ваш выбор, ниже пример)
	//    Собираем cartItems:
	var cartItems []map[string]interface{}
	for _, item := range cart.Items { // cart.Items = []CartItem
		cartItems = append(cartItems, map[string]interface{}{
			"id":    item.ID, // или item.ProductID
			"name":  item.ProductName,
			"price": item.Price,
		})
	}

	// Предположим, мы достаём User из таблицы пользователей
	var user Database.User
	Database.DB.First(&user, userID) // если находим по ID=1

	body := map[string]interface{}{
		// (Если хотите передавать transaction_id, добавьте:
		"transaction_id": tx.ID,
		"cartItems":      cartItems,
		"customer": map[string]interface{}{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
		},
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		http.Error(w, "Failed to marshal JSON", http.StatusInternalServerError)
		return
	}

	// 4. POST-запрос на микросервис (http://localhost:8081/payment)
	microserviceURL := "http://localhost:8081/payment"
	req, err := http.NewRequest(http.MethodPost, microserviceURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		http.Error(w, "Failed to create request to microservice", http.StatusInternalServerError)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Failed to contact microservice", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	// 5. Читаем ответ микросервиса
	var result struct {
		Success      bool   `json:"success"`
		ErrorMessage string `json:"error_message"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		http.Error(w, "Invalid response from microservice", http.StatusBadGateway)
		return
	}

	// 6. Если success=true => meняем статус транзакции на "paid" и корзину на "paid".
	//    Иначе => "declined"
	if result.Success {
		tx.Status = "paid"
		cart.Status = "paid"
		Database.DB.Save(&tx)
		Database.DB.Save(&cart)

		w.Write([]byte("Оплата прошла успешно (через микросервис)!"))
	} else {
		tx.Status = "declined"
		Database.DB.Save(&tx)
		w.WriteHeader(http.StatusPaymentRequired)
		w.Write([]byte("Оплата отклонена: " + result.ErrorMessage))
	}
}

func UpdateQuantityHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var input struct {
		ItemID   int `json:"item_id"`
		Quantity int `json:"quantity"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if input.Quantity < 1 {
		http.Error(w, "Quantity must be >= 1", http.StatusBadRequest)
		return
	}

	var cartItem Database.CartItem
	if err := Database.DB.First(&cartItem, input.ItemID).Error; err != nil {
		http.Error(w, "Item not found", http.StatusNotFound)
		return
	}

	cartItem.Quantity = input.Quantity
	if err := Database.DB.Save(&cartItem).Error; err != nil {
		http.Error(w, "Failed to update item quantity", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Количество обновлено"))
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

func AddItemToCartHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Предположим, что userID = 1 (демо), в реальном случае - берём из куки/JWT
	userID := uint(1)

	// Считываем JSON. Допустим, формат:
	// {
	//   "product_id": "42",
	//   "product_name": "Some Book",
	//   "price": 10.99,
	//   "quantity": 1
	// }
	var input struct {
		ProductID   string  `json:"product_id"`
		ProductName string  `json:"product_name"`
		Price       float64 `json:"price"`
		Quantity    int     `json:"quantity"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Или, если хотите принимать обычную форму:
	// productID := r.FormValue("product_id")
	// ...

	if input.Quantity < 1 {
		input.Quantity = 1
	}

	// Получаем (или создаём) корзину пользователя
	cart, err := GetOrCreateCartForUser(userID)
	if err != nil {
		http.Error(w, "Failed to get/create cart", http.StatusInternalServerError)
		return
	}

	// Создаём новую запись в cart_items
	item := Database.CartItem{
		CartID:      cart.ID,
		ProductID:   input.ProductID,
		ProductName: input.ProductName,
		Price:       input.Price,
		Quantity:    input.Quantity,
	}
	if err := Database.DB.Create(&item).Error; err != nil {
		log.Println("Ошибка добавления в корзину:", err)
		http.Error(w, "Failed to add item to cart", http.StatusInternalServerError)
		return
	}

	// Успешный ответ
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Товар добавлен в корзину! (CartID=%d)", cart.ID)
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
