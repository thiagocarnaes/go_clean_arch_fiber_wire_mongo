package controllers

import (
	"user-management/internal/application/dto"
	"user-management/internal/application/usecases/group"
	"user-management/internal/infrastructure/web/validators"

	"github.com/gofiber/fiber/v2"
)

type GroupController struct {
	validator                  *validators.InputValidator
	createGroupUseCase         *group.CreateGroupUseCase
	getGroupUseCase            *group.GetGroupUseCase
	updateGroupUseCase         *group.UpdateGroupUseCase
	deleteGroupUseCase         *group.DeleteGroupUseCase
	listGroupsUseCase          *group.ListGroupsUseCase
	addUserToGroupUseCase      *group.AddUserToGroupUseCase
	removeUserFromGroupUseCase *group.RemoveUserFromGroupUseCase
}

func NewGroupController(createGroup *group.CreateGroupUseCase, getGroup *group.GetGroupUseCase, updateGroup *group.UpdateGroupUseCase, deleteGroup *group.DeleteGroupUseCase, listGroups *group.ListGroupsUseCase, addUserToGroup *group.AddUserToGroupUseCase, removeUserFromGroup *group.RemoveUserFromGroupUseCase) *GroupController {
	return &GroupController{
		validator:                  validators.NewInputValidator(),
		createGroupUseCase:         createGroup,
		getGroupUseCase:            getGroup,
		updateGroupUseCase:         updateGroup,
		deleteGroupUseCase:         deleteGroup,
		listGroupsUseCase:          listGroups,
		addUserToGroupUseCase:      addUserToGroup,
		removeUserFromGroupUseCase: removeUserFromGroup,
	}
}

func (h *GroupController) Create(c *fiber.Ctx) error {
	var createGroupDTO dto.CreateGroupRequestDTO
	
	if err := h.validator.ParseAndValidate(c, &createGroupDTO); err != nil {
		if validationErr, ok := err.(*validators.ValidationError); ok {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": validationErr.Message})
		}
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	groupDTO, err := h.createGroupUseCase.Execute(c.Context(), &createGroupDTO)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(groupDTO)
}

func (h *GroupController) Get(c *fiber.Ctx) error {
	id := c.Params("id")
	groupDTO, err := h.getGroupUseCase.Execute(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Group not found"})
	}
	return c.JSON(groupDTO)
}

func (h *GroupController) Update(c *fiber.Ctx) error {
	var updateGroupDTO dto.CreateGroupRequestDTO
	
	if err := h.validator.ParseAndValidate(c, &updateGroupDTO); err != nil {
		if validationErr, ok := err.(*validators.ValidationError); ok {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": validationErr.Message})
		}
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	groupID := c.Params("id")
	responseDTO, err := h.updateGroupUseCase.Execute(c.Context(), groupID, &updateGroupDTO)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(responseDTO)
}

func (h *GroupController) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := h.deleteGroupUseCase.Execute(c.Context(), id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func (h *GroupController) List(c *fiber.Ctx) error {
	groups, err := h.listGroupsUseCase.Execute(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(groups)
}

func (h *GroupController) AddUser(c *fiber.Ctx) error {
	groupID := c.Params("groupId")
	userID := c.Params("userId")
	if err := h.addUserToGroupUseCase.Execute(c.Context(), groupID, userID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.SendStatus(fiber.StatusOK)
}

func (h *GroupController) RemoveUser(c *fiber.Ctx) error {
	groupID := c.Params("groupId")
	userID := c.Params("userId")
	if err := h.removeUserFromGroupUseCase.Execute(c.Context(), groupID, userID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.SendStatus(fiber.StatusOK)
}
