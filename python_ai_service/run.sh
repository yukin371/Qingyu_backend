#!/bin/bash
# 快速启动脚本

# 检查是否安装了 Poetry
if ! command -v poetry &> /dev/null; then
    echo "Poetry not found. Installing dependencies with pip..."
    pip install -r requirements.txt
else
    echo "Installing dependencies with Poetry..."
    poetry install
fi

# 启动服务
echo "Starting Qingyu AI Service..."
if command -v poetry &> /dev/null; then
    poetry run uvicorn src.main:app --reload --host 0.0.0.0 --port 8000
else
    python -m uvicorn src.main:app --reload --host 0.0.0.0 --port 8000
fi

