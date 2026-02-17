package database

import (
	"go_mysql/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	dsn := "root:password@tcp(db:3306)/nu_db?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("データベース接続に失敗しました")
	}

	DB.AutoMigrate(&models.Member{}, &models.RentalItem{})
}