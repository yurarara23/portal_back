package main

import (
    "go_mysql/database"
    "go_mysql/handlers"
    "github.com/labstack/echo/v4"
    "github.com/labstack/echo/v4/middleware"
)

func main() {

    database.InitDB() // DB初期化

    e := echo.New()
    e.Use(middleware.Logger(), middleware.Recover(), middleware.CORS())

    e.POST("/login", handlers.Login)
    e.POST("/members", handlers.CreateMember)

    e.Logger.Fatal(e.Start(":8080"))
}