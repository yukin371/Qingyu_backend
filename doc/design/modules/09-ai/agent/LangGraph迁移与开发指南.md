# LangGraph迁移与开发指南

> **文档版本**: v1.0  
> **创建时间**: 2025-10-21  
> **适用对象**: 青羽平台AI开发团队

## 📋 文档概述

本文档为青羽平台AI模块从LangChain迁移到LangGraph提供完整指南，包括：
- 为什么选择LangGraph
- 核心概念对比
- 迁移步骤
- 最佳实践
- 常见问题

**阅读前提**：
- 已阅读 [Agent框架技术选型对比](./Agent框架技术选型对比_LangChain_vs_Others.md)
- 了解基础的LangChain概念（可选）
- 熟悉Python异步编程

---

## 🎯 为什么迁移到LangGraph

### 核心理由

| 维度 | LangChain | LangGraph | 青羽项目需求 |
|------|-----------|-----------|------------|
| **工作流类型** | 链式（Sequential） | 图状（Graph-based） | ✅ 需要图状（审核循环、条件分支） |
| **流程控制** | if-else手动控制 | 声明式边和条件路由 | ✅ 需要复杂流程编排 |
| **状态管理** | 隐式（Memory） | 显式（TypedDict State） | ✅ 需要可调试、可持久化 |
| **工具调用** | 手动调用 | ToolNode自动解析 | ✅ 简化工具调用 |
| **生态兼容** | - | 100%继承LangChain | ✅ 保留LangChain工具和RAG |

### 实际场景对比

**场景**：创作Agent需要执行"理解任务 → RAG检索 → 生成内容 → 审核 → 不通过则重新生成"

#### LangChain实现（手动控制）

```python
# ❌ 复杂且难以维护
from langchain.chains import LLMChain

class CreativeAgent:
    def run(self, task):
        # 1. 理解任务
        understanding = self.understand_chain.run(task)
        
        # 2. RAG检索
        rag_results = self.rag_retriever.search(understanding)
        
        # 3. 生成内容
        retry_count = 0
        while retry_count < 3:
            content = self.generation_chain.run(rag_results)
            
            # 4. 审核
            review = self.review_chain.run(content)
            
            if review['passed']:
                return content  # ✅ 审核通过
            else:
                retry_count += 1  # ❌ 重新生成
        
        return content  # 超过最大重试次数
```

**问题**：
- ❌ 手动while循环，难以可视化
- ❌ 状态隐式，调试困难
- ❌ 无法持久化中间状态
- ❌ 代码结构混乱

#### LangGraph实现（声明式）

```python
# ✅ 清晰、可维护、可视化
from langgraph.graph import StateGraph, END
from typing import TypedDict

class AgentState(TypedDict):
    task: str
    understanding: str
    rag_results: list
    content: str
    review: dict
    retry_count: int

# 定义节点
def understand_node(state): ...
def rag_node(state): ...
def generate_node(state): ...
def review_node(state): ...

# 定义条件路由
def should_regenerate(state):
    if state['review']['passed']:
        return 'end'
    elif state['retry_count'] < 3:
        return 'regenerate'
    else:
        return 'end'

# 构建图
workflow = StateGraph(AgentState)
workflow.add_node("understand", understand_node)
workflow.add_node("rag", rag_node)
workflow.add_node("generate", generate_node)
workflow.add_node("review", review_node)

workflow.set_entry_point("understand")
workflow.add_edge("understand", "rag")
workflow.add_edge("rag", "generate")
workflow.add_edge("generate", "review")
workflow.add_conditional_edges(
    "review",
    should_regenerate,
    {
        "regenerate": "generate",  # 循环
        "end": END
    }
)

app = workflow.compile()
```

**优势**：
- ✅ 工作流一目了然（可视化）
- ✅ 显式状态，易于调试
- ✅ 支持持久化（Checkpointer）
- ✅ 代码结构清晰

---

## 📖 核心概念对比

### 1. 状态管理

#### LangChain：隐式Memory

```python
# LangChain方式
from langchain.memory import ConversationBufferMemory

memory = ConversationBufferMemory()
chain = LLMChain(llm=llm, memory=memory)

# ❌ 问题：
# - 状态分散在Memory对象中
# - 不支持类型检查
# - 持久化困难
```

#### LangGraph：显式State

```python
# LangGraph方式
from typing import TypedDict, Annotated
import operator

class AgentState(TypedDict):
    messages: Annotated[list, operator.add]  # 自动合并
    user_id: str
    context: dict
    retry_count: int

# ✅ 优势：
# - 类型安全（TypedDict）
# - 状态集中管理
# - 支持Reducer（如operator.add）
# - 易于持久化
```

### 2. 工作流定义

#### LangChain：链式组合

```python
# LangChain方式
from langchain.chains import SequentialChain

chain = SequentialChain(chains=[
    chain1,  # 理解任务
    chain2,  # RAG检索
    chain3,  # 生成内容
])

# ❌ 问题：
# - 只能顺序执行，无条件分支
# - 循环需要手动while
# - 不支持并行
```

#### LangGraph：图状工作流

```python
# LangGraph方式
workflow = StateGraph(AgentState)
workflow.add_node("understand", node1)
workflow.add_node("rag", node2)
workflow.add_node("generate", node3)
workflow.add_node("review", node4)

workflow.add_edge("understand", "rag")
workflow.add_conditional_edges(
    "review",
    condition_func,
    {
        "pass": END,
        "retry": "generate"  # 循环
    }
)

# ✅ 优势：
# - 支持条件分支
# - 声明式循环
# - 支持并行节点
# - 可视化工作流
```

### 3. 工具调用

#### LangChain：手动调用

```python
# LangChain方式
from langchain.tools import Tool

tool = Tool(
    name="character_create",
    func=create_character,
    description="创建角色卡"
)

# 手动调用
result = tool.run({"name": "林风", "personality": "..."})

# ❌ 问题：
# - 需要手动解析LLM的工具调用请求
# - 错误处理复杂
```

#### LangGraph：ToolNode自动处理

```python
# LangGraph方式
from langgraph.prebuilt import ToolNode

tools = [character_create_tool, rag_retrieval_tool]
tool_node = ToolNode(tools)

# 添加到图中
workflow.add_node("tools", tool_node)

# ✅ 优势：
# - 自动解析LLM的function_call
# - 自动执行工具
# - 自动错误处理
# - 支持并行工具调用
```

### 4. 持久化

#### LangChain：需要自己实现

```python
# LangChain方式
# ❌ 需要手动保存Memory和中间结果
import json

state = {
    "memory": memory.to_dict(),
    "context": context
}
with open("state.json", "w") as f:
    json.dump(state, f)
```

#### LangGraph：内置Checkpointer

```python
# LangGraph方式
from langgraph.checkpoint.postgres import PostgresSaver

checkpointer = PostgresSaver.from_conn_string(
    "postgresql://user:pass@localhost/db"
)

app = workflow.compile(checkpointer=checkpointer)

# ✅ 优势：
# - 自动保存每个节点的状态
# - 支持断点恢复
# - 支持人工介入（human-in-the-loop）
# - 多线程/多会话管理
```

---

## 🚀 迁移步骤

### 阶段1：环境准备（1天）

#### 1.1 安装依赖

```bash
# 安装LangGraph（会自动安装LangChain）
pip install langgraph langchain-openai langchain-community

# 可选：持久化支持
pip install langgraph-checkpoint-postgres

# 可选：可视化支持
pip install langgraph-studio
```

#### 1.2 项目结构调整

```
python_ai_service/
├── agents/
│   ├── __init__.py
│   ├── base_agent.py          # BaseAgent抽象类
│   ├── creative_agent.py      # LangGraph实现
│   ├── analysis_agent.py
│   └── nodes/                  # 节点定义
│       ├── __init__.py
│       ├── understand.py
│       ├── rag.py
│       ├── generate.py
│       └── review.py
├── tools/                      # LangChain工具
│   ├── character_tool.py
│   ├── outline_tool.py
│   └── rag_tool.py
├── graphs/                     # 工作流定义
│   └── creative_workflow.py
└── checkpointers/              # 持久化
    └── postgres_saver.py
```

### 阶段2：定义状态（1天）

#### 2.1 设计State Schema

```python
# agents/states.py
from typing import TypedDict, Annotated, Sequence
from langchain_core.messages import BaseMessage
import operator

class CreativeAgentState(TypedDict):
    """创作Agent状态"""
    # 输入
    task_description: str
    user_id: str
    project_id: str
    
    # 消息历史（自动累积）
    messages: Annotated[Sequence[BaseMessage], operator.add]
    
    # 工作流状态
    understanding: dict
    plan: list[dict]
    current_step: int
    
    # RAG结果
    rag_results: list[str]
    
    # 生成内容
    generated_content: str
    
    # 审核结果
    review_result: dict
    review_passed: bool
    retry_count: int
    
    # 工具调用记录
    tool_calls: list[dict]
    
    # 最终输出
    final_output: str
    reasoning: list[str]
```

#### 2.2 定义Reducer（可选）

```python
# 自定义Reducer示例
def merge_reasoning(current: list[str], new: list[str]) -> list[str]:
    """合并推理过程，限制最大长度"""
    merged = current + new
    return merged[-50:]  # 只保留最近50条

class CreativeAgentState(TypedDict):
    reasoning: Annotated[list[str], merge_reasoning]
```

### 阶段3：迁移节点逻辑（2-3天）

#### 3.1 将LangChain链转换为节点函数

**迁移前（LangChain）**：

```python
# 旧代码
from langchain.chains import LLMChain

understand_chain = LLMChain(
    llm=llm,
    prompt=understand_prompt
)

result = understand_chain.run(task)
```

**迁移后（LangGraph）**：

```python
# agents/nodes/understand.py
from langchain_openai import ChatOpenAI
from langchain_core.messages import HumanMessage

def understand_node(state: CreativeAgentState) -> CreativeAgentState:
    """理解任务节点"""
    llm = ChatOpenAI(model="gpt-4", temperature=0)
    
    prompt = f"""
    分析以下创作任务：
    {state['task_description']}
    
    请提取：
    1. 任务类型
    2. 关键要素
    3. 所需工具
    """
    
    response = llm.invoke([HumanMessage(content=prompt)])
    
    # 更新状态
    return {
        **state,
        'messages': [HumanMessage(content=prompt), response],
        'understanding': {
            'task_type': '...',
            'key_elements': [...],
            'required_tools': [...]
        },
        'reasoning': state['reasoning'] + [f"任务理解完成"]
    }
```

#### 3.2 迁移工具调用

**迁移前（LangChain手动调用）**：

```python
# 旧代码
from langchain.tools import Tool

character_tool = Tool(
    name="character_create",
    func=create_character,
    description="创建角色卡"
)

# 手动调用
result = character_tool.run(params)
```

**迁移后（LangGraph ToolNode）**：

```python
# tools/character_tool.py
from langchain.tools import BaseTool

class CharacterCreateTool(BaseTool):
    name = "character_create"
    description = "创建小说角色卡"
    
    def _run(self, name: str, personality: str, **kwargs) -> str:
        # 调用Go API
        import requests
        response = requests.post(
            f"{GO_API_BASE}/api/v1/projects/{kwargs['project_id']}/characters",
            json={"name": name, "personality": personality},
            headers={"Authorization": f"Bearer {kwargs['token']}"}
        )
        return response.json()
    
    async def _arun(self, *args, **kwargs):
        return self._run(*args, **kwargs)

# graphs/creative_workflow.py
from langgraph.prebuilt import ToolNode

tools = [CharacterCreateTool(), OutlineTool(), RAGTool()]
tool_node = ToolNode(tools)

# 添加到工作流
workflow.add_node("tools", tool_node)
```

### 阶段4：构建工作流（2天）

#### 4.1 定义节点

```python
# graphs/creative_workflow.py
from langgraph.graph import StateGraph, END
from agents.states import CreativeAgentState
from agents.nodes import (
    understand_node,
    plan_node,
    execute_node,
    review_node,
    regenerate_node,
    finalize_node
)

workflow = StateGraph(CreativeAgentState)

# 添加节点
workflow.add_node("understand", understand_node)
workflow.add_node("plan", plan_node)
workflow.add_node("execute", execute_node)
workflow.add_node("tools", tool_node)
workflow.add_node("review", review_node)
workflow.add_node("regenerate", regenerate_node)
workflow.add_node("finalize", finalize_node)
```

#### 4.2 定义边和条件路由

```python
# 定义条件函数
def should_continue(state: CreativeAgentState) -> str:
    """决定是否继续执行步骤"""
    if state['current_step'] < len(state['plan']) - 1:
        return "continue"
    else:
        return "review"

def should_regenerate(state: CreativeAgentState) -> str:
    """决定是否重新生成"""
    if state['review_passed']:
        return "finalize"
    elif state['retry_count'] < 3:
        return "regenerate"
    else:
        return "finalize"  # 强制结束

# 定义边
workflow.set_entry_point("understand")
workflow.add_edge("understand", "plan")
workflow.add_edge("plan", "execute")
workflow.add_edge("execute", "tools")

# 条件边
workflow.add_conditional_edges(
    "tools",
    should_continue,
    {
        "continue": "execute",
        "review": "review"
    }
)

workflow.add_conditional_edges(
    "review",
    should_regenerate,
    {
        "regenerate": "regenerate",
        "finalize": "finalize"
    }
)

workflow.add_edge("regenerate", "review")  # 循环
workflow.add_edge("finalize", END)
```

#### 4.3 编译和持久化

```python
# 添加持久化（可选）
from langgraph.checkpoint.postgres import PostgresSaver

checkpointer = PostgresSaver.from_conn_string(
    "postgresql://user:pass@localhost/qingyu_db"
)

# 编译
app = workflow.compile(checkpointer=checkpointer)

# 导出可视化
app.get_graph().draw_mermaid_png(output_file_path="workflow.png")
```

### 阶段5：集成FastAPI（1天）

#### 5.1 创建API端点

```python
# main.py
from fastapi import FastAPI, HTTPException
from pydantic import BaseModel

app = FastAPI()

class CreativeRequest(BaseModel):
    task_description: str
    user_id: str
    project_id: str
    token: str

@app.post("/api/ai/creative/generate")
async def generate_creative_content(request: CreativeRequest):
    """创作内容生成（流式）"""
    
    # 初始化状态
    initial_state = {
        "task_description": request.task_description,
        "user_id": request.user_id,
        "project_id": request.project_id,
        "messages": [],
        "plan": [],
        "current_step": 0,
        "generated_content": "",
        "tool_calls": [],
        "review_result": {},
        "review_passed": False,
        "retry_count": 0,
        "final_output": "",
        "reasoning": []
    }
    
    # 执行工作流（流式）
    async for event in app.astream_events(
        initial_state,
        config={"configurable": {"thread_id": f"{request.user_id}_{request.project_id}"}}
    ):
        if event['event'] == 'on_chat_model_stream':
            # 流式输出LLM生成的内容
            chunk = event['data']['chunk']
            yield f"data: {chunk}\n\n"
        
        elif event['event'] == 'on_tool_start':
            # 工具调用开始
            yield f"data: {{\"type\": \"tool_start\", \"tool\": \"{event['name']}\"}}\n\n"
        
        elif event['event'] == 'on_tool_end':
            # 工具调用结束
            yield f"data: {{\"type\": \"tool_end\", \"result\": {event['data']}}}\n\n"
    
    # 最终结果
    final_result = await app.ainvoke(initial_state, config=...)
    yield f"data: {{\"type\": \"final\", \"output\": \"{final_result['final_output']}\"}}\n\n"
```

---

## 💡 最佳实践

### 1. 节点设计原则

#### ✅ 单一职责

```python
# ✅ 推荐：每个节点只做一件事
def understand_node(state):
    """只负责理解任务"""
    ...

def rag_node(state):
    """只负责RAG检索"""
    ...

# ❌ 不推荐：一个节点做多件事
def understand_and_rag_node(state):
    """又理解又检索，职责不清"""
    ...
```

#### ✅ 状态不可变性

```python
# ✅ 推荐：返回新状态，不修改原状态
def my_node(state: AgentState) -> AgentState:
    return {
        **state,
        'new_field': 'value'
    }

# ❌ 不推荐：直接修改状态
def my_node(state: AgentState) -> AgentState:
    state['new_field'] = 'value'  # ❌
    return state
```

### 2. 条件路由最佳实践

```python
# ✅ 推荐：清晰的条件逻辑
def route_after_review(state: AgentState) -> str:
    """审核后的路由决策"""
    if state['review_passed']:
        return 'success'
    
    if state['retry_count'] >= 3:
        return 'max_retry_reached'
    
    if state['review_result']['severity'] == 'critical':
        return 'escalate_to_human'
    
    return 'retry'

# 使用时
workflow.add_conditional_edges(
    "review",
    route_after_review,
    {
        "success": "finalize",
        "retry": "regenerate",
        "max_retry_reached": "finalize",
        "escalate_to_human": "human_review"
    }
)
```

### 3. 错误处理

```python
# agents/nodes/generate.py
from langchain_core.messages import HumanMessage

def generate_node(state: AgentState) -> AgentState:
    """生成内容节点"""
    try:
        llm = ChatOpenAI(model="gpt-4")
        response = llm.invoke([HumanMessage(content=state['prompt'])])
        
        return {
            **state,
            'generated_content': response.content,
            'error': None
        }
    
    except Exception as e:
        # 记录错误，继续工作流（降级处理）
        return {
            **state,
            'generated_content': '',
            'error': str(e),
            'reasoning': state['reasoning'] + [f"生成失败：{e}"]
        }
```

### 4. 持久化和恢复

```python
# 保存会话状态
result = await app.ainvoke(
    initial_state,
    config={
        "configurable": {
            "thread_id": "user123_session001"  # 唯一会话ID
        }
    }
)

# 恢复会话（从断点继续）
continued_result = await app.ainvoke(
    None,  # 不需要初始状态，会从checkpointer加载
    config={
        "configurable": {
            "thread_id": "user123_session001"
        }
    }
)
```

### 5. 人工介入（Human-in-the-Loop）

```python
from langgraph.graph import interrupt

def review_node(state: AgentState) -> AgentState:
    """审核节点"""
    review_result = auto_review(state['generated_content'])
    
    if review_result['needs_human']:
        # 中断工作流，等待人工审核
        return interrupt({
            "message": "需要人工审核",
            "content": state['generated_content'],
            "review": review_result
        })
    
    return {
        **state,
        'review_passed': review_result['passed']
    }

# 人工审核后恢复
resumed_result = await app.ainvoke(
    {"human_decision": "approved"},  # 人工决策
    config={"configurable": {"thread_id": "session001"}}
)
```

---

## 🔍 调试和可视化

### 1. 可视化工作流

```python
# 生成Mermaid图
from IPython.display import Image, display

mermaid_png = app.get_graph().draw_mermaid_png()
display(Image(mermaid_png))

# 或保存到文件
with open("workflow.png", "wb") as f:
    f.write(mermaid_png)
```

### 2. 调试执行过程

```python
# 打印每个节点的执行
async for event in app.astream_events(initial_state):
    print(f"Event: {event['event']}")
    print(f"Name: {event['name']}")
    print(f"Data: {event['data']}")
    print("---")
```

### 3. 查看状态历史

```python
# 获取所有checkpoint
from langgraph.checkpoint.postgres import PostgresSaver

checkpointer = PostgresSaver.from_conn_string("...")
history = checkpointer.list(
    config={"configurable": {"thread_id": "session001"}}
)

for checkpoint in history:
    print(f"Step: {checkpoint.step}")
    print(f"State: {checkpoint.state}")
```

---

## ❓ 常见问题

### Q1: LangGraph是否完全兼容LangChain工具？

**A**: 是的！LangGraph 100%兼容LangChain工具。

```python
# 任何LangChain工具都可以直接使用
from langchain_community.tools.tavily_search import TavilySearchResults
from langchain_community.retrievers import WikipediaRetriever
from langgraph.prebuilt import ToolNode

tools = [
    TavilySearchResults(),
    WikipediaRetriever(),
    # ... 你的自定义工具
]

tool_node = ToolNode(tools)
```

### Q2: 如何从LangChain的Memory迁移到LangGraph的State？

**A**: 使用Annotated + operator.add

```python
# LangChain Memory
from langchain.memory import ConversationBufferMemory
memory = ConversationBufferMemory()

# LangGraph State（等价）
from typing import Annotated, Sequence
import operator
from langchain_core.messages import BaseMessage

class State(TypedDict):
    messages: Annotated[Sequence[BaseMessage], operator.add]
```

### Q3: LangGraph性能如何？

**A**: 
- **单次调用**：与LangChain基本一致
- **复杂工作流**：更优（避免重复计算）
- **持久化开销**：增加约10-20ms（可选功能）

建议：
- 开发环境：启用持久化（便于调试）
- 生产环境：根据需求决定

### Q4: 如何处理长时间运行的任务？

**A**: 使用Checkpointer + 异步任务

```python
# 启动长任务
task_id = "task_12345"
asyncio.create_task(
    app.ainvoke(
        initial_state,
        config={"configurable": {"thread_id": task_id}}
    )
)

# 定期查询状态
checkpointer = PostgresSaver.from_conn_string("...")
current_state = checkpointer.get(
    config={"configurable": {"thread_id": task_id}}
)
```

### Q5: LangGraph是否支持并行执行？

**A**: 是的！

```python
# 并行节点（会自动并行执行）
workflow.add_node("rag", rag_node)
workflow.add_node("outline", outline_node)

# 从同一个节点指向多个节点 = 并行
workflow.add_edge("understand", "rag")
workflow.add_edge("understand", "outline")

# 汇聚点
workflow.add_edge("rag", "generate")
workflow.add_edge("outline", "generate")
```

---

## 📚 参考资源

### 官方文档
- [LangGraph官方文档](https://langchain-ai.github.io/langgraph/)
- [LangGraph GitHub](https://github.com/langchain-ai/langgraph)
- [LangChain官方文档](https://python.langchain.com/docs/get_started/introduction)

### 青羽项目文档
- [Agent框架技术选型对比](./Agent框架技术选型对比_LangChain_vs_Others.md) - 为什么选择LangGraph
- [Python AI Agent系统架构设计](./07.Python_AI_Agent系统架构设计.md) - 完整实现示例
- [Agent工具调用集成设计](./09.Agent工具调用集成设计.md) - 工具调用详细设计

### 示例代码
- [LangGraph示例库](https://github.com/langchain-ai/langgraph/tree/main/examples)
- [青羽Agent实现示例](./07.Python_AI_Agent系统架构设计.md#附录langgraph实现示例-)

---

## 📝 迁移检查清单

### 环境准备
- [ ] 安装LangGraph和相关依赖
- [ ] 设置PostgreSQL（如需持久化）
- [ ] 调整项目结构

### 代码迁移
- [ ] 定义State Schema
- [ ] 迁移LangChain链为节点函数
- [ ] 迁移工具为LangChain BaseTool
- [ ] 构建StateGraph工作流
- [ ] 定义条件路由
- [ ] 添加错误处理

### 集成测试
- [ ] 单节点测试
- [ ] 工作流端到端测试
- [ ] 持久化功能测试
- [ ] 流式输出测试
- [ ] 性能测试

### 部署上线
- [ ] 更新Docker镜像
- [ ] 配置数据库连接
- [ ] 配置环境变量
- [ ] 监控和日志
- [ ] 文档更新

---

**文档版本**: v1.0  
**创建时间**: 2025-10-21  
**维护者**: 青羽AI架构组  
**状态**: ✅ 可用

