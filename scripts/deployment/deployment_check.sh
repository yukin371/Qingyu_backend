#!/bin/bash

# MVP部署前检查脚本
# 用途：验证部署准备工作是否完成
# 日期：2025-10-23

set -e

echo "======================================"
echo "MVP部署准备检查"
echo "======================================"
echo ""

# 颜色定义
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 检查结果统计
TOTAL_CHECKS=0
PASSED_CHECKS=0
FAILED_CHECKS=0
WARNING_CHECKS=0

# 检查结果记录函数
check_result() {
    local check_name=$1
    local result=$2
    local message=$3

    TOTAL_CHECKS=$((TOTAL_CHECKS + 1))

    if [ "$result" = "PASS" ]; then
        echo -e "${GREEN}✓${NC} $check_name"
        PASSED_CHECKS=$((PASSED_CHECKS + 1))
    elif [ "$result" = "WARN" ]; then
        echo -e "${YELLOW}⚠${NC} $check_name"
        [ -n "$message" ] && echo -e "  ${YELLOW}→${NC} $message"
        WARNING_CHECKS=$((WARNING_CHECKS + 1))
    else
        echo -e "${RED}✗${NC} $check_name"
        [ -n "$message" ] && echo -e "  ${RED}→${NC} $message"
        FAILED_CHECKS=$((FAILED_CHECKS + 1))
    fi
}

# 1. 代码完整性检查
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "1. 代码完整性检查"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

# 检查关键文件存在性
if [ -f "service/project/autosave_service.go" ]; then
    check_result "自动保存服务文件" "PASS"
else
    check_result "自动保存服务文件" "FAIL" "文件不存在"
fi

if [ -f "service/shared/auth/password_validator.go" ]; then
    check_result "密码验证器文件" "PASS"
else
    check_result "密码验证器文件" "FAIL" "文件不存在"
fi

if [ -f "service/shared/auth/session_service.go" ]; then
    check_result "会话服务文件" "PASS"
else
    check_result "会话服务文件" "FAIL" "文件不存在"
fi

echo ""

# 2. 编译检查
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "2. 编译检查"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

# 清理之前的构建
rm -f qingyu_backend qingyu_backend.exe

# 编译主程序
if go build -o qingyu_backend ./cmd/server/main.go 2>/dev/null; then
    check_result "主程序编译" "PASS"
else
    check_result "主程序编译" "FAIL" "编译失败"
fi

# 检查可执行文件
if [ -f "qingyu_backend" ] || [ -f "qingyu_backend.exe" ]; then
    check_result "可执行文件生成" "PASS"
else
    check_result "可执行文件生成" "FAIL"
fi

echo ""

# 3. 依赖检查
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "3. 依赖检查"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

# 检查go.mod
if [ -f "go.mod" ]; then
    check_result "go.mod存在" "PASS"
else
    check_result "go.mod存在" "FAIL"
fi

# 检查依赖完整性
if go mod tidy 2>/dev/null && go mod verify 2>/dev/null; then
    check_result "依赖完整性" "PASS"
else
    check_result "依赖完整性" "FAIL" "请运行 go mod tidy"
fi

echo ""

# 4. 配置文件检查
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "4. 配置文件检查"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

# 检查配置文件
if [ -f "config/config.yaml" ]; then
    check_result "生产配置文件" "PASS"
else
    check_result "生产配置文件" "FAIL"
fi

if [ -f "config/config.test.yaml" ]; then
    check_result "测试配置文件" "PASS"
else
    check_result "测试配置文件" "WARN" "建议创建测试配置"
fi

if [ -f "config/config.docker.yaml" ]; then
    check_result "Docker配置文件" "PASS"
else
    check_result "Docker配置文件" "WARN" "建议创建Docker配置"
fi

echo ""

# 5. Docker配置检查
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "5. Docker配置检查"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

# 检查Dockerfile
if [ -f "docker/Dockerfile.prod" ]; then
    check_result "生产Dockerfile" "PASS"
else
    check_result "生产Dockerfile" "WARN" "建议创建生产Dockerfile"
fi

# 检查docker-compose
if [ -f "docker/docker-compose.prod.yml" ]; then
    check_result "生产docker-compose" "PASS"
else
    check_result "生产docker-compose" "WARN" "建议创建生产compose配置"
fi

# 检查Docker环境
if command -v docker &> /dev/null; then
    DOCKER_VERSION=$(docker --version | awk '{print $3}' | sed 's/,//')
    check_result "Docker环境 ($DOCKER_VERSION)" "PASS"
else
    check_result "Docker环境" "FAIL" "Docker未安装"
fi

# 检查docker-compose
if command -v docker-compose &> /dev/null; then
    check_result "docker-compose工具" "PASS"
else
    check_result "docker-compose工具" "WARN" "建议安装docker-compose"
fi

echo ""

# 6. 文档检查
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "6. 文档检查"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

# 检查MVP文档
if [ -f "doc/implementation/00进度指导/MVP收尾完成报告_2025-10-23.md" ]; then
    check_result "MVP完成报告" "PASS"
else
    check_result "MVP完成报告" "WARN"
fi

# 检查部署文档
if [ -f "doc/ops/部署指南.md" ]; then
    check_result "部署指南" "PASS"
else
    check_result "部署指南" "WARN" "建议创建部署文档"
fi

# 检查API文档
if [ -f "doc/api/API接口总览.md" ]; then
    check_result "API文档" "PASS"
else
    check_result "API文档" "WARN"
fi

echo ""

# 7. 安全检查
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "7. 安全检查"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

# 检查敏感信息
if grep -r "password.*=.*\".*\"" config/*.yaml 2>/dev/null | grep -v "# " | head -1 | grep -q .; then
    check_result "配置文件密码明文" "WARN" "建议使用环境变量"
else
    check_result "配置文件密码明文" "PASS"
fi

# 检查.gitignore
if [ -f ".gitignore" ] && grep -q "config.local.yaml" .gitignore; then
    check_result ".gitignore配置" "PASS"
else
    check_result ".gitignore配置" "WARN" "建议忽略敏感配置文件"
fi

echo ""

# 8. 数据库准备检查
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "8. 数据库准备检查"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

# 检查数据库迁移脚本
if [ -f "cmd/migrate/main.go" ]; then
    check_result "数据库迁移脚本" "PASS"
else
    check_result "数据库迁移脚本" "WARN"
fi

# 检查种子数据脚本
if [ -d "migration/seeds" ]; then
    check_result "种子数据脚本" "PASS"
else
    check_result "种子数据脚本" "WARN"
fi

echo ""

# 9. 检查清单总结
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "检查清单总结"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

echo "总计检查: $TOTAL_CHECKS"
echo -e "通过: ${GREEN}$PASSED_CHECKS${NC}"
echo -e "警告: ${YELLOW}$WARNING_CHECKS${NC}"
echo -e "失败: ${RED}$FAILED_CHECKS${NC}"

echo ""

# 部署建议
if [ $FAILED_CHECKS -eq 0 ]; then
    echo -e "${GREEN}✓ 部署准备检查通过！${NC}"
    echo ""
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo "部署步骤建议"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo ""
    echo "1. 构建Docker镜像："
    echo -e "   ${BLUE}docker build -t qingyu-backend:mvp -f docker/Dockerfile.prod .${NC}"
    echo ""
    echo "2. 启动服务："
    echo -e "   ${BLUE}docker-compose -f docker/docker-compose.prod.yml up -d${NC}"
    echo ""
    echo "3. 数据库初始化："
    echo -e "   ${BLUE}docker exec -it qingyu-backend go run cmd/migrate/main.go${NC}"
    echo ""
    echo "4. 健康检查："
    echo -e "   ${BLUE}curl http://localhost:8080/health${NC}"
    echo ""
    echo "5. 查看日志："
    echo -e "   ${BLUE}docker logs -f qingyu-backend${NC}"
    echo ""

    if [ $WARNING_CHECKS -gt 0 ]; then
        echo -e "${YELLOW}注意: 有 $WARNING_CHECKS 个警告项，建议部署前处理${NC}"
    fi

    exit 0
else
    echo -e "${RED}✗ 部署准备检查失败！${NC}"
    echo -e "${RED}请修复上述失败项后重试${NC}"
    exit 1
fi

