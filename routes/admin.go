package routes

import (
	"encoding/json"
	"log/slog"
	"receipt_store/config"
	"receipt_store/helper"

	"github.com/gofiber/fiber/v2"
)

func UsersList(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "List of user tokens",
		"data": config.AppConf.Users})
}

func UserCreate(c *fiber.Ctx) error {

	tkn := config.User{}
	// Parse the body
	if err := c.BodyParser(&tkn); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	// check if user or token already exist
	// fail if either does
	foundUser, err := config.AppConf.FindStructByToken(tkn.Password)
	if err != nil {
		slog.Debug("No matching token found", err)
	}
	foundName, err := config.AppConf.FindStructByName(tkn.UserName)
	if err != nil {
		slog.Debug("No matching username found", err)
	}
	if foundName != nil || foundUser != nil {
		slog.Error("Found an existing username or token. Try updating instead.")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Found an existing username or token. Try updating instead"})
	}

	// append the new token to the list
	config.AppConf.Users = append(config.AppConf.Users, tkn)

	// write the new config to file for persisten storage
	if err := config.AppConf.WriteToConf(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error while trying to store the new config. Please try again later"})
	}

	return c.JSON(fiber.Map{"message": "User created",
		"data": config.AppConf.Users})
}

func UserEdit(c *fiber.Ctx) error {
	tkn := config.User{}
	// Parse the body
	if err := c.BodyParser(&tkn); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	// check if user or token already exist
	// fail if either does

	foundUser, err := config.AppConf.FindStructByToken(tkn.UserName)
	if err != nil {
		foundUser, err = config.AppConf.FindStructByName(tkn.UserName)
	}
	if err != nil {
		slog.Debug("No matching token found", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Couldn't find a token with that value. Are you sure it exists?"})
	}

	config.AppConf.UpdateStructInUser(*foundUser, tkn)

	// write the new config to file for persisten storage
	if err := config.AppConf.WriteToConf(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error while trying to store the new config. Please try again later"})
	}

	return c.JSON(fiber.Map{"message": "token updated", "data": ""})
}

func UserDelete(c *fiber.Ctx) error {

	tkn := config.User{}
	// Parse the body
	if err := c.BodyParser(&tkn); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	// check if user or token already exist
	// fail if either does
	_, err := config.AppConf.FindStructByToken(tkn.Password)
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

func UpdateConfig(c *fiber.Ctx) error {

	var conf config.Config

	// Parse the body
	if err := c.BodyParser(&conf); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid request"})
	}
	rq, _ := json.Marshal(conf)

	if err := helper.WriteToFile(config.AppConf.ConfFile, rq); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to update config"})
	}
	return c.JSON(fiber.Map{"message": "Configuration updated successfully"})
}

func GetConfig(c *fiber.Ctx) error {

	return c.JSON(config.AppConf)
}
