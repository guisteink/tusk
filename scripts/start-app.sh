#!/bin/bash

CURRENT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PARENT_DIR="$(dirname "$CURRENT_DIR")"

# Carregar variáveis de ambiente do arquivo .env
while IFS= read -r line; do
  export "$line"
done < "${PARENT_DIR}/.env"

# Executar o aplicativo
go run -gcflags "all=-N -l" "${CURRENT_DIR}/../cmd/main.go"
