# Qingyu Backend Makefile

.PHONY: help proto proto-go proto-python run-ai-service test-ai-service

# 默认目标
help:
	@echo "Qingyu Backend - 可用命令:"
	@echo "  make proto           - 生成所有 Protobuf 代码"
	@echo "  make proto-go        - 生成 Go Protobuf 代码"
	@echo "  make proto-python    - 生成 Python Protobuf 代码"
	@echo "  make run-ai-service  - 运行 Python AI Service"
	@echo "  make test-ai-service - 测试 Python AI Service"
	@echo "  make docker-up       - 启动所有 Docker 服务"
	@echo "  make docker-down     - 停止所有 Docker 服务"

# 生成所有 Protobuf 代码
proto: proto-go proto-python

# 生成 Go Protobuf 代码
proto-go:
	@echo "Generating Go protobuf code..."
	protoc --go_out=. --go-grpc_out=. \
		--go_opt=paths=source_relative \
		--go-grpc_opt=paths=source_relative \
		-I python_ai_service/proto \
		python_ai_service/proto/ai_service.proto
	@echo "Go protobuf code generated in pkg/grpc/pb/"

# 生成 Python Protobuf 代码
proto-python:
	@echo "Generating Python protobuf code..."
	cd python_ai_service && \
	python -m grpc_tools.protoc -I proto \
		--python_out=src/grpc_server \
		--grpc_python_out=src/grpc_server \
		proto/ai_service.proto
	@echo "Python protobuf code generated in python_ai_service/src/grpc_server/"

# 运行 Python AI Service
run-ai-service:
	@echo "Starting Python AI Service..."
	cd python_ai_service && poetry run uvicorn src.main:app --reload --host 0.0.0.0 --port 8000

# 测试 Python AI Service
test-ai-service:
	@echo "Testing Python AI Service..."
	cd python_ai_service && poetry run pytest tests/ -v

# Docker 相关
docker-up:
	@echo "Starting Docker services..."
	cd docker && docker-compose -f docker-compose.dev.yml up -d

docker-down:
	@echo "Stopping Docker services..."
	cd docker && docker-compose -f docker-compose.dev.yml down

docker-logs:
	cd docker && docker-compose -f docker-compose.dev.yml logs -f
