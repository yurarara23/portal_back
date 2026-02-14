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

    // 公開ルート
    e.POST("/login", handlers.Login)
    e.POST("/members", handlers.CreateMember)

    // 制限ルート（グループ化すると管理しやすい）
    r := e.Group("/auth")
    r.Use(echojwt.WithConfig(echojwt.Config{
        SigningKey: []byte("secret"), // Login時と同じ鍵を使う
    }))

    // /auth/members は有効なトークンがないとアクセス不可
    r.GET("/member", handlers.GetMe)

    e.Logger.Fatal(e.Start(":8080"))
}