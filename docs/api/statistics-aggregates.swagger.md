# Statistics Aggregates Swagger Supplement

## 文件
- `docs/api/statistics-aggregates.swagger.yaml`

## 覆盖范围
- `/api/v1/writer/stats/overview`
- `/api/v1/writer/stats/views`
- `/api/v1/writer/stats/subscribers`
- `/api/v1/writer/stats/chapters`
- `/api/v1/writer/stats/today`
- `/api/v1/reader/statistics`
- `/api/v1/reader/statistics/overview`
- `/api/v1/reader/statistics/reading-time`
- `/api/v1/reader/statistics/heatmap`
- `/api/v1/reader/statistics/trends`

## 说明
- 本文件是对现有 Swagger 文档的补充。
- 当前项目执行 `swag init -g cmd/server/main.go -o docs` 时，会被已有注解和 `swaggo/swag` 解析问题阻塞，无法稳定重新生成全量产物。
- 为避免阻塞联调，这里先提供新增统计聚合接口的独立 Swagger YAML，可供前端、测试和后端联调直接参考。

## 全量生成现状
- `swag init -g cmd/server/main.go -o docs --parseDependency=false`:
  - 失败，已有注解引用外部类型时无法解析。
- `swag init -g cmd/server/main.go -o docs --parseDependency --parseInternal`:
  - 失败，`swaggo/swag` 在当前项目上发生递归崩溃。

## 建议
- 短期：联调阶段直接使用补充 YAML。
- 中期：清理历史 Swagger 注解问题后，再恢复全量 `swag init` 自动生成。
