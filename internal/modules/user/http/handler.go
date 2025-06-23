package http

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/vynazevedo/go-modular-monolith/internal/modules/user/app"
)

type UserHandlers struct {
	service *app.UserService
}

func NewUserHandlers(service *app.UserService) *UserHandlers {
	return &UserHandlers{service: service}
}

func (h *UserHandlers) CreateUser(c *fiber.Ctx) error {
	var req CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Error: "Invalid request body"})
	}

	cmd := app.CreateUserCommand{
		Email: req.Email,
		Name:  req.Name,
	}

	user, err := h.service.CreateUser(c.Context(), cmd)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Error: err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(UserResponse{
		ID:     user.ID,
		Email:  user.Email,
		Name:   user.Name,
		Status: user.Status,
	})
}

func (h *UserHandlers) GetUser(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Error: "User ID is required"})
	}

	query := app.GetUserQuery{ID: id}
	user, err := h.service.GetUser(c.Context(), query)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(ErrorResponse{Error: "User not found"})
	}

	return c.JSON(UserResponse{
		ID:     user.ID,
		Email:  user.Email,
		Name:   user.Name,
		Status: user.Status,
	})
}

func (h *UserHandlers) UpdateUser(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Error: "User ID is required"})
	}

	var req UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Error: "Invalid request body"})
	}

	cmd := app.UpdateUserCommand{
		ID:   id,
		Name: req.Name,
	}

	user, err := h.service.UpdateUser(c.Context(), cmd)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Error: err.Error()})
	}

	return c.JSON(UserResponse{
		ID:     user.ID,
		Email:  user.Email,
		Name:   user.Name,
		Status: user.Status,
	})
}

func (h *UserHandlers) DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Error: "User ID is required"})
	}

	cmd := app.DeleteUserCommand{ID: id}
	if err := h.service.DeleteUser(c.Context(), cmd); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Error: err.Error()})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (h *UserHandlers) ActivateUser(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Error: "User ID is required"})
	}

	cmd := app.ActivateUserCommand{ID: id}
	user, err := h.service.ActivateUser(c.Context(), cmd)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Error: err.Error()})
	}

	return c.JSON(UserResponse{
		ID:     user.ID,
		Email:  user.Email,
		Name:   user.Name,
		Status: user.Status,
	})
}

func (h *UserHandlers) DeactivateUser(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Error: "User ID is required"})
	}

	cmd := app.DeactivateUserCommand{ID: id}
	user, err := h.service.DeactivateUser(c.Context(), cmd)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Error: err.Error()})
	}

	return c.JSON(UserResponse{
		ID:     user.ID,
		Email:  user.Email,
		Name:   user.Name,
		Status: user.Status,
	})
}

func (h *UserHandlers) ListUsers(c *fiber.Ctx) error {
	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(c.Query("limit", "10"))
	if err != nil || limit < 1 || limit > 100 {
		limit = 10
	}

	query := app.ListUsersQuery{
		Page:  page,
		Limit: limit,
	}

	users, err := h.service.ListUsers(c.Context(), query)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Error: err.Error()})
	}

	userResponses := make([]UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = UserResponse{
			ID:     user.ID,
			Email:  user.Email,
			Name:   user.Name,
			Status: user.Status,
		}
	}

	return c.JSON(UsersResponse{
		Users: userResponses,
		Page:  page,
		Limit: limit,
		Total: len(users),
	})
}
