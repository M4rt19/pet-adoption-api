package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"pet-adoption-api/internal/middleware"
	"testing"

	"pet-adoption-api/internal/database"
	"pet-adoption-api/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) {
	t.Helper()

	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}

	if err := db.AutoMigrate(&models.User{}, &models.Shelter{}, &models.Pet{}, &models.AdoptionRequest{}); err != nil {
		t.Fatalf("failed to migrate test db: %v", err)
	}

	database.DB = db
}

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	// init JWT + middleware for tests
	InitAuth()
	middleware.InitAuthMiddleware() // if needed for other tests, safe to call

	auth := r.Group("/auth")
	{
		auth.POST("/register", Register)
		auth.POST("/login", Login)
	}

	return r
}

func TestRegisterAndLogin_Success(t *testing.T) {
	// env for JWT secret (if your InitAuth reads it)
	_ = os.Setenv("JWT_SECRET", "test_secret_key")

	setupTestDB(t)
	r := setupTestRouter()

	// --------- Register ---------
	regBody := map[string]interface{}{
		"name":     "Test User",
		"email":    "test@example.com",
		"password": "secret123",
		"role":     "user",
	}
	regJSON, _ := json.Marshal(regBody)

	req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewReader(regJSON))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated && w.Code != http.StatusOK {
		t.Fatalf("expected status 200/201 on register, got %d, body=%s", w.Code, w.Body.String())
	}

	// --------- Login ---------
	loginBody := map[string]interface{}{
		"email":    "test@example.com",
		"password": "secret123",
	}
	loginJSON, _ := json.Marshal(loginBody)

	req2 := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader(loginJSON))
	req2.Header.Set("Content-Type", "application/json")
	w2 := httptest.NewRecorder()

	r.ServeHTTP(w2, req2)

	if w2.Code != http.StatusOK {
		t.Fatalf("expected status 200 on login, got %d, body=%s", w2.Code, w2.Body.String())
	}

	var resp map[string]interface{}
	if err := json.Unmarshal(w2.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse login response: %v", err)
	}

	token, ok := resp["token"].(string)
	if !ok || token == "" {
		t.Fatalf("expected non-empty token in login response, got: %v", resp)
	}
}
