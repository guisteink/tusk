#!/bin/bash

CURRENT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PARENT_DIR="$(dirname "$CURRENT_DIR")"
ENV_FILE="${PARENT_DIR}/.env"

# Verificar se o arquivo .env existe
if [ -f "$ENV_FILE" ]; then
  # Carregar variáveis de ambiente do arquivo .env
  while IFS= read -r line; do
    export "$line"
  done < "$ENV_FILE"
else
  echo "Aviso: O arquivo $ENV_FILE não foi encontrado. As variáveis de ambiente não serão carregadas."
fi

# Executar o aplicativo
go run -gcflags "all=-N -l" "${CURRENT_DIR}/../cmd/main.go"
