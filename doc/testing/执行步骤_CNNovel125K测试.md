# CNNovel125K 测试执行步骤

## 🎯 测试目标

在 test 分支上使用 CNNovel125K 数据集进行功能测试，验证：
- 书店浏览和搜索功能
- 阅读功能（章节、书签、历史）
- 推荐系统功能
- 整体性能表现

---

## ✅ 准备工作（已完成）

- ✅ test 分支已从 dev 分支创建
- ✅ Python 数据加载脚本已创建
- ✅ Go 数据导入服务已创建
- ✅ 自动化测试脚本已创建
- ✅ 完整文档已准备

---

## 📋 执行步骤

### 步骤 1: 确保环境准备就绪

#### 检查 MongoDB
```bash
# Windows: 检查 MongoDB 服务状态
services.msc
# 确保 MongoDB 服务正在运行

# Linux/Mac: 启动 MongoDB
sudo systemctl start mongod
sudo systemctl status mongod
```

#### 检查 Python
```bash
python --version
# 应该显示 Python 3.7+

# 安装 datasets 库
pip install datasets
```

#### 检查 Go
```bash
go version
# 应该显示 Go 1.21+
```

### 步骤 2: 执行一键测试脚本

#### Windows 用户
```bash
# 在项目根目录执行
scripts\test_novel_import.bat
```

#### Linux/Mac 用户
```bash
# 添加执行权限
chmod +x scripts/test_novel_import.sh

# 执行脚本
./scripts/test_novel_import.sh
```

脚本会自动完成：
1. 检查环境
2. 从 Hugging Face 下载数据（首次会较慢）
3. 处理并生成 JSON 文件
4. 验证数据格式
5. 导入到 MongoDB

**预计时间**: 5-15 分钟（首次下载数据集）

### 步骤 3: 验证数据导入

#### 检查 MongoDB
```bash
mongo
use qingyu

# 查看书籍数量
db.books.count()
# 应该显示约 100

# 查看章节数量
db.chapters.count()
# 应该显示约 3000+

# 查看示例数据
db.books.findOne()
```

#### 检查索引
```javascript
// 查看书籍集合索引
db.books.getIndexes()

// 查看章节集合索引
db.chapters.getIndexes()
```

### 步骤 4: 启动服务器

```bash
# 启动后端服务
go run cmd/server/main.go
```

等待服务器启动成功，应该看到：
```
服务器启动在端口 :8080
```

### 步骤 5: 测试书店功能

打开新的终端窗口，执行以下测试：

#### 5.1 获取书籍列表
```bash
curl "http://localhost:8080/api/v1/bookstore/books?page=1&pageSize=10"
```

**预期结果**: 返回 10 本书的列表，包含 title, author, category 等信息

#### 5.2 搜索书籍
```bash
# 从上一步获取一个书名，进行搜索测试
curl "http://localhost:8080/api/v1/bookstore/books/search?keyword=XXX"
```

**预期结果**: 返回匹配的书籍列表

#### 5.3 按分类筛选
```bash
curl "http://localhost:8080/api/v1/bookstore/books?category=玄幻"
```

**预期结果**: 只返回玄幻分类的书籍

#### 5.4 获取书籍详情
```bash
# 从列表中获取一个书籍 ID
curl "http://localhost:8080/api/v1/bookstore/books/{BOOK_ID}"
```

**预期结果**: 返回书籍完整信息

### 步骤 6: 测试阅读功能

#### 6.1 获取章节列表
```bash
curl "http://localhost:8080/api/v1/bookstore/books/{BOOK_ID}/chapters"
```

**预期结果**: 返回该书所有章节列表

#### 6.2 获取章节内容
```bash
# 从章节列表中获取一个章节 ID
curl "http://localhost:8080/api/v1/bookstore/chapters/{CHAPTER_ID}"
```

**预期结果**: 返回章节完整内容

#### 6.3 测试书签功能
```bash
# 需要先登录获取 token
# 创建书签
curl -X POST "http://localhost:8080/api/v1/reading/bookmarks" \
  -H "Authorization: Bearer {TOKEN}" \
  -H "Content-Type: application/json" \
  -d '{
    "book_id": "{BOOK_ID}",
    "chapter_id": "{CHAPTER_ID}",
    "position": 100
  }'
```

**预期结果**: 书签创建成功

#### 6.4 测试阅读历史
```bash
# 记录阅读历史
curl -X POST "http://localhost:8080/api/v1/reading/history" \
  -H "Authorization: Bearer {TOKEN}" \
  -H "Content-Type: application/json" \
  -d '{
    "book_id": "{BOOK_ID}",
    "chapter_id": "{CHAPTER_ID}",
    "progress": 50
  }'
```

**预期结果**: 阅读历史记录成功

### 步骤 7: 测试推荐功能

#### 7.1 获取推荐书籍
```bash
curl "http://localhost:8080/api/v1/recommendation/books"
```

**预期结果**: 返回推荐书籍列表

#### 7.2 获取热门书籍
```bash
curl "http://localhost:8080/api/v1/bookstore/books?is_hot=true"
```

**预期结果**: 返回标记为热门的书籍

#### 7.3 获取精选书籍
```bash
curl "http://localhost:8080/api/v1/bookstore/books?is_featured=true"
```

**预期结果**: 返回标记为精选的书籍

### 步骤 8: 性能测试

#### 8.1 响应时间测试
使用浏览器开发者工具或 Postman 记录各 API 的响应时间

目标指标：
- 书籍列表: < 100ms
- 搜索: < 200ms
- 详情: < 50ms
- 章节: < 100ms

#### 8.2 并发测试（可选）
```bash
# 安装 Apache Bench
# Windows: 下载并安装 Apache
# Linux: sudo apt-get install apache2-utils
# Mac: 系统自带

# 测试书籍列表
ab -n 1000 -c 10 http://localhost:8080/api/v1/bookstore/books

# 查看结果中的:
# - Requests per second (越高越好)
# - Time per request (越低越好)
# - Failed requests (应该为 0)
```

### 步骤 9: 记录测试结果

打开测试报告模板：
```
doc/testing/CNNovel125K测试报告.md
```

填写所有测试结果，包括：
- 数据导入统计
- 功能测试结果
- 性能测试数据
- 发现的问题
- 优化建议

### 步骤 10: 清理测试数据（可选）

测试完成后，如需清理数据：
```bash
go run cmd/migrate/main.go -command=clean-novels
```

系统会要求确认，输入 `yes` 即可删除所有测试数据。

---

## 📊 测试检查清单

使用以下清单确保测试完整：

### 环境准备
- [ ] MongoDB 正常运行
- [ ] Python 环境就绪
- [ ] Go 环境就绪
- [ ] 数据集已下载

### 数据导入
- [ ] Python 脚本执行成功
- [ ] JSON 文件已生成
- [ ] 数据验证通过
- [ ] MongoDB 导入成功
- [ ] 索引创建成功
- [ ] 数据统计正确

### 书店功能
- [ ] 书籍列表正常
- [ ] 分页功能正常
- [ ] 排序功能正常
- [ ] 搜索功能正常
- [ ] 分类筛选正常
- [ ] 书籍详情正常

### 阅读功能
- [ ] 章节列表正常
- [ ] 章节内容正常
- [ ] 书签功能正常
- [ ] 阅读历史正常
- [ ] 进度保存正常

### 推荐功能
- [ ] 推荐列表正常
- [ ] 热门书籍正常
- [ ] 精选书籍正常
- [ ] 推荐算法合理

### 性能测试
- [ ] 响应时间符合预期
- [ ] 并发测试通过
- [ ] 无明显性能瓶颈
- [ ] 索引使用正常

### 文档记录
- [ ] 测试报告已填写
- [ ] 问题已记录
- [ ] 截图已保存
- [ ] 优化建议已记录

---

## 🐛 常见问题处理

### 问题 1: MongoDB 连接失败
```
解决方案:
1. 检查 MongoDB 服务是否启动
2. 检查配置文件中的连接字符串
3. 检查防火墙设置
```

### 问题 2: Python 数据加载很慢
```
解决方案:
1. 首次下载数据集需要时间，耐心等待
2. 数据会缓存在本地，后续会快很多
3. 可以先用 --max-novels 10 测试
```

### 问题 3: 导入数据时内存不足
```
解决方案:
1. 减少导入数量（使用 --max-novels 50）
2. 增加 MongoDB 内存限制
3. 批量导入章节（代码已实现）
```

### 问题 4: API 响应慢
```
解决方案:
1. 检查索引是否创建成功
2. 使用 explain() 分析查询
3. 考虑添加缓存
```

---

## 📚 参考文档

- [快速开始指南](./CNNovel125K快速开始.md)
- [详细导入指南](../../migration/seeds/README_小说导入指南.md)
- [功能总结](./CNNovel125K功能总结.md)
- [测试报告模板](./CNNovel125K测试报告.md)
- [API 文档](../api/)

---

## 🎉 完成测试后

1. ✅ 填写测试报告
2. ✅ 提交测试结果
3. ✅ 记录发现的问题
4. ✅ 提出优化建议
5. ✅ 清理测试数据（可选）

---

**祝测试顺利！** 如有问题，请查看文档或联系开发团队。

