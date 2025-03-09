package model

import (
	"gorm.io/gorm"
)

// type Image struct {
// 	ID   uint
// 	Path string
// }

// type Video struct {
// 	ID   uint
// 	Path string
// }

type Post struct {
	gorm.Model
	UserID    uint      `gorm:"not null"`
	Content   string    `gorm:"type:text"`
	LikeCount []uint    `gorm:"type:json"`         //存储点赞用户的ID
	Media     []Media   `gorm:"foreignKey:PostID"` //这个字段不会存入post表，而是存到media表
	User      User      `gorm:"foreignKey:UserID"`
	Comment   []Comment `gorm:"foreignKey:PostID"`
}

type Media struct {
	gorm.Model
	PostID uint   `gorm:"not null"`
	Uri    string `gorm:"type:varchar(200);not null"`
	Type   string `gorm:"type:varchar(20);not null;ENUM('image','video')"`
	Post   Post   `gorm:"foreignKey:PostID;constraint:OnDelete:CASCADE" json:"-"` //相互关联,因为循环引用，所以用json:"-"表示忽略POST字段的输出

}

type Comment struct {
	gorm.Model
	PostID  uint   `gorm:"not null"`
	UserID  uint   `gorm:"not null"`
	Content string `gorm:"type:text"`
	Post    Post   `gorm:"foreignKey:PostID;constraint:OnDelete:CASCADE" json:"-"` //相互关联
	User    User   `gorm:"foreignKey:UserID" json:"-"`
}
