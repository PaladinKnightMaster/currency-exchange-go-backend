package main

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/paladinknightmaster/currency-exchange-go-backend/cache"
	"github.com/paladinknightmaster/currency-exchange-go-backend/models"
)

func TestLoadEnv_FileExists(t *testing.T) {
	// Save original godotenvLoadFunc and restore after test
	originalLoadFunc := godotenvLoadFunc
	defer func() { godotenvLoadFunc = originalLoadFunc }()

	// Mock godotenvLoadFunc to set environment variables
	godotenvLoadFunc = func(filenames ...string) error {
		os.Setenv("TEST_VAR", "123")
		os.Setenv("ANOTHER_VAR", "456")
		return nil
	}

	// Clear any existing environment variables that may interfere
	os.Unsetenv("TEST_VAR")
	os.Unsetenv("ANOTHER_VAR")

	// Capture log output
	var logOutput bytes.Buffer
	log.SetOutput(&logOutput)
	defer func() {
		log.SetOutput(os.Stderr) // Restore default log output
	}()

	// Call loadEnv()
	loadEnv()

	// Verify that environment variables are set
	if os.Getenv("TEST_VAR") != "123" {
		t.Errorf("Expected TEST_VAR to be '123', got '%s'", os.Getenv("TEST_VAR"))
	}
	if os.Getenv("ANOTHER_VAR") != "456" {
		t.Errorf("Expected ANOTHER_VAR to be '456', got '%s'", os.Getenv("ANOTHER_VAR"))
	}

	// Ensure no log output (since the .env file exists)
	if logOutput.Len() > 0 {
		t.Errorf("Expected no log output, got: %s", logOutput.String())
	}
}

func TestLoadEnv_FileNotExists(t *testing.T) {
	// Save original godotenvLoadFunc and restore after test
	originalLoadFunc := godotenvLoadFunc
	defer func() { godotenvLoadFunc = originalLoadFunc }()

	// Mock godotenvLoadFunc to return an error
	godotenvLoadFunc = func(filenames ...string) error {
		return fmt.Errorf("mocked error: file not found")
	}

	// Set log flags to zero to exclude timestamp
	log.SetFlags(0)
	defer func() {
		log.SetFlags(log.LstdFlags)
	}()

	// Capture log output
	var logOutput bytes.Buffer
	log.SetOutput(&logOutput)
	defer func() {
		log.SetOutput(os.Stderr)
	}()

	// Call loadEnv()
	loadEnv()

	// Verify the log message
	logMessage := logOutput.String()
	expectedMessage := "No .env file found\n"
	if logMessage != expectedMessage {
		t.Errorf("Expected log message '%s', got '%s'", expectedMessage, logMessage)
	}
}

func TestSetupRouter(t *testing.T) {
	// Clear cache before test
	cache.ClearCache()

	// Populate the cache with mock data
	mockRates := &models.ExchangeRates{
		BaseCode: "USD",
		Rates: map[string]float64{
			"USD": 1.0,
			"EUR": 0.85,
			"GBP": 0.75,
		},
	}
	cache.SaveRatesToCache(mockRates)

	router := setupRouter()

	// Define the routes to test
	routes := []struct {
		method       string
		path         string
		expectedCode int
	}{
		{"GET", "/rates", http.StatusOK},
		{"GET", "/convert?from=USD&to=EUR&amount=100", http.StatusOK},
		{"GET", "/convert", http.StatusBadRequest}, // Missing required query parameters
		{"POST", "/rates", http.StatusMethodNotAllowed},
	}

	for _, route := range routes {
		req, err := http.NewRequest(route.method, route.path, nil)
		if err != nil {
			t.Fatalf("Could not create request: %v", err)
		}
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)
		if rr.Code != route.expectedCode {
			t.Errorf("Route %s %s returned wrong status code: got %v want %v",
				route.method, route.path, rr.Code, route.expectedCode)
		}
	}
}

func TestStartServer_Success(t *testing.T) {
	// Mock serverStartFunc to prevent starting an actual server
	serverStarted := false
	serverStartFunc = func(addr string, handler http.Handler) error {
		serverStarted = true
		return nil
	}

	// Call startServer
	err := startServer()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify that serverStartFunc was called
	if !serverStarted {
		t.Errorf("Expected serverStartFunc to be called")
	}

	// Reset serverStartFunc after the test
	serverStartFunc = http.ListenAndServe
}

func TestStartServer_Error(t *testing.T) {
	// Mock serverStartFunc to simulate an error
	expectedError := errors.New("server failed to start")
	serverStartFunc = func(addr string, handler http.Handler) error {
		return expectedError
	}

	// Call startServer and capture the error
	err := startServer()
	if err == nil {
		t.Errorf("Expected error, got nil")
	} else if err.Error() != expectedError.Error() {
		t.Errorf("Expected error '%v', got '%v'", expectedError, err)
	}

	// Reset serverStartFunc after the test
	serverStartFunc = http.ListenAndServe
}
