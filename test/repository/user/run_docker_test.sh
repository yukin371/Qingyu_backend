#!/bin/bash

# Bash 脚本：使用 Docker 运行用户 Repository 集成测试

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

echo -e "${CYAN}========================================"
echo -e "  用户管理模块 Repository 集成测试"
echo -e "========================================${NC}"
echo ""

# 1. 检查 Docker 是否运行
echo -e "${YELLOW}[1/5] 检查 Docker 服务...${NC}"
if ! docker ps &> /dev/null; then
    echo -e "${RED}❌ Docker 未运行或未安装${NC}"
    echo -e "${RED}   请先启动 Docker${NC}"
    exit 1
fi
echo -e "${GREEN}✅ Docker 服务正常${NC}"
echo ""

# 2. 启动数据库服务
echo -e "${YELLOW}[2/5] 启动 Docker 数据库服务...${NC}"
cd ../../../docker
docker-compose -f docker-compose.db-only.yml up -d
cd - > /dev/null
echo -e "${GREEN}✅ Docker 服务已启动${NC}"
echo ""

# 3. 等待 MongoDB 就绪
echo -e "${YELLOW}[3/5] 等待 MongoDB 就绪...${NC}"
retries=0
max_retries=30
ready=false

while [ $retries -lt $max_retries ]; do
    if docker exec qingyu-mongodb mongosh --eval "db.runCommand('ping').ok" --quiet &> /dev/null; then
        ready=true
        break
    fi
    retries=$((retries + 1))
    echo "   等待中... ($retries/$max_retries)"
    sleep 1
done

if [ "$ready" = false ]; then
    echo -e "${RED}❌ MongoDB 启动超时${NC}"
    echo -e "${RED}   请检查: docker logs qingyu-mongodb${NC}"
    exit 1
fi
echo -e "${GREEN}✅ MongoDB 已就绪${NC}"
echo ""

# 4. 运行测试
echo -e "${YELLOW}[4/5] 运行集成测试...${NC}"
echo ""
cd ../../..
go test -v ./test/repository/user/...
test_result=$?
cd - > /dev/null
echo ""

if [ $test_result -eq 0 ]; then
    echo -e "${GREEN}✅ 所有测试通过${NC}"
else
    echo -e "${RED}❌ 部分测试失败${NC}"
fi
echo ""

# 5. 询问是否停止服务
echo -e "${YELLOW}[5/5] 清理...${NC}"
read -p "是否停止 Docker 服务? (y/N): " cleanup
if [ "$cleanup" = "y" ] || [ "$cleanup" = "Y" ]; then
    cd ../../../docker
    docker-compose -f docker-compose.db-only.yml down
    cd - > /dev/null
    echo -e "${GREEN}✅ Docker 服务已停止${NC}"
else
    echo -e "${CYAN}ℹ️  Docker 服务继续运行${NC}"
    echo -e "${CYAN}   手动停止: cd docker && docker-compose -f docker-compose.db-only.yml down${NC}"
fi

echo ""
echo -e "${CYAN}========================================"
echo -e "  测试完成"
echo -e "========================================${NC}"

exit $test_result

