package handlers

import (
	"net/http"
	"strconv"


	"github.com/gin-gonic/gin"
	"pet-adoption-api/internal/models"
	"pet-adoption-api/internal/repository"
)

func GetPets(c *gin.Context){
	pets , err := repository.GetAllPets()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return 
	}
	c.JSON(http.StatusOK, pets)
}

func CreatePet(c *gin.Context){
	var pet models.Pet
	if err := c.ShouldBindJSON(&pet); err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return 
	}
	c.JSON(http.StatusCreated, pet)
}

func DeletePet(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := repository.DeletePet(uint(id)); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"message": "Pet deleted"})
}