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
	"user-management/internal/config"
	"user-management/internal/infrastructure/database"
	"user-management/internal/infrastructure/logger"
	irepos "user-management/internal/infrastructure/repositories"
	"user-management/internal/infrastructure/web"
	"user-management/internal/infrastructure/web/controllers"

	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const (
	testDBName   = "user_management_test"
	testMongoURI = "mongodb://localhost:27017"
	testPort     = ":3001"
)

type IntegrationTestSuite struct {
	suite.Suite
	server      *web.Server
	client      *http.Client
	baseURL     string
	mongoClient *mongo.Client
}

func (suite *IntegrationTestSuite) SetupSuite() {
	// Configurar cliente MongoDB para limpeza
	client, err := mongo.Connect(options.Client().ApplyURI(testMongoURI))
	suite.NoError(err)
	suite.mongoClient = client

	// Configurar aplicação
	cfg := &config.Config{
		MongoURI: testMongoURI,
		MongoDB:  testDBName,
		Port:     testPort,
	}

	// Inicializar dependências
	loggerInstance := logger.NewLogger()
	mongodb, err := database.NewMongoDB(cfg, loggerInstance)
	suite.NoError(err)

	userRepo := irepos.NewUserRepository(mongodb)
	groupRepo := irepos.NewGroupRepository(mongodb)

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
	suite.server = web.NewServer(cfg, userController, groupController, loggerInstance, mongodb)

	suite.client = &http.Client{Timeout: 10 * time.Second}
	suite.baseURL = fmt.Sprintf("http://localhost%s", testPort)

	// Iniciar servidor em goroutine
	go func() {
		suite.server.Start()
	}()

	// Aguardar servidor inicializar
	time.Sleep(2 * time.Second)
}

func (suite *IntegrationTestSuite) TearDownSuite() {
	if suite.mongoClient != nil {
		suite.mongoClient.Disconnect(context.Background())
	}
}

func (suite *IntegrationTestSuite) SetupTest() {
	// Limpar banco de dados antes de cada teste
	ctx := context.Background()
	suite.mongoClient.Database(testDBName).Drop(ctx)
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
