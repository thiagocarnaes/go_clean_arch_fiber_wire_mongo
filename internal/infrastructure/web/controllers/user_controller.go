package controllers

import (
	"user-management/internal/application/dto"
	"user-management/internal/application/usecases/user"
	"user-management/internal/infrastructure/web/validators"

	"github.com/gofiber/fiber/v2"
)

type UserController struct {
	validator         *validators.InputValidator
	createUserUseCase *user.CreateUserUseCase
	getUserUseCase    *user.GetUserUseCase
	updateUserUseCase *user.UpdateUserUseCase
	deleteUserUseCase *user.DeleteUserUseCase
	listUsersUseCase  *user.ListUsersUseCase
}

func NewUserController(createUser *user.CreateUserUseCase, getUser *user.GetUserUseCase, updateUser *user.UpdateUserUseCase, deleteUser *user.DeleteUserUseCase, listUsers *user.ListUsersUseCase) *UserController {
	return &UserController{
		validator:         validators.NewInputValidator(),
		createUserUseCase: createUser,
		getUserUseCase:    getUser,
		updateUserUseCase: updateUser,
		deleteUserUseCase: deleteUser,
		listUsersUseCase:  listUsers,
	}
}

func (h *UserController) Create(c *fiber.Ctx) error {
	var createUserDTO dto.CreateUserRequestDTO
	
	// Parse e valida em uma operação
	if err := h.validator.ParseAndValidate(c, &createUserDTO); err != nil {
		if validationErr, ok := err.(*validators.ValidationError); ok {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": validationErr.Message})
		}
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	responseDTO, err := h.createUserUseCase.Execute(c.Context(), &createUserDTO)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(responseDTO)
}

func (h *UserController) Get(c *fiber.Ctx) error {
	id := c.Params("id")
	userDTO, err := h.getUserUseCase.Execute(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}
	return c.JSON(userDTO)
}

func (h *UserController) Update(c *fiber.Ctx) error {
	var updateUserDTO dto.CreateUserRequestDTO
	
	if err := h.validator.ParseAndValidate(c, &updateUserDTO); err != nil {
		if validationErr, ok := err.(*validators.ValidationError); ok {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": validationErr.Message})
		}
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	userID := c.Params("id")
	responseDTO, err := h.updateUserUseCase.Execute(c.Context(), userID, &updateUserDTO)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(responseDTO)
}

func (h *UserController) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := h.deleteUserUseCase.Execute(c.Context(), id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func (h *UserController) List(c *fiber.Ctx) error {
	users, err := h.listUsersUseCase.Execute(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(users)
}
