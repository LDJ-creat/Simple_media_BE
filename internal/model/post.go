package model

import (
	"gorm.io/gorm"
)

type Image struct {
	ID   uint
	Path string
}

type Video struct {
	ID   uint
	Path string
}

type Post struct {
	gorm.Model
	UserID    uint    `gorm:"not null"`
	Content   string  `gorm:"type:text;not null"`
	Image     []Image `gorm:"type:varchar(200);default:' '"`
	Video     []Video `gorm:"type:varchar(200);default:' '"`
	LikeCount uint    `gorm:"default:0"`
}
