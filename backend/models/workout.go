package models

import "gorm.io/gorm"

type Workout struct {
	gorm.Model
	UserID int    `json:"user_id" gorm:"not null"`
	Title  string `json:"title"`
	Date   string `json:"date"`
	Notes  string `json:"notes"`
}
