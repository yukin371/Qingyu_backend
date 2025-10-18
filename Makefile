# 青羽轻量级读写平台 Makefile

# 变量定义
APP_NAME := qingyu-backend
GO_VERSION := 1.21
DOCKER_IMAGE := $(APP_NAME):latest
BUILD_DIR := ./build
MAIN_FILE := ./cmd/server/main.go

# 默认目标
.DEFAULT_GOAL := help

# 帮助信息
.PHONY: help
help: ## 显示帮助信息
	@echo "青羽轻量级读写平台 - 可用命令："
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

# 开发相关命令
.PHONY: run
run: ## 启动开发服务器
	@echo "启动青羽后端服务..."
	go run $(MAIN_FILE)

.PHONY: dev
dev: ## 启动开发服务器（带热重载）
	@echo "启动开发服务器（热重载模式）..."
	@if command -v air > /dev/null; then \
		air; \
	else \
		echo "请先安装 air: go install github.com/cosmtrek/air@latest"; \
		exit 1; \
	fi

.PHONY: build
build: ## 构建应用程序
	@echo "构建应用程序..."
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(APP_NAME) $(MAIN_FILE)

.PHONY: clean
clean: ## 清理构建文件
	@echo "清理构建文件..."
	@rm -rf $(BUILD_DIR)
	@go clean

# 测试相关命令
.PHONY: test
test: ## 运行所有测试（带竞态检测）
	@echo "运行测试..."
	go test -v -race ./...

.PHONY: test-unit
test-unit: ## 运行单元测试（Service和Repository层）
	@echo "运行单元测试..."
	go test -v -short -race ./service/... ./repository/... ./test/service/... ./test/repository/...

.PHONY: test-integration
test-integration: ## 运行集成测试
	@echo "运行集成测试..."
	go test -v -run Integration ./test/integration/...

.PHONY: test-api
test-api: ## 运行API测试
	@echo "运行API测试..."
	go test -v ./test/api/...

.PHONY: test-coverage
test-coverage: ## 运行测试并生成覆盖率报告
	@echo "生成测试覆盖率报告..."
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "覆盖率报告已生成: coverage.html"
	@go tool cover -func=coverage.out | grep total

.PHONY: test-coverage-check
test-coverage-check: ## 检查覆盖率是否达到80%
	@echo "检查测试覆盖率..."
	@go test -coverprofile=coverage.out ./... > /dev/null 2>&1
	@coverage=$$(go tool cover -func=coverage.out | grep total | awk '{print $$3}' | sed 's/%//'); \
	echo "总覆盖率: $${coverage}%"; \
	threshold=80; \
	result=$$(echo "$${coverage} >= $${threshold}" | bc -l 2>/dev/null || echo "0"); \
	if [ "$$result" = "1" ]; then \
		echo "✅ 覆盖率达标 (≥80%)"; \
	else \
		echo "❌ 覆盖率低于80%"; \
		exit 1; \
	fi

.PHONY: test-gen
test-gen: ## 为指定文件生成测试模板（需要安装gotests）
	@if [ -z "$(file)" ]; then \
		echo "用法: make test-gen file=path/to/file.go"; \
		echo "示例: make test-gen file=service/user/user_service.go"; \
		exit 1; \
	fi
	@if ! command -v gotests > /dev/null; then \
		echo "正在安装 gotests..."; \
		go install github.com/cweill/gotests/gotests@latest; \
	fi
	@echo "为 $(file) 生成测试..."
	gotests -all -w $(file)
	@echo "✅ 测试文件已生成"

.PHONY: test-clean
test-clean: ## 清理测试缓存和覆盖率文件
	@echo "清理测试缓存和覆盖率文件..."
	go clean -testcache
	rm -f coverage.out coverage.html
	@echo "✅ 清理完成"

.PHONY: test-watch
test-watch: ## 监视文件变化并自动运行测试
	@echo "监视文件变化..."
	@while true; do \
		inotifywait -r -e modify . && clear && make test-unit; \
	done

.PHONY: test-fix
test-fix: ## 只运行之前失败的测试
	@echo "运行上次失败的测试..."
	@if [ -f ".test-failures" ]; then \
		while IFS= read -r pkg; do \
			echo "重新测试: $$pkg"; \
			go test -v "$$pkg"; \
		done < ".test-failures"; \
		rm -f ".test-failures"; \
	else \
		echo "没有记录的失败测试"; \
	fi

.PHONY: test-quick
test-quick: ## 快速测试（排除慢速测试，不使用race检测）
	@echo "运行快速测试..."
	go test -short ./...

.PHONY: test-report
test-report: ## 生成详细的测试报告（JSON格式）
	@echo "生成测试报告..."
	@mkdir -p test_reports
	go test -json -coverprofile=coverage.out ./... > test_reports/test-output.json
	go tool cover -func=coverage.out > test_reports/coverage.txt
	go tool cover -html=coverage.out -o test_reports/coverage.html
	@echo "✅ 测试报告已生成在 test_reports/ 目录"
	@echo "  - test-output.json: JSON格式的测试结果"
	@echo "  - coverage.txt: 文本格式的覆盖率报告"
	@echo "  - coverage.html: HTML格式的覆盖率报告"

.PHONY: test-verbose
test-verbose: ## 运行测试并显示详细输出
	@echo "运行详细测试..."
	go test -v -cover -coverprofile=coverage.out ./... 2>&1 | tee test.log

.PHONY: test-bench
test-bench: ## 运行基准测试
	@echo "运行基准测试..."
	go test -bench=. -benchmem ./...

# 代码质量相关命令
.PHONY: lint
lint: ## 运行代码检查
	@echo "运行代码检查..."
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run; \
	else \
		echo "请先安装 golangci-lint: https://golangci-lint.run/usage/install/"; \
		exit 1; \
	fi

.PHONY: fmt
fmt: ## 格式化代码
	@echo "格式化代码..."
	go fmt ./...

.PHONY: vet
vet: ## 运行 go vet
	@echo "运行 go vet..."
	go vet ./...

# Mock 生成相关命令
.PHONY: mock
mock: ## 生成 Mock 文件
	@echo "生成 Mock 文件..."
	@if command -v mockgen > /dev/null; then \
		go generate ./...; \
	else \
		echo "请先安装 mockgen: go install github.com/golang/mock/mockgen@latest"; \
		exit 1; \
	fi

.PHONY: mock-clean
mock-clean: ## 清理 Mock 文件
	@echo "清理 Mock 文件..."
	find . -name "*_mock.go" -type f -delete

# 依赖管理
.PHONY: deps
deps: ## 下载依赖
	@echo "下载依赖..."
	go mod download

.PHONY: deps-update
deps-update: ## 更新依赖
	@echo "更新依赖..."
	go mod tidy
	go mod download

.PHONY: deps-vendor
deps-vendor: ## 创建 vendor 目录
	@echo "创建 vendor 目录..."
	go mod vendor

# Docker 相关命令
.PHONY: docker-build
docker-build: ## 构建 Docker 镜像
	@echo "构建 Docker 镜像..."
	docker build -t $(DOCKER_IMAGE) .

.PHONY: docker-run
docker-run: ## 运行 Docker 容器
	@echo "运行 Docker 容器..."
	docker run -p 8080:8080 $(DOCKER_IMAGE)

.PHONY: docker-clean
docker-clean: ## 清理 Docker 镜像
	@echo "清理 Docker 镜像..."
	docker rmi $(DOCKER_IMAGE) || true

# 数据库相关命令
.PHONY: db-migrate
db-migrate: ## 运行数据库迁移
	@echo "运行数据库迁移..."
	@echo "TODO: 实现数据库迁移逻辑"

.PHONY: db-seed
db-seed: ## 填充测试数据
	@echo "填充测试数据..."
	@echo "TODO: 实现测试数据填充逻辑"

# 部署相关命令
.PHONY: deploy-dev
deploy-dev: ## 部署到开发环境
	@echo "部署到开发环境..."
	@echo "TODO: 实现开发环境部署逻辑"

.PHONY: deploy-prod
deploy-prod: ## 部署到生产环境
	@echo "部署到生产环境..."
	@echo "TODO: 实现生产环境部署逻辑"

# 工具安装
.PHONY: install-tools
install-tools: ## 安装开发工具
	@echo "安装开发工具..."
	go install github.com/cosmtrek/air@latest
	go install github.com/golang/mock/mockgen@latest
	@echo "开发工具安装完成"

# 项目初始化
.PHONY: init
init: deps install-tools ## 初始化项目环境
	@echo "项目环境初始化完成"

# 完整的 CI 流程
.PHONY: ci
ci: fmt vet lint test ## 运行完整的 CI 流程
	@echo "CI 流程执行完成"

# 快速检查
.PHONY: check
check: fmt vet lint test-unit ## 快速代码检查
	@echo "代码检查完成"