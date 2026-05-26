run:
	@if [ ! -d "docs" ]; then \
		echo "⚠️  Aviso: Pasta 'docs' não encontrada. Gerando o Swagger automaticamente..."; \
		swag init; \
	fi
	go run .