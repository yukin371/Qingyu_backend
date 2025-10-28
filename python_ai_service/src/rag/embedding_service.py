"""
向量化服务
阶段 2.1 实现
"""
from typing import List, Tuple
import torch
from sentence_transformers import SentenceTransformer

from ..core import get_logger, EmbeddingError, settings

logger = get_logger(__name__)


class EmbeddingService:
    """文本向量化服务"""

    def __init__(self):
        """初始化向量化服务"""
        self.model_name = settings.embedding_model_name
        self.device = settings.embedding_model_device
        self.batch_size = settings.embedding_batch_size
        self.model: Optional[SentenceTransformer] = None

        logger.info(
            "initializing_embedding_service",
            model=self.model_name,
            device=self.device
        )

    def load_model(self) -> None:
        """加载向量化模型"""
        # TODO: 实现模型加载逻辑
        # 1. 下载模型（如果需要）
        # 2. 加载到内存
        # 3. 移动到指定设备（CPU/GPU）
        logger.info("load_model_not_implemented")
        raise NotImplementedError("Model loading will be implemented in Stage 2.1")

    def embed_texts(self, texts: List[str]) -> List[List[float]]:
        """批量向量化文本

        Args:
            texts: 文本列表

        Returns:
            向量列表
        """
        # TODO: 实现批量向量化逻辑
        # 1. 文本预处理
        # 2. 批量编码
        # 3. 归一化
        logger.info("embed_texts_not_implemented", num_texts=len(texts))
        raise NotImplementedError("Text embedding will be implemented in Stage 2.1")

    def embed_query(self, query: str) -> List[float]:
        """向量化查询文本

        Args:
            query: 查询文本

        Returns:
            查询向量
        """
        # TODO: 实现查询向量化逻辑
        logger.info("embed_query_not_implemented")
        raise NotImplementedError("Query embedding will be implemented in Stage 2.1")

    def get_dimension(self) -> int:
        """获取向量维度

        Returns:
            向量维度
        """
        # BGE-large-zh-v1.5 的维度是 1024
        return 1024

