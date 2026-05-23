# 🚀 快速开始 - 本地测试数据初始化

## ⚡ 一键运行（推荐）

### Windows 用户

```bash
scripts\setup_local_test_data.bat
```

### Linux/Mac 用户

```bash
# 首次运行需要添加执行权限
chmod +x scripts/setup_local_test_data.sh

# 运行脚本
./scripts/setup_local_test_data.sh
```

---

## ✅ 完成后你将拥有

### 📚 小说数据
- **100 本**网络小说
- **约 3000 章**内容
- **多种分类**：玄幻、言情、都市、武侠、仙侠、科幻等
- **完整数据**：标题、作者、简介、封面、评分、字数等

### 👥 内测账号（20个）

| 类型 | 数量 | 用户名示例 | 密码来源 |
|------|------|-----------|------|
| 管理员 | 3 个 | `admin` | `QINGYU_TEST_ADMIN_PASSWORD` 或初始化输出 |
| VIP用户 | 5 个 | `vip_writer01` | `QINGYU_TEST_VIP_PASSWORD` 或初始化输出 |
| 写作用户 | 5 个 | `writer_xuanhuan` | `QINGYU_TEST_WRITER_PASSWORD` 或初始化输出 |
| 阅读用户 | 5 个 | `reader01` | `QINGYU_TEST_READER_PASSWORD` 或初始化输出 |
| 测试用户 | 2 个 | `tester_api` | `QINGYU_TEST_USER_PASSWORD` 或初始化输出 |

---

## 🎮 立即测试

### 1. 启动服务器

```bash
go run cmd/server/main.go
```

服务器将运行在：`http://localhost:8080`

### 2. 登录测试

选择一个账号登录：

```bash
# 使用管理员账号
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"<从环境变量或初始化输出获取>"}'
```

### 3. 测试书城接口

```bash
# 获取书籍列表
curl -X GET "http://localhost:8080/api/v1/bookstore/books?page=1&page_size=20"

# 搜索书籍
curl -X GET "http://localhost:8080/api/v1/bookstore/search?keyword=玄幻"
```

---

## 📖 详细文档

需要更多信息？查看：

| 文档 | 说明 |
|------|------|
| [README_测试数据初始化.md](README_测试数据初始化.md) | 完整使用指南 |
| [内测账号快速参考](doc/testing/内测账号快速参考.md) | 所有账号详细列表 |
| [本地测试数据初始化指南](doc/testing/本地测试数据初始化指南.md) | 详细步骤说明 |

---

## ❓ 遇到问题？

### MongoDB 连接失败

```bash
# Windows 启动 MongoDB
net start MongoDB

# Linux/Mac 启动 MongoDB
sudo systemctl start mongod
```

### 数据已存在

脚本会自动跳过已存在的数据，安全运行。

### 需要清理数据

```bash
go run cmd/migrate/main.go -command=clean-novels
```

---

## 🎯 下一步

初始化完成后，你可以：

1. ✅ 测试写作端功能（使用 `writer_*` 账号）
2. ✅ 测试阅读端功能（使用 `reader_*` 账号）
3. ✅ 测试管理端功能（使用 `admin` 账号）
4. ✅ 测试 VIP 功能（使用 `vip_*` 账号）
5. ✅ 进行 API 接口测试（使用 `tester_api` 账号）

---

**🎉 祝测试愉快！**

如有问题，请查看详细文档或提交 Issue。

