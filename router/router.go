package router

import (
	"github.com/gin-gonic/gin"
	v1 "github.com/media/api/v1"
	"github.com/media/internal/middleware"
)

func InitRouter(r *gin.Engine) {
	api := r.Group("/api/v1")
	{
		api.POST("/register", v1.Register)
		api.POST("/login", v1.Login)
		api.GET("/ws", v1.HandleWebSocket)

		auth := api.Group("")
		auth.Use(middleware.JWT())
		{
			auth.GET("/getUserInfo", v1.GetUserInfo)
			auth.GET("/getMyPosts", v1.GetMyPosts)
			auth.PUT("/updateUser", v1.UpdateUser)
			auth.POST("/update-password", v1.UpdatePassword)
			auth.POST("/forgot-password", v1.ForgotPassword)
			auth.POST("/reset-password", v1.ResetPassword)
			auth.POST("/post", v1.CreatePost)
			auth.PUT("/post/:id", v1.UpdatePost)
			auth.DELETE("/post/:id", v1.DeletePost)
			auth.POST("/addLike/:id", v1.AddLikeCount)
			auth.PUT("/subLike/:id", v1.SubLikeCount)
			auth.GET("/getPosts", v1.GetPosts)
			auth.GET("/postDetails/:id", v1.GetPostByID)
			auth.POST("/addComment", v1.AddComment)
			auth.DELETE("/deleteComment/:id", v1.DeleteComment)
			auth.GET("/notifications/count", v1.GetNotificationsCount)
			auth.GET("/notifications", v1.GetNotifications)
		}
	}

}
