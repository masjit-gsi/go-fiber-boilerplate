package routes

import (
	"github.com/fiber-go-template/app/controllers"
	"github.com/fiber-go-template/app/repository"
	"github.com/fiber-go-template/app/services"
	"github.com/fiber-go-template/database"
)

type Injection struct {
	AuthController   controllers.AuthController
	AuthorController controllers.AuthorController
}

// Define Dependency Injection
func CallDependenciesInjection() Injection {
	DbConnect, _ := database.NewDBConnection()
	// Auth
	userRepository := repository.NewUserRepository(DbConnect)
	userService := services.NewUserService(userRepository)
	authController := controllers.NewAuthController(userService)
	// Author
	authorRepository := repository.NewAuthorRepository(DbConnect)
	authorService := services.NewAuthorService(DbConnect, authorRepository)
	authorController := controllers.NewAuthorController(authorService)

	return Injection{
		AuthController:   authController,
		AuthorController: authorController,
	}
}
