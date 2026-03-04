# Writer Editor V2 交接单（2026-03-04）

## 建议提交切分

### Commit 1: Writer 编辑器核心重构

- `api/v1/writer/editor_api.go`
- `service/writer/document/document_service.go`
- `service/writer/document/paragraph_mapper.go`
- `models/writer/paragraph.go`
- `models/writer/document_content.go`
- `models/dto/content_dto.go`
- `service/writer/document/document_dto.go`
- `repository/mongodb/writer/document_content_repository_mongo.go`
- `repository/interfaces/writer/errors.go`
- `service/mock/document_content_repository_mock.go`
- `service/writer/document/document_service_content_test.go`

### Commit 2: 评论/关键词链路与 Writer 路由补全

- `api/v1/writer/comment_api.go`
- `service/writer/comment_service.go`
- `models/writer/document_comment.go`
- `repository/mongodb/writer/comment_repository_mongo.go`
- `api/v1/writer/comment_api_test.go`
- `service/writer/comment_service_test.go`
- `api/v1/writer/keyword_api.go`
- `api/v1/writer/keyword_api_test.go`
- `router/writer/keyword_router.go`
- `router/writer/writer.go`
- `router/writer/writer_router_init.go`

### Commit 3: 推荐榜模型与路由

- `models/recommendation/table.go`
- `repository/interfaces/recommendation/table_repository.go`
- `repository/mongodb/recommendation/table_repository_mongo.go`
- `service/recommendation/table_service.go`
- `service/recommendation/table_service_test.go`
- `api/v1/recommendation/recommendation_api.go`
- `api/v1/recommendation/table_api.go`
- `router/recommendation/recommendation_router.go`
- `router/recommendation/recommendation_router_test.go`
- `router/enter.go`

### Commit 4: 文档与依赖

- `docs/architecture/data_model.md`
- `docs/architecture/2026-03-04-writer-paragraph-type-design.md`
- `docs/architecture/2026-03-04-writer-editor-v2-compatibility-checklist.md`
- `docs/architecture/2026-03-04-writer-editor-v2-merge-status.md`
- `go.mod`
- `go.sum`

## 提交前验收命令

```bash
go test ./api/v1/writer ./api/v1/recommendation ./service/writer ./service/writer/document ./service/recommendation ./repository/mongodb/writer ./repository/mongodb/recommendation ./router/writer ./router/recommendation
go test ./test/integration
go test ./tests/integration
```

## 当前结论

- 已满足“编辑器重构 + 推荐榜扩展 + 评论段落绑定 + 兼容层保留”的交付目标。
- 当前无阻塞项，可合并。
