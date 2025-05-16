package middleware

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

// Рекомендовано зберігати секретний ключ в змінних оточення, а не прямо в коді!
var jwtSecret = []byte("supersecretkey") // !!! ЗМІНІТЬ ЦЕЙ КЛЮЧ НА НАДІЙНІШИЙ І ЗБЕРІГАЙТЕ В БЕЗПЕЧНОМУ МІСЦІ !!!

// JWTAuthMiddleware перевіряє наявність та валідність JWT токена
func JWTAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("JWTAuthMiddleware called for", r.Method, r.URL.Path)

		// Отримуємо заголовок Authorization
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			log.Println("Missing Authorization header")
			http.Error(w, "Missing token", http.StatusUnauthorized)
			return // Зупиняємо виконання, якщо заголовка немає
		}

		// Перевіряємо, чи заголовок починається з "Bearer "
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader { // Якщо TrimPrefix нічого не видалив, значить "Bearer " не було
			log.Println("Authorization header does not start with 'Bearer '")
			http.Error(w, "Invalid token format", http.StatusUnauthorized)
			return
		}

		// Парсимо та перевіряємо токен
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Перевіряємо метод підпису токена (наприклад, HMAC)
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				log.Printf("Unexpected signing method: %v", token.Header["alg"])
				return nil, http.ErrNotSupported
			}
			// Повертаємо секретний ключ для валідації
			return jwtSecret, nil
		})

		// Перевіряємо помилки парсингу або невалідність токена
		if err != nil || !token.Valid {
			log.Printf("Token parsing or validation failed: %v", err)
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return // Зупиняємо виконання, якщо токен невалідний
		}

		// Отримуємо claims токена
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			log.Println("Failed to get token claims as MapClaims")
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
			return
		}

		// Отримуємо user_id з claims
		// numbers are often decoded as float64 by default JSON unmarshalling
		userIDFloat, ok := claims["user_id"].(float64)
		if !ok {
			log.Println("User ID claim is missing or not a number")
			http.Error(w, "Invalid user ID claim", http.StatusUnauthorized)
			return
		}
		userID := int(userIDFloat)
		log.Printf("Token validated for user ID: %d", userID)

		// Додаємо user_id в контекст запиту, щоб обробники могли його отримати
		ctx := context.WithValue(r.Context(), "userID", userID)

		// Передаємо запит далі по ланцюгу обробників (до наступного middleware або кінцевого handler)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// CORSMiddleware встановлює необхідні заголовки для дозволу крос-оріджин запитів
func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("CORSMiddleware called for", r.Method, r.URL.Path)

		// Отримуємо заголовок Origin з запиту.
		// Він буде присутній у крос-оріджин запитах, які робить браузер.
		origin := r.Header.Get("Origin")
		log.Printf("Request Origin: %s", origin)

		// *** ОСНОВНЕ ВИПРАВЛЕННЯ ДЛЯ РОБОТИ З 'Authorization' HEADER ТА CREDENTIALS ***
		// Якщо запит є крос-оріджин (тобто має заголовок Origin), ми повинні:
		// 1. Встановити Access-Control-Allow-Origin на КОНКРЕТНЕ значення Origin, а не '*'.
		// 2. Встановити Access-Control-Allow-Credentials в 'true'.
		// Це необхідно, тому що Authorization header вважається credentials,
		// і браузери вимагають Allow-Credentials: true, що несумісно з Allow-Origin: *.
		if origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Credentials", "true") // Дозволяє браузеру надсилати credentials (наприклад, Authorization)
			log.Printf("Set Allow-Origin: %s, Allow-Credentials: true", origin)
		} else {
			// Якщо заголовок Origin відсутній (наприклад, запит з того ж орігіну),
			// CORS заголовки технічно не потрібні. Ми можемо нічого не робити,
			// або встановити стандартні заголовки без credentials.
			// Для простоти розробки можна встановити *, але пам'ятайте про обмеження з credentials.
			// Оскільки у вас є роути з credentials, безпечніше завжди обробляти конкретний орігін.
			// Якщо ви ВПЕВНЕНІ, що цей роут ніколи не отримає credentials з крос-оріджина
			// БЕЗ Origin заголовка, тоді можна поставити *:
			// w.Header().Set("Access-Control-Allow-Origin", "*")
			log.Println("No Origin header found (likely same-origin request). CORS headers for cross-origin might not be strictly needed.")
		}

		// Встановлюємо дозволені методи. Додайте всі методи, які використовуються у вашому API (GET, POST, PUT, DELETE тощо)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")

		// Встановлюємо дозволені заголовки. ДОДАЙТЕ СЮДИ ВСІ НЕСТАНДАРТНІ ЗАГОЛОВКИ, які надсилає ваш фронтенд.
		// Content-Type та Authorization є поширеними.
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		log.Printf("CORS headers set for %s %s -> Allow-Origin: %s, Allow-Credentials: %s, Allow-Methods: %s, Allow-Headers: %s",
			r.URL.Path,
			r.Method,
			w.Header().Get("Access-Control-Allow-Origin"),
			w.Header().Get("Access-Control-Allow-Credentials"),
			w.Header().Get("Access-Control-Allow-Methods"),
			w.Header().Get("Access-Control-Allow-Headers"))

		// Обробка preflight OPTIONS запитів
		// Якщо метод запиту OPTIONS, це preflight. Браузер чекає 200 OK відповідь
		// з встановленими вище CORS заголовками. Після цього він надішле основний запит.
		if r.Method == "OPTIONS" {
			log.Println("Handling OPTIONS preflight request. Sending 200 OK.")
			w.WriteHeader(http.StatusOK)
			return // !!! Важливо зупинити виконання ланцюга обробників після обробки OPTIONS !!!
		}

		// Для всіх інших методів (GET, POST, PUT, DELETE тощо) просто викликаємо наступний обробник у ланцюгу.
		// CORS заголовки вже встановлені вище.
		next.ServeHTTP(w, r)
	})
}
