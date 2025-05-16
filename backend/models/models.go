package models

import (
	"time"
)

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password,omitempty"` // Не повертаємо пароль
}

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Grade struct {
	ID            int       `json:"id"`
	Date          time.Time `json:"date"`
	Semester      int       `json:"semester"`
	Subject       string    `json:"subject"`
	Group         string    `json:"group"`
	TotalStudents int       `json:"total_students"`
	Grade5        int       `json:"grade_5"`
	Grade4        int       `json:"grade_4"`
	Grade3        int       `json:"grade_3"`
	Grade2        int       `json:"grade_2"`
	NotPassed     int       `json:"not_passed"`
	AverageScore  float64   `json:"average_score"`
	SuccessRate   float64   `json:"success_rate"`
	QualityRate   float64   `json:"quality_rate"`
	UserID        int       `json:"user_id"`
}

type LoginResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}
