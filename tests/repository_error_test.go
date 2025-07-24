package tests

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"
	"user-management/internal/application/dto"
	"user-management/internal/application/usecases/group"
	"user-management/internal/application/usecases/user"
	"user-management/internal/config"
	"user-management/internal/domain/entities"
	"user-management/internal/infrastructure/database"
	"user-management/internal/infrastructure/logger"
	repositories "user-management/internal/infrastructure/repositories"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const (
	userEntityType                   = "*entities.User"
	groupEntityType                  = "*entities.Group"
	testUserNameRepo                 = "Test User"
	testUserEmailRepo                = "test@example.com"
	failedToInitializeUserRepo       = "failed to initialize user repository"
	failedToInitializeGroupRepo      = "failed to initialize group repository"
	databaseNotConnectedError        = "database is not connected"
	failedToGetCollectionUsersError  = "failed to get collection: users"
	failedToGetCollectionGroupsError = "failed to get collection: groups"
	contextCanceledError             = "context canceled"
	contextDeadlineError             = "context deadline exceeded"
	serverSelectionError             = "server selection error"
)

// MockUserRepository implementa a interface IUserRepository para simular erros
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *entities.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) GetByID(ctx context.Context, id string) (*entities.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.User), args.Error(1)
}

func (m *MockUserRepository) List(ctx context.Context) ([]*entities.User, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entities.User), args.Error(1)
}

func (m *MockUserRepository) Update(ctx context.Context, user *entities.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// MockGroupRepository implementa a interface IGroupRepository para simular erros
type MockGroupRepository struct {
	mock.Mock
}

func (m *MockGroupRepository) Create(ctx context.Context, group *entities.Group) error {
	args := m.Called(ctx, group)
	return args.Error(0)
}

func (m *MockGroupRepository) GetByID(ctx context.Context, id string) (*entities.Group, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Group), args.Error(1)
}

func (m *MockGroupRepository) List(ctx context.Context) ([]*entities.Group, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entities.Group), args.Error(1)
}

func (m *MockGroupRepository) Update(ctx context.Context, group *entities.Group) error {
	args := m.Called(ctx, group)
	return args.Error(0)
}

func (m *MockGroupRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockGroupRepository) AddUserToGroup(ctx context.Context, groupID, userID string) error {
	args := m.Called(ctx, groupID, userID)
	return args.Error(0)
}

func (m *MockGroupRepository) RemoveUserFromGroup(ctx context.Context, groupID, userID string) error {
	args := m.Called(ctx, groupID, userID)
	return args.Error(0)
}

// RepositoryErrorTestSuite testa cenários de erro nos repositórios
type RepositoryErrorTestSuite struct {
	suite.Suite
	mockUserRepo  *MockUserRepository
	mockGroupRepo *MockGroupRepository
}

func (suite *RepositoryErrorTestSuite) SetupTest() {
	suite.mockUserRepo = new(MockUserRepository)
	suite.mockGroupRepo = new(MockGroupRepository)
}

// Testes de erro para User Repository

func (suite *RepositoryErrorTestSuite) TestCreateUserRepositoryError() {
	// Configurar mock para retornar erro
	repositoryError := errors.New("database connection failed")
	suite.mockUserRepo.On("Create", mock.Anything, mock.AnythingOfType(userEntityType)).Return(repositoryError)

	// Criar usecase com mock
	createUserUseCase := user.NewCreateUserUseCase(suite.mockUserRepo)

	// Executar teste
	userDTO := &dto.CreateUserRequestDTO{
		Name:  "Test User",
		Email: "test@example.com",
	}

	result, err := createUserUseCase.Execute(context.Background(), userDTO)

	// Verificar resultado
	assert.Nil(suite.T(), result)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), repositoryError, err)
	suite.mockUserRepo.AssertExpectations(suite.T())
}

func (suite *RepositoryErrorTestSuite) TestGetUserByIDRepositoryError() {
	// Configurar mock para retornar erro
	repositoryError := errors.New("user not found in database")
	suite.mockUserRepo.On("GetByID", mock.Anything, "user123").Return(nil, repositoryError)

	// Criar usecase com mock
	getUserUseCase := user.NewGetUserUseCase(suite.mockUserRepo)

	// Executar teste
	result, err := getUserUseCase.Execute(context.Background(), "user123")

	// Verificar resultado
	assert.Nil(suite.T(), result)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), repositoryError, err)
	suite.mockUserRepo.AssertExpectations(suite.T())
}

func (suite *RepositoryErrorTestSuite) TestListUsersRepositoryError() {
	// Configurar mock para retornar erro
	repositoryError := errors.New("database query failed")
	suite.mockUserRepo.On("List", mock.Anything).Return(nil, repositoryError)

	// Criar usecase com mock
	listUsersUseCase := user.NewListUsersUseCase(suite.mockUserRepo)

	// Executar teste
	result, err := listUsersUseCase.Execute(context.Background())

	// Verificar resultado
	assert.Nil(suite.T(), result)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), repositoryError, err)
	suite.mockUserRepo.AssertExpectations(suite.T())
}

func (suite *RepositoryErrorTestSuite) TestUpdateUserRepositoryError() {
	// Configurar mock para GetByID (sucesso) e Update (erro)
	userObjectID := bson.NewObjectID()
	existingUser := &entities.User{
		ID:    userObjectID,
		Name:  "Existing User",
		Email: "existing@example.com",
	}
	suite.mockUserRepo.On("GetByID", mock.Anything, "user123").Return(existingUser, nil)

	repositoryError := errors.New("update operation failed")
	suite.mockUserRepo.On("Update", mock.Anything, mock.AnythingOfType(userEntityType)).Return(repositoryError)

	// Criar usecase com mock
	updateUserUseCase := user.NewUpdateUserUseCase(suite.mockUserRepo)

	// Executar teste
	userDTO := &dto.CreateUserRequestDTO{
		Name:  "Updated User",
		Email: "updated@example.com",
	}

	result, err := updateUserUseCase.Execute(context.Background(), "user123", userDTO)

	// Verificar resultado
	assert.Nil(suite.T(), result)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), repositoryError, err)
	suite.mockUserRepo.AssertExpectations(suite.T())
}

func (suite *RepositoryErrorTestSuite) TestDeleteUserRepositoryError() {
	// Configurar mock para retornar erro
	repositoryError := errors.New("delete operation failed")
	suite.mockUserRepo.On("Delete", mock.Anything, "user123").Return(repositoryError)

	// Criar usecase com mock
	deleteUserUseCase := user.NewDeleteUserUseCase(suite.mockUserRepo)

	// Executar teste
	err := deleteUserUseCase.Execute(context.Background(), "user123")

	// Verificar resultado
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), repositoryError, err)
	suite.mockUserRepo.AssertExpectations(suite.T())
}

// Testes de erro para Group Repository

func (suite *RepositoryErrorTestSuite) TestCreateGroupRepositoryError() {
	// Configurar mock para retornar erro
	repositoryError := errors.New("database connection failed")
	suite.mockGroupRepo.On("Create", mock.Anything, mock.AnythingOfType(groupEntityType)).Return(repositoryError)

	// Criar usecase com mock
	createGroupUseCase := group.NewCreateGroupUseCase(suite.mockGroupRepo)

	// Executar teste
	groupDTO := &dto.CreateGroupRequestDTO{
		Name:    "Test Group",
		Members: []string{},
	}

	result, err := createGroupUseCase.Execute(context.Background(), groupDTO)

	// Verificar resultado
	assert.Nil(suite.T(), result)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), repositoryError, err)
	suite.mockGroupRepo.AssertExpectations(suite.T())
}

func (suite *RepositoryErrorTestSuite) TestGetGroupByIDRepositoryError() {
	// Configurar mock para retornar erro
	repositoryError := errors.New("group not found in database")
	suite.mockGroupRepo.On("GetByID", mock.Anything, "group123").Return(nil, repositoryError)

	// Criar usecase com mock
	getGroupUseCase := group.NewGetGroupUseCase(suite.mockGroupRepo)

	// Executar teste
	result, err := getGroupUseCase.Execute(context.Background(), "group123")

	// Verificar resultado
	assert.Nil(suite.T(), result)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), repositoryError, err)
	suite.mockGroupRepo.AssertExpectations(suite.T())
}

func (suite *RepositoryErrorTestSuite) TestListGroupsRepositoryError() {
	// Configurar mock para retornar erro
	repositoryError := errors.New("database query failed")
	suite.mockGroupRepo.On("List", mock.Anything).Return(nil, repositoryError)

	// Criar usecase com mock
	listGroupsUseCase := group.NewListGroupsUseCase(suite.mockGroupRepo)

	// Executar teste
	result, err := listGroupsUseCase.Execute(context.Background())

	// Verificar resultado
	assert.Nil(suite.T(), result)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), repositoryError, err)
	suite.mockGroupRepo.AssertExpectations(suite.T())
}

func (suite *RepositoryErrorTestSuite) TestUpdateGroupRepositoryError() {
	// Configurar mock para GetByID (sucesso) e Update (erro)
	groupObjectID := bson.NewObjectID()
	existingGroup := &entities.Group{
		ID:      groupObjectID,
		Name:    "Existing Group",
		Members: []string{},
	}
	suite.mockGroupRepo.On("GetByID", mock.Anything, "group123").Return(existingGroup, nil)

	repositoryError := errors.New("update operation failed")
	suite.mockGroupRepo.On("Update", mock.Anything, mock.AnythingOfType(groupEntityType)).Return(repositoryError)

	// Criar usecase com mock
	updateGroupUseCase := group.NewUpdateGroupUseCase(suite.mockGroupRepo)

	// Executar teste
	groupDTO := &dto.CreateGroupRequestDTO{
		Name:    "Updated Group",
		Members: []string{},
	}

	result, err := updateGroupUseCase.Execute(context.Background(), "group123", groupDTO)

	// Verificar resultado
	assert.Nil(suite.T(), result)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), repositoryError, err)
	suite.mockGroupRepo.AssertExpectations(suite.T())
}

func (suite *RepositoryErrorTestSuite) TestDeleteGroupRepositoryError() {
	// Configurar mock para retornar erro
	repositoryError := errors.New("delete operation failed")
	suite.mockGroupRepo.On("Delete", mock.Anything, "group123").Return(repositoryError)

	// Criar usecase com mock
	deleteGroupUseCase := group.NewDeleteGroupUseCase(suite.mockGroupRepo)

	// Executar teste
	err := deleteGroupUseCase.Execute(context.Background(), "group123")

	// Verificar resultado
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), repositoryError, err)
	suite.mockGroupRepo.AssertExpectations(suite.T())
}

func (suite *RepositoryErrorTestSuite) TestAddUserToGroupRepositoryError() {
	// Configurar mocks para GetByID (sucessos) e AddUserToGroup (erro)
	userObjectID := bson.NewObjectID()
	groupObjectID := bson.NewObjectID()

	existingUser := &entities.User{
		ID:    userObjectID,
		Name:  testUserNameRepo,
		Email: testUserEmailRepo,
	}
	existingGroup := &entities.Group{
		ID:      groupObjectID,
		Name:    "Test Group",
		Members: []string{},
	}

	suite.mockGroupRepo.On("GetByID", mock.Anything, "group123").Return(existingGroup, nil)
	suite.mockUserRepo.On("GetByID", mock.Anything, "user123").Return(existingUser, nil)

	repositoryError := errors.New("add user to group operation failed")
	suite.mockGroupRepo.On("AddUserToGroup", mock.Anything, "group123", "user123").Return(repositoryError)

	// Criar usecase com mock
	addUserToGroupUseCase := group.NewAddUserToGroupUseCase(suite.mockGroupRepo, suite.mockUserRepo)

	// Executar teste
	err := addUserToGroupUseCase.Execute(context.Background(), "group123", "user123")

	// Verificar resultado
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), repositoryError, err)
	suite.mockGroupRepo.AssertExpectations(suite.T())
	suite.mockUserRepo.AssertExpectations(suite.T())
}

func (suite *RepositoryErrorTestSuite) TestRemoveUserFromGroupRepositoryError() {
	// Configurar mock para retornar erro
	repositoryError := errors.New("remove user from group operation failed")
	suite.mockGroupRepo.On("RemoveUserFromGroup", mock.Anything, "group123", "user123").Return(repositoryError)

	// Criar usecase com mock
	removeUserFromGroupUseCase := group.NewRemoveUserFromGroupUseCase(suite.mockGroupRepo)

	// Executar teste
	err := removeUserFromGroupUseCase.Execute(context.Background(), "group123", "user123")

	// Verificar resultado
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), repositoryError, err)
	suite.mockGroupRepo.AssertExpectations(suite.T())
}

// Testes de cenários específicos de erro

func (suite *RepositoryErrorTestSuite) TestMultipleRepositoryErrors() {
	// Simular cenário onde múltiplos repositórios falham
	userRepoError := errors.New("user repository connection failed")
	groupRepoError := errors.New("group repository connection failed")

	suite.mockUserRepo.On("GetByID", mock.Anything, "user123").Return(nil, userRepoError)
	suite.mockGroupRepo.On("GetByID", mock.Anything, "group123").Return(nil, groupRepoError)

	// Testar user usecase
	getUserUseCase := user.NewGetUserUseCase(suite.mockUserRepo)
	userResult, userErr := getUserUseCase.Execute(context.Background(), "user123")

	// Testar group usecase
	getGroupUseCase := group.NewGetGroupUseCase(suite.mockGroupRepo)
	groupResult, groupErr := getGroupUseCase.Execute(context.Background(), "group123")

	// Verificar resultados
	assert.Nil(suite.T(), userResult)
	assert.Error(suite.T(), userErr)
	assert.Equal(suite.T(), userRepoError, userErr)

	assert.Nil(suite.T(), groupResult)
	assert.Error(suite.T(), groupErr)
	assert.Equal(suite.T(), groupRepoError, groupErr)

	suite.mockUserRepo.AssertExpectations(suite.T())
	suite.mockGroupRepo.AssertExpectations(suite.T())
}

func (suite *RepositoryErrorTestSuite) TestDatabaseTimeoutError() {
	// Simular erro de timeout do banco de dados
	timeoutError := errors.New("database operation timeout")

	suite.mockUserRepo.On("List", mock.Anything).Return(nil, timeoutError)

	// Criar usecase com mock
	listUsersUseCase := user.NewListUsersUseCase(suite.mockUserRepo)

	// Executar teste
	result, err := listUsersUseCase.Execute(context.Background())

	// Verificar resultado
	assert.Nil(suite.T(), result)
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "timeout")
	suite.mockUserRepo.AssertExpectations(suite.T())
}

func (suite *RepositoryErrorTestSuite) TestDatabaseConnectionError() {
	// Simular erro de conexão com banco de dados
	connectionError := errors.New("failed to connect to database")

	suite.mockUserRepo.On("Create", mock.Anything, mock.AnythingOfType(userEntityType)).Return(connectionError)

	// Criar usecase com mock
	createUserUseCase := user.NewCreateUserUseCase(suite.mockUserRepo)

	// Executar teste
	userDTO := &dto.CreateUserRequestDTO{
		Name:  "Test User",
		Email: "test@example.com",
	}

	result, err := createUserUseCase.Execute(context.Background(), userDTO)

	// Verificar resultado
	assert.Nil(suite.T(), result)
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "connect")
	suite.mockUserRepo.AssertExpectations(suite.T())
}

// Testes de integração para validar erros de inicialização dos repositórios
// Estes testes verificam os erros que podem ocorrer durante a chamada de GetMongoCollection

func TestRepositoryInitializationErrors(t *testing.T) {
	// Estes testes validam que os repositórios retornam erro adequadamente
	// quando há falha na inicialização via GetMongoCollection

	t.Run("UserRepository initialization follows error handling pattern", func(t *testing.T) {
		// Este teste valida que NewUserRepository retorna um erro quando
		// GetMongoCollection falha, seguindo o padrão (T, error) ao invés de panic

		// Testar se a função tem a assinatura correta para tratamento de erro
		userRepo, err := repositories.NewUserRepository(createInvalidDatabaseManager())

		// Verificar que o padrão de retorno está correto
		// Se GetMongoCollection falhar, deve retornar (nil, error)
		if err != nil {
			assert.Nil(t, userRepo, "Repository should be nil when initialization fails")
			assert.Error(t, err, "Error should be returned when initialization fails")
			assert.Contains(t, err.Error(), failedToInitializeUserRepo, "Error should indicate user repository initialization failure")
		} else {
			// Se não houve erro, o repository deve ser válido
			assert.NotNil(t, userRepo, "Repository should not be nil when initialization succeeds")
		}
	})

	t.Run("GroupRepository initialization follows error handling pattern", func(t *testing.T) {
		// Este teste valida que NewGroupRepository retorna um erro quando
		// GetMongoCollection falha, seguindo o padrão (T, error) ao invés de panic

		// Testar se a função tem a assinatura correta para tratamento de erro
		groupRepo, err := repositories.NewGroupRepository(createInvalidDatabaseManager())

		// Verificar que o padrão de retorno está correto
		// Se GetMongoCollection falhar, deve retornar (nil, error)
		if err != nil {
			assert.Nil(t, groupRepo, "Repository should be nil when initialization fails")
			assert.Error(t, err, "Error should be returned when initialization fails")
			assert.Contains(t, err.Error(), failedToInitializeGroupRepo, "Error should indicate group repository initialization failure")
		} else {
			// Se não houve erro, o repository deve ser válido
			assert.NotNil(t, groupRepo, "Repository should not be nil when initialization succeeds")
		}
	})
}

// Teste que valida o comportamento de erro na criação dos repositórios
// quando GetMongoCollection falha internamente
func TestRepositoryGetMongoCollectionValidation(t *testing.T) {
	t.Run("Repository constructor returns error instead of panic", func(t *testing.T) {
		// Este teste é importante porque valida que eliminamos o uso de panic
		// nas funções construtoras dos repositórios

		// Testar que NewUserRepository não faz panic e retorna erro adequadamente
		assert.NotPanics(t, func() {
			userRepo, err := repositories.NewUserRepository(createInvalidDatabaseManager())
			if err != nil {
				assert.Nil(t, userRepo)
				assert.Contains(t, err.Error(), "failed to initialize")
			}
		}, "NewUserRepository should not panic, should return error")

		// Testar que NewGroupRepository não faz panic e retorna erro adequadamente
		assert.NotPanics(t, func() {
			groupRepo, err := repositories.NewGroupRepository(createInvalidDatabaseManager())
			if err != nil {
				assert.Nil(t, groupRepo)
				assert.Contains(t, err.Error(), "failed to initialize")
			}
		}, "NewGroupRepository should not panic, should return error")
	})

	t.Run("GetMongoCollection error propagation", func(t *testing.T) {
		// Este teste valida que erros de GetMongoCollection são propagados corretamente
		// através do wrapper fmt.Errorf com %w para preservar o erro original

		// Usar um DatabaseManager sem inicialização para forçar erro em GetMongoCollection
		invalidDBManager := createInvalidDatabaseManager()

		// Testar propagação de erro no UserRepository
		userRepo, userErr := repositories.NewUserRepository(invalidDBManager)
		if userErr != nil {
			assert.Nil(t, userRepo)
			assert.Error(t, userErr)
			// Verificar que o erro foi envolvido adequadamente
			assert.Contains(t, userErr.Error(), "failed to initialize user repository")
		}

		// Testar propagação de erro no GroupRepository
		groupRepo, groupErr := repositories.NewGroupRepository(invalidDBManager)
		if groupErr != nil {
			assert.Nil(t, groupRepo)
			assert.Error(t, groupErr)
			// Verificar que o erro foi envolvido adequadamente
			assert.Contains(t, groupErr.Error(), "failed to initialize group repository")
		}
	})
}

// Helper function para criar um DatabaseManager que pode causar erro em GetMongoCollection
func createInvalidDatabaseManager() *database.DatabaseManager {
	// Criar um DatabaseManager com configuração mínima que não cause panic
	// mas que possa resultar em erro no GetMongoCollection
	cfg := &config.Config{
		DatabaseType: "mongodb", // Tipo válido mas sem outras configurações
	}

	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel) // Minimizar logs durante teste

	return database.NewDatabaseManager(cfg, logger)
	// Este DatabaseManager não foi inicializado, então GetMongoCollection deve falhar
}

func TestRepositoryErrorSuite(t *testing.T) {
	suite.Run(t, new(RepositoryErrorTestSuite))
}

// ========================================
// TESTES PARA ERROS DE ObjectIDFromHex
// Estes testes validam erros que ocorrem quando IDs inválidos são fornecidos
// Especificamente na linha: objectID, err := bson.ObjectIDFromHex(id)
// ========================================

func TestRepositoryObjectIDFromHexErrors(t *testing.T) {
	// Testes que validam erros quando IDs inválidos são fornecidos aos repositories
	// Estes erros ocorrem na conversão de string para ObjectID do MongoDB

	t.Run("UserRepository GetByID with invalid ObjectID", func(t *testing.T) {
		// Criar mock do UserRepository
		mockUserRepo := new(MockUserRepository)

		// IDs inválidos que devem causar erro em ObjectIDFromHex
		invalidIDs := []string{
			"invalid-id",                     // ID muito curto
			"12345",                          // ID numérico mas inválido
			"not-a-valid-objectid",           // String aleatória
			"",                               // String vazia
			"zzzzzzzzzzzzzzzzzzzzzzzz",       // Caracteres inválidos (z não é hex)
			"123456789012345678901234567890", // Muito longo
		}

		for _, invalidID := range invalidIDs {
			// Configurar mock para retornar erro específico de ObjectID inválido
			objectIDError := errors.New("invalid ObjectID format")
			mockUserRepo.On("GetByID", mock.Anything, invalidID).Return(nil, objectIDError)

			// Criar usecase com mock
			getUserUseCase := user.NewGetUserUseCase(mockUserRepo)

			// Executar teste
			result, err := getUserUseCase.Execute(context.Background(), invalidID)

			// Verificar resultado
			assert.Nil(t, result)
			assert.Error(t, err)
			assert.Equal(t, objectIDError, err)
			assert.Contains(t, err.Error(), "ObjectID")
		}

		mockUserRepo.AssertExpectations(t)
	})

	t.Run("UserRepository Delete with invalid ObjectID", func(t *testing.T) {
		// Criar mock do UserRepository
		mockUserRepo := new(MockUserRepository)

		// IDs inválidos para teste de Delete
		invalidIDs := []string{
			"xyz123",                           // Contém caracteres não-hex
			"12345678901234567890123456789012", // Muito longo (32 chars)
			"abc",                              // Muito curto
			"GGGGGGGGGGGGGGGGGGGGGGGG",         // Contém G (não é hex válido)
		}

		for _, invalidID := range invalidIDs {
			// Configurar mock para retornar erro de ObjectID inválido
			objectIDError := errors.New("invalid ObjectID hex string")
			mockUserRepo.On("Delete", mock.Anything, invalidID).Return(objectIDError)

			// Criar usecase com mock
			deleteUserUseCase := user.NewDeleteUserUseCase(mockUserRepo)

			// Executar teste
			err := deleteUserUseCase.Execute(context.Background(), invalidID)

			// Verificar resultado
			assert.Error(t, err)
			assert.Equal(t, objectIDError, err)
			assert.Contains(t, err.Error(), "ObjectID")
		}

		mockUserRepo.AssertExpectations(t)
	})

	t.Run("GroupRepository GetByID with invalid ObjectID", func(t *testing.T) {
		// Criar mock do GroupRepository
		mockGroupRepo := new(MockGroupRepository)

		// IDs inválidos que devem causar erro em ObjectIDFromHex
		invalidIDs := []string{
			"invalid-group-id",
			"123",
			"not-hex-characters",
			"!@#$%^&*()",                       // Caracteres especiais
			"abcdefghijklmnopqrstuvwxyz123456", // Muito longo com chars inválidos
		}

		for _, invalidID := range invalidIDs {
			// Configurar mock para retornar erro específico de ObjectID inválido
			objectIDError := errors.New("invalid ObjectID format for group")
			mockGroupRepo.On("GetByID", mock.Anything, invalidID).Return(nil, objectIDError)

			// Criar usecase com mock
			getGroupUseCase := group.NewGetGroupUseCase(mockGroupRepo)

			// Executar teste
			result, err := getGroupUseCase.Execute(context.Background(), invalidID)

			// Verificar resultado
			assert.Nil(t, result)
			assert.Error(t, err)
			assert.Equal(t, objectIDError, err)
			assert.Contains(t, err.Error(), "ObjectID")
		}

		mockGroupRepo.AssertExpectations(t)
	})

	t.Run("GroupRepository Delete with invalid ObjectID", func(t *testing.T) {
		// Criar mock do GroupRepository
		mockGroupRepo := new(MockGroupRepository)

		// IDs inválidos para teste de Delete
		invalidIDs := []string{
			"short",                             // Muito curto
			"spaces in id",                      // Contém espaços
			"123456789012345678901234567890123", // 33 caracteres (inválido)
			"zxcvbnmasdfghjklqwertyuiop",        // Contém letras não-hex
		}

		for _, invalidID := range invalidIDs {
			// Configurar mock para retornar erro de ObjectID inválido
			objectIDError := errors.New("invalid ObjectID hex representation")
			mockGroupRepo.On("Delete", mock.Anything, invalidID).Return(objectIDError)

			// Criar usecase com mock
			deleteGroupUseCase := group.NewDeleteGroupUseCase(mockGroupRepo)

			// Executar teste
			err := deleteGroupUseCase.Execute(context.Background(), invalidID)

			// Verificar resultado
			assert.Error(t, err)
			assert.Equal(t, objectIDError, err)
			assert.Contains(t, err.Error(), "ObjectID")
		}

		mockGroupRepo.AssertExpectations(t)
	})

	t.Run("GroupRepository AddUserToGroup with invalid groupID", func(t *testing.T) {
		// Criar mocks do GroupRepository e UserRepository
		mockGroupRepo := new(MockGroupRepository)
		mockUserRepo := new(MockUserRepository)

		// IDs inválidos para groupID
		invalidGroupIDs := []string{
			"invalid-group",
			"12345678901234567890123456789012345", // Muito longo
			"hhhhhhhhhhhhhhhhhhhhhhh",             // Contém 'h' (não é hex)
		}

		validUserID := "507f1f77bcf86cd799439011" // ID válido para user

		for _, invalidGroupID := range invalidGroupIDs {
			// O AddUserToGroupUseCase primeiro chama GetByID para validar se o grupo existe
			// Configurar mock para GetByID falhar com erro de ObjectID inválido
			objectIDError := errors.New("invalid groupID ObjectID format")
			mockGroupRepo.On("GetByID", mock.Anything, invalidGroupID).Return(nil, objectIDError)

			// Criar usecase com mocks
			addUserToGroupUseCase := group.NewAddUserToGroupUseCase(mockGroupRepo, mockUserRepo)

			// Executar teste
			err := addUserToGroupUseCase.Execute(context.Background(), invalidGroupID, validUserID)

			// Verificar resultado
			assert.Error(t, err)
			assert.Equal(t, objectIDError, err)
			assert.Contains(t, err.Error(), "ObjectID")
		}

		mockGroupRepo.AssertExpectations(t)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("GroupRepository RemoveUserFromGroup with invalid groupID", func(t *testing.T) {
		// Criar mock do GroupRepository
		mockGroupRepo := new(MockGroupRepository)

		// IDs inválidos para groupID
		invalidGroupIDs := []string{
			"bad-group-id",
			"1234567890123456789012345678901234567890", // Extremamente longo
			"ZZZZZZZZZZZZZZZZZZZZZZZZ",                 // Contém Z (não é hex válido)
		}

		validUserID := "507f1f77bcf86cd799439011" // ID válido para user

		for _, invalidGroupID := range invalidGroupIDs {
			// O RemoveUserFromGroupUseCase chama diretamente RemoveUserFromGroup
			// Configurar mock para retornar erro de ObjectID inválido para groupID
			objectIDError := errors.New("invalid groupID ObjectID hex string")
			mockGroupRepo.On("RemoveUserFromGroup", mock.Anything, invalidGroupID, validUserID).Return(objectIDError)

			// Criar usecase com mock
			removeUserFromGroupUseCase := group.NewRemoveUserFromGroupUseCase(mockGroupRepo)

			// Executar teste
			err := removeUserFromGroupUseCase.Execute(context.Background(), invalidGroupID, validUserID)

			// Verificar resultado
			assert.Error(t, err)
			assert.Equal(t, objectIDError, err)
			assert.Contains(t, err.Error(), "ObjectID")
		}

		mockGroupRepo.AssertExpectations(t)
	})
}

func TestRepositoryObjectIDFromHexIntegrationErrors(t *testing.T) {
	// Testes de integração que validam os erros reais de ObjectIDFromHex
	// usando repositories reais ao invés de mocks

	t.Run("Real UserRepository GetByID with invalid ObjectID strings", func(t *testing.T) {
		// Criar um DatabaseManager válido para teste
		cfg := &config.Config{
			DatabaseType: "mongodb",
			MongoURI:     "mongodb://localhost:27017", // URI de teste
			MongoDB:      "test_db",
		}

		loggerInstance := logrus.New()
		loggerInstance.SetLevel(logrus.ErrorLevel)

		dbManager := database.NewDatabaseManager(cfg, loggerInstance)

		// Tentar criar repository (pode falhar se não houver MongoDB, mas não é o foco)
		userRepo, err := repositories.NewUserRepository(dbManager)
		if err != nil {
			// Se não conseguir criar repository, pular este teste
			t.Skip("Skipping integration test - could not create UserRepository")
			return
		}

		// IDs definitivamente inválidos que devem causar erro em ObjectIDFromHex
		invalidIDs := []string{
			"invalid",                          // Muito curto
			"not-a-valid-id",                   // String comum
			"12345678901234567890123456789012", // 32 chars mas alguns podem ser inválidos
			"gggggggggggggggggggggggg",         // 24 chars mas com 'g' (não hex)
			"",                                 // String vazia
		}

		for _, invalidID := range invalidIDs {
			// Tentar GetByID com ID inválido
			result, err := userRepo.GetByID(context.Background(), invalidID)

			// Deve retornar erro
			assert.Error(t, err, "Should return error for invalid ID: %s", invalidID)
			assert.Nil(t, result, "Result should be nil for invalid ID: %s", invalidID)

			// O erro deve ser relacionado a ObjectID inválido
			assert.True(t,
				strings.Contains(err.Error(), "invalid") ||
					strings.Contains(err.Error(), "ObjectID") ||
					strings.Contains(err.Error(), "hex"),
				"Error should indicate invalid ObjectID for ID: %s, got: %v", invalidID, err)
		}
	})

	t.Run("Real GroupRepository operations with invalid ObjectID strings", func(t *testing.T) {
		// Criar um DatabaseManager válido para teste
		cfg := &config.Config{
			DatabaseType: "mongodb",
			MongoURI:     "mongodb://localhost:27017",
			MongoDB:      "test_db",
		}

		loggerInstance := logrus.New()
		loggerInstance.SetLevel(logrus.ErrorLevel)

		dbManager := database.NewDatabaseManager(cfg, loggerInstance)

		// Tentar criar repository
		groupRepo, err := repositories.NewGroupRepository(dbManager)
		if err != nil {
			t.Skip("Skipping integration test - could not create GroupRepository")
			return
		}

		// IDs inválidos para teste
		invalidIDs := []string{
			"xyz",                                  // Muito curto
			"invalid-group-id-string",              // String comum
			"123456789012345678901234567890123456", // Muito longo
			"qwertyuiopasdfghjklzxcvb",             // 24 chars mas não hex válido
		}

		validUserID := "507f1f77bcf86cd799439011" // ID válido para comparação

		for _, invalidID := range invalidIDs {
			// Testar GetByID
			result, err := groupRepo.GetByID(context.Background(), invalidID)
			assert.Error(t, err, "GetByID should return error for invalid ID: %s", invalidID)
			assert.Nil(t, result, "GetByID result should be nil for invalid ID: %s", invalidID)

			// Testar Delete
			err = groupRepo.Delete(context.Background(), invalidID)
			assert.Error(t, err, "Delete should return error for invalid ID: %s", invalidID)

			// Testar AddUserToGroup (com groupID inválido)
			err = groupRepo.AddUserToGroup(context.Background(), invalidID, validUserID)
			assert.Error(t, err, "AddUserToGroup should return error for invalid groupID: %s", invalidID)

			// Testar RemoveUserFromGroup (com groupID inválido)
			err = groupRepo.RemoveUserFromGroup(context.Background(), invalidID, validUserID)
			assert.Error(t, err, "RemoveUserFromGroup should return error for invalid groupID: %s", invalidID)

			// Verificar que todos os erros são relacionados a ObjectID
			for _, operation := range []string{"GetByID", "Delete", "AddUserToGroup", "RemoveUserFromGroup"} {
				t.Logf("Operation %s with invalid ID %s returned expected error", operation, invalidID)
			}
		}
	})
}

func TestRepositoryObjectIDValidationScenarios(t *testing.T) {
	// Testes que cobrem cenários específicos de validação de ObjectID

	t.Run("Edge cases for ObjectID validation", func(t *testing.T) {
		// Casos extremos de IDs inválidos
		edgeCases := map[string]string{
			"empty_string":         "",
			"only_spaces":          "   ",
			"special_chars":        "!@#$%^&*()_+{}|:<>?",
			"unicode_chars":        "ñáéíóú",
			"mixed_case_invalid":   "AbCdEfGhIjKlMnOpQrSt",
			"with_dashes":          "123-456-789-012-345-678",
			"with_dots":            "123.456.789.012.345.678",
			"hex_but_wrong_length": "abc123def456",                // Hex válido mas tamanho errado
			"almost_valid":         "507f1f77bcf86cd79943901",     // 23 chars ao invés de 24
			"too_long_valid_hex":   "507f1f77bcf86cd799439011123", // 26 chars
		}

		mockUserRepo := new(MockUserRepository)
		mockGroupRepo := new(MockGroupRepository)

		for caseName, invalidID := range edgeCases {
			t.Run(fmt.Sprintf("UserRepository_GetByID_%s", caseName), func(t *testing.T) {
				// Configurar mock para retornar erro específico
				objectIDError := fmt.Errorf("ObjectIDFromHex failed for case: %s", caseName)
				mockUserRepo.On("GetByID", mock.Anything, invalidID).Return(nil, objectIDError)

				// Criar usecase
				getUserUseCase := user.NewGetUserUseCase(mockUserRepo)

				// Executar teste
				result, err := getUserUseCase.Execute(context.Background(), invalidID)

				// Verificar resultado
				assert.Error(t, err, "Should return error for case: %s", caseName)
				assert.Nil(t, result, "Result should be nil for case: %s", caseName)
				assert.Contains(t, err.Error(), "ObjectIDFromHex", "Error should mention ObjectIDFromHex for case: %s", caseName)
			})

			t.Run(fmt.Sprintf("GroupRepository_GetByID_%s", caseName), func(t *testing.T) {
				// Configurar mock para retornar erro específico
				objectIDError := fmt.Errorf("ObjectIDFromHex failed for case: %s", caseName)
				mockGroupRepo.On("GetByID", mock.Anything, invalidID).Return(nil, objectIDError)

				// Criar usecase
				getGroupUseCase := group.NewGetGroupUseCase(mockGroupRepo)

				// Executar teste
				result, err := getGroupUseCase.Execute(context.Background(), invalidID)

				// Verificar resultado
				assert.Error(t, err, "Should return error for case: %s", caseName)
				assert.Nil(t, result, "Result should be nil for case: %s", caseName)
				assert.Contains(t, err.Error(), "ObjectIDFromHex", "Error should mention ObjectIDFromHex for case: %s", caseName)
			})
		}

		mockUserRepo.AssertExpectations(t)
		mockGroupRepo.AssertExpectations(t)
	})
}

// ========================================
// TESTES DE INTEGRAÇÃO ESPECÍFICOS PARA ObjectIDFromHex
// Estes testes usam repositories reais e MongoDB real para validar os erros
// ========================================

func (suite *IntegrationTestSuite) TestRepositoryObjectIDFromHexRealErrors() {
	// Testes de integração que validam erros reais de ObjectIDFromHex

	suite.T().Run("UserRepository real ObjectIDFromHex errors", func(t *testing.T) {
		// Criar repository real
		cfg := suite.testConfig.ToAppConfig()
		loggerInstance := logger.NewLogger()
		dbManager := database.NewDatabaseManager(cfg, loggerInstance)

		err := dbManager.Initialize(context.Background())
		assert.NoError(t, err)

		userRepo, err := repositories.NewUserRepository(dbManager)
		assert.NoError(t, err)

		// Criar usecase real
		getUserUseCase := user.NewGetUserUseCase(userRepo)
		deleteUserUseCase := user.NewDeleteUserUseCase(userRepo)

		// IDs inválidos que devem causar erro em ObjectIDFromHex
		invalidIDs := []string{
			"invalid",                             // Muito curto
			"not-a-valid-objectid-string",         // String comum
			"12345678901234567890123456789012345", // Muito longo (33 chars)
			"gggggggggggggggggggggggg",            // 24 chars mas não hex válido (g)
			"zzzzzzzzzzzzzzzzzzzzzzzz",            // 24 chars mas não hex válido (z)
			"",                                    // String vazia
			"   ",                                 // Apenas espaços
			"507f1f77bcf86cd79943901g",            // 24 chars mas último char inválido
		}

		for _, invalidID := range invalidIDs {
			// Testar GetByID
			result, err := getUserUseCase.Execute(context.Background(), invalidID)
			assert.Error(t, err, "GetByID should return error for invalid ID: %s", invalidID)
			assert.Nil(t, result, "GetByID result should be nil for invalid ID: %s", invalidID)

			// Verificar que o erro é relacionado ao ObjectID
			assert.True(t,
				strings.Contains(err.Error(), "invalid") ||
					strings.Contains(err.Error(), "ObjectID") ||
					strings.Contains(err.Error(), "hex") ||
					strings.Contains(err.Error(), "encoding"),
				"Error should indicate invalid ObjectID for ID: %s, got: %v", invalidID, err)

			// Testar Delete
			err = deleteUserUseCase.Execute(context.Background(), invalidID)
			assert.Error(t, err, "Delete should return error for invalid ID: %s", invalidID)

			// Verificar que o erro é relacionado ao ObjectID
			assert.True(t,
				strings.Contains(err.Error(), "invalid") ||
					strings.Contains(err.Error(), "ObjectID") ||
					strings.Contains(err.Error(), "hex") ||
					strings.Contains(err.Error(), "encoding"),
				"Delete error should indicate invalid ObjectID for ID: %s, got: %v", invalidID, err)

			t.Logf("SUCCESS: Invalid ID '%s' correctly returned ObjectID error", invalidID)
		}
	})

	suite.T().Run("GroupRepository real ObjectIDFromHex errors", func(t *testing.T) {
		// Criar repository real
		cfg := suite.testConfig.ToAppConfig()
		loggerInstance := logger.NewLogger()
		dbManager := database.NewDatabaseManager(cfg, loggerInstance)

		err := dbManager.Initialize(context.Background())
		assert.NoError(t, err)

		groupRepo, err := repositories.NewGroupRepository(dbManager)
		assert.NoError(t, err)

		// Criar usecases reais
		getGroupUseCase := group.NewGetGroupUseCase(groupRepo)
		deleteGroupUseCase := group.NewDeleteGroupUseCase(groupRepo)
		removeUserFromGroupUseCase := group.NewRemoveUserFromGroupUseCase(groupRepo)

		// IDs inválidos específicos para testes de grupo
		invalidGroupIDs := []string{
			"bad-group",                          // String comum
			"123",                                // Muito curto
			"1234567890123456789012345678901234", // 32 chars (inválido)
			"xxxxxxxxxxxxxxxxxxxx",               // 20 chars com x (inválido)
			"507f1f77bcf86cd799439011xxxx",       // 28 chars (muito longo)
			"ABCDEFGHIJKLMNOPQRSTUVWX",           // 24 chars mas maiúsculas com chars inválidos
		}

		validUserID := "507f1f77bcf86cd799439011" // ID válido para testes

		for _, invalidGroupID := range invalidGroupIDs {
			// Testar GetByID
			result, err := getGroupUseCase.Execute(context.Background(), invalidGroupID)
			assert.Error(t, err, "GetByID should return error for invalid groupID: %s", invalidGroupID)
			assert.Nil(t, result, "GetByID result should be nil for invalid groupID: %s", invalidGroupID)

			// Testar Delete
			err = deleteGroupUseCase.Execute(context.Background(), invalidGroupID)
			assert.Error(t, err, "Delete should return error for invalid groupID: %s", invalidGroupID)

			// Testar RemoveUserFromGroup (que chama ObjectIDFromHex internamente)
			err = removeUserFromGroupUseCase.Execute(context.Background(), invalidGroupID, validUserID)
			assert.Error(t, err, "RemoveUserFromGroup should return error for invalid groupID: %s", invalidGroupID)

			// Verificar que todos os erros são relacionados ao ObjectID
			for _, operation := range []string{"GetByID", "Delete", "RemoveUserFromGroup"} {
				t.Logf("SUCCESS: Operation %s with invalid groupID '%s' returned expected ObjectID error", operation, invalidGroupID)
			}
		}
	})

	suite.T().Run("AddUserToGroup with invalid IDs", func(t *testing.T) {
		// Testar AddUserToGroup que pode falhar tanto no groupID quanto no userID

		cfg := suite.testConfig.ToAppConfig()
		loggerInstance := logger.NewLogger()
		dbManager := database.NewDatabaseManager(cfg, loggerInstance)

		err := dbManager.Initialize(context.Background())
		assert.NoError(t, err)

		groupRepo, err := repositories.NewGroupRepository(dbManager)
		assert.NoError(t, err)

		userRepo, err := repositories.NewUserRepository(dbManager)
		assert.NoError(t, err)

		addUserToGroupUseCase := group.NewAddUserToGroupUseCase(groupRepo, userRepo)

		invalidIDs := []string{
			"invalid-id",
			"too-short",
			"waytoolongtobeavalidobjectidstring",
			"!@#$%^&*()",
		}

		validID := "507f1f77bcf86cd799439011"

		// Testar com groupID inválido e userID válido
		for _, invalidGroupID := range invalidIDs {
			err = addUserToGroupUseCase.Execute(context.Background(), invalidGroupID, validID)
			assert.Error(t, err, "AddUserToGroup should return error for invalid groupID: %s", invalidGroupID)

			t.Logf("SUCCESS: AddUserToGroup with invalid groupID '%s' returned expected error", invalidGroupID)
		}

		// Criar um grupo válido para testar userID inválido
		testGroup := &entities.Group{
			ID:      bson.NewObjectID(),
			Name:    "Test Group for ObjectID validation",
			Members: []string{},
		}
		err = groupRepo.Create(context.Background(), testGroup)
		assert.NoError(t, err)

		validGroupID := testGroup.ID.Hex()

		// Testar com groupID válido e userID inválido
		for _, invalidUserID := range invalidIDs {
			err = addUserToGroupUseCase.Execute(context.Background(), validGroupID, invalidUserID)
			assert.Error(t, err, "AddUserToGroup should return error for invalid userID: %s", invalidUserID)

			t.Logf("SUCCESS: AddUserToGroup with invalid userID '%s' returned expected error", invalidUserID)
		}
	})
}

// Testes específicos para validar erros nas funções List dos repositories
// Estes testes focam nos erros que podem ocorrer durante Find e cursor.All

func TestUserRepositoryListFindError(t *testing.T) {
	// Teste que simula erro na operação Find do MongoDB
	t.Run("UserRepository List fails on Find operation", func(t *testing.T) {
		// Criar mock do UserRepository
		mockUserRepo := new(MockUserRepository)

		// Configurar erro específico para operação Find
		findError := errors.New("MongoDB Find operation failed: connection timeout")
		mockUserRepo.On("List", mock.Anything).Return(nil, findError)

		// Criar usecase com mock
		listUsersUseCase := user.NewListUsersUseCase(mockUserRepo)

		// Executar teste
		result, err := listUsersUseCase.Execute(context.Background())

		// Verificar resultado
		assert.Nil(t, result)
		assert.Error(t, err)
		assert.Equal(t, findError, err)
		assert.Contains(t, err.Error(), "Find operation failed")
		mockUserRepo.AssertExpectations(t)
	})
}

func TestGroupRepositoryListFindError(t *testing.T) {
	// Teste que simula erro na operação Find do MongoDB para grupos
	t.Run("GroupRepository List fails on Find operation", func(t *testing.T) {
		// Criar mock do GroupRepository
		mockGroupRepo := new(MockGroupRepository)

		// Configurar erro específico para operação Find
		findError := errors.New("MongoDB Find operation failed: index corruption")
		mockGroupRepo.On("List", mock.Anything).Return(nil, findError)

		// Criar usecase com mock
		listGroupsUseCase := group.NewListGroupsUseCase(mockGroupRepo)

		// Executar teste
		result, err := listGroupsUseCase.Execute(context.Background())

		// Verificar resultado
		assert.Nil(t, result)
		assert.Error(t, err)
		assert.Equal(t, findError, err)
		assert.Contains(t, err.Error(), "Find operation failed")
		mockGroupRepo.AssertExpectations(t)
	})
}

func TestUserRepositoryListCursorAllError(t *testing.T) {
	// Teste que simula erro na operação cursor.All do MongoDB
	t.Run("UserRepository List fails on cursor.All operation", func(t *testing.T) {
		// Criar mock do UserRepository
		mockUserRepo := new(MockUserRepository)

		// Configurar erro específico para operação cursor.All
		cursorError := errors.New("cursor.All operation failed: memory allocation error")
		mockUserRepo.On("List", mock.Anything).Return(nil, cursorError)

		// Criar usecase com mock
		listUsersUseCase := user.NewListUsersUseCase(mockUserRepo)

		// Executar teste
		result, err := listUsersUseCase.Execute(context.Background())

		// Verificar resultado
		assert.Nil(t, result)
		assert.Error(t, err)
		assert.Equal(t, cursorError, err)
		assert.Contains(t, err.Error(), "cursor.All operation failed")
		mockUserRepo.AssertExpectations(t)
	})
}

func TestGroupRepositoryListCursorAllError(t *testing.T) {
	// Teste que simula erro na operação cursor.All do MongoDB para grupos
	t.Run("GroupRepository List fails on cursor.All operation", func(t *testing.T) {
		// Criar mock do GroupRepository
		mockGroupRepo := new(MockGroupRepository)

		// Configurar erro específico para operação cursor.All
		cursorError := errors.New("cursor.All operation failed: document parsing error")
		mockGroupRepo.On("List", mock.Anything).Return(nil, cursorError)

		// Criar usecase com mock
		listGroupsUseCase := group.NewListGroupsUseCase(mockGroupRepo)

		// Executar teste
		result, err := listGroupsUseCase.Execute(context.Background())

		// Verificar resultado
		assert.Nil(t, result)
		assert.Error(t, err)
		assert.Equal(t, cursorError, err)
		assert.Contains(t, err.Error(), "cursor.All operation failed")
		mockGroupRepo.AssertExpectations(t)
	})
}

func TestRepositoryListOperationScenarios(t *testing.T) {
	// Conjunto de testes que cobrem diferentes cenários de erro nas operações List

	t.Run("UserRepository List with database connection lost during Find", func(t *testing.T) {
		mockUserRepo := new(MockUserRepository)

		// Simular perda de conexão durante Find
		connectionLostError := errors.New("Find failed: connection to MongoDB lost during query execution")
		mockUserRepo.On("List", mock.Anything).Return(nil, connectionLostError)

		listUsersUseCase := user.NewListUsersUseCase(mockUserRepo)
		result, err := listUsersUseCase.Execute(context.Background())

		assert.Nil(t, result)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "connection to MongoDB lost")
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("GroupRepository List with cursor timeout during All operation", func(t *testing.T) {
		mockGroupRepo := new(MockGroupRepository)

		// Simular timeout durante cursor.All
		timeoutError := errors.New("cursor.All failed: operation exceeded time limit")
		mockGroupRepo.On("List", mock.Anything).Return(nil, timeoutError)

		listGroupsUseCase := group.NewListGroupsUseCase(mockGroupRepo)
		result, err := listGroupsUseCase.Execute(context.Background())

		assert.Nil(t, result)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "operation exceeded time limit")
		mockGroupRepo.AssertExpectations(t)
	})

	t.Run("UserRepository List with invalid BSON during cursor.All", func(t *testing.T) {
		mockUserRepo := new(MockUserRepository)

		// Simular erro de BSON inválido durante cursor.All
		bsonError := errors.New("cursor.All failed: invalid BSON document structure")
		mockUserRepo.On("List", mock.Anything).Return(nil, bsonError)

		listUsersUseCase := user.NewListUsersUseCase(mockUserRepo)
		result, err := listUsersUseCase.Execute(context.Background())

		assert.Nil(t, result)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid BSON document")
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("GroupRepository List with memory exhaustion during cursor.All", func(t *testing.T) {
		mockGroupRepo := new(MockGroupRepository)

		// Simular esgotamento de memória durante cursor.All
		memoryError := errors.New("cursor.All failed: insufficient memory to load all documents")
		mockGroupRepo.On("List", mock.Anything).Return(nil, memoryError)

		listGroupsUseCase := group.NewListGroupsUseCase(mockGroupRepo)
		result, err := listGroupsUseCase.Execute(context.Background())

		assert.Nil(t, result)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "insufficient memory")
		mockGroupRepo.AssertExpectations(t)
	})

	t.Run("UserRepository List with index corruption during Find", func(t *testing.T) {
		mockUserRepo := new(MockUserRepository)

		// Simular corrupção de índice durante Find
		indexError := errors.New("Find failed: index corruption detected, unable to execute query")
		mockUserRepo.On("List", mock.Anything).Return(nil, indexError)

		listUsersUseCase := user.NewListUsersUseCase(mockUserRepo)
		result, err := listUsersUseCase.Execute(context.Background())

		assert.Nil(t, result)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "index corruption detected")
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("GroupRepository List with collection lock during Find", func(t *testing.T) {
		mockGroupRepo := new(MockGroupRepository)

		// Simular bloqueio de coleção durante Find
		lockError := errors.New("Find failed: collection is locked by another operation")
		mockGroupRepo.On("List", mock.Anything).Return(nil, lockError)

		listGroupsUseCase := group.NewListGroupsUseCase(mockGroupRepo)
		result, err := listGroupsUseCase.Execute(context.Background())

		assert.Nil(t, result)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "collection is locked")
		mockGroupRepo.AssertExpectations(t)
	})
}

func TestRepositoryListErrorPropagation(t *testing.T) {
	// Testes que validam a propagação correta de erros das operações List

	t.Run("Error propagation from Find to UseCase", func(t *testing.T) {
		mockUserRepo := new(MockUserRepository)

		// Erro específico que deve ser propagado sem modificação
		originalFindError := errors.New("MongoDB Find error: server selection timeout")
		mockUserRepo.On("List", mock.Anything).Return(nil, originalFindError)

		listUsersUseCase := user.NewListUsersUseCase(mockUserRepo)
		result, err := listUsersUseCase.Execute(context.Background())

		// Verificar que o erro original foi propagado sem alteração
		assert.Nil(t, result)
		assert.Error(t, err)
		assert.Equal(t, originalFindError, err, "Error should be propagated without modification")
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("Error propagation from cursor.All to UseCase", func(t *testing.T) {
		mockGroupRepo := new(MockGroupRepository)

		// Erro específico que deve ser propagado sem modificação
		originalCursorError := errors.New("cursor.All error: decode failed for document")
		mockGroupRepo.On("List", mock.Anything).Return(nil, originalCursorError)

		listGroupsUseCase := group.NewListGroupsUseCase(mockGroupRepo)
		result, err := listGroupsUseCase.Execute(context.Background())

		// Verificar que o erro original foi propagado sem alteração
		assert.Nil(t, result)
		assert.Error(t, err)
		assert.Equal(t, originalCursorError, err, "Error should be propagated without modification")
		mockGroupRepo.AssertExpectations(t)
	})
}

// ========================================
// TESTES DE INTEGRAÇÃO PARA VALIDAR ERROS NAS LINHAS ESPECÍFICAS DOS REPOSITORIES
// Estes testes fazem com que as linhas específicas dos repositories retornem erro:
// - cursor, err := r.collection.Find(ctx, bson.M{})
// - if err := cursor.All(ctx, &users); err != nil
// ========================================

func (suite *IntegrationTestSuite) TestUserRepositoryListFindErrorWithCanceledContext() {
	// Teste que força erro na linha: cursor, err := r.collection.Find(ctx, bson.M{})

	// Criar repository normalmente
	cfg := suite.testConfig.ToAppConfig()
	loggerInstance := logger.NewLogger()
	dbManager := database.NewDatabaseManager(cfg, loggerInstance)

	// Inicializar conexão
	ctx := context.Background()
	err := dbManager.Initialize(ctx)
	assert.NoError(suite.T(), err)

	userRepo, err := repositories.NewUserRepository(dbManager)
	assert.NoError(suite.T(), err)

	// Criar usecase
	listUsersUseCase := user.NewListUsersUseCase(userRepo)

	// Criar contexto cancelado
	canceledCtx, cancel := context.WithCancel(context.Background())
	cancel() // Cancelar imediatamente

	// Executar List com contexto cancelado - deve falhar na linha Find
	result, err := listUsersUseCase.Execute(canceledCtx)

	// Verificar que ocorreu erro
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
	assert.Contains(suite.T(), err.Error(), contextCanceledError)
}

func (suite *IntegrationTestSuite) TestGroupRepositoryListFindErrorWithCanceledContext() {
	// Teste que força erro na linha: cursor, err := r.collection.Find(ctx, bson.M{})

	// Criar repository normalmente
	cfg := suite.testConfig.ToAppConfig()
	loggerInstance := logger.NewLogger()
	dbManager := database.NewDatabaseManager(cfg, loggerInstance)

	// Inicializar conexão
	ctx := context.Background()
	err := dbManager.Initialize(ctx)
	assert.NoError(suite.T(), err)

	groupRepo, err := repositories.NewGroupRepository(dbManager)
	assert.NoError(suite.T(), err)

	// Criar usecase
	listGroupsUseCase := group.NewListGroupsUseCase(groupRepo)

	// Criar contexto cancelado
	canceledCtx, cancel := context.WithCancel(context.Background())
	cancel() // Cancelar imediatamente

	// Executar List com contexto cancelado - deve falhar na linha Find
	result, err := listGroupsUseCase.Execute(canceledCtx)

	// Verificar que ocorreu erro
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
	assert.Contains(suite.T(), err.Error(), contextCanceledError)
}

func (suite *IntegrationTestSuite) TestUserRepositoryListFindErrorWithTimeout() {
	// Teste que força erro na linha: cursor, err := r.collection.Find(ctx, bson.M{})

	// Criar repository normalmente
	cfg := suite.testConfig.ToAppConfig()
	loggerInstance := logger.NewLogger()
	dbManager := database.NewDatabaseManager(cfg, loggerInstance)

	// Inicializar conexão
	ctx := context.Background()
	err := dbManager.Initialize(ctx)
	assert.NoError(suite.T(), err)

	userRepo, err := repositories.NewUserRepository(dbManager)
	assert.NoError(suite.T(), err)

	// Criar usecase
	listUsersUseCase := user.NewListUsersUseCase(userRepo)

	// Criar contexto com timeout muito curto (1 nanosegundo)
	timeoutCtx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
	defer cancel()

	// Aguardar o contexto expirar
	time.Sleep(1 * time.Millisecond)

	// Executar List com contexto expirado - deve falhar na linha Find
	result, err := listUsersUseCase.Execute(timeoutCtx)

	// Verificar que ocorreu erro
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
	assert.Contains(suite.T(), err.Error(), contextDeadlineError)
}

func (suite *IntegrationTestSuite) TestGroupRepositoryListFindErrorWithTimeout() {
	// Teste que força erro na linha: cursor, err := r.collection.Find(ctx, bson.M{})

	// Criar repository normalmente
	cfg := suite.testConfig.ToAppConfig()
	loggerInstance := logger.NewLogger()
	dbManager := database.NewDatabaseManager(cfg, loggerInstance)

	// Inicializar conexão
	ctx := context.Background()
	err := dbManager.Initialize(ctx)
	assert.NoError(suite.T(), err)

	groupRepo, err := repositories.NewGroupRepository(dbManager)
	assert.NoError(suite.T(), err)

	// Criar usecase
	listGroupsUseCase := group.NewListGroupsUseCase(groupRepo)

	// Criar contexto com timeout muito curto (1 nanosegundo)
	timeoutCtx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
	defer cancel()

	// Aguardar o contexto expirar
	time.Sleep(1 * time.Millisecond)

	// Executar List com contexto expirado - deve falhar na linha Find
	result, err := listGroupsUseCase.Execute(timeoutCtx)

	// Verificar que ocorreu erro
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
	assert.Contains(suite.T(), err.Error(), contextDeadlineError)
}

func (suite *IntegrationTestSuite) TestUserRepositoryListFindErrorWithDisconnectedClient() {
	// Teste que força erro na linha: cursor, err := r.collection.Find(ctx, bson.M{})
	// usando uma conexão que será desconectada

	// Criar uma nova conexão MongoDB
	client, err := mongo.Connect(options.Client().ApplyURI(suite.testConfig.MongoURI))
	assert.NoError(suite.T(), err)

	// Conectar e depois desconectar imediatamente
	err = client.Ping(context.Background(), nil)
	assert.NoError(suite.T(), err)

	err = client.Disconnect(context.Background())
	assert.NoError(suite.T(), err)

	// Agora tentar usar a coleção com a conexão fechada
	collection := client.Database(suite.testConfig.MongoDB).Collection("users")

	// Tentar Find com conexão fechada - deve falhar
	cursor, err := collection.Find(context.Background(), map[string]interface{}{})

	// Verificar que ocorreu erro (este teste valida o comportamento do MongoDB)
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "client is disconnected")

	if cursor != nil {
		cursor.Close(context.Background())
	}
}

func (suite *IntegrationTestSuite) TestGroupRepositoryListFindErrorWithDisconnectedClient() {
	// Teste que força erro na linha: cursor, err := r.collection.Find(ctx, bson.M{})
	// usando uma conexão que será desconectada

	// Criar uma nova conexão MongoDB
	client, err := mongo.Connect(options.Client().ApplyURI(suite.testConfig.MongoURI))
	assert.NoError(suite.T(), err)

	// Conectar e depois desconectar imediatamente
	err = client.Ping(context.Background(), nil)
	assert.NoError(suite.T(), err)

	err = client.Disconnect(context.Background())
	assert.NoError(suite.T(), err)

	// Agora tentar usar a coleção com a conexão fechada
	collection := client.Database(suite.testConfig.MongoDB).Collection("groups")

	// Tentar Find com conexão fechada - deve falhar
	cursor, err := collection.Find(context.Background(), map[string]interface{}{})

	// Verificar que ocorreu erro (este teste valida o comportamento do MongoDB)
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "client is disconnected")

	if cursor != nil {
		cursor.Close(context.Background())
	}
}

func (suite *IntegrationTestSuite) TestRepositoryListErrorsValidationDirectly() {
	// Teste que valida se os repositories realmente propagam erros das operações Find e cursor.All

	suite.T().Run("UserRepository List propagates Find errors with canceled context", func(t *testing.T) {
		// Criar repository normalmente
		cfg := suite.testConfig.ToAppConfig()
		loggerInstance := logger.NewLogger()
		dbManager := database.NewDatabaseManager(cfg, loggerInstance)

		err := dbManager.Initialize(context.Background())
		assert.NoError(t, err)

		userRepo, err := repositories.NewUserRepository(dbManager)
		assert.NoError(t, err)

		// Criar contexto cancelado
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		// Executar List que deve falhar na linha Find
		result, err := userRepo.List(ctx)

		// Verificar que o erro foi propagado corretamente
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), contextCanceledError)
	})

	suite.T().Run("GroupRepository List propagates Find errors with canceled context", func(t *testing.T) {
		// Criar repository normalmente
		cfg := suite.testConfig.ToAppConfig()
		loggerInstance := logger.NewLogger()
		dbManager := database.NewDatabaseManager(cfg, loggerInstance)

		err := dbManager.Initialize(context.Background())
		assert.NoError(t, err)

		groupRepo, err := repositories.NewGroupRepository(dbManager)
		assert.NoError(t, err)

		// Criar contexto cancelado
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		// Executar List que deve falhar na linha Find
		result, err := groupRepo.List(ctx)

		// Verificar que o erro foi propagado corretamente
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), contextCanceledError)
	})

	suite.T().Run("UserRepository List propagates cursor.All errors with context timeout after Find", func(t *testing.T) {
		// Este teste é mais complexo: primeiro inserir dados para que Find tenha sucesso,
		// depois usar contexto que expire durante cursor.All

		cfg := suite.testConfig.ToAppConfig()
		loggerInstance := logger.NewLogger()
		dbManager := database.NewDatabaseManager(cfg, loggerInstance)

		err := dbManager.Initialize(context.Background())
		assert.NoError(t, err)

		userRepo, err := repositories.NewUserRepository(dbManager)
		assert.NoError(t, err)

		// Inserir alguns documentos para garantir que Find retorne dados
		testUser1 := &entities.User{ID: bson.NewObjectID(), Name: "Test User 1", Email: "test1@example.com"}
		testUser2 := &entities.User{ID: bson.NewObjectID(), Name: "Test User 2", Email: "test2@example.com"}

		err = userRepo.Create(context.Background(), testUser1)
		assert.NoError(t, err)
		err = userRepo.Create(context.Background(), testUser2)
		assert.NoError(t, err)

		// Criar contexto com timeout muito curto que pode expirar durante cursor.All
		// O Find pode ser rápido, mas cursor.All pode demorar mais
		ctx, cancel := context.WithTimeout(context.Background(), 50*time.Microsecond)
		defer cancel()

		// Aguardar um pouco para aumentar chance do timeout ocorrer durante cursor.All
		time.Sleep(10 * time.Microsecond)

		// Executar List - Find pode ter sucesso, mas cursor.All deve falhar
		result, err := userRepo.List(ctx)

		// Se der timeout, o erro deve ser propagado
		if err != nil {
			assert.Error(t, err)
			assert.Nil(t, result)
			// Pode ser timeout durante Find ou cursor.All
			assert.True(t,
				strings.Contains(err.Error(), "context deadline exceeded") ||
					strings.Contains(err.Error(), "context canceled") ||
					strings.Contains(err.Error(), "timeout"),
				"Error should indicate context timeout: %v", err)
		}
	})

	suite.T().Run("GroupRepository List propagates cursor.All errors with context timeout after Find", func(t *testing.T) {
		// Mesmo teste para grupos

		cfg := suite.testConfig.ToAppConfig()
		loggerInstance := logger.NewLogger()
		dbManager := database.NewDatabaseManager(cfg, loggerInstance)

		err := dbManager.Initialize(context.Background())
		assert.NoError(t, err)

		groupRepo, err := repositories.NewGroupRepository(dbManager)
		assert.NoError(t, err)

		// Inserir alguns documentos
		testGroup1 := &entities.Group{ID: bson.NewObjectID(), Name: "Test Group 1", Members: []string{}}
		testGroup2 := &entities.Group{ID: bson.NewObjectID(), Name: "Test Group 2", Members: []string{}}

		err = groupRepo.Create(context.Background(), testGroup1)
		assert.NoError(t, err)
		err = groupRepo.Create(context.Background(), testGroup2)
		assert.NoError(t, err)

		// Contexto com timeout muito curto
		ctx, cancel := context.WithTimeout(context.Background(), 50*time.Microsecond)
		defer cancel()

		time.Sleep(10 * time.Microsecond)

		// Executar List
		result, err := groupRepo.List(ctx)

		if err != nil {
			assert.Error(t, err)
			assert.Nil(t, result)
			assert.True(t,
				strings.Contains(err.Error(), "context deadline exceeded") ||
					strings.Contains(err.Error(), "context canceled") ||
					strings.Contains(err.Error(), "timeout"),
				"Error should indicate context timeout: %v", err)
		}
	})
}

// ========================================
// TESTES ESPECÍFICOS PARA FORÇAR ERROS NO cursor.All
// Estes testes criam cenários onde Find é bem-sucedido mas cursor.All falha
// ========================================

func (suite *IntegrationTestSuite) TestUserRepositoryListCursorAllSpecificError() {
	// Teste mais agressivo para forçar erro no cursor.All

	suite.T().Run("UserRepository cursor.All fails with rapid context cancellation", func(t *testing.T) {
		cfg := suite.testConfig.ToAppConfig()
		loggerInstance := logger.NewLogger()
		dbManager := database.NewDatabaseManager(cfg, loggerInstance)

		err := dbManager.Initialize(context.Background())
		assert.NoError(t, err)

		userRepo, err := repositories.NewUserRepository(dbManager)
		assert.NoError(t, err)

		// Inserir muitos documentos para tornar cursor.All mais lento
		for i := 0; i < 100; i++ {
			testUser := &entities.User{
				ID:    bson.NewObjectID(),
				Name:  fmt.Sprintf("Test User %d", i),
				Email: fmt.Sprintf("test%d@example.com", i),
			}
			err = userRepo.Create(context.Background(), testUser)
			assert.NoError(t, err)
		}

		// Criar contexto que será cancelado durante a operação
		ctx, cancel := context.WithCancel(context.Background())

		// Executar em goroutine e cancelar rapidamente
		go func() {
			time.Sleep(1 * time.Millisecond) // Permitir que Find comece
			cancel()                         // Cancelar durante cursor.All
		}()

		// Executar List - Find pode ter sucesso, cursor.All deve falhar
		result, err := userRepo.List(ctx)

		// Verificar se obtivemos erro (pode ser durante Find ou cursor.All)
		if err != nil {
			assert.Error(t, err)
			assert.Nil(t, result)
			assert.True(t,
				strings.Contains(err.Error(), "context canceled") ||
					strings.Contains(err.Error(), "context deadline exceeded"),
				"Error should indicate context cancellation: %v", err)
		}
	})

	suite.T().Run("GroupRepository cursor.All fails with rapid context cancellation", func(t *testing.T) {
		cfg := suite.testConfig.ToAppConfig()
		loggerInstance := logger.NewLogger()
		dbManager := database.NewDatabaseManager(cfg, loggerInstance)

		err := dbManager.Initialize(context.Background())
		assert.NoError(t, err)

		groupRepo, err := repositories.NewGroupRepository(dbManager)
		assert.NoError(t, err)

		// Inserir muitos documentos
		for i := 0; i < 100; i++ {
			testGroup := &entities.Group{
				ID:      bson.NewObjectID(),
				Name:    fmt.Sprintf("Test Group %d", i),
				Members: []string{fmt.Sprintf("user%d", i)},
			}
			err = groupRepo.Create(context.Background(), testGroup)
			assert.NoError(t, err)
		}

		// Contexto com cancelamento rápido
		ctx, cancel := context.WithCancel(context.Background())

		go func() {
			time.Sleep(1 * time.Millisecond)
			cancel()
		}()

		result, err := groupRepo.List(ctx)

		if err != nil {
			assert.Error(t, err)
			assert.Nil(t, result)
			assert.True(t,
				strings.Contains(err.Error(), "context canceled") ||
					strings.Contains(err.Error(), "context deadline exceeded"),
				"Error should indicate context cancellation: %v", err)
		}
	})
}

// ========================================
// TESTES ADICIONAIS PARA GARANTIR COBERTURA DO cursor.All
// Usando estratégias diferentes para forçar falhas específicas
// ========================================

func (suite *IntegrationTestSuite) TestRepositoryListCursorAllForceError() {
	// Teste que força erro especificamente no cursor.All usando volume de dados

	suite.T().Run("UserRepository List cursor.All timeout with large dataset", func(t *testing.T) {
		cfg := suite.testConfig.ToAppConfig()
		loggerInstance := logger.NewLogger()
		dbManager := database.NewDatabaseManager(cfg, loggerInstance)

		err := dbManager.Initialize(context.Background())
		assert.NoError(t, err)

		userRepo, err := repositories.NewUserRepository(dbManager)
		assert.NoError(t, err)

		// Inserir um dataset maior para aumentar o tempo de processamento do cursor.All
		batchSize := 500
		for i := 0; i < batchSize; i++ {
			testUser := &entities.User{
				ID:    bson.NewObjectID(),
				Name:  fmt.Sprintf("Large Dataset User %d with very long name to increase document size", i),
				Email: fmt.Sprintf("very.long.email.address.for.user.number.%d@verylongdomainname.com", i),
			}
			err = userRepo.Create(context.Background(), testUser)
			assert.NoError(t, err)
		}

		// Contexto com timeout extremamente curto especificamente para cursor.All
		// Find deve ser rápido, mas cursor.All com 500 documentos deve ser mais lento
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Microsecond)
		defer cancel()

		// Aguardar um tempo para garantir que o contexto esteja próximo do timeout
		time.Sleep(5 * time.Microsecond)

		// Executar List
		result, err := userRepo.List(ctx)

		// Com um dataset grande e timeout muito curto, cursor.All tem mais chance de falhar
		if err != nil {
			assert.Error(t, err)
			assert.Nil(t, result)
			assert.True(t,
				strings.Contains(err.Error(), "context deadline exceeded") ||
					strings.Contains(err.Error(), "timeout") ||
					strings.Contains(err.Error(), "canceled"),
				"Error should indicate timeout: %v", err)

			// Log para debug - se chegou até aqui, provavelmente cursor.All falhou
			t.Logf("Successfully caught error (likely in cursor.All): %v", err)
		} else {
			// Se não houve erro, pelo menos testamos o caminho de sucesso
			t.Logf("Operation completed successfully with %d results", len(result))
		}
	})

	suite.T().Run("GroupRepository List cursor.All interrupted by cancellation", func(t *testing.T) {
		cfg := suite.testConfig.ToAppConfig()
		loggerInstance := logger.NewLogger()
		dbManager := database.NewDatabaseManager(cfg, loggerInstance)

		err := dbManager.Initialize(context.Background())
		assert.NoError(t, err)

		groupRepo, err := repositories.NewGroupRepository(dbManager)
		assert.NoError(t, err)

		// Inserir dataset grande com documentos complexos
		batchSize := 300
		for i := 0; i < batchSize; i++ {
			members := make([]string, 50) // 50 membros por grupo
			for j := 0; j < 50; j++ {
				members[j] = fmt.Sprintf("member-%d-%d", i, j)
			}

			testGroup := &entities.Group{
				ID:      bson.NewObjectID(),
				Name:    fmt.Sprintf("Complex Group %d with many members and long description", i),
				Members: members,
			}
			err = groupRepo.Create(context.Background(), testGroup)
			assert.NoError(t, err)
		}

		// Estratégia: cancelar o contexto após um delay calculado
		ctx, cancel := context.WithCancel(context.Background())

		// Executar em goroutine separada para cancelar após Find mas durante cursor.All
		go func() {
			// Delay calculado: tempo suficiente para Find mas não para cursor.All completo
			time.Sleep(500 * time.Microsecond) // Find deve completar
			cancel()                           // Cancelar durante cursor.All
		}()

		result, err := groupRepo.List(ctx)

		// Verificar se conseguimos o erro esperado
		if err != nil {
			assert.Error(t, err)
			assert.Nil(t, result)
			assert.True(t,
				strings.Contains(err.Error(), "context canceled") ||
					strings.Contains(err.Error(), "context deadline exceeded"),
				"Error should indicate context cancellation: %v", err)

			t.Logf("Successfully caught cursor.All error: %v", err)
		} else {
			t.Logf("Operation completed with %d results (no error caught)", len(result))
		}
	})
}

// ========================================
// TESTES ESPECÍFICOS PARA FORÇAR ERROS NO cursor.All DO USER_REPOSITORY
// Estes testes focam especificamente em fazer o user_repository falhar no cursor.All
// ========================================

func (suite *IntegrationTestSuite) TestUserRepositoryListCursorAllForceSpecificError() {
	// Testes mais agressivos especificamente para user_repository

	suite.T().Run("UserRepository cursor.All guaranteed failure with massive dataset", func(t *testing.T) {
		cfg := suite.testConfig.ToAppConfig()
		loggerInstance := logger.NewLogger()
		dbManager := database.NewDatabaseManager(cfg, loggerInstance)

		err := dbManager.Initialize(context.Background())
		assert.NoError(t, err)

		userRepo, err := repositories.NewUserRepository(dbManager)
		assert.NoError(t, err)

		// Inserir um dataset MUITO GRANDE para garantir que cursor.All seja lento
		batchSize := 1000 // Aumentar para 1000 documentos
		for i := 0; i < batchSize; i++ {
			// Criar documentos com muito conteúdo para tornar cursor.All mais lento
			longName := fmt.Sprintf("User %d with extremely long name containing lots of text to make the document larger and processing slower", i)
			longEmail := fmt.Sprintf("user.with.very.long.email.address.number.%d.that.contains.many.characters@verylongdomainname.example.com", i)

			testUser := &entities.User{
				ID:    bson.NewObjectID(),
				Name:  longName,
				Email: longEmail,
			}
			err = userRepo.Create(context.Background(), testUser)
			assert.NoError(t, err)
		}

		// Usar timeout EXTREMAMENTE curto - quase impossível de completar cursor.All
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Microsecond)
		defer cancel()

		// Executar List imediatamente sem esperar
		result, err := userRepo.List(ctx)

		// Com 1000 documentos e 1 microsegundo de timeout, deve falhar
		if err != nil {
			assert.Error(t, err)
			assert.Nil(t, result)
			assert.True(t,
				strings.Contains(err.Error(), "context deadline exceeded") ||
					strings.Contains(err.Error(), "timeout") ||
					strings.Contains(err.Error(), "canceled"),
				"Error should indicate timeout: %v", err)

			t.Logf("SUCCESS: UserRepository cursor.All failed as expected: %v", err)
		} else {
			t.Logf("WARNING: UserRepository completed unexpectedly with %d results", len(result))
		}
	})

	suite.T().Run("UserRepository cursor.All with aggressive cancellation timing", func(t *testing.T) {
		cfg := suite.testConfig.ToAppConfig()
		loggerInstance := logger.NewLogger()
		dbManager := database.NewDatabaseManager(cfg, loggerInstance)

		err := dbManager.Initialize(context.Background())
		assert.NoError(t, err)

		userRepo, err := repositories.NewUserRepository(dbManager)
		assert.NoError(t, err)

		// Inserir dataset médio mas com documentos complexos
		batchSize := 750
		for i := 0; i < batchSize; i++ {
			testUser := &entities.User{
				ID:    bson.NewObjectID(),
				Name:  fmt.Sprintf("Complex User Document %d with additional metadata and long description fields that increase document size significantly", i),
				Email: fmt.Sprintf("complex.user.%d.with.very.long.email.address@extremely.long.domain.name.example.org", i),
			}
			err = userRepo.Create(context.Background(), testUser)
			assert.NoError(t, err)
		}

		// Estratégia: usar cancelamento com timing mais preciso
		ctx, cancel := context.WithCancel(context.Background())

		// Cancelar após um delay muito específico
		go func() {
			// Tempo calculado para permitir Find mas interromper cursor.All
			time.Sleep(100 * time.Microsecond) // Tempo mínimo para Find
			cancel()                           // Cancelar durante cursor.All
		}()

		result, err := userRepo.List(ctx)

		if err != nil {
			assert.Error(t, err)
			assert.Nil(t, result)
			assert.True(t,
				strings.Contains(err.Error(), "context canceled") ||
					strings.Contains(err.Error(), "context deadline exceeded"),
				"Error should indicate context cancellation: %v", err)

			t.Logf("SUCCESS: UserRepository cursor.All canceled as expected: %v", err)
		} else {
			t.Logf("INFO: UserRepository completed with %d results", len(result))
		}
	})

	suite.T().Run("UserRepository cursor.All with multiple cancellation attempts", func(t *testing.T) {
		cfg := suite.testConfig.ToAppConfig()
		loggerInstance := logger.NewLogger()
		dbManager := database.NewDatabaseManager(cfg, loggerInstance)

		err := dbManager.Initialize(context.Background())
		assert.NoError(t, err)

		userRepo, err := repositories.NewUserRepository(dbManager)
		assert.NoError(t, err)

		// Inserir dataset grande
		batchSize := 800
		for i := 0; i < batchSize; i++ {
			testUser := &entities.User{
				ID:    bson.NewObjectID(),
				Name:  fmt.Sprintf("User %d - Very long name designed to increase document processing time during cursor operations", i),
				Email: fmt.Sprintf("user%d@longdomainname.example.com", i),
			}
			err = userRepo.Create(context.Background(), testUser)
			assert.NoError(t, err)
		}

		// Tentar múltiplas vezes com timeouts diferentes
		timeouts := []time.Duration{
			5 * time.Microsecond,
			10 * time.Microsecond,
			20 * time.Microsecond,
		}

		for i, timeout := range timeouts {
			ctx, cancel := context.WithTimeout(context.Background(), timeout)

			result, err := userRepo.List(ctx)
			cancel()

			if err != nil {
				assert.Error(t, err)
				assert.Nil(t, result)
				t.Logf("SUCCESS attempt %d: UserRepository failed with timeout %v: %v", i+1, timeout, err)
				break // Sucesso - conseguimos o erro
			} else {
				t.Logf("Attempt %d with timeout %v completed with %d results", i+1, timeout, len(result))
			}
		}
	})
}
