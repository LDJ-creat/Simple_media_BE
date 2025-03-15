package main

import (
	"github.com/gin-gonic/gin"
	"github.com/media/pkg/database"
	"github.com/media/router"
)

func main() {
	r := gin.Default()
	r.Static("/static", "./static") // 配置静态文件服务
	// 初始化数据库连接
	if err := database.InitDatabase(); err != nil {
		panic("数据库连接失败: " + err.Error())
	}

	// 初始化路由
	router.InitRouter(r)

	// 启动服务器，监听所有地址
	r.Run("0.0.0.0:8081")

}
