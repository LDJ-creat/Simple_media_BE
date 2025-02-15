package v1

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/media/internal/model"
	"github.com/media/pkg/database"
)

func CreatePost(c *gin.Context) {
	userID := c.GetUint("userID")
	// var newPost newPost
	// if err:=c.ShouldBindJSON(&newPost);err!=nil{
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	// }
	body := c.PostForm("content")
	imagesPath := []model.Image{}
	videosPath := []model.Video{}
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get multipart form"})
		return
	}
	images := form.File["images"]
	videos := form.File["videos"]
	for _, image := range images {
		var imageItem model.Image
		imageItem.Path = MediaUpload(c, image, "image")
		// 获取文件id
		src, err := image.Open()
		if err != nil {
			c.JSON(500, gin.H{"error": fmt.Sprintf("打开文件 %s 失败: %v", image.Filename, err)})
			return
		}
		defer src.Close()

		// 这里需要解析前端传递的额外信息（id、type 等）
		var mediaInfo struct {
			ID   int    `json:"id"`
			Type string `json:"type"`
		}
		err = json.Unmarshal([]byte(image.Filename), &mediaInfo)
		if err != nil {
			c.JSON(500, gin.H{"error": fmt.Sprintf("解析文件信息 %s 失败: %v", image.Filename, err)})
			return
		}
		imageItem.ID = uint(mediaInfo.ID)
		imagesPath = append(imagesPath, imageItem)
	}
	for _, video := range videos {
		var videoItem model.Video
		videoItem.Path = MediaUpload(c, video, "video")

		// 获取文件id
		src, err := video.Open()
		if err != nil {
			c.JSON(500, gin.H{"error": fmt.Sprintf("打开文件 %s 失败: %v", video.Filename, err)})
			return
		}
		defer src.Close()
		//解析出id
		var mediaInfo struct {
			ID   int    `json:"id"`
			Type string `json:"type"`
		}
		err = json.Unmarshal([]byte(video.Filename), &mediaInfo)
		if err != nil {
			c.JSON(500, gin.H{"error": fmt.Sprintf("解析文件信息 %s 失败: %v", video.Filename, err)})
			return
		}
		videoItem.ID = uint(mediaInfo.ID)
		videosPath = append(videosPath, videoItem)
	}

	post := model.Post{
		UserID:  userID,
		Content: body,
		Image:   imagesPath,
		Video:   videosPath,
	}

	if err := database.DB.Create(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create post"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Post created successfully"})

}

func ModifyPost(c *gin.Context) {
	userID := c.GetUint("userID")
	postID := c.Param("id")
	body := c.PostForm("content")
	imagesPath := []model.Image{}
	oddImages := []model.Image{}
	//将数据库中的图片信息拿出来
	database.DB.Where("id = ? AND user_id = ?", postID, userID).Find(&oddImages)
	fmt.Println(oddImages)
	videosPath := []model.Video{}
	oddVideos := []model.Video{}
	//将数据库中的视频信息拿出来
	database.DB.Where("id = ? AND user_id = ?", postID, userID).Find(&oddVideos)
	fmt.Println(oddVideos)
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get multipart form"})
		return
	}
	images := form.File["images"]
	videos := form.File["videos"]
	//处理图片
	for _, image := range images {
		var imageItem model.Image
		// 获取文件id
		src, err := image.Open()
		if err != nil {
			c.JSON(500, gin.H{"error": fmt.Sprintf("打开文件 %s 失败: %v", image.Filename, err)})
			return
		}
		defer src.Close()

		// 这里需要解析前端传递的额外信息（id、type 等）
		var mediaInfo struct {
			ID   int    `json:"id"`
			Type string `json:"type"`
		}
		err = json.Unmarshal([]byte(image.Filename), &mediaInfo)
		if err != nil {
			c.JSON(500, gin.H{"error": fmt.Sprintf("解析文件信息 %s 失败: %v", image.Filename, err)})
			return
		}
		imageItem.ID = uint(mediaInfo.ID)
		isExist := false
		for _, oddImage := range oddImages {
			if oddImage.ID == imageItem.ID {
				isExist = true
				imagesPath = append(imagesPath, oddImage)
				break
			}
		}
		if !isExist {
			imageItem.Path = MediaUpload(c, image, "image")
			imagesPath = append(imagesPath, imageItem)
		}

	}
	//处理视频
	for _, video := range videos {
		var videoItem model.Video
		// 获取文件id
		src, err := video.Open()
		if err != nil {
			c.JSON(500, gin.H{"error": fmt.Sprintf("打开文件 %s 失败: %v", video.Filename, err)})
			return
		}
		defer src.Close()
		//解析出id
		var mediaInfo struct {
			ID   int    `json:"id"`
			Type string `json:"type"`
		}
		err = json.Unmarshal([]byte(video.Filename), &mediaInfo)
		if err != nil {
			c.JSON(500, gin.H{"error": fmt.Sprintf("解析文件信息 %s 失败: %v", video.Filename, err)})
			return
		}
		videoItem.ID = uint(mediaInfo.ID)
		isExist := false
		for _, oldVideo := range oddVideos {
			if oldVideo.ID == videoItem.ID {
				isExist = true
				videosPath = append(videosPath, oldVideo)
				break
			}
		}
		if !isExist {
			videoItem.Path = MediaUpload(c, video, "video")
			videosPath = append(videosPath, videoItem)
		}
	}

	post := model.Post{
		UserID:  userID,
		Content: body,
		Image:   imagesPath,
		Video:   videosPath,
	}

	database.DB.Model(&model.Post{}).Where("id = ? AND user_id = ?", postID, userID).Updates(post)
	c.JSON(http.StatusOK, gin.H{"message": "Post updated successfully"})
}

func DeletePost(c *gin.Context) {
	userID := c.GetUint("userID")
	postID := c.Param("id")
	database.DB.Where("id = ? AND user_id = ?", postID, userID).Delete(&model.Post{})
	c.JSON(http.StatusOK, gin.H{"message": "Post deleted successfully"})
}
