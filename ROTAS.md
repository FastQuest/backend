# 📚 Documentação das Rotas - FlashQuest API

**Base URL:** `http://localhost:8080`

---

## 📋 Índice
1. [Questões (Questions)](#questões-questions)
2. [Respostas (Answers)](#respostas-answers)
3. [Listas de Exercícios (Question Sets)](#listas-de-exercícios-question-sets)
4. [Fontes de Exame (Sources)](#fontes-de-exame-sources)
5. [Exames (Exams)](#exames-exams)
6. [Inteligência Artificial (AI)](#inteligência-artificial-ai)

---

## 🎯 Questões (Questions)

### 1. Criar Questão(ões)
**`POST /questions`**

Cria uma ou múltiplas questões. Pode enviar um objeto único ou um array.

**Headers:**
```
Content-Type: application/json
```

**Body (Uma única questão):**
```json
{
  "statement": "Qual é a capital do Brasil?",
  "subject_id": 1,
  "user_id": 1,
  "source_exam_instance_id": 1,
  "answers": [
    {
      "text": "Brasília",
      "is_correct": true
    },
    {
      "text": "Rio de Janeiro",
      "is_correct": false
    },
    {
      "text": "São Paulo",
      "is_correct": false
    }
  ]
}
```

**Body (Array de questões - batch):**
```json
[
  {
    "statement": "Pergunta 1?",
    "subject_id": 1,
    "user_id": 1
  },
  {
    "statement": "Pergunta 2?",
    "subject_id": 2,
    "user_id": 2
  }
]
```

**Response (201 Created) - Questão única:**
```json
{
  "id": 11,
  "statement": "Qual é a capital do Brasil?",
  "subject_id": 1,
  "user_id": 1,
  "source_exam_instance_id": 1,
  "created_at": "2025-04-17T10:30:00Z",
  "updated_at": "2025-04-17T10:30:00Z",
  "subject": null,
  "user": null,
  "source": null,
  "answers": null
}
```

**Response (201 Created) - Array de questões:**
```json
[
  {
    "id": 11,
    "statement": "Pergunta 1?",
    "subject_id": 1,
    "user_id": 1,
    "created_at": "2025-04-17T10:30:00Z",
    "updated_at": "2025-04-17T10:30:00Z"
  },
  {
    "id": 12,
    "statement": "Pergunta 2?",
    "subject_id": 2,
    "user_id": 2,
    "created_at": "2025-04-17T10:31:00Z",
    "updated_at": "2025-04-17T10:31:00Z"
  }
]
```

---

### 2. Listar Questões (Com Filtros e Paginação)
**`GET /questions`**

Lista todas as questões com suporte a paginação, filtros e includes.

**Query Parameters:**
- `page` (default: 1) - Número da página
- `perPage` (default: 10, máximo: 100) - Itens por página
- `orderBy` (default: "created_at desc") - Ordenação
- `include` (comma-separated) - Inclui relacionamentos: "answers", "user", "subject", "source"
- `subject_id` - Filtro por disciplina
- `user_id` - Filtro por usuário

**Exemplos:**
```
GET http://localhost:8080/questions
GET http://localhost:8080/questions?page=1&perPage=5
GET http://localhost:8080/questions?include=answers,user&subject_id=1
GET http://localhost:8080/questions?orderBy=id asc&user_id=1
GET http://localhost:8080/questions?include=answers,user,subject,source&page=2&perPage=20
```

**Response (200 OK):**
```json
{
  "data": [
    {
      "id": 1,
      "created_at": "2025-04-16T15:20:00Z",
      "updated_at": "2025-04-16T15:20:00Z",
      "statement": "A Constituição Federal de 1988 é considerada uma constituição:",
      "subject": {
        "id": 1,
        "name": "Direito Constitucional"
      },
      "user": {
        "id": 1,
        "name": "João Silva",
        "email": "joao.silva@test.com"
      },
      "source": {
        "id": 1,
        "name": "OAB - Exame de Ordem Unificado",
        "type": "Exame",
        "metadata": {
          "year": 2024,
          "edition": 35,
          "phase": 1
        }
      },
      "answers": [
        {
          "id": 1,
          "text": "Rígida",
          "is_correct": true,
          "question_id": 1
        }
      ]
    }
  ],
  "pagination": {
    "page": 1,
    "per_page": 10,
    "total": 10,
    "total_pages": 1
  }
}
```

---

### 3. Buscar Questões por IDs
**`POST /questions/by-ids`**

Retorna múltiplas questões a partir de uma lista de IDs.

**Body:**
```json
{
  "ids": [1, 3, 5, 7]
}
```

**Response (200 OK):**
```json
[
  {
    "id": 1,
    "statement": "A Constituição Federal de 1988 é considerada uma constituição:",
    "subject_id": 1,
    "user_id": 1,
    "created_at": "2025-04-16T15:20:00Z",
    "updated_at": "2025-04-16T15:20:00Z"
  },
  {
    "id": 3,
    "statement": "O homicídio é punido com pena de:",
    "subject_id": 2,
    "user_id": 2,
    "created_at": "2025-04-16T15:21:00Z",
    "updated_at": "2025-04-16T15:21:00Z"
  }
]
```

---

### 4. Obter Questão Específica
**`GET /questions/{id}`**

Retorna uma questão específica com opção de incluir relacionamentos.

**Query Parameters:**
- `include` (comma-separated) - "answers", "user", "subject", "source"

**Exemplos:**
```
GET http://localhost:8080/questions/1
GET http://localhost:8080/questions/1?include=answers,user,subject,source
```

**Response (200 OK):**
```json
{
  "id": 1,
  "created_at": "2025-04-16T15:20:00Z",
  "updated_at": "2025-04-16T15:20:00Z",
  "statement": "A Constituição Federal de 1988 é considerada uma constituição:",
  "subject": {
    "id": 1,
    "name": "Direito Constitucional"
  },
  "user": {
    "id": 1,
    "name": "João Silva",
    "email": "joao.silva@test.com"
  },
  "source": {
    "id": 1,
    "name": "OAB - Exame de Ordem Unificado",
    "type": "Exame",
    "metadata": {
      "year": 2024,
      "edition": 35,
      "phase": 1
    }
  },
  "answers": [
    {
      "id": 1,
      "text": "Rígida",
      "is_correct": true,
      "question_id": 1
    },
    {
      "id": 2,
      "text": "Flexível",
      "is_correct": false,
      "question_id": 1
    }
  ]
}
```

---

### 5. Deletar Questão
**`DELETE /questions/{id}`**

Remove uma questão do banco de dados.

**Exemplo:**
```
DELETE http://localhost:8080/questions/11
```

**Response (204 No Content ou 200 OK):**
```json
{
  "message": "Question deleted successfully"
}
```

---

## 💬 Respostas (Answers)

### 1. Adicionar Respostas a uma Questão
**`POST /questions/{id}/answers`**

Cria múltiplas respostas para uma questão existente.

**Body:**
```json
[
  {
    "text": "Alternativa A",
    "is_correct": true
  },
  {
    "text": "Alternativa B",
    "is_correct": false
  },
  {
    "text": "Alternativa C",
    "is_correct": false
  },
  {
    "text": "Alternativa D",
    "is_correct": false
  }
]
```

**Response (201 Created):**
```json
{
  "message": "Answers created successfully",
  "count": 4,
  "ids": [41, 42, 43, 44]
}
```

---

### 2. Listar Respostas de uma Questão
**`GET /questions/{id}/answers`**

Retorna todas as respostas associadas a uma questão.

**Exemplo:**
```
GET http://localhost:8080/questions/1/answers
```

**Response (200 OK):**
```json
[
  {
    "id": 1,
    "text": "Rígida",
    "is_correct": true,
    "question_id": 1
  },
  {
    "id": 2,
    "text": "Flexível",
    "is_correct": false,
    "question_id": 1
  },
  {
    "id": 3,
    "text": "Semi-rígida",
    "is_correct": false,
    "question_id": 1
  },
  {
    "id": 4,
    "text": "Mutável",
    "is_correct": false,
    "question_id": 1
  }
]
```

---

### 3. Buscar Respostas por IDs
**`POST /answers/by-ids`**

Retorna múltiplas respostas a partir de uma lista de IDs.

**Body:**
```json
{
  "answer_ids": [1, 2, 5, 9]
}
```

**Response (200 OK):**
```json
[
  {
    "id": 1,
    "text": "Rígida",
    "is_correct": true,
    "question_id": 1
  },
  {
    "id": 2,
    "text": "Flexível",
    "is_correct": false,
    "question_id": 1
  },
  {
    "id": 5,
    "text": "Supremo Tribunal Federal (STF)",
    "is_correct": true,
    "question_id": 2
  }
]
```

---

## 📝 Listas de Exercícios (Question Sets)

### 1. Criar Lista de Exercícios
**`POST /question-sets`**

Cria uma nova lista de exercícios. Pode associar questões existentes ou criar novas.

**Body (Com questões existentes):**
```json
{
  "name": "OAB 2024 - Simulado 1",
  "type": "Simulado",
  "description": "Primeiro simulado completo para OAB 2024",
  "is_private": false,
  "user_id": 1,
  "questions": [1, 2, 3, 4, 5]
}
```

**Body (Criar com questões novas):**
```json
{
  "name": "Nova Lista com Questões",
  "type": "Lista de Exercícios",
  "description": "Cria questões e adiciona à lista em uma única operação",
  "is_private": true,
  "user_id": 2,
  "questions": [
    {
      "statement": "Pergunta criada inline 1?",
      "subject_id": 1,
      "user_id": 2
    },
    {
      "statement": "Pergunta criada inline 2?",
      "subject_id": 2,
      "user_id": 2
    }
  ]
}
```

**Response (201 Created):**
```json
{
  "id": 6,
  "name": "OAB 2024 - Simulado 1",
  "description": "Primeiro simulado completo para OAB 2024",
  "type": "Simulado",
  "user": null,
  "questions": null,
  "created_at": "2025-04-17T11:15:00Z",
  "is_private": false
}
```

---

### 2. Listar Listas de Exercícios (Paginado)
**`GET /question-sets`**

Lista todas as listas de exercícios com paginação e includes opcionais.

**Query Parameters:**
- `page` - Número da página
- `perPage` - Itens por página
- `include` - "user", "questions"

**Exemplos:**
```
GET http://localhost:8080/question-sets
GET http://localhost:8080/question-sets?page=1&perPage=5
GET http://localhost:8080/question-sets?page=1&perPage=5&include=questions
GET http://localhost:8080/question-sets?include=user,questions
```

**Response (200 OK):**
```json
[
  {
    "id": 1,
    "name": "Simulado OAB - Primeira Fase 2024",
    "description": "Questões de todas as disciplinas para primeira fase da OAB",
    "type": "Simulado",
    "user": {
      "id": 1,
      "name": "João Silva",
      "email": "joao.silva@test.com"
    },
    "questions": null,
    "created_at": "2025-04-16T16:00:00Z",
    "is_private": false
  },
  {
    "id": 2,
    "name": "Direito Constitucional - Fundamentals",
    "description": "Exercícios básicos de Direito Constitucional",
    "type": "Lista de Exercícios",
    "user": {
      "id": 2,
      "name": "Maria Santos",
      "email": "maria.santos@test.com"
    },
    "questions": null,
    "created_at": "2025-04-16T16:05:00Z",
    "is_private": true
  }
]
```

---

### 3. Obter Lista Específica com Questões
**`GET /question-sets/{id}`**

Retorna uma lista de exercícios completa com opção de incluir relacionamentos.

**Query Parameters:**
- `include` - "user", "questions"

**Exemplos:**
```
GET http://localhost:8080/question-sets/1
GET http://localhost:8080/question-sets/1?include=user,questions
```

**Response (200 OK):**
```json
{
  "id": 1,
  "name": "Simulado OAB - Primeira Fase 2024",
  "description": "Questões de todas as disciplinas para primeira fase da OAB",
  "type": "Simulado",
  "user": {
    "id": 1,
    "name": "João Silva",
    "email": "joao.silva@test.com"
  },
  "questions": [
    {
      "id": 1,
      "created_at": "2025-04-16T15:20:00Z",
      "updated_at": "2025-04-16T15:20:00Z",
      "statement": "A Constituição Federal de 1988 é considerada uma constituição:",
      "subject": null,
      "user": null,
      "source": null,
      "answers": null
    },
    {
      "id": 2,
      "created_at": "2025-04-16T15:21:00Z",
      "updated_at": "2025-04-16T15:21:00Z",
      "statement": "Qual é o órgão máximo do judiciário no Brasil?",
      "subject": null,
      "user": null,
      "source": null,
      "answers": null
    }
  ],
  "created_at": "2025-04-16T16:00:00Z",
  "is_private": false
}
```

---

### 4. Obter IDs das Questões de uma Lista
**`GET /question-sets/{id}/questions?fields=id`**

Retorna apenas os IDs das questões da lista.

**Exemplo:**
```
GET http://localhost:8080/question-sets/1/questions?fields=id
```

**Response (200 OK):**
```json
[1, 2, 3, 4, 5, 6, 7, 8, 9, 10]
```

---

### 5. Obter Questões Completas de uma Lista
**`GET /question-sets/{id}/questions`**

Retorna as questões completas da lista de exercícios.

**Exemplo:**
```
GET http://localhost:8080/question-sets/1/questions
```

**Response (200 OK):**
```json
[
  {
    "id": 1,
    "created_at": "2025-04-16T15:20:00Z",
    "updated_at": "2025-04-16T15:20:00Z",
    "statement": "A Constituição Federal de 1988 é considerada uma constituição:",
    "subject": null,
    "user": null,
    "source": null,
    "answers": null
  },
  {
    "id": 2,
    "created_at": "2025-04-16T15:21:00Z",
    "updated_at": "2025-04-16T15:21:00Z",
    "statement": "Qual é o órgão máximo do judiciário no Brasil?",
    "subject": null,
    "user": null,
    "source": null,
    "answers": null
  }
]
```

---

## 🗂️ Fontes de Exame (Sources)

### 1. Criar Fonte de Exame
**`POST /sources`**

Cria uma nova fonte de exame com sua instância.

**Body:**
```json
{
  "name": "CESPE 2024 - Magistratura",
  "type": "Banca Examinadora",
  "edition": 1,
  "phase": 1,
  "year": 2024
}
```

**Response (201 Created):**
```json
{
  "message": "Source Instance created successfully",
  "id": 6
}
```

---

## 🎓 Exames (Exams)

### 1. Criar Exame com Lista de Exercícios e Questões
**`POST /exam`**

Cria um exame completo com lista de exercícios e questões em uma única operação.

**Body:**
```json
{
  "exam": {
    "source_id": 1,
    "edition": 35,
    "phase": 2,
    "year": 2024
  },
  "list": {
    "name": "Exame Completo OAB 2024",
    "type": "Exame",
    "description": "Simulado com questões da OAB 2024",
    "is_private": false,
    "user_id": 1,
    "questions": [
      {
        "statement": "Qual é a base da República Federativa do Brasil?",
        "subject_id": 1,
        "user_id": 1,
        "answers": [
          {
            "text": "A soberania",
            "is_correct": true
          },
          {
            "text": "A democracia",
            "is_correct": false
          }
        ]
      }
    ]
  }
}
```

**Response (201 Created):**
```json
{
  "id": 3,
  "source_id": 1,
  "edition": 35,
  "phase": 2,
  "year": 2024,
  "created_at": "2025-04-17T12:00:00Z"
}
```

---

## 🤖 Inteligência Artificial (AI)

### 1. Gerar Questão via IA (Gemini)
**`POST /ai/gen-question`**

Gera uma questão utilizando a API do Gemini (requer GEMINI_API_KEY configurada).

**Importante:** A resposta é processada de forma assíncrona.

**Body:**
```json
{
  "text": "Gere uma questão sobre mandado de segurança no direito administrativo com 4 alternativas"
}
```

**Response:**
```
(Processamento assíncrono - não retorna resposta imediata)
```

---

### 2. Gerar Lista de Exercícios via IA
**`POST /ai/gen-questionset`**

Gera uma lista completa de exercícios utilizando a API do Gemini (requer GEMINI_API_KEY configurada).

**Importante:** A resposta é processada de forma assíncrona.

**Body:**
```json
{
  "text": "Crie um simulado com 10 questões sobre direito constitucional, incluindo tópicos como: constituição rígida, direitos fundamentais e organização dos poderes"
}
```

**Response:**
```
(Processamento assíncrono - não retorna resposta imediata)
```

---

## 🧪 Roteiro Recomendado de Testes

Execute os testes nesta ordem para garantir que tudo funcione corretamente:

### 1️⃣ Teste de Criação de Questão
```
POST http://localhost:8080/questions
```

### 2️⃣ Teste de Adição de Respostas
```
POST http://localhost:8080/questions/1/answers
```

### 3️⃣ Teste de Listagem de Questões
```
GET http://localhost:8080/questions
GET http://localhost:8080/questions?include=answers,user,subject,source
```

### 4️⃣ Teste de Busca por ID
```
GET http://localhost:8080/questions/1?include=answers
```

### 5️⃣ Teste de Busca por Array de IDs
```
POST http://localhost:8080/questions/by-ids
```

### 6️⃣ Teste de Obtenção de Respostas
```
GET http://localhost:8080/questions/1/answers
```

### 7️⃣ Teste de Criação de Lista
```
POST http://localhost:8080/question-sets
```

### 8️⃣ Teste de Listagem de Listas
```
GET http://localhost:8080/question-sets
GET http://localhost:8080/question-sets?include=user,questions
```

### 9️⃣ Teste de Obtenção de Lista Específica
```
GET http://localhost:8080/question-sets/1?include=user,questions
```

### 🔟 Teste de Questões da Lista
```
GET http://localhost:8080/question-sets/1/questions
GET http://localhost:8080/question-sets/1/questions?fields=id
```

---

## 💡 Dicas Importantes

1. **Includes:** Use `?include=answers,user,subject,source` para trazer dados relacionados
2. **Paginação:** O máximo de itens por página é 100
3. **Ordenação:** Use `?orderBy=id asc` ou `?orderBy=created_at desc`
4. **Filtros:** Use query params como `?subject_id=1` ou `?user_id=1`
5. **Array de IDs:** Quando buscar múltiplos registros, use POST com array de IDs
6. **Respostas em batch:** Pode criar múltiplas questões em uma única requisição
7. **IA:** As features de IA funcionam de forma assíncrona, não espere resposta imediata

---

**Última atualização:** 17 de abril de 2026
