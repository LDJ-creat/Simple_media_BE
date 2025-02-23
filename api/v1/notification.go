package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/media/internal/model"
	"github.com/media/pkg/database"
)

// 获取未读通知数量
func GetNotificationsCount(c *gin.Context) {
	userID := c.GetUint("userID") // 从JWT获取当前用户
	var count int64
	database.DB.Model(&model.Notification{}).Where("receiver_id = ? AND read = ?", userID, false).Count(&count)
	c.JSON(200, gin.H{"count": count})
}

// 获取通知列表
func GetNotifications(c *gin.Context) {
	userID := c.GetUint("userID")
	var notifications []model.Notification
	database.DB.Preload("Sender").Preload("Post").
		Where("receiver_id = ?", userID).
		Order("created_at DESC").
		Find(&notifications)

	// 标记为已读
	database.DB.Model(&model.Notification{}).Where("receiver_id = ?", userID).Update("read", true)

	c.JSON(200, notifications)
}
