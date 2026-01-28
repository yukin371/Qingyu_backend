// 慢查询测试数据生成脚本
// 用于生成测试慢查询，验证分析工具的功能

print("=== 生成测试慢查询数据 ===");
print("生成时间: " + new Date().toISOString());
print("");

// 确保profiler已启用
var profilingStatus = db.getProfilingStatus();
print("当前Profiler状态:");
printjson(profilingStatus);
print("");

if (profilingStatus.was === 0) {
    print("⚠️ Profiler未启用，正在启用...");
    db.setProfilingLevel(1, { slowms: 100 });
    print("✅ Profiler已启用（级别1，阈值100ms）");
    print("");
}

// 创建测试集合
var testCollectionName = "test_slow_query";

// 检查集合是否存在，如果存在则删除
if (db.getCollectionNames().indexOf(testCollectionName) >= 0) {
    print("删除已存在的测试集合: " + testCollectionName);
    db[testCollectionName].drop();
}

// 创建测试集合并插入数据
print("创建测试集合并插入数据...");
var testData = [];
for (var i = 0; i < 1000; i++) {
    testData.push({
        user_id: i % 100,
        book_id: i % 50,
        chapter_id: i,
        title: "Test Book " + (i % 50),
        content: "Test content " + i,
        created_at: new Date(),
        status: i % 3
    });
}
db[testCollectionName].insertMany(testData);
print("✅ 插入 " + testData.length + " 条测试数据");
print("");

// 执行各种类型的慢查询
print("执行测试查询...");

// 1. 全表扫描查询（慢查询）
print("1. 执行全表扫描查询...");
for (var j = 0; j < 5; j++) {
    db[testCollectionName].find({
        $or: [
            { title: "Test Book " + j },
            { content: "Test content " + (j * 10) }
        ]
    }).toArray();
}

// 2. 未使用索引的查询
print("2. 执行未使用索引的查询...");
for (var k = 0; k < 10; k++) {
    db[testCollectionName].find({
        status: k % 3,
        title: "Test Book " + (k % 50)
    }).toArray();
}

// 3. 排序查询（可能较慢）
print("3. 执行排序查询...");
db[testCollectionName].find({}).sort({ created_at: -1 }).toArray();

// 4. 复杂查询
print("4. 执行复杂查询...");
db[testCollectionName].find({
    $and: [
        { user_id: { $gte: 10, $lte: 50 } },
        { status: { $ne: 2 } }
    ]
}).toArray();

print("✅ 测试查询执行完成");
print("");

// 验证慢查询是否被记录
var slowQueryCount = db.system.profile.countDocuments({
    ns: db.getName() + "." + testCollectionName
});
print("记录的慢查询数量: " + slowQueryCount);
print("");

if (slowQueryCount > 0) {
    print("✅ 测试数据生成成功！");
    print("现在可以运行以下命令分析慢查询:");
    print("  mongosh " + db.getName() + " < scripts/db/analyze_slow_queries.js");
    print("  mongosh " + db.getName() + " < scripts/db/auto_analyze_slow_queries.js");
    print("");
    print("或者直接在mongosh中:");
    print("  load('scripts/db/analyze_slow_queries.js')");
    print("  load('scripts/db/auto_analyze_slow_queries.js')");
} else {
    print("⚠️ 没有检测到慢查询，可能原因:");
    print("  1. 查询执行速度太快（低于阈值）");
    print("  2. Profiler配置问题");
    print("");
    print("建议检查:");
    print("  db.system.profile.find().sort({ts: -1}).limit(5)");
}

print("");
print("=== 测试数据生成完成 ===");
print("");
print("清理说明:");
print("测试完成后，可以使用以下命令清理测试数据:");
print("  db." + testCollectionName + ".drop()");
