"""RAG Service - RAG服务

管理RAG系统的索引和检索
"""

from typing import Any, Dict, List, Optional

from core.logger import get_logger
from rag.rag_pipeline import RAGPipeline

logger = get_logger(__name__)


class RAGService:
    """RAG服务

    特性：
    - 文档索引
    - 混合检索
    - 健康检查
    """

    def __init__(self):
        """初始化服务"""
        self.rag_pipeline: Optional[RAGPipeline] = None
        self._initialized = False
        logger.info("RAGService created")

    async def initialize(self) -> None:
        """初始化服务"""
        if self._initialized:
            return

        logger.info("Initializing RAGService...")

        # 初始化RAG Pipeline
        self.rag_pipeline = RAGPipeline()
        await self.rag_pipeline.initialize()

        self._initialized = True
        logger.info("RAGService initialized successfully")

    async def search(
        self,
        query_text: str,
        project_id: str,
        user_id: Optional[str] = None,
        content_types: Optional[List[str]] = None,
        top_k: int = 5,
    ) -> List[Dict[str, Any]]:
        """检索相关文档

        Args:
            query_text: 查询文本
            project_id: 项目ID
            user_id: 用户ID
            content_types: 内容类型过滤
            top_k: 返回结果数量

        Returns:
            检索结果列表
        """
        if not self._initialized:
            await self.initialize()

        logger.info(
            "RAG search",
            query_length=len(query_text),
            project_id=project_id,
            top_k=top_k,
            content_types=content_types,
        )

        try:
            # 调用RAG Pipeline
            results = await self.rag_pipeline.search(
                query_text=query_text,
                project_id=project_id,
                content_types=content_types,
                top_k=top_k,
            )

            logger.info(f"RAG search completed", results_count=len(results))
            return results

        except Exception as e:
            logger.error(f"RAG search failed: {e}", exc_info=True)
            raise

    async def index(
        self,
        document_id: str,
        content: str,
        metadata: Dict[str, Any],
        project_id: str,
        user_id: Optional[str] = None,
    ) -> bool:
        """索引文档

        Args:
            document_id: 文档ID
            content: 文档内容
            metadata: 元数据
            project_id: 项目ID
            user_id: 用户ID

        Returns:
            是否成功
        """
        if not self._initialized:
            await self.initialize()

        logger.info(
            "Indexing document",
            document_id=document_id,
            content_length=len(content),
            project_id=project_id,
        )

        try:
            # 调用RAG Pipeline的索引方法
            await self.rag_pipeline.index_document(
                document_id=document_id,
                content=content,
                metadata=metadata,
                project_id=project_id,
            )

            logger.info(f"Document indexed successfully", document_id=document_id)
            return True

        except Exception as e:
            logger.error(f"Document indexing failed: {e}", exc_info=True)
            return False

    async def delete(
        self,
        document_id: str,
        project_id: str,
    ) -> bool:
        """删除文档索引

        Args:
            document_id: 文档ID
            project_id: 项目ID

        Returns:
            是否成功
        """
        if not self._initialized:
            await self.initialize()

        logger.info("Deleting document index", document_id=document_id)

        try:
            # 调用RAG Pipeline的删除方法
            await self.rag_pipeline.delete_document(
                document_id=document_id,
                project_id=project_id,
            )

            logger.info(f"Document index deleted", document_id=document_id)
            return True

        except Exception as e:
            logger.error(f"Document deletion failed: {e}", exc_info=True)
            return False

    async def health_check(self) -> Dict[str, Any]:
        """健康检查

        Returns:
            健康状态
        """
        if not self._initialized:
            return {"healthy": False, "reason": "Not initialized"}

        try:
            # 检查RAG Pipeline健康状态
            pipeline_healthy = await self.rag_pipeline.health_check()

            return {
                "healthy": pipeline_healthy,
                "components": {
                    "embedding": pipeline_healthy,
                    "vector_db": pipeline_healthy,
                },
            }
        except Exception as e:
            logger.error(f"Health check failed: {e}")
            return {"healthy": False, "error": str(e)}

    async def close(self) -> None:
        """关闭服务"""
        if self.rag_pipeline:
            await self.rag_pipeline.close()
        logger.info("RAGService closed")

