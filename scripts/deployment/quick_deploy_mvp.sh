#!/bin/bash

# MVP一键部署脚本
# 用途：快速部署MVP到内测环境
# 日期：2025-10-23

set -e

echo "======================================"
echo "🚀 MVP一键部署脚本"
echo "======================================"
echo ""

# 颜色定义
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 步骤计数
STEP=1

# 打印步骤
print_step() {
    echo ""
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo -e "${BLUE}步骤 $STEP: $1${NC}"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    STEP=$((STEP + 1))
}

# 错误处理
error_exit() {
    echo -e "${RED}✗ 错误: $1${NC}"
    exit 1
}

# 1. 环境检查
print_step "环境检查"

# 检查Docker
if ! command -v docker &> /dev/null; then
    error_exit "Docker未安装，请先安装Docker"
fi
echo -e "${GREEN}✓${NC} Docker已安装: $(docker --version)"

# 检查docker-compose
if ! command -v docker-compose &> /dev/null; then
    error_exit "docker-compose未安装，请先安装docker-compose"
fi
echo -e "${GREEN}✓${NC} docker-compose已安装: $(docker-compose --version)"

# 2. 清理旧部署
print_step "清理旧部署"

if docker-compose -f docker/docker-compose.prod.yml ps | grep -q "Up"; then
    echo "正在停止旧服务..."
    docker-compose -f docker/docker-compose.prod.yml down
    echo -e "${GREEN}✓${NC} 旧服务已停止"
else
    echo "没有运行中的旧服务"
fi

# 清理旧镜像（可选）
# echo "清理旧镜像..."
# docker image rm qingyu-backend:mvp 2>/dev/null || true

# 3. 构建Docker镜像
print_step "构建Docker镜像"

echo "开始构建镜像（这可能需要几分钟）..."
if docker build -t qingyu-backend:mvp -f docker/Dockerfile.prod . ; then
    echo -e "${GREEN}✓${NC} 镜像构建成功"
else
    error_exit "镜像构建失败"
fi

# 4. 启动服务
print_step "启动服务容器"

echo "启动MongoDB、Redis和后端服务..."
if docker-compose -f docker/docker-compose.prod.yml up -d; then
    echo -e "${GREEN}✓${NC} 服务启动成功"
else
    error_exit "服务启动失败"
fi

# 等待服务就绪
echo "等待服务就绪（10秒）..."
sleep 10

# 检查服务状态
echo ""
echo "服务状态："
docker-compose -f docker/docker-compose.prod.yml ps

# 5. 数据库初始化
print_step "数据库初始化"

echo "等待MongoDB就绪..."
MAX_RETRIES=30
RETRY_COUNT=0

while ! docker exec qingyu-mongodb mongosh --quiet --eval "db.runCommand({ ping: 1 })" &> /dev/null; do
    RETRY_COUNT=$((RETRY_COUNT + 1))
    if [ $RETRY_COUNT -ge $MAX_RETRIES ]; then
        error_exit "MongoDB启动超时"
    fi
    echo "等待MongoDB... ($RETRY_COUNT/$MAX_RETRIES)"
    sleep 2
done

echo -e "${GREEN}✓${NC} MongoDB已就绪"

# 运行数据库迁移
echo "运行数据库迁移..."
if docker exec qingyu-backend go run cmd/migrate/main.go 2>/dev/null; then
    echo -e "${GREEN}✓${NC} 数据库迁移完成"
else
    echo -e "${YELLOW}⚠${NC} 数据库迁移失败或已执行"
fi

# 6. 创建测试账号（可选）
print_step "创建测试账号"

echo "创建内测账号..."
# 暂时跳过，因为种子数据脚本可能不存在
# docker exec qingyu-backend go run migration/seeds/create_test_users.go 2>/dev/null || echo "种子数据脚本不存在"
echo -e "${YELLOW}⚠${NC} 测试账号创建已跳过（需手动创建）"

# 7. 健康检查
print_step "健康检查"

echo "等待API就绪..."
MAX_RETRIES=30
RETRY_COUNT=0

while ! curl -s http://localhost:8080/health &> /dev/null; do
    RETRY_COUNT=$((RETRY_COUNT + 1))
    if [ $RETRY_COUNT -ge $MAX_RETRIES ]; then
        error_exit "API启动超时"
    fi
    echo "等待API... ($RETRY_COUNT/$MAX_RETRIES)"
    sleep 2
done

echo -e "${GREEN}✓${NC} API已就绪"

# 执行健康检查
echo ""
echo "健康检查结果："
curl -s http://localhost:8080/health | python -m json.tool 2>/dev/null || curl -s http://localhost:8080/health

# 8. 部署验证
print_step "部署验证"

echo "验证关键功能..."

# 验证API版本
echo "- 检查API版本..."
if curl -s http://localhost:8080/api/v1/system/version &> /dev/null; then
    echo -e "${GREEN}✓${NC} API版本接口正常"
else
    echo -e "${YELLOW}⚠${NC} API版本接口异常"
fi

# 9. 部署总结
print_step "部署总结"

echo ""
echo -e "${GREEN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${GREEN}✓ MVP部署成功！${NC}"
echo -e "${GREEN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""

echo "📍 服务地址："
echo "  - API服务: http://localhost:8080"
echo "  - 健康检查: http://localhost:8080/health"
echo "  - API文档: http://localhost:8080/swagger/index.html"
echo ""

echo "📊 监控命令："
echo "  - 查看日志: docker logs -f qingyu-backend"
echo "  - 查看状态: docker-compose -f docker/docker-compose.prod.yml ps"
echo "  - 停止服务: docker-compose -f docker/docker-compose.prod.yml down"
echo ""

echo "🧪 测试建议："
echo "  1. 运行集成测试: bash scripts/mvp_integration_test.sh"
echo "  2. 运行冒烟测试: bash scripts/mvp_smoke_test.sh"
echo "  3. 手动功能验证（参考部署指南）"
echo ""

echo "📖 文档参考："
echo "  - 部署指南: doc/ops/MVP部署指南_2025-10-23.md"
echo "  - API文档: doc/api/API接口总览.md"
echo "  - 内测指南: doc/usage/内测指南_v1.0.md"
echo ""

echo -e "${BLUE}🎉 祝部署顺利！${NC}"

