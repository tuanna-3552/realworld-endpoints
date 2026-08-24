package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"realworld-endpoints/internal/auth"
	"realworld-endpoints/internal/models"
	"realworld-endpoints/internal/repository"

	"github.com/labstack/echo/v4"
)

func TestRegister_Success(t *testing.T) {
	e := echo.New()
	mockRepo := repository.NewMockUserRepository()
	handler := NewUserHandler(mockRepo, "secret")

	body := `{"user":{"username":"test","email":"test@test.com","password":"password"}}`
	req := httptest.NewRequest(http.MethodPost, "/api/users", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if err := handler.Register(c); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if rec.Code != http.StatusCreated {
		t.Errorf("expected status %d, got %d", http.StatusCreated, rec.Code)
	}

	var res map[string]interface{}
	if err := json.Unmarshal(rec.Body.Bytes(), &res); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	user, ok := res["user"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected 'user' in response")
	}

	if user["token"] == "" {
		t.Errorf("expected token in response")
	}
}

func TestRegister_MissingFields(t *testing.T) {
	e := echo.New()
	mockRepo := repository.NewMockUserRepository()
	handler := NewUserHandler(mockRepo, "secret")

	body := `{"user":{"email":"test@test.com","password":"password"}}`
	req := httptest.NewRequest(http.MethodPost, "/api/users", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if err := handler.Register(c); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if rec.Code != http.StatusUnprocessableEntity {
		t.Errorf("expected status %d, got %d", http.StatusUnprocessableEntity, rec.Code)
	}
}

func TestRegister_DuplicateUsername(t *testing.T) {
	e := echo.New()
	mockRepo := repository.NewMockUserRepository()
	mockRepo.Users = append(mockRepo.Users, models.User{Username: "test", Email: "other@test.com"})
	handler := NewUserHandler(mockRepo, "secret")

	body := `{"user":{"username":"test","email":"test@test.com","password":"password"}}`
	req := httptest.NewRequest(http.MethodPost, "/api/users", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if err := handler.Register(c); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if rec.Code != http.StatusUnprocessableEntity {
		t.Errorf("expected status %d, got %d", http.StatusUnprocessableEntity, rec.Code)
	}
}

func TestLogin_Success(t *testing.T) {
	e := echo.New()
	mockRepo := repository.NewMockUserRepository()
	hashed, _ := auth.HashPassword("password123")
	mockRepo.Users = append(mockRepo.Users, models.User{ID: 1, Username: "test", Email: "test@test.com", Password: hashed})
	handler := NewUserHandler(mockRepo, "secret")

	body := `{"user":{"email":"test@test.com","password":"password123"}}`
	req := httptest.NewRequest(http.MethodPost, "/api/users/login", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if err := handler.Login(c); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var res map[string]interface{}
	if err := json.Unmarshal(rec.Body.Bytes(), &res); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	user, ok := res["user"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected 'user' in response")
	}

	if user["token"] == "" {
		t.Errorf("expected token in response")
	}
}

func TestLogin_WrongPassword(t *testing.T) {
	e := echo.New()
	mockRepo := repository.NewMockUserRepository()
	hashed, _ := auth.HashPassword("password123")
	mockRepo.Users = append(mockRepo.Users, models.User{ID: 1, Username: "test", Email: "test@test.com", Password: hashed})
	handler := NewUserHandler(mockRepo, "secret")

	body := `{"user":{"email":"test@test.com","password":"wrongpassword"}}`
	req := httptest.NewRequest(http.MethodPost, "/api/users/login", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if err := handler.Login(c); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if rec.Code != http.StatusUnprocessableEntity {
		t.Errorf("expected status %d, got %d", http.StatusUnprocessableEntity, rec.Code)
	}
}

func TestLogin_UserNotFound(t *testing.T) {
	e := echo.New()
	mockRepo := repository.NewMockUserRepository()
	handler := NewUserHandler(mockRepo, "secret")

	body := `{"user":{"email":"test@test.com","password":"password123"}}`
	req := httptest.NewRequest(http.MethodPost, "/api/users/login", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if err := handler.Login(c); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if rec.Code != http.StatusUnprocessableEntity {
		t.Errorf("expected status %d, got %d", http.StatusUnprocessableEntity, rec.Code)
	}
}

func TestGetCurrentUser_Success(t *testing.T) {
	e := echo.New()
	mockRepo := repository.NewMockUserRepository()
	mockRepo.Users = append(mockRepo.Users, models.User{ID: 1, Username: "test", Email: "test@test.com"})
	handler := NewUserHandler(mockRepo, "secret")

	req := httptest.NewRequest(http.MethodGet, "/api/user", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("user_id", uint(1))

	if err := handler.GetCurrentUser(c); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
}

func TestGetUsers_Success(t *testing.T) {
	e := echo.New()
	mockRepo := repository.NewMockUserRepository()
	mockRepo.Users = append(mockRepo.Users, models.User{ID: 1, Username: "test", Email: "test@test.com"})
	handler := NewUserHandler(mockRepo, "secret")

	req := httptest.NewRequest(http.MethodGet, "/api/users", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if err := handler.GetUsers(c); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
}
