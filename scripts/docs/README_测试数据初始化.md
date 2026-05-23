# 青羽写作系统 - 测试数据快速初始化

## 🎯 一键初始化

### Windows 系统

```bash
scripts\setup_local_test_data.bat
```

### Linux/Mac 系统

```bash
# 添加执行权限（首次运行）
chmod +x scripts/setup_local_test_data.sh

# 运行脚本
./scripts/setup_local_test_data.sh
```

---

## 📦 初始化内容

运行脚本后将自动完成：

### ✅ 1. 导入小说数据
- 📚 100 本小说
- 📄 约 3000 章节
- 🏷️ 多种分类（玄幻、言情、都市、武侠等）
- 💯 完整的评分数据

### ✅ 2. 创建内测用户
- 👨‍💼 3 个管理员账号
- 💎 5 个 VIP 用户
- ✍️ 5 个写作用户
- 📖 5 个阅读用户
- 🧪 2 个测试用户
- **总计 20 个测试账号**

---

## 🔑 快速登录

| 类型 | 用户名 | 密码来源 | 用途 |
|------|--------|------|------|
| 管理员 | `admin` | `QINGYU_TEST_ADMIN_PASSWORD` 或本地初始化输出 | 系统管理 |
| VIP作家 | `vip_writer01` | `QINGYU_TEST_VIP_PASSWORD` 或本地初始化输出 | 高级写作 |
| VIP读者 | `vip_reader01` | `QINGYU_TEST_VIP_PASSWORD` 或本地初始化输出 | 付费阅读 |
| 普通作家 | `writer_xuanhuan` | `QINGYU_TEST_WRITER_PASSWORD` 或本地初始化输出 | 玄幻创作 |
| 普通读者 | `reader01` | `QINGYU_TEST_READER_PASSWORD` 或本地初始化输出 | 基础阅读 |

**完整账号列表**: 查看 [内测账号快速参考](doc/testing/内测账号快速参考.md)

---

## 📚 详细文档

- 📖 [本地测试数据初始化指南](doc/testing/本地测试数据初始化指南.md) - 完整使用说明
- 🔑 [内测账号快速参考](doc/testing/内测账号快速参考.md) - 所有测试账号列表
- 📚 [小说导入详细指南](migration/seeds/README_小说导入指南.md) - 数据导入细节

---

## ⚡ 前置要求

### 必须满足：
- ✅ Go 1.21+
- ✅ MongoDB 已启动（端口 27017）
- ✅ 配置文件正确（`config/config.local.yaml`）

### 可选：
- Python 3.7+（仅在需要生成新数据时）

---

## 🚀 启动服务器

初始化完成后，启动服务器：

```bash
go run cmd/server/main.go
```

服务器将运行在：`http://localhost:8080`

---

## 🧪 测试示例

### 1. API 登录测试

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"<从环境变量或初始化输出获取>"}'
```

### 2. 获取书城数据

```bash
curl -X GET "http://localhost:8080/api/v1/bookstore/books?page=1&page_size=20"
```

### 3. 搜索书籍

```bash
curl -X GET "http://localhost:8080/api/v1/bookstore/search?keyword=玄幻"
```

---

## 🔧 进阶操作

### 仅导入小说数据

```bash
go run cmd/migrate/main.go -command=import-novels -file=data/novels_100.json
```

### 仅创建内测用户

```bash
go run cmd/create_beta_users/main.go
```

### 清理测试数据

```bash
go run cmd/migrate/main.go -command=clean-novels
```

### 导入更多小说

```bash
# 生成 500 本小说
python scripts/import_novels.py --max-novels 500 --output data/novels_500.json

# 导入到数据库
go run cmd/migrate/main.go -command=import-novels -file=data/novels_500.json
```

---

## ❓ 常见问题

### Q: MongoDB 连接失败？

**解决方案**：
```bash
# Windows
net start MongoDB

# Linux/Mac
sudo systemctl start mongod
```

### Q: 数据文件不存在？

项目已包含 `data/novels_100.json`，如果缺失可以重新生成：

```bash
python scripts/import_novels.py --max-novels 100 --output data/novels_100.json
```

### Q: 如何重置所有数据？

```bash
# 1. 清理小说数据
go run cmd/migrate/main.go -command=clean-novels

# 2. 重新运行初始化脚本
scripts\setup_local_test_data.bat
```

---

## ⚠️ 安全提醒

### ✅ 开发环境
- 可以使用测试账号
- 定期清理测试数据

### ❌ 生产环境
- **绝对不要**使用这些测试账号
- 使用强密码（至少 16 位）
- 启用双因素认证

---

## 📂 项目结构

```
Qingyu_backend/
├── data/
│   └── novels_100.json           # 100本小说数据
├── scripts/
│   ├── setup_local_test_data.bat # Windows 初始化脚本
│   ├── setup_local_test_data.sh  # Linux/Mac 初始化脚本
│   └── import_novels.py          # Python 数据生成脚本
├── migration/
│   └── seeds/
│       ├── create_beta_users.go  # 内测用户创建脚本
│       └── import_novels.go      # 小说导入脚本
├── doc/
│   └── testing/
│       ├── 本地测试数据初始化指南.md
│       └── 内测账号快速参考.md
└── README_测试数据初始化.md      # 本文档
```

---

## 🎓 学习资源

- [项目架构文档](doc/architecture/)
- [API 接口文档](doc/api/)
- [测试规范](doc/testing/)

---

## 📞 技术支持

遇到问题？
1. 查看 [详细文档](doc/testing/本地测试数据初始化指南.md)
2. 检查日志 `startup.log`
3. 提交 Issue

---

**最后更新**: 2025-10-24  
**维护者**: 青羽后端开发团队  
**版本**: v1.0

