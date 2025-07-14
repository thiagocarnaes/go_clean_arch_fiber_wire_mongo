package handlers

import (
	"github.com/gofiber/fiber/v2"
	"user-management/internal/application/dto"
	"user-management/internal/application/usecases/user"
)

type UserHandler struct {
	createUserUseCase *user.CreateUserUseCase
	getUserUseCase    *user.GetUserUseCase
	updateUserUseCase *user.UpdateUserUseCase
	deleteUserUseCase *user.DeleteUserUseCase
	listUsersUseCase  *user.ListUsersUseCase
}

func NewUserHandler(createUser *user.CreateUserUseCase, getUser *user.GetUserUseCase, updateUser *user.UpdateUserUseCase, deleteUser *user.DeleteUserUseCase, listUsers *user.ListUsersUseCase) *UserHandler {
	return &UserHandler{
		createUserUseCase: createUser,
		getUserUseCase:    getUser,
		updateUserUseCase: updateUser,
		deleteUserUseCase: deleteUser,
		listUsersUseCase:  listUsers,
	}
}

func (h *UserHandler) Create(c *fiber.Ctx) error {
	var userDTO dto.UserDTO
	if err := c.BodyParser(&userDTO); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}
	if err := h.createUserUseCase.Execute(c.Context(), &userDTO); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(userDTO)
}

func (h *UserHandler) Get(c *fiber.Ctx) error {
	id := c.Params("id")
	userDTO, err := h.getUserUseCase.Execute(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}
	return c.JSON(userDTO)
}

func (h *UserHandler) Update(c *fiber.Ctx) error {
	var userDTO dto.UserDTO
	if err := c.BodyParser(&userDTO); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}
	if err := h.updateUserUseCase.Execute(c.Context(), &userDTO); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(userDTO)
}

func (h *UserHandler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := h.deleteUserUseCase.Execute(c.Context(), id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func (h *UserHandler) List(c *fiber.Ctx) error {
	users, err := h.listUsersUseCase.Execute(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(users)
}
