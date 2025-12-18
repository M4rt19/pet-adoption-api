package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// Note: No more setup function here. It uses the TestMain from shelter_test.go

// TestRegisterAndLogin_Success uses the testRouter and DB initialized in shelter_test.go's TestMain
func TestRegisterAndLogin_Success(t *testing.T) {
	// Use a unique email for registration to not conflict with base data
	registerData := gin.H{
		"name":     "Test User Reg",
		"email":    "register@test.com",
		"password": "password123",
	}
	body, _ := json.Marshal(registerData)

	// Test Registration
	wReg := httptest.NewRecorder()
	reqReg, _ := http.NewRequest("POST", "/auth/register", bytes.NewBuffer(body))
	reqReg.Header.Set("Content-Type", "application/json")
	testRouter.ServeHTTP(wReg, reqReg)

	assert.Equal(t, http.StatusCreated, wReg.Code)

	// Test Login with the user created above
	loginData := gin.H{
		"email":    "register@test.com",
		"password": "password123",
	}
	body, _ = json.Marshal(loginData)

	wLogin := httptest.NewRecorder()
	reqLogin, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(body))
	reqLogin.Header.Set("Content-Type", "application/json")
	testRouter.ServeHTTP(wLogin, reqLogin)

	assert.Equal(t, http.StatusOK, wLogin.Code)

	var response map[string]string
	json.Unmarshal(wLogin.Body.Bytes(), &response)
	assert.Contains(t, response, "token", "Login response should contain a token")
}

// Test that you can't register the same email twice
func TestRegister_DuplicateEmail(t *testing.T) {
	// The user owner@test.com is already created in the TestMain setup
	registerData := gin.H{
		"name":     "Another Owner",
		"email":    "owner@test.com", // This email already exists
		"password": "password123",
	}
	body, _ := json.Marshal(registerData)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/auth/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	testRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
