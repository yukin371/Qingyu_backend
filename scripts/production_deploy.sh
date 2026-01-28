#!/bin/bash
# 生产环境部署脚本
# Block3 数据库优化索引部署

set -e  # 任何错误立即退出

# 配置
BACKUP_DIR="${BACKUP_DIR:-/tmp/qingyu_backup}"
TIMESTAMP=$(date +%Y%m%d_%H%M%S)
MONGO_HOST="${MONGO_HOST:-localhost}"
MONGO_PORT="${MONGO_PORT:-27017}"
MONGO_DB="${MONGO_DB:-qingyu}"
PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 日志函数
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}Block3 数据库优化 - 生产部署${NC}"
echo -e "${GREEN}========================================${NC}"
echo ""
log_info "开始时间: $(date '+%Y-%m-%d %H:%M:%S')"
log_info "项目根目录: ${PROJECT_ROOT}"
log_info "目标数据库: ${MONGO_DB}"
echo ""

# ============================================================================
# 步骤 1: 备份数据库
# ============================================================================
echo -e "${YELLOW}========================================${NC}"
echo -e "${YELLOW}[步骤 1/4] 数据库备份${NC}"
echo -e "${YELLOW}========================================${NC}"
log_info "执行备份脚本..."
BACKUP_NAME=$("${PROJECT_ROOT}/scripts/backup_database.sh")
BACKUP_PATH="${BACKUP_DIR}/qingyu_backup_${TIMESTAMP}"

if [ ! -d "${BACKUP_PATH}" ]; then
    log_error "备份失败: 备份目录不存在"
    exit 1
fi

log_success "数据库备份完成"
log_info "备份位置: ${BACKUP_PATH}"
echo ""

# ============================================================================
# 步骤 2: 验证备份
# ============================================================================
echo -e "${YELLOW}========================================${NC}"
echo -e "${YELLOW}[步骤 2/4] 验证备份完整性${NC}"
echo -e "${YELLOW}========================================${NC}"
log_info "执行备份验证..."

"${PROJECT_ROOT}/scripts/verify_backup.sh" "${BACKUP_PATH}"

if [ $? -ne 0 ]; then
    log_error "备份验证失败"
    log_error "部署已中止，请检查备份"
    exit 1
fi

log_success "备份验证通过"
echo ""

# ============================================================================
# 步骤 3: 执行索引迁移
# ============================================================================
echo -e "${YELLOW}========================================${NC}"
echo -e "${YELLOW}[步骤 3/4] 执行索引迁移${NC}"
echo -e "${YELLOW}========================================${NC}"
log_info "运行索引迁移..."

cd "${PROJECT_ROOT}"
go run cmd/migrate/main.go

if [ $? -ne 0 ]; then
    log_error "索引迁移失败"
    log_error "部署已中止，备份仍保存在: ${BACKUP_PATH}"
    exit 1
fi

log_success "索引迁移完成"
echo ""

# ============================================================================
# 步骤 4: 验证索引创建
# ============================================================================
echo -e "${YELLOW}========================================${NC}"
echo -e "${YELLOW}[步骤 4/4] 验证索引创建${NC}"
echo -e "${YELLOW}========================================${NC}"
log_info "运行索引验证测试..."

cd "${PROJECT_ROOT}"
go test -v ./scripts -run TestVerifyIndexes

if [ $? -ne 0 ]; then
    log_warning "部分索引验证失败"
    log_warning "请手动检查索引状态"
    log_warning "备份仍保存在: ${BACKUP_PATH}"
else
    log_success "所有索引验证通过"
fi
echo ""

# ============================================================================
# 部署完成
# ============================================================================
echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}部署完成！${NC}"
echo -e "${GREEN}========================================${NC}"
log_info "完成时间: $(date '+%Y-%m-%d %H:%M:%S')"
echo ""
log_success "✓ 数据库备份: ${BACKUP_PATH}"
log_success "✓ 索引迁移: 完成"
log_success "✓ 索引验证: 完成"
echo ""
echo -e "${YELLOW}后续操作:${NC}"
echo "1. 检查应用日志，确认性能改善"
echo "2. 监控数据库查询性能"
echo "3. 7天后确认无问题可删除备份: rm -rf ${BACKUP_PATH}"
echo ""
echo -e "${YELLOW}如需回滚:${NC}"
echo "mongorestore --host=${MONGO_HOST}:${MONGO_PORT} --db=${MONGO_DB} --gzip --drop ${BACKUP_PATH}/${MONGO_DB}"
