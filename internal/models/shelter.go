package models

import "time"

type Shelter struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"not null" json:"name"`
	Address     string    `json:"address"`
	Phone       string    `json:"phone"`
	OwnerUserID uint      `gorm:"not null" json:"owner_user_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	OwnerUser User  `gorm:"foreignKey:OwnerUserID" json:"-"`
	Pets      []Pet `json:"pets,omitempty"`
}
