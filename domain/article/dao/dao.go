package dao

import "time"

type Article struct {
	ID uint `gorm:"primaryKey"`
	Title string `gorm:"column:title" json:"title"`
	Content string `gorm:"column:content" json:"content"`
	Is_active bool `gorm:"column:is_active" json:"is_active"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
