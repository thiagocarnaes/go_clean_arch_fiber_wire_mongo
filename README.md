# User Management API

Uma implementação de Clean Architecture usando Go, Fiber, Wire e MongoDB para gerenciamento de usuários e grupos.

[![Tests](https://github.com/thiagocarnaes/go_clean_arch_fiber_wire_mongo/actions/workflows/tests.yml/badge.svg)](https://github.com/thiagocarnaes/go_clean_arch_fiber_wire_mongo/actions/workflows/tests.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/thiagocarnaes/go_clean_arch_fiber_wire_mongo)](https://goreportcard.com/report/github.com/thiagocarnaes/go_clean_arch_fiber_wire_mongo)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

## 📋 Índice

- [Características](#-características)
- [Arquitetura](#️-arquitetura)
- [Tecnologias](#️-tecnologias)
- [Como Executar](#-como-executar)
- [Docker e Containerização](#-docker-e-containerização)
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

## 🚀 Como Executar

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

## 🐳 Docker e Containerização

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

## 📊 Logs e Monitoramento

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

#### 👤 Operações de Usuários

##### Criar Usuário
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

##### Buscar Usuário por ID
```bash
curl -X GET http://localhost:3000/api/v1/users/60d5ec49eb1d2c001f5e4b1a
```

**Resposta:**
```json
{
  "id": "60d5ec49eb1d2c001f5e4b1a",
  "name": "João Silva",
  "email": "joao@example.com"
}
```

##### Atualizar Usuário
```bash
curl -X PUT http://localhost:3000/api/v1/users/60d5ec49eb1d2c001f5e4b1a \
  -H "Content-Type: application/json" \
  -d '{
    "name": "João Santos",
    "email": "joao.santos@example.com"
  }'
```

**Resposta:**
```json
{
  "id": "60d5ec49eb1d2c001f5e4b1a",
  "name": "João Santos",
  "email": "joao.santos@example.com"
}
```

##### Excluir Usuário
```bash
curl -X DELETE http://localhost:3000/api/v1/users/60d5ec49eb1d2c001f5e4b1a
```

**Resposta:** Status 204 (No Content)

##### Listar Todos os Usuários
```bash
curl -X GET http://localhost:3000/api/v1/users/
```

**Resposta:**
```json
{
  "users": [
    {
      "id": "60d5ec49eb1d2c001f5e4b1a",
      "name": "João Silva",
      "email": "joao@example.com"
    },
    {
      "id": "60d5ec49eb1d2c001f5e4b1b",
      "name": "Maria Santos",
      "email": "maria@example.com"
    }
  ],
  "meta": {
    "total": 2,
    "per_page": 10,
    "page": 1,
    "total_pages": 1
  }
}
```

##### Listar Usuários com Paginação
```bash
# Página 2, 5 usuários por página
curl -X GET "http://localhost:3000/api/v1/users/?page=2&limit=5"
```

##### Buscar Usuários por Nome/Email
```bash
# Buscar usuários que contenham "joão" no nome ou email
curl -X GET "http://localhost:3000/api/v1/users/?search=joão"
```

**Resposta:**
```json
{
  "users": [
    {
      "id": "60d5ec49eb1d2c001f5e4b1a",
      "name": "João Silva",
      "email": "joao@example.com"
    }
  ],
  "meta": {
    "total": 1,
    "per_page": 10,
    "page": 1,
    "total_pages": 1
  }
}
```

#### 👥 Operações de Grupos

##### Criar Grupo
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
  "id": "60d5ec49eb1d2c001f5e4b1c",
  "name": "Desenvolvedores",
  "members": []
}
```

##### Buscar Grupo por ID
```bash
curl -X GET http://localhost:3000/api/v1/groups/60d5ec49eb1d2c001f5e4b1c
```

**Resposta:**
```json
{
  "id": "60d5ec49eb1d2c001f5e4b1c",
  "name": "Desenvolvedores",
  "members": [
    {
      "id": "60d5ec49eb1d2c001f5e4b1a",
      "name": "João Silva",
      "email": "joao@example.com"
    }
  ]
}
```

##### Atualizar Grupo
```bash
curl -X PUT http://localhost:3000/api/v1/groups/60d5ec49eb1d2c001f5e4b1c \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Desenvolvedores Senior"
  }'
```

**Resposta:**
```json
{
  "id": "60d5ec49eb1d2c001f5e4b1c",
  "name": "Desenvolvedores Senior",
  "members": []
}
```

##### Excluir Grupo
```bash
curl -X DELETE http://localhost:3000/api/v1/groups/60d5ec49eb1d2c001f5e4b1c
```

**Resposta:** Status 204 (No Content)

##### Listar Todos os Grupos
```bash
curl -X GET http://localhost:3000/api/v1/groups/
```

**Resposta:**
```json
{
  "groups": [
    {
      "id": "60d5ec49eb1d2c001f5e4b1c",
      "name": "Desenvolvedores",
      "members": []
    },
    {
      "id": "60d5ec49eb1d2c001f5e4b1d",
      "name": "Designers",
      "members": []
    }
  ],
  "meta": {
    "total": 2,
    "per_page": 10,
    "page": 1,
    "total_pages": 1
  }
}
```

##### Listar Grupos com Paginação
```bash
# Página 2, 5 grupos por página
curl -X GET "http://localhost:3000/api/v1/groups/?page=2&limit=5"
```

#### 🔗 Gerenciamento de Membros de Grupos

##### Adicionar Usuário ao Grupo
```bash
curl -X POST http://localhost:3000/api/v1/groups/60d5ec49eb1d2c001f5e4b1c/members/60d5ec49eb1d2c001f5e4b1a
```

**Resposta:**
```json
{
  "id": "60d5ec49eb1d2c001f5e4b1c",
  "name": "Desenvolvedores",
  "members": [
    {
      "id": "60d5ec49eb1d2c001f5e4b1a",
      "name": "João Silva",
      "email": "joao@example.com"
    }
  ]
}
```

##### Remover Usuário do Grupo
```bash
curl -X DELETE http://localhost:3000/api/v1/groups/60d5ec49eb1d2c001f5e4b1c/members/60d5ec49eb1d2c001f5e4b1a
```

**Resposta:**
```json
{
  "id": "60d5ec49eb1d2c001f5e4b1c",
  "name": "Desenvolvedores",
  "members": []
}
```

#### 🚫 Exemplos de Respostas de Erro

##### Usuário Não Encontrado
```bash
curl -X GET http://localhost:3000/api/v1/users/invalid-id
```

**Resposta:** Status 404
```json
{
  "error": "User not found"
}
```

##### Dados Inválidos
```bash
curl -X POST http://localhost:3000/api/v1/users/ \
  -H "Content-Type: application/json" \
  -d '{
    "name": "",
    "email": "invalid-email"
  }'
```

**Resposta:** Status 400
```json
{
  "error": "Validation failed",
  "details": [
    "Name is required",
    "Email must be a valid email address"
  ]
}
```

##### Grupo Não Encontrado
```bash
curl -X GET http://localhost:3000/api/v1/groups/invalid-id
```

**Resposta:** Status 404
```json
{
  "error": "Group not found"
}
```

#### 📝 Notas Importantes

- **Base URL**: Use `http://localhost:8080` se estiver executando via Docker
- **Content-Type**: Sempre inclua `Content-Type: application/json` para requests POST/PUT
- **IDs**: Substitua os IDs de exemplo pelos IDs reais retornados pelas APIs
- **Paginação**: Por padrão, a API retorna 10 itens por página (máximo 100)
- **Busca**: O parâmetro `search` funciona para nome e email de usuários (case-insensitive)
- **Metadados**: As respostas de listagem incluem um objeto `meta` com informações de paginação:
  - `total`: Total de registros encontrados
  - `per_page`: Número de itens por página
  - `page`: Página atual
  - `total_pages`: Total de páginas disponíveis

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

