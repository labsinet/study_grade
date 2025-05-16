package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var DB *sql.DB

func InitDB() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file:", err)
	}

	// Read environment variables
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	if dbUser == "" || dbPassword == "" || dbHost == "" || dbPort == "" || dbName == "" {
		log.Fatal("Missing required database environment variables")
	}

	// Construct DSN
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", dbUser, dbPassword, dbHost, dbPort, dbName)

	var err error
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Test the connection
	if err = DB.Ping(); err != nil {
		log.Fatal("Database ping failed:", err)
	}

	// Create users table
	_, err = DB.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id INT AUTO_INCREMENT PRIMARY KEY,
			username VARCHAR(50) UNIQUE NOT NULL,
			password VARCHAR(255) NOT NULL
		)
	`)
	if err != nil {
		log.Fatal("Failed to create users table:", err)
	}

	// Create grades table
	_, err = DB.Exec(`
		CREATE TABLE IF NOT EXISTS grades (
			id INT AUTO_INCREMENT PRIMARY KEY,
			date DATE NOT NULL,
			semester INT NOT NULL,
			subject VARCHAR(100) NOT NULL,
			group_name VARCHAR(50) NOT NULL,
			total_students INT NOT NULL,
			grade_5 INT NOT NULL,
			grade_4 INT NOT NULL,
			grade_3 INT NOT NULL,
			grade_2 INT NOT NULL,
			not_passed INT NOT NULL,
			average_score FLOAT NOT NULL,
			success_rate FLOAT NOT NULL,
			quality_rate FLOAT NOT NULL,
			user_id INT NOT NULL,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		)
	`)
	if err != nil {
		log.Fatal("Failed to create grades table:", err)
	}
}
