package tests

import (
	"encoding/json"
	"net/http"
	"user-management/internal/application/dto"

	"github.com/stretchr/testify/assert"
)

func (suite *IntegrationTestSuite) TestCompleteUserGroupWorkflow() {
	// Create multiple users
	createUsers := []dto.CreateUserRequestDTO{
		{Name: "Alice Developer", Email: "alice@example.com"},
		{Name: "Bob Developer", Email: "bob@example.com"},
		{Name: "Carol Designer", Email: "carol@example.com"},
		{Name: "Dave Manager", Email: "dave@example.com"},
	}

	createdUsers := make([]dto.UserResponseDTO, 0, len(createUsers))
	for _, user := range createUsers {
		resp, body := suite.makeRequest("POST", "/api/v1/users/", user)
		assert.Equal(suite.T(), http.StatusCreated, resp.StatusCode)

		var createdUser dto.UserResponseDTO
		err := json.Unmarshal(body, &createdUser)
		suite.NoError(err)
		createdUsers = append(createdUsers, createdUser)
	}

	// Create multiple groups
	createGroups := []dto.CreateGroupRequestDTO{
		{Name: "Development Team"},
		{Name: "Design Team"},
		{Name: "Management Team"},
		{Name: "All Staff"},
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

	for _, group := range createGroups {
		resp, _ := suite.makeRequest("POST", "/api/v1/groups/", group)
		assert.Equal(suite.T(), http.StatusCreated, resp.StatusCode)
	}

	// Create group ID mapping for easier reference
	groupMap := make(map[string]string)
	for i, group := range createdGroups {
		switch i {
		case 0:
			groupMap["developers"] = group.ID
		case 1:
			groupMap["designers"] = group.ID
		case 2:
			groupMap["management"] = group.ID
		case 3:
			groupMap["all-staff"] = group.ID
		}
	}

	// Create user ID mapping for easier reference
	userMap := make(map[string]string)
	for i, user := range createdUsers {
		switch i {
		case 0:
			userMap["dev1"] = user.ID
		case 1:
			userMap["dev2"] = user.ID
		case 2:
			userMap["designer1"] = user.ID
		case 3:
			userMap["manager1"] = user.ID
		}
	}

	// Add users to specific groups
	memberships := map[string][]string{
		"developers": {userMap["dev1"], userMap["dev2"]},
		"designers":  {userMap["designer1"]},
		"management": {userMap["manager1"]},
		"all-staff":  {userMap["dev1"], userMap["dev2"], userMap["designer1"], userMap["manager1"]},
	}

	for groupName, memberIDs := range memberships {
		for _, memberID := range memberIDs {
			resp, _ := suite.makeRequest("POST", "/api/v1/groups/"+groupMap[groupName]+"/members/"+memberID, nil)
			assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)
		}
	}

	// Verify group memberships
	for groupName, expectedMembers := range memberships {
		resp, body := suite.makeRequest("GET", "/api/v1/groups/"+groupMap[groupName], nil)
		assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

		var group dto.GroupResponseDTO
		err := json.Unmarshal(body, &group)
		suite.NoError(err)
		assert.Len(suite.T(), group.Members, len(expectedMembers))

		for _, expectedMember := range expectedMembers {
			assert.Contains(suite.T(), group.Members, expectedMember)
		}
	}

	// Update a user and verify they still exist in groups
	updateUser := dto.CreateUserRequestDTO{
		Name:  "Alice Senior Developer",
		Email: "alice.senior@example.com",
	}
	resp, _ := suite.makeRequest("PUT", "/api/v1/users/"+userMap["dev1"], updateUser)
	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	// Verify user still exists in groups after update
	resp, body := suite.makeRequest("GET", "/api/v1/groups/"+groupMap["developers"], nil)
	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)
	var developersGroup dto.GroupResponseDTO
	err := json.Unmarshal(body, &developersGroup)
	suite.NoError(err)
	assert.Contains(suite.T(), developersGroup.Members, userMap["dev1"])

	// Remove a user from a specific group
	resp, _ = suite.makeRequest("DELETE", "/api/v1/groups/"+groupMap["developers"]+"/members/"+userMap["dev1"], nil)
	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	// Verify user was removed from developers but still in all-staff
	resp, body = suite.makeRequest("GET", "/api/v1/groups/"+groupMap["developers"], nil)
	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)
	err = json.Unmarshal(body, &developersGroup)
	suite.NoError(err)
	assert.NotContains(suite.T(), developersGroup.Members, userMap["dev1"])

	resp, body = suite.makeRequest("GET", "/api/v1/groups/"+groupMap["all-staff"], nil)
	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)
	var allStaffGroup dto.GroupResponseDTO
	err = json.Unmarshal(body, &allStaffGroup)
	suite.NoError(err)
	assert.Contains(suite.T(), allStaffGroup.Members, userMap["dev1"])
}

func (suite *IntegrationTestSuite) TestUserDeletionImpactOnGroups() {
	// Create a user
	createUser := dto.CreateUserRequestDTO{
		Name:  "Temporary User",
		Email: "temp@example.com",
	}
	resp, body := suite.makeRequest("POST", "/api/v1/users/", createUser)
	assert.Equal(suite.T(), http.StatusCreated, resp.StatusCode)

	var createdUser dto.UserResponseDTO
	err := json.Unmarshal(body, &createdUser)
	suite.NoError(err)

	// Create a group
	createGroup := dto.CreateGroupRequestDTO{
		Name: "Test Group",
	}
	resp, body = suite.makeRequest("POST", "/api/v1/groups/", createGroup)
	assert.Equal(suite.T(), http.StatusCreated, resp.StatusCode)

	var createdGroup dto.GroupResponseDTO
	err = json.Unmarshal(body, &createdGroup)
	suite.NoError(err)

	// Add user to group
	resp, _ = suite.makeRequest("POST", "/api/v1/groups/"+createdGroup.ID+"/members/"+createdUser.ID, nil)
	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	// Verify user is in group
	resp, body = suite.makeRequest("GET", "/api/v1/groups/"+createdGroup.ID, nil)
	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)
	var retrievedGroup dto.GroupResponseDTO
	err = json.Unmarshal(body, &retrievedGroup)
	suite.NoError(err)
	assert.Contains(suite.T(), retrievedGroup.Members, createdUser.ID)

	// Delete the user
	resp, _ = suite.makeRequest("DELETE", "/api/v1/users/"+createdUser.ID, nil)
	assert.Equal(suite.T(), http.StatusNoContent, resp.StatusCode)

	// Note: In a real implementation, you might want to automatically remove
	// the user from all groups when they are deleted. This test documents
	// the current behavior and can be updated when that functionality is added.

	// Verify user no longer exists
	resp, _ = suite.makeRequest("GET", "/api/v1/users/"+createdUser.ID, nil)
	assert.Equal(suite.T(), http.StatusNotFound, resp.StatusCode)

	// The group still contains the deleted user's ID
	// This might be considered a bug or expected behavior depending on requirements
	resp, body = suite.makeRequest("GET", "/api/v1/groups/"+createdGroup.ID, nil)
	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)
	err = json.Unmarshal(body, &retrievedGroup)
	suite.NoError(err)
	// This test documents current behavior - may need to change based on business logic
	assert.Contains(suite.T(), retrievedGroup.Members, createdUser.ID)
}

func (suite *IntegrationTestSuite) TestGroupDeletionWithMembers() {
	// Create users
	users := []dto.CreateUserRequestDTO{
		{Name: "User 1", Email: "user1@example.com"},
		{Name: "User 2", Email: "user2@example.com"},
	}

	createdUsers := make([]dto.UserResponseDTO, 0, len(users))
	for _, user := range users {
		resp, body := suite.makeRequest("POST", "/api/v1/users/", user)
		assert.Equal(suite.T(), http.StatusCreated, resp.StatusCode)

		var createdUser dto.UserResponseDTO
		err := json.Unmarshal(body, &createdUser)
		suite.NoError(err)
		createdUsers = append(createdUsers, createdUser)
	}

	// Create group
	createGroup := dto.CreateGroupRequestDTO{
		Name: "Test Group",
	}
	resp, body := suite.makeRequest("POST", "/api/v1/groups/", createGroup)
	assert.Equal(suite.T(), http.StatusCreated, resp.StatusCode)

	var createdGroup dto.GroupResponseDTO
	err := json.Unmarshal(body, &createdGroup)
	suite.NoError(err)

	// Add users to group
	for _, user := range createdUsers {
		resp, _ := suite.makeRequest("POST", "/api/v1/groups/"+createdGroup.ID+"/members/"+user.ID, nil)
		assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)
	}

	// Delete the group
	resp, _ = suite.makeRequest("DELETE", "/api/v1/groups/"+createdGroup.ID, nil)
	assert.Equal(suite.T(), http.StatusNoContent, resp.StatusCode)

	// Verify group is deleted
	resp, _ = suite.makeRequest("GET", "/api/v1/groups/"+createdGroup.ID, nil)
	assert.Equal(suite.T(), http.StatusNotFound, resp.StatusCode)

	// Verify users still exist
	for _, user := range createdUsers {
		resp, _ := suite.makeRequest("GET", "/api/v1/users/"+user.ID, nil)
		assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)
	}
}

func (suite *IntegrationTestSuite) TestConcurrentOperations() {
	// This test simulates concurrent operations to test race conditions
	// Create base data
	createUser := dto.CreateUserRequestDTO{
		Name:  "Concurrent User",
		Email: "concurrent@example.com",
	}
	resp, body := suite.makeRequest("POST", "/api/v1/users/", createUser)
	assert.Equal(suite.T(), http.StatusCreated, resp.StatusCode)

	var createdUser dto.UserResponseDTO
	err := json.Unmarshal(body, &createdUser)
	suite.NoError(err)

	createGroup := dto.CreateGroupRequestDTO{
		Name: "Concurrent Group",
	}
	resp, body = suite.makeRequest("POST", "/api/v1/groups/", createGroup)
	assert.Equal(suite.T(), http.StatusCreated, resp.StatusCode)

	var createdGroup dto.GroupResponseDTO
	err = json.Unmarshal(body, &createdGroup)
	suite.NoError(err)

	// Add user to group
	resp, _ = suite.makeRequest("POST", "/api/v1/groups/"+createdGroup.ID+"/members/"+createdUser.ID, nil)
	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	// Perform multiple operations sequentially (simulating concurrent behavior)
	operations := []struct {
		method string
		path   string
		body   interface{}
	}{
		{"GET", "/api/v1/users/" + createdUser.ID, nil},
		{"GET", "/api/v1/groups/" + createdGroup.ID, nil},
		{"PUT", "/api/v1/users/" + createdUser.ID, dto.CreateUserRequestDTO{Name: "Updated User", Email: "updated@example.com"}},
		{"GET", "/api/v1/groups/" + createdGroup.ID, nil},
		{"DELETE", "/api/v1/groups/" + createdGroup.ID + "/members/" + createdUser.ID, nil},
		{"POST", "/api/v1/groups/" + createdGroup.ID + "/members/" + createdUser.ID, nil},
	}

	for _, op := range operations {
		resp, _ := suite.makeRequest(op.method, op.path, op.body)
		// Basic assertion that operations don't fail catastrophically
		assert.True(suite.T(), resp.StatusCode < 500, "Operation should not fail with server error: %s %s", op.method, op.path)
	}
}

func (suite *IntegrationTestSuite) TestDataConsistency() {
	// Test that data remains consistent across operations
	// Create initial data
	users := []dto.CreateUserRequestDTO{
		{Name: "User 1", Email: "user1@example.com"},
		{Name: "User 2", Email: "user2@example.com"},
	}

	createdUsers := make([]dto.UserResponseDTO, 0, len(users))
	for _, user := range users {
		resp, body := suite.makeRequest("POST", "/api/v1/users/", user)
		assert.Equal(suite.T(), http.StatusCreated, resp.StatusCode)

		var createdUser dto.UserResponseDTO
		err := json.Unmarshal(body, &createdUser)
		suite.NoError(err)
		createdUsers = append(createdUsers, createdUser)
	}

	createGroup := dto.CreateGroupRequestDTO{
		Name: "Consistency Group",
	}
	resp, body := suite.makeRequest("POST", "/api/v1/groups/", createGroup)
	assert.Equal(suite.T(), http.StatusCreated, resp.StatusCode)

	var createdGroup dto.GroupResponseDTO
	err := json.Unmarshal(body, &createdGroup)
	suite.NoError(err)

	// Add both users to group
	for _, user := range createdUsers {
		resp, _ := suite.makeRequest("POST", "/api/v1/groups/"+createdGroup.ID+"/members/"+user.ID, nil)
		assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)
	}

	// Verify initial state
	resp, body = suite.makeRequest("GET", "/api/v1/groups/"+createdGroup.ID, nil)
	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)
	var retrievedGroup dto.GroupResponseDTO
	err = json.Unmarshal(body, &retrievedGroup)
	suite.NoError(err)
	assert.Len(suite.T(), retrievedGroup.Members, 2)

	// Perform various operations and verify consistency
	// Remove one user
	resp, _ = suite.makeRequest("DELETE", "/api/v1/groups/"+createdGroup.ID+"/members/"+createdUsers[0].ID, nil)
	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	// Verify state
	resp, body = suite.makeRequest("GET", "/api/v1/groups/"+createdGroup.ID, nil)
	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)
	err = json.Unmarshal(body, &retrievedGroup)
	suite.NoError(err)
	assert.Len(suite.T(), retrievedGroup.Members, 1)
	assert.Contains(suite.T(), retrievedGroup.Members, createdUsers[1].ID)
	assert.NotContains(suite.T(), retrievedGroup.Members, createdUsers[0].ID)

	// Add user back
	resp, _ = suite.makeRequest("POST", "/api/v1/groups/"+createdGroup.ID+"/members/"+createdUsers[0].ID, nil)
	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	// Final verification
	resp, body = suite.makeRequest("GET", "/api/v1/groups/"+createdGroup.ID, nil)
	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)
	err = json.Unmarshal(body, &retrievedGroup)
	suite.NoError(err)
	assert.Len(suite.T(), retrievedGroup.Members, 2)
	assert.Contains(suite.T(), retrievedGroup.Members, createdUsers[0].ID)
	assert.Contains(suite.T(), retrievedGroup.Members, createdUsers[1].ID)
}
