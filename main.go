package main

import (
	"go_mysql/database"
	"go_mysql/handlers"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
    database.InitDB()
    e := echo.New()
    e.Use(middleware.Logger(), middleware.Recover(), middleware.CORS())

    // --- 公開ルート ---
    e.POST("/login", handlers.Login)
    e.POST("/members", handlers.CreateMember)
    e.GET("/rentals", handlers.GetRentals) 

    // --- 制限ルート ---
    r := e.Group("/auth")
    r.Use(echojwt.WithConfig(echojwt.Config{
        SigningKey: []byte("secret"),
    }))

    r.GET("/member", handlers.GetMe)
    r.POST("/rentals/:id/toggle", handlers.ToggleRental) 

    e.Logger.Fatal(e.Start(":8080"))
}