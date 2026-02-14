package handlers

import (
	"net/http"
	"go_mysql/database"
	"go_mysql/models"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

func Login(c echo.Context) error {
	type LoginRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	req := new(LoginRequest)
	if err := c.Bind(req); err != nil {
		return err
	}

	var member models.Member

	if err := database.DB.Where("username = ?", req.Username).First(&member).Error; err != nil {
        return echo.ErrUnauthorized
    }

	err := bcrypt.CompareHashAndPassword([]byte(member.Password), []byte(req.Password))
    if err != nil {
        return echo.ErrUnauthorized
    }

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "ログイン成功",
		"user":    member,
	})
}