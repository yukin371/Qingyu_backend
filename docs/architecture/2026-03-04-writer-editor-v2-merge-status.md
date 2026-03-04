# Writer Editor V2 重构状态（2026-03-04）

## 完成项

- 路由收口到 `/api/v1/writer/*`，旧 `/content` 接口保留兼容并添加 deprecated 提示头。
- 段落轻量类型 `Paragraph` 已引入，持久化仍落在 `document_content`。
- `ReplaceDocumentContents` 增加边界校验：
- 空段落拒绝
- 重复段落 ID 拒绝
- 重复顺序拒绝
- 评论链路按 `paragraph_id` 工作：
- 创建强制 `paragraphId`
- 回复继承父评论段落
- 查询支持段落过滤
- 推荐榜支持自动榜与手动榜，路由冲突已修复并有回归测试。
- Repository/Service 分层第二轮优化已落地（分页归一化在 service，repo 专注数据访问）。

## 测试结果

- `go test ./api/v1/writer ./api/v1/recommendation ./service/writer ./service/writer/document ./service/recommendation ./repository/mongodb/writer ./repository/mongodb/recommendation ./router/writer ./router/recommendation` 通过
- `go test ./test/integration` 通过
- `go test ./tests/integration` 通过

## 可合并前阻塞项

当前无阻塞项。

## 残余风险（非阻塞）

- `api/v1/recommendation` 与 `repository/mongodb/recommendation` 当前测试覆盖较少（主要由 service/router 层间接覆盖）。
- 旧兼容接口尚未下线，需要按兼容清单完成观察窗口后再删除。

## 建议合并策略

- 先合并本轮重构（保持兼容层）。
- 观察期内统计旧接口调用。
- 观察期结束后按兼容清单执行删除旧接口与旧字段。
