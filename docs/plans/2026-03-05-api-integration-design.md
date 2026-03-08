# 前后端 API 集成标准化设计方案 (修订版 v2)

**创建日期**: 2026-03-05
**修订日期**: 2026-03-05
**修订原因**: 根据项目现状审查报告进行调整
**目标**: 解决前后端数据类型不一致问题，建立统一的 API 集成标准

---

## 修订记录

| 版本  | 日期         | 修订内容                                           | 修订人  |
| --- | ---------- | ---------------------------------------------- | ---- |
| v1  | 2026-03-05 | 初始版本                                           | Kore |
| v2  | 2026-03-05 | 基于审查报告调整：修正OpenAPI结构假设、重新评估Transformer必要性、调整工期 | Kore |

---

## 0. 执行摘要

**核心问题**：前后端存在28个P0/P1级别的类型不一致问题

**解决方案**：
1. **保持现有swagger.yaml为主**，增加模块化OpenAPI扩展
2. **后端DTO统一**：渐进式迁移到`models/dto/`
3. **前端Wrapper优化**：基于现有Orval配置，选择性引入Transformer
4. **建立契约测试**：CI集成自动化验证
5. **Feature Flag迁移**：支持并行双轨切换

**关键调整**：
- ✅ 后端DTO已使用camelCase JSON标签，前端无需额外命名转换
- ✅ 基于现有swagger.yaml，不强制模块化拆分
- ✅ 工期从10周调整为16周（增加60%缓冲）
- ✅ Writer模块优先创建DTO定义

**预计工期**：16周（4个月）

---

## 1. 背景与问题

### 1.1 核心问题

根据 `docs/reports/2026-03-04-frontend-backend-data-type-comparison-report.md`，当前存在 28 个 P0/P1 级别的类型不一致问题：

- Writer 模块: 4 个 P0, 3 个 P1
- Reader 模块: 5 个 P0, 2 个 P1
- Bookstore 模块: 3 个 P0, 1 个 P1
- Social 模块: 2 个 P0, 2 个 P1
- Admin 模块: 1 个 P0, 3 个 P1
- Recommendation 模块: 2 个 P0

### 1.2 现有标准

项目已建立以下标准（位于 `Qingyu_backend/docs/standards/`）：
- `api-status-code-standard.md`: HTTP 状态码使用规范
- `error_code_standard.md`: 4 位错误码体系
- `api/API设计规范.md`: URL 结构和响应格式规范

### 1.3 现有基础设施（实际状态）

**前端**:
- ✅ Orval 已配置 (`orval.config.ts`)，覆盖8个模块
- ✅ 统一 mutator: `src/core/config/orval-mutator.ts`
- ✅ 生成模式: `mode: 'single'`，输出到各模块 `api/generated/`
- ✅ 统一 HTTP 客户端: `core/services/http.service.ts` (3万+行)
- ⚠️ Wrapper 直接导出 Orval 生成的 API，无额外转换层

**后端**:
- ✅ 现有 swagger.yaml (18433行)
- ✅ DTO 部分已统一到 `models/dto/` (reader, bookstore, user, audit, content)
- ⚠️ DTO 仍分散在 `api/v1/{domain}/dto/`, `service/{domain}/dto/`
- ❌ Writer 模块缺少统一的 DTO 定义
- ✅ 后端 DTO 已使用 camelCase JSON 标签

---

## 2. 设计原则

1. **尊重现有结构**: 基于现有 swagger.yaml，不强制重构
2. **渐进式迁移**: DTO 统一、Wrapper 优化均采用渐进式
3. **避免过度设计**: 后端已使用 camelCase，前端无需重复转换
4. **契约驱动测试**: 建立 CI 集成的自动化验证
5. **Feature Flag 降级**: 支持快速回滚

**关键调整**：
- ❌ ~~前端负责命名转换~~ → ✅ 后端负责（DTO层已使用camelCase）
- ❌ ~~模块化 OpenAPI 拆分~~ → ✅ 保持现有 swagger.yaml，按需扩展
- ❌ ~~强制 Transformer 模式~~ → ✅ 选择性使用，仅处理复杂类型转换

---

## 3. 整体架构（修订版）

```
┌─────────────────────────────────────────────────────────────────┐
│                        前端层 (Frontend)                        │
│  ┌────────────┐    ┌────────────┐    ┌────────────┐            │
│  │ Components │───▶│  Wrapper   │───▶│ Orval Gen. │            │
│  └────────────┘    └────────────┘    └────────────┘            │
│                           │                  │                  │
│                           │         ┌────────┴─────────┐        │
│                           │         │  swagger.yaml    │        │
│                           │         │  (主要数据源)    │        │
│                           │         └─────────────────┘        │
│                           │                                     │
│                    ┌──────┴──────┐                              │
│                    │ Transformer │ (可选，仅复杂类型)           │
│                    └─────────────┘                              │
└─────────────────────────────────────────────────────────────────┘
                                                                 │
┌────────────────────────────────────────────────────────────────┤
│                        后端层 (Backend)                        │
│  ┌────────────┐    ┌────────────┐    ┌────────┴───────────┐   │
│  │   Router   │───▶│   Handler  │───▶│   DTO Layer        │   │
│  └────────────┘    └────────────┘    │   models/dto/       │   │
│                                      │   (camelCase JSON)  │   │
│                                      └─────────────────────┘   │
└─────────────────────────────────────────────────────────────────┘
```

**架构变化说明**：
- 后端 DTO 层负责命名转换（snake_case → camelCase）
- 前端 Wrapper 简化为直接使用 Orval 生成代码
- Transformer 仅用于复杂业务逻辑转换（如金额、时间格式）

---

## 4. OpenAPI 规范层（修订）

### 4.1 现状与策略

**现状**:
- ✅ 存在: `docs/swagger.yaml` (18433行，完整的 API 规范)
- ✅ 存在: `docs/swagger/swagger.yaml`
- ✅ 存在: `docs/api/document/openapi.yaml`
- ❌ 不存在: `docs/openapi/{domain}/openapi.yaml` 模块化结构

**策略**:
1. **保持 swagger.yaml 为主**: 作为单一真相来源
2. **可选模块化扩展**: 按需创建 `docs/openapi/{domain}/` 用于增量开发
3. **CI 自动合并**: 模块化规范可合并回主 swagger.yaml

### 4.2 目录结构（修订版）

```
docs/
├── swagger.yaml                      # 主要 OpenAPI 规范（保持）
├── swagger/
│   └── swagger.yaml                  # 备份位置
├── openapi/                          # 新增：模块化扩展（可选）
│   ├── writer/
│   │   └── openapi.yaml              # Writer 模块增量规范
│   ├── common/
│   │   ├── common.yaml               # 通用组件
│   │   ├── errors.yaml               # 错误码定义
│   │   └── pagination.yaml           # 分页模型
│   └── README.md                     # 规范编写指南
└── api/
    └── document/
        └── openapi.yaml              # 现有文档
```

### 4.3 通用响应格式（保持）

参考现有 swagger.yaml 的响应格式，确保一致性。

---

## 5. 后端 DTO 统一（优先级提升）

### 5.1 现状分析

| 位置 | 模块 | 状态 | 优先级 |
|------|------|------|--------|
| models/dto/ | reader, bookstore, user, audit, content | ✅ 已统一 | - |
| models/dto/ | writer | ❌ **缺失** | **P0** |
| api/v1/user/dto/ | user | ⚠️ 混合 | P1 |
| api/v1/social/dto/ | social | ⚠️ 未迁移 | P2 |
| api/v1/notifications/dto/ | notifications | ⚠️ 未迁移 | P2 |
| service/ai/dto/ | ai | ⚠️ 未迁移 | P2 |

### 5.2 DTO 定义规范（保持）

```go
// Qingyu_backend/models/dto/writer_dto.go
// **新增文件** - Writer 模块 DTO 定义

package dto

import "time"

// CreateProjectRequest 创建项目请求
type CreateProjectRequest struct {
    Title    string   `json:"title" validate:"required,min=1,max=100"`
    Summary  string   `json:"summary,omitempty" validate:"max=500"`
    CoverURL string   `json:"coverUrl,omitempty" validate:"url,max=500"`  // camelCase
    Tags     []string `json:"tags,omitempty" validate:"max=10,dive,min=1,max=50"`
}

// UpdateProjectRequest 更新项目请求
type UpdateProjectRequest struct {
    Title    *string   `json:"title,omitempty" validate:"omitempty,min=1,max=100"`
    Summary  *string   `json:"summary,omitempty" validate:"omitempty,max=500"`
    CoverURL *string   `json:"coverUrl,omitempty" validate:"omitempty,url,max=500"`
    Tags     *[]string `json:"tags,omitempty" validate:"omitempty,max=10,dive,min=1,max=50"`
}

// ProjectResponse 项目响应
type ProjectResponse struct {
    ID        string    `json:"id"`
    Title     string    `json:"title"`
    Summary   string    `json:"summary"`
    CoverURL  string    `json:"coverUrl"`      // camelCase
    Tags      []string  `json:"tags"`
    CreatedAt time.Time `json:"createdAt"`     // camelCase
    UpdatedAt time.Time `json:"updatedAt"`     // camelCase
}

// CreateDocumentRequest 创建文档请求
type CreateDocumentRequest struct {
    ProjectID string  `json:"projectId" validate:"required"`      // camelCase
    Title     string  `json:"title" validate:"required,min=1,max=200"`
    Content   string  `json:"content,omitempty"`
}

// DocumentResponse 文档响应
type DocumentResponse struct {
    ID        string    `json:"id"`
    ProjectID string    `json:"projectId"`      // camelCase
    Title     string    `json:"title"`
    Content   string    `json:"content"`
    WordCount int       `json:"wordCount"`      // camelCase
    CreatedAt time.Time `json:"createdAt"`      // camelCase
    UpdatedAt time.Time `json:"updatedAt"`      // camelCase
}
```

### 5.3 命名转换责任明确

**原则**: 后端 DTO 层负责所有命名转换

```go
// MongoDB Model (BSON snake_case)
type Project struct {
    ID        primitive.ObjectID `bson:"_id"`
    Title     string             `bson:"title"`
    CoverURL  string             `bson:"cover_url"`      // BSON: snake_case
    CreatedAt time.Time          `bson:"created_at"`     // BSON: snake_case
}

// DTO (JSON camelCase)
type ProjectResponse struct {
    ID        string    `json:"id"`
    Title     string    `json:"title"`
    CoverURL  string    `json:"coverUrl"`      // JSON: camelCase
    CreatedAt time.Time `json:"createdAt"`     // JSON: camelCase
}

// 转换函数
func ToProjectResponse(p *Project) ProjectResponse {
    return ProjectResponse{
        ID:        p.ID.Hex(),
        Title:     p.Title,
        CoverURL:  p.CoverURL,      // 字段名相同，无需转换
        CreatedAt: p.CreatedAt,
    }
}
```

---

## 6. 前端 Wrapper 优化（修订）

### 6.1 当前模式分析

**现有 wrapper.ts**:
```typescript
// Qingyu_fronted/src/modules/writer/api/wrapper.ts (当前)
import { getApi } from './generated/api'

export const writerApi = getApi()
export const getProject = writerApi.getApiV1ProjectsProjectId
export const createProject = writerApi.postApiV1Projects
// ... 直接导出 Orval 生成的函数
```

**优点**:
- ✅ 简单直接，无额外层级
- ✅ 类型安全，由 Orval 保证
- ✅ 后端已使用 camelCase，无需转换

**缺点**:
- ⚠️ 无法统一处理错误
- ⚠️ 无法统一处理响应格式
- ⚠️ 难以添加业务逻辑

### 6.2 优化策略（渐进式）

**Phase 1: 保持现状 + 统一错误处理** (Week 1-2)
```typescript
// Qingyu_fronted/src/core/config/orval-mutator.ts (增强)
export const orvalMutator = async (args: MutatorApiArgs) => {
  // 现有逻辑...

  // 新增：统一错误处理
  if (response.data?.code !== undefined && response.data.code !== 200) {
    throw new APIError(
      response.data.code,
      response.data.message,
      response.data.request_id
    )
  }

  return response
}
```

**Phase 2: 可选 Wrapper 封装** (Week 3-4)
```typescript
// Qingyu_fronted/src/modules/writer/api/wrapper.ts (渐进式)
import { getApi } from './generated/api'

const api = getApi()

// 导出原始 API（保持兼容）
export const rawApi = api

// 可选：添加业务封装
export class WriterAPI {
  async createProject(data: CreateProjectRequest): Promise<Project> {
    const response = await api.postApiV1Projects({ body: data })
    // 统一处理响应格式
    return response.data as Project
  }

  async getProject(id: string): Promise<Project> {
    const response = await api.getApiV1ProjectsProjectId({ path: { id } })
    return response.data as Project
  }
}

// 根据环境变量决定导出模式
export const writerApi = process.env.VUE_APP_USE_WRAPPER
  ? new WriterAPI()
  : api
```

### 6.3 Transformer 策略（修订）

**原则**: 仅在需要时使用

**不需要 Transformer 的场景**:
- ✅ 命名转换（后端已处理）
- ✅ 简单类型映射

**需要 Transformer 的场景**:
- 金额单位转换（分 → 元）
- 时间格式转换（时间戳 → Date 对象）
- 复杂对象组装

```typescript
// Qingyu_fronted/src/modules/writer/api/transformers.ts (仅复杂场景)
/**
 * 金额转换：后端返回分，前端显示元
 */
export function toYuan(cents: number): number {
  return cents / 100
}

/**
 * 时间转换：后端返回 ISO 字符串，前端转为 Date
 */
export function toDate(isoString: string): Date {
  return new Date(isoString)
}
```

---

## 7. Feature Flags 与迁移策略（修订）

### 7.1 现状分析

**现状**: 前端项目中**没有** Feature Flag 相关实现

**决策**: 简化为环境变量控制，降低复杂度

### 7.2 简化的 Feature Flag 方案

```typescript
// Qingyu_fronted/src/.env.development
VUE_APP_API_V2_WRITER=false
VUE_APP_API_V2_READER=false

// Qingyu_fronted/src/.env.production
VUE_APP_API_V2_WRITER=true
VUE_APP_API_V2_READER=true

// Qingyu_fronted/src/modules/writer/api/index.ts
const useV2 = import.meta.env.VUE_APP_API_V2_WRITER === 'true'

export const writerApi = useV2
  ? new WriterAPI()      // 新版本
  : rawApi               // 原始版本（Orval 直接导出）
```

### 7.3 回滚策略

**紧急回滚**:
1. 修改环境变量: `VUE_APP_API_V2_WRITER=false`
2. 重新部署前端
3. 无需修改代码

**灰度发布**:
1. 按用户 ID 分流: `useV2 = userId % 10 < 3` (30%灰度)
2. 按地域分流: `useV2 = userRegion === 'CN'`
3. 按时间分流: `useV2 = new Date().getHours() >= 10` (白天开)

---

## 8. 类型同步与代码生成（修订）

### 8.1 简化的类型同步流程

```
┌─────────────────────────────────────────────────────────────┐
│                     swagger.yaml                            │
│                  (单一真相来源)                              │
└────────────┬────────────────────────────────────────────────┘
             │
             ├──────────────────┬──────────────────┐
             ▼                  ▼                  ▼
    ┌────────────────┐  ┌──────────────┐  ┌──────────────┐
    │  后端 Handler  │  │ Orval Generate│  │ 契约测试    │
    │  (已使用 DTO)  │  │ (前端类型)    │  │ (CI验证)     │
    └────────────────┘  └──────────────┘  └──────────────┘
```

### 8.2 类型验证脚本（简化版）

```bash
# Qingyu_backend/scripts/validate-api-contract.sh
#!/bin/bash

set -e

echo "🔍 开始验证 API 契约..."

# 1. 检查 swagger.yaml 格式
echo "📋 检查 swagger.yaml 格式..."
npx swagger-parser validate docs/swagger.yaml || echo "⚠️ 需要安装 swagger-parser: npm i -g swagger-parser"

# 2. 生成 TypeScript 类型（前端）
echo "📝 生成前端类型..."
cd ../Qingyu_fronted
npx orval --config orval.config.ts

# 3. 类型检查
echo "🔍 运行类型检查..."
npm run type-check

echo "✅ API 契约验证通过！"
```

---

## 9. 测试策略（修订）

### 9.1 契约测试（CI 集成）

```yaml
# .github/workflows/api-contract-test.yml
name: API Contract Test

on:
  pull_request:
    paths:
      - 'Qingyu_backend/api/v1/**'
      - 'Qingyu_backend/models/dto/**'
      - 'docs/swagger.yaml'
      - 'Qingyu_fronted/orval.config.ts'

jobs:
  contract-test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3

      - name: 验证 swagger.yaml
        run: |
          npx swagger-parser validate docs/swagger.yaml

      - name: 生成前端类型
        run: |
          cd Qingyu_fronted
          npm ci
          npx orval

      - name: 前端类型检查
        run: |
          cd Qingyu_fronted
          npm run type-check
```

### 9.2 E2E 测试（关键场景）

```typescript
// Qingyu_fronted/test/e2e/api-contract.spec.ts
import { test, expect } from '@playwright/test'

test.describe('API 契约测试', () => {
  test('Writer API - 创建项目', async ({ request }) => {
    const response = await request.post('/api/v1/projects', {
      data: { title: '测试项目' }
    })

    expect(response.status()).toBe(200)

    const body = await response.json()
    expect(body).toMatchObject({
      code: 201,
      message: 'created',
      data: {
        id: expect.any(String),
        title: '测试项目',
        createdAt: expect.any(String),
      },
      timestamp: expect.any(String),
      request_id: expect.any(String),
    })
  })
})
```

---

## 10. 实施计划（修订版 v2）

### 10.1 阶段划分（调整）

| Phase | 内容 | 工期 | 交付物 | 关键风险 |
|-------|------|------|--------|----------|
| 1 | 后端 DTO 统一 | Week 1-3 | Writer DTO + swagger 更新 | Writer API 调用方兼容 |
| 2 | 前端 Wrapper 优化 | Week 4-5 | 统一错误处理 | 响应格式不一致 |
| 3 | 契约测试建立 | Week 6 | CI 集成 + 类型验证 | CI 配置复杂 |
| 4 | Writer 试点验证 | Week 7-8 | E2E 测试 + 灰度发布 | 回滚方案验证 |
| 5 | 全域推广 | Week 9-14 | 其他 5 个模块 | 模块间依赖 |
| 6 | 收尾清理 | Week 15-16 | 移除 Feature Flags | 遗留问题修复 |

**工期调整理由**:
- Phase 1: 增加了 Writer DTO 创建（原计划未包含）
- Phase 2: 简化了 Transformer，但增加了统一错误处理
- Phase 3: 新增契约测试建立
- Phase 5-6: 增加了 60% 缓冲时间

### 10.2 Phase 1: 后端 DTO 统一（Week 1-3）

**Week 1: Writer DTO 创建**

| 任务 | 负责人 | 交付物 | 验收标准 |
|------|--------|--------|----------|
| 创建 `models/dto/writer_dto.go` | 后端 | DTO 定义文件 | 包含 Project/Document 等 |
| 更新 swagger.yaml | 后端 | 更新的规范 | 包含 Writer API 定义 |
| 单元测试 | 后端 | `writer_dto_test.go` | 覆盖率 > 80% |

**Week 2: Writer API 迁移**

| 任务 | 负责人 | 交付物 | 验收标准 |
|------|--------|--------|----------|
| 更新 Handler 使用新 DTO | 后端 | 修改后的 handler | 保持 API 兼容 |
| 集成测试 | 后端 | 集成测试用例 | 现有功能不破坏 |
| 更新 API 文档 | 后端 | swagger.yaml | 与 DTO 一致 |

**Week 3: 验证与修复**

| 任务 | 负责人 | 交付物 | 验收标准 |
|------|--------|--------|----------|
| 运行现有测试套件 | 后端 | 测试报告 | 无回归 |
| 手动验证关键 API | QA | 验证报告 | 功能正常 |
| 修复发现的问题 | 后端 | 修复补丁 | 问题关闭 |

### 10.3 Phase 2: 前端 Wrapper 优化（Week 4-5）

**Week 4: 统一错误处理**

| 任务 | 负责人 | 交付物 | 验收标准 |
|------|--------|--------|----------|
| 增强 `orval-mutator.ts` | 前端 | 更新的 mutator | 统一错误处理 |
| 更新 API Error 类 | 前端 | `api-error.ts` | 支持标准错误码 |
| 单元测试 | 前端 | mutator 测试 | 覆盖率 > 80% |

**Week 5: 前端类型验证**

| 任务 | 负责人 | 交付物 | 验收标准 |
|------|--------|--------|----------|
| 运行 Orval 生成 | 前端 | 更新的 API 类型 | 无生成错误 |
| 类型检查 | 前端 | type-check 结果 | 无类型错误 |
| 修复类型问题 | 前端 | 修复补丁 | type-check 通过 |

### 10.4 Phase 3: 契约测试建立（Week 6）

| 任务 | 负责人 | 交付物 | 验收标准 |
|------|--------|--------|----------|
| 创建契约测试脚本 | DevOps | `validate-api-contract.sh` | 能独立运行 |
| CI 配置 | DevOps | GitHub Actions workflow | PR 触发运行 |
| 类型比对测试 | 前端 | types-compare.test.ts | 能检测差异 |
| 文档更新 | 后端 | README 更新 | 流程说明清晰 |

### 10.5 Phase 4: Writer 试点验证（Week 7-8）

**Week 7: 灰度发布准备**

| 任务 | 负责人 | 交付物 | 验收标准 |
|------|--------|--------|----------|
| 环境变量配置 | 前端 | `.env.*` 文件 | 可切换版本 |
| E2E 测试 | QA | Playwright 测试 | 关键场景覆盖 |
| 回滚演练 | DevOps | 回滚文档 | 5分钟内完成 |

**Week 8: 灰度发布与监控**

| 任务 | 负责人 | 交付物 | 验收标准 |
|------|--------|--------|----------|
| 10% 灰度 | DevOps | 生产部署 | 错误率 < 0.1% |
| 50% 灰度 | DevOps | 生产部署 | 无新增问题 |
| 100% 全量 | DevOps | 生产部署 | 稳定运行 24h |

### 10.6 Phase 5: 全域推广（Week 9-14）

**优先级排序**:
1. Reader (Week 9-10): P0 问题最多
2. Social (Week 11): P0 问题较多
3. Bookstore (Week 12): 已有 DTO，风险低
4. Admin (Week 13): 影响面小
5. Recommendation (Week 14): 依赖最少

### 10.7 Phase 6: 收尾清理（Week 15-16）

| 任务 | 负责人 | 交付物 | 验收标准 |
|------|--------|--------|----------|
| 移除 Feature Flags | 前端 | 清理环境变量 | 代码无残留 |
| 文档完善 | 后端 | API 文档 | 完整准确 |
| 知识分享会 | Team | 分享材料 | 团队理解 |
| 复盘会议 | Team | 复盘报告 | 经验总结 |

---

## 11. 风险与缓解措施（补充）

### 11.1 风险登记册

| 风险 | 影响 | 概率 | 缓解措施 | 负责人 |
|------|------|------|----------|--------|
| **后端 DTO 统一困难** | 🔴 高 | 🟠 中 | 渐进式迁移，保留兼容层 | 后端 Lead |
| **API 调用方不兼容** | 🔴 高 | 🟡 低 | 提前通知，提供迁移指南 | 后端 Lead |
| **Orval 生成类型不匹配** | 🟠 中 | 🟠 中 | 手动补充 types.ts | 前端 Lead |
| **Feature Flag 回滚失败** | 🟠 中 | 🟡 低 | 定期回滚演练 | DevOps |
| **性能退化** | 🔴 高 | 🟡 低 | 基准测试，性能监控 | 后端 Lead |
| **CI 契约测试误报** | 🟠 中 | 🟠 中 | 白名单机制，人工审核 | DevOps |
| **团队学习曲线** | 🟠 中 | 🟠 中 | 详细文档，知识分享会 | Tech Lead |
| **工期延误** | 🟠 中 | 🟠 中 | 已增加 60% 缓冲 | PM |
| **swagger.yaml 冲突** | 🟡 低 | 🟠 中 | 合并策略，定期同步 | 后端 Lead |
| **前端类型覆盖率不足** | 🟡 低 | 🟠 中 | 渐进式补充，优先 P0 | 前端 Lead |

### 11.2 回滚方案

**触发条件**:
- 错误率 > 1%
- 响应时间增加 > 50%
- 关键功能不可用

**回滚步骤**:
1. 修改环境变量: `VUE_APP_API_V2_*=false`
2. 重新部署前端 (5分钟)
3. 验证关键功能 (10分钟)
4. 如问题持续，回滚后端 (15分钟)

**总回滚时间**: < 30分钟

---

## 12. 成功指标

### 12.1 量化指标

| 指标 | 基线 | 目标 | 测量方式 |
|------|------|------|----------|
| API 类型一致性 | - | 100% | 类型检查通过率 |
| 契约测试覆盖率 | 0% | > 90% | CI 测试统计 |
| P0/P1 问题数 | 28 | 0 | 类型比对报告 |
| API 调用代码量 | 基线 | -20% | 代码行数统计 |
| 新 API 开发时间 | 基线 | -30% | 开发工时统计 |
| 错误处理一致性 | ~60% | 100% | 错误码覆盖率 |

### 12.2 质量指标

- ✅ 无运行时类型错误（生产环境）
- ✅ API 文档自动生成，准确性 100%
- ✅ 契约测试 CI 集成，100% 通过
- ✅ 回滚演练成功，< 30分钟

### 12.3 交付物清单

**后端**:
- [ ] `models/dto/writer_dto.go`
- [ ] `models/dto/reader_dto.go` (完善)
- [ ] `models/dto/social_dto.go` (新增)
- [ ] 更新的 `swagger.yaml`
- [ ] 契约测试脚本

**前端**:
- [ ] 增强的 `orval-mutator.ts`
- [ ] 统一的 `api-error.ts`
- [ ] 环境变量配置
- [ ] 类型比对测试

**DevOps**:
- [ ] CI 契约测试 workflow
- [ ] 监控告警配置
- [ ] 回滚文档

---

## 13. 文件结构总览（修订版）

```
docs/
├── swagger.yaml                      # 主要 OpenAPI 规范（保持）
├── swagger/
│   └── swagger.yaml                  # 备份
├── openapi/                          # 新增：模块化扩展（可选）
│   ├── common/
│   │   ├── common.yaml
│   │   ├── errors.yaml
│   │   └── pagination.yaml
│   └── README.md
│
Qingyu_backend/
├── models/dto/                       # 统一 DTO 定义
│   ├── common.go                     # 通用响应
│   ├── user_dto.go                   # 已有
│   ├── reader_dto.go                 # 已有
│   ├── bookstore_dto.go              # 已有
│   ├── writer_dto.go                 # **新增**
│   ├── social_dto.go                 # **迁移**
│   └── content_dto.go                # 已有
├── api/v1/
│   ├── writer/
│   │   └── writer_handler.go         # 使用 models/dto
│   └── social/
│       └── social_handler.go         # 迁移使用 models/dto
└── scripts/
    └── validate-api-contract.sh      # 新增
│
Qingyu_fronted/
├── src/
│   ├── core/
│   │   └── config/
│   │       └── orval-mutator.ts      # 增强：统一错误处理
│   ├── modules/
│   │   └── {domain}/api/
│   │       ├── generated/            # Orval 生成（不编辑）
│   │       │   └── api.ts
│   │       ├── wrapper.ts            # 可选封装
│   │       └── transformers.ts       # 仅复杂场景
│   └── utils/
│       └── api-error.ts              # 新增：统一错误类
├── .env.development                  # 新增：环境变量
├── .env.production
└── orval.config.ts                   # 保持现有配置
```

---

## 14. 附录

### A. 参考文档

- `docs/standards/api-status-code-standard.md`
- `Qingyu_backend/docs/standards/error_code_standard.md`
- `Qingyu_backend/docs/standards/api/API设计规范.md`
- `docs/reports/2026-03-04-frontend-backend-data-type-comparison-report.md`
- **新增**: `docs/plans/2026-03-05-api-integration-design-review.md` (审查报告)

### B. 设计决策记录（更新）

| 决策 | v1 | v2 (修订) | 理由 |
|------|-----|----------|------|
| OpenAPI 结构 | 模块化拆分 | 保持 swagger.yaml 为主 | 现有结构已完善，避免破坏性变更 |
| 命名转换责任 | 前端 | 后端 | 后端 DTO 已使用 camelCase JSON 标签 |
| Transformer 模式 | 强制使用 | 可选使用 | 大多数场景不需要，增加复杂度 |
| Feature Flag | 专用系统 | 环境变量 | 简化实现，降低维护成本 |
| 工期估算 | 10 周 | 16 周 | 增加缓冲，考虑实际复杂度 |

### C. 迁移检查清单

**后端 DTO 统一**:
- [ ] 创建 `models/dto/writer_dto.go`
- [ ] 更新 Writer Handler 使用新 DTO
- [ ] 更新 swagger.yaml
- [ ] 运行单元测试
- [ ] 运行集成测试
- [ ] 更新 API 文档

**前端 Wrapper 优化**:
- [ ] 增强 `orval-mutator.ts`
- [ ] 创建 `api-error.ts`
- [ ] 运行 Orval 生成
- [ ] 运行类型检查
- [ ] 手动验证关键 API

**契约测试**:
- [ ] 创建验证脚本
- [ ] 配置 CI workflow
- [ ] 运行首次验证
- [ ] 修复误报问题
- [ ] 编写使用文档

### D. 通信计划

| 阶段 | 通知对象 | 内容 | 频率 |
|------|----------|------|------|
| Phase 1 开始 | 后端团队 | DTO 创建计划 | 一次性 |
| Phase 2 开始 | 前端团队 | Wrapper 优化说明 | 一次性 |
| Phase 4 灰度 | 全体 | 灰度计划 + 回滚方案 | 每周 |
| Phase 5 推广 | 全体 | 各模块时间表 | 每周 |
| Phase 6 完成 | 全体 | 总结报告 | 一次性 |

---

## 修订说明

**v2 主要变更**:
1. **OpenAPI 结构**: 从模块化拆分改为保持现有 swagger.yaml
2. **命名转换**: 明确由后端 DTO 层负责，前端不再处理
3. **Transformer**: 从强制使用改为可选使用
4. **Feature Flag**: 从专用系统简化为环境变量
5. **工期**: 从 10 周增加到 16 周
6. **Writer DTO**: 补充创建任务（原计划遗漏）
7. **风险**: 补充详细的风险登记册和回滚方案

**审查状态**: ✅ 已通过现状审查，调整为与实际项目情况一致
