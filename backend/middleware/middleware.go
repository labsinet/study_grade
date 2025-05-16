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
		log.Printf("--> JWTAuthMiddleware called for %s %s", r.Method, r.URL.Path) // Детальний лог входу

		// Отримуємо заголовок Authorization
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			log.Println("    JWT: Missing Authorization header. Returning 401.") // Лог причини помилки
			// CORS заголовки мали бути встановлені CORSMiddleware раніше.
			// http.Error автоматично встановлює Content-Type: text/plain
			http.Error(w, "Missing token", http.StatusUnauthorized)
			log.Println("<-- JWTAuthMiddleware exited (Missing header)") // Лог виходу
			return                                                       // Зупиняємо виконання, якщо заголовка немає
		}

		// Перевіряємо, чи заголовок починається з "Bearer "
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader { // Якщо TrimPrefix нічого не видалив, значить "Bearer " не було
			log.Println("    JWT: Authorization header does not start with 'Bearer '. Returning 401.") // Лог причини помилки
			http.Error(w, "Invalid token format", http.StatusUnauthorized)
			log.Println("<-- JWTAuthMiddleware exited (Invalid format)") // Лог виходу
			return
		}
		log.Println("    JWT: Extracted token string.")

		// Парсимо та перевіряємо токен
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Перевіряємо метод підпису токена (наприклад, HMAC)
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				log.Printf("    JWT: Unexpected signing method: %v. Returning error.", token.Header["alg"]) // Лог причини помилки
				return nil, http.ErrNotSupported                                                            // Або інша відповідна помилка
			}
			// Повертаємо секретний ключ для валідації
			return jwtSecret, nil
		})

		// Перевіряємо помилки парсингу або невалідність токена
		if err != nil || !token.Valid {
			log.Printf("    JWT: Token parsing or validation failed: %v. Returning 401.", err) // Лог причини помилки
			// Тут можна деталізувати помилку в логах сервера, але клієнту краще дати загальне повідомлення "Invalid token".
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			log.Println("<-- JWTAuthMiddleware exited (Invalid token)") // Лог виходу
			return                                                      // Зупиняємо виконання, якщо токен невалідний
		}
		log.Println("    JWT: Token parsed and is valid.")

		// Отримуємо claims токена
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			log.Println("    JWT: Failed to get token claims as MapClaims. Returning 401.") // Лог причини помилки
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
			log.Println("<-- JWTAuthMiddleware exited (Invalid claims)") // Лог виходу
			return
		}
		log.Println("    JWT: Claims extracted.")

		// Отримуємо user_id з claims
		// numbers are often decoded as float64 by default JSON unmarshalling
		userIDFloat, ok := claims["user_id"].(float64)
		if !ok {
			log.Println("    JWT: User ID claim is missing or not a number (float64). Returning 401.") // Лог причини помилки
			http.Error(w, "Invalid user ID claim type", http.StatusUnauthorized)
			log.Println("<-- JWTAuthMiddleware exited (Invalid UserID type)") // Лог виходу
			return
		}
		userID := int(userIDFloat)
		log.Printf("    JWT: Token validated successfully for user ID: %d", userID)

		// Додаємо user_id в контекст запиту, щоб обробники могли його отримати
		ctx := context.WithValue(r.Context(), "userID", userID)
		log.Println("    JWT: User ID added to context. Proceeding to next handler.") // Лог успішного проходження

		// Передаємо запит далі по ланцюгу обробників (до кінцевого handler)
		next.ServeHTTP(w, r.WithContext(ctx))
		log.Println("<-- JWTAuthMiddleware finished processing") // Лог нормального виходу
	})
}

// CORSMiddleware залишається без змін з попереднього кроку
func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("--> CORSMiddleware called for %s %s", r.Method, r.URL.Path) // Детальний лог входу

		origin := r.Header.Get("Origin")
		log.Printf("    CORS: Request Origin header: %s", origin)

		// Встановлення заголовків Access-Control
		if origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Credentials", "true") // Вмикаємо підтримку credentials (наприклад, Authorization)
			log.Printf("    CORS: Set Allow-Origin: %s, Allow-Credentials: true", origin)
		} else {
			// Для запитів з того ж орігіну, або якщо Origin відсутній
			log.Println("    CORS: No Origin header found (likely same-origin).")
			// Тут можна вирішити, чи встановлювати *, чи нічого.
			// Якщо credentials використовуються (Auth header), * без credentials=true не працює.
			// Найбезпечніше, якщо Origin присутній, ставити його + credentials.
			// Якщо ні, нічого не ставити або ставити * без credentials, якщо ви впевнені, що такі запити можливі і потрібні.
			// Поточна логіка (тільки якщо Origin присутній, встановлюємо конкретний Origin + Credentials) є правильною для крос-оріджин запитів з credentials.
		}

		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE") // Дозволені методи
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")     // Дозволені заголовки

		// Логування встановлених заголовків перед продовженням
		log.Printf("    CORS: Set Headers: Allow-Origin='%s', Allow-Credentials='%s', Allow-Methods='%s', Allow-Headers='%s'",
			w.Header().Get("Access-Control-Allow-Origin"),
			w.Header().Get("Access-Control-Allow-Credentials"),
			w.Header().Get("Access-Control-Allow-Methods"),
			w.Header().Get("Access-Control-Allow-Headers"))

		// Обробка preflight OPTIONS запитів
		if r.Method == "OPTIONS" {
			log.Println("    CORS: Handling OPTIONS preflight request. Sending 200 OK.")
			w.WriteHeader(http.StatusOK)
			log.Println("<-- CORSMiddleware exited (OPTIONS handled)") // Лог виходу
			return                                                     // !!! Зупиняємо виконання для OPTIONS !!!
		}

		// Передаємо запит далі для інших методів
		log.Println("    CORS: Not OPTIONS, calling next handler.")
		next.ServeHTTP(w, r)
		log.Println("<-- CORSMiddleware finished processing") // Лог нормального виходу
	})
}
