package dao

import "time"

type User struct {
	ID uint `gorm:"primaryKey"`
	Username string `gorm:"column:username" json:"username"`
	Email string `gorm:"column:email" json:"email"`
	Age int `gorm:"column:age" json:"age"`
	Is_active bool `gorm:"column:is_active" json:"is_active"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
