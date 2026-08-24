package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"realworld-endpoints/internal/auth"
)

func TestJWTMiddleware_ValidToken(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	secret := "test-secret"
	token, _ := auth.GenerateToken(123, "test@example.com", secret)
	req.Header.Set("Authorization", "Token "+token)

	c := e.NewContext(req, rec)

	handlerCalled := false
	var capturedUserID uint
	var capturedEmail string

	handler := func(c echo.Context) error {
		handlerCalled = true
		if uid, ok := c.Get("user_id").(uint); ok {
			capturedUserID = uid
		}
		if email, ok := c.Get("user_email").(string); ok {
			capturedEmail = email
		}
		return c.String(200, "ok")
	}

	mw := JWTMiddleware(secret)
	err := mw(handler)(c)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !handlerCalled {
		t.Errorf("expected handler to be called")
	}
	if capturedUserID != 123 {
		t.Errorf("expected userID 123, got %d", capturedUserID)
	}
	if capturedEmail != "test@example.com" {
		t.Errorf("expected email test@example.com, got %s", capturedEmail)
	}
}

func TestJWTMiddleware_MissingHeader(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mw := JWTMiddleware("secret")
	_ = mw(func(c echo.Context) error { return nil })(c)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("expected 401 status, got %d", rec.Code)
	}
}

func TestJWTMiddleware_InvalidFormat(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "InvalidPrefix token")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mw := JWTMiddleware("secret")
	_ = mw(func(c echo.Context) error { return nil })(c)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("expected 401 status, got %d", rec.Code)
	}
}

func TestJWTMiddleware_InvalidToken(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Token garbage")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mw := JWTMiddleware("secret")
	_ = mw(func(c echo.Context) error { return nil })(c)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("expected 401 status, got %d", rec.Code)
	}
}

func TestJWTMiddleware_BearerPrefix(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	secret := "test-secret"
	token, _ := auth.GenerateToken(456, "bearer@example.com", secret)
	req.Header.Set("Authorization", "Bearer "+token)

	c := e.NewContext(req, rec)

	handlerCalled := false
	handler := func(c echo.Context) error {
		handlerCalled = true
		return c.String(200, "ok")
	}

	mw := JWTMiddleware(secret)
	err := mw(handler)(c)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !handlerCalled {
		t.Errorf("expected handler to be called")
	}
}
