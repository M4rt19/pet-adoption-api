package repository

import (
	"pet-adoption-api/internal/database"
	"pet-adoption-api/internal/models"
)

func GetAllPets() ([]models.Pet, error){
	var pets []models.Pet
	result := database.DB.Find(&pets)
	return pets, result.Error
}

func CreatePet(pet *models.Pet) error{
	return database.DB.Create(pet).Error
}

func DeletePet(id uint) error{
	return database.DB.Delete(&models.Pet{}, id).Error
}