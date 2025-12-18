package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"pet-adoption-api/internal/database"
	"pet-adoption-api/internal/middleware"
	"pet-adoption-api/internal/models"
	"pet-adoption-api/internal/worker"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var testRouter *gin.Engine
var adminToken string
var shelterOwnerUser models.User

// TestMain sets up a centralized test environment for the entire handlers package
func TestMain(m *testing.M) {
	// Set up an in-memory SQLite database for testing
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to test database: " + err.Error())
	}
	database.DB = db

	// Run migrations for all models
	db.AutoMigrate(&models.User{}, &models.Pet{}, &models.Shelter{}, &models.AdoptionRequest{})

	// Set up the router
	gin.SetMode(gin.TestMode)
	testRouter = gin.Default()

	// Initialize auth and middleware
	InitAuth()
	middleware.InitAuthMiddleware()

	// Initialize a dummy channel for adoption events
	AdoptionEvents = make(chan worker.AdoptionEvent, 10)

	// --- Centralized Route Setup ---
	authRoutes := testRouter.Group("/auth")
	{
		authRoutes.POST("/register", Register)
		authRoutes.POST("/login", Login)
	}

	shelterRoutes := testRouter.Group("/shelters")
	{
		shelterRoutes.GET("/", GetShelters)
		shelterRoutes.GET("/:id", GetShelterByID)
		shelterRoutes.POST("/", middleware.AuthMiddleware(), middleware.AdminOnly(), CreateShelter)
		shelterRoutes.PUT("/:id", middleware.AuthMiddleware(), middleware.AdminOnly(), UpdateShelter)
		shelterRoutes.DELETE("/:id", middleware.AuthMiddleware(), middleware.AdminOnly(), DeleteShelter)
	}

	adoptionRoutes := testRouter.Group("/adoptions", middleware.AuthMiddleware())
	{
		adoptionRoutes.POST("/:petID/apply", ApplyForAdoption)
	}

	// Create base test data (users, tokens, a shelter, a pet)
	createBaseTestData()

	// Run all tests in the package
	code := m.Run()

	// Clean up
	sqlDB, _ := database.DB.DB()
	sqlDB.Close()

	os.Exit(code)
}

func createBaseTestData() {
	// Create an admin user
	hash, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	adminUser := models.User{Name: "Admin User", Email: "admin@test.com", PasswordHash: string(hash), Role: "admin"}
	database.DB.Create(&adminUser)

	// Create a shelter owner user
	shelterOwnerUser = models.User{Name: "Shelter Owner", Email: "owner@test.com", PasswordHash: string(hash), Role: "shelter"}
	database.DB.Create(&shelterOwnerUser)

	// Create a base shelter for other tests to use
	testShelter := models.Shelter{Name: "Base Shelter", OwnerUserID: shelterOwnerUser.ID}
	database.DB.Create(&testShelter) // This will have ID 1

	// Create a base pet for adoption tests
	testPet := models.Pet{Name: "Base Pet", Species: "Dog", ShelterID: testShelter.ID}
	database.DB.Create(&testPet) // This will have ID 1

	// Generate token for the admin user
	adminToken, _ = jwtManager.Generate(adminUser.ID, string(adminUser.Role))
}

// Test GetShelters endpoint
func TestGetShelters(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/shelters/", nil)
	testRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string][]models.Shelter
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.NotEmpty(t, response["shelters"]) // Should have at least the base shelter
}

// Test CreateShelter endpoint
func TestCreateShelter(t *testing.T) {
	t.Run("Admin can create a shelter", func(t *testing.T) {
		shelterData := gin.H{
			"name":          "New Shelter by Admin",
			"address":       "456 Admin Ave",
			"owner_user_id": shelterOwnerUser.ID,
		}
		body, _ := json.Marshal(shelterData)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/shelters/", bytes.NewBuffer(body))
		req.Header.Set("Authorization", "Bearer "+adminToken)
		req.Header.Set("Content-Type", "application/json")
		testRouter.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code, "Expected status 201 for admin creation")

		var response map[string]models.Shelter
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, "New Shelter by Admin", response["shelter"].Name)
	})

	t.Run("Fails without token", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/shelters/", nil)
		testRouter.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

// Test GetShelterByID endpoint
func TestGetShelterByID(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/shelters/1", nil) // Check for the base shelter with ID 1
	testRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]models.Shelter
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Base Shelter", response["shelter"].Name)
}

// Test UpdateShelter endpoint
func TestUpdateShelter(t *testing.T) {
	updateData := gin.H{"name": "Updated Base Shelter"}
	body, _ := json.Marshal(updateData)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/shelters/1", bytes.NewBuffer(body)) // Update base shelter
	req.Header.Set("Authorization", "Bearer "+adminToken)
	req.Header.Set("Content-Type", "application/json")
	testRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]models.Shelter
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Updated Base Shelter", response["shelter"].Name)
}

// Test DeleteShelter endpoint
func TestDeleteShelter(t *testing.T) {
    // Create a shelter specifically for this test so we don't break other tests
    tempShelter := models.Shelter{Name: "To Be Deleted", OwnerUserID: shelterOwnerUser.ID}
    database.DB.Create(&tempShelter)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/shelters/"+strconv.Itoa(int(tempShelter.ID)), nil)
	req.Header.Set("Authorization", "Bearer "+adminToken)
	testRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// Verify it's actually deleted
	var deletedShelter models.Shelter
	err := database.DB.First(&deletedShelter, tempShelter.ID).Error
	assert.NotNil(t, err)
}
