# Python测试脚本说明

## 📋 概述

为避免Windows批处理文件和Go文件中的中文编码问题，我们提供了Python版本的测试脚本。这些脚本确保在所有平台上都能正确处理中文字符。

## 🐍 Python脚本列表

| 脚本 | 功能 | 说明 |
|-----|------|-----|
| `setup_integration_tests.py` | 一键测试准备 | 完整的测试数据准备和执行工具 |
| `cleanup_database.py` | 数据库清理 | 清理测试数据（替代.bat版本） |
| `run_tests.py` | 运行测试 | 执行集成测试（替代.bat版本） |

## 🚀 使用方法

### 前提条件

```bash
# 安装Python依赖
pip install requests
```

### 方法1：一键准备和测试（推荐）

```bash
# Windows
python scripts/testing/setup_integration_tests.py

# Linux/Mac
python3 scripts/testing/setup_integration_tests.py
```

这个脚本会：
1. 询问是否准备测试数据
2. 清理旧数据
3. 导入测试用户
4. 导入小说数据
5. 检查服务器状态
6. 运行集成测试

### 方法2：分步执行

#### 1. 清理数据库

```bash
python scripts/testing/cleanup_database.py
```

#### 2. 导入测试用户

```bash
# 这个仍然使用Go脚本（因为需要连接MongoDB）
go run scripts/testing/import_test_users.go
```

#### 3. 导入小说数据

```bash
# 使用Go脚本
go run cmd/migrate/main.go -command=import-novels -file=data/novels_100.json
```

#### 4. 运行测试

```bash
python scripts/testing/run_tests.py
```

## 📚 详细功能说明

### setup_integration_tests.py

**主要功能**：
- ✅ 自动检查MongoDB服务
- ✅ 交互式数据清理
- ✅ 自动导入测试用户
- ✅ 自动导入小说数据
- ✅ 检查和启动服务器
- ✅ 选择性运行测试
- ✅ 彩色输出，易于阅读

**使用示例**：

```bash
$ python scripts/testing/setup_integration_tests.py

============================================================
青羽后端 - 集成测试准备工具
============================================================

项目根目录: E:\Github\Qingyu\Qingyu_backend

是否需要准备测试数据？(首次运行选yes) [yes/no]: yes

============================================================
第1步: 清理测试数据
============================================================

➜ 检查 MongoDB 服务...
✓ MongoDB 服务正常运行
确认清理测试数据？这将删除所有测试数据 (yes/no): yes
...
```

### cleanup_database.py

**主要功能**：
- ✅ 清理所有测试数据
- ✅ 保留系统必需用户
- ✅ 显示清理前后统计
- ✅ 安全确认机制

**使用示例**：

```bash
$ python scripts/testing/cleanup_database.py

============================================================
青羽后端 - 测试数据清理
============================================================

[1/2] 检查 MongoDB 服务...
[成功] MongoDB 服务正常运行

[警告] 此操作将清空以下数据:
  - books (书籍)
  - chapters (章节)
  ...

确认清理？(yes/no): yes
```

### run_tests.py

**主要功能**：
- ✅ 自动检查服务器状态
- ✅ 可选启动服务器
- ✅ 交互式选择测试场景
- ✅ 显示测试结果
- ✅ 查看服务器日志

**使用示例**：

```bash
$ python scripts/testing/run_tests.py

============================================================
青羽后端 - 集成测试运行工具
============================================================

✓ 服务器正在运行

可用的测试场景:
  1. 书城流程测试
  2. 搜索功能测试
  ...

请选择要执行的测试 (1-8): 1
```

## 🎨 特色功能

### 1. 彩色输出

Python脚本使用ANSI颜色代码提供友好的输出：

- 🟢 绿色：成功信息
- 🔵 蓝色：提示信息
- 🟡 黄色：警告信息
- 🔴 红色：错误信息
- 🟣 紫色：标题

### 2. 跨平台支持

所有脚本在Windows、Linux和Mac上都能正常运行，自动处理：
- 命令差异
- 路径分隔符
- 文件编码
- 服务启动方式

### 3. 智能错误处理

- 自动检测MongoDB连接
- 验证文件存在性
- 超时控制
- 友好的错误提示

### 4. 交互式操作

- 确认危险操作
- 选择测试场景
- 可选功能
- 实时反馈

## 🔧 配置说明

### MongoDB连接

脚本默认连接到：
- 数据库：`qingyu_test`
- 连接：通过`mongosh`命令

如需修改，编辑脚本中的连接字符串。

### 服务器端口

默认端口：`8080`

修改位置：
```python
response = requests.get("http://localhost:8080/api/v1/system/health")
```

## 🐛 故障排除

### 问题1：找不到mongosh命令

**解决方法**：
1. 确保MongoDB已安装
2. 将MongoDB bin目录添加到PATH
3. 或使用完整路径

### 问题2：编码错误

**解决方法**：
Python脚本已强制使用UTF-8编码：
```python
encoding='utf-8'
```

### 问题3：权限错误

**Linux/Mac**：
```bash
chmod +x scripts/testing/*.py
```

### 问题4：导入失败

**检查**：
1. MongoDB服务是否运行
2. 数据文件是否存在
3. Go环境是否正确

## 📊 对比：Python vs 批处理

| 特性 | Python脚本 | 批处理脚本 |
|-----|-----------|-----------|
| 编码支持 | ✅ UTF-8 | ⚠️ 可能有问题 |
| 跨平台 | ✅ 是 | ❌ 仅Windows |
| 彩色输出 | ✅ 是 | ❌ 否 |
| 错误处理 | ✅ 完善 | ⚠️ 基础 |
| 交互性 | ✅ 强 | ⚠️ 一般 |
| 依赖 | Python 3.6+ | Windows |

## 💡 推荐使用

**推荐顺序**：

1. **首次使用**：`setup_integration_tests.py` 一键完成所有操作
2. **日常测试**：`run_tests.py` 快速运行测试
3. **重置数据**：`cleanup_database.py` 清理后重新导入

**最佳实践**：

```bash
# 完整流程（首次）
python scripts/testing/setup_integration_tests.py

# 日常测试
python scripts/testing/run_tests.py

# 重置环境
python scripts/testing/cleanup_database.py
python scripts/testing/setup_integration_tests.py
```

## 📝 注意事项

1. ⚠️ Python脚本需要`requests`库
2. ⚠️ 确保MongoDB服务已启动
3. ⚠️ 数据清理操作不可逆
4. ⚠️ 建议在测试环境使用

## 🔗 相关文档

- [集成测试使用指南](../../doc/testing/集成测试使用指南.md)
- [测试套件说明](../../test/integration/README_集成测试说明.md)

---

**最后更新**: 2025-10-24  
**维护者**: 青羽后端测试团队



