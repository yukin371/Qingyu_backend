"""Agent Service - Agent服务

管理Agent工作流的执行
"""

from typing import Any, AsyncGenerator, Dict, List, Optional

from agents.workflows.creative import create_creative_workflow, execute_creative_workflow
from core.logger import get_logger
from services.tool_service import ToolService

logger = get_logger(__name__)


class AgentExecutionResult:
    """Agent执行结果"""

    def __init__(
        self,
        output: str,
        tool_calls: List[Dict[str, Any]],
        status: str,
        reasoning: List[str],
        metadata: Dict[str, Any],
    ):
        self.output = output
        self.tool_calls = tool_calls
        self.status = status
        self.reasoning = reasoning
        self.metadata = metadata


class AgentService:
    """Agent服务

    特性：
    - Workflow管理
    - 同步/流式执行
    - 健康检查
    """

    def __init__(self, tool_service: Optional[ToolService] = None):
        """初始化服务

        Args:
            tool_service: 工具服务实例（可选）
        """
        self.tool_service = tool_service or ToolService()
        self.workflows = {}
        self._initialized = False
        logger.info("AgentService created")

    async def initialize(self) -> None:
        """初始化服务"""
        if self._initialized:
            return

        logger.info("Initializing AgentService...")

        # 初始化工具服务
        await self.tool_service.initialize()

        # 创建工作流
        self.workflows["creative"] = create_creative_workflow()

        self._initialized = True
        logger.info(
            "AgentService initialized successfully",
            workflows=list(self.workflows.keys()),
        )

    async def execute(
        self,
        agent_type: str,
        task: str,
        context: Dict[str, Any],
        tools: List[str],
        user_id: Optional[str] = None,
        project_id: Optional[str] = None,
    ) -> AgentExecutionResult:
        """执行Agent（同步）

        Args:
            agent_type: Agent类型（creative, outline等）
            task: 任务描述
            context: 上下文信息
            tools: 工具列表
            user_id: 用户ID
            project_id: 项目ID

        Returns:
            执行结果
        """
        if not self._initialized:
            await self.initialize()

        logger.info(
            "Executing agent",
            agent_type=agent_type,
            task_length=len(task),
            tools=tools,
            user_id=user_id,
            project_id=project_id,
        )

        try:
            # 目前只支持creative类型
            if agent_type != "creative":
                raise ValueError(f"Unsupported agent type: {agent_type}")

            # 获取RAG工具
            rag_tool = self.tool_service.get_tool("rag_tool")

            # 执行工作流
            final_state = await execute_creative_workflow(
                task=task,
                user_id=user_id or "",
                project_id=project_id or "",
                constraints=context.get("constraints", {}),
                context=context,
                max_retries=3,
                rag_tool=rag_tool,
            )

            # 构建结果
            result = AgentExecutionResult(
                output=final_state.get("final_output", ""),
                tool_calls=final_state.get("tool_calls", []),
                status=final_state.get("current_step", "unknown"),
                reasoning=final_state.get("reasoning", []),
                metadata=final_state.get("output_metadata", {}),
            )

            logger.info(
                "Agent execution completed",
                agent_type=agent_type,
                status=result.status,
                output_length=len(result.output),
            )

            return result

        except Exception as e:
            logger.error(f"Agent execution failed: {e}", exc_info=True)
            raise

    async def execute_stream(
        self,
        agent_type: str,
        task: str,
        context: Dict[str, Any],
        tools: List[str],
    ) -> AsyncGenerator[Dict[str, Any], None]:
        """执行Agent（流式）

        Args:
            agent_type: Agent类型
            task: 任务描述
            context: 上下文信息
            tools: 工具列表

        Yields:
            事件流
        """
        if not self._initialized:
            await self.initialize()

        logger.info(
            "Executing agent stream",
            agent_type=agent_type,
            task_length=len(task),
        )

        # 获取工作流
        workflow = self.workflows.get(agent_type)
        if not workflow:
            raise ValueError(f"Unknown agent type: {agent_type}")

        # 准备初始状态
        from agents.states.creative_state import create_initial_creative_state

        initial_state = create_initial_creative_state(
            task=task,
            user_id=context.get("user_id", ""),
            project_id=context.get("project_id", ""),
            constraints=context.get("constraints", {}),
            context=context,
        )

        # 流式执行
        try:
            async for event in workflow.astream_events(initial_state, version="v1"):
                # 发送事件
                yield {
                    "event": event.get("event"),
                    "name": event.get("name"),
                    "data": event.get("data"),
                }
        except Exception as e:
            logger.error(f"Agent stream execution failed: {e}", exc_info=True)
            raise

    async def health_check(self) -> Dict[str, Any]:
        """健康检查

        Returns:
            健康状态
        """
        tool_health = await self.tool_service.health_check()

        return {
            "healthy": self._initialized and tool_health["healthy"],
            "workflows": list(self.workflows.keys()),
            "tools": tool_health,
        }

    async def close(self) -> None:
        """关闭服务"""
        await self.tool_service.close()
        logger.info("AgentService closed")

