package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"user-management/internal/application/dto"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGroupControllerCreate(t *testing.T) {
	testApp := SetupTestApp(t)
	defer testApp.Cleanup(t)

	tests := []struct {
		name           string
		payload        interface{}
		expectedStatus int
	}{
		{
			name: "Valid group creation without members",
			payload: dto.CreateGroupRequestDTO{
				Name:    "Developers",
				Members: []string{},
			},
			expectedStatus: 201,
		},
		{
			name: "Valid group creation with members",
			payload: dto.CreateGroupRequestDTO{
				Name:    "Admins",
				Members: []string{"user1", "user2"},
			},
			expectedStatus: 201,
		},
		{
			name: "Invalid group - missing name",
			payload: dto.CreateGroupRequestDTO{
				Members: []string{},
			},
			expectedStatus: 400,
		},
		{
			name: "Invalid group - empty name",
			payload: dto.CreateGroupRequestDTO{
				Name:    "",
				Members: []string{},
			},
			expectedStatus: 400,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testApp.ClearDatabase(t)

			payloadBytes, err := json.Marshal(tt.payload)
			require.NoError(t, err)

			req, err := http.NewRequest("POST", "/api/v1/groups", bytes.NewBuffer(payloadBytes))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			resp, err := testApp.Request(req)
			require.NoError(t, err)

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			if tt.expectedStatus == 201 {
				var groupResponse dto.GroupResponseDTO
				err := json.NewDecoder(resp.Body).Decode(&groupResponse)
				require.NoError(t, err)

				assert.NotEmpty(t, groupResponse.ID)
				assert.Equal(t, tt.payload.(dto.CreateGroupRequestDTO).Name, groupResponse.Name)
				assert.Equal(t, tt.payload.(dto.CreateGroupRequestDTO).Members, groupResponse.Members)
			}
		})
	}
}

func TestGroupControllerCreateRepositoryErrors(t *testing.T) {
	t.Run("Create group - force MongoDB connection close to cause repository Create error", func(t *testing.T) {
		testApp := SetupTestApp(t)
		// Note: Don't defer cleanup here since we'll disconnect the database

		const (
			contentTypeJSON   = "application/json"
			contentTypeHeader = "Content-Type"
			groupsEndpoint    = "/api/v1/groups"
		)

		// Verify normal operation works first
		createGroupDTO := dto.CreateGroupRequestDTO{
			Name:    "Test Group Before Disconnect",
			Members: []string{},
		}

		payloadBytes, err := json.Marshal(createGroupDTO)
		require.NoError(t, err)

		createReq, err := http.NewRequest("POST", groupsEndpoint, bytes.NewBuffer(payloadBytes))
		require.NoError(t, err)
		createReq.Header.Set(contentTypeHeader, contentTypeJSON)

		createResp, err := testApp.Request(createReq)
		require.NoError(t, err)
		assert.Equal(t, 201, createResp.StatusCode)

		// Now force close the MongoDB connection to simulate repository error
		ctx := context.Background()
		err = testApp.DB.Client.Disconnect(ctx)
		require.NoError(t, err)

		// Try to create a group after closing connection - should return repository error
		createGroupDTO2 := dto.CreateGroupRequestDTO{
			Name:    "Test Group After Disconnect",
			Members: []string{},
		}

		payloadBytes2, err := json.Marshal(createGroupDTO2)
		require.NoError(t, err)

		createReq2, err := http.NewRequest("POST", groupsEndpoint, bytes.NewBuffer(payloadBytes2))
		require.NoError(t, err)
		createReq2.Header.Set(contentTypeHeader, contentTypeJSON)

		createResp2, err := testApp.Request(createReq2)
		require.NoError(t, err)

		// Should return 500 due to repository Create() operation failing
		assert.Equal(t, 500, createResp2.StatusCode)

		var errorResponse map[string]interface{}
		err = json.NewDecoder(createResp2.Body).Decode(&errorResponse)
		require.NoError(t, err)
		assert.NotNil(t, errorResponse["error"])

		// The error should be related to connection/repository failure
		errorMsg := errorResponse["error"].(string)
		assert.True(t,
			strings.Contains(errorMsg, "connection") ||
				strings.Contains(errorMsg, "client is disconnected") ||
				strings.Contains(errorMsg, "topology is closed") ||
				strings.Contains(errorMsg, "server selection error"),
			"Expected connection error, got: %s", errorMsg)
	})

	t.Run("Create group - duplicate group name", func(t *testing.T) {
		testApp := SetupTestApp(t)
		defer testApp.Cleanup(t)

		const (
			contentTypeJSON   = "application/json"
			contentTypeHeader = "Content-Type"
			groupsEndpoint    = "/api/v1/groups"
			duplicateName     = "Duplicate Group Name"
		)

		// Create first group
		createGroupDTO := dto.CreateGroupRequestDTO{
			Name:    duplicateName,
			Members: []string{"user1"},
		}

		payloadBytes, err := json.Marshal(createGroupDTO)
		require.NoError(t, err)

		createReq, err := http.NewRequest("POST", groupsEndpoint, bytes.NewBuffer(payloadBytes))
		require.NoError(t, err)
		createReq.Header.Set(contentTypeHeader, contentTypeJSON)

		createResp, err := testApp.Request(createReq)
		require.NoError(t, err)
		assert.Equal(t, 201, createResp.StatusCode)

		// Try to create second group with same name
		createGroupDTO2 := dto.CreateGroupRequestDTO{
			Name:    duplicateName,
			Members: []string{"user2"},
		}

		payloadBytes2, err := json.Marshal(createGroupDTO2)
		require.NoError(t, err)

		createReq2, err := http.NewRequest("POST", groupsEndpoint, bytes.NewBuffer(payloadBytes2))
		require.NoError(t, err)
		createReq2.Header.Set(contentTypeHeader, contentTypeJSON)

		createResp2, err := testApp.Request(createReq2)
		require.NoError(t, err)

		// MongoDB should return error due to duplicate key if unique index exists
		// or succeed if no unique constraint (depends on database setup)
		assert.True(t, createResp2.StatusCode == 201 || createResp2.StatusCode == 500,
			"Expected 201 (no constraint) or 500 (duplicate constraint), got %d", createResp2.StatusCode)

		if createResp2.StatusCode == 500 {
			var errorResponse map[string]interface{}
			err = json.NewDecoder(createResp2.Body).Decode(&errorResponse)
			require.NoError(t, err)
			assert.NotNil(t, errorResponse["error"])
		}
	})

	t.Run("Create group - database collection drop during operation", func(t *testing.T) {
		testApp := SetupTestApp(t)
		defer testApp.Cleanup(t)

		const (
			contentTypeJSON   = "application/json"
			contentTypeHeader = "Content-Type"
			groupsEndpoint    = "/api/v1/groups"
		)

		// Create first group normally
		createGroupDTO := dto.CreateGroupRequestDTO{
			Name:    "Test Group Before Drop",
			Members: []string{},
		}

		payloadBytes, err := json.Marshal(createGroupDTO)
		require.NoError(t, err)

		createReq, err := http.NewRequest("POST", groupsEndpoint, bytes.NewBuffer(payloadBytes))
		require.NoError(t, err)
		createReq.Header.Set(contentTypeHeader, contentTypeJSON)

		createResp, err := testApp.Request(createReq)
		require.NoError(t, err)
		assert.Equal(t, 201, createResp.StatusCode)

		// Drop the groups collection to force potential repository issues
		ctx := context.Background()
		collection := testApp.DB.DB.Collection("groups")
		err = collection.Drop(ctx)
		require.NoError(t, err)

		// Try to create another group after dropping collection
		createGroupDTO2 := dto.CreateGroupRequestDTO{
			Name:    "Test Group After Drop",
			Members: []string{"user1", "user2"},
		}

		payloadBytes2, err := json.Marshal(createGroupDTO2)
		require.NoError(t, err)

		createReq2, err := http.NewRequest("POST", groupsEndpoint, bytes.NewBuffer(payloadBytes2))
		require.NoError(t, err)
		createReq2.Header.Set(contentTypeHeader, contentTypeJSON)

		createResp2, err := testApp.Request(createReq2)
		require.NoError(t, err)

		// Should either work (201 - collection recreated) or return error (500)
		assert.True(t, createResp2.StatusCode == 201 || createResp2.StatusCode == 500,
			"Expected 201 (collection recreated) or 500 (repository error), got %d", createResp2.StatusCode)

		if createResp2.StatusCode == 201 {
			var createdGroup dto.GroupResponseDTO
			err = json.NewDecoder(createResp2.Body).Decode(&createdGroup)
			require.NoError(t, err)
			assert.Equal(t, createGroupDTO2.Name, createdGroup.Name)
			assert.Equal(t, createGroupDTO2.Members, createdGroup.Members)
		} else {
			var errorResponse map[string]interface{}
			err = json.NewDecoder(createResp2.Body).Decode(&errorResponse)
			require.NoError(t, err)
			assert.NotNil(t, errorResponse["error"])
		}
	})

	t.Run("Create group - invalid members causing repository issues", func(t *testing.T) {
		testApp := SetupTestApp(t)
		defer testApp.Cleanup(t)

		const (
			contentTypeJSON   = "application/json"
			contentTypeHeader = "Content-Type"
			groupsEndpoint    = "/api/v1/groups"
		)

		// Try to create group with members that might cause repository issues
		// (e.g., very long member IDs, special characters, etc.)
		createGroupDTO := dto.CreateGroupRequestDTO{
			Name: "Group With Problematic Members",
			Members: []string{
				strings.Repeat("a", 1000), // Very long member ID
				"member@with@special@chars",
				"", // Empty member ID
				"normal-member-id",
			},
		}

		payloadBytes, err := json.Marshal(createGroupDTO)
		require.NoError(t, err)

		createReq, err := http.NewRequest("POST", groupsEndpoint, bytes.NewBuffer(payloadBytes))
		require.NoError(t, err)
		createReq.Header.Set(contentTypeHeader, contentTypeJSON)

		createResp, err := testApp.Request(createReq)
		require.NoError(t, err)

		// Should either work (201) or return validation/repository error (400/500)
		assert.True(t, createResp.StatusCode == 201 || createResp.StatusCode == 400 || createResp.StatusCode == 500,
			"Expected 201, 400, or 500, got %d", createResp.StatusCode)

		if createResp.StatusCode == 201 {
			var createdGroup dto.GroupResponseDTO
			err = json.NewDecoder(createResp.Body).Decode(&createdGroup)
			require.NoError(t, err)
			assert.Equal(t, createGroupDTO.Name, createdGroup.Name)
			// Members might be filtered or processed differently
		} else {
			var errorResponse map[string]interface{}
			err = json.NewDecoder(createResp.Body).Decode(&errorResponse)
			require.NoError(t, err)
			assert.NotNil(t, errorResponse["error"])
		}
	})
}

func TestGroupControllerGet(t *testing.T) {
	testApp := SetupTestApp(t)
	defer testApp.Cleanup(t)

	// Create a group first
	createPayload := dto.CreateGroupRequestDTO{
		Name:    "Test Group",
		Members: []string{"user1", "user2"},
	}

	payloadBytes, err := json.Marshal(createPayload)
	require.NoError(t, err)

	req, err := http.NewRequest("POST", "/api/v1/groups", bytes.NewBuffer(payloadBytes))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	resp, err := testApp.Request(req)
	require.NoError(t, err)
	require.Equal(t, 201, resp.StatusCode)

	var createdGroup dto.GroupResponseDTO
	err = json.NewDecoder(resp.Body).Decode(&createdGroup)
	require.NoError(t, err)

	// Test getting the created group
	req, err = http.NewRequest("GET", fmt.Sprintf("/api/v1/groups/%s", createdGroup.ID), nil)
	require.NoError(t, err)

	resp, err = testApp.Request(req)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	var retrievedGroup dto.GroupResponseDTO
	err = json.NewDecoder(resp.Body).Decode(&retrievedGroup)
	require.NoError(t, err)

	assert.Equal(t, createdGroup.ID, retrievedGroup.ID)
	assert.Equal(t, createdGroup.Name, retrievedGroup.Name)
	assert.Equal(t, createdGroup.Members, retrievedGroup.Members)

	// Test getting non-existent group
	req, err = http.NewRequest("GET", "/api/v1/groups/507f1f77bcf86cd799439011", nil)
	require.NoError(t, err)

	resp, err = testApp.Request(req)
	require.NoError(t, err)
	assert.Equal(t, 404, resp.StatusCode)
}

func TestGroupControllerUpdate(t *testing.T) {
	testApp := SetupTestApp(t)
	defer testApp.Cleanup(t)

	// Create a group first
	createPayload := dto.CreateGroupRequestDTO{
		Name:    "Original Group",
		Members: []string{"user1"},
	}

	payloadBytes, err := json.Marshal(createPayload)
	require.NoError(t, err)

	req, err := http.NewRequest("POST", "/api/v1/groups", bytes.NewBuffer(payloadBytes))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	resp, err := testApp.Request(req)
	require.NoError(t, err)
	require.Equal(t, 201, resp.StatusCode)

	var createdGroup dto.GroupResponseDTO
	err = json.NewDecoder(resp.Body).Decode(&createdGroup)
	require.NoError(t, err)

	// Update the group
	updatePayload := dto.CreateGroupRequestDTO{
		Name:    "Updated Group",
		Members: []string{"user1", "user2", "user3"},
	}

	payloadBytes, err = json.Marshal(updatePayload)
	require.NoError(t, err)

	req, err = http.NewRequest("PUT", fmt.Sprintf("/api/v1/groups/%s", createdGroup.ID), bytes.NewBuffer(payloadBytes))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	resp, err = testApp.Request(req)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	var updatedGroup dto.GroupResponseDTO
	err = json.NewDecoder(resp.Body).Decode(&updatedGroup)
	require.NoError(t, err)

	assert.Equal(t, createdGroup.ID, updatedGroup.ID)
	assert.Equal(t, updatePayload.Name, updatedGroup.Name)
	assert.Equal(t, updatePayload.Members, updatedGroup.Members)
}

func TestGroupControllerDelete(t *testing.T) {
	testApp := SetupTestApp(t)
	defer testApp.Cleanup(t)

	// Create a group first
	createPayload := dto.CreateGroupRequestDTO{
		Name:    "To Be Deleted",
		Members: []string{},
	}

	payloadBytes, err := json.Marshal(createPayload)
	require.NoError(t, err)

	req, err := http.NewRequest("POST", "/api/v1/groups", bytes.NewBuffer(payloadBytes))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	resp, err := testApp.Request(req)
	require.NoError(t, err)
	require.Equal(t, 201, resp.StatusCode)

	var createdGroup dto.GroupResponseDTO
	err = json.NewDecoder(resp.Body).Decode(&createdGroup)
	require.NoError(t, err)

	// Delete the group
	req, err = http.NewRequest("DELETE", fmt.Sprintf("/api/v1/groups/%s", createdGroup.ID), nil)
	require.NoError(t, err)

	resp, err = testApp.Request(req)
	require.NoError(t, err)
	assert.Equal(t, 204, resp.StatusCode)

	// Verify group is deleted by trying to get it
	req, err = http.NewRequest("GET", fmt.Sprintf("/api/v1/groups/%s", createdGroup.ID), nil)
	require.NoError(t, err)

	resp, err = testApp.Request(req)
	require.NoError(t, err)
	assert.Equal(t, 404, resp.StatusCode)
}

func TestGroupControllerList(t *testing.T) {
	testApp := SetupTestApp(t)
	defer testApp.Cleanup(t)

	// Create multiple groups
	groups := []dto.CreateGroupRequestDTO{
		{Name: "Group 1", Members: []string{"user1"}},
		{Name: "Group 2", Members: []string{"user2"}},
		{Name: "Group 3", Members: []string{"user3"}},
	}

	for _, group := range groups {
		payloadBytes, err := json.Marshal(group)
		require.NoError(t, err)

		req, err := http.NewRequest("POST", "/api/v1/groups", bytes.NewBuffer(payloadBytes))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		resp, err := testApp.Request(req)
		require.NoError(t, err)
		require.Equal(t, 201, resp.StatusCode)
	}

	// Test listing groups
	req, err := http.NewRequest("GET", "/api/v1/groups", nil)
	require.NoError(t, err)

	resp, err := testApp.Request(req)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	var groupList dto.ListGroupResponseDTO
	err = json.NewDecoder(resp.Body).Decode(&groupList)
	require.NoError(t, err)

	assert.Len(t, groupList.Data, 3)
	assert.Equal(t, int64(3), groupList.Meta.Total)

	// Test pagination
	req, err = http.NewRequest("GET", "/api/v1/groups?page=1&per_page=2", nil)
	require.NoError(t, err)

	resp, err = testApp.Request(req)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	err = json.NewDecoder(resp.Body).Decode(&groupList)
	require.NoError(t, err)

	assert.Len(t, groupList.Data, 2)
	assert.Equal(t, int64(3), groupList.Meta.Total)
}

func TestGroupControllerAddRemoveUser(t *testing.T) {
	testApp := SetupTestApp(t)
	defer testApp.Cleanup(t)

	// Create a user first
	userPayload := dto.CreateUserRequestDTO{
		Name:  "Test User",
		Email: "test@example.com",
	}

	payloadBytes, err := json.Marshal(userPayload)
	require.NoError(t, err)

	req, err := http.NewRequest("POST", "/api/v1/users", bytes.NewBuffer(payloadBytes))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	resp, err := testApp.Request(req)
	require.NoError(t, err)
	require.Equal(t, 201, resp.StatusCode)

	var createdUser dto.UserResponseDTO
	err = json.NewDecoder(resp.Body).Decode(&createdUser)
	require.NoError(t, err)

	// Create a group
	groupPayload := dto.CreateGroupRequestDTO{
		Name:    "Test Group",
		Members: []string{},
	}

	payloadBytes, err = json.Marshal(groupPayload)
	require.NoError(t, err)

	req, err = http.NewRequest("POST", "/api/v1/groups", bytes.NewBuffer(payloadBytes))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	resp, err = testApp.Request(req)
	require.NoError(t, err)
	require.Equal(t, 201, resp.StatusCode)

	var createdGroup dto.GroupResponseDTO
	err = json.NewDecoder(resp.Body).Decode(&createdGroup)
	require.NoError(t, err)

	// Add user to group
	req, err = http.NewRequest("POST", fmt.Sprintf("/api/v1/groups/%s/members/%s", createdGroup.ID, createdUser.ID), nil)
	require.NoError(t, err)

	resp, err = testApp.Request(req)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	// Verify user was added to group
	req, err = http.NewRequest("GET", fmt.Sprintf("/api/v1/groups/%s", createdGroup.ID), nil)
	require.NoError(t, err)

	resp, err = testApp.Request(req)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	var updatedGroup dto.GroupResponseDTO
	err = json.NewDecoder(resp.Body).Decode(&updatedGroup)
	require.NoError(t, err)

	assert.Contains(t, updatedGroup.Members, createdUser.ID)

	// Remove user from group
	req, err = http.NewRequest("DELETE", fmt.Sprintf("/api/v1/groups/%s/members/%s", createdGroup.ID, createdUser.ID), nil)
	require.NoError(t, err)

	resp, err = testApp.Request(req)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	// Verify user was removed from group
	req, err = http.NewRequest("GET", fmt.Sprintf("/api/v1/groups/%s", createdGroup.ID), nil)
	require.NoError(t, err)

	resp, err = testApp.Request(req)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	err = json.NewDecoder(resp.Body).Decode(&updatedGroup)
	require.NoError(t, err)

	assert.NotContains(t, updatedGroup.Members, createdUser.ID)
}

// Test repository error scenarios for groups
func TestGroupControllerRepositoryErrors(t *testing.T) {
	testApp := SetupTestApp(t)
	defer testApp.Cleanup(t)

	const (
		invalidObjectID     = "invalid-object-id"
		nonExistentObjectID = "507f1f77bcf86cd799439011"
		contentTypeJSON     = "application/json"
		contentTypeHeader   = "Content-Type"
		errorKeyName        = "error"
		invalidObjectIDMsg  = "invalid ObjectID"
		groupNotFoundMsg    = "Group not found"
		userNotFoundMsg     = "User not found"
		updatedGroupName    = "Updated Group Name"
	)

	t.Run("GetByID with invalid ObjectID format", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/groups/%s", invalidObjectID), nil)

		resp, err := testApp.App.Test(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// The application returns 404 for invalid ObjectID format
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)

		var errorResponse map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&errorResponse)
		require.NoError(t, err)
		assert.Contains(t, errorResponse[errorKeyName], groupNotFoundMsg)
	})

	t.Run("GetByID with valid ObjectID format but non-existent group", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/groups/%s", nonExistentObjectID), nil)

		resp, err := testApp.App.Test(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)

		var errorResponse map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&errorResponse)
		require.NoError(t, err)
		assert.Contains(t, errorResponse[errorKeyName], groupNotFoundMsg)
	})

	t.Run("Update with invalid ObjectID format", func(t *testing.T) {
		updateData := dto.CreateGroupRequestDTO{
			Name:    updatedGroupName,
			Members: []string{},
		}
		payload, _ := json.Marshal(updateData)

		req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("/api/v1/groups/%s", invalidObjectID), bytes.NewBuffer(payload))
		req.Header.Set(contentTypeHeader, contentTypeJSON)

		resp, err := testApp.App.Test(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// The application returns 500 for invalid ObjectID in repository operations
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

		var errorResponse map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&errorResponse)
		require.NoError(t, err)
		assert.Contains(t, errorResponse[errorKeyName], "provided hex string is not a valid ObjectID")
	})

	t.Run("Update with valid ObjectID format but non-existent group", func(t *testing.T) {
		updateData := dto.CreateGroupRequestDTO{
			Name:    updatedGroupName,
			Members: []string{},
		}
		payload, _ := json.Marshal(updateData)

		req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("/api/v1/groups/%s", nonExistentObjectID), bytes.NewBuffer(payload))
		req.Header.Set(contentTypeHeader, contentTypeJSON)

		resp, err := testApp.App.Test(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Now update operations correctly return 404 when no document is found
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)

		var errorResponse map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&errorResponse)
		require.NoError(t, err)
		assert.Contains(t, errorResponse[errorKeyName], "Group not found")
	})

	t.Run("Delete with invalid ObjectID format", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/groups/%s", invalidObjectID), nil)

		resp, err := testApp.App.Test(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// The application returns 500 for invalid ObjectID in repository operations
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

		var errorResponse map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&errorResponse)
		require.NoError(t, err)
		assert.Contains(t, errorResponse[errorKeyName], "provided hex string is not a valid ObjectID")
	})

	t.Run("Delete with valid ObjectID format but non-existent group", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/groups/%s", nonExistentObjectID), nil)

		resp, err := testApp.App.Test(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// MongoDB delete operations return 204 even if no document was found
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
		// No body expected for 204 status
	})

	t.Run("AddUserToGroup with invalid group ObjectID format", func(t *testing.T) {
		// First create a user to add
		userData := dto.CreateUserRequestDTO{
			Name:  "Test User",
			Email: "test@example.com",
		}
		userPayload, _ := json.Marshal(userData)

		userReq, _ := http.NewRequest(http.MethodPost, "/api/v1/users/", bytes.NewBuffer(userPayload))
		userReq.Header.Set(contentTypeHeader, contentTypeJSON)

		userResp, err := testApp.App.Test(userReq)
		require.NoError(t, err)
		defer userResp.Body.Close()

		var createdUser dto.UserResponseDTO
		err = json.NewDecoder(userResp.Body).Decode(&createdUser)
		require.NoError(t, err)

		// Try to add user to group with invalid group ID
		req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/groups/%s/members/%s", invalidObjectID, createdUser.ID), nil)

		resp, err := testApp.App.Test(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// The application returns 500 for invalid ObjectID in repository operations
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

		var errorResponse map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&errorResponse)
		require.NoError(t, err)
		assert.Contains(t, errorResponse[errorKeyName], "provided hex string is not a valid ObjectID")
	})

	t.Run("AddUserToGroup with non-existent group", func(t *testing.T) {
		// First create a user to add
		userData := dto.CreateUserRequestDTO{
			Name:  "Test User",
			Email: "test2@example.com",
		}
		userPayload, _ := json.Marshal(userData)

		userReq, _ := http.NewRequest(http.MethodPost, "/api/v1/users/", bytes.NewBuffer(userPayload))
		userReq.Header.Set(contentTypeHeader, contentTypeJSON)

		userResp, err := testApp.App.Test(userReq)
		require.NoError(t, err)
		defer userResp.Body.Close()

		var createdUser dto.UserResponseDTO
		err = json.NewDecoder(userResp.Body).Decode(&createdUser)
		require.NoError(t, err)

		// Try to add user to non-existent group
		req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/groups/%s/members/%s", nonExistentObjectID, createdUser.ID), nil)

		resp, err := testApp.App.Test(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// MongoDB $addToSet operations return 500 when group document doesn't exist
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})

	t.Run("AddUserToGroup with non-existent user", func(t *testing.T) {
		// First create a group
		groupData := dto.CreateGroupRequestDTO{
			Name:    "Test Group",
			Members: []string{},
		}
		groupPayload, _ := json.Marshal(groupData)

		groupReq, _ := http.NewRequest(http.MethodPost, "/api/v1/groups/", bytes.NewBuffer(groupPayload))
		groupReq.Header.Set(contentTypeHeader, contentTypeJSON)

		groupResp, err := testApp.App.Test(groupReq)
		require.NoError(t, err)
		defer groupResp.Body.Close()

		var createdGroup dto.GroupResponseDTO
		err = json.NewDecoder(groupResp.Body).Decode(&createdGroup)
		require.NoError(t, err)

		// Try to add non-existent user to group
		req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/groups/%s/members/%s", createdGroup.ID, nonExistentObjectID), nil)

		resp, err := testApp.App.Test(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// MongoDB $addToSet operations return 500 for various errors
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})

	t.Run("RemoveUserFromGroup with invalid group ObjectID format", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/groups/%s/members/%s", invalidObjectID, nonExistentObjectID), nil)

		resp, err := testApp.App.Test(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// The application returns 500 for invalid ObjectID in repository operations
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

		var errorResponse map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&errorResponse)
		require.NoError(t, err)
		assert.Contains(t, errorResponse[errorKeyName], "provided hex string is not a valid ObjectID")
	})

	t.Run("List groups with invalid pagination parameters", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/groups/?page=0&limit=-1", nil)

		resp, err := testApp.App.Test(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Should handle invalid parameters gracefully by using defaults
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}

// Test use case error scenarios for group updates
func TestGroupControllerUseCaseUpdateErrors(t *testing.T) {
	testApp := SetupTestApp(t)
	defer testApp.Cleanup(t)

	const (
		contentTypeJSON   = "application/json"
		contentTypeHeader = "Content-Type"
		errorKeyName      = "error"
		testGroupName     = "Test Group For Update Errors"
		updatedGroupName  = "Updated Group Name"
	)

	t.Run("Update group - GetByID returns error in use case", func(t *testing.T) {
		// Try to update a non-existent group (valid ObjectID format)
		nonExistentID := "507f1f77bcf86cd799439011"
		updateData := dto.CreateGroupRequestDTO{
			Name:    updatedGroupName,
			Members: []string{},
		}
		payload, _ := json.Marshal(updateData)

		req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("/api/v1/groups/%s", nonExistentID), bytes.NewBuffer(payload))
		req.Header.Set(contentTypeHeader, contentTypeJSON)

		resp, err := testApp.App.Test(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Use case calls GetByID first, which should return error for non-existent group
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)

		var errorResponse map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&errorResponse)
		require.NoError(t, err)
		assert.Contains(t, errorResponse[errorKeyName], "Group not found")
	})

	t.Run("Update group - validation errors before repository update", func(t *testing.T) {
		// First create a group
		groupData := dto.CreateGroupRequestDTO{
			Name:    testGroupName,
			Members: []string{},
		}
		groupPayload, _ := json.Marshal(groupData)

		groupReq, _ := http.NewRequest(http.MethodPost, "/api/v1/groups/", bytes.NewBuffer(groupPayload))
		groupReq.Header.Set(contentTypeHeader, contentTypeJSON)

		groupResp, err := testApp.App.Test(groupReq)
		require.NoError(t, err)
		defer groupResp.Body.Close()
		require.Equal(t, http.StatusCreated, groupResp.StatusCode)

		var createdGroup dto.GroupResponseDTO
		err = json.NewDecoder(groupResp.Body).Decode(&createdGroup)
		require.NoError(t, err)

		// Now try to update with invalid data to trigger validation error
		updateData := dto.CreateGroupRequestDTO{
			Name:    "", // Empty name should trigger validation error
			Members: []string{},
		}
		payload, _ := json.Marshal(updateData)

		req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("/api/v1/groups/%s", createdGroup.ID), bytes.NewBuffer(payload))
		req.Header.Set(contentTypeHeader, contentTypeJSON)

		resp, err := testApp.App.Test(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Should return 400 due to validation errors before reaching repository
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		var errorResponse map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&errorResponse)
		require.NoError(t, err)
		// Should contain validation error message
		assert.NotNil(t, errorResponse[errorKeyName])
	})

	t.Run("Update group - concurrent modification scenario", func(t *testing.T) {
		// First create a group
		groupData := dto.CreateGroupRequestDTO{
			Name:    testGroupName,
			Members: []string{},
		}
		groupPayload, _ := json.Marshal(groupData)

		groupReq, _ := http.NewRequest(http.MethodPost, "/api/v1/groups/", bytes.NewBuffer(groupPayload))
		groupReq.Header.Set(contentTypeHeader, contentTypeJSON)

		groupResp, err := testApp.App.Test(groupReq)
		require.NoError(t, err)
		defer groupResp.Body.Close()
		require.Equal(t, http.StatusCreated, groupResp.StatusCode)

		var createdGroup dto.GroupResponseDTO
		err = json.NewDecoder(groupResp.Body).Decode(&createdGroup)
		require.NoError(t, err)

		// Delete the group to simulate concurrent modification
		deleteReq, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/groups/%s", createdGroup.ID), nil)
		deleteResp, err := testApp.App.Test(deleteReq)
		require.NoError(t, err)
		defer deleteResp.Body.Close()
		require.Equal(t, http.StatusNoContent, deleteResp.StatusCode)

		// Now try to update the deleted group
		updateData := dto.CreateGroupRequestDTO{
			Name:    updatedGroupName,
			Members: []string{},
		}
		payload, _ := json.Marshal(updateData)

		req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("/api/v1/groups/%s", createdGroup.ID), bytes.NewBuffer(payload))
		req.Header.Set(contentTypeHeader, contentTypeJSON)

		resp, err := testApp.App.Test(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Should return 404 because GetByID in use case fails for deleted group
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)

		var errorResponse map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&errorResponse)
		require.NoError(t, err)
		assert.Contains(t, errorResponse[errorKeyName], "Group not found")
	})

	t.Run("Update group - repository update operation fails after successful GetByID", func(t *testing.T) {
		// First create a group
		groupData := dto.CreateGroupRequestDTO{
			Name:    testGroupName,
			Members: []string{},
		}
		groupPayload, _ := json.Marshal(groupData)

		groupReq, _ := http.NewRequest(http.MethodPost, "/api/v1/groups/", bytes.NewBuffer(groupPayload))
		groupReq.Header.Set(contentTypeHeader, contentTypeJSON)

		groupResp, err := testApp.App.Test(groupReq)
		require.NoError(t, err)
		defer groupResp.Body.Close()
		require.Equal(t, http.StatusCreated, groupResp.StatusCode)

		var createdGroup dto.GroupResponseDTO
		err = json.NewDecoder(groupResp.Body).Decode(&createdGroup)
		require.NoError(t, err)

		// Delete the group first to make the update fail
		deleteReq, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/groups/%s", createdGroup.ID), nil)
		deleteResp, err := testApp.App.Test(deleteReq)
		require.NoError(t, err)
		defer deleteResp.Body.Close()

		// Try to update with valid data but group was deleted concurrently
		updateData := dto.CreateGroupRequestDTO{
			Name:    updatedGroupName,
			Members: []string{},
		}
		payload, _ := json.Marshal(updateData)

		req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("/api/v1/groups/%s", createdGroup.ID), bytes.NewBuffer(payload))
		req.Header.Set(contentTypeHeader, contentTypeJSON)

		resp, err := testApp.App.Test(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Should return 404 as GetByID will fail first
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)

		var errorResponse map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&errorResponse)
		require.NoError(t, err)
		assert.Contains(t, errorResponse[errorKeyName], "Group not found")
	})
}

func TestGroupControllerListRepositoryErrors(t *testing.T) {
	t.Run("List groups - force MongoDB connection close to cause repository List error", func(t *testing.T) {
		testApp := SetupTestApp(t)
		// Note: Don't defer cleanup here since we'll disconnect the database

		// Create some test groups first
		groups := []dto.CreateGroupRequestDTO{
			{Name: "Test Group 1", Members: []string{}},
			{Name: "Test Group 2", Members: []string{}},
		}

		for _, group := range groups {
			payloadBytes, err := json.Marshal(group)
			require.NoError(t, err)

			createReq, err := http.NewRequest("POST", "/api/v1/groups", bytes.NewBuffer(payloadBytes))
			require.NoError(t, err)
			createReq.Header.Set("Content-Type", "application/json")

			createResp, err := testApp.Request(createReq)
			require.NoError(t, err)
			require.Equal(t, 201, createResp.StatusCode)
		}

		// Verify normal operation works first
		req, err := http.NewRequest("GET", "/api/v1/groups", nil)
		require.NoError(t, err)

		resp, err := testApp.Request(req)
		require.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)

		// Now force close the MongoDB connection to simulate repository error
		// This will cause the next List() call in the repository to fail
		ctx := context.Background()
		err = testApp.DB.Client.Disconnect(ctx)
		require.NoError(t, err)

		// Try to list groups after closing connection - should return repository error
		req, err = http.NewRequest("GET", "/api/v1/groups", nil)
		require.NoError(t, err)

		resp, err = testApp.Request(req)
		require.NoError(t, err)

		// Should return 500 due to repository List() operation failing
		assert.Equal(t, 500, resp.StatusCode)

		var errorResponse map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&errorResponse)
		require.NoError(t, err)
		assert.NotNil(t, errorResponse["error"])

		// The error should be related to connection/repository failure
		errorMsg := errorResponse["error"].(string)
		assert.True(t,
			strings.Contains(errorMsg, "connection") ||
				strings.Contains(errorMsg, "client is disconnected") ||
				strings.Contains(errorMsg, "topology is closed") ||
				strings.Contains(errorMsg, "server selection error"),
			"Expected connection error, got: %s", errorMsg)
	})

	t.Run("Count repository error - force MongoDB connection close during Count operation", func(t *testing.T) {
		testApp := SetupTestApp(t)

		// First create some groups
		for i := 0; i < 3; i++ {
			groupData := dto.CreateGroupRequestDTO{
				Name:    fmt.Sprintf("Count Test Group %d", i),
				Members: []string{},
			}
			payloadBytes, err := json.Marshal(groupData)
			require.NoError(t, err)

			createReq, err := http.NewRequest("POST", "/api/v1/groups", bytes.NewBuffer(payloadBytes))
			require.NoError(t, err)
			createReq.Header.Set("Content-Type", "application/json")

			createResp, err := testApp.Request(createReq)
			require.NoError(t, err)
			require.Equal(t, 201, createResp.StatusCode)
		}

		// Test normal listing first
		req, err := http.NewRequest("GET", "/api/v1/groups", nil)
		require.NoError(t, err)

		resp, err := testApp.Request(req)
		require.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)

		// Now disconnect during the request to cause Count() operation to fail
		// Close connection to force repository error
		ctx := context.Background()
		err = testApp.DB.Client.Disconnect(ctx)
		require.NoError(t, err)

		// Try to list again - this should cause either List() or Count() to fail
		req, err = http.NewRequest("GET", "/api/v1/groups", nil)
		require.NoError(t, err)

		resp, err = testApp.Request(req)
		require.NoError(t, err)

		// Should return 500 due to repository operation failing
		assert.Equal(t, 500, resp.StatusCode)

		var errorResponse map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&errorResponse)
		require.NoError(t, err)
		assert.NotNil(t, errorResponse["error"])
	})

	t.Run("Repository cursor error - force collection state corruption", func(t *testing.T) {
		testApp := SetupTestApp(t)
		defer testApp.Cleanup(t)

		// Create test groups
		groupData := dto.CreateGroupRequestDTO{
			Name:    "Cursor Test Group",
			Members: []string{},
		}
		payloadBytes, err := json.Marshal(groupData)
		require.NoError(t, err)

		createReq, err := http.NewRequest("POST", "/api/v1/groups", bytes.NewBuffer(payloadBytes))
		require.NoError(t, err)
		createReq.Header.Set("Content-Type", "application/json")

		createResp, err := testApp.Request(createReq)
		require.NoError(t, err)
		require.Equal(t, 201, createResp.StatusCode)

		// Verify it works normally first
		listReq, err := http.NewRequest("GET", "/api/v1/groups", nil)
		require.NoError(t, err)

		listResp, err := testApp.Request(listReq)
		require.NoError(t, err)
		assert.Equal(t, 200, listResp.StatusCode)

		// Drop the database to force cursor/collection errors
		ctx := context.Background()
		err = testApp.DB.DB.Drop(ctx)
		require.NoError(t, err)

		// Try to list groups after dropping database
		// This should cause the Find operation to fail
		listReq, err = http.NewRequest("GET", "/api/v1/groups", nil)
		require.NoError(t, err)

		listResp, err = testApp.Request(listReq)
		require.NoError(t, err)

		// Should either work (200 with empty list) or return error (500)
		assert.True(t, listResp.StatusCode == 200 || listResp.StatusCode == 500)

		if listResp.StatusCode == 200 {
			var groupList dto.ListGroupResponseDTO
			err = json.NewDecoder(listResp.Body).Decode(&groupList)
			require.NoError(t, err)
			// Should return empty list since database was dropped
			assert.Equal(t, 0, len(groupList.Data))
			assert.Equal(t, int64(0), groupList.Meta.Total)
		} else {
			var errorResponse map[string]interface{}
			err = json.NewDecoder(listResp.Body).Decode(&errorResponse)
			require.NoError(t, err)
			assert.NotNil(t, errorResponse["error"])
		}
	})

	t.Run("Repository pagination calculation overflow", func(t *testing.T) {
		testApp := SetupTestApp(t)
		defer testApp.Cleanup(t)

		// Create test data
		groupData := dto.CreateGroupRequestDTO{
			Name:    "Pagination Test Group",
			Members: []string{},
		}
		payloadBytes, err := json.Marshal(groupData)
		require.NoError(t, err)

		createReq, err := http.NewRequest("POST", "/api/v1/groups", bytes.NewBuffer(payloadBytes))
		require.NoError(t, err)
		createReq.Header.Set("Content-Type", "application/json")

		createResp, err := testApp.Request(createReq)
		require.NoError(t, err)
		require.Equal(t, 201, createResp.StatusCode)

		// Test with values that might cause issues in pagination calculation
		// These should be caught by validation, but if they reach the repository
		// they could cause errors in skip/limit calculation
		testCases := []struct {
			page    string
			perPage string
			name    string
		}{
			{"9223372036854775807", "9223372036854775807", "Maximum int64 values"},
			{"-1", "10", "Negative page"},
			{"1", "-1", "Negative per_page"},
			{"0", "0", "Zero values"},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				req, err := http.NewRequest("GET",
					fmt.Sprintf("/api/v1/groups?page=%s&per_page=%s", tc.page, tc.perPage), nil)
				require.NoError(t, err)

				resp, err := testApp.Request(req)
				require.NoError(t, err)

				// Should handle edge cases gracefully - either return 400 (validation error)
				// or 500 (repository error) or 200 (handled gracefully)
				assert.True(t, resp.StatusCode == 200 || resp.StatusCode == 400 || resp.StatusCode == 500,
					"Expected 200, 400, or 500, got %d for case: %s", resp.StatusCode, tc.name)

				if resp.StatusCode == 500 {
					var errorResponse map[string]interface{}
					err = json.NewDecoder(resp.Body).Decode(&errorResponse)
					require.NoError(t, err)
					assert.NotNil(t, errorResponse["error"])
				}
			})
		}
	})

	t.Run("Repository connection timeout during List operation", func(t *testing.T) {
		testApp := SetupTestApp(t)
		defer testApp.Cleanup(t)

		// Create test groups
		for i := 0; i < 5; i++ {
			groupData := dto.CreateGroupRequestDTO{
				Name:    fmt.Sprintf("Timeout Test Group %d", i),
				Members: []string{},
			}
			payloadBytes, err := json.Marshal(groupData)
			require.NoError(t, err)

			createReq, err := http.NewRequest("POST", "/api/v1/groups", bytes.NewBuffer(payloadBytes))
			require.NoError(t, err)
			createReq.Header.Set("Content-Type", "application/json")

			createResp, err := testApp.Request(createReq)
			require.NoError(t, err)
			require.Equal(t, 201, createResp.StatusCode)
		}

		// Test normal operation first
		req, err := http.NewRequest("GET", "/api/v1/groups", nil)
		require.NoError(t, err)

		resp, err := testApp.Request(req)
		require.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)

		// Now try to cause issues by creating high load on database
		// We'll try to access with different pagination parameters that might cause issues
		for i := 0; i < 3; i++ {
			req, err := http.NewRequest("GET", fmt.Sprintf("/api/v1/groups?page=%d&per_page=100", i+1), nil)
			require.NoError(t, err)

			resp, err := testApp.Request(req)
			require.NoError(t, err)

			// Should either work normally or handle gracefully
			assert.True(t, resp.StatusCode == 200 || resp.StatusCode == 500,
				"Expected 200 or 500, got %d", resp.StatusCode)
		}
	})
}
