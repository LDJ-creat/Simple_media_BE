package model

import (
	"gorm.io/gorm"
)

type Post struct {
	gorm.Model
	UserID    uint   `gorm:"not null"`
	Content   string `gorm:"type:text;not null"`
	Image     string `gorm:"type:varchar(200);default:' '"`
	Video     string `gorm:"type:varchar(200);default:' '"`
	LikeCount uint   `gorm:"default:0"`
}
