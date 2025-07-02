package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"calico-go-project/database"
	"calico-go-project/handlers"
	"calico-go-project/models"
	"calico-go-project/utils"
)

var testServer *httptest.Server

func TestMain(m *testing.M) {
	// Setup test environment
	setup()

	// Run tests
	code := m.Run()

	// Cleanup
	cleanup()

	os.Exit(code)
}

func setup() {
	// Load test environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using default test values")
	}

	// Initialize test database
	if err := database.InitDB(); err != nil {
		log.Fatalf("Failed to initialize test database: %v", err)
	}

	// Setup test server
	router := handlers.SetupRoutes()
	testServer = httptest.NewServer(router)
}

func cleanup() {
	if testServer != nil {
		testServer.Close()
	}

	// Clean up test data
	if database.DB != nil {
		database.DB.Exec("DELETE FROM users WHERE email LIKE '%@test.com'")
		database.CloseDB()
	}
}

func TestUserRegistrationAndLogin(t *testing.T) {
	testEmail := "testuser@test.com"
	testPassword := "testpassword123"
	testName := "Test User"

	// Clean up any existing test data
	database.DB.Exec("DELETE FROM users WHERE email = $1", testEmail)

	t.Run("Register New User", func(t *testing.T) {
		user := models.User{
			Name:     testName,
			Email:    testEmail,
			Password: testPassword,
		}

		jsonData, err := json.Marshal(user)
		if err != nil {
			t.Fatalf("Failed to marshal user data: %v", err)
		}

		resp, err := http.Post(testServer.URL+"/api/register", "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			t.Fatalf("Failed to make register request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			t.Errorf("Expected status %d, got %d", http.StatusCreated, resp.StatusCode)
		}

		var response models.Response
		if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if !response.Success {
			t.Errorf("Expected success=true, got %v. Message: %v", response.Success, response.Message)
		}

		if response.Message != "User registered successfully" {
			t.Errorf("Expected success message, got: %v", response.Message)
		}
	})

	t.Run("Prevent Duplicate Registration", func(t *testing.T) {
		user := models.User{
			Name:     testName,
			Email:    testEmail,
			Password: testPassword,
		}

		jsonData, err := json.Marshal(user)
		if err != nil {
			t.Fatalf("Failed to marshal user data: %v", err)
		}

		resp, err := http.Post(testServer.URL+"/api/register", "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			t.Fatalf("Failed to make register request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusConflict {
			t.Errorf("Expected status %d, got %d", http.StatusConflict, resp.StatusCode)
		}

		var response models.Response
		if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if response.Success {
			t.Error("Expected success=false for duplicate registration")
		}
	})

	t.Run("Login with Valid Credentials", func(t *testing.T) {
		loginReq := models.LoginRequest{
			Email:    testEmail,
			Password: testPassword,
		}

		jsonData, err := json.Marshal(loginReq)
		if err != nil {
			t.Fatalf("Failed to marshal login data: %v", err)
		}

		resp, err := http.Post(testServer.URL+"/api/login", "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			t.Fatalf("Failed to make login request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, resp.StatusCode)
		}

		var response models.Response
		if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if !response.Success {
			t.Errorf("Expected success=true, got %v. Message: %v", response.Success, response.Message)
		}

		if response.UserID == 0 {
			t.Error("Expected UserID to be set")
		}

		if response.Name != testName {
			t.Errorf("Expected name=%s, got %s", testName, response.Name)
		}

		if response.Message != "Login successful" {
			t.Errorf("Expected login success message, got: %v", response.Message)
		}
	})

	t.Run("Login with Invalid Credentials", func(t *testing.T) {
		loginReq := models.LoginRequest{
			Email:    testEmail,
			Password: "wrongpassword",
		}

		jsonData, err := json.Marshal(loginReq)
		if err != nil {
			t.Fatalf("Failed to marshal login data: %v", err)
		}

		resp, err := http.Post(testServer.URL+"/api/login", "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			t.Fatalf("Failed to make login request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, resp.StatusCode)
		}

		var response models.Response
		if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if response.Success {
			t.Error("Expected success=false for invalid credentials")
		}
	})

	t.Run("Login with Non-existent User", func(t *testing.T) {
		loginReq := models.LoginRequest{
			Email:    "nonexistent@test.com",
			Password: "anypassword",
		}

		jsonData, err := json.Marshal(loginReq)
		if err != nil {
			t.Fatalf("Failed to marshal login data: %v", err)
		}

		resp, err := http.Post(testServer.URL+"/api/login", "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			t.Fatalf("Failed to make login request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, resp.StatusCode)
		}

		var response models.Response
		if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if response.Success {
			t.Error("Expected success=false for non-existent user")
		}
	})
}

func TestDatabaseOperations(t *testing.T) {
	testEmail := "dbtest@test.com"
	testName := "Database Test User"
	testPassword := "testpassword"

	// Clean up before test
	database.DB.Exec("DELETE FROM users WHERE email = $1", testEmail)

	t.Run("UserExists - Non-existent User", func(t *testing.T) {
		exists, err := database.UserExists(testEmail)
		if err != nil {
			t.Fatalf("UserExists failed: %v", err)
		}
		if exists {
			t.Error("Expected user to not exist")
		}
	})

	t.Run("CreateUser", func(t *testing.T) {
		hashedPassword, err := utils.HashPassword(testPassword)
		if err != nil {
			t.Fatalf("Failed to hash password: %v", err)
		}

		err = database.CreateUser(testName, testEmail, hashedPassword)
		if err != nil {
			t.Fatalf("CreateUser failed: %v", err)
		}
	})

	t.Run("UserExists - Existing User", func(t *testing.T) {
		exists, err := database.UserExists(testEmail)
		if err != nil {
			t.Fatalf("UserExists failed: %v", err)
		}
		if !exists {
			t.Error("Expected user to exist")
		}
	})

	t.Run("GetUserByEmail", func(t *testing.T) {
		id, name, hashedPassword, err := database.GetUserByEmail(testEmail)
		if err != nil {
			t.Fatalf("GetUserByEmail failed: %v", err)
		}

		if id == 0 {
			t.Error("Expected valid user ID")
		}

		if name != testName {
			t.Errorf("Expected name=%s, got %s", testName, name)
		}

		if hashedPassword == "" {
			t.Error("Expected hashed password to be returned")
		}

		// Verify password hash
		if !utils.CheckPasswordHash(testPassword, hashedPassword) {
			t.Error("Password hash verification failed")
		}
	})

	t.Run("GetUserByEmail - Non-existent User", func(t *testing.T) {
		_, _, _, err := database.GetUserByEmail("nonexistent@test.com")
		if err != sql.ErrNoRows {
			t.Errorf("Expected sql.ErrNoRows, got %v", err)
		}
	})
}

func TestPasswordHashing(t *testing.T) {
	password := "testpassword123"

	t.Run("HashPassword", func(t *testing.T) {
		hashedPassword, err := utils.HashPassword(password)
		if err != nil {
			t.Fatalf("HashPassword failed: %v", err)
		}

		if hashedPassword == "" {
			t.Error("Expected non-empty hashed password")
		}

		if hashedPassword == password {
			t.Error("Hashed password should not equal original password")
		}
	})

	t.Run("CheckPasswordHash", func(t *testing.T) {
		hashedPassword, err := utils.HashPassword(password)
		if err != nil {
			t.Fatalf("HashPassword failed: %v", err)
		}

		// Test correct password
		if !utils.CheckPasswordHash(password, hashedPassword) {
			t.Error("Password verification should succeed")
		}

		// Test incorrect password
		if utils.CheckPasswordHash("wrongpassword", hashedPassword) {
			t.Error("Password verification should fail for wrong password")
		}
	})
}

func TestHTTPMethodValidation(t *testing.T) {
	t.Run("Register - Invalid Method", func(t *testing.T) {
		req, err := http.NewRequest("GET", testServer.URL+"/api/register", nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusMethodNotAllowed {
			t.Errorf("Expected status %d, got %d", http.StatusMethodNotAllowed, resp.StatusCode)
		}
	})

	t.Run("Login - Invalid Method", func(t *testing.T) {
		req, err := http.NewRequest("GET", testServer.URL+"/api/login", nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusMethodNotAllowed {
			t.Errorf("Expected status %d, got %d", http.StatusMethodNotAllowed, resp.StatusCode)
		}
	})
}

func TestInvalidJSONHandling(t *testing.T) {
	t.Run("Register - Invalid JSON", func(t *testing.T) {
		invalidJSON := `{"name": "Test", "email": "test@test.com", "password": }`

		resp, err := http.Post(testServer.URL+"/api/register", "application/json", bytes.NewBufferString(invalidJSON))
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", http.StatusBadRequest, resp.StatusCode)
		}
	})

	t.Run("Login - Invalid JSON", func(t *testing.T) {
		invalidJSON := `{"email": "test@test.com", "password": }`

		resp, err := http.Post(testServer.URL+"/api/login", "application/json", bytes.NewBufferString(invalidJSON))
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", http.StatusBadRequest, resp.StatusCode)
		}
	})
}

// Benchmark tests
func BenchmarkUserRegistration(b *testing.B) {
	// Setup
	setup()
	defer cleanup()

	user := models.User{
		Name:     "Benchmark User",
		Email:    "benchmark@test.com",
		Password: "benchmarkpassword",
	}

	jsonData, _ := json.Marshal(user)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Clean up before each iteration
		database.DB.Exec("DELETE FROM users WHERE email = $1", fmt.Sprintf("benchmark%d@test.com", i))

		user.Email = fmt.Sprintf("benchmark%d@test.com", i)
		jsonData, _ = json.Marshal(user)

		resp, err := http.Post(testServer.URL+"/api/register", "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			b.Fatalf("Failed to make request: %v", err)
		}
		resp.Body.Close()
	}
}

func BenchmarkPasswordHashing(b *testing.B) {
	password := "benchmarkpassword123"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := utils.HashPassword(password)
		if err != nil {
			b.Fatalf("HashPassword failed: %v", err)
		}
	}
}
