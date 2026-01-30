#!/bin/bash
set -e

echo "Starting batch replacement of shared calls..."

# List of files to process
files=(
  "api/v1/recommendation/recommendation_api.go"
  "api/v1/ai/writing_api.go"
  "api/v1/ai/writing_assistant_api.go"
  "api/v1/system/health_api.go"
)

for file in "${files[@]}"; do
  if [ ! -f "$file" ]; then
    echo "Skipping $file (not found)"
    continue
  fi
  
  echo "Processing $file..."
  
  # Create backup
  cp "$file" "${file}.backup"
  
  # Replace shared.Error calls
  sed -i 's/shared\.Error(c, http\.StatusUnauthorized, "未授权", "请先登录")/response.Unauthorized(c, "请先登录")/g' "$file"
  sed -i 's/shared\.Error(c, http\.StatusUnauthorized, "UNAUTHORIZED", "未授权访问")/response.Unauthorized(c, "未授权访问")/g' "$file"
  sed -i 's/shared\.Error(c, http\.StatusForbidden, "FORBIDDEN", "无权访问")/response.Forbidden(c, "无权访问")/g' "$file"
  
  # Replace shared.Success calls with gin.H
  sed -i 's/shared\.Success(c, http\.StatusCreated, \([^,]*\), gin\.H{/response.Created(c, gin.H{/g' "$file"
  sed -i 's/shared\.Success(c, http\.StatusOK, \([^,]*\), gin\.H{/response.SuccessWithMessage(c, \1, gin.H{/g' "$file"
  sed -i 's/shared\.Success(c, http\.StatusOk, \([^,]*\), gin\.H{/response.SuccessWithMessage(c, \1, gin.H{/g' "$file"
  sed -i 's/shared\.Success(c, http\.Status200, \([^,]*\), gin\.H{/response.SuccessWithMessage(c, \1, gin.H{/g' "$file"
  sed -i 's/shared\.Success(c, http\..StatusOK, \([^,]*\), gin\.H{/response.SuccessWithMessage(c, \1, gin.H{/g' "$file"
  
  # Replace shared.SuccessData calls
  sed -i 's/shared\.SuccessData(c, \([^)]*\))/response.Success(c, \1)/g' "$file"
  
  # Replace remaining shared.Success calls
  sed -i 's/shared\.Success(c, http\.StatusCreated, \([^,]*\), \([^)]*\))/response.Created(c, \2)/g' "$file"
  sed -i 's/shared\.Success(c, http\.StatusOK, \([^,]*\), \([^)]*\))/response.SuccessWithMessage(c, \1, \2)/g' "$file"
  sed -i 's/shared\.Success(c, http\.StatusOk, \([^,]*\), \([^)]*\))/response.SuccessWithMessage(c, \1, \2)/g' "$file"
  sed -i 's/shared\.Success(c, http\.Status200, \([^,]*\), \([^)]*\))/response.SuccessWithMessage(c, \1, \2)/g' "$file"
  sed -i 's/shared\.Success(c, http\.StatusOK, \([^,]*\), \([^)]*\))/response.SuccessWithMessage(c, \1, \2)/g' "$file"
  sed -i 's/shared\.Success(c, http\.Status200, \([^,]*\), \([^)]*\))/response.SuccessWithMessage(c, \1, \2)/g' "$file"
  sed -i 's/shared\.Success(c, http\.StatusOk, \([^,]*\), \([^)]*\))/response.SuccessWithMessage(c, \1, \2)/g' "$file"
  sed -i 's/shared\.Success(c, http\.Status200, \([^,]*\), \([^)]*\))/response.SuccessWithMessage(c, \1, \2)/g' "$file"
  
  echo "✓ Completed $file"
done

echo ""
echo "Batch replacement complete!"
echo "Please review and test the changes."
