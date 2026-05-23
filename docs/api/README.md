# 青羽后端 API 文档入口

> 最后整理: 2026-05-22

## Page Role

- current-hub
- current-owner: `docs/api/`
- current-bounded: 当前后端 API 文档主入口，负责接口总览、Swagger 导出物和专题 API 导航

## Recommended Read Path

1. 先读 `API接口总览.md`。
2. 再读 `usage_guide.md`。
3. 按业务域进入对应专题目录或专题文档。

## Quick Section Map

- 当前入口说明
- First Read Path
- 按主题查找
- Topic Directories
- Swagger 与导出物
- 与其他目录的关系
- 已知缺口

## Quick Takeaways

- 这是当前 API 文档 canonical owner。
- 历史设计和实施说明都不应反向取代这里。

## Skip Guide

- 只看历史设计：跳去 `../design/`。
- 只看落地实施：跳去 `../implementation/`。

## 当前入口说明

- `Qingyu_backend/docs/api/` 是当前 API 文档、Swagger 导出物和专题接口说明的 canonical owner。
- `Qingyu_backend/docs/guides/api/` 仅保留历史镜像与迁移提示，不再作为主入口维护。
- 面向前端或联调同学的“从哪开始看”，优先走本 README，而不是旧的 `frontend/`、`postman/` 子目录占位链接。

## First Read Path

1. [API接口总览](./API接口总览.md)
2. [usage_guide.md](./usage_guide.md)
3. [shared/统一响应处理指南.md](./shared/统一响应处理指南.md)
4. [system/用户系统API文档.md](./system/用户系统API文档.md)
5. [bookstore/书城系统API文档.md](./bookstore/书城系统API文档.md)
6. [reader/阅读器系统API文档.md](./reader/阅读器系统API文档.md)

## 按主题查找

| 主题 | 入口文档 | 说明 |
|------|----------|------|
| 总览与使用 | [API接口总览](./API接口总览.md) / [usage_guide.md](./usage_guide.md) | 快速了解路由分组、使用方式和联调建议 |
| 统一响应与共享语义 | [shared/统一响应处理指南.md](./shared/统一响应处理指南.md) / [shared/共享服务API接口文档.md](./shared/共享服务API接口文档.md) | 统一响应结构、共享服务接口口径 |
| 用户与认证 | [system/用户系统API文档.md](./system/用户系统API文档.md) / [用户管理API使用指南.md](./用户管理API使用指南.md) | 用户、权限、账号相关接口 |
| 阅读与书城 | [bookstore/书城系统API文档.md](./bookstore/书城系统API文档.md) / [reader/阅读器系统API文档.md](./reader/阅读器系统API文档.md) / [阅读端API使用文档.md](./阅读端API使用文档.md) | 书城、阅读器、章节与阅读流程 |
| 推荐与统计 | [recommendation/推荐系统API文档.md](./recommendation/推荐系统API文档.md) / [统计API文档.md](./统计API文档.md) | 推荐链路与统计类接口 |
| 写作与 AI | [写作端API完整文档.md](./写作端API完整文档.md) / [Phase3创作API文档.md](./Phase3创作API文档.md) / [AI_WRITING_ASSISTANT_APIS.md](./AI_WRITING_ASSISTANT_APIS.md) | 写作工作流、创作与 AI 辅助接口 |
| 管理与专题 | [管理员API文档.md](./管理员API文档.md) / [admin/配置管理API文档.md](./admin/配置管理API文档.md) / [审核API文档.md](./审核API文档.md) | 管理后台、配置、审核 |
| 文档类 API | [document/document_api.md](./document/document_api.md) / [document/character_api.md](./document/character_api.md) | 文档、角色等专题接口 |

## Topic Directories

- [admin/README.md](./admin/README.md)
- [ai/README.md](./ai/README.md)
- [bookstore/README.md](./bookstore/README.md)
- [document/README.md](./document/README.md)
- [reader/README.md](./reader/README.md)
- [recommendation/README.md](./recommendation/README.md)
- [shared/README.md](./shared/README.md)
- [system/README.md](./system/README.md)

## Swagger 与导出物

- [swagger.yaml](./swagger.yaml)
- [swagger.json](./swagger.json)
- [SWAGGER_API_导出说明.md](./SWAGGER_API_导出说明.md)
- [document/openapi.yaml](./document/openapi.yaml)

## 与其他目录的关系

- 设计稿请回看 [../design/reader/README.md](../design/reader/README.md)。
- 实施落地请回看 [../implementation/03-reading/READING_STATS_IMPLEMENTATION.md](../implementation/03-reading/READING_STATS_IMPLEMENTATION.md) 与 [../implementation/06-bookstore/BOOKSTORE_API_IMPLEMENTATION_SUMMARY.md](../implementation/06-bookstore/BOOKSTORE_API_IMPLEMENTATION_SUMMARY.md)。
- 测试入口请回看 [../testing/README.md](../testing/README.md)。

## 已知缺口

- 旧版 `frontend/` 子目录入口已不存在；如果需要“前端专用 API 速查模板”，当前记为 `TBD`。
  - 确认路径：优先检查父仓 `docs/plans/submodules/frontend/` 与前端仓 `Qingyu_fronted` 中是否已有稳定 owner，再决定是否在本目录补新的联调索引。
- 旧版 Postman 集合链接已失效；如果仍需集合文件，当前记为 `TBD`。
  - 确认路径：优先检查根仓或 API 生成链路是否已有可导出产物，不在本目录重复制造第二套入口。
