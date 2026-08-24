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

func TestGetArticles_Success(t *testing.T) {
	e := echo.New()
	mockArticleRepo := repository.NewMockArticleRepository()
	mockUserRepo := repository.NewMockUserRepository()
	mockArticleRepo.Articles = append(mockArticleRepo.Articles, models.Article{ID: 1, Title: "Test", Slug: "test"})
	handler := NewArticleHandler(mockArticleRepo, mockUserRepo, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/articles", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if err := handler.GetArticles(c); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var res map[string]interface{}
	if err := json.Unmarshal(rec.Body.Bytes(), &res); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if _, ok := res["articles"]; !ok {
		t.Errorf("expected articles in response")
	}
	if _, ok := res["articlesCount"]; !ok {
		t.Errorf("expected articlesCount in response")
	}
}

func TestGetArticles_WithAuth(t *testing.T) {
	e := echo.New()
	mockArticleRepo := repository.NewMockArticleRepository()
	mockUserRepo := repository.NewMockUserRepository()
	mockArticleRepo.Articles = append(mockArticleRepo.Articles, models.Article{ID: 1, Title: "Test", Slug: "test"})
	handler := NewArticleHandler(mockArticleRepo, mockUserRepo, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/articles", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("user_id", uint(1))

	if err := handler.GetArticles(c); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
}

func TestGetArticleBySlug_Success(t *testing.T) {
	e := echo.New()
	mockArticleRepo := repository.NewMockArticleRepository()
	mockUserRepo := repository.NewMockUserRepository()
	mockArticleRepo.Articles = append(mockArticleRepo.Articles, models.Article{ID: 1, Title: "Test", Slug: "test-slug"})
	handler := NewArticleHandler(mockArticleRepo, mockUserRepo, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/articles/test-slug", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("slug")
	c.SetParamValues("test-slug")

	if err := handler.GetArticleBySlug(c); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
}

func TestGetArticleBySlug_NotFound(t *testing.T) {
	e := echo.New()
	mockArticleRepo := repository.NewMockArticleRepository()
	mockUserRepo := repository.NewMockUserRepository()
	handler := NewArticleHandler(mockArticleRepo, mockUserRepo, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/articles/unknown-slug", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("slug")
	c.SetParamValues("unknown-slug")

	if err := handler.GetArticleBySlug(c); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if rec.Code != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, rec.Code)
	}
}

func TestGetFeed_Success(t *testing.T) {
	e := echo.New()
	mockArticleRepo := repository.NewMockArticleRepository()
	mockUserRepo := repository.NewMockUserRepository()
	mockUserRepo.Follow(1, 2)
	mockArticleRepo.Articles = append(mockArticleRepo.Articles, models.Article{ID: 1, Title: "Test", Slug: "test", AuthorID: 2})
	handler := NewArticleHandler(mockArticleRepo, mockUserRepo, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/articles/feed", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("user_id", uint(1))

	if err := handler.GetFeed(c); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
}

func TestFavoriteArticle_Success(t *testing.T) {
	e := echo.New()
	mockArticleRepo := repository.NewMockArticleRepository()
	mockUserRepo := repository.NewMockUserRepository()
	mockArticleRepo.Articles = append(mockArticleRepo.Articles, models.Article{ID: 1, Title: "Test", Slug: "test-slug"})
	handler := NewArticleHandler(mockArticleRepo, mockUserRepo, nil)

	req := httptest.NewRequest(http.MethodPost, "/api/articles/test-slug/favorite", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("slug")
	c.SetParamValues("test-slug")
	c.Set("user_id", uint(1))

	if err := handler.FavoriteArticle(c); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var res map[string]interface{}
	if err := json.Unmarshal(rec.Body.Bytes(), &res); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	article, ok := res["article"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected 'article' in response")
	}

	if favorited, ok := article["favorited"].(bool); !ok || !favorited {
		t.Errorf("expected favorited to be true")
	}
}

func TestUnfavoriteArticle_Success(t *testing.T) {
	e := echo.New()
	mockArticleRepo := repository.NewMockArticleRepository()
	mockUserRepo := repository.NewMockUserRepository()
	mockArticleRepo.Articles = append(mockArticleRepo.Articles, models.Article{ID: 1, Title: "Test", Slug: "test-slug"})
	mockArticleRepo.Favorite(1, 1)
	handler := NewArticleHandler(mockArticleRepo, mockUserRepo, nil)

	req := httptest.NewRequest(http.MethodDelete, "/api/articles/test-slug/favorite", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("slug")
	c.SetParamValues("test-slug")
	c.Set("user_id", uint(1))

	if err := handler.UnfavoriteArticle(c); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var res map[string]interface{}
	if err := json.Unmarshal(rec.Body.Bytes(), &res); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	article, ok := res["article"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected 'article' in response")
	}

	if favorited, ok := article["favorited"].(bool); !ok || favorited {
		t.Errorf("expected favorited to be false")
	}
}
