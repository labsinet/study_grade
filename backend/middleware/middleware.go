package middleware

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

var jwtSecret = []byte("supersecretkey")

func JWTAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("JWTAuthMiddleware called for", r.URL.Path)
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing token", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "Invalid claims", http.StatusUnauthorized)
			return
		}

		userID := int(claims["user_id"].(float64))
		ctx := context.WithValue(r.Context(), "userID", userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("CORSMiddleware called for", r.URL.Path, r.Method)

		// Set CORS headers for all responses
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000") // Specific origin instead of wildcard
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Max-Age", "3600") // Cache preflight for 1 hour

		log.Println("CORS headers set for origin: http://localhost:3000")

		// Handle preflight OPTIONS requests
		if r.Method == http.MethodOptions {
			log.Println("Handling OPTIONS request for", r.URL.Path)
			w.WriteHeader(http.StatusNoContent) // 204 is more appropriate for OPTIONS
			return
		}

		// Proceed to the next handler
		next.ServeHTTP(w, r)
	})
}
