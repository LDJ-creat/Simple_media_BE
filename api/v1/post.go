package v1

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/media/internal/model"
	"github.com/media/pkg/database"
	"gorm.io/gorm"
)

// 定义上传目录的常量
const UPLOAD_DIR = "static/uploads"

func init() {
	// 确保上传目录存在
	if err := os.MkdirAll(UPLOAD_DIR, 0755); err != nil {
		panic(fmt.Sprintf("Failed to create upload directory: %v", err))
	}
}

func CreatePost(c *gin.Context) {
	// 获取当前用户ID
	userID := c.GetUint("userID")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	// 解析multipart form
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无法解析表单数据"})
		return
	}

	// 获取内容
	content := ""
	if contentValues := form.Value["content"]; len(contentValues) > 0 {
		content = contentValues[0]
	}

	// 获取媒体文件
	files := form.File["media[]"]
	media := []model.Media{}

	// 处理每个上传的文件
	for _, file := range files {
		// 生成唯一的文件名
		ext := filepath.Ext(file.Filename)
		filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)

		// 使用绝对路径保存文件
		absPath := filepath.Join(UPLOAD_DIR, filename)
		if err := c.SaveUploadedFile(file, absPath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "保存文件失败"})
			return
		}

		// 判断文件类型
		mediaType := "image"
		ext = strings.ToLower(ext)
		if ext == ".mp4" || ext == ".mov" || ext == ".avi" || ext == ".wmv" ||
			strings.HasPrefix(strings.ToLower(file.Header.Get("Content-Type")), "video/") {
			mediaType = "video"
		}

		media = append(media, model.Media{
			Uri:  "/static/uploads/" + filename,
			Type: mediaType,
		})
	}

	// 创建帖子记录
	post := model.Post{
		UserID:  userID,
		Content: content,
		Media:   media,
	}

	if err := database.DB.Create(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建帖子失败"})
		return
	}

	// 添加WebSocket通知
	message := gin.H{
		"type": "new_post",
		"post": post,
	}
	messageBytes, _ := json.Marshal(message)
	SendNewPostToUser(userID, messageBytes)

	c.JSON(http.StatusOK, gin.H{
		"message": "帖子创建成功",
		"post":    post,
	})
}

func UpdatePost(c *gin.Context) {
	postID := c.Param("id")
	//获取现有帖子
	var post model.Post
	if err := database.DB.First(&post, postID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
	}

	//解析更新内容
	form, _ := c.MultipartForm()
	content := form.Value["content"][0]
	keepMediaIDs := form.Value["keepMediaIDs"] //需要保留的媒体ID
	newMedias := form.File["media[]"]          //新上传的媒体

	//更新帖子内容
	post.Content = content

	//处理媒体
	oldMediaIDs := []uint{}
	for _, id := range keepMediaIDs {
		idInt, err := strconv.Atoi(id)
		if err != nil {
			continue // 跳过无效的ID
		}
		oldMediaIDs = append(oldMediaIDs, uint(idInt))
	}

	//删除不在使用的媒体
	database.DB.Where("post_id=? AND id NOT IN (?)", postID, oldMediaIDs).Delete(&model.Media{})

	//处理新上传的媒体
	for _, file := range newMedias {
		//检查文件大小
		if file.Size > 10<<20 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "File size too large"})
			return
		}

		//生成唯一文件名
		fileExt := filepath.Ext(file.Filename)
		fileName := uuid.New().String() + fileExt

		//保存文件
		dst := filepath.Join("static/uploads", fileName)
		if err := c.SaveUploadedFile(file, dst); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		}

		mediaType := "image"
		if strings.Contains(file.Header.Get("Content-Type"), "video") {
			mediaType = "video"
		}
		media := model.Media{
			Uri:  dst,
			Type: mediaType,
		}
		database.DB.Create(&media)

	}

	//保存更新
	if err := database.DB.Save(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update post"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Post updated successfully"})
}

func DeletePost(c *gin.Context) {
	userID := c.GetUint("userID")
	postID := c.Param("id")
	database.DB.Where("id = ? AND user_id = ?", postID, userID).Delete(&model.Post{})
	c.JSON(http.StatusOK, gin.H{"message": "Post deleted successfully"})
}

func AddLikeCount(c *gin.Context) {
	userID := c.GetUint("userID")
	postID := c.Param("id")

	var post model.Post
	if err := database.DB.First(&post, postID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "帖子未找到"})
		return
	}

	// 检查是否已经点赞
	for _, id := range post.LikeCount {
		if id == userID {
			c.JSON(http.StatusBadRequest, gin.H{"error": "已经点赞过了"})
			return
		}
	}

	// 如果 LikeCount 为 nil，初始化它
	if post.LikeCount == nil {
		post.LikeCount = make([]uint, 0)
	}

	// 添加新的点赞
	post.LikeCount = append(post.LikeCount, userID)

	// 保存更新
	if err := database.DB.Save(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "点赞失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "点赞成功"})
}

func SubLikeCount(c *gin.Context) {
	userID := c.GetUint("userID")
	postID := c.Param("id")

	var post model.Post
	if err := database.DB.First(&post, postID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "帖子未找到"})
		return
	}

	// 如果 LikeCount 为 nil，初始化它
	if post.LikeCount == nil {
		post.LikeCount = make([]uint, 0)
		c.JSON(http.StatusBadRequest, gin.H{"error": "还没有点赞"})
		return
	}

	// 移除点赞
	newLikeCount := make([]uint, 0)
	found := false
	for _, id := range post.LikeCount {
		if id != userID {
			newLikeCount = append(newLikeCount, id)
		} else {
			found = true
		}
	}

	if !found {
		c.JSON(http.StatusBadRequest, gin.H{"error": "还没有点赞"})
		return
	}

	post.LikeCount = newLikeCount

	// 保存更新
	if err := database.DB.Save(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "取消点赞失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "取消点赞成功"})
}

// 用于主页面的渲染，只返回media中的一个
func GetPosts(c *gin.Context) {
	lastID := c.DefaultQuery("last_id", "0")
	pageSize := 10

	var posts []model.Post
	query := database.DB.
		Preload("Media").
		Preload("User").
		Order("id DESC").
		Limit(pageSize)

	if lastID != "0" {
		query = query.Where("id < ?", lastID)
	} else {
		query = query.Where("id > ?", 0)
	}

	if err := query.Find(&posts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var newCursor string
	if len(posts) > 0 {
		newCursor = fmt.Sprintf("%d", posts[len(posts)-1].ID)
	}
	c.Header("X-Next-Cursor", newCursor)

	c.JSON(http.StatusOK, gin.H{"posts": posts})
}

func GetPostByID(c *gin.Context) {
	postID := c.Param("id")
	var post model.Post

	// 修改查询以确保加载所有关联数据
	if err := database.DB.
		Preload("Media").
		Preload("User").
		Preload("Comment", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at DESC") // 按创建时间倒序排列评论
		}).
		Preload("Comment.User", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, username, avatar") // 只选择需要的用户字段
		}).
		First(&post, postID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "帖子未找到"})
		return
	}

	// 确保评论数据完整性
	for i := range post.Comment {
		if post.Comment[i].User.ID == 0 {
			// 如果评论的用户信息缺失，尝试重新获取
			var user model.User
			if err := database.DB.Select("id, username, avatar").
				First(&user, post.Comment[i].UserID).Error; err == nil {
				post.Comment[i].User = user
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"post":    post,
		"success": true,
	})
}

func AddComment(c *gin.Context) {
	// 修改接收结构体以匹配前端数据结构
	var request struct {
		PostID  json.Number `json:"PostID"`
		Content string      `json:"Content"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  "解析请求数据失败",
			"detail": err.Error(),
		})
		return
	}

	// 转换 PostID 为 uint
	postIDInt, err := request.PostID.Int64()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  "无效的帖子ID",
			"detail": err.Error(),
		})
		return
	}
	postID := uint(postIDInt)

	// 验证帖子是否存在
	var post model.Post
	if err := database.DB.First(&post, postID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "帖子不存在"})
		return
	}

	userID := c.GetUint("userID")

	comment := model.Comment{
		PostID:  postID,
		Content: request.Content,
		UserID:  userID,
	}

	// 保存评论
	if err := database.DB.Create(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "添加评论失败",
			"detail": err.Error(),
		})
		return
	}

	// 获取帖子作者ID
	receiverID := post.UserID

	// 如果评论者不是帖子作者，则创建通知
	if comment.UserID != receiverID {
		notification := model.Notification{
			SenderID:   comment.UserID,
			ReceiverID: receiverID,
			PostID:     postID,
			Type:       "comment",
		}
		database.DB.Create(&notification)

		message := `{"type": "notification", "count": 1}`
		SendNotificationToUser(receiverID, []byte(message))
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "评论添加成功",
		"comment": comment,
	})
}

func DeleteComment(c *gin.Context) {
	userID := c.GetUint("userID")
	postID := c.Param("id")

	// 修改 SQL 查询中的列名 postID 为 post_id
	if err := database.DB.Where("post_id = ? AND user_id = ?", postID, userID).Delete(&model.Comment{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除评论失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "评论删除成功"})
}

// 获取我的帖子
func GetMyPosts(c *gin.Context) {
	lastID := c.DefaultQuery("last_id", "0")
	pageSize := 10
	userID := c.GetUint("userID")

	var posts []model.Post
	query := database.DB.
		Preload("Media", func(db *gorm.DB) *gorm.DB {
			return db.Order("id DESC").Limit(1) // 按ID降序取第一条媒体
		}).
		Preload("User").
		Order("id DESC").
		Limit(pageSize)

	if lastID != "0" {
		query = query.Where("id < ? AND user_id = ?", lastID, userID)
	} else {
		query = query.Where("id > ? AND user_id = ?", 0, userID)
	}

	if err := query.Find(&posts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var newCursor string
	if len(posts) > 0 {
		newCursor = fmt.Sprintf("%d", posts[len(posts)-1].ID)
	}
	c.Header("X-Next-Cursor", newCursor)

	c.JSON(http.StatusOK, gin.H{"posts": posts})
}
