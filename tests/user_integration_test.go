package tests

import (
	"encoding/json"
	"net/http"
	"user-management/internal/application/dto"

	"github.com/stretchr/testify/assert"
)

func (suite *IntegrationTestSuite) TestUserCRUD() {
	// Test Create User
	createUserDTO := dto.CreateUserRequestDTO{
		Name:  "John Doe",
		Email: "john@example.com",
	}

	resp, body := suite.makeRequest("POST", "/api/v1/users/", createUserDTO)
	assert.Equal(suite.T(), http.StatusCreated, resp.StatusCode)

	var createdUser dto.UserResponseDTO
	err := json.Unmarshal(body, &createdUser)
	suite.NoError(err)
	assert.NotEmpty(suite.T(), createdUser.ID)
	assert.Equal(suite.T(), createUserDTO.Name, createdUser.Name)
	assert.Equal(suite.T(), createUserDTO.Email, createdUser.Email)

	userID := createdUser.ID

	// Test Get User
	resp, body = suite.makeRequest("GET", "/api/v1/users/"+userID, nil)
	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	var retrievedUser dto.UserResponseDTO
	err = json.Unmarshal(body, &retrievedUser)
	suite.NoError(err)
	assert.Equal(suite.T(), createdUser.ID, retrievedUser.ID)
	assert.Equal(suite.T(), createUserDTO.Name, retrievedUser.Name)
	assert.Equal(suite.T(), createUserDTO.Email, retrievedUser.Email)

	// Test Update User
	updateUserDTO := dto.CreateUserRequestDTO{
		Name:  "John Updated",
		Email: "john.updated@example.com",
	}

	resp, body = suite.makeRequest("PUT", "/api/v1/users/"+userID, updateUserDTO)
	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	var updatedUser dto.UserResponseDTO
	err = json.Unmarshal(body, &updatedUser)
	suite.NoError(err)
	assert.Equal(suite.T(), updateUserDTO.Name, updatedUser.Name)
	assert.Equal(suite.T(), updateUserDTO.Email, updatedUser.Email)

	// Test List Users
	resp, body = suite.makeRequest("GET", "/api/v1/users/", nil)
	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	var users []dto.UserResponseDTO
	err = json.Unmarshal(body, &users)
	suite.NoError(err)
	assert.Len(suite.T(), users, 1)
	assert.Equal(suite.T(), updateUserDTO.Name, users[0].Name)

	// Test Delete User
	resp, _ = suite.makeRequest("DELETE", "/api/v1/users/"+userID, nil)
	assert.Equal(suite.T(), http.StatusNoContent, resp.StatusCode)

	// Verify user is deleted
	resp, _ = suite.makeRequest("GET", "/api/v1/users/"+userID, nil)
	assert.Equal(suite.T(), http.StatusNotFound, resp.StatusCode)
}

func (suite *IntegrationTestSuite) TestUserNotFound() {
	resp, _ := suite.makeRequest("GET", "/api/v1/users/nonexistent", nil)
	assert.Equal(suite.T(), http.StatusNotFound, resp.StatusCode)
}

func (suite *IntegrationTestSuite) TestCreateUserInvalidData() {
	// Test with invalid JSON
	invalidJSON := "{"
	resp, _ := suite.makeRequest("POST", "/api/v1/users/", invalidJSON)
	assert.Equal(suite.T(), http.StatusBadRequest, resp.StatusCode)

	// Test with missing required fields
	invalidUser := map[string]interface{}{
		"invalid_field": "data",
	}
	resp2, _ := suite.makeRequest("POST", "/api/v1/users/", invalidUser)
	assert.Equal(suite.T(), http.StatusBadRequest, resp2.StatusCode)
}

func (suite *IntegrationTestSuite) TestListUsersEmpty() {
	resp, body := suite.makeRequest("GET", "/api/v1/users/", nil)
	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	var users []dto.UserResponseDTO
	err := json.Unmarshal(body, &users)
	suite.NoError(err)
	assert.Empty(suite.T(), users)
}

func (suite *IntegrationTestSuite) TestMultipleUsers() {
	// Create multiple users
	createRequests := []dto.CreateUserRequestDTO{
		{Name: "User 1", Email: "user1@example.com"},
		{Name: "User 2", Email: "user2@example.com"},
		{Name: "User 3", Email: "user3@example.com"},
	}

	createdUsers := make([]dto.UserResponseDTO, 0, len(createRequests))
	for _, req := range createRequests {
		resp, body := suite.makeRequest("POST", "/api/v1/users/", req)
		assert.Equal(suite.T(), http.StatusCreated, resp.StatusCode)

		var createdUser dto.UserResponseDTO
		err := json.Unmarshal(body, &createdUser)
		suite.NoError(err)
		createdUsers = append(createdUsers, createdUser)
	}

	// List all users
	resp, body := suite.makeRequest("GET", "/api/v1/users/", nil)
	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	var retrievedUsers []dto.UserResponseDTO
	err := json.Unmarshal(body, &retrievedUsers)
	suite.NoError(err)
	assert.Len(suite.T(), retrievedUsers, 3)

	// Verify each user exists
	for _, user := range createdUsers {
		resp, body := suite.makeRequest("GET", "/api/v1/users/"+user.ID, nil)
		assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

		var retrievedUser dto.UserResponseDTO
		err := json.Unmarshal(body, &retrievedUser)
		suite.NoError(err)
		assert.Equal(suite.T(), user.ID, retrievedUser.ID)
		assert.Equal(suite.T(), user.Name, retrievedUser.Name)
		assert.Equal(suite.T(), user.Email, retrievedUser.Email)
	}
}
