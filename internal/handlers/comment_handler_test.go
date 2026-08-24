package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"realworld-endpoints/internal/models"
	"realworld-endpoints/internal/repository"

	"github.com/labstack/echo/v4"
)

func TestGetComments_Success(t *testing.T) {
	e := echo.New()
	mockCommentRepo := repository.NewMockCommentRepository()
	mockArticleRepo := repository.NewMockArticleRepository()
	mockUserRepo := repository.NewMockUserRepository()
	
	mockCommentRepo.SlugToArticleID["test-slug"] = 1
	mockCommentRepo.Comments = append(mockCommentRepo.Comments, models.Comment{ID: 1, ArticleID: 1, Body: "Test"})
	
	handler := NewCommentHandler(mockCommentRepo, mockArticleRepo, mockUserRepo)

	req := httptest.NewRequest(http.MethodGet, "/api/articles/test-slug/comments", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("slug")
	c.SetParamValues("test-slug")

	if err := handler.GetComments(c); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
}

func TestCreateComment_Success(t *testing.T) {
	e := echo.New()
	mockCommentRepo := repository.NewMockCommentRepository()
	mockArticleRepo := repository.NewMockArticleRepository()
	mockUserRepo := repository.NewMockUserRepository()
	
	mockArticleRepo.Articles = append(mockArticleRepo.Articles, models.Article{ID: 1, Slug: "test-slug"})
	
	handler := NewCommentHandler(mockCommentRepo, mockArticleRepo, mockUserRepo)

	body := `{"comment":{"body":"This is a test comment"}}`
	req := httptest.NewRequest(http.MethodPost, "/api/articles/test-slug/comments", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("slug")
	c.SetParamValues("test-slug")
	c.Set("user_id", uint(1))

	if err := handler.CreateComment(c); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if rec.Code != http.StatusCreated {
		t.Errorf("expected status %d, got %d", http.StatusCreated, rec.Code)
	}
}

func TestCreateComment_EmptyBody(t *testing.T) {
	e := echo.New()
	mockCommentRepo := repository.NewMockCommentRepository()
	mockArticleRepo := repository.NewMockArticleRepository()
	mockUserRepo := repository.NewMockUserRepository()
	
	mockArticleRepo.Articles = append(mockArticleRepo.Articles, models.Article{ID: 1, Slug: "test-slug"})
	
	handler := NewCommentHandler(mockCommentRepo, mockArticleRepo, mockUserRepo)

	body := `{"comment":{"body":""}}`
	req := httptest.NewRequest(http.MethodPost, "/api/articles/test-slug/comments", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("slug")
	c.SetParamValues("test-slug")
	c.Set("user_id", uint(1))

	if err := handler.CreateComment(c); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if rec.Code != http.StatusUnprocessableEntity {
		t.Errorf("expected status %d, got %d", http.StatusUnprocessableEntity, rec.Code)
	}
}

func TestDeleteComment_Success(t *testing.T) {
	e := echo.New()
	mockCommentRepo := repository.NewMockCommentRepository()
	mockArticleRepo := repository.NewMockArticleRepository()
	mockUserRepo := repository.NewMockUserRepository()
	
	mockCommentRepo.Comments = append(mockCommentRepo.Comments, models.Comment{ID: 1, AuthorID: 1})
	
	handler := NewCommentHandler(mockCommentRepo, mockArticleRepo, mockUserRepo)

	req := httptest.NewRequest(http.MethodDelete, "/api/articles/test-slug/comments/1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("slug", "id")
	c.SetParamValues("test-slug", "1")
	c.Set("user_id", uint(1))

	if err := handler.DeleteComment(c); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
}

func TestDeleteComment_NotOwner(t *testing.T) {
	e := echo.New()
	mockCommentRepo := repository.NewMockCommentRepository()
	mockArticleRepo := repository.NewMockArticleRepository()
	mockUserRepo := repository.NewMockUserRepository()
	
	mockCommentRepo.Comments = append(mockCommentRepo.Comments, models.Comment{ID: 1, AuthorID: 2})
	
	handler := NewCommentHandler(mockCommentRepo, mockArticleRepo, mockUserRepo)

	req := httptest.NewRequest(http.MethodDelete, "/api/articles/test-slug/comments/1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("slug", "id")
	c.SetParamValues("test-slug", "1")
	c.Set("user_id", uint(1))

	if err := handler.DeleteComment(c); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if rec.Code != http.StatusForbidden {
		t.Errorf("expected status %d, got %d", http.StatusForbidden, rec.Code)
	}
}

func TestDeleteComment_NotFound(t *testing.T) {
	e := echo.New()
	mockCommentRepo := repository.NewMockCommentRepository()
	mockArticleRepo := repository.NewMockArticleRepository()
	mockUserRepo := repository.NewMockUserRepository()
	
	handler := NewCommentHandler(mockCommentRepo, mockArticleRepo, mockUserRepo)

	req := httptest.NewRequest(http.MethodDelete, "/api/articles/test-slug/comments/99", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("slug", "id")
	c.SetParamValues("test-slug", "99")
	c.Set("user_id", uint(1))

	if err := handler.DeleteComment(c); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if rec.Code != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, rec.Code)
	}
}
