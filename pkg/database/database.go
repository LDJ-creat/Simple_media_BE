package database

import (
	"fmt"

	"github.com/media/internal/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDatabase() error {
	db, err := gorm.Open(mysql.Open("root:268968&&ABc@tcp(localhost:3306)/media?charset=utf8mb4&parseTime=True&loc=Local"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	//如果不存在先创建
	err = db.Exec("CREATE DATABASE IF NOT EXISTS media").Error
	if err != nil {
		return fmt.Errorf("创建数据库失败: %v", err)
	}

	//自动迁移
	err = db.AutoMigrate(&model.User{}, &model.Post{}, &model.Media{})
	if err != nil {
		return fmt.Errorf("自动迁移失败: %v", err)
	}
	DB = db
	return nil
}
