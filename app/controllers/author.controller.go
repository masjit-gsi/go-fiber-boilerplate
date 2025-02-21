package controllers

import (
	"strconv"
	"time"

	"github.com/fiber-go-template/app/models"
	"github.com/fiber-go-template/app/services"
	"github.com/fiber-go-template/config/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofrs/uuid"
)

type AuthorController struct {
	AuthorService services.AuthorService
}

func NewAuthorController(service services.AuthorService) AuthorController {
	return AuthorController{
		AuthorService: service,
	}
}

// ResolveAll list all Author.
// @Summary Get list all Author.
// @Description endpoint get all data with pagination.
// @Tags Author
// @Produce json
// @Param q query string false "Keyword search"
// @Param pageSize query int true "Set pageSize data"
// @Param pageNumber query int true "Set page number"
// @Param sortBy query string false "Set sortBy parameter is one of [ nama ]"
// @Param sortType query string false "Set sortType with asc or desc"
// @Success 200 {object} response.Base
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Security ApiKeyAuth
// @Router /v1/authors [get]
func (h *AuthorController) ResolveAll(c *fiber.Ctx) error {
	keyword := c.Query("keyword")
	pageSizeStr := c.Query("pageSize")
	pageNumberStr := c.Query("pageNumber")
	sortBy := c.Query("sortBy")
	if sortBy == "" {
		sortBy = "createdAt"
	}

	sortType := c.Query("sortType")
	if sortType == "" {
		sortType = "DESC"
	}
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil {
		return err
	}

	pageNumber, err := strconv.Atoi(pageNumberStr)
	if err != nil {
		return err
	}

	req := models.StandardRequest{
		Keyword:    keyword,
		PageSize:   pageSize,
		PageNumber: pageNumber,
		SortBy:     sortBy,
		SortType:   sortType,
	}

	data, err := h.AuthorService.ResolveAll(req)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"data": data,
	})
}

// GetAll func gets all exists authors.
// @Description Get all exists authors.
// @Summary get all exists authors
// @Tags Author
// @Accept json
// @Produce json
// @Success 200 {object} response.Base{models.Author}
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Security ApiKeyAuth
// @Router /v1/authors/all [get]
func (h *AuthorController) GetAll(c *fiber.Ctx) error {
	// Get all data
	data, err := h.AuthorService.GetAll()
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Data were not found",
			"data":    nil,
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": nil,
		"data":    data,
	})
}

// FindByID func gets author by given ID or 404 error.
// @Description Get author by given ID.
// @Summary get author by given ID
// @Tags Author
// @Accept json
// @Produce json
// @Param id path string true "Author ID"
// @Success 200 {object} response.Base{models.Author}
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Security ApiKeyAuth
// @Router /v1/author/{id} [get]
func (h *AuthorController) FindByID(c *fiber.Ctx) error {
	// Catch data ID from URL.
	id, err := uuid.FromString(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	// Get data by ID.
	data, err := h.AuthorService.FindByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Data with the given ID is not found",
			"data":    nil,
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": nil,
		"data":    data,
	})
}

// Create func for creates a new author.
// @Description Create a new author.
// @Summary create a new author
// @Tags Author
// @Accept json
// @Produce json
// @Param data body models.AuthorRequest true "Author"
// @Success 200 {object} response.Base{data=models.Author}
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Security ApiKeyAuth
// @Router /v1/author [post]
func (h *AuthorController) Create(c *fiber.Ctx) error {
	now := time.Now().Unix()
	// Get claims from JWT.
	claims, err := utils.ExtractTokenMetadata(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	// Checking, if now time greather than expiration from JWT.
	expires := claims.Expires
	if now > expires {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"error":   "unauthorized, check expiration time of your token",
		})
	}

	var request models.AuthorRequest

	// Check, if received JSON data is valid.
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	// Validate author fields.
	validate := utils.NewValidator()
	if err := validate.Struct(request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   utils.ValidatorErrors(err),
		})
	}

	// Create author by given model.
	request.UserID = claims.UserID
	data, err := h.AuthorService.Create(request)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Create data successfully",
		"data":    data,
	})
}

// Update func for updates author by given ID.
// @Description Update author.
// @Summary update author
// @Tags Author
// @Accept json
// @Produce json
// @Param id path string true "Author ID"
// @Param data body models.AuthorRequest true "Author"
// @Success 200 {object} response.Base{models.Author}
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Security ApiKeyAuth
// @Router /v1/author/{id} [put]
func (h *AuthorController) Update(c *fiber.Ctx) error {
	id, err := uuid.FromString(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	// Get claims from JWT.
	now := time.Now().Unix()
	claims, err := utils.ExtractTokenMetadata(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	// Checking, if now time greather than expiration from JWT.
	expires := claims.Expires
	if now > expires {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"error":   "unauthorized, check expiration time of your token",
		})
	}

	var request models.AuthorRequest

	// Check, if received JSON data is valid.
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	// Checking, if author with given ID is exists.
	foundedAuthor, err := h.AuthorService.FindByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Data with this ID not found",
		})
	}

	// Validate author fields.
	validate := utils.NewValidator()
	if err := validate.Struct(request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   utils.ValidatorErrors(err),
		})
	}

	// Update author by given ID.
	request.ID = id
	request.UserID = claims.UserID
	data, err := h.AuthorService.Update(foundedAuthor.ID, request)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "Update data successfully",
		"data":    data,
	})
}

// Delete func for delete author by given ID.
// @Description Delete author by given ID.
// @Summary delete author by given ID
// @Tags Author
// @Accept json
// @Produce json
// @Param id path string true "Author ID"
// @Success 200 {object} response.Base
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Security ApiKeyAuth
// @Router /v1/author/{id} [delete]
func (h *AuthorController) Delete(c *fiber.Ctx) error {
	id, err := uuid.FromString(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	// Get claims from JWT.
	now := time.Now().Unix()
	claims, err := utils.ExtractTokenMetadata(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	// Checking, if now time greather than expiration from JWT.
	expires := claims.Expires
	if now > expires {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"error":   "unauthorized, check expiration time of your token",
		})
	}

	// Checking, if author with given ID is exists.
	foundedAuthor, err := h.AuthorService.FindByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Data with this ID not found",
		})
	}

	// Set user ID from JWT data of current user.
	err = h.AuthorService.Delete(foundedAuthor.ID, claims.UserID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "Delete data successfully",
	})
}
