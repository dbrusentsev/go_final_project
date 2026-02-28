package api

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

// ключ для подписи JWT токена
var jwtKey = []byte("super_secret_key")

// структура для запроса авторизации
type signInRequest struct {
	Password string `json:"password"`
}

// signInHandler обрабатывает запрос на авторизацию
func signInHandler(w http.ResponseWriter, r *http.Request) {
	var req signInRequest

	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeJSON(w, map[string]string{"error": "ошибка чтения запроса"})
		return
	}

	err = json.Unmarshal(body, &req)
	if err != nil {
		writeJSON(w, map[string]string{"error": "ошибка десериализации JSON"})
		return
	}

	// получаем пароль из переменной окружения
	pass := os.Getenv("TODO_PASSWORD")
	if pass == "" {
		writeJSON(w, map[string]string{"error": "авторизация не требуется"})
		return
	}

	// проверяем пароль
	if req.Password != pass {
		writeJSON(w, map[string]string{"error": "Неверный пароль"})
		return
	}

	// создаем хэш пароля для токена
	hash := sha256.Sum256([]byte(pass))
	hashStr := hex.EncodeToString(hash[:])

	// создаем JWT токен
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"hash": hashStr,
	})

	// подписываем токен
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		writeJSON(w, map[string]string{"error": "ошибка создания токена"})
		return
	}

	writeJSON(w, map[string]string{"token": tokenString})
}

// auth - middleware для проверки аутентификации
func auth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// смотрим наличие пароля
		pass := os.Getenv("TODO_PASSWORD")
		if len(pass) > 0 {
			var tokenStr string
			// получаем куку
			cookie, err := r.Cookie("token")
			if err == nil {
				tokenStr = cookie.Value
			}

			// если куки нет, проверяем заголовок Authorization
			if tokenStr == "" {
				tokenStr = r.Header.Get("Authorization")
			}

			if tokenStr == "" {
				http.Error(w, "Authentification required", http.StatusUnauthorized)
				return
			}

			// парсим и проверяем токен
			token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
				return jwtKey, nil
			})

			if err != nil || !token.Valid {
				http.Error(w, "Authentification required", http.StatusUnauthorized)
				return
			}

			// проверяем хэш пароля
			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				http.Error(w, "Authentification required", http.StatusUnauthorized)
				return
			}

			hashFromToken, ok := claims["hash"].(string)
			if !ok {
				http.Error(w, "Authentification required", http.StatusUnauthorized)
				return
			}

			// вычисляем хэш текущего пароля
			hash := sha256.Sum256([]byte(pass))
			hashStr := hex.EncodeToString(hash[:])

			// сравниваем хэши
			if hashFromToken != hashStr {
				http.Error(w, "Authentification required", http.StatusUnauthorized)
				return
			}
		}
		next(w, r)
	})
}
