package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"pet-adoption-api/internal/database"
	"pet-adoption-api/internal/models"

	"github.com/gin-gonic/gin"
)

func setupAdoptionTestData(t *testing.T) (user models.User, shelter models.Shelter, pet models.Pet) {
	t.Helper()

	// DB should already be initialized by setupTestDB from auth_test.go
	if database.DB == nil {
		setupTestDB(t)
	}

	user = models.User{
		Name:      "Requester",
		Email:     "req@example.com",
		Role:      "user",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := database.DB.Create(&user).Error; err != nil {
		t.Fatalf("failed to create test user: %v", err)
	}

	shelter = models.Shelter{
		Name:        "Test Shelter",
		Address:     "Somewhere",
		Phone:       "12345",
		OwnerUserID: user.ID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	if err := database.DB.Create(&shelter).Error; err != nil {
		t.Fatalf("failed to create test shelter: %v", err)
	}

	pet = models.Pet{
		ShelterID:   shelter.ID,
		Name:        "Buddy",
		Species:     "dog",
		Breed:       "test breed",
		Age:         2,
		Description: "friendly dog",
		Status:      models.PetStatusAvailable,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	if err := database.DB.Create(&pet).Error; err != nil {
		t.Fatalf("failed to create test pet: %v", err)
	}

	return
}

func TestApplyForAdoption_CreatesRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// use same helper as in auth_test.go
	setupTestDB(t)

	user, _, pet := setupAdoptionTestData(t)

	// real Gin router so route params work
	r := gin.Default()

	// minimal route that injects userID into context and calls handler
	r.POST("/adoptions/:petID/apply", func(c *gin.Context) {
		c.Set("userID", user.ID)
		ApplyForAdoption(c)
	})

	body := map[string]string{
		"message": "I love this dog",
	}
	bodyJSON, _ := json.Marshal(body)

	url := "/adoptions/" + strconv.Itoa(int(pet.ID)) + "/apply"

	req := httptest.NewRequest(http.MethodPost, url, bytes.NewReader(bodyJSON))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d, body=%s", w.Code, w.Body.String())
	}

	var ar models.AdoptionRequest
	if err := database.DB.First(&ar).Error; err != nil {
		t.Fatalf("adoption request not persisted: %v", err)
	}

	if ar.UserID != user.ID || ar.PetID != pet.ID {
		t.Fatalf("adoption request has wrong user or pet: %+v", ar)
	}
}
