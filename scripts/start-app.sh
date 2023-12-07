#!/bin/bash

CURRENT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PARENT_DIR="$(dirname "$CURRENT_DIR")"

source "${PARENT_DIR}/.env"

DATABASE_URI="${DATABASE_URI}" PORT="${PORT}" go run -gcflags "all=-N -l" "${CURRENT_DIR}/../cmd/app/main.go"
