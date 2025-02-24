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
	Code      string    //忘记密码时邮箱验证存储验证码
	CodeTime  time.Time //验证码过期时间
	Avatar    string    `gorm:"type:varchar(200);default:'/static/avatars/default.png'"`
	Signature string    `gorm:"type:varchar(200);default:' '"`
}
