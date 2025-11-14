package models

type User struct {
    ID       uint   `json:"id" gorm:"primaryKey"`
    Email    string `json:"email" gorm:"unique"`
    Password string `json:"-"`        
    Role     string `json:"role"`   
}
