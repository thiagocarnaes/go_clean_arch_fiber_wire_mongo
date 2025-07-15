# API de Gerenciamento de Usuários

Esta é uma API de gerenciamento de usuários construída com Go, seguindo os princípios da Arquitetura Limpa (Clean Architecture). Utiliza o Fiber v2 como framework web, MongoDB como banco de dados, Wire para injeção de dependências e Logrus para logs com integração ao Datadog. A API oferece endpoints para criar e recuperar usuários, com encerramento gracioso (graceful shutdown) e testes de integração.

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

