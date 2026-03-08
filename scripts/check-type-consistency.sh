#!/bin/bash
# 类型一致性检查脚本
# 用途：检测重复的类型定义，确保类型定义的一致性

echo "======================================"
echo "  Qingyu 类型一致性检查"
echo "======================================"
echo ""

EXIT_CODE=0

# 1. 检查重复的枚举类型定义
echo "[1] 检查重复的枚举类型定义..."
ENUM_TYPES=$(grep -r "type.*Status.*string\|type.*Type.*string\|type.*State.*string" --include="*.go" internal/ models/ 2>/dev/null | grep -v "// " | cut -d: -f1 | sort | uniq -d)

if [ -n "$ENUM_TYPES" ]; then
    echo "⚠️  发现重复的枚举类型定义："
    echo "$ENUM_TYPES"
    grep -rn "type.*Status.*string" --include="*.go" internal/ models/ 2>/dev/null | grep -v "// "
    EXIT_CODE=1
else
    echo "✓ 未发现重复的枚举类型定义"
fi
echo ""

# 2. 检查BookStatus定义
echo "[2] 检查BookStatus定义..."
BOOKSTATUS_COUNT=$(grep -r "type BookStatus" --include="*.go" . 2>/dev/null | grep -v "// " | wc -l)

if [ "$BOOKSTATUS_COUNT" -gt 1 ]; then
    echo "⚠️  发现多处BookStatus定义（应为1处）："
    grep -rn "type BookStatus" --include="*.go" . 2>/dev/null | grep -v "// "
    EXIT_CODE=1
else
    echo "✓ BookStatus定义唯一"
fi
echo ""

# 3. 检查CategoryIDs类型
echo "[3] 检查CategoryIDs类型定义..."
CATEGORY_IDS_DEFS=$(grep -rn "CategoryIDs.*\[\]" --include="*.go" models/ 2>/dev/null)

if [ -n "$CATEGORY_IDS_DEFS" ]; then
    echo "ℹ️  CategoryIDs定义："
    echo "$CATEGORY_IDS_DEFS"

    # 检查是否有不一致的类型
    STRING_COUNT=$(echo "$CATEGORY_IDS_DEFS" | grep "\[\]string" | wc -l)
    OBJECTID_COUNT=$(echo "$CATEGORY_IDS_DEFS" | grep "ObjectID" | wc -l)

    if [ "$STRING_COUNT" -gt 0 ] && [ "$OBJECTID_COUNT" -gt 0 ]; then
        echo "⚠️  发现CategoryIDs类型不一致（既有string又有ObjectID）"
        EXIT_CODE=1
    else
        echo "✓ CategoryIDs类型一致"
    fi
else
    echo "ℹ️  未找到CategoryIDs定义"
fi
echo ""

# 4. 检查ID字段类型
echo "[4] 检查ID字段类型一致性..."
echo "ℹ️  当前ID字段使用情况："
grep -rn "ID.*string.*\`bson" --include="*.go" models/ | head -5
echo "..."
grep -rn "ID.*primitive.ObjectID.*\`bson" --include="*.go" models/ | head -5
echo ""

# 5. 检查domain层是否作为类型定义源
echo "[5] 检查domain层类型定义..."
if [ -d "internal/domain" ]; then
    DOMAIN_TYPES=$(find internal/domain -name "*.go" -exec grep -l "^type.*enum\|^type.*const" {} \; 2>/dev/null)
    if [ -n "$DOMAIN_TYPES" ]; then
        echo "✓ Domain层发现的类型定义："
        echo "$DOMAIN_TYPES"
    else
        echo "⚠️  Domain层未找到类型定义文件"
    fi
else
    echo "⚠️  未找到internal/domain目录"
fi
echo ""

# 6. 检查models层是否定义了枚举（应该引用domain层）
echo "[6] 检查models层是否有独立枚举定义..."
MODEL_ENUMS=$(grep -rn "^const (" --include="*.go" models/ 2>/dev/null | grep -A 10 "Status\|Type\|State" | grep "string" | wc -l)

if [ "$MODEL_ENUMS" -gt 0 ]; then
    echo "⚠️  models层发现枚举定义，应该引用domain层定义："
    grep -rn "^const (" --include="*.go" models/ 2>/dev/null | head -20
    EXIT_CODE=1
else
    echo "✓ models层未发现独立枚举定义"
fi
echo ""

# 总结
echo "======================================"
if [ $EXIT_CODE -eq 0 ]; then
    echo "✓ 类型一致性检查通过"
else
    echo "⚠️  发现类型一致性问题，请修复后再提交"
fi
echo "======================================"

exit $EXIT_CODE
