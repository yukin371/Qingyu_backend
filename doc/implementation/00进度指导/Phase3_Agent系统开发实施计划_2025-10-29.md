# Phase3 Agent系统开发实施计划

**开始日期**: 2025-10-29  
**预计完成**: 2025-11-15 (3周)  
**当前阶段**: 阶段3 - Agent核心功能开发

---

## 📋 执行概览

### 基础状态
- ✅ **阶段1**: 基础架构100%完成（Python微服务、gRPC、Milvus）
- ✅ **阶段2**: RAG系统95%完成（向量化、检索、索引）
- ⏳ **阶段3**: Agent系统开发（当前阶段）

### 核心目标
构建基于v2.0设计的智能协作Agent生态，包括：
1. **WorkspaceContextTool** - 上下文感知工具
2. **增强审核Agent** - 结构化诊断和反思循环
3. **元调度器** - 智能修正和迭代控制
4. **专业Agent** - 大纲、角色、情节Agent
5. **LangGraph工作流** - Agent协作编排

---

## 🎯 阶段3实施路线图

### Week 1: 核心工具和基础Agent（2025-10-29 ~ 11-04）

#### Day 1-2: WorkspaceContextTool实现
**优先级**: P0（基础工具，其他Agent依赖）

**任务**:
- [ ] 设计WorkspaceContextTool接口
- [ ] 实现任务类型识别（continue_writing, create_chapter, review_content）
- [ ] 实现结构化上下文构建
  - [ ] 前序内容提取
  - [ ] 相关角色卡加载
  - [ ] 大纲节点获取
- [ ] 集成RAG检索能力
- [ ] 单元测试

**交付物**:
```
python_ai_service/src/tools/workspace/
├── __init__.py
├── workspace_context_tool.py     # 核心工具实现
├── context_builder.py             # 上下文构建逻辑
└── task_analyzer.py               # 任务类型分析
```

#### Day 3-4: 基础Agent框架升级
**优先级**: P0

**任务**:
- [ ] 设计PipelineStateV2（支持反思循环）
- [ ] 实现BaseAgent抽象类
- [ ] 集成WorkspaceContextTool到Agent基类
- [ ] 更新现有Agent节点（generation, review）

**交付物**:
```
python_ai_service/src/agents/
├── base_agent.py                  # Agent基类
├── states/
│   └── pipeline_state_v2.py       # 新状态Schema
└── nodes/
    ├── generation_v2.py           # 升级生成节点
    └── review_v2.py               # 升级审核节点
```

#### Day 5-7: MCP工具框架
**优先级**: P1

**任务**:
- [ ] 实现MCP标准化工具接口
- [ ] 封装LangChain工具适配器
- [ ] 实现工具注册和发现机制
- [ ] 创建CharacterTool, OutlineTool
- [ ] 工具执行引擎

**交付物**:
```
python_ai_service/src/tools/mcp/
├── __init__.py
├── base_mcp_tool.py               # MCP工具基类
├── langchain_adapter.py           # LangChain适配器
├── tool_registry.py               # 工具注册表
└── tools/
    ├── character_tool.py          # 角色卡工具
    └── outline_tool.py            # 大纲工具
```

---

### Week 2: 反思循环和专业Agent（2025-11-05 ~ 11-11）

#### Day 8-10: 增强审核Agent
**优先级**: P0（v2.0核心）

**任务**:
- [ ] 设计DiagnosticReport Schema
- [ ] 实现深度诊断分析逻辑
- [ ] 生成结构化诊断报告（JSON）
- [ ] 问题根因分析
- [ ] 修正指令生成
- [ ] 受影响Agent识别

**交付物**:
```
python_ai_service/src/agents/nodes/
├── review_agent_v2.py             # 增强审核Agent
└── diagnostic/
    ├── __init__.py
    ├── report_schema.py           # 诊断报告Schema
    ├── analyzer.py                # 问题分析器
    └── correction_planner.py      # 修正计划生成
```

#### Day 11-12: 元调度器（Meta-Scheduler）
**优先级**: P0

**任务**:
- [ ] 实现meta_scheduler_node
- [ ] 诊断报告解析
- [ ] 智能定位问题Agent
- [ ] 修正Prompt生成（增强版）
- [ ] 修正策略选择（regenerate vs incremental_fix）
- [ ] 迭代次数控制和自动降级

**交付物**:
```
python_ai_service/src/agents/nodes/
├── meta_scheduler.py              # 元调度器
└── correction/
    ├── __init__.py
    ├── prompt_enhancer.py         # Prompt增强
    ├── strategy_selector.py       # 策略选择器
    └── iteration_controller.py    # 迭代控制
```

#### Day 13-14: 专业Agent实现
**优先级**: P1

**任务**:
- [ ] OutlineAgent v2（集成WorkspaceContext）
- [ ] CharacterAgent v2
- [ ] PlotAgent v2
- [ ] 为每个Agent配置专属工具集

**交付物**:
```
python_ai_service/src/agents/specialized/
├── __init__.py
├── outline_agent.py               # 大纲Agent
├── character_agent.py             # 角色Agent
└── plot_agent.py                  # 情节Agent
```

---

### Week 3: LangGraph工作流和集成测试（2025-11-12 ~ 11-15）

#### Day 15-17: LangGraph工作流
**优先级**: P0

**任务**:
- [ ] 设计v2.0工作流架构
- [ ] 实现反思循环路由
- [ ] 实现元调度器集成
- [ ] 动态Agent路由
- [ ] 错误处理和恢复

**交付物**:
```
python_ai_service/src/agents/workflows/
├── creative_v2.py                 # v2.0创作工作流
├── reflection_loop.py             # 反思循环
└── routers/
    ├── meta_router.py             # 元路由器
    └── correction_router.py       # 修正路由器
```

#### Day 18-20: 集成测试和优化
**优先级**: P0

**任务**:
- [ ] 端到端工作流测试
- [ ] 反思循环验证
- [ ] 性能优化
- [ ] 文档完善
- [ ] 示例和演示

**交付物**:
```
python_ai_service/tests/
├── test_workspace_context_tool.py
├── test_review_agent_v2.py
├── test_meta_scheduler.py
├── test_specialized_agents.py
└── test_creative_workflow_v2.py
```

---

## 📊 优先级矩阵

| 组件 | 优先级 | 依赖关系 | 预估时间 |
|------|-------|---------|---------|
| WorkspaceContextTool | P0 | RAG系统 | 2天 |
| BaseAgent框架 | P0 | WorkspaceContext | 2天 |
| MCP工具框架 | P1 | BaseAgent | 3天 |
| 增强审核Agent | P0 | BaseAgent | 3天 |
| 元调度器 | P0 | 审核Agent | 2天 |
| 专业Agent | P1 | 元调度器 | 2天 |
| LangGraph工作流 | P0 | 所有Agent | 3天 |
| 集成测试 | P0 | 工作流 | 3天 |

---

## 🎯 里程碑

### 里程碑1: 基础工具完成（Day 7）
- ✅ WorkspaceContextTool可用
- ✅ MCP工具框架搭建
- ✅ BaseAgent框架升级

### 里程碑2: 反思循环实现（Day 14）
- ✅ 增强审核Agent完成
- ✅ 元调度器实现
- ✅ 专业Agent升级

### 里程碑3: 完整系统集成（Day 20）
- ✅ LangGraph工作流完成
- ✅ 端到端测试通过
- ✅ 文档和示例完成

---

## 📝 技术栈

### 核心框架
- **LangChain**: Agent框架和工具集成
- **LangGraph**: 工作流编排和状态管理
- **Pydantic**: 数据验证和Schema定义

### AI模型
- **主模型**: GPT-4-turbo-preview（审核、规划）
- **辅助模型**: GPT-3.5-turbo（简单任务）
- **向量模型**: BGE-large-zh-v1.5

### 基础设施
- **Milvus**: 向量检索
- **Redis**: 缓存和状态管理
- **Go Backend**: 数据和业务逻辑

---

## 🔧 开发规范

### 代码结构
```
python_ai_service/src/
├── agents/                # Agent系统
│   ├── base_agent.py      # Agent基类
│   ├── specialized/       # 专业Agent
│   ├── nodes/             # Agent节点
│   ├── states/            # 状态定义
│   └── workflows/         # 工作流
├── tools/                 # 工具系统
│   ├── mcp/               # MCP工具
│   ├── workspace/         # 工作区工具
│   └── langchain/         # LangChain工具
├── rag/                   # RAG系统（已完成）
└── services/              # 服务层
```

### 命名规范
- Agent类: `XXXAgent` (如 `OutlineAgent`)
- 节点函数: `xxx_agent_node` (如 `outline_agent_node`)
- 工具类: `XXXTool` (如 `WorkspaceContextTool`)
- 状态类: `XXXState` (如 `PipelineStateV2`)

### 文档规范
- 每个模块必须有docstring
- 复杂逻辑需要注释说明
- 公开接口必须有类型注解
- 关键决策需要在代码中说明原因

---

## ✅ 验收标准

### 功能验收
- [ ] WorkspaceContextTool正确提取上下文
- [ ] 审核Agent生成结构化诊断报告
- [ ] 元调度器智能定位问题并修正
- [ ] 反思循环成功减少迭代次数
- [ ] 专业Agent输出质量符合预期
- [ ] LangGraph工作流稳定运行

### 性能验收
- [ ] 单次Agent调用 < 5秒
- [ ] 完整工作流 < 60秒
- [ ] 内存使用 < 2GB
- [ ] 并发支持 ≥ 10个请求

### 质量验收
- [ ] 单元测试覆盖率 ≥ 80%
- [ ] 集成测试全部通过
- [ ] 无严重bug
- [ ] 文档完整

---

## 📚 参考文档

### 设计文档
- [A2A流水线v2.0设计](../../design/ai/phase3/05.A2A创作流水线Agent设计_v2.0_智能协作生态.md)
- [v2.0改进总结](../../design/ai/phase3/A2A流水线v2.0改进总结.md)

### 实施文档
- [Phase3实施进度](./计划/Phase3-v2.0/实施进度_2025-10-28.md)
- [Phase3行动指南](../phase3_行动指南.md)

### 理论基础
- Reflexion: Language Agents with Verbal Reinforcement Learning
- ReAct: Synergizing Reasoning and Acting in Language Models
- Plan-and-Solve Prompting

---

## 🚀 启动指令

```bash
# 1. 确保环境就绪
cd python_ai_service
poetry install

# 2. 启动依赖服务
cd ../docker
docker-compose -f docker-compose.dev.yml up -d milvus redis

# 3. 运行开发服务器
cd ../python_ai_service
poetry run uvicorn src.main:app --reload --port 8000

# 4. 运行测试
poetry run pytest tests/ -v
```

---

**创建时间**: 2025-10-29  
**负责人**: AI Development Team  
**当前阶段**: Week 1 - Day 1

