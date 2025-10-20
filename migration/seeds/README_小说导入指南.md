# CNNovel125K 小说导入指南

## 概述

本指南介绍如何从 Hugging Face 的 CNNovel125K 数据集导入小说数据到青羽后端系统，用于功能测试。

## 前置要求

### Python 环境
- Python 3.7+
- pip 包管理器

### 安装依赖
```bash
pip install datasets
```

### Go 环境
- Go 1.21+
- MongoDB 连接配置正确

## 导入流程

### 步骤 1: 使用 Python 脚本加载数据

导入少量数据（100本）用于快速测试：
```bash
python scripts/import_novels.py --max-novels 100 --output data/novels_100.json
```

导入中等数量数据（500本）：
```bash
python scripts/import_novels.py --max-novels 500 --output data/novels_500.json
```

导入较多数据（1000本）：
```bash
python scripts/import_novels.py --max-novels 1000 --output data/novels_1000.json
```

#### 参数说明
- `--max-novels`: 最大导入小说数量（默认1000）
- `--chapter-size`: 每章字数（默认3000）
- `--output`: 输出 JSON 文件路径

### 步骤 2: 验证数据（试运行）

在正式导入前，建议先进行数据验证：
```bash
go run cmd/migrate/main.go -command=import-novels -file=data/novels_100.json -dry-run=true
```

这会验证数据格式但不实际写入数据库。

### 步骤 3: 正式导入数据

确认数据无误后，执行正式导入：
```bash
go run cmd/migrate/main.go -command=import-novels -file=data/novels_100.json
```

导入过程会：
1. 读取 JSON 文件
2. 验证数据完整性
3. 批量插入书籍记录
4. 批量插入章节记录（每批500章）
5. 创建数据库索引
6. 显示统计信息

### 步骤 4: 验证导入结果

导入完成后，可以：

1. **查看统计信息**（在导入输出中会显示）
   - 书籍总数
   - 章节总数

2. **通过 API 验证**
   - 访问书籍列表 API
   - 测试搜索功能
   - 验证章节阅读

3. **直接查询 MongoDB**
```bash
# 连接到 MongoDB
mongo

# 切换到数据库
use qingyu

# 查看书籍数量
db.books.count()

# 查看章节数量
db.chapters.count()

# 查看示例书籍
db.books.findOne()
```

## 数据映射关系

### CNNovel125K → 青羽模型

| CNNovel125K 字段 | 青羽字段 | 说明 |
|-----------------|---------|------|
| title | Book.Title | 书名 |
| author | Book.Author | 作者名 |
| content | 拆分为多个 Chapter | 按章节拆分 |
| category | Book.Categories | 分类（映射处理） |
| word_count | Book.WordCount | 字数统计 |
| rating | 推荐标记 | >=4.5 推荐，>=4.8 精选 |

### 自动生成字段

- **BookID**: MongoDB ObjectID 自动生成
- **Cover**: 默认占位图片
- **Status**: 标记为"已完结"（completed）
- **PublishedAt**: 假设30天前发布
- **IsFree**: 默认免费
- **IsRecommended**: 评分 >= 4.5
- **IsFeatured**: 评分 >= 4.8
- **IsHot**: 字数 > 50万

## 清理测试数据

### 方法 1: 使用清理脚本

```bash
go run cmd/migrate/main.go -command=clean-novels
```

### 方法 2: 手动清理

```bash
mongo
use qingyu
db.books.deleteMany({})
db.chapters.deleteMany({})
```

⚠️ **警告**: 这会删除所有书籍和章节数据！

## 性能建议

### 小规模测试（100本）
- 适用于：功能验证、界面测试
- 导入时间：约 10-30 秒
- 数据量：~100 书籍，~3000 章节

### 中等规模测试（500本）
- 适用于：综合功能测试、性能初步评估
- 导入时间：约 1-3 分钟
- 数据量：~500 书籍，~15000 章节

### 大规模测试（1000本）
- 适用于：性能测试、压力测试
- 导入时间：约 3-10 分钟
- 数据量：~1000 书籍，~30000 章节

## 常见问题

### Q1: Python 脚本加载数据集很慢？
**A**: 首次下载数据集需要时间，数据会缓存在本地。后续加载会快很多。

### Q2: 导入过程中断了怎么办？
**A**: 需要清理已导入的数据，然后重新导入。建议先用少量数据测试。

### Q3: 如何只导入特定分类的小说？
**A**: 可以修改 Python 脚本，在 `load_and_process()` 方法中添加分类过滤。

### Q4: 导入的书籍没有封面？
**A**: 当前使用占位图片。可以后续集成封面生成服务或爬取真实封面。

### Q5: 如何处理敏感词？
**A**: Python 脚本有基础的敏感词过滤，但词库需要完善。建议使用专业的敏感词库。

## 测试建议

### 基础功能测试
1. 书店浏览（分页、排序）
2. 书籍搜索（标题、作者）
3. 分类筛选
4. 书籍详情展示

### 阅读功能测试
1. 章节列表展示
2. 章节内容阅读
3. 阅读进度保存
4. 书签功能

### 推荐系统测试
1. 基于分类推荐
2. 基于评分推荐
3. 热门书籍推荐

### 性能测试
1. 列表查询性能
2. 搜索响应时间
3. 章节加载速度
4. 并发用户测试

## 许可证说明

CNNovel125K 数据集使用 Apache-2.0 许可证，可商用。

⚠️ **注意**: 
- 数据为原始爬取，未脱敏
- 需要自行进行敏感词和版权二次过滤
- 仅用于测试，不建议直接用于生产环境

## 技术支持

如有问题，请查看：
- [数据集文档](https://huggingface.co/datasets/RyokoAI/CNNovel125K)
- [项目测试文档](../../doc/testing/)
- [架构设计文档](../../doc/architecture/)

