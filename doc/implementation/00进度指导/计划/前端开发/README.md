# 青羽前端开发文档中心

**欢迎来到青羽写作平台前端开发文档中心！**

本目录包含前端开发的全部规划、指南和参考文档。

---

## 📚 文档导航

### 🚀 快速开始

1. **[前端快速开始指南](./前端快速开始指南_v1.0.md)** ⚡️ **必读**
   - 30分钟完成项目初始化
   - 第一天任务清单
   - 常见问题解答
   - **新手优先阅读此文档**

2. **[前端开发阶段规划](./前端开发阶段规划_v1.0.md)** 📋 **核心文档**
   - 17周完整开发规划
   - 与后端完全同步
   - 详细任务分解
   - 验收标准清晰

3. **[前端API对接清单](./前端API对接清单_v1.0.md)** 🔗 **参考手册**
   - 165+ 个API端点清单
   - 按Phase和优先级组织
   - 对接进度跟踪
   - TypeScript类型定义

---

## 🎯 文档使用指南

### 我是新加入的前端开发者

**推荐阅读顺序**:

1. 📖 阅读 [前端快速开始指南](./前端快速开始指南_v1.0.md)
   - 了解项目技术栈
   - 完成开发环境搭建
   - 运行第一个页面

2. 📋 查看 [前端开发阶段规划](./前端开发阶段规划_v1.0.md)
   - 了解整体开发规划
   - 找到自己负责的Phase
   - 熟悉验收标准

3. 🔗 参考 [前端API对接清单](./前端API对接清单_v1.0.md)
   - 查看需要对接的API
   - 了解API优先级
   - 对接时查阅

4. 📚 阅读后端API文档（在`doc/api/frontend/`目录）
   - [前端集成指南](../../../api/frontend/前端集成指南.md)
   - [API快速参考](../../../api/frontend/API快速参考.md)
   - 各模块API参考文档

---

### 我需要了解后端API

**后端API文档位置**: `doc/api/frontend/`

**核心文档**:

| 文档 | 字数 | 状态 | 说明 |
|------|------|------|------|
| [前端集成指南](../../../api/frontend/前端集成指南.md) | ~8000 | ✅ v1.3 | 完整集成教程 |
| [API快速参考](../../../api/frontend/API快速参考.md) | ~6000 | ✅ v1.3 | 一页纸速查表 |
| [认证API参考](../../../api/frontend/认证API参考.md) | ~8500 | ✅ v1.3 | JWT认证完整文档 |
| [书城API参考](../../../api/frontend/书城API参考.md) | ~12000 | ✅ v1.3 | 书城系统全部API |
| [阅读器API参考](../../../api/frontend/阅读器API参考.md) | ~12000 | ✅ v1.3 | 阅读器完整API |
| [推荐系统API参考](../../../api/frontend/推荐系统API参考.md) | ~9500 | ✅ v1.3 | 推荐算法API |
| [写作系统API参考](../../../api/frontend/写作系统API参考.md) | ~12000 | ✅ v1.2 | 项目/文档/编辑器 |
| [共享服务API参考](../../../api/frontend/共享服务API参考.md) | ~9000 | ✅ v1.3 | 钱包/存储/管理员 |

**Swagger在线文档**:
```
http://localhost:8080/swagger/index.html
```

---

### 我需要开始开发

#### Phase 0: 读者端 (Week 1-2) 🚀 **立即开始**

**任务**: 
- Day 1: 项目初始化
- Day 2-3: 用户登录注册
- Day 4-5: 书城首页与列表
- Day 6-7: 书籍详情与阅读器
- Day 8-9: 评论与推荐
- Day 10: 集成测试

**API依赖**: ✅ **后端已全部完成**
- 用户系统 (7个API)
- 书城系统 (20个API)
- 阅读器 (21个API)
- 推荐系统 (6个API)

**参考文档**:
- [前端快速开始指南](./前端快速开始指南_v1.0.md)
- [前端开发阶段规划 - Phase 0](./前端开发阶段规划_v1.0.md#phase-0-前端基础框架搭建-week-1-2)
- [前端API对接清单 - Phase 0](./前端API对接清单_v1.0.md#phase-0-读者端核心api-54个)

---

#### Phase 1: 写作端 (Week 3-4)

**任务**:
- Day 11-12: 项目管理
- Day 13-15: 文档管理与编辑器
- Day 16-17: 数据统计
- Day 18-19: 钱包系统
- Day 20: 集成测试

**API依赖**: ✅ **后端已全部完成**
- 项目管理 (6个API)
- 文档管理 (12个API)
- 编辑器 (8个API)
- 统计 (8个API)
- 钱包 (7个API)

---

#### Phase 2: 核心功能增强 (Week 5-6)

**任务**:
- Day 21-22: 文件上传与管理
- Day 23-24: 搜索优化
- Day 25: 消息通知
- Day 26-30: 编辑器增强

**API依赖**: 
- ✅ 文件存储 (15个API) - **已完成**
- 🚀 消息通知 (5个API) - **开发中，2025-11-04完成**

---

#### Phase 3: AI能力提升 (Week 7-14) 🔥 **核心功能**

**任务**:
- Week 7-8: AI辅助写作基础
- Week 9-10: RAG检索增强
- Week 11-12: AI Agent工具
- Week 13-14: 设定百科系统

**API依赖**: ⏳ **后端2025-11-30开始开发**
- AI基础服务 (10个API)
- RAG系统 (15个API)
- Agent工具 (12个API)
- 设定百科 (13个API)

**⚠️ 重要**: Phase 3是产品核心竞争力，需要：
- 充分的UI/UX设计
- 流畅的交互体验
- 详细的测试验证

---

## 📊 进度总览

| Phase | 时间 | 后端状态 | 前端状态 | 优先级 |
|-------|------|---------|---------|--------|
| **Phase 0** | Week 1-2 | ✅ 100% | ⏳ 待开始 | P0 🔥 |
| **Phase 1** | Week 3-4 | ✅ 100% | ⏳ 待开始 | P0 🔥 |
| **Phase 2** | Week 5-6 | ✅ 85% | ⏳ 待开始 | P1 |
| **Phase 3** | Week 7-14 | ⏳ 待开发 | ⏳ 待开始 | P0 🔥 |
| **Phase 4** | Week 15-17 | ⏳ 待开发 | ⏳ 待开始 | P1 |
| **Phase 5** | 持续 | 🔄 持续 | 🔄 持续 | P2 |

**当前状态**: 🚀 **可以立即启动Phase 0开发**

**关键里程碑**:
- ✅ 2025-10-28: 前端项目启动
- 🎯 2025-11-08: Phase 0完成（读者端MVP）
- 🎯 2025-11-22: Phase 1完成（写作端基础）
- 🎯 2025-12-06: Phase 2完成（核心功能增强）
- 🔥 2026-01-24: Phase 3完成（**AI核心功能**）

---

## 🛠️ 技术栈

### 核心框架
- **前端框架**: Vue 3 + TypeScript
- **构建工具**: Vite
- **路由**: Vue Router 4
- **状态管理**: Pinia
- **HTTP客户端**: Axios

### UI与样式
- **UI组件库**: Element Plus / Ant Design Vue
- **CSS框架**: Tailwind CSS
- **图表库**: ECharts
- **图谱库**: G6 / Cytoscape.js

### 开发工具
- **编辑器**: Vditor / TinyMCE
- **Markdown**: markdown-it
- **代码高亮**: highlight.js
- **日期处理**: dayjs

### 测试
- **单元测试**: Vitest
- **E2E测试**: Cypress / Playwright
- **组件测试**: @vue/test-utils

---

## 📋 开发规范

### 代码规范

**命名规范**:
- 组件文件: PascalCase (`UserProfile.vue`)
- 普通文件: camelCase (`userService.ts`)
- 常量: UPPER_SNAKE_CASE (`API_BASE_URL`)

**目录规范**:
```
src/
├── api/              # 小写，按模块划分
├── components/       # 小写，按功能划分
├── pages/            # 小写，按路由划分
├── stores/           # 小写，按模块划分
└── utils/            # 小写，工具函数
```

**组件规范**:
- 使用组合式API (Composition API)
- 使用`<script setup>`语法
- 使用TypeScript类型定义
- Props使用接口定义

**示例**:
```vue
<script setup lang="ts">
import { ref } from 'vue';

interface Props {
  title: string;
  count?: number;
}

const props = withDefaults(defineProps<Props>(), {
  count: 0
});

const emit = defineEmits<{
  (e: 'update', value: number): void;
}>();
</script>
```

---

### Git提交规范

**提交格式**:
```
<type>(<scope>): <subject>

<body>

<footer>
```

**Type类型**:
- `feat`: 新功能
- `fix`: 修复bug
- `docs`: 文档更新
- `style`: 代码格式（不影响功能）
- `refactor`: 重构
- `perf`: 性能优化
- `test`: 测试相关
- `chore`: 构建/工具相关

**示例**:
```
feat(auth): 添加用户登录功能

- 实现登录页面
- 对接登录API
- 添加Token管理

Closes #123
```

---

## 🎯 质量标准

### 代码质量
- [ ] TypeScript类型完整
- [ ] 无ESLint错误
- [ ] 代码格式统一（Prettier）
- [ ] 组件复用性高
- [ ] 注释清晰

### 功能质量
- [ ] 所有功能正常工作
- [ ] 错误处理完善
- [ ] Loading状态友好
- [ ] 空状态处理
- [ ] 边界情况考虑

### 用户体验
- [ ] 响应式适配（桌面、平板、手机）
- [ ] 页面加载速度快（首屏<3s）
- [ ] 交互流畅（无卡顿）
- [ ] 提示信息友好
- [ ] 符合无障碍规范

### 测试覆盖
- [ ] 核心功能单元测试
- [ ] 关键流程E2E测试
- [ ] 测试覆盖率>60%

---

## 🚨 常见问题

### Q1: 如何启动后端服务？

**A**: 
```bash
# 进入后端目录
cd Qingyu_backend

# 启动MongoDB和Redis（Docker）
docker-compose -f docker/docker-compose.db-only.yml up -d

# 启动后端服务
go run cmd/server/main.go
```

后端服务: `http://localhost:8080`  
Swagger文档: `http://localhost:8080/swagger/index.html`

---

### Q2: 如何解决跨域问题？

**A**: 配置Vite代理

**`vite.config.ts`**:
```typescript
export default defineConfig({
  server: {
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true
      }
    }
  }
})
```

---

### Q3: 后端API变更了怎么办？

**A**: 
1. 查看后端API文档更新日志
2. 更新TypeScript类型定义
3. 更新API调用代码
4. 运行测试验证

---

### Q4: 如何与后端团队协作？

**A**:
1. **每周同步会议**: 周一下午
2. **API变更通知**: 后端提前1天通知
3. **联调时间**: 每周五下午
4. **问题追踪**: GitHub Issues

---

## 📞 联系方式

### 技术支持
- **后端团队**: backend@qingyu.com
- **前端团队**: frontend@qingyu.com
- **GitHub Issues**: [提交问题](https://github.com/qingyu/qingyu-frontend/issues)

### 文档反馈
发现文档错误或有改进建议？
- 提交GitHub Issue
- 发送邮件到 docs@qingyu.com

---

## 📅 更新日志

### v1.0 (2025-10-27) 🎉 **初始版本**

**新增文档**:
- ✅ 前端开发阶段规划 v1.0
- ✅ 前端API对接清单 v1.0
- ✅ 前端快速开始指南 v1.0
- ✅ 前端开发文档中心 README

**核心内容**:
- 📋 17周完整开发规划
- 🔗 165+ API端点清单
- 🚀 30分钟快速启动指南
- 📚 完整后端API文档链接

**关键特性**:
- ✅ 与后端进度完全同步
- ✅ 按Phase和优先级组织
- ✅ 详细任务分解和验收标准
- ✅ TypeScript类型定义支持

---

## 🎯 下一步行动

### 立即行动

1. **阅读快速开始指南**
   - 📖 [前端快速开始指南](./前端快速开始指南_v1.0.md)

2. **搭建开发环境**
   - 安装Node.js 16+
   - 创建Vue 3 + TypeScript项目
   - 配置开发工具

3. **启动第一个任务**
   - Phase 0, Day 1: 项目初始化
   - 完成Axios封装
   - 实现登录功能

### 本周目标 (Week 1)

- [ ] 完成项目初始化
- [ ] 完成用户认证功能
- [ ] 完成书城首页

---

## 🌟 项目愿景

**我们要打造什么？**

一个功能完整、技术先进、用户体验优秀的**AI辅助写作平台**：

✨ **读者端**: 流畅的阅读体验，智能推荐，社交互动  
✨ **写作端**: 强大的编辑器，智能统计，便捷管理  
✨ **AI功能**: 智能续写、RAG增强、Agent工具、设定百科  
✨ **用户体验**: 快速响应、优雅交互、响应式设计

---

**🚀 Let's Build Something Amazing!**

**最后更新**: 2025-10-27  
**文档维护**: 前端团队负责人  
**版本**: v1.0

