package v1

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	// "github.com/media/internal/model"
	// "github.com/media/pkg/database"
)

func MediaUpload(c *gin.Context, file *multipart.FileHeader, filePath string) string {
	// 获取上传的文件
	// file, err := c.FormFile(filePath)
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "请选择文件"})
	// 	return ""
	// }

	// 验证文件大小（例如限制为2MB）
	if file.Size > 2<<20 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "文件大小不能超过2MB"})
		return ""
	}

	// // 验证文件类型
	ext := strings.ToLower(path.Ext(file.Filename))
	// if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".gif" {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "只支持jpg、jpeg、png、gif格式"})
	// 	return
	// }

	// 生成随机文件名
	fileName := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	mediaPath := fmt.Sprintf("static/%s/%s", filePath, fileName)

	// 确保目录存在
	if err := os.MkdirAll(fmt.Sprintf("static/%s", filePath), 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建目录失败"})
		return ""
	}

	// 保存文件
	if err := c.SaveUploadedFile(file, mediaPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "保存文件失败"})
		return ""
	}

	return "/" + mediaPath // 返回文件URL
}
