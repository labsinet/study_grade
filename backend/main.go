package main

import (
	"log"
	"net/http"
	"study_grade/db"
	"study_grade/handlers"
	"study_grade/middleware"

	"github.com/gorilla/mux"
)

func main() {
	db.InitDB()
	r := mux.NewRouter()

	// Налаштування CORS
	r.Use(middleware.CORSMiddleware)
	log.Println("CORSMiddleware підключено")

	// Публічні маршрути
	r.HandleFunc("/api/register", handlers.Register).Methods("POST")
	r.HandleFunc("/api/login", handlers.Login).Methods("POST")
	log.Println("Registered public routes: /api/register (POST), /api/login (POST)")

	// Захищені маршрути з JWT
	protected := r.PathPrefix("/api").Subrouter()
	protected.Use(middleware.JWTAuthMiddleware)
	protected.HandleFunc("/grades", handlers.CreateGrade).Methods("POST")
	protected.HandleFunc("/grades", handlers.GetGrades).Methods("GET")
	log.Println("Registered protected routes: /api/grades (POST, GET)")

	log.Println("Backend server starting on :8080...")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
