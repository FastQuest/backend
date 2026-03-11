# FastQuest — Backend API

Backend da plataforma **FastQuest**, um sistema de simulados e banco de questões para preparação de estudantes de Direito para o **Exame da Ordem dos Advogados do Brasil (OAB)**.

A API é responsável por processar as regras de negócio da plataforma, gerenciar questões, listas de exercícios, respostas dos usuários e estatísticas de desempenho.

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
├── database/
│   └── database.go          # Configuração da conexão com o banco
│
├── docs/
│   ├── docs.go              # Arquivo gerado pelo Swagger
│   ├── swagger.json         # Documentação OpenAPI
│   └── swagger.yaml
│
├── handlers/
│   ├── h_answers.go         # Endpoints de respostas
│   ├── h_questions.go       # Endpoints de questões
│   ├── h_question_set.go    # Endpoints de listas de questões
│   └── utils.go             # Funções auxiliares
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

# Instalação e Configuração

## Pré-requisitos

- Go 1.20+
- PostgreSQL
- Docker
- Git

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





