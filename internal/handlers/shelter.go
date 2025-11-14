package handlers

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "pet-adoption-api/internal/models"
    "pet-adoption-api/internal/repository"
)

func GetShelters(c *gin.Context) {
	shelters, err := repository.GetAllShelters()
	if err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return 
	}
	c.JSON(http.StatusOK, shelters)
}

func CreateShelter(c *gin.Context){
	var shelter models.Shelter
	if err := c.ShouldBindJSON(&shelter); err!= nil{
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := repository.CreateShelter(&shelter); err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, shelter)
}

