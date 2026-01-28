#!/bin/bash

# 批量更新Handler文件的响应函数脚本

FILE=$1

if [ -z "$FILE" ]; then
    echo "Usage: $0 <file_path>"
    exit 1
fi

echo "Processing: $FILE"

# 备份文件
cp "$FILE" "$FILE.bak"

# 1. 替换import
sed -i 's|"Qingyu_backend/api/v1/shared"|"Qingyu_backend/pkg/response"|g' "$FILE"

# 2. 替换Success (200 OK)
sed -i 's|shared\.Success(c, http\.StatusOK, \([^,]*\), \(.*\))|response.Success(c, \2)|g' "$FILE"

# 3. 替换Success (201 Created)
sed -i 's|shared\.Success(c, http\.StatusCreated, \([^,]*\), \(.*\))|response.Created(c, \2)|g' "$FILE"

# 4. 替换Error (400 BadRequest)
sed -i 's|shared\.Error(c, http\.StatusBadRequest, \([^,]*\), \(.*\))|response.BadRequest(c, \1, \2)|g' "$FILE"

# 5. 替换Error (401 Unauthorized)
sed -i 's|shared\.Error(c, http\.StatusUnauthorized, \([^,]*\), \(.*\))|response.Unauthorized(c, \1)|g' "$FILE"

# 6. 替换Error (404 NotFound)
sed -i 's|shared\.Error(c, http\.StatusNotFound, \([^,]*\), \(.*\))|response.NotFound(c, \1)|g' "$FILE"

# 7. 替换Error (500 InternalServerError)
sed -i 's|shared\.Error(c, http\.StatusInternalServerError, \([^,]*\), \(.*\))|response.InternalError(c, \2)|g' "$FILE"

# 8. 移除net/http import（如果不再使用）
sed -i '/^import ($/,/^)$/{ /"net\/http"/d }' "$FILE"

# 9. 修复InternalError调用（传递err而不是err.Error()）
sed -i 's|response\.InternalError(c, err\.Error())|response.InternalError(c, err)|g' "$FILE"
sed -i 's|response\.InternalError(c, errMsg)|response.InternalError(c, err)|g' "$FILE"

echo "Completed: $FILE"
