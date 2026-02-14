package handlers

import (
	"net/http"
	"go_mysql/database"
	"go_mysql/models"
	"time"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"github.com/golang-jwt/jwt/v5"
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

	claims := &jwt.MapClaims{
        "name":  member.Username,
        "admin": true, // 権限など
        "exp":   time.Now().Add(time.Hour * 1).Unix(), // 有効期限（1時間）
    }

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

    // 署名（秘密鍵は本来環境変数などで管理する）
    t, err := token.SignedString([]byte("secret")) 
    if err != nil {
        return err
    }

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "ログイン成功",
		"token": t, //複数行なら,ないとエラーになる
	})
}