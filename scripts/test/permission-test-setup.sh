#!/bin/bash

# ============================================
# Qingyu Backend - 权限系统测试环境准备脚本
# ============================================
#
# 功能:
#   1. 检查MongoDB连接状态
#   2. 检查Redis连接状态
#   3. 创建/重置测试数据库
#   4. 初始化测试数据
#
# 使用方法:
#   bash scripts/test/permission-test-setup.sh
#   bash scripts/test/permission-test-setup.sh --skip-data
#   bash scripts/test/permission-test-setup.sh --db-only
#
# ============================================

set -e  # 遇到错误立即退出

# ==================== 颜色定义 ====================
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# ==================== 配置参数 ====================
DB_NAME="${QINGYU_TEST_DB_NAME:-qingyu_permission_test}"
MONGO_HOST="${QINGYU_MONGO_HOST:-localhost}"
MONGO_PORT="${QINGYU_MONGO_PORT:-27017}"
MONGO_URI="mongodb://${MONGO_HOST}:${MONGO_PORT}"

REDIS_HOST="${QINGYU_REDIS_HOST:-localhost}"
REDIS_PORT="${QINGYU_REDIS_PORT:-6379}"
REDIS_PASSWORD="${QINGYU_REDIS_PASSWORD:-}"

SKIP_DATA=false
DB_ONLY=false

# ==================== 工具函数 ====================

print_header() {
    echo ""
    echo "========================================"
    echo "$1"
    echo "========================================"
    echo ""
}

print_step() {
    echo -e "${CYAN}[$1/$2] $2...${NC}"
}

print_success() {
    echo -e "${GREEN}✓ $1${NC}"
}

print_error() {
    echo -e "${RED}✗ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠ $1${NC}"
}

print_info() {
    echo -e "${BLUE}ℹ $1${NC}"
}

# ==================== 参数解析 ====================

parse_args() {
    while [[ $# -gt 0 ]]; do
        case $1 in
            --skip-data)
                SKIP_DATA=true
                shift
                ;;
            --db-only)
                DB_ONLY=true
                shift
                ;;
            --help)
                show_help
                exit 0
                ;;
            *)
                print_error "未知参数: $1"
                show_help
                exit 1
                ;;
        esac
    done
}

show_help() {
    cat << EOF
使用方法: $0 [选项]

选项:
  --skip-data    跳过测试数据填充
  --db-only      仅准备数据库（不检查Redis）
  --help         显示此帮助信息

环境变量:
  QINGYU_TEST_DB_NAME      测试数据库名称 (默认: qingyu_permission_test)
  QINGYU_MONGO_HOST        MongoDB主机 (默认: localhost)
  QINGYU_MONGO_PORT        MongoDB端口 (默认: 27017)
  QINGYU_REDIS_HOST        Redis主机 (默认: localhost)
  QINGYU_REDIS_PORT        Redis端口 (默认: 6379)
  QINGYU_REDIS_PASSWORD    Redis密码 (默认: 空)

示例:
  # 完整设置（包括数据）
  $0

  # 仅检查数据库
  $0 --db-only

  # 仅设置数据库，不填充数据
  $0 --skip-data

EOF
}

# ==================== 检查函数 ====================

check_mongodb() {
    print_step 1 4 "检查MongoDB连接"

    # 检查MongoDB是否运行
    if command -v mongosh &> /dev/null; then
        MONGO_CMD="mongosh"
    elif command -v mongo &> /dev/null; then
        MONGO_CMD="mongo"
    else
        print_error "未找到MongoDB客户端 (mongosh 或 mongo)"
        print_info "请安装MongoDB客户端: https://www.mongodb.com/try/download"
        return 1
    fi

    # 尝试连接MongoDB
    if $MONGO_CMD --quiet --host "$MONGO_HOST" --port "$MONGO_PORT" --eval "db.version()" > /dev/null 2>&1; then
        local version=$($MONGO_CMD --quiet --host "$MONGO_HOST" --port "$MONGO_PORT" --eval "db.version()")
        print_success "MongoDB连接成功 (版本: $version)"
        print_info "  主机: $MONGO_HOST:$MONGO_PORT"
        return 0
    else
        print_error "无法连接到MongoDB"
        print_info "请确保MongoDB服务已启动:"
        print_info "  Linux:   sudo systemctl start mongod"
        print_info "  Mac:     brew services start mongodb-community"
        print_info "  Windows: net start MongoDB"
        return 1
    fi
}

check_redis() {
    if [ "$DB_ONLY" = true ]; then
        print_info "跳过Redis检查 (--db-only 模式)"
        return 0
    fi

    print_step 2 4 "检查Redis连接"

    # 检查redis-cli是否可用
    if ! command -v redis-cli &> /dev/null; then
        print_warning "未找到redis-cli命令"
        print_info "请安装Redis客户端或使用Docker启动Redis:"
        print_info "  docker run -d -p 6379:6379 redis:alpine"
        return 1
    fi

    # 尝试连接Redis
    if [ -z "$REDIS_PASSWORD" ]; then
        REDIStest=$(redis-cli -h "$REDIS_HOST" -p "$REDIS_PORT" ping 2>/dev/null || echo "")
    else
        REDIStest=$(redis-cli -h "$REDIS_HOST" -p "$REDIS_PORT" -a "$REDIS_PASSWORD" ping 2>/dev/null || echo "")
    fi

    if [ "$REDIStest" = "PONG" ]; then
        print_success "Redis连接成功"
        print_info "  主机: $REDIS_HOST:$REDIS_PORT"
        return 0
    else
        print_error "无法连接到Redis"
        print_info "请确保Redis服务已启动:"
        print_info "  Linux:   sudo systemctl start redis"
        print_info "  Mac:     brew services start redis"
        print_info "  Windows: redis-server.exe"
        print_info "  Docker:  docker run -d -p 6379:6379 redis:alpine"
        return 1
    fi
}

setup_database() {
    print_step 3 4 "准备测试数据库"

    print_info "数据库名称: $DB_NAME"

    # 检查数据库是否存在
    local db_exists=$($MONGO_CMD --quiet --host "$MONGO_HOST" --port "$MONGO_PORT" --eval "db.getSiblingDB('$DB_NAME').getName()")

    if [ "$db_exists" = "$DB_NAME" ]; then
        print_warning "数据库 '$DB_NAME' 已存在"

        # 询问是否删除并重建
        if [ "$FORCE_RECREATE" = "true" ]; then
            print_info "强制重建数据库..."
            $MONGO_CMD --quiet --host "$MONGO_HOST" --port "$MONGO_PORT" --eval "db.getSiblingDB('$DB_NAME').dropDatabase()"
            print_success "数据库已删除"
        else
            print_info "使用现有数据库 (使用 FORCE_RECREATE=true 来重建)"
        fi
    fi

    # 创建必要的集合和索引
    print_info "创建集合和索引..."

    # 创建roles集合
    $MONGO_CMD --quiet --host "$MONGO_HOST" --port "$MONGO_PORT" "$DB_NAME" --eval "
    db.createCollection('roles');
    db.roles.createIndex({ name: 1 }, { unique: true });
    db.roles.createIndex({ is_system: 1 });
    db.roles.createIndex({ is_default: 1 });
    " > /dev/null 2>&1

    # 创建users集合（如果不存在）
    $MONGO_CMD --quiet --host "$MONGO_HOST" --port "$MONGO_PORT" "$DB_NAME" --eval "
    if (!db.getCollectionNames().includes('users')) {
        db.createCollection('users');
        db.users.createIndex({ username: 1 }, { unique: true });
        db.users.createIndex({ email: 1 }, { unique: true });
    }
    " > /dev/null 2>&1

    print_success "数据库准备完成"
    return 0
}

init_test_data() {
    if [ "$SKIP_DATA" = true ]; then
        print_info "跳过测试数据填充 (--skip-data)"
        return 0
    fi

    print_step 4 4 "初始化测试数据"

    # 检查Go是否安装
    if ! command -v go &> /dev/null; then
        print_error "未找到Go命令"
        print_info "请安装Go: https://golang.org/dl/"
        return 1
    fi

    # 检查是否在项目根目录
    if [ ! -f "go.mod" ]; then
        print_error "请在项目根目录运行此脚本"
        return 1
    fi

    print_info "运行测试数据填充脚本..."

    # 运行Go测试数据脚本
    if go run scripts/test/permission-test-data.go --db="$DB_NAME"; then
        print_success "测试数据填充完成"
        return 0
    else
        print_error "测试数据填充失败"
        return 1
    fi
}

print_summary() {
    print_header "测试环境准备完成！"

    echo -e "${GREEN}测试环境信息:${NC}"
    echo ""
    echo "  数据库: $DB_NAME"
    echo "  MongoDB: $MONGO_HOST:$MONGO_PORT"
    echo "  Redis: $REDIS_HOST:$REDIS_PORT"
    echo ""

    if [ "$SKIP_DATA" = false ]; then
        echo -e "${GREEN}测试账号:${NC}"
        echo ""
        echo "  管理员: admin@test.com / Admin@123"
        echo "  作者:   author@test.com / Author@123"
        echo "  读者:   reader@test.com / Reader@123"
        echo "  编辑:   editor@test.com / Editor@123"
        echo ""
    fi

    echo -e "${GREEN}下一步:${NC}"
    echo ""
    echo "  1. 启动测试服务器:"
    echo "     export QINGYU_DATABASE_NAME=$DB_NAME"
    echo "     go run cmd/server/main.go"
    echo ""
    echo "  2. 运行权限测试:"
    echo "     go test ./internal/middleware/auth/... -v"
    echo ""

    if [ "$DB_ONLY" = false ] && [ "$SKIP_DATA" = false ]; then
        echo -e "${CYAN}快速测试命令:${NC}"
        echo ""
        echo "  # 测试管理员登录"
        echo "  curl -X POST http://localhost:8080/api/v1/auth/login \\"
        echo "    -H \"Content-Type: application/json\" \\"
        echo "    -d '{\"username\":\"admin@test.com\",\"password\":\"Admin@123\"}'"
        echo ""
    fi

    echo "========================================"
    echo ""
}

# ==================== 主流程 ====================

main() {
    print_header "Qingyu Backend - 权限系统测试环境准备"

    # 解析参数
    parse_args "$@"

    # 检查MongoDB
    if ! check_mongodb; then
        print_error "MongoDB检查失败，退出"
        exit 1
    fi

    # 检查Redis
    if ! check_redis; then
        if [ "$DB_ONLY" = false ]; then
            print_warning "Redis检查失败，但继续进行..."
        fi
    fi

    # 准备数据库
    if ! setup_database; then
        print_error "数据库准备失败，退出"
        exit 1
    fi

    # 初始化测试数据
    if ! init_test_data; then
        print_error "测试数据初始化失败，退出"
        exit 1
    fi

    # 打印摘要
    print_summary

    print_success "所有步骤完成！"
}

# 运行主流程
main "$@"
