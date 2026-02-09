package handlers

import (
	"net/http"
	"go_mysql/database"
	"go_mysql/models"
	"github.com/labstack/echo/v4"
)

func ToggleRental(c echo.Context) error {
	type ToggleRequest struct {
		MemberID uint   `json:"member_id"`
		ItemName string `json:"item_name"`
	}
	req := new(ToggleRequest)
	if err := c.Bind(req); err != nil {
		return err
	}

	var rental models.Rental
	result := database.DB.Where("member_id = ? AND item_name = ?", req.MemberID, req.ItemName).Last(&rental)

	if result.Error != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "記録が見つかりません"})
	}

	rental.IsRented = !rental.IsRented
	database.DB.Save(&rental)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"isRented":  rental.IsRented,
		"item_name": rental.ItemName,
	})
}

func GetRentalStatus(c echo.Context) error {
	itemName := c.QueryParam("item_name")
	var rental models.Rental
	result := database.DB.Where("item_name = ?", itemName).Last(&rental)

	if result.Error != nil {
		return c.JSON(http.StatusOK, map[string]interface{}{"isRented": false, "username": ""})
	}

	var member models.Member
	database.DB.First(&member, rental.MemberID)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"isRented": rental.IsRented,
		"username": member.UserName,
	})
}