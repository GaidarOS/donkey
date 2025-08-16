package routes

import (
	"donkey/auth"
	"donkey/config"
	"donkey/database"
	"log/slog"

	"github.com/gofiber/fiber/v2"
)

func Register(c *fiber.Ctx) error {
	var req config.AuthRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	user := config.User{
		UserName:     req.Username,
		PasswordHash: auth.GeneratePassword(req.Password),
	}
	res := database.DB.Create(&user)
	if res.Error != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": res.Error.Error(),
		})
	}
	return c.Status(201).JSON(fiber.Map{
		"message": "user created",
	})
}

func Login(c *fiber.Ctx) error {
	var req config.AuthRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	var user config.User
	res := database.DB.Where("username = ?", req.Username).First(&user)
	if res.Error != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "user not found",
		})
	}

	if !auth.ComparePassword(user.PasswordHash, req.Password) {
		return c.Status(400).JSON(fiber.Map{
			"message": "incorrect password",
		})
	}

	token, err := auth.GenerateToken(user.ID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	valid, err := auth.VerifyToken(token)
	slog.Info("Token Validattion", slog.Any("verify", valid))
	if err != nil {
		slog.Error("Invalid token", slog.Any("err", err))
		return c.Status(500).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"token": token,
	})
}
