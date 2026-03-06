# Publication Flow E2E

最小联调脚本：

- [e2e_publication_flow.ps1](/E:/Github/Qingyu/_wt_qy_backend_publication_mvp/scripts/e2e_publication_flow.ps1)

覆盖链路：

1. 作者提交项目发布
2. 作者提交文档发布
3. 管理员查询待审核发布单
4. 管理员审批项目发布单
5. 读侧通过 `bookstore` 查询已发布书籍和章节
6. 读侧通过 `reader` 路由验证章节读取

示例：

```powershell
powershell -ExecutionPolicy Bypass -File .\scripts\e2e_publication_flow.ps1 `
  -BaseUrl "http://localhost:8080" `
  -AuthorToken "<author-jwt>" `
  -AdminToken "<admin-jwt>" `
  -ProjectId "<project-id>" `
  -DocumentId "<document-id>" `
  -ChapterNumber 1
```

前提：

- 服务已启动
- `ProjectId` 和 `DocumentId` 对应同一作者项目
- `DocumentId` 对应的文档内容已存在于 `document_contents`
- `AuthorToken` 具备作者身份
- `AdminToken` 具备管理员身份

当前限制：

- 审批脚本默认审批项目发布单；文档发布单只做提交，不自动审批
- reader 校验使用当前实现的现有读侧接口，不补额外适配
- 如果读侧章节接口要求的 `chapterId` 与 `project document id` 不一致，需要先从 `bookstore/books/:id/chapters` 响应里取真实章节 ID 再做下一跳
