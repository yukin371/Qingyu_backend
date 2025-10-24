# 青羽后端集成测试套件

## 📋 概述

本目录包含青羽后端的完整集成测试套件，涵盖从数据准备到业务流程的所有测试场景。

## 🗂️ 测试文件结构

```
test/integration/
├── scenario_bookstore_test.go      # 书城首页和榜单测试
├── scenario_search_test.go         # 书籍搜索功能测试
├── scenario_reading_test.go        # 阅读流程测试
├── scenario_ai_generation_test.go  # AI文本生成测试
├── scenario_auth_test.go           # 用户认证和权限测试
├── scenario_writing_test.go        # 写作流程测试
├── scenario_interaction_test.go    # 互动功能测试
└── README_集成测试说明.md         # 本文件
```

## 🚀 快速开始

### 一键执行（推荐）

```batch
scripts\testing\run_integration_tests.bat
```

按提示操作：
1. 选择是否准备测试数据（首次运行选yes）
2. 选择要执行的测试场景（1-8）

### 手动执行

#### 1. 准备环境

```batch
# 清理旧数据
scripts\testing\cleanup_test_data.bat

# 导入测试用户
go run scripts\testing\import_test_users.go

# 导入小说数据
scripts\testing\import_novels.bat
```

#### 2. 启动服务器

```batch
go run cmd/server/main.go
```

#### 3. 运行测试

```batch
# 单个测试
go test -v ./test/integration/scenario_bookstore_test.go

# 全部测试
go test -v ./test/integration/scenario_*.go
```

## 📦 测试场景详情

### 1. 书城流程测试 ✅

**文件**: `scenario_bookstore_test.go`

**测试内容**:
- ✓ 获取书城首页数据
- ✓ 获取分类树
- ✓ 获取推荐书籍
- ✓ 获取精选书籍
- ✓ 获取活动Banner
- ✓ 获取实时榜/周榜/月榜/新人榜
- ✓ 数据库统计验证

### 2. 搜索功能测试 🔍

**文件**: `scenario_search_test.go`

**测试内容**:
- ✓ 按标题关键词搜索
- ✓ 按作者搜索
- ✓ 组合搜索
- ✓ 排序功能（最新、最热、字数）
- ✓ 分页功能
- ✓ 无结果场景处理

### 3. 阅读流程测试 📖

**文件**: `scenario_reading_test.go`

**测试内容**:
- ✓ 获取书籍详情
- ✓ 获取章节列表
- ✓ 阅读章节内容
- ✓ 保存阅读进度
- ✓ 添加书签
- ✓ 添加笔记
- ✓ 管理个人书架

**需要**: 登录测试用户

### 4. AI生成测试 🤖

**文件**: `scenario_ai_generation_test.go`

**测试内容**:
- ✓ 文本续写
- ✓ 文本改写
- ✓ 文本扩写
- ✓ 文本润色
- ✓ Token使用统计
- ✓ 错误处理

**需要**: 
- 配置Gemini API Key
- 网络连接

**API配置**: 
```yaml
# config/config.yaml
external_api:
  default_provider: "gemini"
  providers:
    gemini:
      api_key: "YOUR_API_KEY"
      base_url: "https://generativelanguage.googleapis.com"
      enabled: true
```

### 5. 认证流程测试 🔐

**文件**: `scenario_auth_test.go`

**测试内容**:
- ✓ 普通用户登录
- ✓ VIP用户登录
- ✓ 管理员登录
- ✓ Token验证
- ✓ 权限控制
- ✓ 用户注册
- ✓ 错误处理（错误密码、不存在用户）

### 6. 写作流程测试 ✍️

**文件**: `scenario_writing_test.go`

**测试内容**:
- ✓ 创建写作项目
- ✓ 获取项目列表
- ✓ 创建文档
- ✓ 保存草稿
- ✓ 版本管理
- ✓ 发布文档

**需要**: 登录测试用户

### 7. 互动功能测试 💬

**文件**: `scenario_interaction_test.go`

**测试内容**:
- ✓ 收藏书籍
- ✓ 发表评论
- ✓ 点赞/取消点赞
- ✓ 阅读历史
- ✓ 个人书架
- ✓ 互动数据统计

**需要**: 登录测试用户

## 👥 测试账号

| 角色 | 邮箱 | 密码 | 用途 |
|-----|------|------|-----|
| 管理员 | admin@qingyu.com | Admin@123456 | 管理员测试 |
| VIP用户 | vip01@qingyu.com | Vip@123456 | VIP功能测试 |
| VIP用户 | vip02@qingyu.com | Vip@123456 | VIP功能测试 |
| 普通用户 | test01@qingyu.com | Test@123456 | 普通用户测试 |
| 普通用户 | test02@qingyu.com | Test@123456 | 普通用户测试 |
| 普通用户 | test03@qingyu.com | Test@123456 | 普通用户测试 |
| 普通用户 | test04@qingyu.com | Test@123456 | 普通用户测试 |
| 普通用户 | test05@qingyu.com | Test@123456 | 普通用户测试 |

## 📊 测试数据

### 小说数据
- **来源**: CNNovel125K数据集
- **数量**: 100本小说（约3000章）
- **文件**: `data/novels_100.json`

### 数据准备脚本

```batch
# Windows
scripts\testing\cleanup_test_data.bat      # 清理数据
scripts\testing\import_test_users.go       # 导入用户
scripts\testing\import_novels.bat          # 导入小说
scripts\testing\run_integration_tests.bat  # 一键测试
```

## ⚙️ 配置说明

### 测试环境配置

测试使用独立的配置文件：`config/config.test.yaml`

关键配置：
```yaml
mongodb:
  database: "qingyu_test"  # 使用测试数据库

redis:
  db: 1  # 使用独立的Redis DB

external_api:
  providers:
    gemini:
      api_key: "YOUR_API_KEY"  # Gemini API Key
```

## 🐛 常见问题

### 服务器连接失败

**症状**: 测试提示"无法连接到服务器"

**解决**:
```batch
# 1. 检查服务器是否运行
curl http://localhost:8080/api/v1/system/health

# 2. 启动服务器
go run cmd/server/main.go
```

### 数据库无数据

**症状**: 测试提示"数据库中没有书籍"

**解决**:
```batch
# 重新导入数据
scripts\testing\import_novels.bat
```

### AI测试失败

**症状**: AI测试返回错误

**解决**:
1. 检查API Key配置
2. 检查网络连接
3. 确认API配额

### 登录失败

**症状**: 测试提示"用户不存在"

**解决**:
```batch
# 重新导入测试用户
go run scripts\testing\import_test_users.go
```

## 📈 测试报告

生成测试报告：

```batch
# 运行测试并保存报告
go test -v ./test/integration/scenario_*.go > test_report.txt 2>&1

# 查看报告
type test_report.txt
```

## 🔄 持续集成

测试可集成到CI/CD流程中：

```yaml
# .github/workflows/integration-test.yml
name: Integration Tests
on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    services:
      mongodb:
        image: mongo:4.4
        ports:
          - 27017:27017
    
    steps:
      - uses: actions/checkout@v2
      - uses: actions/setup-go@v2
      
      - name: Prepare test data
        run: |
          go run scripts/testing/import_test_users.go
          bash scripts/testing/import_novels.sh
      
      - name: Run tests
        run: go test -v ./test/integration/...
```

## 📚 相关文档

- [集成测试使用指南](../../doc/testing/集成测试使用指南.md)
- [项目架构文档](../../doc/architecture/)
- [API文档](../../doc/api/)

## ✅ 测试检查清单

运行测试前确认：

- [ ] MongoDB 正在运行
- [ ] 测试数据已导入
- [ ] 测试用户已创建
- [ ] 服务器正在运行
- [ ] Gemini API Key已配置（AI测试需要）
- [ ] 网络连接正常

## 🎯 测试目标

- ✅ 验证所有核心业务流程
- ✅ 确保API接口正常工作
- ✅ 验证用户认证和权限控制
- ✅ 测试AI集成功能
- ✅ 验证数据完整性

---

**最后更新**: 2025-10-24  
**维护者**: 青羽后端测试团队


