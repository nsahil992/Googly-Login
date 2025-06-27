package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
)

var DB *sql.DB

// InitDB initializes database connection and runs migrations
func InitDB() error {
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		dbUser, dbPassword, dbHost, dbPort, dbName)

	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("error opening DB: %w", err)
	}

	if err := DB.Ping(); err != nil {
		return fmt.Errorf("error pinging DB: %w", err)
	}

	log.Println("✅ Connected to database")

	if err := runMigrations(DB); err != nil {
		return fmt.Errorf("migration error: %w", err)
	}

	if err := runMigrations(DB); err != nil {
		log.Println("Error running migrations:", err)
		return err
	}

	return nil
}

// runMigrations ensures required tables exist
func runMigrations(db *sql.DB) error {
	log.Println("🛠️  Running DB migrations...")

	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			name TEXT NOT NULL,
			email TEXT UNIQUE NOT NULL,
			password TEXT NOT NULL
		);

		CREATE INDEX IF NOT EXISTS idx_users_email ON users (email);
	`)
	if err != nil {
		return err
	}

	log.Println("✅ Migrations completed successfully")
	return nil
}

// CloseDB closes the connection
func CloseDB() {
	if DB != nil {
		DB.Close()
	}
}

// UserExists checks if a user with given email exists
func UserExists(email string) (bool, error) {
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)"
	err := DB.QueryRow(query, email).Scan(&exists)
	return exists, err
}

// CreateUser inserts a user into the DB
func CreateUser(name, email, hashedPassword string) error {
	query := "INSERT INTO users (name, email, password) VALUES ($1, $2, $3)"
	_, err := DB.Exec(query, name, email, hashedPassword)
	return err
}

// GetUserByEmail retrieves login info by email
func GetUserByEmail(email string) (int, string, string, error) {
	var id int
	var name, password string
	query := "SELECT id, name, password FROM users WHERE email = $1"
	err := DB.QueryRow(query, email).Scan(&id, &name, &password)
	return id, name, password, err
}
