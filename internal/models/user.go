package models

import (
	"database/sql"
	"time"

	"gorm.io/gorm"
)

type User struct {
	Id           string         `json:"id" gorm:"type:uuid;column:id;primaryKey"`
	Name         string         `json:"name" gorm:"type:varchar(255);not null"`
	Document     string         `json:"document" gorm:"type:varchar(14);unique;not null"`
	Email        string         `json:"email" gorm:"type:varchar(50);unique;not null"`
	Phone        string         `json:"phone" gorm:"type:varchar(11);not null"`
	PasswordHash string         `json:"password_hash" gorm:"type:varchar(255);not null"`
	CreatedAt    time.Time      `json:"created_at" gorm:"not null"`
	UpdatedAt    sql.NullTime   `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"deleted_at"`
}
