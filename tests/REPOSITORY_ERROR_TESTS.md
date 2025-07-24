# Testes de Erro de Repository

Este documento descreve os testes de erro de repository implementados para validar o comportamento da aplicação quando os repositórios falham.

## Arquivos de Teste

### `repository_error_test.go`
Este arquivo contém testes unitários que simulam erros nos repositórios usando mocks. Os testes verificam se os usecases propagam corretamente os erros dos repositórios.

### `integration_repository_error_test.go`
Este arquivo contém testes de integração que simulam erros de repository dentro do contexto da suite de integração existente.

## Cenários de Teste Implementados

### Testes de User Repository

1. **TestCreateUserRepositoryError**
   - Simula erro durante criação de usuário
   - Verifica propagação do erro do repository

2. **TestGetUserByIDRepositoryError**
   - Simula erro durante busca de usuário por ID
   - Verifica tratamento de usuário não encontrado

3. **TestListUsersRepositoryError**
   - Simula erro durante listagem de usuários
   - Verifica comportamento com falha na query

4. **TestUpdateUserRepositoryError**
   - Simula erro durante atualização de usuário
   - Testa cenário onde GetByID funciona mas Update falha

5. **TestDeleteUserRepositoryError**
   - Simula erro durante exclusão de usuário
   - Verifica propagação do erro de delete

### Testes de Group Repository

1. **TestCreateGroupRepositoryError**
   - Simula erro durante criação de grupo
   - Verifica propagação do erro do repository

2. **TestGetGroupByIDRepositoryError**
   - Simula erro durante busca de grupo por ID
   - Verifica tratamento de grupo não encontrado

3. **TestListGroupsRepositoryError**
   - Simula erro durante listagem de grupos
   - Verifica comportamento com falha na query

4. **TestUpdateGroupRepositoryError**
   - Simula erro during atualização de grupo
   - Testa cenário onde GetByID funciona mas Update falha

5. **TestDeleteGroupRepositoryError**
   - Simula erro durante exclusão de grupo
   - Verifica propagação do erro de delete

6. **TestAddUserToGroupRepositoryError**
   - Simula erro ao adicionar usuário ao grupo
   - Testa cenário onde validações passam mas operação falha

7. **TestRemoveUserFromGroupRepositoryError**
   - Simula erro ao remover usuário do grupo
   - Verifica propagação do erro de remoção

### Testes de Cenários Específicos

1. **TestMultipleRepositoryErrors**
   - Testa múltiplos repositórios falhando simultaneamente
   - Verifica independência dos erros

2. **TestDatabaseTimeoutError**
   - Simula timeout do banco de dados
   - Verifica tratamento de timeouts

3. **TestDatabaseConnectionError**
   - Simula falha de conexão com banco
   - Verifica tratamento de erros de conectividade

4. **TestEntityNotFoundRepositoryError**
   - Simula erro de entidade não encontrada
   - Verifica mensagens de erro apropriadas

5. **TestDatabaseConstraintViolationError**
   - Simula violação de constraints (ex: email duplicado)
   - Verifica tratamento de erros de integridade

## Status dos Testes

✅ **Todos os testes estão passando com sucesso!**

Resultado da execução mais recente:
```
PASS: TestIntegrationSuite (1.48s)
PASS: TestRepositoryErrorSuite (0.00s)
PASS
ok      user-management/tests   1.502s
```

### Execução dos Testes Corrigidos

Os testes que estavam falhando foram corrigidos:

1. **TestAddUserToGroupWithRepositoryError** ✅
   - Problema: Mock não configurado para GetByID
   - Solução: Configurado GetByID para user e group antes do AddUserToGroup

2. **TestUserUpdateWithRepositoryError** ✅  
   - Problema: Mock não configurado para GetByID
   - Solução: Configurado GetByID com sucesso antes do Update com erro

## Como Executar os Testes

### Executar apenas os testes de erro de repository:
```bash
go test ./tests -v -run TestRepositoryErrorSuite
```

### Executar todos os testes:
```bash
go test ./tests -v
```

### Executar com timeout específico:
```bash
go test ./tests -v -run TestRepositoryErrorSuite -timeout 30s
```

## Estrutura dos Mocks

Os testes utilizam mocks das interfaces `IUserRepository` e `IGroupRepository` implementados com a biblioteca `testify/mock`. Os mocks permitem:

- Configurar retornos específicos para métodos
- Simular erros de diferentes tipos
- Verificar se os métodos foram chamados corretamente
- Validar os parâmetros passados

## Tipos de Erro Testados

1. **Erros de Conexão**: Simulam problemas de conectividade com o banco
2. **Erros de Timeout**: Simulam operações que demoram demais
3. **Erros de Entidade Não Encontrada**: Simulam buscas que não retornam dados
4. **Erros de Constraint**: Simulam violações de integridade do banco
5. **Erros Gerais de Repository**: Simulam falhas genéricas nas operações

## Cobertura de Teste

Os testes cobrem todos os métodos das interfaces de repository:
- `Create()`
- `GetByID()`
- `List()`
- `Update()`
- `Delete()`
- `AddUserToGroup()` (apenas GroupRepository)
- `RemoveUserFromGroup()` (apenas GroupRepository)

## Benefícios

1. **Validação de Error Handling**: Garante que erros são tratados apropriadamente
2. **Cobertura de Cenários Críticos**: Testa falhas que podem ocorrer em produção
3. **Documentação de Comportamento**: Serve como documentação de como a aplicação deve se comportar em caso de erro
4. **Prevenção de Regressões**: Evita que mudanças quebrem o tratamento de erros
5. **Confiabilidade**: Aumenta a confiança na robustez da aplicação
