package models

import "time"

type PetStatus string

const (
	PetStatusAvailable PetStatus = "available"
	PetStatusReserved  PetStatus = "reserved"
	PetStatusAdopted   PetStatus = "adopted"
)

type Pet struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	ShelterID   uint      `gorm:"not null" json:"shelter_id"`
	Name        string    `gorm:"not null" json:"name"`
	Species     string    `gorm:"not null" json:"species"` // dog, cat, etc.
	Breed       string    `json:"breed"`
	Age         int       `json:"age"`
	Description string    `json:"description"`
	Status      PetStatus `gorm:"type:varchar(20);not null;default:'available'" json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	Shelter Shelter `gorm:"foreignKey:ShelterID" json:"-"`
}
