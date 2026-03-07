# Writer Editor V2 兼容层下线清单（2026-03-04）

## 目标

在删除旧接口 `/content` 及旧字段（`chapterId`/`position`）之前，确保前后端已完全切换到新方案：
- 文档内容：`/contents`
- 批注锚点：`paragraphId`

## A. 接口切换确认

- 前端已不再调用：
- `GET /api/v1/writer/documents/:id/content`
- `PUT /api/v1/writer/documents/:id/content`
- 前端已切换到：
- `GET /api/v1/writer/documents/:id/contents`
- `PUT /api/v1/writer/documents/:id/contents`
- `POST /api/v1/writer/documents/:id/contents/reindex`

## B. 批注模型切换确认

- 创建批注请求仅发送 `paragraphId`。
- 不再发送 `chapterId` 与 `position`。
- 回复批注后，子评论正确继承父评论 `paragraph_id`。
- 评论列表查询按 `document_id + paragraph_id` 可过滤。

## C. 数据一致性确认

- 同一文档段落顺序无重复（`paragraph_order` 唯一于文档内）。
- `ReplaceDocumentContents` 提交不存在空段落、重复段落 ID、重复顺序。
- 评论中的 `paragraph_id` 能在对应文档段落集合中找到（抽样核查）。

## D. 观测与灰度

- 兼容接口保留期内，统计旧接口调用次数（建议加网关/日志指标）。
- 连续 7 天旧接口调用为 0，再进入删除流程。
- 删除前打发布快照（可快速回滚）。

## E. 删除动作（最终）

- 删除 `EditorApi.GetDocumentContent/UpdateDocumentContent` 路由绑定。
- 删除批注请求中 `chapterId/position` 兼容校验逻辑与相关文案。
- 更新 Swagger/Apifox 文档，仅保留 `/contents` 与 `paragraphId` 方案。

## F. 回滚策略

- 删除动作在独立分支完成并保留 tag。
- 回滚优先恢复路由绑定，再恢复请求兼容字段解析。
