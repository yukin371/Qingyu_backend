// MongoDB Profiler配置脚本
// 用途: 配置数据库慢查询监控
// 使用方法: mongosh qingyu --file scripts/enable_profiling.js

print("=== 配置MongoDB Profiler ===");

// 获取当前数据库名称
const dbName = db.getName();
print(`当前数据库: ${dbName}`);

// 配置慢查询阈值和profiling级别
// level: 0=off, 1=slow only, 2=all
// slowms: 慢查询阈值（毫秒）
const profilingLevel = 1;
const slowMs = 100;

print(`设置profiling级别为 ${profilingLevel} (仅记录慢查询)`);
print(`设置慢查询阈值为 ${slowMs}ms`);

try {
    db.setProfilingLevel(profilingLevel, { slowms: slowMs });
    print("✓ Profiling级别配置成功");
} catch (e) {
    print(`✗ 配置profiling级别失败: ${e}`);
    quit(1);
}

// 限制profiler collection大小（防止磁盘溢出）
const profilerSizeMB = 100;
const profilerSizeBytes = profilerSizeMB * 1024 * 1024;

print(`设置profiler存储上限为 ${profilerSizeMB}MB (循环覆盖)`);

try {
    // 检查system.profile集合是否存在
    const collections = db.getCollectionNames();
    if (collections.includes("system.profile")) {
        // 如果已存在，转换为capped collection
        db.system.profile.convertToCapped({
            size: profilerSizeBytes
        });
        print("✓ 已存在system.profile集合，转换为capped collection");
    } else {
        // 如果不存在，创建新的capped collection
        db.createCollection("system.profile", {
            capped: true,
            size: profilerSizeBytes
        });
        print("✓ 创建新的system.profile集合");
    }
} catch (e) {
    // 如果集合已经是capped的，可能会报错，这是正常的
    if (e.message.includes("already capped")) {
        print("✓ system.profile集合已经是capped collection");
    } else {
        print(`⚠ 设置profiler存储大小时出现警告: ${e}`);
    }
}

print("\n=== Profiler配置完成 ===");
print(`   级别: ${profilingLevel} (仅记录慢查询)`);
print(`   阈值: ${slowMs}ms`);
print(`   存储: ${profilerSizeMB}MB (循环覆盖)`);
print("");

// 验证配置
print("=== 验证当前配置 ===");
const status = db.getProfilingStatus();
print("当前配置状态:");
printjson(status);

// 显示profiler集合统计信息
try {
    const stats = db.system.profile.stats();
    print("\nsystem.profile集合统计:");
    print(`   文档数量: ${stats.count}`);
    print(`   存储大小: ${(stats.size / 1024 / 1024).toFixed(2)}MB`);
    print(`   总大小: ${(stats.storageSize / 1024 / 1024).toFixed(2)}MB`);
} catch (e) {
    print(`⚠ 无法获取集合统计: ${e}`);
}

// 显示示例查询：查找最慢的查询
print("\n=== 示例查询 ===");
print("查看最慢的5个查询:");
print('db.system.profile.find().sort({millis: -1}).limit(5).pretty()');
print("\n查看特定collection的慢查询:");
print('db.system.profile.find({ns: "qingyu.users"}).sort({millis: -1}).limit(5)');
print("\n查看慢于200ms的查询:");
print('db.system.profile.find({millis: {$gt: 200}}).sort({millis: -1}).pretty()');

print("\n✅ 配置完成！Profiler现在将记录超过" + slowMs + "ms的查询。");
