#!/bin/bash

echo "╔════════════════════════════════════════╗"
echo "║   青羽写作平台 - 测试数据更新工具    ║"
echo "╚════════════════════════════════════════╝"
echo ""

cd "$(dirname "$0")/.." || exit 1
echo "当前目录: $(pwd)"
echo ""

echo "正在编译并运行数据更新工具..."
echo ""

go run cmd/seed_data/main.go

if [ $? -ne 0 ]; then
    echo ""
    echo "❌ 运行失败"
    echo ""
    exit 1
fi

echo ""
echo "✓ 运行完成"
echo ""
