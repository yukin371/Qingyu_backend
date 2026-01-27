#!/bin/bash
# 备份验证脚本
# 使用mongorestore --dry-run验证备份完整性

set -e  # 任何错误立即退出

# 配置
BACKUP_DIR="${1:?错误: 请提供备份目录路径}"
MONGO_HOST="${MONGO_HOST:-localhost}"
MONGO_PORT="${MONGO_PORT:-27017}"
MONGO_DB="${MONGO_DB:-qingyu}"

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}Qingyu 备份验证脚本${NC}"
echo -e "${GREEN}========================================${NC}"
echo ""

# 1. 检查备份目录
echo -e "${YELLOW}[步骤 1/3] 检查备份目录...${NC}"
if [ ! -d "${BACKUP_DIR}" ]; then
    echo -e "${RED}错误: 备份目录不存在: ${BACKUP_DIR}${NC}"
    exit 1
fi
echo -e "${GREEN}✓ 备份目录存在${NC}"
echo ""

# 2. 检查mongorestore
echo -e "${YELLOW}[步骤 2/3] 检查mongorestore工具...${NC}"
if ! command -v mongorestore &> /dev/null; then
    echo -e "${RED}错误: mongorestore未安装或不在PATH中${NC}"
    exit 1
fi
echo -e "${GREEN}✓ mongorestore已就绪${NC}"
echo ""

# 3. 执行dry-run验证
echo -e "${YELLOW}[步骤 3/3] 执行备份验证(dry-run)...${NC}"
echo "备份目录: ${BACKUP_DIR}"
echo "目标数据库: ${MONGO_DB}"
echo ""

# 检查数据库子目录
DB_BACKUP_PATH="${BACKUP_DIR}/${MONGO_DB}"
if [ ! -d "${DB_BACKUP_PATH}" ]; then
    echo -e "${RED}错误: 数据库备份子目录不存在: ${DB_BACKUP_PATH}${NC}"
    exit 1
fi

# 执行dry-run验证
mongorestore \
    --host="${MONGO_HOST}:${MONGO_PORT}" \
    --db="${MONGO_DB}" \
    --gzip \
    --drop \
    --dry-run \
    "${DB_BACKUP_PATH}"

echo ""
echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}备份验证通过！${NC}"
echo -e "${GREEN}========================================${NC}"
echo "备份目录: ${BACKUP_DIR}"
echo "验证时间: $(date '+%Y-%m-%d %H:%M:%S')"
echo ""
echo -e "${GREEN}✓ 备份可以安全用于恢复${NC}"
