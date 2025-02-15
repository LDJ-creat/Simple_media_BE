package v1

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
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
	var req UserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//通过邮箱和用户名检查用户是否存在
	var existUser model.User
	if err := database.DB.Where("username=?", req.Username).First(&existUser).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户名已存在"})
		return
	}

	//加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "密码加密失败"})
		return
	}
	user := model.User{
		Username: req.Username,
		Email:    req.Email,
		Password: string(hashedPassword),
	}

	//创建用户到数据库
	if err := database.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建用户失败"})
		return
	}

	//生成token
	token, err := jwt.GenerateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成token失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "用户创建成功", "token": token})

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
	c.JSON(http.StatusOK, gin.H{"message": "登录成功", "token": token})
}

// 处理头像
func UploadAvatar(c *gin.Context) {
	// 获取上传的文件
	file, err := c.FormFile("avatar")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请选择文件"})
		return
	}

	// 验证文件大小（例如限制为2MB）
	if file.Size > 2<<20 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "文件大小不能超过2MB"})
		return
	}

	// 验证文件类型
	ext := strings.ToLower(path.Ext(file.Filename))
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".gif" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "只支持jpg、jpeg、png、gif格式"})
		return
	}

	// 生成随机文件名
	fileName := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	avatarPath := fmt.Sprintf("static/avatars/%s", fileName)

	// 确保目录存在
	if err := os.MkdirAll("static/avatars", 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建目录失败"})
		return
	}

	// 保存文件
	if err := c.SaveUploadedFile(file, avatarPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "保存文件失败"})
		return
	}

	// 更新数据库中的头像URL
	avatarURL := "/" + avatarPath // 存储相对路径
	if err := database.DB.Model(&model.User{}).Where("id = ?", c.GetUint("userID")).Update("avatar", avatarURL).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新头像失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "头像上传成功",
		"url":     avatarURL,
	})

}

// 修改用户信息
func UpdateUser(c *gin.Context) {
	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//更新到数据库
	if err := database.DB.Model(&model.User{}).Where("id=?", c.GetUint("userID")).Updates(req).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新用户信息失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "用户信息更新成功"})
	UploadAvatar(c)
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

// 上传头像
