#!/bin/bash

# Criar estrutura completa do projeto com usecases
echo "Criando estrutura completa do projeto com usecases..."

# Criar estrutura de pastas
mkdir -p cmd
mkdir -p internal/application/dto
mkdir -p internal/application/mappers
mkdir -p internal/application/usecases/user
mkdir -p internal/application/usecases/group
mkdir -p internal/domain/entities
mkdir -p internal/domain/interfaces/repositories
mkdir -p internal/infrastructure/database
mkdir -p internal/infrastructure/repositories
mkdir -p internal/infrastructure/web/handlers
mkdir -p internal/infrastructure/web/routes
mkdir -p internal/infrastructure/web/middelware
mkdir -p internal/infrastructure/logger
mkdir -p internal/config
mkdir -p tests

# Criar arquivos raiz
touch .env
touch go.mod
touch go.sum
touch main.go
touch wire.go
touch wire_gen.go

# Criar arquivo cmd
touch cmd/root.go

# Criar DTOs
touch internal/application/dto/user.go
touch internal/application/dto/group.go

# Criar Mappers
touch internal/application/mappers/user_mapper.go
touch internal/application/mappers/group_mapper.go

# Criar interfaces dos usecases - User
touch internal/domain/interfaces/usecases/user/create_user_usecase.go
touch internal/domain/interfaces/usecases/user/get_user_usecase.go
touch internal/domain/interfaces/usecases/user/update_user_usecase.go
touch internal/domain/interfaces/usecases/user/delete_user_usecase.go
touch internal/domain/interfaces/usecases/user/list_users_usecase.go

# Criar interfaces dos usecases - Group
touch internal/domain/interfaces/usecases/group/create_group_usecase.go
touch internal/domain/interfaces/usecases/group/get_group_usecase.go
touch internal/domain/interfaces/usecases/group/update_group_usecase.go
touch internal/domain/interfaces/usecases/group/delete_group_usecase.go
touch internal/domain/interfaces/usecases/group/list_groups_usecase.go
touch internal/domain/interfaces/usecases/group/add_user_to_group_usecase.go
touch internal/domain/interfaces/usecases/group/remove_user_from_group_usecase.go

# Criar implementações dos usecases - User
touch internal/application/usecases/user/create_user_usecase.go
touch internal/application/usecases/user/get_user_usecase.go
touch internal/application/usecases/user/update_user_usecase.go
touch internal/application/usecases/user/delete_user_usecase.go
touch internal/application/usecases/user/list_users_usecase.go

# Criar implementações dos usecases - Group
touch internal/application/usecases/group/create_group_usecase.go
touch internal/application/usecases/group/get_group_usecase.go
touch internal/application/usecases/group/update_group_usecase.go
touch internal/application/usecases/group/delete_group_usecase.go
touch internal/application/usecases/group/list_groups_usecase.go
touch internal/application/usecases/group/add_user_to_group_usecase.go
touch internal/application/usecases/group/remove_user_from_group_usecase.go

# Criar entidades do domínio
touch internal/domain/entities/user.go
touch internal/domain/entities/group.go

# Criar interfaces de repositórios
touch internal/domain/interfaces/repositories/user_repository.go
touch internal/domain/interfaces/repositories/group_repository.go

# Criar infraestrutura de banco
touch internal/infrastructure/database/mongodb.go

# Criar implementações de repositórios
touch internal/infrastructure/repository/user_repository.go
touch internal/infrastructure/repository/group_repository.go

# Criar logger
touch internal/infrastructure/logger/logger.go

# Criar handlers web
touch internal/infrastructure/web/handlers/user_handler.go
touch internal/infrastructure/web/handlers/group_handler.go

# Criar rotas
touch internal/infrastructure/web/routes/routes.go

# Criar servidor web
touch internal/infrastructure/web/server.go

# Criar configurações
touch internal/config/config.go

echo "Estrutura completa criada com sucesso!"
echo ""
echo "Estrutura criada:"
echo "user-management/"
echo "├── .env"
echo "├── go.mod"
echo "├── go.sum"
echo "├── main.go"
echo "├── wire.go"
echo "├── wire_gen.go"
echo "├── cmd/"
echo "│   └── root.go"
echo "├── internal/"
echo "│   ├── application/"
echo "│   │   ├── dto/"
echo "│   │   │   ├── user.go"
echo "│   │   │   └── group.go"
echo "│   │   ├── mappers/"
echo "│   │   │   ├── user_mapper.go"
echo "│   │   │   └── group_mapper.go"
echo "│   │   └── usecases/"
echo "│   │       ├── user/"
echo "│   │       │   ├── create_user_usecase.go"
echo "│   │       │   ├── get_user_usecase.go"
echo "│   │       │   ├── update_user_usecase.go"
echo "│   │       │   ├── delete_user_usecase.go"
echo "│   │       │   └── list_users_usecase.go"
echo "│   │       └── group/"
echo "│   │           ├── create_group_usecase.go"
echo "│   │           ├── get_group_usecase.go"
echo "│   │           ├── update_group_usecase.go"
echo "│   │           ├── delete_group_usecase.go"
echo "│   │           ├── list_groups_usecase.go"
echo "│   │           ├── add_user_to_group_usecase.go"
echo "│   │           └── remove_user_from_group_usecase.go"
echo "│   ├── domain/"
echo "│   │   ├── entities/"
echo "│   │   │   ├── user.go"
echo "│   │   │   └── group.go"
echo "│   │   └── interfaces/"
echo "│   │       ├── repositories/"
echo "│   │       │   ├── user_repository.go"
echo "│   │       │   └── group_repository.go"
echo "│   ├── infrastructure/"
echo "│   │   ├── database/"
echo "│   │   │   └── mongodb.go"
echo "│   │   ├── repositories/"
echo "│   │   │   ├── user_repository.go"
echo "│   │   │   └── group_repository.go"
echo "│   │   └── web/"
echo "│   │       ├── handlers/"
echo "│   │       │   ├── user_handler.go"
echo "│   │       │   └── group_handler.go"
echo "│   │       ├── routes/"
echo "│   │       │   └── routes.go"
echo "│   │       ├── middleware/"
echo "│   │       └── server.go"
echo "│   └── config/"
echo "│       └── config.go"
echo "│   └── tests/"
echo ""
echo "Próximos passos:"
echo "1. Implementar as interfaces dos usecases"
echo "2. Implementar os usecases concretos"
echo "3. Implementar entidades do domínio"
echo "4. Implementar repositórios"
echo "5. Implementar handlers"
echo "6. Configurar rotas"
echo "7. Configurar injeção de dependências (wire.go)"
echo "8. Configurar go.mod com dependências necessárias"