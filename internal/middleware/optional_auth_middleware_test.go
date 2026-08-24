package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"realworld-endpoints/internal/auth"
)

func TestOptionalJWT_WithValidToken(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	
	secret := "test-secret"
	token, _ := auth.GenerateToken(999, "opt@example.com", secret)
	req.Header.Set("Authorization", "Token "+token)
	
	c := e.NewContext(req, rec)
	
	handlerCalled := false
	var capturedUserID uint
	handler := func(c echo.Context) error {
		handlerCalled = true
		if uid, ok := c.Get("user_id").(uint); ok {
			capturedUserID = uid
		}
		return c.String(200, "ok")
	}
	
	mw := OptionalJWTMiddleware(secret)
	err := mw(handler)(c)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	
	if !handlerCalled {
		t.Errorf("expected handler to be called")
	}
	if capturedUserID != 999 {
		t.Errorf("expected userID 999, got %d", capturedUserID)
	}
}

func TestOptionalJWT_WithoutToken(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	
	handlerCalled := false
	hasUserID := false
	handler := func(c echo.Context) error {
		handlerCalled = true
		if c.Get("user_id") != nil {
			hasUserID = true
		}
		return c.String(200, "ok")
	}
	
	mw := OptionalJWTMiddleware("secret")
	err := mw(handler)(c)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	
	if !handlerCalled {
		t.Errorf("expected handler to be called")
	}
	if hasUserID {
		t.Errorf("expected context to NOT have user_id")
	}
}

func TestOptionalJWT_WithInvalidToken(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Token invalid")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	
	handlerCalled := false
	hasUserID := false
	handler := func(c echo.Context) error {
		handlerCalled = true
		if c.Get("user_id") != nil {
			hasUserID = true
		}
		return c.String(200, "ok")
	}
	
	mw := OptionalJWTMiddleware("secret")
	err := mw(handler)(c)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	
	if !handlerCalled {
		t.Errorf("expected handler to be called")
	}
	if hasUserID {
		t.Errorf("expected context to NOT have user_id")
	}
}
