package tests

import (
	"context"
	"errors"
	"user-management/internal/application/dto"
	"user-management/internal/application/usecases/group"
	"user-management/internal/application/usecases/user"
	"user-management/internal/domain/entities"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/v2/bson"
)

const (
	userEntityTypeIntegration  = "*entities.User"
	groupEntityTypeIntegration = "*entities.Group"
	testUserName               = "Test User"
	testUserEmail              = "test@example.com"
	testGroupName              = "Test Group"
	updatedUserEmail           = "updated@example.com"
	nonexistentGroupID         = "nonexistent-group-id"
	nonexistentUserID          = "nonexistent-user-id"
)

// Testes de integração com simulação de erros de repository

func (suite *IntegrationTestSuite) TestUserCRUDWithRepositoryErrors() {
	// Teste que simula erro no repository durante operações CRUD

	// Criar um mock do UserRepository
	mockUserRepo := new(MockUserRepository)

	// Configurar erro no Create
	repositoryError := errors.New("database connection failed during user creation")
	mockUserRepo.On("Create", mock.Anything, mock.AnythingOfType(userEntityTypeIntegration)).Return(repositoryError)

	// Criar usecase com o mock
	createUserUseCase := user.NewCreateUserUseCase(mockUserRepo)

	// Executar teste de criação que deve falhar
	createUserDTO := dto.CreateUserRequestDTO{
		Name:  testUserName,
		Email: testUserEmail,
	}

	_, err := createUserUseCase.Execute(context.Background(), &createUserDTO)

	// Verificar que o erro foi propagado corretamente
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), repositoryError, err)
	mockUserRepo.AssertExpectations(suite.T())
}

func (suite *IntegrationTestSuite) TestUserGetWithRepositoryError() {
	// Teste que simula erro no repository durante busca por ID

	mockUserRepo := new(MockUserRepository)

	// Configurar erro no GetByID
	repositoryError := errors.New("user not found in database")
	mockUserRepo.On("GetByID", mock.Anything, "nonexistent-id").Return(nil, repositoryError)

	// Criar usecase com o mock
	getUserUseCase := user.NewGetUserUseCase(mockUserRepo)

	// Executar teste
	result, err := getUserUseCase.Execute(context.Background(), "nonexistent-id")

	// Verificar resultado
	assert.Nil(suite.T(), result)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), repositoryError, err)
	mockUserRepo.AssertExpectations(suite.T())
}

func (suite *IntegrationTestSuite) TestUserListWithRepositoryError() {
	// Teste que simula erro no repository durante listagem

	mockUserRepo := new(MockUserRepository)

	// Configurar erro no List
	repositoryError := errors.New("database query failed")
	mockUserRepo.On("List", mock.Anything).Return(nil, repositoryError)

	// Criar usecase com o mock
	listUsersUseCase := user.NewListUsersUseCase(mockUserRepo)

	// Executar teste
	result, err := listUsersUseCase.Execute(context.Background())

	// Verificar resultado
	assert.Nil(suite.T(), result)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), repositoryError, err)
	mockUserRepo.AssertExpectations(suite.T())
}

func (suite *IntegrationTestSuite) TestUserUpdateWithRepositoryError() {
	// Teste que simula erro no repository durante atualização

	mockUserRepo := new(MockUserRepository)

	// Configurar GetByID para sucesso (necessário para o usecase)
	userObjectID := bson.NewObjectID()
	existingUser := &entities.User{
		ID:    userObjectID,
		Name:  "Existing User",
		Email: "existing@example.com",
	}
	mockUserRepo.On("GetByID", mock.Anything, "user123").Return(existingUser, nil)

	// Configurar erro no Update
	repositoryError := errors.New("update operation failed")
	mockUserRepo.On("Update", mock.Anything, mock.AnythingOfType(userEntityTypeIntegration)).Return(repositoryError)

	// Criar usecase com o mock
	updateUserUseCase := user.NewUpdateUserUseCase(mockUserRepo)

	// Executar teste
	userDTO := &dto.CreateUserRequestDTO{
		Name:  "Updated User",
		Email: updatedUserEmail,
	}

	result, err := updateUserUseCase.Execute(context.Background(), "user123", userDTO)

	// Verificar resultado
	assert.Nil(suite.T(), result)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), repositoryError, err)
	mockUserRepo.AssertExpectations(suite.T())
}

func (suite *IntegrationTestSuite) TestUserDeleteWithRepositoryError() {
	// Teste que simula erro no repository durante exclusão

	mockUserRepo := new(MockUserRepository)

	// Configurar erro no Delete
	repositoryError := errors.New("delete operation failed")
	mockUserRepo.On("Delete", mock.Anything, "user123").Return(repositoryError)

	// Criar usecase com o mock
	deleteUserUseCase := user.NewDeleteUserUseCase(mockUserRepo)

	// Executar teste
	err := deleteUserUseCase.Execute(context.Background(), "user123")

	// Verificar resultado
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), repositoryError, err)
	mockUserRepo.AssertExpectations(suite.T())
}

func (suite *IntegrationTestSuite) TestGroupCRUDWithRepositoryErrors() {
	// Teste que simula erro no repository durante operações CRUD de grupos

	mockGroupRepo := new(MockGroupRepository)

	// Configurar erro no Create
	repositoryError := errors.New("database connection failed during group creation")
	mockGroupRepo.On("Create", mock.Anything, mock.AnythingOfType(groupEntityTypeIntegration)).Return(repositoryError)

	// Criar usecase com o mock
	createGroupUseCase := group.NewCreateGroupUseCase(mockGroupRepo)

	// Executar teste de criação que deve falhar
	createGroupDTO := dto.CreateGroupRequestDTO{
		Name:    testGroupName,
		Members: []string{},
	}

	result, err := createGroupUseCase.Execute(context.Background(), &createGroupDTO)

	// Verificar que o erro foi propagado corretamente
	assert.Nil(suite.T(), result)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), repositoryError, err)
	mockGroupRepo.AssertExpectations(suite.T())
}

func (suite *IntegrationTestSuite) TestGroupGetWithRepositoryError() {
	// Teste que simula erro no repository durante busca por ID de grupo

	mockGroupRepo := new(MockGroupRepository)

	// Configurar erro no GetByID
	repositoryError := errors.New("group not found in database")
	mockGroupRepo.On("GetByID", mock.Anything, nonexistentGroupID).Return(nil, repositoryError)

	// Criar usecase com o mock
	getGroupUseCase := group.NewGetGroupUseCase(mockGroupRepo)

	// Executar teste
	result, err := getGroupUseCase.Execute(context.Background(), nonexistentGroupID)

	// Verificar resultado
	assert.Nil(suite.T(), result)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), repositoryError, err)
	mockGroupRepo.AssertExpectations(suite.T())
}

func (suite *IntegrationTestSuite) TestAddUserToGroupWithRepositoryError() {
	// Teste que simula erro no repository durante adição de usuário ao grupo

	mockGroupRepo := new(MockGroupRepository)
	mockUserRepo := new(MockUserRepository)

	// Configurar GetByID para sucesso em ambos repositórios (necessário para o usecase)
	userObjectID := bson.NewObjectID()
	groupObjectID := bson.NewObjectID()

	existingUser := &entities.User{
		ID:    userObjectID,
		Name:  testUserName,
		Email: testUserEmail,
	}
	existingGroup := &entities.Group{
		ID:      groupObjectID,
		Name:    testGroupName,
		Members: []string{},
	}

	mockGroupRepo.On("GetByID", mock.Anything, "group123").Return(existingGroup, nil)
	mockUserRepo.On("GetByID", mock.Anything, "user123").Return(existingUser, nil)

	// Configurar erro no AddUserToGroup
	repositoryError := errors.New("failed to add user to group")
	mockGroupRepo.On("AddUserToGroup", mock.Anything, "group123", "user123").Return(repositoryError)

	// Criar usecase com o mock
	addUserToGroupUseCase := group.NewAddUserToGroupUseCase(mockGroupRepo, mockUserRepo)

	// Executar teste
	err := addUserToGroupUseCase.Execute(context.Background(), "group123", "user123")

	// Verificar resultado
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), repositoryError, err)
	mockGroupRepo.AssertExpectations(suite.T())
	mockUserRepo.AssertExpectations(suite.T())
}

func (suite *IntegrationTestSuite) TestRemoveUserFromGroupWithRepositoryError() {
	// Teste que simula erro no repository durante remoção de usuário do grupo

	mockGroupRepo := new(MockGroupRepository)

	// Configurar erro no RemoveUserFromGroup
	repositoryError := errors.New("failed to remove user from group")
	mockGroupRepo.On("RemoveUserFromGroup", mock.Anything, "group123", "user123").Return(repositoryError)

	// Criar usecase com o mock
	removeUserFromGroupUseCase := group.NewRemoveUserFromGroupUseCase(mockGroupRepo)

	// Executar teste
	err := removeUserFromGroupUseCase.Execute(context.Background(), "group123", "user123")

	// Verificar resultado
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), repositoryError, err)
	mockGroupRepo.AssertExpectations(suite.T())
}

func (suite *IntegrationTestSuite) TestDatabaseTimeoutSimulation() {
	// Teste que simula timeout do banco de dados

	mockUserRepo := new(MockUserRepository)

	// Simular timeout
	timeoutError := errors.New("context deadline exceeded: database operation timeout")
	mockUserRepo.On("Create", mock.Anything, mock.AnythingOfType(userEntityTypeIntegration)).Return(timeoutError)

	// Criar usecase com o mock
	createUserUseCase := user.NewCreateUserUseCase(mockUserRepo)

	// Executar teste
	userDTO := &dto.CreateUserRequestDTO{
		Name:  testUserName,
		Email: testUserEmail,
	}

	result, err := createUserUseCase.Execute(context.Background(), userDTO)

	// Verificar resultado
	assert.Nil(suite.T(), result)
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "timeout")
	mockUserRepo.AssertExpectations(suite.T())
}

func (suite *IntegrationTestSuite) TestDatabaseConnectionFailureSimulation() {
	// Teste que simula falha de conexão com banco de dados

	mockUserRepo := new(MockUserRepository)

	// Simular falha de conexão
	connectionError := errors.New("no reachable servers: connection refused")
	mockUserRepo.On("List", mock.Anything).Return(nil, connectionError)

	// Criar usecase com o mock
	listUsersUseCase := user.NewListUsersUseCase(mockUserRepo)

	// Executar teste
	result, err := listUsersUseCase.Execute(context.Background())

	// Verificar resultado
	assert.Nil(suite.T(), result)
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "connection")
	mockUserRepo.AssertExpectations(suite.T())
}

func (suite *IntegrationTestSuite) TestMultipleRepositoryErrorsSimulation() {
	// Teste que simula múltiplos erros de repository em sequência

	mockUserRepo := new(MockUserRepository)
	mockGroupRepo := new(MockGroupRepository)

	// Configurar erros em ambos repositórios
	userError := errors.New("user repository error")
	groupError := errors.New("group repository error")

	mockUserRepo.On("Create", mock.Anything, mock.AnythingOfType(userEntityTypeIntegration)).Return(userError)
	mockGroupRepo.On("Create", mock.Anything, mock.AnythingOfType(groupEntityTypeIntegration)).Return(groupError)

	// Criar usecases com mocks
	createUserUseCase := user.NewCreateUserUseCase(mockUserRepo)
	createGroupUseCase := group.NewCreateGroupUseCase(mockGroupRepo)

	// Executar testes
	userDTO := &dto.CreateUserRequestDTO{
		Name:  testUserName,
		Email: testUserEmail,
	}

	groupDTO := &dto.CreateGroupRequestDTO{
		Name:    testGroupName,
		Members: []string{},
	}

	userResult, userErr := createUserUseCase.Execute(context.Background(), userDTO)
	groupResult, groupErr := createGroupUseCase.Execute(context.Background(), groupDTO)

	// Verificar resultados
	assert.Nil(suite.T(), userResult)
	assert.Error(suite.T(), userErr)
	assert.Equal(suite.T(), userError, userErr)

	assert.Nil(suite.T(), groupResult)
	assert.Error(suite.T(), groupErr)
	assert.Equal(suite.T(), groupError, groupErr)

	mockUserRepo.AssertExpectations(suite.T())
	mockGroupRepo.AssertExpectations(suite.T())
}

// Teste de erro específico para operações de busca que não encontram dados
func (suite *IntegrationTestSuite) TestEntityNotFoundRepositoryError() {
	mockUserRepo := new(MockUserRepository)

	// Simular erro de entidade não encontrada
	notFoundError := errors.New("entity not found")
	mockUserRepo.On("GetByID", mock.Anything, "invalid-id").Return(nil, notFoundError)

	// Criar usecase com o mock
	getUserUseCase := user.NewGetUserUseCase(mockUserRepo)

	// Executar teste
	result, err := getUserUseCase.Execute(context.Background(), "invalid-id")

	// Verificar resultado
	assert.Nil(suite.T(), result)
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "not found")
	mockUserRepo.AssertExpectations(suite.T())
}

// Teste de erro para operações que violam constraints do banco
func (suite *IntegrationTestSuite) TestDatabaseConstraintViolationError() {
	mockUserRepo := new(MockUserRepository)

	// Simular erro de violação de constraint (ex: email duplicado)
	constraintError := errors.New("duplicate key error: email already exists")
	mockUserRepo.On("Create", mock.Anything, mock.AnythingOfType(userEntityTypeIntegration)).Return(constraintError)

	// Criar usecase com o mock
	createUserUseCase := user.NewCreateUserUseCase(mockUserRepo)

	// Executar teste
	userDTO := &dto.CreateUserRequestDTO{
		Name:  testUserName,
		Email: "duplicate@example.com",
	}

	result, err := createUserUseCase.Execute(context.Background(), userDTO)

	// Verificar resultado
	assert.Nil(suite.T(), result)
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "duplicate")
	mockUserRepo.AssertExpectations(suite.T())
}

// Testes específicos para erros de GetByID nos use cases de update

func (suite *IntegrationTestSuite) TestUserUpdateWithGetByIDRepositoryError() {
	// Teste que simula erro no GetByID durante operação de update de usuário

	mockUserRepo := new(MockUserRepository)

	// Configurar erro no GetByID - usuário não encontrado
	repositoryError := errors.New("user not found: GetByID failed")
	mockUserRepo.On("GetByID", mock.Anything, nonexistentUserID).Return(nil, repositoryError)

	// Criar usecase com o mock
	updateUserUseCase := user.NewUpdateUserUseCase(mockUserRepo)

	// Executar teste
	userDTO := &dto.CreateUserRequestDTO{
		Name:  "Updated User Name",
		Email: updatedUserEmail,
	}

	result, err := updateUserUseCase.Execute(context.Background(), nonexistentUserID, userDTO)

	// Verificar resultado - erro deve ser propagado do GetByID
	assert.Nil(suite.T(), result)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), repositoryError, err)
	assert.Contains(suite.T(), err.Error(), "GetByID failed")

	// Verificar que Update não foi chamado devido ao erro no GetByID
	mockUserRepo.AssertExpectations(suite.T())
	mockUserRepo.AssertNotCalled(suite.T(), "Update")
}

func (suite *IntegrationTestSuite) TestGroupUpdateWithGetByIDRepositoryError() {
	// Teste que simula erro no GetByID durante operação de update de grupo

	mockGroupRepo := new(MockGroupRepository)

	// Configurar erro no GetByID - grupo não encontrado
	repositoryError := errors.New("group not found: GetByID failed")
	mockGroupRepo.On("GetByID", mock.Anything, nonexistentGroupID).Return(nil, repositoryError)

	// Criar usecase com o mock
	updateGroupUseCase := group.NewUpdateGroupUseCase(mockGroupRepo)

	// Executar teste
	groupDTO := &dto.CreateGroupRequestDTO{
		Name:    "Updated Group Name",
		Members: []string{"user1", "user2"},
	}

	result, err := updateGroupUseCase.Execute(context.Background(), nonexistentGroupID, groupDTO)

	// Verificar resultado - erro deve ser propagado do GetByID
	assert.Nil(suite.T(), result)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), repositoryError, err)
	assert.Contains(suite.T(), err.Error(), "GetByID failed")

	// Verificar que Update não foi chamado devido ao erro no GetByID
	mockGroupRepo.AssertExpectations(suite.T())
	mockGroupRepo.AssertNotCalled(suite.T(), "Update")
}

func (suite *IntegrationTestSuite) TestUserUpdateWithGetByIDDatabaseConnectionError() {
	// Teste que simula erro de conexão com banco durante GetByID no update de usuário

	mockUserRepo := new(MockUserRepository)

	// Configurar erro de conexão no GetByID
	connectionError := errors.New("connection lost: unable to reach database server")
	mockUserRepo.On("GetByID", mock.Anything, "user123").Return(nil, connectionError)

	// Criar usecase com o mock
	updateUserUseCase := user.NewUpdateUserUseCase(mockUserRepo)

	// Executar teste
	userDTO := &dto.CreateUserRequestDTO{
		Name:  "Updated User",
		Email: updatedUserEmail,
	}

	result, err := updateUserUseCase.Execute(context.Background(), "user123", userDTO)

	// Verificar resultado
	assert.Nil(suite.T(), result)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), connectionError, err)
	assert.Contains(suite.T(), err.Error(), "connection lost")

	mockUserRepo.AssertExpectations(suite.T())
	mockUserRepo.AssertNotCalled(suite.T(), "Update")
}

func (suite *IntegrationTestSuite) TestGroupUpdateWithGetByIDTimeoutError() {
	// Teste que simula timeout durante GetByID no update de grupo

	mockGroupRepo := new(MockGroupRepository)

	// Configurar erro de timeout no GetByID
	timeoutError := errors.New("context deadline exceeded: GetByID operation timeout")
	mockGroupRepo.On("GetByID", mock.Anything, "group123").Return(nil, timeoutError)

	// Criar usecase com o mock
	updateGroupUseCase := group.NewUpdateGroupUseCase(mockGroupRepo)

	// Executar teste
	groupDTO := &dto.CreateGroupRequestDTO{
		Name:    "Updated Group",
		Members: []string{"member1", "member2"},
	}

	result, err := updateGroupUseCase.Execute(context.Background(), "group123", groupDTO)

	// Verificar resultado
	assert.Nil(suite.T(), result)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), timeoutError, err)
	assert.Contains(suite.T(), err.Error(), "timeout")

	mockGroupRepo.AssertExpectations(suite.T())
	mockGroupRepo.AssertNotCalled(suite.T(), "Update")
}

func (suite *IntegrationTestSuite) TestUserUpdateWithGetByIDPermissionError() {
	// Teste que simula erro de permissão durante GetByID no update de usuário

	mockUserRepo := new(MockUserRepository)

	// Configurar erro de permissão no GetByID
	permissionError := errors.New("access denied: insufficient permissions to read user")
	mockUserRepo.On("GetByID", mock.Anything, "restricted-user-id").Return(nil, permissionError)

	// Criar usecase com o mock
	updateUserUseCase := user.NewUpdateUserUseCase(mockUserRepo)

	// Executar teste
	userDTO := &dto.CreateUserRequestDTO{
		Name:  "Should Not Update",
		Email: "shouldnot@example.com",
	}

	result, err := updateUserUseCase.Execute(context.Background(), "restricted-user-id", userDTO)

	// Verificar resultado
	assert.Nil(suite.T(), result)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), permissionError, err)
	assert.Contains(suite.T(), err.Error(), "access denied")

	mockUserRepo.AssertExpectations(suite.T())
	mockUserRepo.AssertNotCalled(suite.T(), "Update")
}

func (suite *IntegrationTestSuite) TestGroupUpdateWithGetByIDCorruptedDataError() {
	// Teste que simula erro de dados corrompidos durante GetByID no update de grupo

	mockGroupRepo := new(MockGroupRepository)

	// Configurar erro de dados corrompidos no GetByID
	corruptionError := errors.New("data corruption detected: unable to deserialize group data")
	mockGroupRepo.On("GetByID", mock.Anything, "corrupted-group-id").Return(nil, corruptionError)

	// Criar usecase com o mock
	updateGroupUseCase := group.NewUpdateGroupUseCase(mockGroupRepo)

	// Executar teste
	groupDTO := &dto.CreateGroupRequestDTO{
		Name:    "Cannot Update Corrupted",
		Members: []string{},
	}

	result, err := updateGroupUseCase.Execute(context.Background(), "corrupted-group-id", groupDTO)

	// Verificar resultado
	assert.Nil(suite.T(), result)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), corruptionError, err)
	assert.Contains(suite.T(), err.Error(), "corruption")

	mockGroupRepo.AssertExpectations(suite.T())
	mockGroupRepo.AssertNotCalled(suite.T(), "Update")
}
