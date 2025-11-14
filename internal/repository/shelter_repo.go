package repository

import(
	"pet-adoption-api/internal/database"
	"pet-adoption-api/internal/models"
)

func GetAllShelters() ([]models.Shelter, error){
	var shelters []models.Shelter
	result := database.DB.Preload("Pets").Find(&shelters)
	return shelters, result.Error
}

func CreateShelter(shelter *models.Shelter) error{
	return database.DB.Create(shelter).Error
}
