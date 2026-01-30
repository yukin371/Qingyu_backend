#!/bin/bash
set -e

echo "Fixing imports..."

# Remove unused shared imports
sed -i '/^[[:space:]]*"Qingyu_backend\/api\/v1\/shared"[[:space:]]*$/d' api/v1/recommendation/recommendation_api.go
sed -i '/^[[:space:]]*"Qingyu_backend\/api\/v1\/shared"[[:space:]]*$/d' api/v1/stats/reading_stats_api.go
sed -i '/^[[:space:]]*"Qingyu_backend\/api\/v1\/shared"[[:space:]]*$/d' api/v1/ai/writing_api.go
sed -i '/^[[:space:]]*"Qingyu_backend\/api\/v1\/shared"[[:space:]]*$/d' api/v1/ai/writing_assistant_api.go
sed -i '/^[[:space:]]*"net\/http"[[:space:]]*$/d' api/v1/ai/writing_api.go

# Add missing response import to admin files
for file in api/v1/admin/config_api.go api/v1/admin/permission_api.go; do
  if ! grep -q '"Qingyu_backend/pkg/response"' "$file"; then
    sed -i '/^[[:space:]]*"github\/com\/gin-gonic\/gin"[[:space:]]*$/a\t"Qingyu_backend/pkg/response"' "$file"
    echo "Added response import to $file"
  fi
done

echo "✓ Imports fixed"
