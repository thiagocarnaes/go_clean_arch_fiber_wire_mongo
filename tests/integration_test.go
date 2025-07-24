package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"
	"user-management/internal/application/usecases/group"
	"user-management/internal/application/usecases/user"
	"user-management/internal/infrastructure/database"
	"user-management/internal/infrastructure/logger"
	irepos "user-management/internal/infrastructure/repositories"
	"user-management/internal/infrastructure/web"
	"user-management/internal/infrastructure/web/controllers"

	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type IntegrationTestSuite struct {
	suite.Suite
	server      *web.Server
	client      *http.Client
	baseURL     string
	mongoClient *mongo.Client
	testConfig  *TestConfig
}

func (suite *IntegrationTestSuite) SetupSuite() {
	ctx := context.Background()

	// Configurar teste com Testcontainers
	testConfig, err := SetupTestContainer(ctx)
	suite.NoError(err)
	suite.testConfig = testConfig

	// Configurar cliente MongoDB para limpeza
	client, err := mongo.Connect(options.Client().ApplyURI(testConfig.MongoURI))
	suite.NoError(err)
	suite.mongoClient = client

	// Configurar aplicação
	cfg := testConfig.ToAppConfig()

	// Inicializar dependências
	loggerInstance := logger.NewLogger()
	dbManager := database.NewDatabaseManager(cfg, loggerInstance)

	// Inicializar conexão com o banco
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = dbManager.Initialize(ctx)
	suite.NoError(err)

	userRepo, err := irepos.NewUserRepository(dbManager)
	suite.NoError(err)

	groupRepo, err := irepos.NewGroupRepository(dbManager)
	suite.NoError(err)

	// Criar casos de uso
	createUserUseCase := user.NewCreateUserUseCase(userRepo)
	getUserUseCase := user.NewGetUserUseCase(userRepo)
	updateUserUseCase := user.NewUpdateUserUseCase(userRepo)
	deleteUserUseCase := user.NewDeleteUserUseCase(userRepo)
	listUsersUseCase := user.NewListUsersUseCase(userRepo)

	createGroupUseCase := group.NewCreateGroupUseCase(groupRepo)
	getGroupUseCase := group.NewGetGroupUseCase(groupRepo)
	updateGroupUseCase := group.NewUpdateGroupUseCase(groupRepo)
	deleteGroupUseCase := group.NewDeleteGroupUseCase(groupRepo)
	listGroupsUseCase := group.NewListGroupsUseCase(groupRepo)
	addUserToGroupUseCase := group.NewAddUserToGroupUseCase(groupRepo, userRepo)
	removeUserFromGroupUseCase := group.NewRemoveUserFromGroupUseCase(groupRepo)

	// Criar controladores
	userController := controllers.NewUserController(createUserUseCase, getUserUseCase, updateUserUseCase, deleteUserUseCase, listUsersUseCase)
	groupController := controllers.NewGroupController(createGroupUseCase, getGroupUseCase, updateGroupUseCase, deleteGroupUseCase, listGroupsUseCase, addUserToGroupUseCase, removeUserFromGroupUseCase)

	// Criar server
	suite.server = web.NewServer(cfg, userController, groupController, loggerInstance, dbManager)

	suite.client = &http.Client{Timeout: 10 * time.Second}
	suite.baseURL = fmt.Sprintf("http://localhost%s", testConfig.Port)

	// Iniciar servidor em goroutine
	errChan := make(chan error)
	go func() {
		err := suite.server.Start()
		errChan <- err
	}()

	// Aguardar servidor inicializar e verificar se está respondendo
	ready := false
	for i := 0; i < 10; i++ {
		resp, err := http.Get(fmt.Sprintf("%s/health", suite.baseURL))
		if err == nil {
			resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				ready = true
				break
			}
		}
		time.Sleep(100 * time.Millisecond)
	}
	suite.True(ready, "Server did not start within timeout")
}

func (suite *IntegrationTestSuite) TearDownSuite() {
	ctx := context.Background()

	if suite.mongoClient != nil {
		suite.mongoClient.Disconnect(ctx)
	}

	// Stop test container
	if suite.testConfig != nil {
		suite.testConfig.StopMongoContainer(ctx)
	}
}

func (suite *IntegrationTestSuite) SetupTest() {
	// Limpar banco de dados antes de cada teste
	ctx := context.Background()
	suite.mongoClient.Database(suite.testConfig.MongoDB).Drop(ctx)
}

func (suite *IntegrationTestSuite) makeRequest(method, path string, body interface{}) (*http.Response, []byte) {
	var bodyReader io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		suite.NoError(err)
		bodyReader = bytes.NewBuffer(jsonBody)
	}

	req, err := http.NewRequest(method, suite.baseURL+path, bodyReader)
	suite.NoError(err)

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := suite.client.Do(req)
	suite.NoError(err)

	responseBody, err := io.ReadAll(resp.Body)
	suite.NoError(err)
	resp.Body.Close()

	return resp, responseBody
}

func TestIntegrationSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}
