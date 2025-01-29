package auth

import (
	_ "errors"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Секретный ключ для подписи JWT (можно из .env или оставить напрямую)
var jwtSecret = []byte(os.Getenv("zsYZS/6Rgnlc2G/S1aBLvTFSzfkWyYCW8Shudrmxgtc="))

// Если не используешь .env, то просто оставь:
var jwtSecretFallback = []byte("zsYZS/6Rgnlc2G/S1aBLvTFSzfkWyYCW8Shudrmxgtc=")

// Структура Claims (данные внутри токена)
type Claims struct {
	Email string `json:"email"`
	Role  string `json:"role"`
	jwt.RegisteredClaims
}

// Генерация JWT-токена
func GenerateJWT(email, role string) (string, error) {
	claims := Claims{
		Email: email,
		Role:  role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // Токен на 24 часа
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := jwtSecret
	if len(secret) == 0 {
		secret = jwtSecretFallback
	}
	return token.SignedString(secret)
}

// Middleware для проверки JWT
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Получаем заголовок Authorization
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing token", http.StatusUnauthorized)
			return
		}

		// Проверяем формат заголовка
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			http.Error(w, "Invalid token format", http.StatusUnauthorized)
			return
		}

		// Разбираем токен
		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		// Проверяем валидность токена
		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Если токен валиден, передаем управление обработчику
		next(w, r)
	}
}
