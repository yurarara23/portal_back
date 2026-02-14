package handlers

import (
	"net/http"
	"go_mysql/database"
	"go_mysql/models"
	"github.com/labstack/echo/v4"
)

func CreateMember(c echo.Context) error {
	m := new(models.Member)
	if err := c.Bind(m); err != nil {
		return err
	}

	if err := m.HashPassword(); err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{"message": "ハッシュ化に失敗しました"})
    }

	database.DB.Create(&m)

	return c.JSON(http.StatusCreated, m)
}

func GetMe(c echo.Context) error {
    // 1. トークンからユーザー情報（Claims）を取り出す
    user := c.Get("user").(*jwt.Token)
    claims := user.Claims.(jwt.MapClaims)
    
    // Login時に "name" というキーで保存した値を取り出す
    username := claims["name"].(string)

    // 2. そのユーザー名でDBを検索
    var member models.Member
    if err := database.DB.Where("username = ?", username).First(&member).Error; err != nil {
        return c.JSON(http.StatusNotFound, map[string]string{"message": "ユーザーが見つかりません"})
    }

    // パスワードは隠して返す
    member.Password = ""
    return c.JSON(http.StatusOK, member)
}