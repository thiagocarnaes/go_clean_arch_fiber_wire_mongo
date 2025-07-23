# User Management API

Uma implementação de Clean Architecture usando Go, Fiber, Wire e MongoDB para gerenciamento de usuários e grupos.

[![Tests](https://github.com/thiagocarnaes/go_clean_arch_fiber_wire_mongo/actions/workflows/tests.yml/badge.svg)](https://github.com/thiagocarnaes/go_clean_arch_fiber_wire_mongo/actions/workflows/tests.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/thiagocarnaes/go_clean_arch_fiber_wire_mongo)](https://goreportcard.com/report/github.com/thiagocarnaes/go_clean_arch_fiber_wire_mongo)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

## 📋 Índice

- [Características](#-características)
- [Arquitetura](#-arquitetura)
- [Tecnologias](#-tecnologias)
- [Instalação](#-instalação)
- [Uso](#-uso)
- [API Endpoints](#-api-endpoints)
- [Testes](#-testes)
- [Desenvolvimento](#-desenvolvimento)
- [Docker](#-docker)
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

## 🔧 Instalação

### Pré-requisitos

- Go 1.24+
- MongoDB 7.0+
- Make (opcional, mas recomendado)
- Docker e Docker Compose (opcional)

### Setup Rápido

```bash
# Clone o repositório
git clone https://github.com/thiagocarnaes/go_clean_arch_fiber_wire_mongo.git
cd go_clean_arch_fiber_wire_mongo

# Setup completo do ambiente de desenvolvimento
make setup

# Copie e configure o arquivo de ambiente
cp .env.example .env
# Edite o .env com suas configurações

# Inicie o MongoDB (via Docker)
make mongo-start

# Execute a aplicação
make run
```

### Instalação Manual

```bash
# Instale as dependências
go mod download

# Instale as ferramentas necessárias
go install github.com/google/wire/cmd/wire@latest
go install github.com/air-verse/air@latest

# Gere as dependências do Wire
cd cmd && wire && cd ..

# Execute a aplicação
go run .
```

## 🚀 Uso

### Variáveis de Ambiente

Crie um arquivo `.env` baseado no `.env.example`:

```env
# MongoDB Configuration
MONGO_URI=mongodb://localhost:27017
MONGO_DB=user_management

# Server Configuration  
PORT=:8080

# Test Configuration (opcional)
TEST_MONGO_URI=mongodb://localhost:27017
TEST_MONGO_DB=user_management_test
TEST_PORT=:3001
```

### Executando a Aplicação

```bash
# Desenvolvimento (com hot reload)
make dev

# Produção
make run

# Via Docker Compose
make up
```

A API estará disponível em `http://localhost:8080`

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
curl -X POST http://localhost:8080/api/v1/users/ \
  -H "Content-Type: application/json" \
  -d '{
    "id": "user1",
    "name": "João Silva",
    "email": "joao@example.com"
  }'
```

#### Criar Grupo
```bash
curl -X POST http://localhost:8080/api/v1/groups/ \
  -H "Content-Type: application/json" \
  -d '{
    "id": "developers",
    "name": "Desenvolvedores",
    "members": []
  }'
```

#### Adicionar Usuário ao Grupo
```bash
curl -X POST http://localhost:8080/api/v1/groups/developers/members/user1
```

## 🧪 Testes

O projeto possui uma suite completa de testes de integração que testa todo o fluxo da aplicação.

### Executar Testes

```bash
# Todos os testes
make test

# Apenas testes de integração (MongoDB deve estar rodando)
make test-integration

# Testes de integração com MongoDB via Docker
make test-integration-docker

# Testes com relatório de cobertura
make test-coverage
```

### Estrutura dos Testes

- **`tests/integration_test.go`** - Setup da suite de testes
- **`tests/user_integration_test.go`** - Testes CRUD de usuários
- **`tests/group_integration_test.go`** - Testes CRUD de grupos
- **`tests/complex_scenarios_test.go`** - Cenários complexos e workflows

Para mais detalhes, consulte [tests/README.md](tests/README.md).

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

## 🐳 Docker

### Docker Compose (Recomendado)

```bash
# Iniciar todos os serviços
make up

# Ver logs
make logs

# Parar serviços
make down

# Rebuild e restart
make rebuild
```

O Docker Compose inclui:
- Aplicação Go
- MongoDB
- Mongo Express (interface web para MongoDB)

### Docker Manual

```bash
# Build da imagem
make docker-build

# Executar container
make docker-run
```

### Acessos

- **API**: http://localhost:8080
- **Mongo Express**: http://localhost:8081 (admin/admin)

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

