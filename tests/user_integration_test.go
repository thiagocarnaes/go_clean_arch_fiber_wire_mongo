package tests

import (
	"encoding/json"
	"net/http"
	"user-management/internal/application/dto"

	"github.com/stretchr/testify/assert"
)

func (suite *IntegrationTestSuite) TestUserCRUD() {
	// Test Create User
	createUserDTO := dto.UserDTO{
		ID:    "user1",
		Name:  "John Doe",
		Email: "john@example.com",
	}

	resp, body := suite.makeRequest("POST", "/api/v1/users/", createUserDTO)
	assert.Equal(suite.T(), http.StatusCreated, resp.StatusCode)

	var createdUser dto.UserDTO
	err := json.Unmarshal(body, &createdUser)
	suite.NoError(err)
	assert.Equal(suite.T(), createUserDTO.ID, createdUser.ID)
	assert.Equal(suite.T(), createUserDTO.Name, createdUser.Name)
	assert.Equal(suite.T(), createUserDTO.Email, createdUser.Email)

	// Test Get User
	resp, body = suite.makeRequest("GET", "/api/v1/users/user1", nil)
	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	var retrievedUser dto.UserDTO
	err = json.Unmarshal(body, &retrievedUser)
	suite.NoError(err)
	assert.Equal(suite.T(), createUserDTO.ID, retrievedUser.ID)
	assert.Equal(suite.T(), createUserDTO.Name, retrievedUser.Name)
	assert.Equal(suite.T(), createUserDTO.Email, retrievedUser.Email)

	// Test Update User
	updateUserDTO := dto.UserDTO{
		ID:    "user1",
		Name:  "John Updated",
		Email: "john.updated@example.com",
	}

	resp, body = suite.makeRequest("PUT", "/api/v1/users/user1", updateUserDTO)
	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	var updatedUser dto.UserDTO
	err = json.Unmarshal(body, &updatedUser)
	suite.NoError(err)
	assert.Equal(suite.T(), updateUserDTO.Name, updatedUser.Name)
	assert.Equal(suite.T(), updateUserDTO.Email, updatedUser.Email)

	// Test List Users
	resp, body = suite.makeRequest("GET", "/api/v1/users/", nil)
	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	var users []dto.UserDTO
	err = json.Unmarshal(body, &users)
	suite.NoError(err)
	assert.Len(suite.T(), users, 1)
	assert.Equal(suite.T(), updateUserDTO.Name, users[0].Name)

	// Test Delete User
	resp, _ = suite.makeRequest("DELETE", "/api/v1/users/user1", nil)
	assert.Equal(suite.T(), http.StatusNoContent, resp.StatusCode)

	// Verify user is deleted
	resp, _ = suite.makeRequest("GET", "/api/v1/users/user1", nil)
	assert.Equal(suite.T(), http.StatusNotFound, resp.StatusCode)
}

func (suite *IntegrationTestSuite) TestUserNotFound() {
	resp, _ := suite.makeRequest("GET", "/api/v1/users/nonexistent", nil)
	assert.Equal(suite.T(), http.StatusNotFound, resp.StatusCode)
}

func (suite *IntegrationTestSuite) TestCreateUserInvalidData() {
	invalidUser := map[string]interface{}{
		"invalid": "data",
	}

	resp, _ := suite.makeRequest("POST", "/api/v1/users/", invalidUser)
	assert.Equal(suite.T(), http.StatusCreated, resp.StatusCode) // Still creates with empty fields
}

func (suite *IntegrationTestSuite) TestListUsersEmpty() {
	resp, body := suite.makeRequest("GET", "/api/v1/users/", nil)
	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	var users []dto.UserDTO
	err := json.Unmarshal(body, &users)
	suite.NoError(err)
	assert.Empty(suite.T(), users)
}

func (suite *IntegrationTestSuite) TestMultipleUsers() {
	// Create multiple users
	users := []dto.UserDTO{
		{ID: "user1", Name: "User 1", Email: "user1@example.com"},
		{ID: "user2", Name: "User 2", Email: "user2@example.com"},
		{ID: "user3", Name: "User 3", Email: "user3@example.com"},
	}

	for _, user := range users {
		resp, _ := suite.makeRequest("POST", "/api/v1/users/", user)
		assert.Equal(suite.T(), http.StatusCreated, resp.StatusCode)
	}

	// List all users
	resp, body := suite.makeRequest("GET", "/api/v1/users/", nil)
	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	var retrievedUsers []dto.UserDTO
	err := json.Unmarshal(body, &retrievedUsers)
	suite.NoError(err)
	assert.Len(suite.T(), retrievedUsers, 3)

	// Verify each user exists
	for _, user := range users {
		resp, body := suite.makeRequest("GET", "/api/v1/users/"+user.ID, nil)
		assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

		var retrievedUser dto.UserDTO
		err := json.Unmarshal(body, &retrievedUser)
		suite.NoError(err)
		assert.Equal(suite.T(), user.ID, retrievedUser.ID)
		assert.Equal(suite.T(), user.Name, retrievedUser.Name)
		assert.Equal(suite.T(), user.Email, retrievedUser.Email)
	}
}
