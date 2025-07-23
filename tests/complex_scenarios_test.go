package tests

import (
	"encoding/json"
	"net/http"
	"user-management/internal/application/dto"

	"github.com/stretchr/testify/assert"
)

func (suite *IntegrationTestSuite) TestCompleteUserGroupWorkflow() {
	// Create multiple users
	users := []dto.UserDTO{
		{ID: "dev1", Name: "Alice Developer", Email: "alice@example.com"},
		{ID: "dev2", Name: "Bob Developer", Email: "bob@example.com"},
		{ID: "designer1", Name: "Carol Designer", Email: "carol@example.com"},
		{ID: "manager1", Name: "Dave Manager", Email: "dave@example.com"},
	}

	for _, user := range users {
		resp, _ := suite.makeRequest("POST", "/api/v1/users/", user)
		assert.Equal(suite.T(), http.StatusCreated, resp.StatusCode)
	}

	// Create multiple groups
	groups := []dto.GroupDTO{
		{ID: "developers", Name: "Development Team", Members: []string{}},
		{ID: "designers", Name: "Design Team", Members: []string{}},
		{ID: "management", Name: "Management Team", Members: []string{}},
		{ID: "all-staff", Name: "All Staff", Members: []string{}},
	}

	for _, group := range groups {
		resp, _ := suite.makeRequest("POST", "/api/v1/groups/", group)
		assert.Equal(suite.T(), http.StatusCreated, resp.StatusCode)
	}

	// Add users to specific groups
	memberships := map[string][]string{
		"developers": {"dev1", "dev2"},
		"designers":  {"designer1"},
		"management": {"manager1"},
		"all-staff":  {"dev1", "dev2", "designer1", "manager1"},
	}

	for groupID, memberIDs := range memberships {
		for _, memberID := range memberIDs {
			resp, _ := suite.makeRequest("POST", "/api/v1/groups/"+groupID+"/members/"+memberID, nil)
			assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)
		}
	}

	// Verify group memberships
	for groupID, expectedMembers := range memberships {
		resp, body := suite.makeRequest("GET", "/api/v1/groups/"+groupID, nil)
		assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

		var group dto.GroupDTO
		err := json.Unmarshal(body, &group)
		suite.NoError(err)
		assert.Len(suite.T(), group.Members, len(expectedMembers))
		
		for _, expectedMember := range expectedMembers {
			assert.Contains(suite.T(), group.Members, expectedMember)
		}
	}

	// Update a user and verify they still exist in groups
	updatedUser := dto.UserDTO{
		ID:    "dev1",
		Name:  "Alice Senior Developer",
		Email: "alice.senior@example.com",
	}
	resp, _ := suite.makeRequest("PUT", "/api/v1/users/dev1", updatedUser)
	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	// Verify user still exists in groups after update
	resp, body := suite.makeRequest("GET", "/api/v1/groups/developers", nil)
	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)
	var developersGroup dto.GroupDTO
	err := json.Unmarshal(body, &developersGroup)
	suite.NoError(err)
	assert.Contains(suite.T(), developersGroup.Members, "dev1")

	// Remove a user from a specific group
	resp, _ = suite.makeRequest("DELETE", "/api/v1/groups/developers/members/dev1", nil)
	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	// Verify user was removed from developers but still in all-staff
	resp, body = suite.makeRequest("GET", "/api/v1/groups/developers", nil)
	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)
	err = json.Unmarshal(body, &developersGroup)
	suite.NoError(err)
	assert.NotContains(suite.T(), developersGroup.Members, "dev1")

	resp, body = suite.makeRequest("GET", "/api/v1/groups/all-staff", nil)
	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)
	var allStaffGroup dto.GroupDTO
	err = json.Unmarshal(body, &allStaffGroup)
	suite.NoError(err)
	assert.Contains(suite.T(), allStaffGroup.Members, "dev1")
}

func (suite *IntegrationTestSuite) TestUserDeletionImpactOnGroups() {
	// Create a user
	user := dto.UserDTO{
		ID:    "temp-user",
		Name:  "Temporary User",
		Email: "temp@example.com",
	}
	resp, _ := suite.makeRequest("POST", "/api/v1/users/", user)
	assert.Equal(suite.T(), http.StatusCreated, resp.StatusCode)

	// Create a group
	group := dto.GroupDTO{
		ID:      "test-group",
		Name:    "Test Group",
		Members: []string{},
	}
	resp, _ = suite.makeRequest("POST", "/api/v1/groups/", group)
	assert.Equal(suite.T(), http.StatusCreated, resp.StatusCode)

	// Add user to group
	resp, _ = suite.makeRequest("POST", "/api/v1/groups/test-group/members/temp-user", nil)
	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	// Verify user is in group
	resp, body := suite.makeRequest("GET", "/api/v1/groups/test-group", nil)
	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)
	var retrievedGroup dto.GroupDTO
	err := json.Unmarshal(body, &retrievedGroup)
	suite.NoError(err)
	assert.Contains(suite.T(), retrievedGroup.Members, "temp-user")

	// Delete the user
	resp, _ = suite.makeRequest("DELETE", "/api/v1/users/temp-user", nil)
	assert.Equal(suite.T(), http.StatusNoContent, resp.StatusCode)

	// Note: In a real implementation, you might want to automatically remove
	// the user from all groups when they are deleted. This test documents
	// the current behavior and can be updated when that functionality is added.
	
	// Verify user no longer exists
	resp, _ = suite.makeRequest("GET", "/api/v1/users/temp-user", nil)
	assert.Equal(suite.T(), http.StatusNotFound, resp.StatusCode)

	// The group still contains the deleted user's ID
	// This might be considered a bug or expected behavior depending on requirements
	resp, body = suite.makeRequest("GET", "/api/v1/groups/test-group", nil)
	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)
	err = json.Unmarshal(body, &retrievedGroup)
	suite.NoError(err)
	// This test documents current behavior - may need to change based on business logic
	assert.Contains(suite.T(), retrievedGroup.Members, "temp-user")
}

func (suite *IntegrationTestSuite) TestGroupDeletionWithMembers() {
	// Create users
	users := []dto.UserDTO{
		{ID: "user1", Name: "User 1", Email: "user1@example.com"},
		{ID: "user2", Name: "User 2", Email: "user2@example.com"},
	}

	for _, user := range users {
		resp, _ := suite.makeRequest("POST", "/api/v1/users/", user)
		assert.Equal(suite.T(), http.StatusCreated, resp.StatusCode)
	}

	// Create group with members
	group := dto.GroupDTO{
		ID:      "temp-group",
		Name:    "Temporary Group",
		Members: []string{},
	}
	resp, _ := suite.makeRequest("POST", "/api/v1/groups/", group)
	assert.Equal(suite.T(), http.StatusCreated, resp.StatusCode)

	// Add users to group
	for _, user := range users {
		resp, _ := suite.makeRequest("POST", "/api/v1/groups/temp-group/members/"+user.ID, nil)
		assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)
	}

	// Verify group has members
	resp, body := suite.makeRequest("GET", "/api/v1/groups/temp-group", nil)
	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)
	var retrievedGroup dto.GroupDTO
	err := json.Unmarshal(body, &retrievedGroup)
	suite.NoError(err)
	assert.Len(suite.T(), retrievedGroup.Members, 2)

	// Delete the group
	resp, _ = suite.makeRequest("DELETE", "/api/v1/groups/temp-group", nil)
	assert.Equal(suite.T(), http.StatusNoContent, resp.StatusCode)

	// Verify group no longer exists
	resp, _ = suite.makeRequest("GET", "/api/v1/groups/temp-group", nil)
	assert.Equal(suite.T(), http.StatusNotFound, resp.StatusCode)

	// Verify users still exist
	for _, user := range users {
		resp, _ := suite.makeRequest("GET", "/api/v1/users/"+user.ID, nil)
		assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)
	}
}

func (suite *IntegrationTestSuite) TestConcurrentOperations() {
	// This test simulates concurrent operations to test race conditions
	// Create base data
	user := dto.UserDTO{
		ID:    "concurrent-user",
		Name:  "Concurrent User",
		Email: "concurrent@example.com",
	}
	resp, _ := suite.makeRequest("POST", "/api/v1/users/", user)
	assert.Equal(suite.T(), http.StatusCreated, resp.StatusCode)

	group := dto.GroupDTO{
		ID:      "concurrent-group",
		Name:    "Concurrent Group",
		Members: []string{},
	}
	resp, _ = suite.makeRequest("POST", "/api/v1/groups/", group)
	assert.Equal(suite.T(), http.StatusCreated, resp.StatusCode)

	// Add user to group
	resp, _ = suite.makeRequest("POST", "/api/v1/groups/concurrent-group/members/concurrent-user", nil)
	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	// Perform multiple operations sequentially (simulating concurrent behavior)
	operations := []struct {
		method string
		path   string
		body   interface{}
	}{
		{"GET", "/api/v1/users/concurrent-user", nil},
		{"GET", "/api/v1/groups/concurrent-group", nil},
		{"PUT", "/api/v1/users/concurrent-user", dto.UserDTO{ID: "concurrent-user", Name: "Updated User", Email: "updated@example.com"}},
		{"GET", "/api/v1/groups/concurrent-group", nil},
		{"DELETE", "/api/v1/groups/concurrent-group/members/concurrent-user", nil},
		{"POST", "/api/v1/groups/concurrent-group/members/concurrent-user", nil},
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
	users := []dto.UserDTO{
		{ID: "consistency-user1", Name: "User 1", Email: "user1@example.com"},
		{ID: "consistency-user2", Name: "User 2", Email: "user2@example.com"},
	}

	for _, user := range users {
		resp, _ := suite.makeRequest("POST", "/api/v1/users/", user)
		assert.Equal(suite.T(), http.StatusCreated, resp.StatusCode)
	}

	group := dto.GroupDTO{
		ID:      "consistency-group",
		Name:    "Consistency Group",
		Members: []string{},
	}
	resp, _ := suite.makeRequest("POST", "/api/v1/groups/", group)
	assert.Equal(suite.T(), http.StatusCreated, resp.StatusCode)

	// Add both users to group
	for _, user := range users {
		resp, _ := suite.makeRequest("POST", "/api/v1/groups/consistency-group/members/"+user.ID, nil)
		assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)
	}

	// Verify initial state
	resp, body := suite.makeRequest("GET", "/api/v1/groups/consistency-group", nil)
	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)
	var retrievedGroup dto.GroupDTO
	err := json.Unmarshal(body, &retrievedGroup)
	suite.NoError(err)
	assert.Len(suite.T(), retrievedGroup.Members, 2)

	// Perform various operations and verify consistency
	// Remove one user
	resp, _ = suite.makeRequest("DELETE", "/api/v1/groups/consistency-group/members/consistency-user1", nil)
	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	// Verify state
	resp, body = suite.makeRequest("GET", "/api/v1/groups/consistency-group", nil)
	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)
	err = json.Unmarshal(body, &retrievedGroup)
	suite.NoError(err)
	assert.Len(suite.T(), retrievedGroup.Members, 1)
	assert.Contains(suite.T(), retrievedGroup.Members, "consistency-user2")
	assert.NotContains(suite.T(), retrievedGroup.Members, "consistency-user1")

	// Add user back
	resp, _ = suite.makeRequest("POST", "/api/v1/groups/consistency-group/members/consistency-user1", nil)
	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	// Final verification
	resp, body = suite.makeRequest("GET", "/api/v1/groups/consistency-group", nil)
	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)
	err = json.Unmarshal(body, &retrievedGroup)
	suite.NoError(err)
	assert.Len(suite.T(), retrievedGroup.Members, 2)
	assert.Contains(suite.T(), retrievedGroup.Members, "consistency-user1")
	assert.Contains(suite.T(), retrievedGroup.Members, "consistency-user2")
}
