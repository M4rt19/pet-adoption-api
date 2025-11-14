package repository

import (
	"pet-adoption-api/internal/database"
    "pet-adoption-api/internal/models"
)

func CreateUser(user *models.User) error{
	return database.DB.Create(user).Error
}

func GetUserByEmail(email string) (*models.User, error){
	var u models.User
	err := database.DB.Where("email= ?", email).First(&u).Error
	return &u, err
}