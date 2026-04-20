# Refactor para DDD (Domain-Driven Design)

## Objetivo

Este documento descreve o refactor realizado no backend para uma arquitetura orientada a domínio (DDD), **sem alterar a lógica de negócio existente, o fluxo das rotas, os contratos HTTP e os comportamentos já existentes**.

O foco foi reorganizar o código para melhorar:
- separação de responsabilidades,
- manutenção,
- testabilidade,
- legibilidade,
- evolução por domínio.

---

## Princípios aplicados

1. **Preservação de comportamento**
- Mesmos endpoints, métodos HTTP e payloads.
- Mesmas validações e códigos de status existentes.
- Mesma ordem de execução do fluxo principal.

2. **Separação em camadas por contexto**
Cada contexto em `internal/` foi dividido em:
- `handler.go`: camada HTTP (entrada/saída, parse JSON, status code)
- `service.go`: regras de negócio / orquestração
- `repository.go`: acesso a dados (GORM/queries)
- `dto.go`: estruturas de request/response específicas do contexto
- `model.go`: aliases/modelos de domínio quando aplicável

3. **Infraestrutura compartilhada centralizada**
- Conexão de banco consolidada em `internal/platform/database/database.go`.
- `database/database.go` mantido como camada de compatibilidade (delegate), evitando quebra imediata de imports antigos.

---

## Estrutura final

```text
internal/
├── ai/
│   ├── dto.go
│   ├── handler.go
│   ├── model.go
│   ├── repository.go
│   └── service.go
│
├── answer/
│   ├── dto.go
│   ├── handler.go
│   ├── model.go
│   ├── repository.go
│   └── service.go
│
├── exam/
│   ├── dto.go
│   ├── handler.go
│   ├── model.go
│   ├── repository.go
│   └── service.go
│
├── platform/
│   └── database/
│       └── database.go
│
├── question/
│   ├── dto.go
│   ├── handler.go
│   ├── model.go
│   ├── repository.go
│   └── service.go
│
├── questionset/
│   ├── dto.go
│   ├── handler.go
│   ├── model.go
│   ├── repository.go
│   └── service.go
│
└── source/
    ├── dto.go
    ├── handler.go
    ├── model.go
    ├── repository.go
    └── service.go
```

---

## Mapeamento de rotas (antes e depois)

As rotas **não foram alteradas**. Apenas os handlers foram redirecionados para os novos contextos.

### Question
- `POST /questions` -> `internal/question.CreateQuestion`
- `GET /questions` -> `internal/question.GetQuestions`
- `POST /questions/by-ids` -> `internal/question.GetQuestionsByArray`
- `GET /questions/{id}` -> `internal/question.GetQuestion`
- `DELETE /questions/{id}` -> `internal/question.DeleteQuestion`

### Answer
- `POST /questions/{id}/answers` -> `internal/answer.PostAnswers`
- `GET /questions/{id}/answers` -> `internal/answer.GetAnswers`
- `POST /answers/by-ids` -> `internal/answer.GetAnswersByIDArray`

### Question Set
- `POST /question-sets` -> `internal/questionset.CreateQuestionSet`
- `GET /question-sets` -> `internal/questionset.GetLists`
- `GET /question-sets/{id}` -> `internal/questionset.GetQuestionSet`
- `GET /question-sets/{id}/questions` -> `internal/questionset.GetQuestionsFromSet`

### Source
- `POST /sources` -> `internal/source.CreateSource`

### AI
- `POST /ai/gen-question` -> `internal/ai.PostAIGenQuestion`
- `POST /ai/gen-questionset` -> `internal/ai.PostAIGenQuestionSet`

### Exam
- `POST /exam` -> `internal/exam.CreateExam`

---

## Fluxo por camada (exemplo)

### Exemplo: `POST /questions`
1. **Handler (`question/handler.go`)**
   - Lê body,
   - valida formato (objeto único ou array),
   - valida campos obrigatórios,
   - envia para persistência.

2. **Repository (`question/repository.go`)**
   - Executa criação no banco via GORM.

3. **Service (`question/service.go`)**
   - Expõe funções de negócio reutilizáveis para criação em lote e validações internas.

4. **Response**
   - Retorna payload JSON idêntico ao comportamento anterior.

---

## Arquivos alterados no bootstrap da aplicação

### `main.go`
- Inicialização de IA passou a usar `internal/ai.InitGemini()`.
- Inicialização de banco continua via `database.InitDB()` (compatibilidade preservada).

### `router.go`
- Imports de `handlers` substituídos por imports dos contextos em `internal/...`.
- Paths e métodos HTTP mantidos sem alterações.

### `database/database.go`
- Agora delega para `internal/platform/database`.

---

## Remoção de código legado

O pacote antigo `handlers/` foi removido após o roteamento ser migrado para `internal/`.

Isso evita duplicidade de lógica e reduz risco de manutenção paralela em duas arquiteturas.

---

## Boas práticas de Go aplicadas

- Organização por domínio e responsabilidade clara.
- Funções menores e focadas.
- Evita acoplamento de HTTP com query SQL sempre que possível.
- Reaproveitamento de modelos existentes em `pkg/models` para não quebrar compatibilidade.
- Continuidade de uso de `gorm` com acesso centralizado ao DB.
- `gofmt` aplicado.

---

## Observações de compatibilidade

- O projeto permanece com comportamento funcional equivalente ao original.
- Não houve mudança intencional de semântica das regras de negócio.
- A base ficou preparada para evolução incremental, com possibilidade de:
  - introduzir interfaces de repositório por domínio,
  - aumentar cobertura de testes unitários por camada,
  - refinar contratos DTO sem tocar no core da regra.

---

## Como validar localmente

```bash
go test ./...
```

Resultado esperado atual:
- execução bem-sucedida da compilação dos pacotes,
- sem arquivos de teste (`[no test files]`) nos pacotes atuais.

---

## Próximos passos recomendados

1. Adicionar testes unitários por domínio (`service` e `repository`).
2. Introduzir interfaces para facilitar mocks em testes.
3. Separar melhor regras transversais (logging, observabilidade, erros padronizados).
4. Evoluir para um padrão de injeção de dependências (constructors explícitos) para reduzir dependência de singleton global.
5. Se desejado, criar contexto `user/` completo com o mesmo padrão (`handler/service/repository/dto/model`).

