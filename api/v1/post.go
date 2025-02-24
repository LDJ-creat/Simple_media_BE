package v1

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/media/internal/model"
	"github.com/media/pkg/database"
	"gorm.io/gorm"
)

// func CreatePost(c *gin.Context) {
// 	userID := c.GetUint("userID")
// 	// var newPost newPost
// 	// if err:=c.ShouldBindJSON(&newPost);err!=nil{
// 	// 	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 	// }
// 	content := c.PostForm("content")
// 	imagesPath := []model.Media{}
// 	videosPath := []model.Media{}
// 	form, err := c.MultipartForm()
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get multipart form"})
// 		return
// 	}
// 	images := form.File["images"]
// 	videos := form.File["videos"]
// 	for _, image := range images {
// 		var imageItem model.Image
// 		imageItem.Path = MediaUpload(c, image, "image")
// 		// 获取文件id
// 		src, err := image.Open()
// 		if err != nil {
// 			c.JSON(500, gin.H{"error": fmt.Sprintf("打开文件 %s 失败: %v", image.Filename, err)})
// 			return
// 		}
// 		defer src.Close()

// 		// 这里需要解析前端传递的额外信息（id、type 等）
// 		var mediaInfo struct {
// 			ID   int    `json:"id"`
// 			Type string `json:"type"`
// 		}
// 		err = json.Unmarshal([]byte(image.Filename), &mediaInfo)
// 		if err != nil {
// 			c.JSON(500, gin.H{"error": fmt.Sprintf("解析文件信息 %s 失败: %v", image.Filename, err)})
// 			return
// 		}
// 		imageItem.ID = uint(mediaInfo.ID)
// 		imagesPath = append(imagesPath, imageItem)
// 	}
// 	for _, video := range videos {
// 		var videoItem model.Video
// 		videoItem.Path = MediaUpload(c, video, "video")

// 		// 获取文件id
// 		src, err := video.Open()
// 		if err != nil {
// 			c.JSON(500, gin.H{"error": fmt.Sprintf("打开文件 %s 失败: %v", video.Filename, err)})
// 			return
// 		}
// 		defer src.Close()
// 		//解析出id
// 		var mediaInfo struct {
// 			ID   int    `json:"id"`
// 			Type string `json:"type"`
// 		}
// 		err = json.Unmarshal([]byte(video.Filename), &mediaInfo)
// 		if err != nil {
// 			c.JSON(500, gin.H{"error": fmt.Sprintf("解析文件信息 %s 失败: %v", video.Filename, err)})
// 			return
// 		}
// 		videoItem.ID = uint(mediaInfo.ID)
// 		videosPath = append(videosPath, videoItem)
// 	}

// 	post := model.Post{
// 		UserID:  userID,
// 		Content: body,
// 		Image:   imagesPath,
// 		Video:   videosPath,
// 	}

// 	if err := database.DB.Create(&post).Error; err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create post"})
// 		return
// 	}
// 	c.JSON(http.StatusOK, gin.H{"message": "Post created successfully"})

// }

// func ModifyPost(c *gin.Context) {
// 	userID := c.GetUint("userID")
// 	postID := c.Param("id")
// 	body := c.PostForm("content")
// 	imagesPath := []model.Image{}
// 	oddImages := []model.Image{}
// 	//将数据库中的图片信息拿出来
// 	database.DB.Where("id = ? AND user_id = ?", postID, userID).Find(&oddImages)
// 	fmt.Println(oddImages)
// 	videosPath := []model.Video{}
// 	oddVideos := []model.Video{}
// 	//将数据库中的视频信息拿出来
// 	database.DB.Where("id = ? AND user_id = ?", postID, userID).Find(&oddVideos)
// 	fmt.Println(oddVideos)
// 	form, err := c.MultipartForm()
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get multipart form"})
// 		return
// 	}
// 	images := form.File["images"]
// 	videos := form.File["videos"]
// 	//处理图片
// 	for _, image := range images {
// 		var imageItem model.Image
// 		// 获取文件id
// 		src, err := image.Open()
// 		if err != nil {
// 			c.JSON(500, gin.H{"error": fmt.Sprintf("打开文件 %s 失败: %v", image.Filename, err)})
// 			return
// 		}
// 		defer src.Close()

// 		// 这里需要解析前端传递的额外信息（id、type 等）
// 		var mediaInfo struct {
// 			ID   int    `json:"id"`
// 			Type string `json:"type"`
// 		}
// 		err = json.Unmarshal([]byte(image.Filename), &mediaInfo)
// 		if err != nil {
// 			c.JSON(500, gin.H{"error": fmt.Sprintf("解析文件信息 %s 失败: %v", image.Filename, err)})
// 			return
// 		}
// 		imageItem.ID = uint(mediaInfo.ID)
// 		isExist := false
// 		for _, oddImage := range oddImages {
// 			if oddImage.ID == imageItem.ID {
// 				isExist = true
// 				imagesPath = append(imagesPath, oddImage)
// 				break
// 			}
// 		}
// 		if !isExist {
// 			imageItem.Path = MediaUpload(c, image, "image")
// 			imagesPath = append(imagesPath, imageItem)
// 		}

// 	}
// 	//处理视频
// 	for _, video := range videos {
// 		var videoItem model.Video
// 		// 获取文件id
// 		src, err := video.Open()
// 		if err != nil {
// 			c.JSON(500, gin.H{"error": fmt.Sprintf("打开文件 %s 失败: %v", video.Filename, err)})
// 			return
// 		}
// 		defer src.Close()
// 		//解析出id
// 		var mediaInfo struct {
// 			ID   int    `json:"id"`
// 			Type string `json:"type"`
// 		}
// 		err = json.Unmarshal([]byte(video.Filename), &mediaInfo)
// 		if err != nil {
// 			c.JSON(500, gin.H{"error": fmt.Sprintf("解析文件信息 %s 失败: %v", video.Filename, err)})
// 			return
// 		}
// 		videoItem.ID = uint(mediaInfo.ID)
// 		isExist := false
// 		for _, oldVideo := range oddVideos {
// 			if oldVideo.ID == videoItem.ID {
// 				isExist = true
// 				videosPath = append(videosPath, oldVideo)
// 				break
// 			}
// 		}
// 		if !isExist {
// 			videoItem.Path = MediaUpload(c, video, "video")
// 			videosPath = append(videosPath, videoItem)
// 		}
// 	}

// 	post := model.Post{
// 		UserID:  userID,
// 		Content: body,
// 		Image:   imagesPath,
// 		Video:   videosPath,
// 	}

// 	database.DB.Model(&model.Post{}).Where("id = ? AND user_id = ?", postID, userID).Updates(post)
// 	c.JSON(http.StatusOK, gin.H{"message": "Post updated successfully"})
// }

func CreatePost(c *gin.Context) {
	userID := c.GetUint("userID")
	form, _ := c.MultipartForm()
	content := form.Value["content"][0]
	post := model.Post{
		Content: content,
		UserID:  userID,
	}

	files := form.File["media"]
	medias := []model.Media{}
	for _, file := range files {
		//检查文件大小
		if file.Size > 10<<20 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "File size too large"})
			return
		}

		//生成唯一文件名
		fileExt := filepath.Ext(file.Filename)
		fileName := uuid.New().String() + fileExt

		// 确保目录存在
		if err := os.MkdirAll("static/uploads", 0755); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "创建目录失败"})
			return
		}

		//保存文件
		dst := filepath.Join("static/uploads", fileName)
		if err := c.SaveUploadedFile(file, dst); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		}

		//保存文件信息到数据库
		mediaType := "image"
		if strings.Contains(file.Header.Get("Content-Type"), "video") {
			mediaType = "video"
		}
		media := model.Media{
			Path: dst,
			Type: mediaType,
			// PostID:post.ID,会自动添加，无需手动管理
		}
		medias = append(medias, media)

	}
	post.Media = medias
	if err := database.DB.Create(&post).Error; err != nil { //因为post表和media表通过外键postID关联，所以创建post时，会自动创建media并填充media表中的postID
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create post"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Post created successfully", "postID": post.ID})
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
	newMedias := form.File["newMedias"]        //新上传的媒体

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
			Path:   dst,
			Type:   mediaType,
			PostID: post.ID,
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
	userID := c.GetUint("userID") //点赞用户的ID
	postID := c.Param("id")
	// database.DB.Model(&model.Post{}).Where("id = ? AND user_id = ?", postID, userID).Update("like_count", gorm.Expr("like_count + 1"))
	var post model.Post
	result := database.DB.Where("id=? AND user_id=?", postID, userID).First(&post)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}
	post.LikeCount = append(post.LikeCount, userID)
	database.DB.Save(&post)
	c.JSON(http.StatusOK, gin.H{"message": "Like count added successfully"})
}

func SubLikeCount(c *gin.Context) {
	userID := c.GetUint("userID")
	postID := c.Param("id")
	var post model.Post
	result := database.DB.Where("id=? AND user_id=?", postID, userID).First(&post)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	// 手动过滤掉当前用户ID
	newLikeCount := make([]uint, 0)
	for _, id := range post.LikeCount {
		if id != userID {
			newLikeCount = append(newLikeCount, id)
		}
	}
	post.LikeCount = newLikeCount

	database.DB.Save(&post)
	c.JSON(http.StatusOK, gin.H{"message": "Like count subtracted successfully"})
}

// 用于主页面的渲染，只返回media中的一个
func GetPosts(c *gin.Context) {
	lastID := c.DefaultQuery("last_id", "0")
	pageSize := 10

	var posts []model.Post
	query := database.DB.
		Preload("Media", func(db *gorm.DB) *gorm.DB {
			return db.Order("id DESC").Limit(1) // 按ID降序取第一条媒体
		}).
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
	if err := database.DB.Preload("Media").Preload("User").Preload("Comment").Preload("Comment.User").First(&post, postID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
	}
	c.JSON(http.StatusOK, gin.H{"post": post})
}

func AddComment(c *gin.Context) {
	var receiveComment struct {
		PostID  uint   `json:"post_id"`
		Content string `json:"content"`
	}
	if err := c.ShouldBindJSON(&receiveComment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetUint("userID")

	comment := model.Comment{
		PostID:  receiveComment.PostID,
		Content: receiveComment.Content,
		UserID:  userID,
	}

	// 保存评论
	if err := database.DB.Create(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add comment"})
		return
	}

	// 获取帖子作者ID
	var post model.Post
	if err := database.DB.Where("id = ?", receiveComment.PostID).Select("user_id").First(&post).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}
	receiverId := post.UserID

	// 如果评论者不是帖子作者，则创建通知
	if comment.UserID != receiverId {
		notification := model.Notification{
			SenderID:   comment.UserID,
			ReceiverID: receiverId,
			PostID:     receiveComment.PostID,
			Type:       "comment",
		}
		database.DB.Create(&notification)

		// SendNotificationToUser(userID, []byte(fmt.Sprintf("你收到了来自%d的评论", comment.UserID)))
		message := `{"type": "notification", "count": 1}`
		SendNotificationToUser(receiverId, []byte(message))
	}
	c.JSON(http.StatusOK, gin.H{"message": "Comment added successfully"})
}

func DeleteComment(c *gin.Context) {
	userID := c.GetUint("userID")
	postID := c.Param("postID")

	if err := database.DB.Where("postID=? AND user_id=?", postID, userID).Delete(&model.Comment{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete comment"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Comment deleted successfully"})
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
