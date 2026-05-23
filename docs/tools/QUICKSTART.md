# 🚀 青羽后端脚本 - 快速开始

## ⚡ 常用操作

### 1. 初始化本地开发环境

**推荐使用 Python 版本**（跨平台兼容）：

```bash
# 初始化测试数据（导入小说 + 创建用户）
python scripts/init/setup_local_test_data.py

# 仅导入小说数据
python scripts/init/setup_local_test_data.py --skip-users

# 仅创建用户
python scripts/init/setup_local_test_data.py --skip-novels

# 使用自定义数据文件
python scripts/init/setup_local_test_data.py --novels-file data/novels_500.json
```

**或使用平台特定脚本**：

```bash
# Windows
scripts\init\setup_local_test_data.bat

# Linux/Mac
./scripts/init/setup_local_test_data.sh
```

### 2. 快速验证项目状态

```bash
# 推荐：Python 版本
python scripts/testing/quick_verify.py

# 详细输出
python scripts/testing/quick_verify.py -v

# 跳过编译检查（加快验证）
python scripts/testing/quick_verify.py --skip-build
```

**或使用平台特定脚本**：

```bash
# Windows
scripts\testing\quick_verify.bat

# Linux/Mac
./scripts/testing/quick_verify.sh
```

### 3. 运行测试

```bash
# 推荐：Python 版本
python scripts/testing/run_tests.py

# 生成覆盖率报告
python scripts/testing/run_tests.py --coverage

# 测试特定包
python scripts/testing/run_tests.py --package=./service/...

# 详细输出
python scripts/testing/run_tests.py -v
```

### 4. 数据管理

```bash
# 导入小说数据（从 Hugging Face）
python scripts/data/import_novels.py --max-novels 100 --output data/novels_100.json

# 测试小说导入
python scripts/data/test_novel_import.py

# 仅验证数据文件
python scripts/data/test_novel_import.py --validate-only
```

## 🎯 典型工作流

### 新环境搭建

```bash
# 1. 快速验证环境
python scripts/testing/quick_verify.py

# 2. 初始化测试数据
python scripts/init/setup_local_test_data.py

# 3. 启动服务器
go run cmd/server/main.go

# 4. 测试功能
python scripts/testing/test_reading_features.py
```

### 日常开发

```bash
# 1. 快速验证
python scripts/testing/quick_verify.py

# 2. 运行测试
python scripts/testing/run_tests.py

# 3. 提交代码前
python scripts/testing/run_tests.py --coverage
```

### 数据管理

```bash
# 导入更多测试数据
python scripts/data/import_novels.py --max-novels 500 --output data/novels_500.json
python scripts/data/test_novel_import.py --file data/novels_500.json

# 创建额外的测试用户
go run cmd/create_beta_users/main.go
```

## 🐍 为什么推荐 Python 版本？

### 优点

✅ **跨平台兼容** - Windows/Linux/Mac 都能用  
✅ **更好的错误处理** - 清晰的错误信息  
✅ **丰富的功能** - 命令行参数、进度显示、彩色输出  
✅ **易于维护** - 代码清晰，易于扩展

### 要求

- Python 3.7+
- 标准库即可（无需额外依赖）

### 使用建议

1. **优先使用** Python 版本脚本
2. **保留** Shell/Batch 版本作为备选
3. **自动化** 任务使用 Python 版本

## 📚 脚本分类

### init/ - 初始化脚本

| 脚本 | 说明 |
|------|------|
| `setup_local_test_data.py` | 本地测试数据初始化（推荐） |
| `setup_local_test_data.bat/sh` | 平台特定版本 |
| `setup-test-env.sh` | 测试环境设置 |

### testing/ - 测试脚本

| 脚本 | 说明 |
|------|------|
| `quick_verify.py` | 快速验证（推荐） |
| `run_tests.py` | 运行测试（推荐） |
| `test_reading_features.py` | 阅读功能测试 |
| `mvp_smoke_test.sh` | MVP 冒烟测试 |

### deployment/ - 部署脚本

| 脚本 | 说明 |
|------|------|
| `quick_deploy_mvp.sh` | 快速部署 MVP |
| `deployment_check.sh` | 部署检查 |

### data/ - 数据处理脚本

| 脚本 | 说明 |
|------|------|
| `import_novels.py` | 从 Hugging Face 导入小说 |
| `test_novel_import.py` | 测试小说导入（推荐） |

### utils/ - 工具脚本

| 脚本 | 说明 |
|------|------|
| `fix_swagger_types.py` | 修复 Swagger 类型 |

## ❓ 常见问题

### Q: 如何选择使用哪个版本的脚本？

**A**: 优先使用 Python 版本（`.py`）：
- 跨平台兼容
- 功能更丰富
- 错误处理更好

仅在以下情况使用 Shell/Batch 版本：
- Python 不可用
- CI/CD 集成需要
- 系统级脚本调用

### Q: Python 脚本运行失败怎么办？

**A**: 检查 Python 版本：

```bash
# 检查版本（需要 3.7+）
python --version

# 或尝试使用 python3
python3 scripts/init/setup_local_test_data.py
```

### Q: Shell 脚本权限错误？

**A**: 添加执行权限（Linux/Mac）：

```bash
chmod +x scripts/init/*.sh
chmod +x scripts/testing/*.sh
chmod +x scripts/deployment/*.sh
```

### Q: 如何添加新脚本？

**A**: 
1. 确定分类（init/testing/deployment/data/utils）
2. 优先创建 Python 版本
3. 更新 `README.md`
4. 如需要，创建 Shell/Batch 版本

## 📞 获取帮助

### 查看脚本帮助

所有 Python 脚本都支持 `--help` 参数：

```bash
python scripts/init/setup_local_test_data.py --help
python scripts/testing/quick_verify.py --help
python scripts/testing/run_tests.py --help
python scripts/data/test_novel_import.py --help
```

### 查看文档

- [完整 README](README.md) - 详细文档
- [测试数据初始化指南](../testing/本地测试数据初始化指南.md)
- [项目测试文档](../testing/README.md)

---

**最后更新**: 2025-10-24  
**维护者**: 青羽后端开发团队

