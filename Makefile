# Qingyu Backend Makefile

.PHONY: help proto proto-go proto-python run-ai-service test-ai-service
.PHONY: build test test-unit test-integration test-api test-all test-coverage lint fmt vet clean run
.PHONY: docker-up docker-down docker-logs benchmark check security deps deps-update

# 默认目标
help:
	@echo "Qingyu Backend - 可用命令:"
	@echo ""
	@echo "开发命令:"
	@echo "  make build           - 编译项目"
	@echo "  make run             - 运行开发服务器"
	@echo "  make deps            - 下载依赖"
	@echo "  make deps-update     - 更新依赖"
	@echo ""
	@echo "测试命令:"
	@echo "  make test            - 运行所有测试"
	@echo "  make test-unit       - 运行单元测试"
	@echo "  make test-integration - 运行集成测试"
	@echo "  make test-api        - 运行API测试"
	@echo "  make test-all        - 运行所有测试"
	@echo "  make test-coverage   - 生成测试覆盖率报告"
	@echo "  make benchmark       - 运行性能基准测试"
	@echo ""
	@echo "代码质量:"
	@echo "  make fmt             - 格式化代码"
	@echo "  make vet             - 运行 go vet 检查"
	@echo "  make lint            - 运行代码质量检查"
	@echo "  make check           - 运行所有代码检查"
	@echo "  make security        - 运行安全扫描"
	@echo ""
	@echo "Protobuf & AI:"
	@echo "  make proto           - 生成所有 Protobuf 代码"
	@echo "  make proto-go        - 生成 Go Protobuf 代码"
	@echo "  make proto-python    - 生成 Python Protobuf 代码"
	@echo "  make run-ai-service  - 运行 Python AI Service"
	@echo "  make test-ai-service - 测试 Python AI Service"
	@echo ""
	@echo "Docker:"
	@echo "  make docker-up       - 启动所有 Docker 服务"
	@echo "  make docker-down     - 停止所有 Docker 服务"
	@echo "  make docker-logs     - 查看 Docker 服务日志"
	@echo ""
	@echo "清理:"
	@echo "  make clean           - 清理构建文件"

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

# ========== 新增测试和开发命令 ==========

# 编译项目
build:
	@echo "编译项目..."
	go build -o bin/server cmd/server/main.go
	@echo "编译完成: bin/server"

# 运行开发服务器
run:
	@echo "启动开发服务器..."
	go run cmd/server/main.go

# 下载依赖
deps:
	@echo "下载依赖..."
	go mod download
	go mod verify
	@echo "依赖下载完成！"

# 更新依赖
deps-update:
	@echo "更新依赖..."
	go get -u ./...
	go mod tidy
	@echo "依赖更新完成！"

# 运行所有测试
test:
	@echo "运行所有测试..."
	go test -v -race ./...

# 运行单元测试
test-unit:
	@echo "运行单元测试..."
	go test -v -race -count=1 \
		-coverprofile=coverage.out \
		-covermode=atomic \
		./service/... ./api/... ./pkg/... ./middleware/...

# 运行集成测试
test-integration:
	@echo "运行集成测试..."
	go test -v -race -count=1 \
		-tags=integration \
		-coverprofile=integration_coverage.out \
		-covermode=atomic \
		./test/integration/...

# 运行API测试
test-api:
	@echo "运行API测试..."
	go test -v -race -count=1 \
		./test/api/...

# 运行所有测试（包括集成测试）
test-all: test-unit test-integration test-api
	@echo "所有测试完成！"

# 生成测试覆盖率报告
test-coverage:
	@echo "生成测试覆盖率报告..."
	@rm -f coverage.out coverage.html
	go test -v -race -count=1 \
		-coverprofile=coverage.out \
		-covermode=atomic \
		./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "覆盖率报告已生成: coverage.html"
	@echo ""
	@echo "=== 覆盖率摘要 ==="
	@go tool cover -func=coverage.out | tail -1

# 运行性能基准测试
benchmark:
	@echo "运行性能基准测试..."
	go test -bench=. -benchmem -count=5 \
		-run=^$$ \
		./... | tee benchmark.txt

# 格式化代码
fmt:
	@echo "格式化代码..."
	go fmt ./...
	@if command -v goimports >/dev/null 2>&1; then \
		goimports -w .; \
	fi
	@echo "代码格式化完成！"

# 运行 go vet 检查
vet:
	@echo "运行 go vet..."
	go vet ./...
	@echo "go vet 检查完成！"

# 运行代码质量检查
lint:
	@echo "运行 golangci-lint..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run --timeout=5m; \
	else \
		echo "golangci-lint 未安装，跳过..."; \
	fi
	@echo "代码质量检查完成！"

# 运行所有代码检查
check: fmt vet lint
	@echo "所有代码检查完成！"

# 运行安全扫描
security:
	@echo "运行安全扫描..."
	@if command -v gosec >/dev/null 2>&1; then \
		gosec ./...; \
	fi
	@echo "安全扫描完成！"

# 清理构建文件
clean:
	@echo "清理构建文件..."
	@rm -rf bin/
	@rm -f coverage.out coverage.html
	@rm -f integration_coverage.out integration_coverage.html
	@rm -f benchmark.txt
	@echo "清理完成！"
