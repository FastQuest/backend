GOPATH := $(shell go env GOPATH)
SWAG := $(GOPATH)/bin/swag
GOOSE := $(GOPATH)/bin/goose

check-env:
	@if [ ! -f ".env" ]; then \
		echo "Erro: Arquivo '.env' não encontrado!"; \
		echo "Dica: Crie um arquivo '.env' na raiz do projeto com as credenciais do banco, goose e chave da api google."; \
		exit 1; \
	fi

setup:
	@echo "Instalando Swagger e Goose..."
	go install github.com/swaggo/swag/cmd/swag@latest
	go install github.com/pressly/goose/v3/cmd/goose@latest
	@echo "Ferramentas instaladas com sucesso!"

db-up: check-env
	@echo "Rodando migrações..."
	$(GOOSE) up

run: check-env
	@if [ ! -d "docs" ]; then \
		echo "⚠️  Aviso: Pasta 'docs' não encontrada. Gerando o Swagger automaticamente..."; \
		$(SWAG) init; \
	fi
	go run .