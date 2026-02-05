package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Item struct {
	gorm.Model
	ItemName string `json:"itemName"`
	IsRented bool `json:"isRented"`
	RenterName string `json:"renterName`
}

func main() {
	dsn := "root:password@tcp(127.0.0.1:3306)/rental_db?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("データベース接続に失敗しました")
	}

	db.AutoMigrate(&Item{})

	e.echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, Go Backend!")
	})

	fmt.Println("Server started at :8080")
	e.Logger.Fatal(e.Start(":8080"))
}