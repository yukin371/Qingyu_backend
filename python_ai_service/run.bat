@echo off
REM Windows 快速启动脚本

REM 检查 Poetry
where poetry >nul 2>nul
if %errorlevel% neq 0 (
    echo Poetry not found. Installing dependencies with pip...
    pip install -r requirements.txt
) else (
    echo Installing dependencies with Poetry...
    poetry install
)

REM 启动服务
echo Starting Qingyu AI Service...
where poetry >nul 2>nul
if %errorlevel% neq 0 (
    python -m uvicorn src.main:app --reload --host 0.0.0.0 --port 8000
) else (
    poetry run uvicorn src.main:app --reload --host 0.0.0.0 --port 8000
)

