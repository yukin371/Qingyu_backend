#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"
cd "$REPO_ROOT"

SWAGGER_DIRS="api/v1,pkg/response,models,models/dto,service/interfaces,service/ai/dto,service/shared/storage,service/shared/stats"

echo "Generating Swagger artifacts..."
swag init \
  -g swagger.go \
  -d "$SWAGGER_DIRS" \
  --parseDependency=false \
  -o docs

echo
echo "Swagger artifacts updated:"
echo "  docs/docs.go"
echo "  docs/swagger.json"
echo "  docs/swagger.yaml"
