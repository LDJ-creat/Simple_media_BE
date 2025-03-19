package model

import (
	"gorm.io/gorm"
)

type Notification struct {
	gorm.Model
	SenderID   uint   `gorm:"not null"` // 触发通知的用户（评论者）
	ReceiverID uint   `gorm:"not null"` // 接收通知的用户（被评论帖子的作者）
	PostID     uint   `gorm:"not null"` // 关联的帖子ID
	Type       string `gorm:"type:varchar(20);ENUM('comment','like')"`
	IsRead     bool   `gorm:"column:is_read;default:false"` // 是否已读
}
