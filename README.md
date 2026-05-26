# FastQuest — Backend API

Backend da plataforma **FastQuest**, um sistema de simulados e banco de questões para preparação de estudantes de Direito para o **Exame da Ordem dos Advogados do Brasil (OAB)**.

A API é responsável por processar as regras de negócio da plataforma, gerenciar questões, listas de exercícios, respostas dos usuários e estatísticas de desempenho.

---

# Como Rodar o Projeto

Este projeto utiliza um **Makefile** para automatizar as tarefas de configuração, migração de banco de dados e inicialização do servidor. Abaixo estão as instruções detalhadas para colocar a aplicação para rodar na sua máquina.

---

## Pré-requisitos

Antes de começar, certifique-se de ter instalado em sua máquina:
* **Go** (versão 1.20 ou superior)
* **Make** (nativo no Linux/Mac; necessário instalar no Windows)
* Um banco de dados ativo (PostgreSQL, conforme configurado)

---

## Passo a Passo para Configuração

### 1. Clonar o Repositório
Abra o terminal na pasta onde deseja salvar o projeto e execute:

```bash
git clone [https://github.com/FastQuest/backend.git](https://github.com/FastQuest/backend.git)
cd backend
```

### 2. Configurar as Variáveis de Ambiente

O projeto exige um arquivo `.env` na raiz para se conectar ao banco de dados e carregar outras configurações.

Crie um arquivo chamado `.env` na raiz do projeto e adicione as suas credenciais seguindo o arquivo `.env.example`. Exemplo:

```env
DB_NAME=postgres
DB_HOST=localhost
DB_PASSWORD=password
DB_USER=postgres
DB_PORT=5432

GOOSE_DRIVER=postgres
GOOSE_DBSTRING=postgres://postgres:password@localhost:5432/postgres
GOOSE_MIGRATION_DIR=./migrations

GEMINI_API_KEY=api_pass

```

### 3. Instalar as Ferramentas Globais (`swag` e `goose`)

O Makefile possui um comando dedicado para instalar o gerador de documentação (**Swagger**) e o gerenciador de migrações (**Goose**) diretamente no seu ambiente Go:

```bash
make setup
```

---

## Migrações do Banco de Dados

Antes de subir o servidor pela primeira vez (ou sempre que houver novas tabelas), execute as migrações para estruturar o banco de dados de forma automatizada:

```bash
make db-up
```

> 🛑 *Nota: Se o seu arquivo `.env` não for encontrado, o comando será abortado imediatamente com um aviso amigável.*

---

## Inicializando a Aplicação

Para compilar, gerar a documentação e rodar o servidor de desenvolvimento, basta executar:

```bash
make run
```

💡 **O que o `make run` faz por debaixo dos panos?**

1. Valida se o **Go** está instalado no sistema.
2. Valida se o arquivo **`.env`** está presente.
3. Verifica se a pasta `docs/` (gerada pelo Swagger) existe. Se não existir, ele executa o `swag init` automaticamente para você.
4. Roda o comando nativo `go run .`.

A documentação interativa da API estará disponível em: `http://localhost:8080/swagger/index.html` (ajuste a porta conforme configurado no seu `.env`).

---

## Resumo de Comandos Disponíveis

| Comando | Descrição |
| --- | --- |
| `make setup` | Instala o `swag` e o `goose` na pasta de binários do Go (`GOPATH`). |
| `make db-up` | Executa todas as migrações SQL pendentes na pasta `/migrations`. |
| `make run` | Executa as validações, gera o Swagger (se necessário) e inicia a API. |

---

# Tecnologias

| Tecnologia | Finalidade |
| :--- | :--- |
| Go (Golang) | Linguagem principal do backend |
| Gorilla Mux | Roteamento HTTP da API |
| GORM | ORM para interação com banco de dados |
| PostgreSQL | Banco de dados relacional |
| Swagger / OpenAPI | Documentação automática da API |
| Goose | Ferramenta de gerenciamento de migrações de banco de dados. |
| JSON REST API | Comunicação entre frontend e backend |

---

# Funcionalidades

- **Questões**
  - Criar novas questões
  - Listar questões
  - Buscar questão por ID
  - Remover questões

- **Respostas**
  - Envio de respostas dos usuários
  - Consulta de respostas por questão

- **Listas de Questões**
  - Criar listas personalizadas
  - Associar questões às listas
  - Visualizar listas de exercícios

- **Simulados**
  - Responder listas completas como simulados
  - Registro de histórico de respostas

- **Estatísticas**
  - Taxa de acertos por disciplina
  - Histórico de desempenho do usuário

- **Integração com Frontend**
  - API REST consumida pelo frontend em **Vue.js**

---

# Estrutura do Projeto

```text
fastquest-backend/
├── docs/
│   ├── docs.go              # Arquivo gerado pelo Swagger
│   ├── swagger.json         # Documentação OpenAPI
│   └── swagger.yaml
│
├── internal/
│   ├── platform/
│   │   └── database/        # Conexão com banco (GORM)
│   ├── question/            # Contexto de questões (handler/service/repository/dto/model)
│   ├── answer/              # Contexto de respostas (handler/service/repository/dto/model)
│   ├── questionset/         # Contexto de listas (handler/service/repository/dto/model)
│   ├── source/              # Contexto de fontes (handler/service/repository/dto/model)
│   ├── exam/                # Contexto de simulados/exam (handler/service/repository/dto/model)
│   └── ai/                  # Contexto de geração por IA (handler/service/repository/dto/model)
│
├── migrations/              # Scripts de migração do banco
│
├── pkg/
│   ├── models/
│   │   ├── answers.go
│   │   ├── comment.go
│   │   ├── pagination.go
│   │   ├── question.go
│   │   ├── question_set.go
│   │   ├── questionSource.go
│   │   ├── source.go
│   │   ├── subject.go
│   │   ├── topic.go
│   │   └── user.go
│   │
│   └── filtersMap.go        # Mapeamento de filtros para queries
│
├── router.go                # Definição das rotas da API
├── main.go                  # Ponto de entrada do servidor
├── go.mod                   # Dependências do projeto
└── go.sum                   # Checksum das dependências
```

---

# DataBase Scheme

<img width="1008" height="769" alt="image" src="https://github.com/user-attachments/assets/a41e68f4-5ce0-4e41-84bb-d160eb7e75e1" />

---

# Arquitetura

A arquitetura do FastQuest segue um modelo **API REST com separação entre frontend e backend**.

Frontend (Vue.js + Typescript)
↓
API REST (Go + Mux + GORM)
↓
PostgreSQL Database

---

### Fluxo de funcionamento

1. O usuário interage com o frontend
2. O frontend envia requisições HTTP para a API
3. A API processa as regras de negócio
4. O banco de dados armazena ou recupera informações
5. A resposta é retornada ao frontend em formato JSON

---

# Banco de Dados

O sistema utiliza **PostgreSQL** com modelo relacional.

### Entidades principais

| Entidade | Descrição |
| :--- | :--- |
| User | Usuários cadastrados na plataforma |
| Subject | Disciplinas da prova da OAB |
| Topic | Tópicos dentro das disciplinas |
| Question | Questões cadastradas |
| Answer | Alternativas de resposta |
| Question_Set | Listas de questões |
| Source | Fonte da questão (ex: prova específica) |
| Comment | Comentários de usuários |
| User_Response | Histórico de respostas |

### Relacionamentos importantes

- `Question → Answer` (1:N)
- `Question → Topic` (N:N)
- `Question → Source` (N:N)
- `Question_Set → Question` (N:N)
- `User → User_Response` (1:N)

---

# Endpoints principais 

## Questões

| Método | Endpoint           | Descrição                        |
| :----- | :----------------- | :------------------------------- |
| POST   | `/questions`       | Criar nova questão               |
| GET    | `/questions`       | Listar questões                  |
| POST   | `/questions/array` | Buscar questões por lista de IDs |
| GET    | `/questions/{id}`  | Buscar questão específica        |
| DELETE | `/questions/{id}`  | Remover questão                  |

## Respostas

| Método | Endpoint                  | Descrição                   |
| :----- | :------------------------ | :-------------------------- |
| POST   | `/questions/{id}/answers` | Enviar resposta             |
| GET    | `/questions/{id}/answers` | Listar respostas da questão |

## Listas de questões

| Método | Endpoint                           | Descrição         |
| :----- | :--------------------------------- | :---------------- |
| POST   | `/question-sets`                   | Criar nova lista  |
| GET    | `/question-sets`                   | Listar listas     |
| GET    | `/question-sets/{id}`              | Detalhe da lista  |
| GET    | `/question-sets/{id}/questions`    | Questões da lista |
| GET    | `/question-sets/{id}/question-ids` | IDs das questões  |

---

# Roadmap

| Fase   | Status | Descrição                          |
| :----- | :----: | :--------------------------------- |
| Fase 1 |    ✅   | Estrutura inicial do backend       |
| Fase 2 |    ✅   | Banco de questões                  |
| Fase 3 |   🚀   | Sistema de simulados               |
| Fase 4 |   📊   | Estatísticas de desempenho         |
| Fase 5 |   🔍   | Sistema de busca avançada          |
| Fase 6 |   🧠   | Recomendações de estudo            |
| Fase 7 |   🤖   | Integração com IA para explicações |





