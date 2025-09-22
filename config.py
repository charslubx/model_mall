"""
网关系统配置文件
"""
import os
from typing import Dict, Any
from dataclasses import dataclass


@dataclass
class GatewayConfig:
    """网关配置类"""
    
    # 服务器配置
    host: str = "0.0.0.0"
    port: int = 8000
    debug: bool = True
    
    # 静态文件配置
    static_dir: str = "/workspace/static"
    
    # 日志配置
    log_level: str = "INFO"
    log_format: str = "%(asctime)s - %(name)s - %(levelname)s - %(message)s"
    
    # CORS配置
    cors_origins: list = None
    cors_allow_credentials: bool = True
    cors_allow_methods: list = None
    cors_allow_headers: list = None
    
    # 超时配置
    request_timeout: int = 30
    connection_timeout: int = 10
    
    # 处理器配置
    max_file_size: int = 10 * 1024 * 1024  # 10MB
    
    def __post_init__(self):
        """初始化后处理"""
        if self.cors_origins is None:
            self.cors_origins = ["*"]
        if self.cors_allow_methods is None:
            self.cors_allow_methods = ["*"]
        if self.cors_allow_headers is None:
            self.cors_allow_headers = ["*"]


def get_config() -> GatewayConfig:
    """获取配置实例"""
    return GatewayConfig(
        host=os.getenv("GATEWAY_HOST", "0.0.0.0"),
        port=int(os.getenv("GATEWAY_PORT", "8000")),
        debug=os.getenv("GATEWAY_DEBUG", "true").lower() == "true",
        static_dir=os.getenv("GATEWAY_STATIC_DIR", "/workspace/static"),
        log_level=os.getenv("GATEWAY_LOG_LEVEL", "INFO"),
        request_timeout=int(os.getenv("GATEWAY_REQUEST_TIMEOUT", "30")),
        connection_timeout=int(os.getenv("GATEWAY_CONNECTION_TIMEOUT", "10")),
        max_file_size=int(os.getenv("GATEWAY_MAX_FILE_SIZE", str(10 * 1024 * 1024)))
    )