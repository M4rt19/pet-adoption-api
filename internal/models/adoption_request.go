package models

import "time"

type AdoptionStatus string

const (
	AdoptionStatusPending   AdoptionStatus = "pending"
	AdoptionStatusApproved  AdoptionStatus = "approved"
	AdoptionStatusRejected  AdoptionStatus = "rejected"
	AdoptionStatusCancelled AdoptionStatus = "cancelled"
	AdoptionStatusExpired   AdoptionStatus = "expired"
)

type AdoptionRequest struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	UserID    uint           `gorm:"not null" json:"user_id"`
	PetID     uint           `gorm:"not null" json:"pet_id"`
	Status    AdoptionStatus `gorm:"type:varchar(20);not null;default:'pending'" json:"status"`
	Message   string         `json:"message"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`

	User User `gorm:"foreignKey:UserID" json:"-"`
	Pet  Pet  `gorm:"foreignKey:PetID" json:"-"`
}
