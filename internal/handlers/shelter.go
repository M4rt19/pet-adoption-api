package handlers

import (
	"net/http"
	"strconv"
	"time"

	"pet-adoption-api/internal/database"
	"pet-adoption-api/internal/models"

	"github.com/gin-gonic/gin"
)

// GET /shelters
func GetShelters(c *gin.Context) {
	var shelters []models.Shelter

	if err := database.DB.Find(&shelters).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch shelters"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"shelters": shelters})
}

// GET /shelters/:id
func GetShelterByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid shelter id"})
		return
	}

	var shelter models.Shelter
	if err := database.DB.First(&shelter, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "shelter not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"shelter": shelter})
}

type createShelterRequest struct {
	Name    string `json:"name" binding:"required"`
	Address string `json:"address"`
	Phone   string `json:"phone"`
	// Optional: owner_user_id; or we can assign current user if role == shelter
	OwnerUserID uint `json:"owner_user_id" binding:"required"`
}

// POST /shelters  (admin only in routes)
func CreateShelter(c *gin.Context) {
	var req createShelterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	shelter := models.Shelter{
		Name:        req.Name,
		Address:     req.Address,
		Phone:       req.Phone,
		OwnerUserID: req.OwnerUserID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := database.DB.Create(&shelter).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create shelter"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"shelter": shelter})
}

type updateShelterRequest struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	Phone   string `json:"phone"`
}

// PUT /shelters/:id
func UpdateShelter(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid shelter id"})
		return
	}

	var shelter models.Shelter
	if err := database.DB.First(&shelter, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "shelter not found"})
		return
	}

	var req updateShelterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	// Update fields from request
	shelter.Name = req.Name
	shelter.Address = req.Address
	shelter.Phone = req.Phone
	shelter.UpdatedAt = time.Now()

	if err := database.DB.Save(&shelter).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update shelter"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"shelter": shelter})
}

// DELETE /shelters/:id
func DeleteShelter(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid shelter id"})
		return
	}

	if err := database.DB.Delete(&models.Shelter{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete shelter"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "shelter deleted successfully"})
}
