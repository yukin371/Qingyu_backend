#!/bin/bash

# 测试环境设置脚本
# 用于本地开发和CI/CD环境准备

set -e  # 遇到错误立即退出

echo "========================================"
echo "🔧 青羽平台测试环境设置"
echo "========================================"

# ========== 颜色定义 ==========
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# ========== 辅助函数 ==========
print_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

print_error() {
    echo -e "${RED}❌ $1${NC}"
}

print_info() {
    echo -e "ℹ️  $1"
}

# ========== 1. 检查Go环境 ==========
echo ""
echo "📦 1. 检查Go环境..."

if ! command -v go &> /dev/null; then
    print_error "Go未安装，请先安装Go 1.21或更高版本"
    exit 1
fi

GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
print_success "Go版本: $GO_VERSION"

# ========== 2. 安装测试工具 ==========
echo ""
echo "🛠️  2. 安装测试工具..."

# gotests - 测试代码生成工具
if ! command -v gotests &> /dev/null; then
    print_info "安装gotests..."
    go install github.com/cweill/gotests/gotests@latest
    print_success "gotests安装完成"
else
    print_success "gotests已安装"
fi

# mockgen - Mock代码生成工具（可选）
if ! command -v mockgen &> /dev/null; then
    print_info "安装mockgen..."
    go install github.com/golang/mock/mockgen@latest
    print_success "mockgen安装完成"
else
    print_success "mockgen已安装"
fi

# golangci-lint - 代码质量检查工具（可选）
if ! command -v golangci-lint &> /dev/null; then
    print_warning "golangci-lint未安装（可选工具）"
    print_info "可通过以下命令安装: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"
else
    print_success "golangci-lint已安装"
fi

# ========== 3. 检查测试服务 ==========
echo ""
echo "🗄️  3. 检查测试服务状态..."

# 检查MongoDB
if command -v mongosh &> /dev/null || command -v mongo &> /dev/null; then
    if nc -z localhost 27017 2>/dev/null; then
        print_success "MongoDB服务运行中 (localhost:27017)"
    else
        print_warning "MongoDB服务未运行"
        print_info "可通过Docker启动: docker run -d -p 27017:27017 --name test-mongo mongo:5.0"
    fi
else
    print_warning "MongoDB客户端未安装"
fi

# 检查Redis
if command -v redis-cli &> /dev/null; then
    if redis-cli -p 6379 ping &> /dev/null; then
        print_success "Redis服务运行中 (localhost:6379)"
    else
        print_warning "Redis服务未运行"
        print_info "可通过Docker启动: docker run -d -p 6379:6379 --name test-redis redis:6.2-alpine"
    fi
else
    print_warning "Redis客户端未安装"
fi

# ========== 4. 设置环境变量 ==========
echo ""
echo "🌍 4. 设置测试环境变量..."

export GO_ENV=test
export MONGODB_URI=${MONGODB_URI:-"mongodb://test:test123@localhost:27017/qingyu_test"}
export REDIS_ADDR=${REDIS_ADDR:-"localhost:6379"}
export REDIS_DB=${REDIS_DB:-1}

print_success "环境变量已设置:"
echo "  GO_ENV=$GO_ENV"
echo "  MONGODB_URI=$MONGODB_URI"
echo "  REDIS_ADDR=$REDIS_ADDR"
echo "  REDIS_DB=$REDIS_DB"

# ========== 5. 下载Go依赖 ==========
echo ""
echo "📦 5. 下载Go模块依赖..."

go mod download
go mod verify
print_success "依赖下载完成"

# ========== 6. 创建测试目录 ==========
echo ""
echo "📁 6. 检查测试目录结构..."

TEST_DIRS=(
    "test/testutil"
    "test/fixtures"
    "test/examples"
    "test/integration"
    "test/api"
)

for dir in "${TEST_DIRS[@]}"; do
    if [ ! -d "$dir" ]; then
        mkdir -p "$dir"
        print_info "创建目录: $dir"
    fi
done

print_success "测试目录结构完整"

# ========== 7. 快速测试验证 ==========
echo ""
echo "🧪 7. 运行快速测试验证..."

if go test -short -v ./... &> /dev/null; then
    print_success "快速测试验证通过"
else
    print_warning "快速测试验证失败（可能是因为测试服务未启动或测试代码有问题）"
fi

# ========== 完成 ==========
echo ""
echo "========================================"
print_success "测试环境设置完成！"
echo "========================================"
echo ""
echo "📝 后续步骤:"
echo "  1. 启动测试服务（如果尚未启动）："
echo "     docker-compose -f docker-compose.test.yml up -d"
echo ""
echo "  2. 运行测试："
echo "     make test          # 运行所有测试"
echo "     make test-unit     # 运行单元测试"
echo "     make test-coverage # 生成覆盖率报告"
echo ""
echo "  3. 生成测试模板："
echo "     make test-gen file=service/user/user_service.go"
echo ""
echo "  4. 查看更多命令："
echo "     make help"
echo ""

