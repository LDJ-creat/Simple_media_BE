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
		auth := api.Group("")
		auth.Use(middleware.JWT())
		{
			api.POST("/update", v1.UpdateUser)
			api.POST("/update-password", v1.UpdatePassword)
			// api.POST("/forgot-password", v1.ForgotPassword)
			// api.POST("/reset-password", v1.ResetPassword)
		}
	}
}
