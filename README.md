# StudyFox API — Documentação

**Versão:** v1  
**Base URL:** `https://studyfoxapi.onrender.com/api/v1`  
**Repositório:** [github.com/salsapunk/StudyFoxAPI](https://github.com/salsapunk/StudyFoxAPI)  
**Linguagem:** Go · Framework: Gin · Banco de dados: PostgreSQL

---

## Visão Geral

A StudyFox API é uma REST API desenvolvida em Go para o projeto **StudyFox** (AV2 de Front-End II). Ela permite gerenciar usuários, matérias e tarefas acadêmicas, com autenticação via JWT armazenado em cookie.

### Tecnologias

| Tecnologia | Finalidade |
|---|---|
| Go + Gin | Framework Web HTTP |
| PostgreSQL + pgx | Banco de dados relacional |
| JWT (golang-jwt/jwt) | Autenticação |
| bcrypt | Hash de senhas |

---

## Autenticação

A API utiliza **JWT em cookie**. Após o login, o servidor define um cookie `Authorization` com validade de 30 dias. Todas as rotas protegidas exigem esse cookie.

```
Cookie: Authorization=<jwt_token>
```

O cookie é configurado com `HttpOnly`, `Secure` e `SameSite=None`.

---

## Formato de Resposta

Todas as respostas seguem o padrão:

**Sucesso:**
```json
{
  "success": true,
  "data": <conteúdo>
}
```

**Erro:**
```json
{
  "success": false,
  "error": {
    "message": "descrição do erro"
  }
}
```

---

## Modelos de Dados

### Usuário (`Usuario`)

| Campo | Tipo | JSON | Obrigatório |
|---|---|---|---|
| Matricula | int | `matricula_usuario` | X (gerado) |
| Email | string | `email` | ✅ |
| Senha | string | `senha` | ✅ |
| Tema | string | `tema` | X |

### Matéria (`Materia`)

| Campo | Tipo | JSON | Obrigatório |
|---|---|---|---|
| Codigo | int | `codigo_materia` | X (gerado) |
| Nome | string | `nome` | ✅ |
| Matricula | int | `matricula_usuario` | ✅ |

### Tarefa (`Tarefa`)

| Campo | Tipo | JSON | Obrigatório |
|---|---|---|---|
| Id | int | `id_tarefa` | X (gerado) |
| Nome | string | `nome` | ✅ |
| Prazo | `pgtype.Date` | `prazo` | X |
| Anotacao | string | `anotacao` | X |
| Codigo | int | `codigo_materia` | ✅ |
| Status | int | `status` | X |

> ⚠️ **Formato de `prazo`:** o campo é serializado como uma **string simples no formato `"YYYY-MM-DD"`** (ex.: `"2025-07-15"`), ou `null` quando não definido. O parser aceita **exclusivamente** esse formato — enviar um objeto ou um datetime completo (`2025-07-15T00:00:00Z`) causa erro de bind no body.

---

## Endpoints

---

### Autenticação & Usuário

#### `POST /signup`
Cadastra um novo usuário.

**Autenticação:** Não requerida

**Body:**
```json
{
  "email": "user@email.com",
  "senha": "minhasenha123"
}
```

**Resposta de sucesso (200):**
```json
{
  "success": true,
  "data": 42
}
```
`data` é a matrícula gerada para o novo usuário.

---

#### `POST /login`
Autentica o usuário e define o cookie JWT.

**Autenticação:** Não requerida

**Body:**
```json
{
  "email": "user@email.com",
  "senha": "minhasenha123"
}
```

**Resposta de sucesso (200):**
```json
{
  "success": true,
  "data": "Login executed successfully"
}
```
O cookie `Authorization` é definido automaticamente na resposta.

---

#### `GET /validate`
Valida o token JWT e retorna os dados do usuário autenticado.

**Autenticação:** ✅ Requerida

**Resposta de sucesso (200):**
```json
{
  "success": true,
  "data": {
    "matricula_usuario": 42,
    "email": "user@email.com",
    "senha": "",
    "tema": "dark"
  }
}
```

---

#### `PUT /email`
Altera o e-mail do usuário autenticado.

**Autenticação:** ✅ Requerida

**Body:**
```json
{
  "email": "novo@email.com"
}
```

**Resposta de sucesso (200):**
```json
{
  "success": true,
  "data": "Email do usuário atualizado"
}
```

---

#### `PUT /senha`
Altera a senha do usuário autenticado.

**Autenticação:** ✅ Requerida

**Body:**
```json
{
  "senha": "novasenha456"
}
```

**Resposta de sucesso (200):**
```json
{
  "success": true,
  "data": "Senha do usuário atualizada"
}
```

---

#### `PUT /tema`
Altera o tema (preferência visual) do usuário autenticado.

**Autenticação:** ✅ Requerida

**Body:**
```json
{
  "tema": "escuro"
}
```

**Resposta de sucesso (200):**
```json
{
  "success": true,
  "data": "Tema do usuário atualizada"
}
```

---

#### `DELETE /usuario`
Remove permanentemente a conta do usuário autenticado.

**Autenticação:** ✅ Requerida

**Resposta de sucesso (200):**
```json
{
  "success": true,
  "data": "Usuário deletado"
}
```

---

### Matérias

#### `POST /materia`
Cria uma nova matéria para o usuário autenticado.

**Autenticação:** ✅ Requerida

**Body:**
```json
{
  "nome": "Cálculo I"
}
```

**Resposta de sucesso (200):**
```json
{
  "success": true,
  "data": 7
}
```
`data` é o código gerado para a matéria.

---

#### `GET /materias`
Lista todas as matérias do usuário autenticado.

**Autenticação:** ✅ Requerida

**Resposta de sucesso (200):**
```json
{
  "success": true,
  "data": [
    {
      "codigo_materia": 7,
      "nome": "Cálculo I",
      "matricula_usuario": 42
    },
    {
      "codigo_materia": 8,
      "nome": "Estrutura de Dados",
      "matricula_usuario": 42
    }
  ]
}
```

---

#### `GET /materia/:codigo`
Retorna uma matéria específica pelo código.

**Autenticação:** ✅ Requerida

**Parâmetros de rota:**

| Parâmetro | Tipo | Descrição |
|---|---|---|
| `codigo` | int | Código da matéria |

**Resposta de sucesso (200):**
```json
{
  "success": true,
  "data": {
    "codigo_materia": 7,
    "nome": "Cálculo I",
    "matricula_usuario": 42
  }
}
```

---

#### `PUT /materia/:codigo`
Altera o nome de uma matéria.

**Autenticação:** ✅ Requerida

**Parâmetros de rota:**

| Parâmetro | Tipo | Descrição |
|---|---|---|
| `codigo` | int | Código da matéria |

**Body:**
```json
{
  "nome": "Cálculo II"
}
```

**Resposta de sucesso (200):**
```json
{
  "success": true,
  "data": "Nome da matéria atualizado"
}
```

---

#### `DELETE /materia/:codigo`
Remove uma matéria e todas as suas tarefas (via cascade no banco).

**Autenticação:** ✅ Requerida

**Parâmetros de rota:**

| Parâmetro | Tipo | Descrição |
|---|---|---|
| `codigo` | int | Código da matéria |

**Resposta de sucesso (200):**
```json
{
  "success": true,
  "data": "Matéria deletada"
}
```

---

### Tarefas

#### `POST /materia/:codigo/tarefa`
Cria uma nova tarefa em uma matéria.

**Autenticação:** ✅ Requerida

**Parâmetros de rota:**

| Parâmetro | Tipo | Descrição |
|---|---|---|
| `codigo` | int | Código da matéria |

**Body:**
```json
{
  "nome": "Lista 1",
  "anotacao": "Exercícios de limites",
  "prazo": "2025-07-15"
}
```

**Resposta de sucesso (200):**
```json
{
  "success": true,
  "data": 101
}
```
`data` é o ID gerado para a tarefa.

> ⚠️ **Nota de implementação:** ao contrário dos demais endpoints, este retorna **`502 Bad Gateway`** (em vez de `400`) quando o parâmetro `:codigo` não é um número válido ou o body é malformado. Isso é uma inconsistência conhecida no handler atual, não uma escolha documentada de design — trate `502` aqui como um sinal de requisição inválida, não de erro no servidor upstream.

---

#### `GET /materia/:codigo/tarefas`
Lista todas as tarefas de uma matéria.

**Autenticação:** ✅ Requerida

**Parâmetros de rota:**

| Parâmetro | Tipo | Descrição |
|---|---|---|
| `codigo` | int | Código da matéria |

**Resposta de sucesso (200):**
```json
{
  "success": true,
  "data": [
    {
      "id_tarefa": 101,
      "nome": "Lista 1",
      "anotacao": "Exercícios de limites",
      "prazo": "2025-07-15",
      "codigo_materia": 7,
      "status": 0
    }
  ]
}
```

---

#### `GET /materia/:codigo/tarefa/:id`
Retorna uma tarefa específica.

**Autenticação:** ✅ Requerida

**Parâmetros de rota:**

| Parâmetro | Tipo | Descrição |
|---|---|---|
| `codigo` | int | Código da matéria |
| `id` | int | ID da tarefa |

**Resposta de sucesso (200):**
```json
{
  "success": true,
  "data": {
    "id_tarefa": 101,
    "nome": "Lista 1",
    "anotacao": "Exercícios de limites",
    "prazo": "2025-07-15",
    "codigo_materia": 7,
    "status": 0
  }
}
```

---

#### `PUT /materia/:codigo/tarefa/:id/nome`
Altera o nome de uma tarefa.

**Autenticação:** ✅ Requerida

**Body:**
```json
{ "nome": "Lista 2" }
```

**Resposta de sucesso (200):**
```json
{ "success": true, "data": "Nome da tarefa atualizada" }
```

---

#### `PUT /materia/:codigo/tarefa/:id/prazo`
Altera o prazo de uma tarefa.

**Autenticação:** ✅ Requerida

**Body:**
```json
{ "prazo": "2025-08-01" }
```

**Resposta de sucesso (200):**
```json
{ "success": true, "data": "Prazo da tarefa atualizada" }
```

---

#### `PUT /materia/:codigo/tarefa/:id/anotacao`
Altera a anotação de uma tarefa.

**Autenticação:** ✅ Requerida

**Body:**
```json
{ "anotacao": "Revisar capítulo 3 antes de resolver" }
```

**Resposta de sucesso (200):**
```json
{ "success": true, "data": "Anotação da tarefa atualizada" }
```

---

#### `PUT /materia/:codigo/tarefa/:id/status`
Altera o status de uma tarefa.

**Autenticação:** ✅ Requerida

**Body:**
```json
{ "status": 1 }
```

**Resposta de sucesso (200):**
```json
{ "success": true, "data": "Status da tarefa atualizada" }
```

---

#### `DELETE /materia/:codigo/tarefa/:id`
Remove uma tarefa específica.

**Autenticação:** ✅ Requerida

**Parâmetros de rota:**

| Parâmetro | Tipo | Descrição |
|---|---|---|
| `codigo` | int | Código da matéria |
| `id` | int | ID da tarefa |

**Resposta de sucesso (200):**
```json
{ "success": true, "data": "Tarefa deletada" }
```

---

### Utilitários

#### `GET /ping`
Verifica se a API está online e o token é válido.

**Autenticação:** ✅ Requerida

**Resposta de sucesso (200):**
```json
{ "message": "pong" }
```

---

## Resumo de Endpoints

| Método | Rota | Auth | Descrição |
|---|---|---|---|
| POST | `/signup` | X | Cadastrar usuário |
| POST | `/login` | X | Login |
| GET | `/validate` | ✅ | Validar token |
| GET | `/ping` | ✅ | Health check |
| PUT | `/email` | ✅ | Alterar e-mail |
| PUT | `/senha` | ✅ | Alterar senha |
| PUT | `/tema` | ✅ | Alterar tema |
| DELETE | `/usuario` | ✅ | Deletar conta |
| POST | `/materia` | ✅ | Criar matéria |
| GET | `/materias` | ✅ | Listar matérias |
| GET | `/materia/:codigo` | ✅ | Ver matéria |
| PUT | `/materia/:codigo` | ✅ | Renomear matéria |
| DELETE | `/materia/:codigo` | ✅ | Deletar matéria |
| POST | `/materia/:codigo/tarefa` | ✅ | Criar tarefa |
| GET | `/materia/:codigo/tarefas` | ✅ | Listar tarefas |
| GET | `/materia/:codigo/tarefa/:id` | ✅ | Ver tarefa |
| PUT | `/materia/:codigo/tarefa/:id/nome` | ✅ | Renomear tarefa |
| PUT | `/materia/:codigo/tarefa/:id/prazo` | ✅ | Alterar prazo |
| PUT | `/materia/:codigo/tarefa/:id/anotacao` | ✅ | Alterar anotação |
| PUT | `/materia/:codigo/tarefa/:id/status` | ✅ | Alterar status |
| DELETE | `/materia/:codigo/tarefa/:id` | ✅ | Deletar tarefa |

---

## Códigos de Status HTTP

| Código | Situação |
|---|---|
| 200 | Operação realizada com sucesso |
| 400 | Dados inválidos ou ausentes no body |
| 401 | Token ausente, inválido ou expirado |
| 500 | Erro interno no servidor |
| 502 | *(apenas em `POST /materia/:codigo/tarefa`)* `:codigo` inválido ou body malformado — inconsistência conhecida, ver nota no endpoint |

---

## Configuração e Deploy

### Variáveis de Ambiente

| Variável | Descrição |
|---|---|
| `DATABASE_URL` | String de conexão PostgreSQL |
| `SECRET` | Chave secreta para assinar os JWTs |
| `PORT` | Porta em que o servidor vai escutar |
| `ALLOWED_ORIGIN` | Origens permitidas para CORS (separadas por vírgula) |

### Executar localmente

```bash
# Clone o repositório
git clone https://github.com/salsapunk/StudyFoxAPI.git
cd StudyFoxAPI

# Configure as variáveis de ambiente
export DATABASE_URL="postgres://user:password@localhost:5432/nomedobanco"
export SECRET="seu_segredo_aqui"
export PORT="8080"
export ALLOWED_ORIGIN="http://localhost:3000"

# Execute
go run ./cmd/api/
```

### Estrutura do Projeto

```
StudyFoxAPI/
├── cmd/
│   └── api/
│       └── main.go        # Entrypoint: rotas, CORS, pool de conexões
├── internal/
│   ├── handler/
│   │   └── handler.go     # Handlers HTTP (controllers)
│   ├── model/
│   │   └── model.go       # Modelos, queries SQL, helpers de resposta
│   ├── repository/        # Acesso ao banco de dados
│   └── service/           # Regras de negócio
├── go.mod
└── go.sum
```

---
