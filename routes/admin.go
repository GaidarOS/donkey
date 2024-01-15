package routes

import (
	"log/slog"
	"receipt_store/config"

	"github.com/gofiber/fiber/v2"
)

func TokensList(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "List of user tokens",
		"data": config.AppConf.Tokens})
}

func TokenCreate(c *fiber.Ctx) error {

	tkn := config.Token{}
	// Parse the body
	if err := c.BodyParser(&tkn); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	// check if user or token already exist
	// fail if either does
	foundToken, err := config.AppConf.FindStructByToken(tkn.Value)
	if err != nil {
		slog.Debug("No matching token found", err)
	}
	foundName, err := config.AppConf.FindStructByName(tkn.UserName)
	if err != nil {
		slog.Debug("No matching username found", err)
	}
	if foundName != nil || foundToken != nil {
		slog.Error("Found an existing username or token. Try updating instead.")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Found an existing username or token. Try updating instead"})
	}

	// append the new token to the list
	config.AppConf.Tokens = append(config.AppConf.Tokens, tkn)

	// write the new config to file for persisten storage
	if err := config.AppConf.WriteToConf(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error while trying to store the new config. Please try again later"})
	}

	return c.JSON(fiber.Map{"message": "Token created",
		"data": config.AppConf.Tokens})
}

func TokenEdit(c *fiber.Ctx) error {
	tkn := config.Token{}
	// Parse the body
	if err := c.BodyParser(&tkn); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	// check if user or token already exist
	// fail if either does
	foundToken, err := config.AppConf.FindStructByToken(tkn.Value)
	if err != nil {
		slog.Debug("No matching token found", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Couldn't find a token with that value. Are you sure it exists?"})
	}

	config.AppConf.UpdateStructInToken(*foundToken, tkn)

	// write the new config to file for persisten storage
	if err := config.AppConf.WriteToConf(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error while trying to store the new config. Please try again later"})
	}

	return c.JSON(fiber.Map{"message": "token updated", "data": ""})
}

func TokenDelete(c *fiber.Ctx) error {

	tkn := config.Token{}
	// Parse the body
	if err := c.BodyParser(&tkn); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	// check if user or token already exist
	// fail if either does
	_, err := config.AppConf.FindStructByToken(tkn.Value)
	if err != nil {
		slog.Debug("No matching token found", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Couldn't find a token with that value. Are you sure it exists?"})
	}

	if err = config.AppConf.DeleteStructFromArray(tkn); err != nil {
		slog.Error("Could not delete the token from the config", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error while trying to delete the new config. No token could be deleted"})
	}

	// write the new config to file for persisten storage
	if err := config.AppConf.WriteToConf(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error while trying to delete the new config. Please try again later"})
	}

	return c.JSON(fiber.Map{"message": "token deleted"})
}
