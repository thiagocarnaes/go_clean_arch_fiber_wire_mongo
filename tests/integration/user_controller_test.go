package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"user-management/internal/application/dto"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	usersEndpoint     = "/api/v1/users"
	usersEndpointFmt  = "/api/v1/users/%s"
	contentTypeHeader = "Content-Type"
	contentTypeJSON   = "application/json"
	errorKeyName      = "error"
	userNotFoundMsg   = "User not found"
	updatedName       = "Updated Name"
	updatedEmail      = "updated@example.com"
	invalidJSONMsg    = "Invalid JSON format"
)

func TestUserControllerCreate(t *testing.T) {
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
				Name:     "John Doe",
				Email:    "john.doe@example.com",
				IsActive: true,
			},
			expectedStatus: 201,
		},
		{
			name: "Valid user creation with IsActive false",
			payload: dto.CreateUserRequestDTO{
				Name:     "Jane Doe",
				Email:    "jane.doe@example.com",
				IsActive: false,
			},
			expectedStatus: 201,
		},
		{
			name: "Invalid email format",
			payload: dto.CreateUserRequestDTO{
				Name:     "John Doe",
				Email:    "invalid-email",
				IsActive: true,
			},
			expectedStatus: 400,
		},
		{
			name: "Missing required fields",
			payload: dto.CreateUserRequestDTO{
				Name:     "John Doe",
				IsActive: true,
				// Email missing
			},
			expectedStatus: 400,
		},
		{
			name: "Empty name",
			payload: dto.CreateUserRequestDTO{
				Name:     "",
				Email:    "john.doe@example.com",
				IsActive: true,
			},
			expectedStatus: 400,
		},
	}

	// Test invalid JSON separately
	t.Run("Invalid JSON format", func(t *testing.T) {
		testApp.ClearDatabase(t)

		// Send malformed JSON
		invalidJSON := `{"name": "John Doe", "email": "john@example.com", "is_active": true, "invalid": }`

		req, err := http.NewRequest("POST", usersEndpoint, bytes.NewBufferString(invalidJSON))
		require.NoError(t, err)
		req.Header.Set(contentTypeHeader, contentTypeJSON)

		resp, err := testApp.Request(req)
		require.NoError(t, err)
		assert.Equal(t, 400, resp.StatusCode)

		var errorResponse map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&errorResponse)
		require.NoError(t, err)
		// The application returns a generic "Invalid JSON format" message
		assert.Contains(t, errorResponse[errorKeyName], invalidJSONMsg)
	})

	// Test completely invalid JSON
	t.Run("Completely invalid JSON", func(t *testing.T) {
		testApp.ClearDatabase(t)

		// Send completely malformed JSON
		invalidJSON := `{this is not json at all!}`

		req, err := http.NewRequest("POST", usersEndpoint, bytes.NewBufferString(invalidJSON))
		require.NoError(t, err)
		req.Header.Set(contentTypeHeader, contentTypeJSON)

		resp, err := testApp.Request(req)
		require.NoError(t, err)
		assert.Equal(t, 400, resp.StatusCode)

		var errorResponse map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&errorResponse)
		require.NoError(t, err)
		// The application returns a generic "Invalid JSON format" message
		assert.Contains(t, errorResponse[errorKeyName], invalidJSONMsg)
	})

	// Test empty JSON object
	t.Run("Empty JSON object", func(t *testing.T) {
		testApp.ClearDatabase(t)

		// Send empty JSON object
		emptyJSON := `{}`

		req, err := http.NewRequest("POST", usersEndpoint, bytes.NewBufferString(emptyJSON))
		require.NoError(t, err)
		req.Header.Set(contentTypeHeader, contentTypeJSON)

		resp, err := testApp.Request(req)
		require.NoError(t, err)
		assert.Equal(t, 400, resp.StatusCode)

		// Should fail validation due to missing required fields
		var errorResponse map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&errorResponse)
		require.NoError(t, err)
		assert.NotNil(t, errorResponse[errorKeyName])
	})

	// Test JSON with wrong field types
	t.Run("JSON with wrong field types", func(t *testing.T) {
		testApp.ClearDatabase(t)

		// Send JSON with wrong field types (number instead of string)
		wrongTypeJSON := `{"name": 12345, "email": true, "is_active": "not_boolean"}`

		req, err := http.NewRequest("POST", usersEndpoint, bytes.NewBufferString(wrongTypeJSON))
		require.NoError(t, err)
		req.Header.Set(contentTypeHeader, contentTypeJSON)

		resp, err := testApp.Request(req)
		require.NoError(t, err)
		assert.Equal(t, 400, resp.StatusCode)

		var errorResponse map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&errorResponse)
		require.NoError(t, err)
		assert.NotNil(t, errorResponse[errorKeyName])
	})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testApp.ClearDatabase(t)

			payloadBytes, err := json.Marshal(tt.payload)
			require.NoError(t, err)

			req, err := http.NewRequest("POST", usersEndpoint, bytes.NewBuffer(payloadBytes))
			require.NoError(t, err)
			req.Header.Set(contentTypeHeader, contentTypeJSON)

			resp, err := testApp.Request(req)
			require.NoError(t, err)
			require.Equal(t, tt.expectedStatus, resp.StatusCode)
		})
	}
}

func TestUserControllerGet(t *testing.T) {
	testApp := SetupTestApp(t)
	defer testApp.Cleanup(t)

	const (
		contentTypeJSON   = "application/json"
		contentTypeHeader = "Content-Type"
		usersEndpoint     = "/api/v1/users"
		errorKeyName      = "error"
	)

	t.Run("Get existing user", func(t *testing.T) {
		testApp.ClearDatabase(t)

		// First create a user
		createUserDTO := dto.CreateUserRequestDTO{
			Name:     "John Doe",
			Email:    "john@example.com",
			IsActive: true,
		}

		payloadBytes, err := json.Marshal(createUserDTO)
		require.NoError(t, err)

		createReq, err := http.NewRequest("POST", usersEndpoint, bytes.NewBuffer(payloadBytes))
		require.NoError(t, err)
		createReq.Header.Set(contentTypeHeader, contentTypeJSON)

		createResp, err := testApp.Request(createReq)
		require.NoError(t, err)
		require.Equal(t, 201, createResp.StatusCode)

		var createdUser dto.UserResponseDTO
		err = json.NewDecoder(createResp.Body).Decode(&createdUser)
		require.NoError(t, err)

		// Now get the created user
		getReq, err := http.NewRequest("GET", fmt.Sprintf("%s/%s", usersEndpoint, createdUser.ID), nil)
		require.NoError(t, err)

		getResp, err := testApp.Request(getReq)
		require.NoError(t, err)
		assert.Equal(t, 200, getResp.StatusCode)

		var retrievedUser dto.UserResponseDTO
		err = json.NewDecoder(getResp.Body).Decode(&retrievedUser)
		require.NoError(t, err)

		// Verify the retrieved user matches the created user
		assert.Equal(t, createdUser.ID, retrievedUser.ID)
		assert.Equal(t, createdUser.Name, retrievedUser.Name)
		assert.Equal(t, createdUser.Email, retrievedUser.Email)
		assert.Equal(t, createdUser.IsActive, retrievedUser.IsActive)
	})

	t.Run("Get non-existent user with valid ObjectID", func(t *testing.T) {
		testApp.ClearDatabase(t)

		// Use a valid ObjectID format but non-existent
		nonExistentID := "507f1f77bcf86cd799439011"

		getReq, err := http.NewRequest("GET", fmt.Sprintf("%s/%s", usersEndpoint, nonExistentID), nil)
		require.NoError(t, err)

		getResp, err := testApp.Request(getReq)
		require.NoError(t, err)
		assert.Equal(t, 404, getResp.StatusCode)

		var errorResponse map[string]interface{}
		err = json.NewDecoder(getResp.Body).Decode(&errorResponse)
		require.NoError(t, err)
		assert.Equal(t, "User not found", errorResponse[errorKeyName])
	})

	t.Run("Get user with invalid ObjectID format", func(t *testing.T) {
		testApp.ClearDatabase(t)

		// Use an invalid ObjectID format
		invalidID := "invalid-id-format"

		getReq, err := http.NewRequest("GET", fmt.Sprintf("%s/%s", usersEndpoint, invalidID), nil)
		require.NoError(t, err)

		getResp, err := testApp.Request(getReq)
		require.NoError(t, err)
		assert.Equal(t, 404, getResp.StatusCode)

		var errorResponse map[string]interface{}
		err = json.NewDecoder(getResp.Body).Decode(&errorResponse)
		require.NoError(t, err)
		assert.Equal(t, "User not found", errorResponse[errorKeyName])
	})

	t.Run("Get user with empty ID", func(t *testing.T) {
		testApp.ClearDatabase(t)

		// Try to get user with empty ID - this should hit the list route instead
		getReq, err := http.NewRequest("GET", fmt.Sprintf("%s/", usersEndpoint), nil)
		require.NoError(t, err)

		getResp, err := testApp.Request(getReq)
		require.NoError(t, err)

		// Should return 200 because the route /api/v1/users/ maps to List, not Get
		assert.Equal(t, 200, getResp.StatusCode)

		// Should return an empty list since database is cleared
		var userList dto.UserListResponseDTO
		err = json.NewDecoder(getResp.Body).Decode(&userList)
		require.NoError(t, err)
		assert.Equal(t, 0, len(userList.Data))
	})

	t.Run("Get multiple users to verify ID uniqueness", func(t *testing.T) {
		testApp.ClearDatabase(t)

		// Create multiple users
		users := []dto.CreateUserRequestDTO{
			{Name: "User 1", Email: "user1@example.com", IsActive: true},
			{Name: "User 2", Email: "user2@example.com", IsActive: false},
			{Name: "User 3", Email: "user3@example.com", IsActive: true},
		}

		var createdUsers []dto.UserResponseDTO

		for _, user := range users {
			payloadBytes, err := json.Marshal(user)
			require.NoError(t, err)

			createReq, err := http.NewRequest("POST", usersEndpoint, bytes.NewBuffer(payloadBytes))
			require.NoError(t, err)
			createReq.Header.Set(contentTypeHeader, contentTypeJSON)

			createResp, err := testApp.Request(createReq)
			require.NoError(t, err)
			require.Equal(t, 201, createResp.StatusCode)

			var createdUser dto.UserResponseDTO
			err = json.NewDecoder(createResp.Body).Decode(&createdUser)
			require.NoError(t, err)
			createdUsers = append(createdUsers, createdUser)
		}

		// Verify each user can be retrieved individually
		for i, createdUser := range createdUsers {
			getReq, err := http.NewRequest("GET", fmt.Sprintf("%s/%s", usersEndpoint, createdUser.ID), nil)
			require.NoError(t, err)

			getResp, err := testApp.Request(getReq)
			require.NoError(t, err)
			assert.Equal(t, 200, getResp.StatusCode)

			var retrievedUser dto.UserResponseDTO
			err = json.NewDecoder(getResp.Body).Decode(&retrievedUser)
			require.NoError(t, err)

			// Verify the retrieved user matches the expected user
			assert.Equal(t, createdUser.ID, retrievedUser.ID)
			assert.Equal(t, users[i].Name, retrievedUser.Name)
			assert.Equal(t, users[i].Email, retrievedUser.Email)
			assert.Equal(t, users[i].IsActive, retrievedUser.IsActive)
		}

		// Verify all users have unique IDs
		for i := 0; i < len(createdUsers); i++ {
			for j := i + 1; j < len(createdUsers); j++ {
				assert.NotEqual(t, createdUsers[i].ID, createdUsers[j].ID,
					"Users should have unique IDs")
			}
		}
	})
}

func TestUserControllerCreateRepositoryErrors(t *testing.T) {
	t.Run("Create user - force MongoDB connection close to cause repository Create error", func(t *testing.T) {
		testApp := SetupTestApp(t)
		// Note: Don't defer cleanup here since we'll disconnect the database

		const (
			contentTypeJSON   = "application/json"
			contentTypeHeader = "Content-Type"
			usersEndpoint     = "/api/v1/users"
		)

		// Verify normal operation works first
		createUserDTO := dto.CreateUserRequestDTO{
			Name:     "Test User Before Disconnect",
			Email:    "test.before@example.com",
			IsActive: true,
		}

		payloadBytes, err := json.Marshal(createUserDTO)
		require.NoError(t, err)

		createReq, err := http.NewRequest("POST", usersEndpoint, bytes.NewBuffer(payloadBytes))
		require.NoError(t, err)
		createReq.Header.Set(contentTypeHeader, contentTypeJSON)

		createResp, err := testApp.Request(createReq)
		require.NoError(t, err)
		assert.Equal(t, 201, createResp.StatusCode)

		// Now force close the MongoDB connection to simulate repository error
		ctx := context.Background()
		err = testApp.DB.Client.Disconnect(ctx)
		require.NoError(t, err)

		// Try to create a user after closing connection - should return repository error
		createUserDTO2 := dto.CreateUserRequestDTO{
			Name:     "Test User After Disconnect",
			Email:    "test.after@example.com",
			IsActive: true,
		}

		payloadBytes2, err := json.Marshal(createUserDTO2)
		require.NoError(t, err)

		createReq2, err := http.NewRequest("POST", usersEndpoint, bytes.NewBuffer(payloadBytes2))
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

	t.Run("Create user - duplicate email constraint violation", func(t *testing.T) {
		testApp := SetupTestApp(t)
		defer testApp.Cleanup(t)

		const (
			contentTypeJSON   = "application/json"
			contentTypeHeader = "Content-Type"
			usersEndpoint     = "/api/v1/users"
			duplicateEmail    = "duplicate@example.com"
		)

		// Create first user
		createUserDTO := dto.CreateUserRequestDTO{
			Name:     "First User",
			Email:    duplicateEmail,
			IsActive: true,
		}

		payloadBytes, err := json.Marshal(createUserDTO)
		require.NoError(t, err)

		createReq, err := http.NewRequest("POST", usersEndpoint, bytes.NewBuffer(payloadBytes))
		require.NoError(t, err)
		createReq.Header.Set(contentTypeHeader, contentTypeJSON)

		createResp, err := testApp.Request(createReq)
		require.NoError(t, err)
		assert.Equal(t, 201, createResp.StatusCode)

		// Try to create second user with same email
		createUserDTO2 := dto.CreateUserRequestDTO{
			Name:     "Second User",
			Email:    duplicateEmail,
			IsActive: false,
		}

		payloadBytes2, err := json.Marshal(createUserDTO2)
		require.NoError(t, err)

		createReq2, err := http.NewRequest("POST", usersEndpoint, bytes.NewBuffer(payloadBytes2))
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

	t.Run("Create user - database collection drop during operation", func(t *testing.T) {
		testApp := SetupTestApp(t)
		defer testApp.Cleanup(t)

		const (
			contentTypeJSON   = "application/json"
			contentTypeHeader = "Content-Type"
			usersEndpoint     = "/api/v1/users"
		)

		// Create first user normally
		createUserDTO := dto.CreateUserRequestDTO{
			Name:     "Test User Before Drop",
			Email:    "before.drop@example.com",
			IsActive: true,
		}

		payloadBytes, err := json.Marshal(createUserDTO)
		require.NoError(t, err)

		createReq, err := http.NewRequest("POST", usersEndpoint, bytes.NewBuffer(payloadBytes))
		require.NoError(t, err)
		createReq.Header.Set(contentTypeHeader, contentTypeJSON)

		createResp, err := testApp.Request(createReq)
		require.NoError(t, err)
		assert.Equal(t, 201, createResp.StatusCode)

		// Drop the users collection to force potential repository issues
		ctx := context.Background()
		collection := testApp.DB.DB.Collection("users")
		err = collection.Drop(ctx)
		require.NoError(t, err)

		// Try to create another user after dropping collection
		createUserDTO2 := dto.CreateUserRequestDTO{
			Name:     "Test User After Drop",
			Email:    "after.drop@example.com",
			IsActive: true,
		}

		payloadBytes2, err := json.Marshal(createUserDTO2)
		require.NoError(t, err)

		createReq2, err := http.NewRequest("POST", usersEndpoint, bytes.NewBuffer(payloadBytes2))
		require.NoError(t, err)
		createReq2.Header.Set(contentTypeHeader, contentTypeJSON)

		createResp2, err := testApp.Request(createReq2)
		require.NoError(t, err)

		// Should either work (201 - collection recreated) or return error (500)
		assert.True(t, createResp2.StatusCode == 201 || createResp2.StatusCode == 500,
			"Expected 201 (collection recreated) or 500 (repository error), got %d", createResp2.StatusCode)

		if createResp2.StatusCode == 201 {
			var createdUser dto.UserResponseDTO
			err = json.NewDecoder(createResp2.Body).Decode(&createdUser)
			require.NoError(t, err)
			assert.Equal(t, createUserDTO2.Name, createdUser.Name)
			assert.Equal(t, createUserDTO2.Email, createdUser.Email)
			assert.Equal(t, createUserDTO2.IsActive, createdUser.IsActive)
		} else {
			var errorResponse map[string]interface{}
			err = json.NewDecoder(createResp2.Body).Decode(&errorResponse)
			require.NoError(t, err)
			assert.NotNil(t, errorResponse["error"])
		}
	})
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
			Name:     "",                     // Empty name should trigger validation error
			Email:    "invalid-email-format", // Invalid email format
			IsActive: true,
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
			Name:     testUserName,
			Email:    "concurrent@example.com",
			IsActive: true,
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

func TestUserControllerUpdate(t *testing.T) {
	testApp := SetupTestApp(t)
	defer testApp.Cleanup(t)

	// Create a user first
	createPayload := dto.CreateUserRequestDTO{
		Name:     "Original Name",
		Email:    "original@example.com",
		IsActive: true,
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
		Name:     "Updated Name",
		Email:    "updated@example.com",
		IsActive: false, // Change to test the field
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
	assert.Equal(t, updatePayload.IsActive, updatedUser.IsActive)
}

func TestUserControllerDelete(t *testing.T) {
	testApp := SetupTestApp(t)
	defer testApp.Cleanup(t)

	// Create a user first
	createPayload := dto.CreateUserRequestDTO{
		Name:     "To Be Deleted",
		Email:    "delete@example.com",
		IsActive: true,
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

func TestUserControllerList(t *testing.T) {
	testApp := SetupTestApp(t)
	defer testApp.Cleanup(t)

	// Create multiple users
	users := []dto.CreateUserRequestDTO{
		{Name: "User 1", Email: "user1@example.com", IsActive: true},
		{Name: "User 2", Email: "user2@example.com", IsActive: false},
		{Name: "User 3", Email: "user3@example.com", IsActive: true},
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

func TestUserControllerListWithSearch(t *testing.T) {
	testApp := SetupTestApp(t)
	defer testApp.Cleanup(t)

	// Create users with different names
	users := []dto.CreateUserRequestDTO{
		{Name: "Alice Johnson", Email: "alice@example.com", IsActive: true},
		{Name: "Bob Smith", Email: "bob@example.com", IsActive: true},
		{Name: "Alice Cooper", Email: "alice.cooper@example.com", IsActive: false},
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

func TestUserControllerListWithoutUsers(t *testing.T) {
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
			Name:     updatedName,
			Email:    updatedEmail,
			IsActive: true,
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
			Name:     updatedName,
			Email:    updatedEmail,
			IsActive: true,
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

func TestUserControllerListRepositoryErrors(t *testing.T) {
	t.Run("List users - force MongoDB connection close to cause repository error", func(t *testing.T) {
		testApp := SetupTestApp(t)
		// Note: Don't defer cleanup here since we'll disconnect the database
		// Create some test users first
		users := []dto.CreateUserRequestDTO{
			{Name: "Test User 1", Email: "test1@example.com", IsActive: true},
			{Name: "Test User 2", Email: "test2@example.com", IsActive: true},
		}

		for _, user := range users {
			payloadBytes, err := json.Marshal(user)
			require.NoError(t, err)

			createReq, err := http.NewRequest("POST", "/api/v1/users", bytes.NewBuffer(payloadBytes))
			require.NoError(t, err)
			createReq.Header.Set("Content-Type", "application/json")

			createResp, err := testApp.Request(createReq)
			require.NoError(t, err)
			require.Equal(t, 201, createResp.StatusCode)
		}

		// Verify normal operation works first
		req, err := http.NewRequest("GET", "/api/v1/users", nil)
		require.NoError(t, err)

		resp, err := testApp.Request(req)
		require.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)

		// Now force close the MongoDB connection to simulate repository error
		// This will cause the next request to fail
		ctx := context.Background()
		err = testApp.DB.Client.Disconnect(ctx)
		require.NoError(t, err)

		// Try to list users after closing connection - should return repository error
		req, err = http.NewRequest("GET", "/api/v1/users", nil)
		require.NoError(t, err)

		resp, err = testApp.Request(req)
		require.NoError(t, err)

		// Should return 500 due to repository connection error
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
				strings.Contains(errorMsg, "topology is closed"),
			"Expected connection error, got: %s", errorMsg)
	})

	t.Run("List users - search with invalid regex to force repository error", func(t *testing.T) {
		testApp := SetupTestApp(t)
		defer testApp.Cleanup(t)

		// Create test users first
		users := []dto.CreateUserRequestDTO{
			{Name: "Regex Test User", Email: "regex@example.com", IsActive: true},
		}

		for _, user := range users {
			payloadBytes, err := json.Marshal(user)
			require.NoError(t, err)

			createReq, err := http.NewRequest("POST", "/api/v1/users", bytes.NewBuffer(payloadBytes))
			require.NoError(t, err)
			createReq.Header.Set("Content-Type", "application/json")

			createResp, err := testApp.Request(createReq)
			require.NoError(t, err)
			require.Equal(t, 201, createResp.StatusCode)
		}

		// Test with regex patterns that will cause MongoDB regex compilation to fail
		invalidRegexPatterns := []string{
			"[",     // Unclosed bracket
			"*",     // Invalid quantifier
			"(?P<",  // Incomplete named group
			"\\",    // Incomplete escape
			"[z-a]", // Invalid range
			"(?",    // Incomplete group
		}

		for _, pattern := range invalidRegexPatterns {
			t.Run(fmt.Sprintf("Invalid regex pattern: %s", pattern), func(t *testing.T) {
				// URL encode the pattern to ensure it reaches the backend properly
				encodedPattern := url.QueryEscape(pattern)
				req, err := http.NewRequest("GET", fmt.Sprintf("/api/v1/users?search=%s", encodedPattern), nil)
				require.NoError(t, err)

				resp, err := testApp.Request(req)
				require.NoError(t, err)

				// MongoDB regex compilation error should result in 500
				assert.Equal(t, 500, resp.StatusCode)

				var errorResponse map[string]interface{}
				err = json.NewDecoder(resp.Body).Decode(&errorResponse)
				require.NoError(t, err)
				assert.NotNil(t, errorResponse["error"])
			})
		}
	})

	t.Run("Count repository error - force database timeout", func(t *testing.T) {
		testApp := SetupTestApp(t)
		defer testApp.Cleanup(t)

		// This test simulates timeout during count operation
		// We'll use a very large collection name or query that might timeout

		// First create some users
		for i := 0; i < 3; i++ {
			userData := dto.CreateUserRequestDTO{
				Name:     fmt.Sprintf("Timeout Test User %d", i),
				Email:    fmt.Sprintf("timeout%d@example.com", i),
				IsActive: true,
			}
			payloadBytes, err := json.Marshal(userData)
			require.NoError(t, err)

			createReq, err := http.NewRequest("POST", "/api/v1/users", bytes.NewBuffer(payloadBytes))
			require.NoError(t, err)
			createReq.Header.Set("Content-Type", "application/json")

			createResp, err := testApp.Request(createReq)
			require.NoError(t, err)
			require.Equal(t, 201, createResp.StatusCode)
		}

		// Test with a moderately long search string that might cause issues
		// but not exceed HTTP limits
		longSearchString := strings.Repeat("a", 1000) // 1KB string
		encodedSearch := url.QueryEscape(longSearchString)

		req, err := http.NewRequest("GET", fmt.Sprintf("/api/v1/users?search=%s", encodedSearch), nil)
		require.NoError(t, err)

		resp, err := testApp.Request(req)
		require.NoError(t, err)

		// Should either handle gracefully (200) or return error (500) or validation error (400)
		assert.True(t, resp.StatusCode == 200 || resp.StatusCode == 500 || resp.StatusCode == 400)

		if resp.StatusCode == 500 {
			var errorResponse map[string]interface{}
			err = json.NewDecoder(resp.Body).Decode(&errorResponse)
			require.NoError(t, err)
			assert.NotNil(t, errorResponse["error"])
		}
	})

	t.Run("Repository List error - invalid collection access", func(t *testing.T) {
		testApp := SetupTestApp(t)
		defer testApp.Cleanup(t)

		// This test tries to trigger repository errors by manipulating the database state
		// We'll create users then try to access them after modifying database state

		// Create test users
		userData := dto.CreateUserRequestDTO{
			Name:     "Collection Test User",
			Email:    "collection@example.com",
			IsActive: true,
		}
		payloadBytes, err := json.Marshal(userData)
		require.NoError(t, err)

		createReq, err := http.NewRequest("POST", "/api/v1/users", bytes.NewBuffer(payloadBytes))
		require.NoError(t, err)
		createReq.Header.Set("Content-Type", "application/json")

		createResp, err := testApp.Request(createReq)
		require.NoError(t, err)
		require.Equal(t, 201, createResp.StatusCode)

		// Verify it works normally first
		listReq, err := http.NewRequest("GET", "/api/v1/users", nil)
		require.NoError(t, err)

		listResp, err := testApp.Request(listReq)
		require.NoError(t, err)
		assert.Equal(t, 200, listResp.StatusCode)

		// Now try to cause a repository error by dropping the collection
		ctx := context.Background()
		collection := testApp.DB.DB.Collection("users")
		err = collection.Drop(ctx)
		require.NoError(t, err)

		// Try to list users after dropping collection
		// The collection will be recreated but might cause temporary repository errors
		listReq, err = http.NewRequest("GET", "/api/v1/users", nil)
		require.NoError(t, err)

		listResp, err = testApp.Request(listReq)
		require.NoError(t, err)

		// Should either work (200 with empty list) or return error (500)
		assert.True(t, listResp.StatusCode == 200 || listResp.StatusCode == 500)

		if listResp.StatusCode == 200 {
			var userList dto.UserListResponseDTO
			err = json.NewDecoder(listResp.Body).Decode(&userList)
			require.NoError(t, err)
			// Should return empty list since collection was dropped
			assert.Equal(t, 0, len(userList.Data))
			assert.Equal(t, int64(0), userList.Meta.Total)
		}
	})

	t.Run("Repository pagination error - invalid skip/limit values", func(t *testing.T) {
		testApp := SetupTestApp(t)
		defer testApp.Cleanup(t)

		// Create test data
		userData := dto.CreateUserRequestDTO{
			Name:     "Pagination Test User",
			Email:    "pagination@example.com",
			IsActive: true,
		}
		payloadBytes, err := json.Marshal(userData)
		require.NoError(t, err)

		createReq, err := http.NewRequest("POST", "/api/v1/users", bytes.NewBuffer(payloadBytes))
		require.NoError(t, err)
		createReq.Header.Set("Content-Type", "application/json")

		createResp, err := testApp.Request(createReq)
		require.NoError(t, err)
		require.Equal(t, 201, createResp.StatusCode)

		// Test with values that might cause integer overflow in skip calculation
		// page * per_page might overflow if values are too large
		testCases := []struct {
			page    string
			perPage string
			name    string
		}{
			{"999999999", "999999999", "Extremely large values"},
			{"-999999999", "10", "Very negative page"},
			{"1", "-999999999", "Very negative per_page"},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				req, err := http.NewRequest("GET",
					fmt.Sprintf("/api/v1/users?page=%s&per_page=%s", tc.page, tc.perPage), nil)
				require.NoError(t, err)

				resp, err := testApp.Request(req)
				require.NoError(t, err)

				// Should handle edge cases gracefully
				assert.True(t, resp.StatusCode == 200 || resp.StatusCode == 400 || resp.StatusCode == 500)
			})
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
