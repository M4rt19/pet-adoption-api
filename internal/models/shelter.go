package models
import "gorm.io/gorm" 
type Shelter struct {
    gorm.Model
    Name     string
    Location string
    Pets     []Pet `gorm:"foreignKey:ShelterID"`
}