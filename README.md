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


fastquest-backend/
├── database/
│ └── database.go # Configuração da conexão com o banco
│
├── docs/
│ ├── docs.go # Arquivo gerado pelo Swagger
│ ├── swagger.json # Documentação OpenAPI
│ └── swagger.yaml
│
├── handlers/


</div> ```
