package db

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func InitDB() {
	var err error
	DB, err = sql.Open("mysql", "user:password@tcp(db:3306)/educational_db?parseTime=true")
	if err != nil {
		log.Fatal(err)
	}

	// Створення таблиць
	_, err = DB.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id INT AUTO_INCREMENT PRIMARY KEY,
			username VARCHAR(50) UNIQUE NOT NULL,
			password VARCHAR(255) NOT NULL
		);
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
		);
	`)
	if err != nil {
		log.Fatal(err)
	}
}
