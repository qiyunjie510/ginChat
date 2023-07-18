package main

import (
	"ginchat/models"
	"ginchat/utils"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Code  string
	Price uint
}

func main() {
	utils.InitConfig()
	mysqlStr := viper.GetString("mysql.dns")
	db, err := gorm.Open(mysql.Open(mysqlStr), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// 迁移 schema
	//db.AutoMigrate(&models.Message{})
	//db.AutoMigrate(&models.Contact{})
	//db.AutoMigrate(&models.GroupBasic{})
	//db.AutoMigrate(&models.UserBasic{})
	db.AutoMigrate(&models.Community{})

	// Create
	// user := &models.UserBasic{}
	// user.Name = "盛转"
	// db.Create(user)

	// // Read
	// fmt.Println(db.First(user, 1))

	// // Update - 将 product 的 price 更新为 200
	// db.Model(&user).Update("PassWord", 1234)

}
