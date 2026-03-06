# Publication Flow E2E

最小联调脚本：

- [e2e_publication_flow.py](/E:/Github/Qingyu/_wt_qy_backend_publication_mvp/scripts/e2e_publication_flow.py)
- [e2e_publication_flow.ps1](/E:/Github/Qingyu/_wt_qy_backend_publication_mvp/scripts/e2e_publication_flow.ps1) 已废弃，仅输出迁移提示
- [publication_flow_smoke.py](/E:/Github/Qingyu/_wt_qy_backend_publication_mvp/scripts/testing/publication_flow_smoke.py) 用于 CI / 一键 smoke 运行

覆盖链路：

1. 作者提交项目发布
2. 作者提交文档发布
3. 管理员查询待审核发布单
4. 管理员审批项目发布单
5. 读侧通过 `bookstore` 查询已发布书籍和章节
6. 读侧通过 `reader` 路由验证章节读取

示例：

```bash
python .\scripts\e2e_publication_flow.py --base-url "http://localhost:9090"
```

自动审批文档发布单：

```bash
python .\scripts\e2e_publication_flow.py --base-url "http://localhost:9090" --approve-document
```

显式传入 token / ID 的模式仍然支持：

```bash
python .\scripts\e2e_publication_flow.py \
  --base-url "http://localhost:9090" \
  --author-token "<author-jwt>" \
  --admin-token "<admin-jwt>" \
  --project-id "<project-id>" \
  --document-id "<document-id>" \
  --chapter-number 1
```

前提：

- 服务已启动
- 默认种子账号存在：`author_new / Author@123456`、`admin / Admin@123456`
- 默认联调项目存在：`联调发布示例项目`
- 默认联调文档存在：`第1章 风起青川`
- 如果显式传入 `ProjectId` / `DocumentId`，它们需要对应同一作者项目
- 目标文档内容已存在于 `document_contents`

当前限制：

- 默认只审批项目发布单；加 `--approve-document` 后会继续审批文档发布单
- reader 校验使用当前实现的现有读侧接口，不补额外适配
- 脚本会优先读取发布记录里的 `externalId` 作为真实 `bookId`，再从 `bookstore/books/:id/chapters` 响应里解析真实 `chapterId`
- 脚本默认会在执行前检查项目是否已发布；若已发布，会自动调用 `unpublish` 做清场。可通过 `--skip-reset` 关闭

CI 入口：

- workflow: [publication-flow-smoke.yml](/E:/Github/Qingyu/_wt_qy_backend_publication_mvp/.github/workflows/publication-flow-smoke.yml)
- 本地/CI 通用 runner: [publication_flow_smoke.py](/E:/Github/Qingyu/_wt_qy_backend_publication_mvp/scripts/testing/publication_flow_smoke.py)
- runner 会完成：
  - 非交互 seed
  - 构建并启动后端
  - 等待健康检查
  - 调用 `e2e_publication_flow.py --approve-document`
