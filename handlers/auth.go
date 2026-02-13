package handlers

import (
	"net/http"
	"go_mysql/database"
	"go_mysql/models"
	"github.com/labstack/echo/v4"
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

	result := database.DB.Where("username = ? AND password = ?", req.Username, req.Password).First(&member)

	if result.Error != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "ユーザー名またはパスワードが違います"})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "ログイン成功",
		"user":    member,
	})
}