"""
配置管理模块
使用 Pydantic Settings 从环境变量加载配置
"""
from typing import Optional
from pydantic import Field
from pydantic_settings import BaseSettings, SettingsConfigDict


class Settings(BaseSettings):
    """应用配置"""

    model_config = SettingsConfigDict(
        env_file=".env",
        env_file_encoding="utf-8",
        case_sensitive=False,
        extra="ignore"
    )

    # Service
    service_name: str = Field(default="qingyu-ai-service", alias="SERVICE_NAME")
    service_port: int = Field(default=8000, alias="SERVICE_PORT")
    log_level: str = Field(default="INFO", alias="LOG_LEVEL")

    # OpenAI
    openai_api_key: str = Field(default="", alias="OPENAI_API_KEY")
    openai_base_url: str = Field(
        default="https://api.openai.com/v1",
        alias="OPENAI_BASE_URL"
    )
    openai_model: str = Field(default="gpt-4-turbo-preview", alias="OPENAI_MODEL")

    # Anthropic
    anthropic_api_key: str = Field(default="", alias="ANTHROPIC_API_KEY")
    anthropic_model: str = Field(
        default="claude-3-sonnet-20240229",
        alias="ANTHROPIC_MODEL"
    )

    # Milvus
    milvus_host: str = Field(default="localhost", alias="MILVUS_HOST")
    milvus_port: int = Field(default=19530, alias="MILVUS_PORT")
    milvus_collection_name: str = Field(
        default="qingyu_knowledge",
        alias="MILVUS_COLLECTION_NAME"
    )

    # Go Backend gRPC
    go_grpc_host: str = Field(default="localhost", alias="GO_GRPC_HOST")
    go_grpc_port: int = Field(default=50051, alias="GO_GRPC_PORT")

    # Redis
    redis_host: str = Field(default="localhost", alias="REDIS_HOST")
    redis_port: int = Field(default=6379, alias="REDIS_PORT")
    redis_db: int = Field(default=0, alias="REDIS_DB")
    redis_password: str = Field(default="", alias="REDIS_PASSWORD")

    # Embedding
    embedding_provider: str = Field(default="local", alias="EMBEDDING_PROVIDER")  # local, openai, custom
    embedding_model_name: str = Field(
        default="BAAI/bge-large-zh-v1.5",
        alias="EMBEDDING_MODEL_NAME"
    )
    embedding_model_device: str = Field(default="cuda", alias="EMBEDDING_MODEL_DEVICE")
    embedding_batch_size: int = Field(default=32, alias="EMBEDDING_BATCH_SIZE")
    embedding_cache_enabled: bool = Field(default=True, alias="EMBEDDING_CACHE_ENABLED")
    embedding_cache_ttl: int = Field(default=604800, alias="EMBEDDING_CACHE_TTL")  # 7天

    # OpenAI Embedding
    openai_embedding_model: str = Field(
        default="text-embedding-3-small",
        alias="OPENAI_EMBEDDING_MODEL"
    )
    openai_embedding_batch_size: int = Field(default=100, alias="OPENAI_EMBEDDING_BATCH_SIZE")
    openai_embedding_max_retries: int = Field(default=3, alias="OPENAI_EMBEDDING_MAX_RETRIES")

    # Text Splitter
    text_chunk_size: int = Field(default=500, alias="TEXT_CHUNK_SIZE")
    text_chunk_overlap: int = Field(default=50, alias="TEXT_CHUNK_OVERLAP")
    text_splitter_type: str = Field(default="recursive", alias="TEXT_SPLITTER_TYPE")  # recursive, semantic

    # LangSmith
    langchain_tracing_v2: bool = Field(default=False, alias="LANGCHAIN_TRACING_V2")
    langchain_endpoint: Optional[str] = Field(default=None, alias="LANGCHAIN_ENDPOINT")
    langchain_api_key: Optional[str] = Field(default=None, alias="LANGCHAIN_API_KEY")
    langchain_project: str = Field(
        default="qingyu-ai-service",
        alias="LANGCHAIN_PROJECT"
    )

    @property
    def go_grpc_address(self) -> str:
        """Go gRPC 服务地址"""
        return f"{self.go_grpc_host}:{self.go_grpc_port}"

    @property
    def milvus_uri(self) -> str:
        """Milvus 连接 URI"""
        return f"http://{self.milvus_host}:{self.milvus_port}"


# 全局配置实例
settings = Settings()

