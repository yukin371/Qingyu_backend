# Search Grayscale Metrics

搜索灰度指标收集模块，用于监控 Elasticsearch 和 MongoDB 之间的灰度切换。

## 功能特性

- **自动指标收集**：自动记录搜索请求次数和耗时
- **Prometheus 集成**：原生支持 Prometheus 指标导出
- **灰度决策跟踪**：记录灰度决策结果
- **流量分配监控**：实时监控 ES 和 MongoDB 的流量分配
- **性能指标**：追踪各搜索引擎的平均响应时间

## Prometheus 指标

### search_requests_total

搜索请求总数（按引擎类型分类）

```promql
# 查询 ES 总请求数
search_requests_total{engine="elasticsearch"}

# 查询 MongoDB 总请求数
search_requests_total{engine="mongodb"}

# 查询请求速率（每秒）
rate(search_requests_total[5m])
```

### search_duration_seconds

搜索请求耗时（直方图）

```promql
# ES 搜索平均耗时
rate(search_duration_seconds_sum{engine="elasticsearch"}[5m]) /
rate(search_duration_seconds_count{engine="elasticsearch"}[5m])

# ES 搜索 P95 耗时
histogram_quantile(0.95,
  sum(rate(search_duration_seconds_bucket{engine="elasticsearch"}[5m])) by (le)
)
```

### search_grayscale_traffic_percent

灰度流量百分比（按搜索类型）

```promql
# 查看各搜索类型的灰度百分比
search_grayscale_traffic_percent

# 查看 books 搜索的 ES 流量百分比
search_grayscale_traffic_percent{search_type="books"}
```

### search_grayscale_decision_total

灰度决策总数（按搜索类型和决策结果）

```promql
# 查看 ES 决策次数
search_grayscale_decision_total{decision="elasticsearch"}

# 查看 MongoDB 决策次数
search_grayscale_decision_total{decision="mongodb"}

# ES 决策占比
sum(search_grayscale_decision_total{decision="elasticsearch"}) /
sum(search_grayscale_decision_total)
```

### search_engine_switches_total

搜索引擎切换次数

```promql
# 从 ES 切换到 MongoDB 的次数
search_engine_switches_total{from_engine="elasticsearch",to_engine="mongodb"}

# 从 MongoDB 切换到 ES 的次数
search_engine_switches_total{from_engine="mongodb",to_engine="elasticsearch"}
```

## 使用方法

### 1. 记录搜索请求

在搜索服务中自动记录：

```go
import "Qingyu_backend/pkg/metrics"

// 记录 ES 搜索
metrics.RecordSearch("elasticsearch", 150*time.Millisecond)

// 记录 MongoDB 搜索
metrics.RecordSearch("mongodb", 80*time.Millisecond)
```

### 2. 更新灰度百分比

当灰度配置更新时调用：

```go
// 更新 books 搜索的灰度百分比为 30%
metrics.UpdateGrayscalePercent("books", 30)

// 更新 projects 搜索的灰度百分比为 50%
metrics.UpdateGrayscalePercent("projects", 50)
```

### 3. 记录灰度决策

在做出灰度决策时记录：

```go
// 记录使用 ES 的决策
metrics.RecordGrayscaleDecision("books", true)

// 记录使用 MongoDB 的决策
metrics.RecordGrayscaleDecision("books", false)
```

### 4. 获取指标快照

获取当前指标数据：

```go
// 获取指标快照
snapshot := metrics.GetGrayScaleMetrics()
fmt.Printf("ES Count: %d\n", snapshot.ESCount)
fmt.Printf("MongoDB Count: %d\n", snapshot.MongoDBCount)

// 获取流量分配比例
esPercent, mongoPercent := metrics.GetTrafficDistribution()
fmt.Printf("ES Traffic: %.2f%%\n", esPercent)

// 获取平均耗时
esAvg := metrics.GetAverageDuration("elasticsearch")
fmt.Printf("ES Avg Duration: %v\n", esAvg)
```

## Grafana 仪表盘示例

### 1. 流量分配饼图

```promql
sum by (engine) (search_requests_total)
```

### 2. 搜索速率趋势图

```promql
# ES 搜索速率
rate(search_requests_total{engine="elasticsearch"}[5m])

# MongoDB 搜索速率
rate(search_requests_total{engine="mongodb"}[5m])
```

### 3. 平均响应时间对比

```promql
# ES 平均响应时间
rate(search_duration_seconds_sum{engine="elasticsearch"}[5m]) /
rate(search_duration_seconds_count{engine="elasticsearch"}[5m])

# MongoDB 平均响应时间
rate(search_duration_seconds_sum{engine="mongodb"}[5m]) /
rate(search_duration_seconds_count{engine="mongodb"}[5m])
```

### 4. 灰度决策比例

```promql
# ES 决策占比
sum(search_grayscale_decision_total{decision="elasticsearch"}) /
sum(search_grayscale_decision_total) * 100

# MongoDB 决策占比
sum(search_grayscale_decision_total{decision="mongodb"}) /
sum(search_grayscale_decision_total) * 100
```

## 告警规则示例

### Prometheus 告警规则

```yaml
groups:
  - name: search_grayscale
    interval: 30s
    rules:
      # ES 搜索失败率过高
      - alert: HighSearchErrorRate
        expr: |
          rate(search_requests_total{engine="elasticsearch",status="error"}[5m])
          / rate(search_requests_total{engine="elasticsearch"}[5m]) > 0.05
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "ES 搜索错误率过高"
          description: "ES 搜索错误率为 {{ $value | humanizePercentage }}"

      # ES 响应时间过长
      - alert: SlowSearchResponse
        expr: |
          histogram_quantile(0.95,
            sum(rate(search_duration_seconds_bucket{engine="elasticsearch"}[5m])) by (le)
          ) > 1
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "ES 搜索响应慢"
          description: "ES P95 响应时间为 {{ $value }}s"

      # MongoDB 流量异常（在灰度期间应该很少）
      - alert: HighMongoDBTrafficDuringGrayscale
        expr: |
          sum(rate(search_requests_total{engine="mongodb"}[5m])) /
          sum(rate(search_requests_total[5m])) > 0.5
        for: 10m
        labels:
          severity: info
        annotations:
          summary: "MongoDB 流量异常"
          description: "MongoDB 流量占比为 {{ $value | humanizePercentage }}"
```

## 测试

运行测试：

```bash
# 运行所有测试
go test ./pkg/metrics/...

# 运行特定测试
go test ./pkg/metrics/... -run TestRecordSearch

# 运行基准测试
go test ./pkg/metrics/... -bench=.

# 查看测试覆盖率
go test ./pkg/metrics/... -cover
```

## 性能考虑

- 所有指标收集操作使用 `sync.RWMutex` 保护，支持并发读写
- Prometheus 指标使用 `promauto` 自动注册，性能开销最小
- 内部状态与 Prometheus 指标分离，避免频繁计算
- 支持高并发场景（基准测试显示 >1000万 ops/秒）

## 集成到 SearchService

SearchService 已集成灰度指标收集，自动记录：

1. 每次搜索请求的引擎类型和耗时
2. 灰度决策结果
3. 配置更新时的灰度百分比

无需手动调用，SearchService 会自动处理。

## 相关文件

- `pkg/metrics/search_grayscale.go` - 指标收集实现
- `pkg/metrics/search_grayscale_example_test.go` - 测试和示例
- `service/search/search.go` - SearchService 集成
- `service/search/grayscale.go` - 灰度决策器
