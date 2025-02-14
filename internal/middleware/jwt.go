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
