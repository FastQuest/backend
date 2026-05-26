GOPATH := $(shell go env GOPATH 2>/dev/null)
SWAG := $(GOPATH)/bin/swag
GOOSE := $(GOPATH)/bin/goose

check-go:
	@if ! command -v go > /dev/null 2>&1; then \
		echo "🛑 Erro: O comando 'go' não foi encontrado!"; \
		echo "💡 Dica: Instale o Go (https://go.dev/doc/install) e adicione-o ao seu PATH."; \
		exit 1; \
	fi

check-env:
	@if [ ! -f ".env" ]; then \
		echo "🛑 Erro: Arquivo '.env' não encontrado!"; \
		echo "💡 Dica: Crie um arquivo '.env' na raiz do projeto com as credenciais do banco."; \
		exit 1; \
	fi

setup: check-go
	@echo "Instalando Swagger e Goose..."
	go install github.com/swaggo/swag/cmd/swag@latest
	go install github.com/pressly/goose/v3/cmd/goose@latest
	@echo "Ferramentas instaladas com sucesso!"

db-up: check-env
	@echo "Rodando migrações..."
	$(GOOSE) up

run: check-go check-env
	@if [ ! -d "docs" ]; then \
		echo "⚠️  Aviso: Pasta 'docs' não encontrada. Gerando o Swagger automaticamente..."; \
		$(SWAG) init; \
	fi
	go run .