package handlers

import (
	"net/http"
	"strconv"
	"time"

	"pet-adoption-api/internal/database"
	"pet-adoption-api/internal/models"

	"github.com/gin-gonic/gin"
)

// GET /pets
func GetPets(c *gin.Context) {
	var pets []models.Pet

	// Optional filters: ?status=available&species=dog
	status := c.Query("status")
	species := c.Query("species")

	query := database.DB.Model(&models.Pet{})
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if species != "" {
		query = query.Where("species = ?", species)
	}

	if err := query.Find(&pets).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch pets"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"pets": pets})
}

type createPetRequest struct {
	ShelterID   uint   `json:"shelter_id" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Species     string `json:"species" binding:"required"`
	Breed       string `json:"breed"`
	Age         int    `json:"age"`
	Description string `json:"description"`
}

// POST /pets
func CreatePet(c *gin.Context) {
	var req createPetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	// Optional: verify shelter exists
	var shelter models.Shelter
	if err := database.DB.First(&shelter, req.ShelterID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "shelter not found"})
		return
	}

	pet := models.Pet{
		ShelterID:   req.ShelterID,
		Name:        req.Name,
		Species:     req.Species,
		Breed:       req.Breed,
		Age:         req.Age,
		Description: req.Description,
		Status:      models.PetStatusAvailable,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := database.DB.Create(&pet).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create pet"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"pet": pet})
}

// DELETE /pets/:id
func DeletePet(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid pet id"})
		return
	}

	if err := database.DB.Delete(&models.Pet{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete pet"})
		return
	}

	c.Status(http.StatusNoContent)
}
