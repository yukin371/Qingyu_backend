#!/bin/bash
# 数据库备份脚本
# 用于生产环境部署前备份MongoDB数据

set -e  # 任何错误立即退出

# 配置
BACKUP_DIR="${BACKUP_DIR:-/tmp/qingyu_backup}"
TIMESTAMP=$(date +%Y%m%d_%H%M%S)
BACKUP_NAME="qingyu_backup_${TIMESTAMP}"
MONGO_HOST="${MONGO_HOST:-localhost}"
MONGO_PORT="${MONGO_PORT:-27017}"
MONGO_DB="${MONGO_DB:-qingyu}"

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}Qingyu 数据库备份脚本${NC}"
echo -e "${GREEN}========================================${NC}"
echo ""

# 1. 检查mongodump是否可用
echo -e "${YELLOW}[步骤 1/4] 检查mongodump工具...${NC}"
if ! command -v mongodump &> /dev/null; then
    echo -e "${RED}错误: mongodump未安装或不在PATH中${NC}"
    echo "请安装MongoDB工具: sudo apt-get install mongodb-tools"
    exit 1
fi
echo -e "${GREEN}✓ mongodump已就绪${NC}"
echo ""

# 2. 创建备份目录
echo -e "${YELLOW}[步骤 2/4] 创建备份目录...${NC}"
mkdir -p "${BACKUP_DIR}"
echo -e "${GREEN}✓ 备份目录: ${BACKUP_DIR}${NC}"
echo ""

# 3. 执行备份
echo -e "${YELLOW}[步骤 3/4] 开始备份数据库...${NC}"
echo "数据库: ${MONGO_DB}"
echo "主机: ${MONGO_HOST}:${MONGO_PORT}"
echo "目标: ${BACKUP_DIR}/${BACKUP_NAME}"
echo ""

mongodump \
    --host="${MONGO_HOST}:${MONGO_PORT}" \
    --db="${MONGO_DB}" \
    --out="${BACKUP_DIR}/${BACKUP_NAME}" \
    --gzip

echo -e "${GREEN}✓ 备份完成${NC}"
echo ""

# 4. 验证备份
echo -e "${YELLOW}[步骤 4/4] 验证备份...${NC}"
BACKUP_PATH="${BACKUP_DIR}/${BACKUP_NAME}/${MONGO_DB}"

if [ ! -d "${BACKUP_PATH}" ]; then
    echo -e "${RED}错误: 备份目录不存在${NC}"
    exit 1
fi

# 统计备份的集合数量
COLLECTION_COUNT=$(find "${BACKUP_PATH}" -name "*.metadata.json" | wc -l)

# 计算备份大小
BACKUP_SIZE=$(du -sh "${BACKUP_DIR}/${BACKUP_NAME}" | cut -f1)

echo -e "${GREEN}✓ 备份验证通过${NC}"
echo ""
echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}备份成功！${NC}"
echo -e "${GREEN}========================================${NC}"
echo "备份位置: ${BACKUP_DIR}/${BACKUP_NAME}"
echo "备份大小: ${BACKUP_SIZE}"
echo "集合数量: ${COLLECTION_COUNT}"
echo "时间戳: ${TIMESTAMP}"
echo ""
echo -e "${YELLOW}提示: 使用以下命令恢复备份${NC}"
echo "mongorestore --host=${MONGO_HOST}:${MONGO_PORT} --db=${MONGO_DB} --gzip ${BACKUP_DIR}/${BACKUP_NAME}/${MONGO_DB}"
