package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username  string `gorm:"type:varchar(20);not null;unique"`
	Password  string `gorm:"type:varchar(100);not null"`
	Email     string
	Phone     string
	Code      string
	CodeTime  *time.Time `gorm:"default:null"` // 改为指针类型，允许为 null
	Avatar    string     `gorm:"type:varchar(200);default:''"`
	Signature string     `gorm:"type:varchar(200);default:' '"`
}
