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

	// Публічні маршрути
	r.HandleFunc("/api/register", handlers.Register).Methods("POST")
	r.HandleFunc("/api/login", handlers.Login).Methods("POST")

	// Захищені маршрути з JWT
	protected := r.PathPrefix("/api").Subrouter()
	protected.Use(middleware.JWTAuthMiddleware)
	protected.HandleFunc("/grades", handlers.CreateGrade).Methods("POST")
	protected.HandleFunc("/grades", handlers.GetGrades).Methods("GET")

	log.Fatal(http.ListenAndServe(":8080", r))
}
