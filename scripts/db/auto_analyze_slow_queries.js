// MongoDB慢查询自动分析和优化建议工具
// 用于自动分析查询模式并提供索引优化建议

print("=== 慢查询自动分析和优化建议 ===");
print("生成时间: " + new Date().toISOString());
print("");

// 获取慢查询阈值
var profilingStatus = db.getProfilingStatus();
var slowMs = profilingStatus.slowms || 100;

// 从system.profile聚合慢查询
var slowQueries = db.system.profile.aggregate([
    {
        $match: {
            millis: { $gt: slowMs }
        }
    },
    {
        $project: {
            ns: 1,
            op: 1,
            millis: 1,
            query: "$query",
            sort: "$orderby",
            ts: 1,
            execStats: "$execStats",
            hasIndexScan: {
                $cond: [
                    {
                        $or: [
                            { $eq: ["$execStats.stage", "IXSCAN"] },
                            { $ifNull: ["$execStats.stage", false] }
                        ]
                    },
                    true,
                    {
                        $anyElementTrue: {
                            $map: {
                                input: { $ifNull: ["$execStats.inputStages", []] },
                                as: "stage",
                                in: { $eq: ["$$stage.stage", "IXSCAN"] }
                            }
                        }
                    }
                ]
            }
        }
    }
]).toArray();

if (slowQueries.length === 0) {
    print("✅ 当前没有慢查询数据，无需分析！");
    print("");
    print("提示:");
    print("- 确保Profiler已启用: db.setProfilingLevel(1, {slowms: 100})");
    print("- 执行一些查询后再运行此脚本");
    quit(0);
}

// 按查询模式聚合
var queryPatterns = {};
var collectionPatterns = {};

slowQueries.forEach(function(q) {
    // 生成查询模式的唯一标识（标准化查询对象）
    var queryKey = JSON.stringify({
        ns: q.ns,
        query: normalizeQuery(q.query || {}),
        sort: normalizeQuery(q.sort || {})
    });

    if (!queryPatterns[queryKey]) {
        queryPatterns[queryKey] = {
            ns: q.ns,
            query: q.query || {},
            sort: q.sort || {},
            count: 0,
            totalTime: 0,
            maxTime: 0,
            minTime: Number.MAX_VALUE,
            hasIndexScan: true,
            samples: []
        };
    }

    var pattern = queryPatterns[queryKey];
    pattern.count++;
    pattern.totalTime += q.millis;
    pattern.maxTime = Math.max(pattern.maxTime, q.millis);
    pattern.minTime = Math.min(pattern.minTime, q.millis);

    // 检查是否使用了索引
    if (!q.hasIndexScan) {
        pattern.hasIndexScan = false;
    }

    // 保存样本（最多5个）
    if (pattern.samples.length < 5) {
        pattern.samples.push({
            millis: q.millis,
            ts: q.ts,
            execStats: q.execStats
        });
    }

    // 按集合统计
    if (!collectionPatterns[q.ns]) {
        collectionPatterns[q.ns] = {
            count: 0,
            totalTime: 0
        };
    }
    collectionPatterns[q.ns].count++;
    collectionPatterns[q.ns].totalTime += q.millis;
});

// 标准化查询对象（移除特定值以识别模式）
function normalizeQuery(obj) {
    var normalized = {};
    for (var key in obj) {
        if (obj.hasOwnProperty(key)) {
            var value = obj[key];
            if (typeof value === 'object' && value !== null && !Array.isArray(value)) {
                // 对于查询操作符，保留结构但标准化值
                normalized[key] = normalizeQuery(value);
            } else if (typeof value !== 'function') {
                // 用占位符替换实际值
                normalized[key] = "<value>";
            }
        }
    }
    return normalized;
}

// 检查查询是否使用了索引
function checkIndexUsage(execStats) {
    if (!execStats) return false;

    // 检查顶层stage
    if (execStats.stage === "IXSCAN") {
        return true;
    }

    // 检查子stage
    if (execStats.inputStages && execStats.inputStages.length > 0) {
        for (var i = 0; i < execStats.inputStages.length; i++) {
            if (checkIndexUsage(execStats.inputStages[i])) {
                return true;
            }
        }
    }

    return false;
}

// 生成索引建议
function generateIndexSuggestion(query, sort) {
    var indexFields = [];

    // 提取查询条件中的字段
    for (var key in query) {
        if (query.hasOwnProperty(key) && key !== '_id') {
            // 跳过操作符，提取实际字段
            if (key.startsWith('$')) {
                // 处理 {$or: [...]} 等操作符
                continue;
            }
            indexFields.push(key);
        }
    }

    // 提取排序字段
    if (sort) {
        for (var sortKey in sort) {
            if (sort.hasOwnProperty(sortKey)) {
                if (indexFields.indexOf(sortKey) === -1) {
                    indexFields.push(sortKey);
                }
            }
        }
    }

    if (indexFields.length === 0) {
        return { _id: 1 }; // 默认建议
    }

    // 构建索引对象
    var indexSpec = {};
    indexFields.forEach(function(field, index) {
        indexSpec[field] = 1; // 默认升序
    });

    return indexSpec;
}

// 计算优先级
function calculatePriority(pattern) {
    var avgTime = pattern.totalTime / pattern.count;

    // P0: 高频慢查询 (次数>10 && 平均耗时>200ms)
    if (pattern.count > 10 && avgTime > 200) {
        return "P0";
    }
    // P1: 极慢查询 (平均耗时>500ms)
    if (avgTime > 500) {
        return "P1";
    }
    // P2: 中等慢查询 (平均耗时>200ms 或 次数>5)
    if (avgTime > 200 || pattern.count > 5) {
        return "P2";
    }
    // P3: 一般慢查询
    return "P3";
}

// 获取优先级描述和颜色标记
function getPriorityInfo(priority) {
    switch(priority) {
        case "P0":
            return { icon: "🔴", desc: "高频慢查询 - 需要立即优化" };
        case "P1":
            return { icon: "🟠", desc: "极慢查询 - 建议尽快优化" };
        case "P2":
            return { icon: "🟡", desc: "中等慢查询 - 建议优化" };
        case "P3":
            return { icon: "🟢", desc: "一般慢查询 - 可以优化" };
        default:
            return { icon: "⚪", desc: "低优先级" };
    }
}

// 输出分析结果
print("分析查询总数: " + slowQueries.length);
print("发现查询模式: " + Object.keys(queryPatterns).length);
print("");

// 按优先级排序
var sortedPatterns = Object.keys(queryPatterns).map(function(key) {
    var pattern = queryPatterns[key];
    pattern.avgTime = pattern.totalTime / pattern.count;
    pattern.priority = calculatePriority(pattern);
    return pattern;
});

// 排序: P0 > P1 > P2 > P3
var priorityOrder = { "P0": 0, "P1": 1, "P2": 2, "P3": 3 };
sortedPatterns.sort(function(a, b) {
    var priorityDiff = priorityOrder[a.priority] - priorityOrder[b.priority];
    if (priorityDiff !== 0) return priorityDiff;
    // 相同优先级按平均耗时排序
    return b.avgTime - a.avgTime;
});

// 输出每个查询模式的建议
var p0Count = 0;
var p1Count = 0;

sortedPatterns.forEach(function(pattern, index) {
    var priorityInfo = getPriorityInfo(pattern.priority);
    var indexSuggestion = generateIndexSuggestion(pattern.query, pattern.sort);

    print("[" + (index + 1) + "] 查询模式 #" + (index + 1));
    print("集合: " + pattern.ns);
    print("查询: " + JSON.stringify(pattern.query));
    if (Object.keys(pattern.sort).length > 0) {
        print("排序: " + JSON.stringify(pattern.sort));
    }
    print("统计信息:");
    print("  次数: " + pattern.count);
    print("  平均耗时: " + pattern.avgTime.toFixed(2) + "ms");
    print("  最大耗时: " + pattern.maxTime + "ms");
    print("  最小耗时: " + pattern.minTime + "ms");

    // 索引使用情况
    if (pattern.hasIndexScan) {
        print("✅ 索引使用: 已使用索引");
    } else {
        print("⚠️ 索引使用: 未使用索引（检测到全表扫描）");
        print("💡 建议: 为该查询添加索引");
        print("   推荐索引: " + JSON.stringify(indexSuggestion));
        print("   创建命令: db." + pattern.ns.split('.').pop() + ".createIndex(" + JSON.stringify(indexSuggestion) + ")");
    }

    // 优先级标记
    print(priorityInfo.icon + " 优先级: " + pattern.priority + " - " + priorityInfo.desc);

    // 如果索引使用良好但仍然慢，提供其他建议
    if (pattern.hasIndexScan && pattern.avgTime > 200) {
        print("💡 其他优化建议:");
        if (pattern.count > 10) {
            print("   - 考虑添加缓存减少查询频率");
        }
        print("   - 检查返回的文档数量是否过多");
        print("   - 考虑使用投影只返回需要的字段");
    }

    print("");

    // 统计P0和P1问题
    if (pattern.priority === "P0") p0Count++;
    if (pattern.priority === "P1") p1Count++;
});

// 输出集合级别的统计
print("=== 集合级别的慢查询统计 ===");
var sortedCollections = Object.keys(collectionPatterns).map(function(ns) {
    var stats = collectionPatterns[ns];
    stats.ns = ns;
    stats.avgTime = stats.totalTime / stats.count;
    return stats;
});
sortedCollections.sort(function(a, b) { return b.count - a.count; });

sortedCollections.forEach(function(coll) {
    print(coll.ns + ":");
    print("  慢查询次数: " + coll.count);
    print("  总耗时: " + coll.totalTime.toFixed(2) + "ms");
    print("  平均耗时: " + coll.avgTime.toFixed(2) + "ms");
    print("");
});

// 输出总结
print("=== 分析总结 ===");
print("🔴 P0 (高频慢查询): " + p0Count + " 个");
print("🟠 P1 (极慢查询): " + p1Count + " 个");
print("🟡 P2 (中等慢查询): " + (sortedPatterns.filter(function(p) { return p.priority === 'P2'; }).length) + " 个");
print("🟢 P3 (一般慢查询): " + (sortedPatterns.filter(function(p) { return p.priority === 'P3'; }).length) + " 个");
print("");

if (p0Count > 0 || p1Count > 0) {
    print("⚠️ 发现 " + (p0Count + p1Count) + " 个需要优先处理的慢查询问题！");
    print("");
    print("建议操作:");
    print("1. 优先处理P0级别的查询，这些是高频且慢的查询");
    print("2. 检查P1级别的查询，这些是特别慢的查询");
    print("3. 使用 explain() 分析查询执行计划");
    print("4. 根据索引建议创建合适的索引");
    print("");
    print("示例创建索引:");
    print("  // 进入数据库");
    print("  use qingyu_dev");
    print("  // 创建索引（根据具体建议替换）");
    print("  db.collection_name.createIndex({field_name: 1})");
} else {
    print("✅ 没有发现严重的慢查询问题！");
    print("建议继续监控数据库性能。");
}

print("");
print("=== 分析完成 ===");
