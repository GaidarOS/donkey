package middleware

import (
	"log/slog"
	"receipt_store/config"

	"github.com/gofiber/fiber/v2"
)

// TokenMiddleware
// Checks if the requester's token exists and if he has access to
func TokenMiddleware(c *fiber.Ctx) error {

	// The request should contain a token if the token is invalid throw 403

	if c.Get("Token") == "" {
		slog.Error("No token \"Token\" found in the request!")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":     "error",
			"message":    "Malformed: Missing proper headers",
			"suggestion": "Make sure you include the required header(\"Token\") and provide valid values",
		})
	}

	result, err := config.AppConf.FindStructByToken(c.Get("Token"))
	if err != nil {
		slog.Error("No token found", err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "Unauthorized: No user found with that token",
		}) 
	}

	slog.Debug("Found this user from token", slog.Any("config", result))

	slog.Info("Request body", slog.String("body", string(c.Body())))

	slog.Info("Middleware executed before route handler")
	return c.Next()
}

func AdminMiddleware(c *fiber.Ctx) error {

	slog.Info("Verifying credentials!")
	slog.Debug("Request headers", slog.Any("all-headers", c.GetReqHeaders()))

	result, err := config.AppConf.FindStructByToken(c.Get("Token"))
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
