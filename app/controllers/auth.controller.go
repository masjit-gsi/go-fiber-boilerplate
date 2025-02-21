package controllers

import (
	"fmt"
	"time"

	"github.com/fiber-go-template/app/models"
	"github.com/fiber-go-template/app/services"
	"github.com/fiber-go-template/config/utils"
	"github.com/gofiber/fiber/v2"
)

type AuthController struct {
	UserService services.UserService
}

func NewAuthController(service services.UserService) AuthController {
	return AuthController{
		UserService: service,
	}
}

// UserSignIn method to auth user and return access and refresh tokens.
// @Description Auth user and return access and refresh token.
// @Summary auth user and return access and refresh token
// @Tags User
// @Accept json
// @Produce json
// @Param data body models.SignIn true "Data user"
// @Success 200 {string} status "ok"
// @Router /v1/user/login [post]
func (h *AuthController) UserSignIn(c *fiber.Ctx) error {
	signIn := &models.SignIn{}

	// Checking received data from JSON body.
	if err := c.BodyParser(signIn); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	// Get user by email.
	foundedUser, err := h.UserService.GetUserByUsername(signIn.Username)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "user with the given email is not found",
		})
	}

	// Compare given user password with stored in found user.
	compareUserPassword := utils.ComparePasswords(foundedUser.Password, signIn.Password)
	if !compareUserPassword {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "wrong user username address or password",
		})
	}

	// Generate a new pair of access and refresh tokens.
	var credentials []string
	tokens, err := utils.GenerateNewTokens(foundedUser.ID.String(), credentials)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"data": foundedUser,
		"token": fiber.Map{
			"accessToken": tokens.Access,
			"refresh":     tokens.Refresh,
		},
	})
}

// UserSignOut method to de-authorize user and delete refresh token from Redis.
// @Description De-authorize user and delete refresh token from Redis.
// @Summary de-authorize user and delete refresh token from Redis
// @Tags User
// @Accept json
// @Produce json
// @Success 204 {string} status "ok"
// @Security ApiKeyAuth
// @Router /v1/user/logout [post]
func (h *AuthController) UserSignOut(c *fiber.Ctx) error {
	// Get claims from JWT.
	claims, err := utils.ExtractTokenMetadata(c)
	if err != nil {
		// Return status 500 and JWT parse error.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	fmt.Println(claims)

	return c.SendStatus(fiber.StatusNoContent)
}

// RenewTokens method for renew access and refresh tokens.
// @Description Renew access and refresh tokens.
// @Summary renew access and refresh tokens
// @Tags Token
// @Accept json
// @Produce json
// @Param refresh_token body string true "Refresh token"
// @Success 200 {string} status "ok"
// @Security ApiKeyAuth
// @Router /v1/token/renew [post]
func (h *AuthController) RenewTokens(c *fiber.Ctx) error {
	// Get claims from JWT.
	now := time.Now().Unix()
	claims, err := utils.ExtractTokenMetadata(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	// Checking, if now time greather than Access token expiration time.
	expiresAccessToken := claims.Expires
	if now > expiresAccessToken {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"error":   "unauthorized, check expiration time of your token",
		})
	}

	renew := &models.Renew{}

	// Checking received data from JSON body.
	if err := c.BodyParser(renew); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	// Set expiration time from Refresh token of current user.
	expiresRefreshToken, err := utils.ParseRefreshToken(renew.RefreshToken)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	// Checking, if now time greather than Refresh token expiration time.
	if now < expiresRefreshToken {
		userID := claims.UserID.String()

		// Get user by ID.
		foundedUser, err := h.UserService.GetUserByID(userID)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"message": "user with the given ID is not found",
			})
		}

		// Generate JWT Access & Refresh tokens.
		var credentials []string
		tokens, err := utils.GenerateNewTokens(userID, credentials)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"error":   err.Error(),
			})
		}

		return c.JSON(fiber.Map{
			"data": foundedUser,
			"token": fiber.Map{
				"accessToken": tokens.Access,
				"refresh":     tokens.Refresh,
			},
		})
	} else {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "unauthorized, your session was ended earlier",
		})
	}
}
