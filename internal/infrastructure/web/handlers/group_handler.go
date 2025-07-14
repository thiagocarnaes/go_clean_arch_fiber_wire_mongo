package handlers

import (
	"github.com/gofiber/fiber/v2"
	"user-management/internal/application/dto"
	"user-management/internal/application/usecases/group"
)

type GroupHandler struct {
	createGroupUseCase         *group.CreateGroupUseCase
	getGroupUseCase            *group.GetGroupUseCase
	updateGroupUseCase         *group.UpdateGroupUseCase
	deleteGroupUseCase         *group.DeleteGroupUseCase
	listGroupsUseCase          *group.ListGroupsUseCase
	addUserToGroupUseCase      *group.AddUserToGroupUseCase
	removeUserFromGroupUseCase *group.RemoveUserFromGroupUseCase
}

func NewGroupHandler(createGroup *group.CreateGroupUseCase, getGroup *group.GetGroupUseCase, updateGroup *group.UpdateGroupUseCase, deleteGroup *group.DeleteGroupUseCase, listGroups *group.ListGroupsUseCase, addUserToGroup *group.AddUserToGroupUseCase, removeUserFromGroup *group.RemoveUserFromGroupUseCase) *GroupHandler {
	return &GroupHandler{
		createGroupUseCase:         createGroup,
		getGroupUseCase:            getGroup,
		updateGroupUseCase:         updateGroup,
		deleteGroupUseCase:         deleteGroup,
		listGroupsUseCase:          listGroups,
		addUserToGroupUseCase:      addUserToGroup,
		removeUserFromGroupUseCase: removeUserFromGroup,
	}
}

func (h *GroupHandler) Create(c *fiber.Ctx) error {
	var groupDTO dto.GroupDTO
	if err := c.BodyParser(&groupDTO); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}
	if err := h.createGroupUseCase.Execute(c.Context(), &groupDTO); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(groupDTO)
}

func (h *GroupHandler) Get(c *fiber.Ctx) error {
	id := c.Params("id")
	groupDTO, err := h.getGroupUseCase.Execute(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Group not found"})
	}
	return c.JSON(groupDTO)
}

func (h *GroupHandler) Update(c *fiber.Ctx) error {
	var groupDTO dto.GroupDTO
	if err := c.BodyParser(&groupDTO); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}
	if err := h.updateGroupUseCase.Execute(c.Context(), &groupDTO); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(groupDTO)
}

func (h *GroupHandler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := h.deleteGroupUseCase.Execute(c.Context(), id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func (h *GroupHandler) List(c *fiber.Ctx) error {
	groups, err := h.listGroupsUseCase.Execute(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(groups)
}

func (h *GroupHandler) AddUser(c *fiber.Ctx) error {
	groupID := c.Params("groupId")
	userID := c.Params("userId")
	if err := h.addUserToGroupUseCase.Execute(c.Context(), groupID, userID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.SendStatus(fiber.StatusOK)
}

func (h *GroupHandler) RemoveUser(c *fiber.Ctx) error {
	groupID := c.Params("groupId")
	userID := c.Params("userId")
	if err := h.removeUserFromGroupUseCase.Execute(c.Context(), groupID, userID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.SendStatus(fiber.StatusOK)
}
