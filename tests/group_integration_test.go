package tests

import (
	"encoding/json"
	"net/http"
	"user-management/internal/application/dto"

	"github.com/stretchr/testify/assert"
)

func (suite *IntegrationTestSuite) TestGroupCRUD() {
	// Test Create Group
	createGroupDTO := dto.CreateGroupRequestDTO{
		Name:    "Developers",
		Members: []string{},
	}

	resp, body := suite.makeRequest("POST", "/api/v1/groups/", createGroupDTO)
	assert.Equal(suite.T(), http.StatusCreated, resp.StatusCode)

	var createdGroup dto.GroupResponseDTO
	err := json.Unmarshal(body, &createdGroup)
	suite.NoError(err)
	assert.NotEmpty(suite.T(), createdGroup.ID)
	assert.Equal(suite.T(), createGroupDTO.Name, createdGroup.Name)
	assert.Equal(suite.T(), createGroupDTO.Members, createdGroup.Members)

	groupID := createdGroup.ID

	// Test Get Group
	resp, body = suite.makeRequest("GET", "/api/v1/groups/"+groupID, nil)
	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	var retrievedGroup dto.GroupResponseDTO
	err = json.Unmarshal(body, &retrievedGroup)
	suite.NoError(err)
	assert.Equal(suite.T(), createdGroup.ID, retrievedGroup.ID)
	assert.Equal(suite.T(), createGroupDTO.Name, retrievedGroup.Name)

	// Test Update Group
	updateGroupDTO := dto.CreateGroupRequestDTO{
		Name:    "Senior Developers",
		Members: []string{},
	}

	resp, body = suite.makeRequest("PUT", "/api/v1/groups/"+groupID, updateGroupDTO)
	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	var updatedGroup dto.GroupResponseDTO
	err = json.Unmarshal(body, &updatedGroup)
	suite.NoError(err)
	assert.Equal(suite.T(), updateGroupDTO.Name, updatedGroup.Name)

	// Test List Groups
	resp, body = suite.makeRequest("GET", "/api/v1/groups/", nil)
	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	var groups []dto.GroupResponseDTO
	err = json.Unmarshal(body, &groups)
	suite.NoError(err)
	assert.Len(suite.T(), groups, 1)
	assert.Equal(suite.T(), updateGroupDTO.Name, groups[0].Name)

	// Test Delete Group
	resp, _ = suite.makeRequest("DELETE", "/api/v1/groups/"+groupID, nil)
	assert.Equal(suite.T(), http.StatusNoContent, resp.StatusCode)

	// Verify group is deleted
	resp, _ = suite.makeRequest("GET", "/api/v1/groups/"+groupID, nil)
	assert.Equal(suite.T(), http.StatusNotFound, resp.StatusCode)
}

func (suite *IntegrationTestSuite) TestGroupNotFound() {
	resp, _ := suite.makeRequest("GET", "/api/v1/groups/nonexistent", nil)
	assert.Equal(suite.T(), http.StatusNotFound, resp.StatusCode)
}

func (suite *IntegrationTestSuite) TestGroupMemberManagement() {
	// Create a user first
	user := dto.CreateUserRequestDTO{
		Name:  "John Doe",
		Email: "john@example.com",
	}
	resp, body := suite.makeRequest("POST", "/api/v1/users/", user)
	assert.Equal(suite.T(), http.StatusCreated, resp.StatusCode)

	var createdUser dto.UserResponseDTO
	err := json.Unmarshal(body, &createdUser)
	suite.NoError(err)

	// Create a group
	group := dto.CreateGroupRequestDTO{
		Name: "Developers",
	}
	resp, body = suite.makeRequest("POST", "/api/v1/groups/", group)
	assert.Equal(suite.T(), http.StatusCreated, resp.StatusCode)

	var createdGroup dto.GroupResponseDTO
	err = json.Unmarshal(body, &createdGroup)
	suite.NoError(err)

	// Add user to group
	resp, _ = suite.makeRequest("POST", "/api/v1/groups/"+createdGroup.ID+"/members/"+createdUser.ID, nil)
	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	// Verify user was added to group
	resp, body = suite.makeRequest("GET", "/api/v1/groups/"+createdGroup.ID, nil)
	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	var retrievedGroup dto.GroupResponseDTO
	err = json.Unmarshal(body, &retrievedGroup)
	suite.NoError(err)
	assert.Contains(suite.T(), retrievedGroup.Members, createdUser.ID)

	// Remove user from group
	resp, _ = suite.makeRequest("DELETE", "/api/v1/groups/"+createdGroup.ID+"/members/"+createdUser.ID, nil)
	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	// Verify user was removed from group
	resp, body = suite.makeRequest("GET", "/api/v1/groups/"+createdGroup.ID, nil)
	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	err = json.Unmarshal(body, &retrievedGroup)
	suite.NoError(err)
	assert.NotContains(suite.T(), retrievedGroup.Members, createdUser.ID)
}

func (suite *IntegrationTestSuite) TestAddNonExistentUserToGroup() {
	// Create a group
	group := dto.CreateGroupRequestDTO{
		Name: "Developers",
	}
	resp, body := suite.makeRequest("POST", "/api/v1/groups/", group)
	assert.Equal(suite.T(), http.StatusCreated, resp.StatusCode)

	var createdGroup dto.GroupResponseDTO
	err := json.Unmarshal(body, &createdGroup)
	suite.NoError(err)

	// Try to add non-existent user to group
	resp, _ = suite.makeRequest("POST", "/api/v1/groups/"+createdGroup.ID+"/members/nonexistent", nil)
	assert.Equal(suite.T(), http.StatusInternalServerError, resp.StatusCode)
}

func (suite *IntegrationTestSuite) TestAddUserToNonExistentGroup() {
	// Create a user
	user := dto.CreateUserRequestDTO{
		Name:  "John Doe",
		Email: "john@example.com",
	}
	resp, body := suite.makeRequest("POST", "/api/v1/users/", user)
	assert.Equal(suite.T(), http.StatusCreated, resp.StatusCode)

	var createdUser dto.UserResponseDTO
	err := json.Unmarshal(body, &createdUser)
	suite.NoError(err)

	// Try to add user to non-existent group
	resp, _ = suite.makeRequest("POST", "/api/v1/groups/nonexistent/members/"+createdUser.ID, nil)
	assert.Equal(suite.T(), http.StatusInternalServerError, resp.StatusCode)
}

func (suite *IntegrationTestSuite) TestListGroupsEmpty() {
	resp, body := suite.makeRequest("GET", "/api/v1/groups/", nil)
	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	var groups []dto.GroupResponseDTO
	err := json.Unmarshal(body, &groups)
	suite.NoError(err)
	assert.Empty(suite.T(), groups)
}

func (suite *IntegrationTestSuite) TestMultipleGroups() {
	// Create multiple groups
	createGroups := []dto.CreateGroupRequestDTO{
		{Name: "Developers"},
		{Name: "Designers"},
		{Name: "Managers"},
	}

	createdGroups := make([]dto.GroupResponseDTO, 0, len(createGroups))
	for _, group := range createGroups {
		resp, body := suite.makeRequest("POST", "/api/v1/groups/", group)
		assert.Equal(suite.T(), http.StatusCreated, resp.StatusCode)

		var createdGroup dto.GroupResponseDTO
		err := json.Unmarshal(body, &createdGroup)
		suite.NoError(err)
		createdGroups = append(createdGroups, createdGroup)
	}

	// List all groups
	resp, body := suite.makeRequest("GET", "/api/v1/groups/", nil)
	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	var retrievedGroups []dto.GroupResponseDTO
	err := json.Unmarshal(body, &retrievedGroups)
	suite.NoError(err)
	assert.Len(suite.T(), retrievedGroups, 3)

	// Verify each group exists
	for _, group := range createdGroups {
		resp, body := suite.makeRequest("GET", "/api/v1/groups/"+group.ID, nil)
		assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

		var retrievedGroup dto.GroupResponseDTO
		err := json.Unmarshal(body, &retrievedGroup)
		suite.NoError(err)
		assert.Equal(suite.T(), group.ID, retrievedGroup.ID)
		assert.Equal(suite.T(), group.Name, retrievedGroup.Name)
	}
}
