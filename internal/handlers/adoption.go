package handlers

import (
	"net/http"
	"strconv"
	"time"

	"pet-adoption-api/internal/database"
	"pet-adoption-api/internal/models"
	"pet-adoption-api/internal/worker"

	"github.com/gin-gonic/gin"
)

// This is set in main.go: handlers.AdoptionEvents = aw.Events
var AdoptionEvents chan worker.AdoptionEvent

// POST /adoptions/:petID/apply
type applyAdoptionRequest struct {
	Message string `json:"message"`
}

func ApplyForAdoption(c *gin.Context) {
	userIDVal, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found in context"})
		return
	}
	userID, ok := userIDVal.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user id type"})
		return
	}

	petIDStr := c.Param("petID")
	petIDInt, err := strconv.Atoi(petIDStr)
	if err != nil || petIDInt <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid pet id"})
		return
	}
	petID := uint(petIDInt)

	// check pet exists and is available
	var pet models.Pet
	if err := database.DB.First(&pet, petID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "pet not found"})
		return
	}
	if pet.Status != models.PetStatusAvailable {
		c.JSON(http.StatusBadRequest, gin.H{"error": "pet is not available for adoption"})
		return
	}

	var reqBody applyAdoptionRequest
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		// message is optional, but invalid JSON should still error
		reqBody.Message = ""
	}

	ar := models.AdoptionRequest{
		UserID:    userID,
		PetID:     petID,
		Status:    models.AdoptionStatusPending,
		Message:   reqBody.Message,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := database.DB.Create(&ar).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create adoption request"})
		return
	}

	// fire async event to worker (non-blocking)
	if AdoptionEvents != nil {
		select {
		case AdoptionEvents <- worker.AdoptionEvent{
			RequestID: ar.ID,
			UserID:    ar.UserID,
			PetID:     ar.PetID,
			Status:    string(ar.Status),
			Message:   "New adoption request created",
		}:
		default:
			// channel full → skip silently
		}
	}

	c.JSON(http.StatusCreated, gin.H{"adoption_request": ar})
}

// GET /adoptions/my
func GetMyAdoptions(c *gin.Context) {
	userIDVal, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found in context"})
		return
	}
	userID, ok := userIDVal.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user id type"})
		return
	}

	var requests []models.AdoptionRequest

	if err := database.DB.
		Preload("Pet").
		Where("user_id = ?", userID).
		Find(&requests).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch adoption requests"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"adoption_requests": requests})
}

// GET /adoptions/shelter (ShelterOnly)
func GetShelterAdoptions(c *gin.Context) {
	userIDVal, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found in context"})
		return
	}
	userID, ok := userIDVal.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user id type"})
		return
	}

	var requests []models.AdoptionRequest

	// all requests for pets whose shelter owner is this user
	if err := database.DB.
		Joins("JOIN pets ON pets.id = adoption_requests.pet_id").
		Joins("JOIN shelters ON shelters.id = pets.shelter_id").
		Where("shelters.owner_user_id = ?", userID).
		Preload("Pet").
		Find(&requests).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch shelter adoption requests"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"adoption_requests": requests})
}

// PATCH /adoptions/:id/approve
func ApproveAdoption(c *gin.Context) {
	updateAdoptionStatus(c, models.AdoptionStatusApproved)
}

// PATCH /adoptions/:id/reject
func RejectAdoption(c *gin.Context) {
	updateAdoptionStatus(c, models.AdoptionStatusRejected)
}

// helper: approve / reject
func updateAdoptionStatus(c *gin.Context, newStatus models.AdoptionStatus) {
	userIDVal, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found in context"})
		return
	}
	userID, ok := userIDVal.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user id type"})
		return
	}

	idStr := c.Param("id")
	idInt, err := strconv.Atoi(idStr)
	if err != nil || idInt <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid adoption request id"})
		return
	}
	id := uint(idInt)

	var ar models.AdoptionRequest
	if err := database.DB.
		Preload("Pet").
		Preload("Pet.Shelter").
		First(&ar, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "adoption request not found"})
		return
	}

	// check that current user is owner of the shelter for this pet OR admin
	roleVal, _ := c.Get("role")
	role, _ := roleVal.(string)

	isOwner := ar.Pet.Shelter.OwnerUserID == userID
	isAdmin := role == "admin"

	if !isOwner && !isAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "not allowed to update this request"})
		return
	}

	// update request status
	ar.Status = newStatus
	ar.UpdatedAt = time.Now()

	// if approved, set pet to adopted
	if newStatus == models.AdoptionStatusApproved {
		ar.Pet.Status = models.PetStatusAdopted
		if err := database.DB.Save(&ar.Pet).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update pet status"})
			return
		}
	}

	if err := database.DB.Save(&ar).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update adoption request"})
		return
	}

	// fire async event to worker (non-blocking)
	if AdoptionEvents != nil {
		select {
		case AdoptionEvents <- worker.AdoptionEvent{
			RequestID: ar.ID,
			UserID:    ar.UserID,
			PetID:     ar.PetID,
			Status:    string(newStatus),
			Message:   "Adoption request status updated",
		}:
		default:
			// channel full → skip
		}
	}

	c.JSON(http.StatusOK, gin.H{"adoption_request": ar})
}
