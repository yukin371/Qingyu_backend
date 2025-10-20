# CNNovel125K 测试快速开始

## 🚀 一键测试（推荐）

### Windows 用户
```bash
# 运行自动化测试脚本
scripts\test_novel_import.bat
```

### Linux/Mac 用户
```bash
# 添加执行权限
chmod +x scripts/test_novel_import.sh

# 运行自动化测试脚本
./scripts/test_novel_import.sh
```

脚本会自动完成：
1. ✓ 检查 Python 环境
2. ✓ 安装依赖（datasets）
3. ✓ 从 Hugging Face 加载数据（100本小说）
4. ✓ 验证数据格式
5. ✓ 导入到 MongoDB

---

## 📋 手动步骤（详细）

### 前置要求

确保已安装：
- ✅ Python 3.7+
- ✅ Go 1.21+
- ✅ MongoDB（已启动）

### 步骤 1: 准备环境

#### 切换到 test 分支
```bash
git checkout test
git pull origin test
```

#### 安装 Python 依赖
```bash
pip install datasets
```

### 步骤 2: 加载数据

#### 快速测试（100本）
```bash
python scripts/import_novels.py --max-novels 100 --output data/novels_100.json
```

#### 中等规模（500本）
```bash
python scripts/import_novels.py --max-novels 500 --output data/novels_500.json
```

#### 大规模测试（1000本）
```bash
python scripts/import_novels.py --max-novels 1000 --output data/novels_1000.json
```

⏱️ **预计时间**:
- 100本: 2-5 分钟（首次下载数据集会更久）
- 500本: 5-10 分钟
- 1000本: 10-20 分钟

### 步骤 3: 验证数据

试运行模式（不写入数据库）：
```bash
go run cmd/migrate/main.go -command=import-novels -file=data/novels_100.json -dry-run=true
```

✅ 看到 "所有数据验证通过" 即可继续。

### 步骤 4: 导入数据

**确保 MongoDB 正在运行！**

```bash
go run cmd/migrate/main.go -command=import-novels -file=data/novels_100.json
```

导入完成后会显示：
```
✓ 索引创建成功

数据库统计:
  书籍总数: 100
  章节总数: ~3000
```

### 步骤 5: 验证导入

#### 方法 1: MongoDB 直接查询
```bash
mongo
use qingyu

# 查看书籍数量
db.books.count()

# 查看章节数量
db.chapters.count()

# 查看示例书籍
db.books.findOne()

# 查看示例章节
db.chapters.findOne()
```

#### 方法 2: 启动服务器测试
```bash
# 启动服务器
go run cmd/server/main.go

# 在另一个终端测试 API
curl http://localhost:8080/api/v1/bookstore/books
```

---

## 🧪 功能测试

### 测试书店功能

#### 1. 获取书籍列表
```bash
curl "http://localhost:8080/api/v1/bookstore/books?page=1&pageSize=10"
```

#### 2. 搜索书籍
```bash
curl "http://localhost:8080/api/v1/bookstore/books/search?keyword=测试"
```

#### 3. 按分类筛选
```bash
curl "http://localhost:8080/api/v1/bookstore/books?category=玄幻"
```

#### 4. 获取书籍详情
```bash
# 替换 {book_id} 为实际的书籍ID
curl "http://localhost:8080/api/v1/bookstore/books/{book_id}"
```

### 测试阅读功能

#### 1. 获取章节列表
```bash
# 替换 {book_id} 为实际的书籍ID
curl "http://localhost:8080/api/v1/bookstore/books/{book_id}/chapters"
```

#### 2. 获取章节内容
```bash
# 替换 {chapter_id} 为实际的章节ID
curl "http://localhost:8080/api/v1/bookstore/chapters/{chapter_id}"
```

### 测试推荐功能

#### 1. 获取推荐书籍
```bash
curl "http://localhost:8080/api/v1/recommendation/books"
```

#### 2. 获取热门书籍
```bash
curl "http://localhost:8080/api/v1/bookstore/books?is_hot=true"
```

---

## 🧹 清理测试数据

### 清理所有数据
```bash
go run cmd/migrate/main.go -command=clean-novels
```

系统会要求确认：
```
⚠️  WARNING: This will delete ALL books and chapters!
Are you sure? (yes/no):
```

输入 `yes` 确认删除。

### 按分类清理
```bash
go run cmd/migrate/main.go -command=clean-novels -category=玄幻
```

---

## 📊 性能测试

### 使用 Apache Bench

测试书籍列表性能：
```bash
ab -n 1000 -c 10 http://localhost:8080/api/v1/bookstore/books
```

参数说明：
- `-n 1000`: 总请求数
- `-c 10`: 并发数

### 查看 MongoDB 性能

```javascript
// 分析查询性能
db.books.find({}).explain("executionStats")

// 查看索引使用情况
db.books.getIndexes()

// 查看集合统计
db.books.stats()
```

---

## ❓ 常见问题

### Q1: Python 脚本报错 "ModuleNotFoundError: No module named 'datasets'"
```bash
# 安装 datasets 库
pip install datasets
```

### Q2: 导入时报错 "failed to connect to MongoDB"
```bash
# 确保 MongoDB 正在运行
# Windows: 检查服务
services.msc

# Linux/Mac: 启动 MongoDB
sudo systemctl start mongod
```

### Q3: 数据加载很慢
首次运行会下载 CNNovel125K 数据集，需要时间。数据会缓存在本地，后续会快很多。

### Q4: 如何修改每章字数？
```bash
python scripts/import_novels.py --chapter-size 5000 --max-novels 100
```

### Q5: 如何只导入特定分类？
需要修改 Python 脚本，在 `load_and_process()` 中添加分类过滤逻辑。

---

## 📝 测试检查清单

使用此清单确保完整测试：

### 数据导入
- [ ] Python 脚本成功运行
- [ ] JSON 文件生成
- [ ] 数据验证通过
- [ ] MongoDB 导入成功
- [ ] 索引创建成功

### 书店功能
- [ ] 书籍列表正常显示
- [ ] 分页功能正常
- [ ] 搜索功能正常
- [ ] 分类筛选正常
- [ ] 书籍详情正常

### 阅读功能
- [ ] 章节列表正常
- [ ] 章节内容正常
- [ ] 书签功能正常
- [ ] 阅读历史正常

### 推荐功能
- [ ] 推荐列表正常
- [ ] 热门书籍正常
- [ ] 精选书籍正常

### 性能
- [ ] 列表查询 < 100ms
- [ ] 搜索响应 < 200ms
- [ ] 章节加载 < 100ms

---

## 📚 相关文档

- [详细导入指南](../../migration/seeds/README_小说导入指南.md)
- [测试报告模板](./CNNovel125K测试报告.md)
- [API 文档](../api/)

---

## 🆘 需要帮助？

如遇到问题：
1. 查看 [常见问题](#常见问题)
2. 查看详细的 [导入指南](../../migration/seeds/README_小说导入指南.md)
3. 检查日志输出
4. 提交 Issue

---

**祝测试顺利！** 🎉

