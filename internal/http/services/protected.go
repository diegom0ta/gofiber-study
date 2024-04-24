package services

import (
	"strconv"

	db "github.com/diegom0ta/gofiber-study/internal/database"
	"github.com/diegom0ta/gofiber-study/internal/models"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

const pageSize = 10

// GetCredits returns a list of credits
func GetUsers(c *fiber.Ctx) error {
	var (
		users    []models.User
		response models.UserResponse
	)

	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userId := claims["id"]

	page := c.Params("page")
	pageInt, _ := strconv.Atoi(page)

	offset := (pageInt - 1) * pageSize

	var count int64
	db.Db.Model(&models.User{}).Where("id = ?", userId).Count(&count)

	if err := db.Db.Select("name, document, email, phone").Where("id = ?", userId).Offset(offset).Limit(pageSize).Find(&users).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error fetching credits",
			"code":    "error.internal",
		})
	}

	response.Page = pageInt
	response.PageSize = pageSize
	response.Total = int(count)
	response.Data = users

	return c.JSON(response)
}
