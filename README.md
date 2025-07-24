# User Management API

Uma implementação de Clean Architecture usando Go, Fiber, Wire e MongoDB para gerenciamento de usuários e grupos.

[![Tests](https://github.com/thiagocarnaes/go_clean_arch_fiber_wire_mongo/actions/workflows/tests.yml/badge.svg)](https://github.com/thiagocarnaes/go_clean_arch_fiber_wire_mongo/actions/workflows/tests.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/thiagocarnaes/go_clean_arch_fiber_wire_mongo)](https://goreportcard.com/report/github.com/thiagocarnaes/go_clean_arch_fiber_wire_mongo)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

## 📋 Índice

- [Características](#-características)
- [Arquitetura](#-arquitetura)
- [Tecnologias](#-tecnologias)
- [Como Executar](#-como-executar)
- [Logs e Monitoramento](#-logs-e-monitoramento)
- [Testes](#-testes)
- [API Endpoints](#-api-endpoints)
- [Desenvolvimento](#-desenvolvimento)
- [Contribuição](#-contribuição)
- [Licença](#-licença)

## 🚀 Características

- ✅ **Clean Architecture** - Separação clara de responsabilidades
- ✅ **Dependency Injection** - Usando Google Wire
- ✅ **REST API** - Endpoints para gerenciamento de usuários e grupos
- ✅ **MongoDB** - Banco de dados NoSQL
- ✅ **Hot Reload** - Desenvolvimento com Air
- ✅ **Testes de Integração** - Suite completa de testes
- ✅ **Docker Support** - Containerização completa
- ✅ **CI/CD** - GitHub Actions
- ✅ **Logging** - Structured logging com Logrus
- ✅ **Validation** - Validação de dados de entrada
- ✅ **CORS** - Cross-Origin Resource Sharing

## 🏗️ Arquitetura

Este projeto segue os princípios da Clean Architecture, organizando o código em camadas bem definidas:

```
internal/
├── application/          # Camada de Aplicação
│   ├── dto/             # Data Transfer Objects
│   ├── mappers/         # Mapeadores entre entidades e DTOs
│   └── usecases/        # Casos de uso (regras de negócio da aplicação)
├── domain/              # Camada de Domínio
│   ├── entities/        # Entidades de negócio
│   └── interfaces/      # Interfaces/Contratos
├── infrastructure/      # Camada de Infraestrutura
│   ├── database/        # Configuração do banco de dados
│   ├── logger/          # Configuração de logging
│   ├── repositories/    # Implementação dos repositórios
│   └── web/            # Framework web (Fiber)
│       ├── controllers/ # Controladores HTTP
│       ├── middleware/  # Middlewares
│       └── routes/      # Definição de rotas
└── config/              # Configurações da aplicação
```

### Fluxo de Dados

```
HTTP Request → Controller → Use Case → Repository → Database
                   ↓            ↓          ↓
HTTP Response ← Controller ← Use Case ← Repository ← Database
```

## 🛠️ Tecnologias

- **[Go 1.24+](https://golang.org/)** - Linguagem de programação
- **[Fiber v2](https://gofiber.io/)** - Framework web rápido e expressivo
- **[MongoDB](https://www.mongodb.com/)** - Banco de dados NoSQL
- **[Wire](https://github.com/google/wire)** - Dependency injection
- **[Logrus](https://github.com/sirupsen/logrus)** - Structured logging
- **[Testify](https://github.com/stretchr/testify)** - Testing toolkit
- **[Air](https://github.com/air-verse/air)** - Hot reload
- **[Docker](https://www.docker.com/)** - Containerização

## � Como Executar

### Pré-requisitos

- Go 1.24+
- MongoDB 7.0+ (rodando localmente)
- Make (opcional, mas recomendado)

### Configuração

1. **Clone o repositório:**
```bash
git clone https://github.com/thiagocarnaes/go_clean_arch_fiber_wire_mongo.git
cd go_clean_arch_fiber_wire_mongo
```

2. **Configure as variáveis de ambiente:**
```bash
# O arquivo .env já está configurado com valores padrão
# Edite se necessário para seu ambiente
cat .env
```

Arquivo `.env` padrão:
```env
# MongoDB Configuration
MONGO_URI=mongodb://localhost:27017
MONGO_DB=user_management
PORT=:3000

# Datadog Configuration
DD_SOURCE=go
DD_SERVICE=user-management
DD_TAGS=env:dev,app:fiber
```

3. **Instale as dependências:**
```bash
go mod download
```

### Executando a Aplicação

#### Opção 1: Execução Direta (Recomendado)
```bash
# Executar a aplicação (certifique-se que o MongoDB está rodando)
go run main.go initApiServer
```

#### Opção 2: Usando Make
```bash
# Executar a aplicação via Make
make run

# Para executar os testes
make test-integration
```

#### Opção 3: Usando Docker Compose (Com Datadog Agent)
```bash
# 1. Configure as variáveis de ambiente do Datadog
cp .env.example .env
# Edite o arquivo .env e adicione sua DD_API_KEY

# 2. Execute todos os serviços
make docker-up

# 3. Verificar logs da aplicação
make docker-logs-app

# 4. Parar todos os serviços
make docker-down
```

**Serviços incluídos no Docker Compose:**
- **API**: `http://localhost:8080` - Aplicação principal
- **MongoDB**: `localhost:27017` - Banco de dados
- **MongoDB Express**: `http://localhost:8081` - Interface web para MongoDB (admin/admin)
- **Datadog Agent**: Coleta de métricas, logs e traces

A API estará disponível em `http://localhost:3000` (execução direta) ou `http://localhost:8080` (Docker)

## � Docker e Containerização

### Docker Compose

O projeto inclui um arquivo `docker-compose.yml` completo com todos os serviços necessários:

```yaml
services:
  app:                 # Aplicação Go
  mongodb:            # Banco de dados MongoDB 7.0
  mongo-express:      # Interface web para MongoDB
  datadog-agent:      # Agente Datadog para monitoramento
```

### Comandos Docker Disponíveis

```bash
# Construir imagem Docker
make docker-build

# Iniciar todos os serviços
make docker-up

# Parar todos os serviços  
make docker-down

# Ver logs de todos os serviços
make docker-logs

# Ver logs apenas da aplicação
make docker-logs-app

# Reiniciar todos os serviços
make docker-restart

# Limpeza completa (containers, volumes, imagens)
make docker-clean
```

### Configuração do Datadog

Para usar o monitoramento com Datadog, você precisa:

1. **Obter uma API Key do Datadog:**
   - Acesse: https://app.datadoghq.com/organization-settings/api-keys
   - Copie sua API key

2. **Configurar o arquivo .env:**
```bash
# Copie o arquivo de exemplo
cp .env.example .env

# Edite e adicione sua API key
DD_API_KEY=sua_api_key_aqui
DD_SITE=datadoghq.com  # ou datadoghq.eu para EU
```

3. **Iniciar os serviços:**
```bash
make docker-up
```

O agente Datadog coletará automaticamente:
- **Logs** da aplicação e containers
- **Métricas** de sistema e aplicação  
- **Traces** APM (se configurado)
- **Métricas Docker** dos containers

## �📊 Logs e Monitoramento

A aplicação utiliza **Logrus** para logging estruturado e está configurada para integração com **Datadog**.

### Configuração de Logs

As configurações de log são controladas pelas variáveis de ambiente no arquivo `.env`:

```env
# Datadog Configuration
DD_SOURCE=go              # Fonte dos logs
DD_SERVICE=user-management # Nome do serviço
DD_TAGS=env:dev,app:fiber # Tags para filtragem
```

### Exemplo de Logs

```json
{
  "timestamp": "2025-07-23T22:31:09.318Z",
  "level": "info",
  "message": "Successfully connected to MongoDB",
  "ddsource": "go",
  "service": "user-management", 
  "ddtags": "env:dev,app:fiber",
  "uri": "mongodb://localhost:27017/"
}
```

### Monitoramento com Datadog

Para habilitar o monitoramento com Datadog:

1. Configure as variáveis de ambiente apropriadas
2. Instale o Datadog Agent
3. Configure o Agent para coletar logs da aplicação

Os logs estruturados facilitam a análise e debugging da aplicação.

## 🧪 Testes

### GitHub Actions - CI/CD Pipeline

O projeto possui uma pipeline completa de CI/CD configurada no GitHub Actions que executa automaticamente nos seguintes eventos:

- **Push** para branches `main` e `develop`
- **Pull Requests** para `main` e `develop`

#### Steps da Pipeline de Testes

```yaml
name: Tests

on:
  push:
    branches: [ main, develop ]
  pull_request:
    branches: [ main, develop ]

jobs:
  test:
    runs-on: ubuntu-latest

    steps:
      # 1. Checkout do código
      - name: Checkout code
        uses: actions/checkout@v4

      # 2. Setup do Go
      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.24'

      # 3. Cache das dependências Go
      - name: Cache Go modules
        uses: actions/cache@v4
        with:
          path: |
            ~/go/pkg/mod
            ~/.cache/go-build
          key: ${{ runner.os }}-go-${{ hashFiles('**/go.sum') }}
          restore-keys: |
            ${{ runner.os }}-go-

      # 4. Download das dependências
      - name: Install dependencies
        run: go mod download

      # 5. Instalação do Wire
      - name: Install Wire
        run: go install github.com/google/wire/cmd/wire@latest

      # 6. Geração do código Wire
      - name: Generate Wire dependencies
        run: |
          cd cmd
          wire

      # 7. Testes de integração com Testcontainers
      - name: Run integration tests
        run: go test -v ./tests/...

      # 8. Geração de relatório de cobertura
      - name: Generate test coverage
        run: |
          go test -coverprofile=coverage.out -covermode=atomic ./internal/... ./tests/...
          go tool cover -html=coverage.out -o coverage.html

      # 9. Upload para Codecov
      - name: Upload coverage reports
        uses: codecov/codecov-action@v4
        with:
          file: ./coverage.out
          flags: unittests
          name: codecov-umbrella
          fail_ci_if_error: false
```

#### Jobs Executados

1. **🧪 Test Job**: Executa testes unitários e de integração
2. **🏗️ Build Job**: Compila a aplicação e gera artefatos
3. **🐳 Docker Job**: Constrói imagem Docker (apenas na branch main)

#### Configurações Importantes

- **Testcontainers**: Usa Testcontainers para criar instâncias temporárias do MongoDB
- **Cache Otimizado**: Cache das dependências Go para builds mais rápidos
- **Wire Auto-generation**: Gera código Wire automaticamente
- **Coverage Reports**: Upload automático para Codecov
- **Docker Disponível**: Requer Docker para executar Testcontainers

### Executando os Testes de Integração

Os testes de integração usam **Testcontainers** para criar uma instância isolada do MongoDB automaticamente. **Não é necessário ter o MongoDB instalado ou rodando localmente** - o Testcontainers cuida disso para você.

#### Executar todos os testes:
```bash
go test ./tests/ -v
```

#### Executar com cobertura de código:
```bash
make test-integration
```

> **Nota:** Os testes usam Testcontainers e criam automaticamente uma instância temporária do MongoDB. Certifique-se de ter o Docker rodando em sua máquina.

#### Executar teste específico:
```bash
# Executar apenas testes de usuário
go test ./tests/ -v -run TestIntegrationSuite/TestUserCRUD

# Executar apenas testes de grupo  
go test ./tests/ -v -run TestIntegrationSuite/TestGroupCRUD

# Executar cenários complexos
go test ./tests/ -v -run TestIntegrationSuite/TestCompleteUserGroupWorkflow
```

### Estrutura dos Testes

O projeto possui uma suite completa de testes organizados por funcionalidade:

- **`tests/user_integration_test.go`** - Testes CRUD de usuários
  - `TestUserCRUD` - Criar, ler, atualizar, deletar usuário
  - `TestUserNotFound` - Teste de usuário não encontrado
  - `TestCreateUserInvalidData` - Validação de dados inválidos
  - `TestListUsersEmpty` - Lista vazia de usuários
  - `TestMultipleUsers` - Múltiplos usuários

- **`tests/group_integration_test.go`** - Testes CRUD de grupos
  - `TestGroupCRUD` - Criar, ler, atualizar, deletar grupo
  - `TestGroupNotFound` - Teste de grupo não encontrado
  - `TestGroupMemberManagement` - Gerenciamento de membros
  - `TestAddNonExistentUserToGroup` - Adicionar usuário inexistente

- **`tests/complex_scenarios_test.go`** - Cenários complexos
  - `TestCompleteUserGroupWorkflow` - Workflow completo
  - `TestUserDeletionImpactOnGroups` - Impacto da deleção nos grupos
  - `TestConcurrentOperations` - Operações concorrentes
  - `TestDataConsistency` - Consistência de dados

### Cobertura de Código

```bash
# Gerar relatório de cobertura
make test-integration

# Ver cobertura detalhada
go tool cover -html=coverage-integration.out
```

### Monitorando Testes no GitHub Actions

#### Status Badges
O README inclui badges que mostram o status atual dos testes:

[![Tests](https://github.com/thiagocarnaes/go_clean_arch_fiber_wire_mongo/actions/workflows/tests.yml/badge.svg)](https://github.com/thiagocarnaes/go_clean_arch_fiber_wire_mongo/actions/workflows/tests.yml)

#### Visualizando Resultados
1. **Acesse a aba Actions** no repositório GitHub
2. **Clique no workflow "Tests"** para ver execuções recentes
3. **Clique em uma execução específica** para ver detalhes dos jobs
4. **Expand os steps** para ver logs detalhados de cada etapa

#### Artefatos Gerados
- **Coverage Reports**: Relatórios de cobertura em HTML
- **Build Binaries**: Executáveis compilados
- **Test Results**: Resultados detalhados dos testes

#### Notificações
O GitHub enviará notificações por email em caso de:
- ❌ Falhas nos testes
- ✅ Sucesso após correção de falhas
- 🔄 Status de builds em PRs

## 📚 API Endpoints

### Usuários

| Método | Endpoint        | Descrição           |
|--------|----------------|---------------------|
| POST   | `/api/v1/users/` | Criar usuário      |
| GET    | `/api/v1/users/:id` | Buscar usuário   |
| PUT    | `/api/v1/users/:id` | Atualizar usuário |
| DELETE | `/api/v1/users/:id` | Excluir usuário   |
| GET    | `/api/v1/users/` | Listar usuários    |

### Grupos

| Método | Endpoint                        | Descrição                |
|--------|---------------------------------|--------------------------|
| POST   | `/api/v1/groups/`              | Criar grupo              |
| GET    | `/api/v1/groups/:id`           | Buscar grupo             |
| PUT    | `/api/v1/groups/:id`           | Atualizar grupo          |
| DELETE | `/api/v1/groups/:id`           | Excluir grupo            |
| GET    | `/api/v1/groups/`              | Listar grupos            |
| POST   | `/api/v1/groups/:groupId/members/:userId` | Adicionar usuário ao grupo |
| DELETE | `/api/v1/groups/:groupId/members/:userId` | Remover usuário do grupo   |

### Exemplos de Uso

#### Criar Usuário
```bash
curl -X POST http://localhost:3000/api/v1/users/ \
  -H "Content-Type: application/json" \
  -d '{
    "name": "João Silva",
    "email": "joao@example.com"
  }'
```

**Resposta:**
```json
{
  "id": "60d5ec49eb1d2c001f5e4b1a",
  "name": "João Silva", 
  "email": "joao@example.com"
}
```

#### Criar Grupo
```bash
curl -X POST http://localhost:3000/api/v1/groups/ \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Desenvolvedores"
  }'
```

**Resposta:**
```json
{
  "id": "60d5ec49eb1d2c001f5e4b1b",
  "name": "Desenvolvedores",
  "members": []
}
```

#### Adicionar Usuário ao Grupo
```bash
# Substitua {groupId} e {userId} pelos IDs reais obtidos nas respostas das APIs
curl -X POST http://localhost:3000/api/v1/groups/{groupId}/members/{userId}
```

#### Buscar Usuário
```bash
# Substitua {userId} pelo ID real
curl -X GET http://localhost:3000/api/v1/users/{userId}
```

#### Listar Todos os Usuários
```bash
curl -X GET http://localhost:3000/api/v1/users/
```

## 🔧 Desenvolvimento

### Hot Reload

```bash
# Usar Air para hot reload
make dev
```

### Regenerar Wire

```bash
# Após modificar dependências
make wire
```

### Linting e Formatação

```bash
# Executar todas as verificações de qualidade
make check

# Apenas linting
make lint

# Apenas formatação
make fmt
```

### Comandos Úteis

```bash
# Ver todos os comandos disponíveis
make help

# Limpar artefatos de build
make clean

# Atualizar dependências
make deps-update
```

## 🤝 Contribuição

1. Fork o projeto
2. Crie uma branch para sua feature (`git checkout -b feature/AmazingFeature`)
3. Commit suas mudanças (`git commit -m 'Add some AmazingFeature'`)
4. Push para a branch (`git push origin feature/AmazingFeature`)
5. Abra um Pull Request

### Guidelines

- Siga os princípios da Clean Architecture
- Mantenha alta cobertura de testes
- Use conventional commits
- Execute `make check` antes de commit
- Adicione testes para novas funcionalidades

## 📄 Licença

Este projeto está licenciado sob a Licença MIT - veja o arquivo [LICENSE](LICENSE) para detalhes.

## 📞 Suporte

Se você encontrar algum problema ou tiver dúvidas:

1. Verifique a [documentação](README.md)
2. Consulte os [testes de integração](tests/README.md) para exemplos
3. Abra uma [issue](https://github.com/thiagocarnaes/go_clean_arch_fiber_wire_mongo/issues)

## 🙏 Agradecimentos

- [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html) por Uncle Bob
- [Fiber](https://gofiber.io/) pela excelente framework web
- [Wire](https://github.com/google/wire) pela dependency injection
- Comunidade Go pelo suporte e ferramentas incríveis

## Funcionalidades
- **API RESTful**: Criação e recuperação de usuários via `POST /users` e `GET /users/:id`.
- **Arquitetura Limpa**: Separa as responsabilidades em camadas de domínio, casos de uso, interfaces e infraestrutura.
- **MongoDB**: Armazena dados de usuários com validação de conexão e configuração baseada em URI.
- **Encerramento Limpo**: Lida com sinais `SIGINT` e `SIGTERM` para fechar conexões do servidor e banco de dados de forma limpa.
- **Logs**: Logs estruturados em JSON com Logrus, compatíveis com Datadog.
- **Testes de Integração**: Testes para endpoints de usuários e encerramento gracioso usando `dockertest`.

## Pré-requisitos
- **Go**: Versão 1.20 ou superior
- **Docker**: Para executar MongoDB e o agente Datadog
- **MongoDB**: Instância local ou em container Docker
- **Git**: Para clonar o repositório
- **Wire**: Para injeção de dependências (`go install github.com/google/wire/cmd/wire@latest`)

## Estrutura do Projeto
```
user-management/
├── cmd/                    # Ponto de entrada e configuração do Wire
│   ├── root.go             # Comando CLI principal
│   ├── wire.go             # Configuração de injeção de dependências do Wire
│   └── wire_gen.go         # Código gerado pelo Wire
├── internal/               # Código da aplicação
│   ├── config/             # Configurações (URI do MongoDB, porta)
│   ├── domain/             # Entidades e interfaces
│   ├── infrastructure/      # Banco de dados, servidor web e logger
│   ├── interfaces/         # Handlers HTTP
│   └── usecases/           # Lógica de negócio
├── tests/                  # Testes de integração
│   └── integration_test.go # Testes para endpoints de usuários e encerramento
├── go.mod                  # Dependências do módulo Go
└── README.md               # Este arquivo
```

## Instruções de Configuração

### 1. Clonar o Repositório
```bash
git clone https://github.com/thiagocarnaes/go_clean_arch_fiber_wire_mongo.git
cd go_clean_arch_fiber_wire_mongo
```

### 2. Instalar Dependências
Instale as dependências do Go:
```bash
go mod tidy
```

Instale o Wire:
```bash
go install github.com/google/wire/cmd/wire@latest
```

### 3. Gerar Código do Wire
Gere o código de injeção de dependências:
```bash
wire ./cmd
```
Isso cria o arquivo `wire_gen.go` no diretório `cmd`.

### 4. Executar o MongoDB
**Opção 1: MongoDB Local**
Inicie uma instância local do MongoDB:
```bash
mongod
```

**Opção 2: Docker**
Execute o MongoDB em um container Docker:
```bash
docker run -d -p 27017:27017 --name mongo-container mongo:latest
```

Crie o banco de dados `user_management` (necessário para testes):
```bash
mongo
use user_management
db.users.insertOne({"test": "data"})
exit
```

### 5. Executar a Aplicação
Inicie o servidor da API:
```bash
go run main.go initApiServer
```
O servidor estará disponível em `http://localhost:3000`.

### 6. Testar a API
**Criar um Usuário**:
```bash
curl -X POST http://localhost:3000/users -H "Content-Type: application/json" -d '{"id":"123","name":"Usuário de Teste","email":"teste@exemplo.com"}'
```

**Recuperar um Usuário**:
```bash
curl -X GET http://localhost:3000/users/123
```

### 7. Executar Testes de Integração
Os testes de integração em `tests/integration_test.go` usam `dockertest` para iniciar um container MongoDB.

Instale as dependências de teste:
```bash
go get github.com/stretchr/testify
go get github.com/ory/dockertest/v3
go mod tidy
```

Execute os testes:
```bash
cd tests
go test -v
```

Saída esperada:
```
=== RUN   TestIntegration
{"timestamp":"2025-07-14T22:20:05.123Z","level":"info","message":"Iniciando testes de integração","ddsource":"go","service":"user-management","ddtags":"env:test,app:fiber"}
=== RUN   TestIntegration/CreateUser
=== RUN   TestIntegration/GetUser
=== RUN   TestIntegration/GracefulShutdown
--- PASS: TestIntegration (3.45s)
    --- PASS: TestIntegration/CreateUser (0.12s)
    --- PASS: TestIntegration/GetUser (0.08s)
    --- PASS: TestIntegration/GracefulShutdown (2.05s)
PASS
ok      user-management/tests    3.456s
```

### 8. Encerramento Limpo
Para testar o encerramento gracioso, inicie o servidor:
```bash
go run ./cmd
```

Em seguida, envie um sinal `SIGTERM` ou pressione `Ctrl+C`:
```bash
kill -TERM <pid>
```

Logs esperados:
```json
{"timestamp":"2025-07-14T22:20:10.123Z","level":"info","message":"Sinal de encerramento recebido, iniciando encerramento gracioso","ddsource":"go","service":"user-management","ddtags":"env:dev,app:fiber"}
{"timestamp":"2025-07-14T22:20:10.456Z","level":"info","message":"Conexão com o banco de dados validada antes do encerramento","database":"user_management","ddsource":"go","service":"user-management","ddtags":"env:dev,app:fiber"}
{"timestamp":"2025-07-14T22:20:10.789Z","level":"info","message":"Conexão com MongoDB fechada com sucesso","database":"user_management","ddsource":"go","service":"user-management","ddtags":"env:dev,app:fiber"}
{"timestamp":"2025-07-14T22:20:10.890Z","level":"info","message":"Servidor encerrado graciosamente","ddsource":"go","service":"user-management","ddtags":"env:dev,app:fiber"}
```

### 9. Integração com Datadog
Para enviar logs ao Datadog:

**Instalar o Agente Datadog**:
```bash
docker run -d -v /var/run/docker.sock:/var/run/docker.sock:ro \
    -v /proc/:/host/proc/:ro \
    -v /sys/fs/cgroup/:/host/sys/fs/cgroup:ro \
    -e DD_API_KEY=<SUA_CHAVE_API> \
    -e DD_LOGS_ENABLED=true \
    --name datadog-agent \
    datadog/agent:latest
```

**Configurar o Agente**:
Edite `/etc/datadog-agent/datadog.yaml`:
```yaml
logs_enabled: true
```

Crie `/etc/datadog-agent/conf.d/go.d/conf.yaml`:
```yaml
logs:
  - type: docker
    service: user-management
    source: go
    log_processing_rules:
      - type: multi_line
        pattern: '^{'
        name: json_log
```

**Verificar Logs**:
Acesse o painel do Datadog e filtre por `service:user-management` ou `source:go`.

## Solução de Problemas
- **Erros de Conexão com MongoDB**:
    - Verifique se o MongoDB está executando em `localhost:27017` e se o banco `user_management` existe.
    - Consulte os logs para erros como `"MongoDB is not available"`.
- **Erros do Wire**:
    - Execute `wire -v` para depurar problemas de injeção de dependências:
      ```bash
      wire -v
      ```
    - Regenere o `wire_gen.go`:
      ```bash
      cd cmd
      wire
      ```
- **Falhas nos Testes**:
    - Certifique-se de que o Docker está em execução para o `dockertest`.
    - Verifique os logs dos testes para erros de inicialização do MongoDB ou servidor.
- **Panic no Encerramento**:
    - Confirme que `server.go` e `mongodb.go` correspondem às implementações mais recentes.
    - Compartilhe o stack trace do panic para depuração adicional.

