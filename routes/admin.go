package routes

import "github.com/gofiber/fiber/v2"



func TokensList(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "welcome to the jungle"})
}

func TokenCreate(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "welcome to the jungle"})
}

func TokenEdit(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "welcome to the jungle"})
}

func TokenDelete(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "welcome to the jungle"})
}
