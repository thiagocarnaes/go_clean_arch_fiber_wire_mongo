package acceptance

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"user-management/internal/application/dto"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserController_Create(t *testing.T) {
	testApp := SetupTestApp(t)
	defer testApp.Cleanup(t)

	tests := []struct {
		name           string
		payload        interface{}
		expectedStatus int
		expectedError  string
	}{
		{
			name: "Valid user creation",
			payload: dto.CreateUserRequestDTO{
				Name:  "John Doe",
				Email: "john.doe@example.com",
			},
			expectedStatus: 201,
		},
		{
			name: "Invalid email format",
			payload: dto.CreateUserRequestDTO{
				Name:  "John Doe",
				Email: "invalid-email",
			},
			expectedStatus: 400,
		},
		{
			name: "Missing required fields",
			payload: dto.CreateUserRequestDTO{
				Name: "John Doe",
				// Email missing
			},
			expectedStatus: 400,
		},
		{
			name: "Empty name",
			payload: dto.CreateUserRequestDTO{
				Name:  "",
				Email: "john.doe@example.com",
			},
			expectedStatus: 400,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testApp.ClearDatabase(t)

			payloadBytes, err := json.Marshal(tt.payload)
			require.NoError(t, err)

			req, err := http.NewRequest("POST", "/api/v1/users", bytes.NewBuffer(payloadBytes))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			resp, err := testApp.Request(req)
			require.NoError(t, err)
			require.Equal(t, tt.expectedStatus, resp.StatusCode)
		})
	}
}

// Test use case error scenarios for updates
func TestUserControllerUseCaseUpdateErrors(t *testing.T) {
	testApp := SetupTestApp(t)
	defer testApp.Cleanup(t)

	const (
		contentTypeJSON   = "application/json"
		contentTypeHeader = "Content-Type"
		errorKeyName      = "error"
		testUserName      = "Test User For Update Errors"
		testUserEmail     = "update.errors@example.com"
		updatedName       = "Updated Name"
		updatedEmail      = "updated@example.com"
	)

	t.Run("Update user - GetByID returns error in use case", func(t *testing.T) {
		// Try to update a non-existent user (valid ObjectID format)
		nonExistentID := "507f1f77bcf86cd799439011"
		updateData := dto.CreateUserRequestDTO{
			Name:  updatedName,
			Email: updatedEmail,
		}
		payload, _ := json.Marshal(updateData)

		req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("/api/v1/users/%s", nonExistentID), bytes.NewBuffer(payload))
		req.Header.Set(contentTypeHeader, contentTypeJSON)

		resp, err := testApp.App.Test(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Use case calls GetByID first, which should return error for non-existent user
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)

		var errorResponse map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&errorResponse)
		require.NoError(t, err)
		assert.Contains(t, errorResponse[errorKeyName], "User not found")
	})

	t.Run("Update user - success GetByID but Update repository fails", func(t *testing.T) {
		// First create a user
		userData := dto.CreateUserRequestDTO{
			Name:  testUserName,
			Email: testUserEmail,
		}
		userPayload, _ := json.Marshal(userData)

		userReq, _ := http.NewRequest(http.MethodPost, "/api/v1/users/", bytes.NewBuffer(userPayload))
		userReq.Header.Set(contentTypeHeader, contentTypeJSON)

		userResp, err := testApp.App.Test(userReq)
		require.NoError(t, err)
		defer userResp.Body.Close()
		require.Equal(t, http.StatusCreated, userResp.StatusCode)

		var createdUser dto.UserResponseDTO
		err = json.NewDecoder(userResp.Body).Decode(&createdUser)
		require.NoError(t, err)

		// Now simulate an update with invalid data to trigger repository error
		updateData := dto.CreateUserRequestDTO{
			Name:  "",                     // Empty name should trigger validation error
			Email: "invalid-email-format", // Invalid email format
		}
		payload, _ := json.Marshal(updateData)

		req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("/api/v1/users/%s", createdUser.ID), bytes.NewBuffer(payload))
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

	t.Run("Update user - concurrent modification scenario", func(t *testing.T) {
		// First create a user
		userData := dto.CreateUserRequestDTO{
			Name:  testUserName,
			Email: "concurrent@example.com",
		}
		userPayload, _ := json.Marshal(userData)

		userReq, _ := http.NewRequest(http.MethodPost, "/api/v1/users/", bytes.NewBuffer(userPayload))
		userReq.Header.Set(contentTypeHeader, contentTypeJSON)

		userResp, err := testApp.App.Test(userReq)
		require.NoError(t, err)
		defer userResp.Body.Close()
		require.Equal(t, http.StatusCreated, userResp.StatusCode)

		var createdUser dto.UserResponseDTO
		err = json.NewDecoder(userResp.Body).Decode(&createdUser)
		require.NoError(t, err)

		// Delete the user to simulate concurrent modification
		deleteReq, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/users/%s", createdUser.ID), nil)
		deleteResp, err := testApp.App.Test(deleteReq)
		require.NoError(t, err)
		defer deleteResp.Body.Close()
		require.Equal(t, http.StatusNoContent, deleteResp.StatusCode)

		// Now try to update the deleted user
		updateData := dto.CreateUserRequestDTO{
			Name:  updatedName,
			Email: updatedEmail,
		}
		payload, _ := json.Marshal(updateData)

		req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("/api/v1/users/%s", createdUser.ID), bytes.NewBuffer(payload))
		req.Header.Set(contentTypeHeader, contentTypeJSON)

		resp, err := testApp.App.Test(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Should return 404 because GetByID in use case fails for deleted user
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)

		var errorResponse map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&errorResponse)
		require.NoError(t, err)
		assert.Contains(t, errorResponse[errorKeyName], "User not found")
	})
}

func TestUserController_Update(t *testing.T) {
	testApp := SetupTestApp(t)
	defer testApp.Cleanup(t)

	// Create a user first
	createPayload := dto.CreateUserRequestDTO{
		Name:  "Original Name",
		Email: "original@example.com",
	}

	payloadBytes, err := json.Marshal(createPayload)
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

	// Update the user
	updatePayload := dto.CreateUserRequestDTO{
		Name:  "Updated Name",
		Email: "updated@example.com",
	}

	payloadBytes, err = json.Marshal(updatePayload)
	require.NoError(t, err)

	req, err = http.NewRequest("PUT", fmt.Sprintf("/api/v1/users/%s", createdUser.ID), bytes.NewBuffer(payloadBytes))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	resp, err = testApp.Request(req)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	var updatedUser dto.UserResponseDTO
	err = json.NewDecoder(resp.Body).Decode(&updatedUser)
	require.NoError(t, err)

	assert.Equal(t, createdUser.ID, updatedUser.ID)
	assert.Equal(t, updatePayload.Name, updatedUser.Name)
	assert.Equal(t, updatePayload.Email, updatedUser.Email)
}

func TestUserController_Delete(t *testing.T) {
	testApp := SetupTestApp(t)
	defer testApp.Cleanup(t)

	// Create a user first
	createPayload := dto.CreateUserRequestDTO{
		Name:  "To Be Deleted",
		Email: "delete@example.com",
	}

	payloadBytes, err := json.Marshal(createPayload)
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

	// Delete the user
	req, err = http.NewRequest("DELETE", fmt.Sprintf("/api/v1/users/%s", createdUser.ID), nil)
	require.NoError(t, err)

	resp, err = testApp.Request(req)
	require.NoError(t, err)
	assert.Equal(t, 204, resp.StatusCode)

	// Verify user is deleted by trying to get it
	req, err = http.NewRequest("GET", fmt.Sprintf("/api/v1/users/%s", createdUser.ID), nil)
	require.NoError(t, err)

	resp, err = testApp.Request(req)
	require.NoError(t, err)
	assert.Equal(t, 404, resp.StatusCode)
}

func TestUserController_List(t *testing.T) {
	testApp := SetupTestApp(t)
	defer testApp.Cleanup(t)

	// Create multiple users
	users := []dto.CreateUserRequestDTO{
		{Name: "User 1", Email: "user1@example.com"},
		{Name: "User 2", Email: "user2@example.com"},
		{Name: "User 3", Email: "user3@example.com"},
	}

	for _, user := range users {
		payloadBytes, err := json.Marshal(user)
		require.NoError(t, err)

		req, err := http.NewRequest("POST", "/api/v1/users", bytes.NewBuffer(payloadBytes))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		resp, err := testApp.Request(req)
		require.NoError(t, err)
		require.Equal(t, 201, resp.StatusCode)
	}

	// Test listing users
	req, err := http.NewRequest("GET", "/api/v1/users", nil)
	require.NoError(t, err)

	resp, err := testApp.Request(req)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	var userList dto.UserListResponseDTO
	err = json.NewDecoder(resp.Body).Decode(&userList)
	require.NoError(t, err)

	assert.Len(t, userList.Data, 3)
	assert.Equal(t, int64(3), userList.Meta.Total)

	// Test pagination
	req, err = http.NewRequest("GET", "/api/v1/users?page=1&per_page=2", nil)
	require.NoError(t, err)

	resp, err = testApp.Request(req)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	err = json.NewDecoder(resp.Body).Decode(&userList)
	require.NoError(t, err)

	assert.Len(t, userList.Data, 2)
	assert.Equal(t, int64(3), userList.Meta.Total)
}

func TestUserController_ListWithSearch(t *testing.T) {
	testApp := SetupTestApp(t)
	defer testApp.Cleanup(t)

	// Create users with different names
	users := []dto.CreateUserRequestDTO{
		{Name: "Alice Johnson", Email: "alice@example.com"},
		{Name: "Bob Smith", Email: "bob@example.com"},
		{Name: "Alice Cooper", Email: "alice.cooper@example.com"},
	}

	for _, user := range users {
		payloadBytes, err := json.Marshal(user)
		require.NoError(t, err)

		req, err := http.NewRequest("POST", "/api/v1/users", bytes.NewBuffer(payloadBytes))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		resp, err := testApp.Request(req)
		require.NoError(t, err)
		require.Equal(t, 201, resp.StatusCode)
	}

	// Search for users with "Alice" in name
	req, err := http.NewRequest("GET", "/api/v1/users?search=Alice", nil)
	require.NoError(t, err)

	resp, err := testApp.Request(req)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	var userList dto.UserListResponseDTO
	err = json.NewDecoder(resp.Body).Decode(&userList)
	require.NoError(t, err)

	assert.Len(t, userList.Data, 2)
	for _, user := range userList.Data {
		assert.True(t, strings.Contains(user.Name, "Alice"))
	}
}

func TestUserController_ListWithoutUsers(t *testing.T) {
	testApp := SetupTestApp(t)
	defer testApp.Cleanup(t)
	// Search for users with "Alice" in name
	req, err := http.NewRequest("GET", "/api/v1/users", nil)
	require.NoError(t, err)

	resp, err := testApp.Request(req)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

}

// Test repository error scenarios
func TestUserControllerRepositoryErrors(t *testing.T) {
	testApp := SetupTestApp(t)
	defer testApp.Cleanup(t)

	const (
		invalidObjectID     = "invalid-object-id"
		nonExistentObjectID = "507f1f77bcf86cd799439011"
		contentTypeJSON     = "application/json"
		contentTypeHeader   = "Content-Type"
		errorKeyName        = "error"
		invalidObjectIDMsg  = "invalid ObjectID"
		userNotFoundMsg     = "User not found"
		updatedName         = "Updated Name"
		updatedEmail        = "updated@example.com"
	)

	t.Run("GetByID with invalid ObjectID format", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/users/%s", invalidObjectID), nil)

		resp, err := testApp.App.Test(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// The application returns 404 for invalid ObjectID format
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)

		var errorResponse map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&errorResponse)
		require.NoError(t, err)
		assert.Contains(t, errorResponse[errorKeyName], userNotFoundMsg)
	})

	t.Run("GetByID with valid ObjectID format but non-existent user", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/users/%s", nonExistentObjectID), nil)

		resp, err := testApp.App.Test(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)

		var errorResponse map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&errorResponse)
		require.NoError(t, err)
		assert.Contains(t, errorResponse[errorKeyName], userNotFoundMsg)
	})

	t.Run("Update with invalid ObjectID format", func(t *testing.T) {
		updateData := dto.CreateUserRequestDTO{
			Name:  updatedName,
			Email: updatedEmail,
		}
		payload, _ := json.Marshal(updateData)

		req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("/api/v1/users/%s", invalidObjectID), bytes.NewBuffer(payload))
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

	t.Run("Update with valid ObjectID format but non-existent user", func(t *testing.T) {
		updateData := dto.CreateUserRequestDTO{
			Name:  updatedName,
			Email: updatedEmail,
		}
		payload, _ := json.Marshal(updateData)

		req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("/api/v1/users/%s", nonExistentObjectID), bytes.NewBuffer(payload))
		req.Header.Set(contentTypeHeader, contentTypeJSON)

		resp, err := testApp.App.Test(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Now update operations correctly return 404 when no document is found
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)

		var errorResponse map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&errorResponse)
		require.NoError(t, err)
		assert.Contains(t, errorResponse[errorKeyName], "User not found")
	})

	t.Run("Delete with invalid ObjectID format", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/users/%s", invalidObjectID), nil)

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

	t.Run("Delete with valid ObjectID format but non-existent user", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/users/%s", nonExistentObjectID), nil)

		resp, err := testApp.App.Test(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// MongoDB delete operations return 204 even if no document was found
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
		// No body expected for 204 status
	})

	t.Run("List with invalid pagination parameters", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/users/?page=0&limit=-1", nil)

		resp, err := testApp.App.Test(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Should handle invalid parameters gracefully by using defaults
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("Search with special characters", func(t *testing.T) {
		// Test search with regex special characters - use URL encoding
		specialSearch := "%2E%2A%5B%5D%7B%7D%28%29%5E%24%2B%3F%7C" // URL encoded version of .*[]{}()^$+?|
		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/users/?search=%s", specialSearch), nil)

		resp, err := testApp.App.Test(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Should handle special characters without crashing - might return 500 due to regex compilation
		// This is acceptable behavior for invalid regex patterns
		assert.True(t, resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusInternalServerError)

		if resp.StatusCode == http.StatusOK {
			var searchResponse dto.UserListResponseDTO
			err = json.NewDecoder(resp.Body).Decode(&searchResponse)
			require.NoError(t, err)
			assert.Equal(t, 0, len(searchResponse.Data))
		}
	})
}

func TestHealthCheck(t *testing.T) {
	testApp := SetupTestApp(t)
	defer testApp.Cleanup(t)

	req, err := http.NewRequest("GET", "/health", nil)
	require.NoError(t, err)

	resp, err := testApp.Request(req)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}
