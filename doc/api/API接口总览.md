# API 接口总览

统一罗列系统对外 REST API 的分组、路径规范、响应格式与错误码，并链接各专项设计。

## 1. 路径与版本
- 前缀：/api/v1
- 分组：
  - /system/user：注册、登录、个人资料、刷新令牌
  - /document/project：项目 CRUD、索引
  - /document/node：节点增删改查、移动、排序
  - /document/file：文档读写、版本推进
  - /ai：文本/图像生成等

详见：
- 《JWT身份认证设计》./功能设计/JWT身份认证设计.md
- 《API导入设计（AI）》./ai设计/01.api导入设计.md

## 2. 统一响应格式（建议）
```
{
  "code": 0,
  "message": "success",
  "data": {},
  "timestamp": 1234567890
}
```

## 3. 错误码规范（建议）
- 4000x：认证与授权（如 40001 token 过期）
- 4100x：参数与校验
- 4200x：资源不存在/冲突
- 5000x：服务内部错误

## 4. 通用约定

- 认证：Authorization: Bearer `<token>`
- 幂等：PUT/DELETE 的幂等性约定；写操作返回资源最新状态
- 分页：page/size + 排序字段（默认 created_at desc）

## 5. OpenAPI/文档化（规划）
- 后续补充 swagger/openapi 生成与端点列表
