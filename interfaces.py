"""
网关系统的抽象接口定义
"""
from abc import ABC, abstractmethod
from typing import Any, Dict, Optional, Union
from fastapi import Request, Response
from enum import Enum


class RequestType(Enum):
    """请求类型枚举"""
    XHR = "xhr"
    NON_XHR = "non_xhr"
    HTTP = "http"
    NON_HTTP = "non_http"
    REST_API = "rest_api"
    FILE = "file"
    SSE = "sse"
    WEBSOCKET = "websocket"


class GatewayInterface(ABC):
    """网关顶层接口 - 统一处理所有请求"""
    
    @abstractmethod
    async def process_request(self, request: Request) -> Response:
        """
        处理传入的请求
        
        Args:
            request: FastAPI请求对象
            
        Returns:
            Response: 处理后的响应
        """
        pass
    
    @abstractmethod
    def determine_request_type(self, request: Request) -> RequestType:
        """
        确定请求类型
        
        Args:
            request: FastAPI请求对象
            
        Returns:
            RequestType: 请求类型枚举
        """
        pass


class RequestHandlerInterface(ABC):
    """请求处理器接口"""
    
    @abstractmethod
    async def handle(self, request: Request) -> Response:
        """
        处理特定类型的请求
        
        Args:
            request: FastAPI请求对象
            
        Returns:
            Response: 处理后的响应
        """
        pass
    
    @abstractmethod
    def can_handle(self, request: Request) -> bool:
        """
        判断是否能处理该请求
        
        Args:
            request: FastAPI请求对象
            
        Returns:
            bool: 是否能处理
        """
        pass


class XHRHandlerInterface(RequestHandlerInterface):
    """XHR处理器接口"""
    
    @abstractmethod
    async def handle_xhr_request(self, request: Request) -> Response:
        """处理XHR请求的具体逻辑"""
        pass


class NonXHRHandlerInterface(RequestHandlerInterface):
    """非XHR处理器接口"""
    
    @abstractmethod
    async def handle_non_xhr_request(self, request: Request) -> Response:
        """处理非XHR请求的具体逻辑"""
        pass


class HTTPRequestInterface(ABC):
    """HTTP请求处理接口"""
    
    @abstractmethod
    async def process_http_request(self, request: Request) -> Response:
        """处理HTTP请求"""
        pass


class RestAPIInterface(HTTPRequestInterface):
    """REST API处理接口"""
    
    @abstractmethod
    async def handle_rest_api(self, request: Request) -> Response:
        """处理REST API请求"""
        pass
    
    @abstractmethod
    async def handle_get(self, request: Request) -> Response:
        """处理GET请求"""
        pass
    
    @abstractmethod
    async def handle_post(self, request: Request) -> Response:
        """处理POST请求"""
        pass
    
    @abstractmethod
    async def handle_put(self, request: Request) -> Response:
        """处理PUT请求"""
        pass
    
    @abstractmethod
    async def handle_patch(self, request: Request) -> Response:
        """处理PATCH请求"""
        pass
    
    @abstractmethod
    async def handle_delete(self, request: Request) -> Response:
        """处理DELETE请求"""
        pass
    
    @abstractmethod
    async def handle_head(self, request: Request) -> Response:
        """处理HEAD请求"""
        pass
    
    @abstractmethod
    async def handle_options(self, request: Request) -> Response:
        """处理OPTIONS请求"""
        pass
    
    @abstractmethod
    async def validate_request(self, request: Request) -> bool:
        """验证请求有效性"""
        pass
    
    @abstractmethod
    async def authenticate_request(self, request: Request) -> Dict[str, Any]:
        """请求认证"""
        pass
    
    @abstractmethod
    async def authorize_request(self, request: Request, user_info: Dict[str, Any]) -> bool:
        """请求授权"""
        pass


class FileHandlerInterface(HTTPRequestInterface):
    """文件处理接口"""
    
    @abstractmethod
    async def handle_file_request(self, request: Request) -> Response:
        """处理文件请求"""
        pass


class NonHTTPHandlerInterface(ABC):
    """非HTTP请求处理接口"""
    
    @abstractmethod
    async def handle_non_http_request(self, request: Request) -> Any:
        """处理非HTTP请求（如SSE、WebSocket）"""
        pass


class SSEHandlerInterface(NonHTTPHandlerInterface):
    """SSE处理接口"""
    
    @abstractmethod
    async def handle_sse_request(self, request: Request) -> Any:
        """处理SSE请求"""
        pass


class WebSocketHandlerInterface(NonHTTPHandlerInterface):
    """WebSocket处理接口"""
    
    @abstractmethod
    async def handle_websocket_request(self, request: Request) -> Any:
        """处理WebSocket请求"""
        pass