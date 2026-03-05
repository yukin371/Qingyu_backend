#!/bin/bash
# Repository职责检查脚本
# 用途：检测Repository层是否承担了Service层的职责

echo "======================================"
echo "  Qingyu Repository职责检查"
echo "======================================"
echo ""

EXIT_CODE=0

# 查找所有Repository文件
REPO_FILES=$(find . -name "*_repository.go" -o -name "*repository*.go" | grep -v vendor)

if [ -z "$REPO_FILES" ]; then
    echo "ℹ️  未找到Repository文件"
    exit 0
fi

echo "检查以下Repository文件："
echo "$REPO_FILES"
echo ""

# 1. 检查可疑的方法名
echo "[1] 检查可疑的方法名..."
SUSPICIOUS_METHODS="Get.*With\|Calculate\|Process\|Transform\|Aggregate\|Compose"

for file in $REPO_FILES; do
    MATCHES=$(grep -E "func.*\($SUSPICIOUS_METHODS\)" "$file" 2>/dev/null)
    if [ -n "$MATCHES" ]; then
        echo "⚠️  $file 发现可疑方法："
        echo "$MATCHES"
        EXIT_CODE=1
    fi
done

if [ $EXIT_CODE -eq 0 ]; then
    echo "✓ 未发现可疑方法名"
fi
echo ""

# 2. 检查跨Collection操作
echo "[2] 检查跨Collection操作..."
for file in $REPO_FILES; do
    # 检查是否有多个collection变量的使用
    COLLECTION_COUNT=$(grep -o "\.Collection(" "$file" 2>/dev/null | wc -l)
    if [ "$COLLECTION_COUNT" -gt 1 ]; then
        echo "⚠️  $file 可能存在跨Collection操作（发现$COLLECTION_COUNT个Collection调用）"
        grep -n "\.Collection(" "$file" 2>/dev/null
        EXIT_CODE=1
    fi
done

if [ $EXIT_CODE -eq 0 ]; then
    echo "✓ 未发现跨Collection操作"
fi
echo ""

# 3. 检查复杂的数据转换
echo "[3] 检查复杂的数据转换..."
for file in $REPO_FILES; do
    # 检查是否有struct字面量构造（可能是DTO转换）
    DTO_PATTERNS="func.*toDTO\|func.*ToDTO\|&DTO\|&Dto"
    MATCHES=$(grep -E "$DTO_PATTERNS" "$file" 2>/dev/null)
    if [ -n "$MATCHES" ]; then
        echo "ℹ️  $file 发现DTO转换（应该考虑移到Service层）："
        grep -n "$DTO_PATTERNS" "$file" 2>/dev/null
    fi
done
echo ""

# 4. 检查业务逻辑关键字
echo "[4] 检查业务逻辑关键字..."
BUSINESS_KEYWORDS="if.*Status.*==.*Status\|validate\|check.*permission\|calculate.*price"

for file in $REPO_FILES; do
    MATCHES=$(grep -iE "$BUSINESS_KEYWORDS" "$file" 2>/dev/null)
    if [ -n "$MATCHES" ]; then
        echo "⚠️  $file 发现可能的业务逻辑："
        grep -inE "$BUSINESS_KEYWORDS" "$file" 2>/dev/null | head -5
        EXIT_CODE=1
    fi
done

if [ $EXIT_CODE -eq 0 ]; then
    echo "✓ 未发现明显业务逻辑"
fi
echo ""

# 5. 检查方法复杂度
echo "[5] 检查方法复杂度..."
for file in $REPO_FILES; do
    # 统计每个方法的行数
    while IFS= read -r line; do
        if [[ $line =~ ^func[[:space:]]+\([^)]+\)[[:space:]]+([A-Z].+) ]]; then
            FUNC_NAME="${BASH_REMATCH[1]}"
            # 计算到下一个func或文件结束的行数
            LINE_NUM=$(grep -n "^func" "$file" | grep "$line" | cut -d: -f1)
            NEXT_FUNC_LINE=$(grep -n "^func" "$file" | awk -F: -v curr=$LINE_NUM '$1 > curr {print $1; exit}')
            if [ -z "$NEXT_FUNC_LINE" ]; then
                NEXT_FUNC_LINE=$(wc -l < "$file")
            fi
            FUNC_LINES=$((NEXT_FUNC_LINE - LINE_NUM))

            if [ $FUNC_LINES -gt 50 ]; then
                echo "⚠️  $file:$LINE_NUM $FUNC_NAME 方法过长（${FUNC_LINES}行），建议拆分"
                EXIT_CODE=1
            fi
        fi
    done < <(grep -n "^func" "$file")
done

if [ $EXIT_CODE -eq 0 ]; then
    echo "✓ 方法复杂度检查通过"
fi
echo ""

# 总结
echo "======================================"
if [ $EXIT_CODE -eq 0 ]; then
    echo "✓ Repository职责检查通过"
else
    echo "⚠️  发现Repository职责问题，建议："
    echo "   1. 将业务逻辑移到Service层"
    echo "   2. 将DTO转换移到Service层或API层"
    echo "   3. Repository层只负责数据CRUD"
fi
echo "======================================"

exit $EXIT_CODE
