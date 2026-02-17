package handlers

import (
	"go_mysql/database"
	"go_mysql/models"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func GetRentals(c echo.Context) error {
    var equipments []models.RentalItem
    
    // 全件取得
    if err := database.DB.Find(&equipments).Error; err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{"message": "データ取得失敗"})
    }

    return c.JSON(http.StatusOK, equipments)
}

func ToggleRental(c echo.Context) error {
    // パラメータから操作したい備品のIDを取得
    itemID := c.Param("id")

    // ユーザー認証情報から操作者の名前を取得
    userToken := c.Get("user").(*jwt.Token)
    claims := userToken.Claims.(jwt.MapClaims)
    username := claims["name"].(string)

    // 備品の現在の状態を確認
    var item models.RentalItem
    if err := database.DB.First(&item, itemID).Error; err != nil {
        return c.JSON(http.StatusNotFound, map[string]string{"message": "備品がありません"})
    }


    // すでに誰かが借りている場合
    if item.RenterName != "" {
        // 借りているのが自分なら返却する
        if item.RenterName == username {
            database.DB.Model(&item).Update("renter_name", "") // 空文字で返却
            return c.JSON(http.StatusOK, map[string]string{"message": "返却しました"})
        }
        // 借りているのが他人ならエラー
        return c.JSON(http.StatusConflict, map[string]string{"message": "他の人が使用中です"})
    }

    // 誰も借りていない場合は、自分が借りる
    database.DB.Model(&item).Update("renter_name", username)
    
    return c.JSON(http.StatusOK, map[string]string{"message": "借用しました"})
}