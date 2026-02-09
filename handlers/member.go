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
	database.DB.Create(&m)
	return c.JSON(http.StatusCreated, m)
}

func GetMembers(c echo.Context) error {
	var members []models.Member
	database.DB.Find(&members)
	return c.JSON(http.StatusOK, members)
}