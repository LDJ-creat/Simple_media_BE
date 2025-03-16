// package database

// import (
// 	"fmt"

// 	"github.com/media/internal/model"
// 	"gorm.io/driver/mysql"
// 	"gorm.io/gorm"
// )

// var DB *gorm.DB

// func InitDatabase() error {
// 	db, err := gorm.Open(mysql.Open("root:268968&&ABc@tcp(localhost:3306)/media?charset=utf8mb4&parseTime=True&loc=Local"), &gorm.Config{})
// 	if err != nil {
// 		fmt.Println("连接数据库失败: " + err.Error())
// 	}

// 	//如果不存在先创建
// 	err = db.Exec("CREATE DATABASE IF NOT EXISTS media").Error
// 	if err != nil {
// 		return fmt.Errorf("创建数据库失败: %v", err)
// 	}

//		//自动迁移
//		err = db.AutoMigrate(&model.User{}, &model.Post{}, &model.Media{}, &model.Comment{})
//		if err != nil {
//			return fmt.Errorf("自动迁移失败: %v", err)
//		}
//		DB = db
//		return nil
//	}
package database

import (
	"database/sql"
	"fmt"

	"github.com/media/internal/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDatabase() error {
	// 先创建数据库
	createDB, err := sql.Open("mysql", "root:268968ABc@tcp(localhost:3306)/")
	if err != nil {
		return fmt.Errorf("打开数据库连接失败: %v", err)
	}
	defer createDB.Close()

	_, err = createDB.Exec("CREATE DATABASE IF NOT EXISTS media")
	if err != nil {
		return fmt.Errorf("创建数据库失败: %v", err)
	}

	// 再建立gorm连接
	db, err := gorm.Open(mysql.Open("root:268968ABc@tcp(localhost:3306)/media?charset=utf8mb4&parseTime=True&loc=Local"), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("连接数据库失败: %v", err)
	}

	// 自动迁移
	err = db.AutoMigrate(&model.User{}, &model.Post{}, &model.Media{}, &model.Comment{}, &model.Notification{})
	if err != nil {
		return fmt.Errorf("自动迁移失败: %v", err)
	}
	DB = db
	return nil
}
