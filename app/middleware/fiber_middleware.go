package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

// FiberMiddleware provide Fiber's built-in middlewares.
// See: https://docs.gofiber.io/api/middleware
func FiberMiddleware(app *fiber.App) {
	// Add CORS to each route.
	setupCORS(app)
	app.Use(
		// Add simple logger.
		logger.New(),
	)
}
