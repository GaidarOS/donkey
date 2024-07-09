package middleware

import (
	"log/slog"
	"receipt_store/config"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// TokenMiddleware
// Checks if the requester's token exists and if he has access to
func TokenMiddleware(c *fiber.Ctx) error {

	// The request should contain a token if the token is invalid throw 403

	var token string = ""
	if c.Queries()["token"] != "" {
		token = c.Queries()["token"]
	} else {
		token = c.Get("Token")
	}

	if token == "" {
		slog.Error("No token \"Token\" found in the request!")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":     "error",
			"message":    "Malformed: Missing proper headers",
			"suggestion": "Make sure you include the required header(\"Token\") and provide valid values",
		})
	}

	result, err := config.AppConf.FindStructByUser(token)
	if err != nil {
		slog.Error("No token found", err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "Unauthorized: No user found with that token",
		})
	}

	// Check if the user has access to this folder
	param := c.Params("*")
	if len(strings.Split(param, "/")) > 0 {
		param = strings.Split(param, "/")[0]
	}
	base_path := c.Path()
	if len(strings.Split(c.Path(), "/")) > 0 {
		base_path = strings.Split(c.Path(), "/")[3]
	}
	if param != "thumbnails" && base_path != "user" {
		if !(result.AccessPaths[param] || result.Admin) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  "error",
				"message": "Unauthorized: No permissions to access this folder!",
			})
		}
	}

	slog.Debug("Found this user from token", slog.Any("config", result))

	slog.Debug("Request body", slog.String("body", string(c.Body())))

	return c.Next()
}

func AdminMiddleware(c *fiber.Ctx) error {

	slog.Info("Verifying credentials!")
	slog.Debug("Request headers", slog.Any("all-headers", c.GetReqHeaders()))

	result, err := config.AppConf.FindStructByUser(c.Get("Token"))
	if err != nil {
		slog.Error("No token found", err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "Unauthorized: No user found with that token",
		})
	}

	slog.Debug("Found this user from token", slog.Any("config", result))
	if result == nil || !result.Admin {
		slog.Error("A non-admin user requested access")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "Unauthorized: Missing permissions",
		})
	}
	return c.Next()
}
