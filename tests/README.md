# Integration Tests

Este diretório contém os testes de integração para a API de gerenciamento de usuários. Os testes verificam o funcionamento completo da aplicação, incluindo controladores, casos de uso, repositórios e integração com o banco de dados MongoDB.

## Estrutura dos Testes

### Arquivos de Teste

- **`integration_test.go`** - Configuração base da suite de testes de integração
- **`user_integration_test.go`** - Testes CRUD para usuários
- **`group_integration_test.go`** - Testes CRUD para grupos
- **`complex_scenarios_test.go`** - Cenários complexos e testes de workflow
- **`test_config.go`** - Configurações e utilitários para testes

### Cenários Testados

#### Testes de Usuários (`user_integration_test.go`)
- ✅ Criação, leitura, atualização e exclusão de usuários
- ✅ Busca de usuário não existente
- ✅ Listagem de usuários (vazia e com múltiplos usuários)
- ✅ Validação de dados inválidos
- ✅ Operações com múltiplos usuários

#### Testes de Grupos (`group_integration_test.go`) 
- ✅ Criação, leitura, atualização e exclusão de grupos
- ✅ Gerenciamento de membros (adicionar/remover usuários)
- ✅ Busca de grupo não existente
- ✅ Tentativas de adicionar usuários inexistentes a grupos
- ✅ Tentativas de adicionar usuários a grupos inexistentes
- ✅ Listagem de grupos (vazia e com múltiplos grupos)

#### Cenários Complexos (`complex_scenarios_test.go`)
- ✅ Workflow completo de usuários e grupos
- ✅ Impacto da exclusão de usuários nos grupos
- ✅ Exclusão de grupos com membros
- ✅ Operações concorrentes (simuladas)
- ✅ Consistência de dados entre operações

## Pré-requisitos

### MongoDB
Os testes requerem uma instância do MongoDB rodando. Você pode usar:

1. **MongoDB local**:
   ```bash
   # Instalar e iniciar MongoDB localmente
   # Ou usar Docker:
   docker run --name mongo-test -p 27017:27017 -d mongo:7.0
   ```

2. **MongoDB via Docker** (recomendado):
   ```bash
   make mongo-start  # Inicia MongoDB via Docker
   ```

### Dependências Go
```bash
go mod download
```

## Executando os Testes

### Opção 1: Usando o Makefile (Recomendado)

```bash
# Executar testes de integração (MongoDB deve estar rodando)
make test-integration

# Executar testes com MongoDB via Docker (mais conveniente)
make test-integration-docker

# Executar todos os testes (unit + integration)
make test

# Executar testes com relatório de cobertura
make test-coverage
```

### Opção 2: Comando Go Direto

```bash
# MongoDB deve estar rodando em localhost:27017
go test -v -race ./tests/...
```

### Opção 3: Testes Individuais

```bash
# Testar apenas usuários
go test -v -race ./tests/ -run TestIntegrationSuite/TestUser

# Testar apenas grupos  
go test -v -race ./tests/ -run TestIntegrationSuite/TestGroup

# Testar cenários complexos
go test -v -race ./tests/ -run TestIntegrationSuite/TestComplete
```

## Configuração dos Testes

### Variáveis de Ambiente

Os testes podem ser configurados através das seguintes variáveis de ambiente:

```bash
# URI do MongoDB para testes (padrão: mongodb://localhost:27017)
export TEST_MONGO_URI="mongodb://localhost:27017"

# Nome do banco de dados de teste (padrão: user_management_test)
export TEST_MONGO_DB="user_management_test"

# Porta para o servidor de teste (padrão: :3001)
export TEST_PORT=":3001"
```

### Limpeza Automática

- Os testes limpam automaticamente o banco de dados antes de cada teste
- Cada teste é executado em isolamento
- O banco de teste é completamente removido após a suite

## Estrutura da Suite de Testes

```go
type IntegrationTestSuite struct {
    suite.Suite
    server      *web.Server     // Servidor web da aplicação
    client      *http.Client    // Cliente HTTP para requisições
    baseURL     string          // URL base para testes
    mongoClient *mongo.Client   // Cliente MongoDB para limpeza
}
```

### Lifecycle dos Testes

1. **SetupSuite**: Executado uma vez no início
   - Conecta ao MongoDB
   - Inicializa todas as dependências da aplicação
   - Inicia o servidor web
   
2. **SetupTest**: Executado antes de cada teste
   - Limpa o banco de dados de teste
   
3. **TearDownSuite**: Executado uma vez no final
   - Desconecta do MongoDB
   - Para o servidor

## Exemplos de Uso

### Testando Criação de Usuário
```go
func (suite *IntegrationTestSuite) TestUserCreation() {
    user := dto.UserDTO{
        ID:    "test-user",
        Name:  "Test User", 
        Email: "test@example.com",
    }
    
    resp, body := suite.makeRequest("POST", "/api/v1/users/", user)
    assert.Equal(suite.T(), http.StatusCreated, resp.StatusCode)
    
    var createdUser dto.UserDTO
    json.Unmarshal(body, &createdUser)
    assert.Equal(suite.T(), user.Name, createdUser.Name)
}
```

### Testando Workflow Completo
```go  
func (suite *IntegrationTestSuite) TestCompleteWorkflow() {
    // 1. Criar usuário
    // 2. Criar grupo
    // 3. Adicionar usuário ao grupo
    // 4. Verificar membro
    // 5. Remover usuário do grupo
    // 6. Verificar remoção
}
```

## Debugging dos Testes

### Logs Detalhados
```bash
go test -v -race ./tests/ -args -test.v
```

### Executar Teste Específico
```bash
go test -v ./tests/ -run TestIntegrationSuite/TestUserCRUD
```

### Manter MongoDB Após Falha
Para debugging, você pode comentar a limpeza do banco no `TearDownSuite` e inspecionar os dados:

```go
// Comentar temporariamente para debugging
// suite.mongoClient.Database(testDBName).Drop(ctx)
```

## Troubleshooting

### Erro: "connection refused"
- Verifique se o MongoDB está rodando: `docker ps` ou `sudo systemctl status mongod`
- Verifique a URL de conexão na configuração

### Erro: "database not found" 
- Normal - o banco de teste é criado automaticamente

### Erro: "port already in use"
- Altere a variável `TEST_PORT` para uma porta diferente
- Ou pare outros serviços usando a porta 3001

### Testes Lentos
- Considere usar MongoDB em memória para testes mais rápidos
- Verifique se há muitos logs sendo gerados

## Melhorias Futuras

- [ ] Testes de performance com carga
- [ ] Testes de failover do banco de dados  
- [ ] Testes com dados inválidos mais abrangentes
- [ ] Mocks para testes unitários dos controladores
- [ ] Testes de segurança e validação de entrada
- [ ] Testes de limite de rate limiting (se implementado)

## Contribuindo

Ao adicionar novos testes:

1. Siga o padrão de nomenclatura: `TestXXXX` 
2. Use a suite de testes existente
3. Limpe dados entre testes
4. Adicione documentação para cenários complexos
5. Verifique se os testes passam isoladamente e em conjunto
