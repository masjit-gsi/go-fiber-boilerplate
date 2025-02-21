package middleware

import (
	"os"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func setupCORS(app *fiber.App) {
	allowCredentials, _ := strconv.ParseBool(os.Getenv("CORS_ALLOW_CREDENTIALS"))
	maxAge, _ := strconv.Atoi(os.Getenv("CORS_MAX_AGE_SECONDS"))
	app.Use(cors.New(cors.Config{
		AllowCredentials: allowCredentials,
		AllowHeaders:     os.Getenv("CORS_ALLOWED_HEADERS"),
		AllowMethods:     os.Getenv("CORS_ALLOWED_METHODS"),
		AllowOrigins:     os.Getenv("CORS_ALLOWED_ORIGINS"),
		MaxAge:           maxAge,
	}))
}
