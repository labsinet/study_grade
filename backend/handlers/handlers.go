package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
	"study_grade/db"
	"study_grade/models"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)

var validate = validator.New()
var jwtSecret = []byte("supersecretkey") // У продакшені використовуйте змінну середовища

func Register(w http.ResponseWriter, r *http.Request) {
	log.Println("Register handler called for", r.Method, r.URL.Path, "from", r.RemoteAddr)
	var req models.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Println("Failed to decode request body:", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	log.Println("Received register request with username:", req.Username)

	// Валідація даних
	req.Username = strings.TrimSpace(req.Username)
	req.Password = strings.TrimSpace(req.Password)
	if len(req.Username) < 3 || len(req.Username) > 50 {
		log.Println("Validation failed: Username length invalid")
		http.Error(w, "Username must be 3-50 characters", http.StatusBadRequest)
		return
	}
	if len(req.Password) < 8 {
		log.Println("Validation failed: Password too short")
		http.Error(w, "Password must be at least 8 characters", http.StatusBadRequest)
		return
	}

	// Перевірка унікальності
	var exists bool
	if err := db.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE username = ?)", req.Username).Scan(&exists); err != nil {
		log.Println("Database error during username check:", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	if exists {
		log.Println("Username already taken:", req.Username)
		http.Error(w, "Username already taken", http.StatusBadRequest)
		return
	}

	// Хешування пароля
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Failed to hash password:", err)
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	// Збереження користувача
	result, err := db.DB.Exec("INSERT INTO users (username, password) VALUES (?, ?)", req.Username, hashedPassword)
	if err != nil {
		log.Println("Failed to insert user:", err)
		http.Error(w, "Failed to register user", http.StatusInternalServerError)
		return
	}

	userID, _ := result.LastInsertId()
	user := models.User{
		ID:       int(userID),
		Username: req.Username,
	}

	w.WriteHeader(http.StatusCreated)
	log.Println("User registered successfully, ID:", userID, "Status: 201, Response:", user)
	json.NewEncoder(w).Encode(user)
}

func Login(w http.ResponseWriter, r *http.Request) {
	log.Println("Login handler called for", r.Method, r.URL.Path, "from", r.RemoteAddr)

	// Read raw body for debugging
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Failed to read login request body:", err)
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	log.Println("Raw login request body:", string(bodyBytes))

	// Recreate body for decoding
	r.Body = io.NopCloser(strings.NewReader(string(bodyBytes)))

	var input models.User
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		log.Println("Failed to decode login request body:", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	log.Println("Parsed login request:", input)

	// Валідація
	input.Username = strings.TrimSpace(input.Username)
	input.Password = strings.TrimSpace(input.Password)
	if input.Username == "" || input.Password == "" {
		log.Println("Validation failed: Username or password empty")
		http.Error(w, "Username and password required", http.StatusBadRequest)
		return
	}

	// Перевірка користувача
	var user models.User
	var hashedPassword string
	err = db.DB.QueryRow("SELECT id, username, password FROM users WHERE username = ?", input.Username).
		Scan(&user.ID, &user.Username, &hashedPassword)
	if err != nil {
		log.Println("Database error or user not found:", err)
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(input.Password)); err != nil {
		log.Println("Password mismatch for username:", input.Username)
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Генерація JWT
	token, err := generateJWT(user.ID)
	if err != nil {
		log.Println("Failed to generate JWT:", err)
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	response := models.LoginResponse{
		Token: token,
		User:  user,
	}
	log.Println("Login successful for user ID:", user.ID, "Response:", response)
	json.NewEncoder(w).Encode(response)
}

func CreateGrade(w http.ResponseWriter, r *http.Request) {
	var grade models.Grade
	if err := json.NewDecoder(r.Body).Decode(&grade); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Валідація
	var err error
	if err = validate.Struct(grade); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if grade.Date.IsZero() || grade.Semester < 1 || grade.TotalStudents < 1 ||
		grade.Grade5 < 0 || grade.Grade4 < 0 || grade.Grade3 < 0 || grade.Grade2 < 0 || grade.NotPassed < 0 {
		http.Error(w, "Invalid grade data", http.StatusBadRequest)
		return
	}
	if grade.Grade5+grade.Grade4+grade.Grade3+grade.Grade2+grade.NotPassed != grade.TotalStudents {
		http.Error(w, "Sum of grades must equal total students", http.StatusBadRequest)
		return
	}

	// Обчислення показників
	grade.AverageScore = float64(grade.Grade5*5+grade.Grade4*4+grade.Grade3*3+grade.Grade2*2) / float64(grade.TotalStudents)
	grade.SuccessRate = float64(grade.Grade5+grade.Grade4+grade.Grade3) / float64(grade.TotalStudents) * 100
	grade.QualityRate = float64(grade.Grade5+grade.Grade4) / float64(grade.TotalStudents) * 100

	// Отримання userID з JWT
	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	grade.UserID = userID

	// Збереження оцінки
	_, err = db.DB.Exec(
		"INSERT INTO grades (date, semester, subject, group_name, total_students, grade_5, grade_4, grade_3, grade_2, not_passed, average_score, success_rate, quality_rate, user_id) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		grade.Date, grade.Semester, grade.Subject, grade.Group, grade.TotalStudents, grade.Grade5, grade.Grade4, grade.Grade3, grade.Grade2, grade.NotPassed, grade.AverageScore, grade.SuccessRate, grade.QualityRate, grade.UserID,
	)
	if err != nil {
		http.Error(w, "Failed to save grade", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(grade)
}

func GetGrades(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	rows, err := db.DB.Query(
		"SELECT id, date, semester, subject, group_name, total_students, grade_5, grade_4, grade_3, grade_2, not_passed, average_score, success_rate, quality_rate, user_id FROM grades WHERE user_id = ?",
		userID,
	)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var grades []models.Grade
	for rows.Next() {
		var grade models.Grade
		if err := rows.Scan(&grade.ID, &grade.Date, &grade.Semester, &grade.Subject, &grade.Group, &grade.TotalStudents, &grade.Grade5, &grade.Grade4, &grade.Grade3, &grade.Grade2, &grade.NotPassed, &grade.AverageScore, &grade.SuccessRate, &grade.QualityRate, &grade.UserID); err != nil {
			http.Error(w, "Failed to scan grades", http.StatusInternalServerError)
			return
		}
		grades = append(grades, grade)
	}

	json.NewEncoder(w).Encode(grades)
}

func generateJWT(userID int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})
	return token.SignedString(jwtSecret)
}
