package serverConnection

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"golang.org/x/time/rate"
)

// Создаем лимитер с лимитом 60 запросов в минуту
var limiter = rate.NewLimiter(60, 1)

// Функция для логирования HTTP-запросов
func logRequest(r *http.Request) {
	log.Printf("Method: %s, Path: %s, RemoteAddr: %s, Time: %s\n", r.Method, r.URL.Path, r.RemoteAddr, time.Now().Format(time.RFC3339))
}

// Middleware для ограничения запросов
func rateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !limiter.Allow() {
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			log.Printf("Too many requests from %s\n", r.RemoteAddr) // Логирование превышения лимита
			return
		}
		next.ServeHTTP(w, r)
	})
}

// Инициализация логирования в файл
func init() {
	logFile, err := os.OpenFile("server.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(logFile)
}

// Пример обработчика главной страницы
func homeHandler(w http.ResponseWriter, r *http.Request) {
	logRequest(r) // Логирование каждого запроса
	fmt.Fprintf(w, "Hello, World!")
}

// Пример обработчика страницы отправки email
func sendEmailHandler(w http.ResponseWriter, r *http.Request) {
	logRequest(r) // Логирование каждого запроса
	// Логика отправки письма (добавьте свою)
	log.Println("Sending email to users")
	fmt.Fprintf(w, "Email sent!")
}

// Главная функция для настройки сервера
func StartServer() {
	// Обработчики с ограничением запросов
	http.Handle("/sendEmail", rateLimitMiddleware(http.HandlerFunc(sendEmailHandler)))
	http.Handle("/", rateLimitMiddleware(http.HandlerFunc(homeHandler)))

	// Логируем запуск сервера
	log.Println("Server is starting on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
