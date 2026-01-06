// MongoDB全文搜索索引创建脚本
// 用于Phase 2搜索功能
// 执行方式: mongo qingyu_db create_search_indexes.js

// 切换到数据库
db = db.getSiblingDB('qingyu_db');

print('=== 开始创建搜索索引 ===');

// ============ 书籍全文索引 ============

print('\n1. 创建书籍全文索引...');

// 检查是否已存在索引
var existingIndexes = db.books.getIndexes();
var textIndexExists = existingIndexes.some(function(index) {
    return index.name === 'book_text_search';
});

if (textIndexExists) {
    print('   - 书籍全文索引已存在，先删除旧索引');
    db.books.dropIndex('book_text_search');
}

// 创建全文索引
db.books.createIndex(
    {
        "title": "text",
        "author": "text",
        "description": "text",
        "tags": "text"
    },
    {
        name: "book_text_search",
        weights: {
            title: 10,       // 标题权重最高
            author: 5,       // 作者权重次之
            tags: 3,         // 标签权重中等
            description: 1   // 描述权重最低
        },
        default_language: "none",  // 支持中文
        language_override: "language"
    }
);

print('   ✓ 书籍全文索引创建成功');

// ============ 文档全文索引 ============

print('\n2. 创建文档全文索引...');

// 检查documents集合是否存在
var collections = db.getCollectionNames();
if (collections.indexOf('documents') === -1) {
    print('   - documents集合不存在，跳过');
} else {
    var docIndexExists = db.documents.getIndexes().some(function(index) {
        return index.name === 'document_text_search';
    });

    if (docIndexExists) {
        print('   - 文档全文索引已存在，先删除旧索引');
        db.documents.dropIndex('document_text_search');
    }

    db.documents.createIndex(
        {
            "title": "text",
            "content": "text",
            "tags": "text"
        },
        {
            name: "document_text_search",
            weights: {
                title: 10,
                tags: 5,
                content: 1
            },
            default_language: "none"
        }
    );

    print('   ✓ 文档全文索引创建成功');
}

// ============ 其他优化索引 ============

print('\n3. 创建其他优化索引...');

// 书籍分类索引
db.books.createIndex({ "category": 1 }, { name: "book_category_idx" });
print('   ✓ 书籍分类索引创建成功');

// 书籍作者索引
db.books.createIndex({ "author": 1 }, { name: "book_author_idx" });
print('   ✓ 书籍作者索引创建成功');

// 书籍标签索引
db.books.createIndex({ "tags": 1 }, { name: "book_tags_idx" });
print('   ✓ 书籍标签索引创建成功');

// 书籍创建时间索引（用于排序）
db.books.createIndex({ "created_at": -1 }, { name: "book_created_idx" });
print('   ✓ 书籍创建时间索引创建成功');

// 书籍热度索引（用于排序）
if (db.books.findOne({"popularity": {$exists: true}})) {
    db.books.createIndex({ "popularity": -1 }, { name: "book_popularity_idx" });
    print('   ✓ 书籍热度索引创建成功');
}

// 文档项目ID索引
if (collections.indexOf('documents') !== -1) {
    db.documents.createIndex({ "project_id": 1 }, { name: "document_project_idx" });
    print('   ✓ 文档项目ID索引创建成功');
}

// ============ 验证索引 ============

print('\n4. 验证创建的索引...');

print('\n   书籍索引:');
db.books.getIndexes().forEach(function(index) {
    print('   - ' + index.name);
});

if (collections.indexOf('documents') !== -1) {
    print('\n   文档索引:');
    db.documents.getIndexes().forEach(function(index) {
        print('   - ' + index.name);
    });
}

// ============ 测试搜索 ============

print('\n5. 测试全文搜索...');

// 测试书籍搜索
var bookCount = db.books.count();
if (bookCount > 0) {
    print('   测试书籍搜索...');
    var testResult = db.books.find(
        { $text: { $search: "小说" } },
        { score: { $meta: "textScore" } }
    ).sort({ score: { $meta: "textScore" } }).limit(1);

    if (testResult.hasNext()) {
        print('   ✓ 书籍搜索测试通过');
    } else {
        print('   - 书籍搜索测试: 未找到结果（可能是数据库为空）');
    }
} else {
    print('   - 数据库中没有书籍数据，跳过测试');
}

print('\n=== 搜索索引创建完成 ===\n');

// ============ TODO: Elasticsearch集成 ============
// TODO(Phase3): Elasticsearch集成
// 当MongoDB全文索引性能不满足需求时，可切换到Elasticsearch
//
// 1. 部署Elasticsearch集群
// 2. 创建索引映射（Mapping）:
//    - books索引: title, author, description, tags, category
//    - documents索引: title, content, tags, project_id
//
// 3. 数据同步:
//    - 批量导入现有数据
//    - 监听MongoDB变更事件，实时同步
//
// 4. 搜索服务切换:
//    - 实现ElasticsearchBackend
//    - 配置切换: search.backend: "elasticsearch"
//
// 5. 高级功能:
//    - 拼音搜索（IK分词器）
//    - 同义词扩展
//    - 搜索结果高亮
//    - 搜索建议（Completion Suggester）
//
// 优先级: P1
// 预计工时: 3天

