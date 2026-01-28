// MongoDB慢查询分析工具
// 用于分析system.profile集合中的慢查询数据

print("=== 慢查询分析报告 ===");
print("生成时间: " + new Date().toISOString());
print("");

// 获取慢查询阈值
var profilingStatus = db.getProfilingStatus();
var slowMs = profilingStatus.slowms || 100;
print("当前慢查询阈值: " + slowMs + "ms");
print("");

// 1. 统计总慢查询数量
var totalSlowQueries = db.system.profile.countDocuments({
    millis: { $gt: slowMs }
});
print("总慢查询数: " + totalSlowQueries);
print("");

// 2. 按集合分组统计
print("按集合统计:");
var collectionStats = db.system.profile.aggregate([
    {
        $match: {
            millis: { $gt: slowMs }
        }
    },
    {
        $group: {
            _id: "$ns",
            count: { $sum: 1 },
            avgTime: { $avg: "$millis" },
            maxTime: { $max: "$millis" },
            minTime: { $min: "$millis" },
            totalTime: { $sum: "$millis" }
        }
    },
    {
        $sort: { count: -1 }
    }
]).toArray();

if (collectionStats.length > 0) {
    collectionStats.forEach(function(stat) {
        print("  " + stat._id + ":");
        print("    次数: " + stat.count);
        print("    平均耗时: " + stat.avgTime.toFixed(2) + "ms");
        print("    最大耗时: " + stat.maxTime + "ms");
        print("    最小耗时: " + stat.minTime + "ms");
        print("    总耗时: " + stat.totalTime.toFixed(2) + "ms");
        print("");
    });
} else {
    print("  (无慢查询数据)");
    print("");
}

// 3. Top 10最慢查询
print("Top 10最慢查询:");
var topSlowQueries = db.system.profile.aggregate([
    {
        $match: {
            millis: { $gt: slowMs }
        }
    },
    {
        $sort: { millis: -1 }
    },
    {
        $limit: 10
    },
    {
        $project: {
            ns: 1,
            millis: 1,
            query: "$query",
            ts: 1
        }
    }
]).toArray();

if (topSlowQueries.length > 0) {
    topSlowQueries.forEach(function(query, index) {
        print((index + 1) + ". " + query.ns + " - " + query.millis + "ms");
        print("   时间: " + query.ts);
        print("   查询: " + JSON.stringify(query.query));
        print("");
    });
} else {
    print("  (无慢查询数据)");
    print("");
}

// 4. 统计未使用索引的慢查询
// 通过检查execStats中是否包含COLLSCAN（全表扫描）来判断
print("未使用索引的慢查询分析:");

var unindexedQueries = db.system.profile.aggregate([
    {
        $match: {
            millis: { $gt: slowMs },
            "execStats.stage": "COLLSCAN"
        }
    },
    {
        $group: {
            _id: "$ns",
            count: { $sum: 1 },
            avgTime: { $avg: "$millis" },
            maxTime: { $max: "$millis" },
            queries: {
                $push: {
                    query: "$query",
                    millis: "$millis",
                    ts: "$ts"
                }
            }
        }
    },
    {
        $sort: { count: -1 }
    }
]).toArray();

if (unindexedQueries.length > 0) {
    var totalCount = 0;
    unindexedQueries.forEach(function(stat) {
        totalCount += stat.count;
        print("  集合: " + stat._id);
        print("    未使用索引次数: " + stat.count);
        print("    平均耗时: " + stat.avgTime.toFixed(2) + "ms");
        print("    最大耗时: " + stat.maxTime + "ms");
        print("    示例查询:");
        stat.queries.slice(0, 3).forEach(function(q) {
            print("      - " + JSON.stringify(q.query) + " (" + q.millis + "ms)");
        });
        print("");
    });
    print("总计未使用索引的慢查询: " + totalCount + "次");
} else {
    print("  (所有慢查询都使用了索引或无慢查询数据)");
    print("");
}

// 5. 按操作类型统计
print("按操作类型统计:");
var opStats = db.system.profile.aggregate([
    {
        $match: {
            millis: { $gt: slowMs }
        }
    },
    {
        $group: {
            _id: "$op",
            count: { $sum: 1 },
            avgTime: { $avg: "$millis" },
            maxTime: { $max: "$millis" }
        }
    },
    {
        $sort: { count: -1 }
    }
]).toArray();

if (opStats.length > 0) {
    opStats.forEach(function(stat) {
        print("  " + stat._id + ": " + stat.count + "次, 平均" + stat.avgTime.toFixed(2) + "ms, 最大" + stat.maxTime + "ms");
    });
    print("");
}

// 6. 性能建议
print("=== 优化建议 ===");
if (totalSlowQueries === 0) {
    print("✅ 当前没有慢查询，数据库性能良好！");
} else {
    print("⚠️ 发现 " + totalSlowQueries + " 个慢查询，建议:");

    if (unindexedQueries.length > 0) {
        print("1. 为频繁全表扫描的查询添加索引");
        print("2. 检查查询条件是否使用了合适的字段");
    }

    if (collectionStats.length > 0) {
        var worstCollection = collectionStats[0];
        if (worstCollection.avgTime > 500) {
            print("3. 集合 " + worstCollection._id + " 平均耗时超过500ms，需要重点关注");
        }
    }

    print("4. 使用auto_analyze_slow_queries.js获取详细的索引建议");
    print("5. 考虑调整慢查询阈值（当前" + slowMs + "ms）");
}

print("");
print("=== 分析完成 ===");
print("提示: 使用 auto_analyze_slow_queries.js 获取详细的索引优化建议");
