package tests

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
    gin.SetMode(gin.TestMode)
    r := gin.New()
    // register routes here...
    return r
}

func TestRegisterAndLogin(t *testing.T) {
    r := setupRouter()

    // Register
    body, _ := json.Marshal(map[string]string{
        "username": "tester",
        "email":    "tester@example.com",
        "password": "secret123",
    })
    req := httptest.NewRequest("POST", "/auth/register", bytes.NewReader(body))
    req.Header.Set("Content-Type", "application/json")
    w := httptest.NewRecorder()
    r.ServeHTTP(w, req)

    if w.Code != http.StatusCreated {
        t.Fatalf("expected 201, got %d", w.Code)
    }

    // Login
    // ...repeat pattern, extract token, hit /api/profile
}