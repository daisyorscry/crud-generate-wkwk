package dao

import "time"

type Todo struct {
	ID uint `gorm:"primaryKey"`
	Id string `gorm:"column:id" json:"id"`
	Title string `gorm:"column:title" json:"title"`
	Is_done bool `gorm:"column:is_done" json:"is_done"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
