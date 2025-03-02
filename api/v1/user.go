package v1

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/media/internal/model"
	"github.com/media/pkg/database"
	"github.com/media/pkg/email"
	"github.com/media/pkg/jwt"
	"golang.org/x/crypto/bcrypt"
)

type UserRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required"`
}

type UpdateUserRequest struct {
	Username  string `json:"username"`
	Email     string `json:"email"`
	Avatar    string `json:"avatar"`
	Signature string `json:"signature"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required"`
}

type ResetPasswordRequest struct {
	Email       string `json:"email" binding:"required"`
	Code        string `json:"code" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}

type UpdatePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}

// 注册
func Register(c *gin.Context) {
	// 打印原始请求数据
	body, _ := io.ReadAll(c.Request.Body)
	log.Printf("收到注册请求，原始数据: %s", string(body))
	c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

	var req UserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("解析请求数据失败: %v", err)
		log.Printf("请求头 Content-Type: %s", c.GetHeader("Content-Type"))
		c.JSON(http.StatusBadRequest, gin.H{
			"error":    fmt.Sprintf("无效的请求数据: %v", err),
			"received": string(body),
		})
		return
	}

	log.Printf("解析后的请求数据: %+v", req)

	// 检查用户是否存在
	var existUser model.User
	result := database.DB.Where("username=? OR email=?", req.Username, req.Email).First(&existUser)
	log.Printf("检查用户是否存在: %v", result.Error)

	if existUser.ID != 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户已存在"})
		return
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("密码加密失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "密码加密失败"})
		return
	}

	user := model.User{
		Username: req.Username,
		Email:    req.Email,
		Password: string(hashedPassword),
		// CodeTime 字段不设置，将默认为 null
	}

	if err := database.DB.Create(&user).Error; err != nil {
		log.Printf("创建用户失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建用户失败"})
		return
	}

	// 生成token
	token, err := jwt.GenerateToken(user.ID)
	if err != nil {
		log.Printf("生成token失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成token失败"})
		return
	}
	log.Printf("生成token成功: %v", token)
	log.Printf("用户创建成功: %v", user)

	c.JSON(http.StatusOK, gin.H{
		"message":  "用户创建成功",
		"token":    token,
		"userData": user,
	})
}

// 登录
func Login(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//检查用户是否存在
	var user model.User
	if err := database.DB.Where("username=? OR email=?", req.Username, req.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户不存在"})
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户名或密码错误"})
		return
	}

	//生成token
	token, err := jwt.GenerateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成token失败"})
		return
	}
	//返回token和用户信息
	c.JSON(http.StatusOK, gin.H{"message": "登录成功", "token": token, "userData": user})
}

// 修改用户信息
func UpdateUser(c *gin.Context) {
	var req UpdateUserRequest
	form, _ := c.MultipartForm()
	req.Username = form.Value["username"][0]
	req.Email = form.Value["email"][0]
	req.Signature = form.Value["signature"][0]
	file := form.File["avatar"][0]
	if file != nil {
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
		dst := filepath.Join("static/avatars", fileName)
		if err := c.SaveUploadedFile(file, dst); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		}

		req.Avatar = dst

	}

	//更新到数据库
	if err := database.DB.Model(&model.User{}).Where("id=?", c.GetUint("userID")).Updates(req).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新用户信息失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "用户信息更新成功"})

}

// 修改密码
func UpdatePassword(c *gin.Context) {
	var req UpdatePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetUint("userID")
	var user model.User
	if err := database.DB.Where("id=?", userID).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户不存在"})
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.OldPassword)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "旧密码错误"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "密码加密失败"})
		return
	}
	database.DB.Model(&model.User{}).Where("id=?", userID).Update("password", string(hashedPassword))
	c.JSON(http.StatusOK, gin.H{"message": "密码更新成功"})
}

// 忘记密码
func ForgotPassword(c *gin.Context) {
	var req ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 检查邮箱是否存在
	var user model.User
	if err := database.DB.Where("email=?", req.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "邮箱不存在"})
		return
	}

	code, err := email.SendCode(req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "发送验证码失败"})
		return
	}
	database.DB.Model(&model.User{}).Where("email=?", req.Email).Update("code", code)
	//设置验证码过期时间为5分钟
	database.DB.Model(&model.User{}).Where("email=?", req.Email).Update("code_time", time.Now().Add(5*time.Minute))
	c.JSON(http.StatusOK, gin.H{"message": "验证码发送成功"})
}

// 重置密码
func ResetPassword(c *gin.Context) {
	var req ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 查找对应用户
	var user model.User
	if err := database.DB.Where("email=?", req.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "不存在该邮箱"})
		return
	}

	// 检查验证码是否过期
	if user.CodeTime.Before(time.Now()) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "验证码已过期"})
		return
	}

	// 检查验证码是否正确
	if user.Code != req.Code {
		c.JSON(http.StatusBadRequest, gin.H{"error": "验证码错误"})
		return
	}

	// 更新密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "密码加密失败"})
		return
	}
	database.DB.Model(&model.User{}).Where("email=?", req.Email).Update("password", string(hashedPassword))
	c.JSON(http.StatusOK, gin.H{"message": "密码重置成功"})
}

// 获取用户信息
func GetUserInfo(c *gin.Context) {
	userID := c.GetUint("userID")
	fmt.Println("userID:", userID)
	var user model.User
	if err := database.DB.Where("id=?", userID).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户不存在"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"userData": user})
}
