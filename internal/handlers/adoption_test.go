package handlers

import (
	"net/http"
	"net/http/httptest"
	"pet-adoption-api/internal/database"
	"pet-adoption-api/internal/models"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Note: No setup function here. It uses the TestMain from shelter_test.go

// TestApplyForAdoption_CreatesRequest uses the testRouter and DB initialized in shelter_test.go's TestMain
func TestApplyForAdoption_CreatesRequest(t *testing.T) {
	// The base data created in TestMain includes a pet with ID 1.
	// We also have an adminToken for an authenticated user.
	petID := 1

	w := httptest.NewRecorder()
	// The URL is now just the path, as the server is managed by httptest
	req, _ := http.NewRequest("POST", "/adoptions/"+strconv.Itoa(petID)+"/apply", nil)
	req.Header.Set("Authorization", "Bearer "+adminToken) // Use the global admin token

	testRouter.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusCreated, w.Code)

	// Check that an adoption request was created in the database
	var request models.AdoptionRequest
	err := database.DB.First(&request, "pet_id = ?", petID).Error
	assert.Nil(t, err, "Adoption request should be found in the database")
	assert.Equal(t, uint(petID), request.PetID)
}
