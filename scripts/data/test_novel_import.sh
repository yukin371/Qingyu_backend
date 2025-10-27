#!/bin/bash
# CNNovel125K 小说导入测试脚本 (Linux/Mac)
# 用于快速测试小说导入功能

set -e  # 遇到错误立即退出

echo "========================================"
echo "CNNovel125K 小说导入测试"
echo "========================================"
echo ""

# 检查 Python 环境
echo "[1/5] 检查 Python 环境..."
if ! command -v python3 &> /dev/null; then
    echo "错误: 未找到 Python3，请先安装 Python 3.7+"
    exit 1
fi
echo "✓ Python 环境正常"
echo ""

# 检查 datasets 库
echo "[2/5] 检查 Python 依赖..."
if ! python3 -c "import datasets" &> /dev/null; then
    echo "正在安装 datasets 库..."
    pip3 install datasets
fi
echo "✓ Python 依赖已安装"
echo ""

# 创建数据目录
mkdir -p data

# 运行 Python 脚本加载数据
echo "[3/5] 从 Hugging Face 加载小说数据..."
echo "提示: 首次运行会下载数据集，需要一些时间"
echo ""
python3 scripts/import_novels.py --max-novels 100 --output data/novels_test.json
echo ""

# 验证数据（试运行）
echo "[4/5] 验证数据格式..."
go run cmd/migrate/main.go -command=import-novels -file=data/novels_test.json -dry-run=true
echo ""

# 正式导入
echo "[5/5] 导入数据到 MongoDB..."
echo "提示: 请确保 MongoDB 正在运行"
echo ""
go run cmd/migrate/main.go -command=import-novels -file=data/novels_test.json
echo ""

echo "========================================"
echo "✓ 测试完成！"
echo "========================================"
echo ""
echo "下一步:"
echo "1. 启动服务器: go run cmd/server/main.go"
echo "2. 测试书店 API: GET /api/v1/bookstore/books"
echo "3. 测试搜索功能: GET /api/v1/bookstore/books/search?keyword=xxx"
echo "4. 测试章节阅读: GET /api/v1/bookstore/books/{id}/chapters"
echo ""
echo "清理测试数据: go run cmd/migrate/main.go -command=clean-novels"
echo ""

