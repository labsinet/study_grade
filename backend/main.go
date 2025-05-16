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
	r := mux.NewRouter().StrictSlash(true) // Handle trailing slashes

	// Налаштування CORS
	r.Use(middleware.CORSMiddleware)

	// Публічні маршрути
	r.HandleFunc("/api/register", handlers.Register).Methods("POST")
	r.HandleFunc("/api/register", handlers.Register).Methods("OPTIONS")
	r.HandleFunc("/api/register", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Invalid method for /api/register:", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}).Methods("GET")
	r.HandleFunc("/api/login", handlers.Login).Methods("POST")
	r.HandleFunc("/api/login", handlers.Login).Methods("OPTIONS")
	log.Println("Registered public routes: /api/register (POST, GET), /api/login (POST)")

	// Захищені маршрути з JWT
	protected := r.PathPrefix("/api").Subrouter()
	protected.Use(middleware.JWTAuthMiddleware)
	protected.HandleFunc("/grades", handlers.CreateGrade).Methods("POST")
	protected.HandleFunc("/grades", handlers.GetGrades).Methods("GET")
	log.Println("Registered protected routes: /api/grades (POST, GET)")

	// Catch-all for undefined routes
	r.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Received request for undefined route:", r.Method, r.URL.Path)
		// CORS headers are set by CORSMiddleware, so no need to add here
		http.Error(w, "404 page not found", http.StatusNotFound)
	})

	log.Println("Backend server starting on :8080 with StrictSlash enabled...")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
