package acceptance

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
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
