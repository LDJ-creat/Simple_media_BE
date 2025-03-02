package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/media/pkg/jwt"
)

func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		fmt.Printf("Received Authorization header: %s\n", token)

		if token == "" {
			fmt.Println("No Authorization header found")
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "未授权，请先登录",
			})
			c.Abort()
			return
		}

		// 检查并去掉Bearer前缀
		const bearerPrefix = "Bearer "
		if !strings.HasPrefix(token, bearerPrefix) {
			fmt.Printf("Invalid token format, missing Bearer prefix: %s\n", token)
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "无效的token格式",
			})
			c.Abort()
			return
		}
		token = token[len(bearerPrefix):]
		fmt.Printf("Token after removing Bearer prefix: %s\n", token)

		claims, err := jwt.ParseToken(token)
		if err != nil {
			fmt.Printf("Failed to parse token: %v\n", err)
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "token无效或已过期",
			})
			c.Abort()
			return
		}

		fmt.Printf("Token validated successfully for user ID: %d\n", claims.UserID)
		c.Set("userID", claims.UserID)
		c.Next()
	}
}

// AuthMiddleware 中间件
func WebSocketAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Query("token") // 从 URL 参数获取 token
		if token == "" {
			token = c.GetHeader("Authorization")
			// 移除 "Bearer " 前缀
			if len(token) > 7 && token[:7] == "Bearer " {
				token = token[7:]
			}
		}

		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
			c.Abort()
			return
		}

		// 验证 token 并获取用户 ID
		claims, err := jwt.ParseToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "token无效"})
			c.Abort()
			return
		}

		c.Set("userID", claims.UserID)
		c.Next()
	}
}
