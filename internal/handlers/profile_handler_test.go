package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"realworld-endpoints/internal/models"
	"realworld-endpoints/internal/repository"

	"github.com/labstack/echo/v4"
)

func TestGetProfile_Success(t *testing.T) {
	e := echo.New()
	mockUserRepo := repository.NewMockUserRepository()
	mockUserRepo.Users = append(mockUserRepo.Users, models.User{ID: 1, Username: "testuser"})
	handler := NewProfileHandler(mockUserRepo)

	req := httptest.NewRequest(http.MethodGet, "/api/profiles/testuser", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("username")
	c.SetParamValues("testuser")

	if err := handler.GetProfile(c); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var res map[string]interface{}
	if err := json.Unmarshal(rec.Body.Bytes(), &res); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if _, ok := res["profile"]; !ok {
		t.Errorf("expected profile in response")
	}
}

func TestGetProfile_NotFound(t *testing.T) {
	e := echo.New()
	mockUserRepo := repository.NewMockUserRepository()
	handler := NewProfileHandler(mockUserRepo)

	req := httptest.NewRequest(http.MethodGet, "/api/profiles/testuser", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("username")
	c.SetParamValues("testuser")

	if err := handler.GetProfile(c); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if rec.Code != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, rec.Code)
	}
}

func TestFollowUser_Success(t *testing.T) {
	e := echo.New()
	mockUserRepo := repository.NewMockUserRepository()
	mockUserRepo.Users = append(mockUserRepo.Users, models.User{ID: 2, Username: "targetuser"})
	handler := NewProfileHandler(mockUserRepo)

	req := httptest.NewRequest(http.MethodPost, "/api/profiles/targetuser/follow", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("username")
	c.SetParamValues("targetuser")
	c.Set("user_id", uint(1))

	if err := handler.FollowUser(c); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var res map[string]interface{}
	if err := json.Unmarshal(rec.Body.Bytes(), &res); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	profile, ok := res["profile"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected 'profile' in response")
	}

	if following, ok := profile["following"].(bool); !ok || !following {
		t.Errorf("expected following to be true")
	}
}

func TestUnfollowUser_Success(t *testing.T) {
	e := echo.New()
	mockUserRepo := repository.NewMockUserRepository()
	mockUserRepo.Users = append(mockUserRepo.Users, models.User{ID: 2, Username: "targetuser"})
	mockUserRepo.Follow(1, 2)
	handler := NewProfileHandler(mockUserRepo)

	req := httptest.NewRequest(http.MethodDelete, "/api/profiles/targetuser/follow", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("username")
	c.SetParamValues("targetuser")
	c.Set("user_id", uint(1))

	if err := handler.UnfollowUser(c); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var res map[string]interface{}
	if err := json.Unmarshal(rec.Body.Bytes(), &res); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	profile, ok := res["profile"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected 'profile' in response")
	}

	if following, ok := profile["following"].(bool); !ok || following {
		t.Errorf("expected following to be false")
	}
}
