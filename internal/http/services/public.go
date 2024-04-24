package services

import (
	"errors"
	"fmt"
	"log"
	"time"

	db "github.com/diegom0ta/gofiber-study/internal/database"
	"github.com/diegom0ta/gofiber-study/internal/models"
	"github.com/diegom0ta/gofiber-study/internal/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Name                 string `json:"name"`
	Document             string `json:"document"`
	Email                string `json:"email"`
	Celular              string `json:"celular"`
	Password             string `json:"password"`
	PasswordConfirmation string `json:"password_confirmation"`
}

type UserLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

var ErrPasswdNotEqual = errors.New("passwords are not equal")

// Register creates a new user in database
func Register(c *fiber.Ctx) error {
	newUser := new(User)

	if err := c.BodyParser(&newUser); err != nil {
		return fmt.Errorf("error: %v", err)
	}

	var existingUser models.User

	existingUser.Email = newUser.Email
	existingUser.Document = newUser.Document

	if err := db.Db.Where("email = ?", existingUser.Email).Or("document = ?", existingUser.Document).First(&existingUser).Error; err == nil {
		err = c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "email or document already exists"})
		if err != nil {
			log.Printf("error in JSON statement: %v", err)
		}
	}

	passwd, err := validatePasswd(newUser.Password, newUser.PasswordConfirmation)
	if err != nil {
		log.Printf("password not valid: %v", err)
	}

	uuid := uuid.New()

	user := models.User{
		Id:           uuid.String(),
		Name:         newUser.Name,
		Document:     newUser.Document,
		Email:        newUser.Email,
		Phone:        newUser.Celular,
		PasswordHash: passwd,
		CreatedAt:    time.Now().UTC(),
	}

	var success bool

	if result := db.Db.Create(&user); result.Error != nil {
		success = false

		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": success})
	} else {
		success = true
	}

	return c.JSON(fiber.Map{"success": success})
}

func hashPasswd(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("could not hash password: %v", err)
	}

	return string(hashedPassword), nil
}

func validatePasswd(password, passwdConfirmation string) (string, error) {
	var hashedPasswd string

	if password == passwdConfirmation {
		hash, err := hashPasswd(password)
		if err != nil {
			return "", err
		}

		hashedPasswd = hash
	} else {
		return "", ErrPasswdNotEqual
	}

	return hashedPasswd, nil
}

// Login authenticates a user
func Login(c *fiber.Ctx) error {
	var (
		existingUser models.User
		success      bool
	)

	login := new(UserLogin)

	if err := c.BodyParser(&login); err != nil {
		success = false

		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Error parsing JSON",
			"code":    "error.internal",
			"success": success,
		})
	}

	existingUser.Email = login.Email

	if err := db.Db.Where("email = ?", existingUser.Email).First(&existingUser).Error; err != nil {
		success = false

		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "Invalid email or password",
			"code":    "error.invalid-credentials",
			"success": success,
		})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(existingUser.PasswordHash), []byte(login.Password)); err != nil {
		success = false

		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "Invalid email or password",
			"code":    "error.invalid-credentials",
			"success": success,
		})
	}

	token, err := utils.GenerateToken(existingUser.Email)
	if err != nil {
		success = false

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error generating token",
			"code":    "error.internal",
			"success": success,
		})
	}

	success = true

	return c.JSON(fiber.Map{
		"success": success,
		"token":   token,
		"data": map[string]interface{}{
			"name": existingUser.Name,
		},
	})
}
