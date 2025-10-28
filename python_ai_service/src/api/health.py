"""
健康检查 API
"""
from typing import Dict, Any
from fastapi import APIRouter, status
from datetime import datetime

from ..core import settings, get_logger

router = APIRouter(tags=["Health"])
logger = get_logger(__name__)


@router.get(
    "/health",
    status_code=status.HTTP_200_OK,
    summary="健康检查",
    description="检查服务是否正常运行"
)
async def health_check() -> Dict[str, Any]:
    """健康检查端点"""
    logger.info("health_check_called")

    return {
        "status": "healthy",
        "service": settings.service_name,
        "timestamp": datetime.utcnow().isoformat(),
        "version": "0.1.0"
    }


@router.get(
    "/health/ready",
    status_code=status.HTTP_200_OK,
    summary="就绪检查",
    description="检查服务是否就绪（所有依赖服务连接正常）"
)
async def readiness_check() -> Dict[str, Any]:
    """就绪检查端点"""
    logger.info("readiness_check_called")

    # TODO: 检查 Milvus、gRPC、Redis 连接状态
    checks = {
        "milvus": "not_implemented",
        "grpc": "not_implemented",
        "redis": "not_implemented"
    }

    all_ready = all(status == "ok" for status in checks.values() if status != "not_implemented")

    return {
        "status": "ready" if all_ready else "not_ready",
        "service": settings.service_name,
        "timestamp": datetime.utcnow().isoformat(),
        "checks": checks
    }


@router.get(
    "/health/live",
    status_code=status.HTTP_200_OK,
    summary="存活检查",
    description="检查服务进程是否存活"
)
async def liveness_check() -> Dict[str, Any]:
    """存活检查端点"""
    return {
        "status": "alive",
        "service": settings.service_name,
        "timestamp": datetime.utcnow().isoformat()
    }

