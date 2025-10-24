# 青羽后端 - 脚本工具集

本目录包含项目开发、测试、部署相关的自动化脚本。

## 📁 目录结构

```
scripts/
├── init/           # 初始化和环境设置脚本
├── testing/        # 测试相关脚本
├── deployment/     # 部署相关脚本
├── data/           # 数据处理脚本
├── utils/          # 工具脚本
└── docs/           # 文档
```

## 🚀 快速开始

### 本地开发环境初始化

推荐使用 Python 版本（跨平台兼容）：

```bash
# 初始化本地测试数据（推荐）
python scripts/init/setup_local_test_data.py

# 或使用平台特定脚本
scripts/init/setup_local_test_data.bat    # Windows
scripts/init/setup_local_test_data.sh     # Linux/Mac
```

### 运行测试

```bash
# 快速验证（推荐）
python scripts/testing/quick_verify.py

# 完整测试套件
python scripts/testing/run_tests.py
```

## 📂 分类说明

### init/ - 初始化脚本

环境搭建和初始化相关脚本。

| 脚本 | 说明 | 平台 |
|------|------|------|
| `setup_local_test_data.py` | 本地测试数据初始化（推荐） | 跨平台 |
| `setup_local_test_data.bat` | 本地测试数据初始化 | Windows |
| `setup_local_test_data.sh` | 本地测试数据初始化 | Linux/Mac |
| `setup-test-env.sh` | 设置测试环境 | Linux/Mac |

**使用示例**：
```bash
# 推荐：使用 Python 版本
python scripts/init/setup_local_test_data.py

# Windows
scripts\init\setup_local_test_data.bat

# Linux/Mac
./scripts/init/setup_local_test_data.sh
```

### testing/ - 测试脚本

自动化测试相关脚本。

| 脚本 | 说明 | 平台 |
|------|------|------|
| `quick_verify.py` | 快速验证（推荐） | 跨平台 |
| `quick_verify.bat/sh` | 快速验证 | Windows/Linux |
| `run_tests.py` | 运行测试套件（推荐） | 跨平台 |
| `run_tests.sh` | 运行测试套件 | Linux/Mac |
| `run_tests_with_docker.bat/sh` | Docker 环境测试 | Windows/Linux |
| `test_reading_features.py` | 阅读功能测试 | 跨平台 |
| `test_reading_features.bat/sh` | 阅读功能测试 | Windows/Linux |
| `mvp_smoke_test.sh` | MVP 冒烟测试 | Linux/Mac |
| `mvp_integration_test.sh` | MVP 集成测试 | Linux/Mac |
| `验证项目修复.bat` | 验证项目修复 | Windows |

**使用示例**：
```bash
# 快速验证（推荐）
python scripts/testing/quick_verify.py

# 运行完整测试
python scripts/testing/run_tests.py

# 测试阅读功能
python scripts/testing/test_reading_features.py
```

### deployment/ - 部署脚本

部署和发布相关脚本。

| 脚本 | 说明 | 平台 |
|------|------|------|
| `quick_deploy_mvp.sh` | 快速部署 MVP | Linux/Mac |
| `deployment_check.sh` | 部署检查 | Linux/Mac |

**使用示例**：
```bash
./scripts/deployment/quick_deploy_mvp.sh
./scripts/deployment/deployment_check.sh
```

### data/ - 数据处理脚本

数据导入、处理相关脚本。

| 脚本 | 说明 | 平台 |
|------|------|------|
| `import_novels.py` | 导入小说数据 | 跨平台 |
| `test_novel_import.bat/sh` | 测试小说导入 | Windows/Linux |

**使用示例**：
```bash
# 导入 100 本小说
python scripts/data/import_novels.py --max-novels 100 --output data/novels_100.json

# 测试导入功能
python scripts/data/test_novel_import.py
```

### utils/ - 工具脚本

通用工具脚本。

| 脚本 | 说明 | 平台 |
|------|------|------|
| `fix_swagger_types.py` | 修复 Swagger 类型定义 | 跨平台 |

**使用示例**：
```bash
python scripts/utils/fix_swagger_types.py
```

## 🐍 Python 版本 vs Shell/Batch 版本

### 推荐使用 Python 版本

**优点**：
- ✅ 跨平台兼容（Windows/Linux/Mac）
- ✅ 更好的错误处理
- ✅ 丰富的标准库支持
- ✅ 易于维护和扩展

**要求**：
- Python 3.7+
- 安装依赖：`pip install -r requirements.txt`（如有）

### Shell/Batch 版本

**使用场景**：
- 系统级脚本调用
- CI/CD 集成
- 特定平台优化

## 📋 常用命令速查

### 初始化开发环境

```bash
# 1. 初始化测试数据（推荐）
python scripts/init/setup_local_test_data.py

# 2. 快速验证
python scripts/testing/quick_verify.py

# 3. 启动服务
go run cmd/server/main.go
```

### 测试工作流

```bash
# 1. 运行单元测试
python scripts/testing/run_tests.py

# 2. 测试特定功能
python scripts/testing/test_reading_features.py

# 3. Docker 环境测试
scripts/testing/run_tests_with_docker.sh
```

### 数据管理

```bash
# 1. 导入小说数据
python scripts/data/import_novels.py --max-novels 100

# 2. 测试导入功能
python scripts/data/test_novel_import.py

# 3. 创建测试用户
go run cmd/create_beta_users/main.go
```

## 🔧 开发指南

### 添加新脚本

1. **确定分类**：选择合适的目录（init/testing/deployment/data/utils）
2. **提供多版本**：
   - 优先提供 Python 版本（跨平台）
   - 可选提供 Shell/Batch 版本（特定场景）
3. **添加文档**：在本 README 中添加说明

### 脚本命名规范

- **Python**：`snake_case.py`
- **Shell**：`kebab-case.sh`
- **Batch**：`snake_case.bat`

### 脚本模板

#### Python 脚本模板

```python
#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
脚本简要说明

使用方法:
    python script_name.py [options]
"""

import argparse
import sys
import os

def main():
    """主函数"""
    parser = argparse.ArgumentParser(description='脚本说明')
    parser.add_argument('--option', type=str, help='选项说明')
    args = parser.parse_args()
    
    # 脚本逻辑
    print("脚本执行中...")

if __name__ == '__main__':
    try:
        main()
    except Exception as e:
        print(f"错误: {e}", file=sys.stderr)
        sys.exit(1)
```

## 📚 相关文档

- [本地测试数据初始化指南](docs/README_测试数据初始化.md)
- [快速开始指南](docs/QUICKSTART_测试数据.md)
- [内测账号快速参考](../doc/testing/内测账号快速参考.md)

## ⚠️ 注意事项

1. **环境要求**：
   - Python 3.7+ （Python 脚本）
   - Go 1.21+ （Go 相关脚本）
   - MongoDB 运行中

2. **权限问题**：
   - Linux/Mac 脚本需要执行权限：`chmod +x script.sh`
   - Windows 可能需要管理员权限

3. **路径问题**：
   - 所有脚本都应在项目根目录运行
   - 使用相对路径时注意当前工作目录

4. **安全提示**：
   - 测试脚本仅用于开发环境
   - 生产环境请使用正式部署流程
   - 不要在脚本中硬编码敏感信息

## 🐛 问题排查

### 常见问题

1. **Python 脚本执行失败**
   ```bash
   # 检查 Python 版本
   python --version  # 应该是 3.7+
   
   # 使用 python3 命令
   python3 scripts/init/setup_local_test_data.py
   ```

2. **MongoDB 连接失败**
   ```bash
   # 检查 MongoDB 服务状态
   # Windows
   net start MongoDB
   
   # Linux/Mac
   sudo systemctl status mongod
   ```

3. **权限错误（Linux/Mac）**
   ```bash
   # 添加执行权限
   chmod +x scripts/init/*.sh
   chmod +x scripts/testing/*.sh
   ```

## 📞 技术支持

遇到问题？
1. 查看对应脚本的帮助信息：`python script.py --help`
2. 查看相关文档
3. 提交 Issue

---

**最后更新**: 2025-10-24  
**维护者**: 青羽后端开发团队
