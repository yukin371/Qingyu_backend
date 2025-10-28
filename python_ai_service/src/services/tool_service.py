"""Tool Service - 工具服务

管理Tool的注册、执行和查询
"""

from typing import Any, Dict, List, Optional

from core.logger import get_logger
from core.tools.base import BaseTool, ToolMetadata, ToolResult
from core.tools.langchain import CharacterTool, OutlineTool, RAGTool
from core.tools.registry import ToolRegistry
from infrastructure.go_api import GoAPIClient

logger = get_logger(__name__)


class ToolService:
    """工具服务

    特性：
    - Tool注册和管理
    - Tool执行（带权限检查）
    - Tool列表查询
    """

    def __init__(self):
        """初始化服务"""
        self.registry = ToolRegistry()
        self.go_api_client: Optional[GoAPIClient] = None
        self._initialized = False
        logger.info("ToolService created")

    async def initialize(self) -> None:
        """初始化服务"""
        if self._initialized:
            return

        logger.info("Initializing ToolService...")

        # 初始化Go API客户端
        self.go_api_client = GoAPIClient()
        await self.go_api_client.initialize()

        # 注册内置工具
        await self._register_builtin_tools()

        self._initialized = True
        logger.info("ToolService initialized successfully")

    async def _register_builtin_tools(self) -> None:
        """注册内置工具"""
        logger.info("Registering builtin tools...")

        # 注册RAGTool
        rag_tool = RAGTool()
        self.registry.register_tool_instance(rag_tool)
        logger.info("Registered RAGTool")

        # 注册CharacterTool
        character_tool = CharacterTool(go_api_client=self.go_api_client)
        self.registry.register_tool_instance(character_tool)
        logger.info("Registered CharacterTool")

        # 注册OutlineTool
        outline_tool = OutlineTool(go_api_client=self.go_api_client)
        self.registry.register_tool_instance(outline_tool)
        logger.info("Registered OutlineTool")

        logger.info(f"Registered {len(self.registry.list_tools())} builtin tools")

    async def execute_tool(
        self,
        tool_name: str,
        params: Dict[str, Any],
        user_id: Optional[str] = None,
        project_id: Optional[str] = None,
        agent_call_id: Optional[str] = None,
    ) -> ToolResult:
        """执行工具

        Args:
            tool_name: 工具名称
            params: 工具参数
            user_id: 用户ID
            project_id: 项目ID
            agent_call_id: Agent调用ID

        Returns:
            工具执行结果
        """
        if not self._initialized:
            await self.initialize()

        # 获取工具
        tool = self.registry.get_tool(tool_name)
        if not tool:
            logger.error(f"Tool not found: {tool_name}")
            return ToolResult(success=False, error=f"Tool not found: {tool_name}")

        # 执行工具
        logger.info(f"Executing tool: {tool_name}", params=params)

        result = await tool.execute(
            params=params,
            user_id=user_id,
            project_id=project_id,
            agent_call_id=agent_call_id,
        )

        logger.info(
            f"Tool execution completed: {tool_name}",
            success=result.success,
            duration_ms=result.duration_ms,
        )

        return result

    def list_tools(self) -> List[ToolMetadata]:
        """列出所有工具

        Returns:
            工具元数据列表
        """
        return self.registry.list_tools()

    def get_tool(self, tool_name: str) -> Optional[BaseTool]:
        """获取工具实例

        Args:
            tool_name: 工具名称

        Returns:
            工具实例（如果存在）
        """
        return self.registry.get_tool(tool_name)

    def get_tools_by_category(self, category: str) -> List[BaseTool]:
        """按分类获取工具

        Args:
            category: 分类名称

        Returns:
            工具列表
        """
        return self.registry.get_tools_by_category(category)

    async def health_check(self) -> Dict[str, Any]:
        """健康检查

        Returns:
            健康状态
        """
        return {
            "healthy": self._initialized,
            "tools_count": len(self.registry.list_tools()),
            "tools": [meta.name for meta in self.registry.list_tools()],
        }

    async def close(self) -> None:
        """关闭服务"""
        if self.go_api_client:
            await self.go_api_client.close()
        logger.info("ToolService closed")

