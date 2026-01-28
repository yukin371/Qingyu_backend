// MongoDB Profiler配置脚本
// 用于启用和配置MongoDB慢查询监控

print("=== 配置MongoDB Profiler ===");

// 设置profiling级别为1 (仅记录慢查询)，阈值为100ms
db.setProfilingLevel(1, { slowms: 100 });

// 限制profiler集合大小为100MB
db.system.profile.convertToCapped({
    size: 104857600  // 100MB
});

print("✅ Profiler配置完成");
print("   级别: 1 (仅记录慢查询)");
print("   阈值: 100ms");
print("   存储: 100MB (循环覆盖)");

// 验证配置
var status = db.getProfilingStatus();
print("\n当前配置:");
printjson(status);

// 显示如何查询慢查询
print("\n查询慢查询示例:");
print("db.system.profile.find().sort({ts: -1}).limit(10)");
print("\n解释说明:");
print("- was: 0 = 关闭, 1 = 仅慢查询, 2 = 全部查询");
print("- slowms: 慢查询阈值(毫秒)");
print("- system.profile是固定大小集合，旧数据会被自动覆盖");
