package model

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"type:varchar(20);not null;unique"`
	Password string `gorm:"type:varchar(100);not null"`
	Email    string
	Avatar   string `gorm:"type:varchar(200);default:'/static/avatars/default.png'"`
}
