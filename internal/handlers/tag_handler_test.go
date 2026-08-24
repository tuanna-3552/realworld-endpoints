package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"realworld-endpoints/internal/models"
	"realworld-endpoints/internal/repository"

	"github.com/labstack/echo/v4"
)

func TestGetTags_Success(t *testing.T) {
	e := echo.New()
	mockTagRepo := repository.NewMockTagRepository()
	mockTagRepo.Tags = append(mockTagRepo.Tags, models.Tag{ID: 1, Name: "test-tag"})
	handler := NewTagHandler(mockTagRepo, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/tags", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if err := handler.GetTags(c); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var res map[string]interface{}
	if err := json.Unmarshal(rec.Body.Bytes(), &res); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	tags, ok := res["tags"].([]interface{})
	if !ok {
		t.Fatalf("expected 'tags' array in response")
	}

	if len(tags) != 1 || tags[0] != "test-tag" {
		t.Errorf("expected [\"test-tag\"], got %v", tags)
	}
}

func TestGetTags_Empty(t *testing.T) {
	e := echo.New()
	mockTagRepo := repository.NewMockTagRepository()
	handler := NewTagHandler(mockTagRepo, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/tags", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if err := handler.GetTags(c); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var res map[string]interface{}
	if err := json.Unmarshal(rec.Body.Bytes(), &res); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	tags, ok := res["tags"].([]interface{})
	if !ok {
		t.Fatalf("expected 'tags' array in response")
	}

	if len(tags) != 0 {
		t.Errorf("expected empty tags array, got %v", tags)
	}
}

func TestGetTags_Error(t *testing.T) {
	e := echo.New()
	mockTagRepo := repository.NewMockTagRepository()
	mockTagRepo.ErrFindAll = errors.New("db error")
	handler := NewTagHandler(mockTagRepo, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/tags", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if err := handler.GetTags(c); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("expected status %d, got %d", http.StatusInternalServerError, rec.Code)
	}
}
