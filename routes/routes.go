package routes

import (
	"github.com/fiber-go-template/app/middleware"
	"github.com/gofiber/fiber/v2"
	swagger "github.com/gofiber/swagger"
)

func SetupRoutes(a *fiber.App, c Injection) {
	SwaggerRoute(a)
	// Create routes group.
	route := a.Group("/api/v1")

	// AUTH
	userController := c.AuthController
	route.Post("/user/login", userController.UserSignIn)
	route.Post("/user/logout", middleware.JWTProtected(), userController.UserSignOut)
	route.Post("/token/renew", middleware.JWTProtected(), userController.RenewTokens)

	// BOOK
	authorController := c.AuthorController
	route.Get("/authors", middleware.JWTProtected(), authorController.ResolveAll)
	route.Get("/authors/all", middleware.JWTProtected(), authorController.GetAll)
	route.Get("/author/:id", authorController.FindByID)
	route.Post("/author", middleware.JWTProtected(), authorController.Create)
	route.Put("/author/:id", middleware.JWTProtected(), authorController.Update)
	route.Delete("/author/:id", middleware.JWTProtected(), authorController.Delete)
}

func SwaggerRoute(a *fiber.App) {
	route := a.Group("/swagger")
	route.Get("*", swagger.HandlerDefault)
}

func NotFoundRoute(a *fiber.App) {
	a.Use(
		func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": true,
				"msg":   "sorry, endpoint is not found",
			})
		},
	)
}
