package models

import "gorm.io/gorm" 

type Pet struct {
    gorm.Model
    Name      string
    Breed     string
    Age       int
    ShelterID uint
}
