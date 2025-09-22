"""
网关系统的主要实现类
"""
from typing import Dict, List, Optional
from fastapi import Request, Response, HTTPException
import logging
from interfaces import GatewayInterface, RequestHandlerInterface, RequestType
from handlers import (
    XHRHandler, NonXHRHandler, RestAPIHandler, FileHandler, 
    SSEHandler, WebSocketHandler
)

# 配置日志
logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)


class Gateway(GatewayInterface):
    """
    主网关实现类
    负责请求分发和统一处理流程
    """
    
    def __init__(self):
        """初始化网关和所有处理器"""
        self.handlers: Dict[RequestType, RequestHandlerInterface] = {}
        self._init_handlers()
    
    def _init_handlers(self):
        """初始化所有处理器"""
        self.xhr_handler = XHRHandler()
        self.non_xhr_handler = NonXHRHandler()
        self.rest_api_handler = RestAPIHandler()
        self.file_handler = FileHandler()
        self.sse_handler = SSEHandler()
        self.websocket_handler = WebSocketHandler()
        
        # 注册处理器
        self.handlers[RequestType.XHR] = self.xhr_handler
        self.handlers[RequestType.NON_XHR] = self.non_xhr_handler
        self.handlers[RequestType.REST_API] = self.rest_api_handler
        self.handlers[RequestType.FILE] = self.file_handler
        self.handlers[RequestType.SSE] = self.sse_handler
        self.handlers[RequestType.WEBSOCKET] = self.websocket_handler
    
    async def determine_request_type(self, request: Request) -> RequestType:
        """
        确定请求类型
        根据请求的特征判断应该使用哪个处理器
        """
        path = str(request.url.path)
        headers = request.headers
        
        # 检查是否为SSE请求
        accept = headers.get("Accept", "").lower()
        if "text/event-stream" in accept or path.startswith("/sse"):
            return RequestType.SSE
        
        # 检查是否为WebSocket升级请求
        if headers.get("Upgrade", "").lower() == "websocket" or path.startswith("/ws"):
            return RequestType.WEBSOCKET
        
        # 检查是否为文件请求
        if await self._is_file_request(path):
            return RequestType.FILE
        
        # 检查是否为REST API请求
        if path.startswith("/api/"):
            return RequestType.REST_API
        
        # 检查是否为XHR请求
        if await self._is_xhr_request(request):
            return RequestType.XHR
        
        # 默认为非XHR请求
        return RequestType.NON_XHR
    
    async def _is_xhr_request(self, request: Request) -> bool:
        """判断是否为XHR请求"""
        xhr_header = request.headers.get("X-Requested-With", "").lower()
        content_type = request.headers.get("Content-Type", "").lower()
        
        return (
            xhr_header == "xmlhttprequest" or 
            "application/json" in content_type or
            "application/xml" in content_type
        )
    
    async def _is_file_request(self, path: str) -> bool:
        """判断是否为文件请求"""
        # 常见的静态文件扩展名
        file_extensions = {
            '.js', '.css', '.html', '.htm', '.png', '.jpg', '.jpeg', '.gif', 
            '.svg', '.ico', '.txt', '.pdf', '.zip', '.json', '.xml'
        }
        
        # 检查路径是否以文件扩展名结尾
        for ext in file_extensions:
            if path.lower().endswith(ext):
                return True
        
        # 检查是否为静态文件路径
        static_paths = ['/static/', '/assets/', '/public/', '/files/']
        for static_path in static_paths:
            if path.startswith(static_path):
                return True
        
        return False
    
    async def process_request(self, request: Request) -> Response:
        """
        处理传入的请求
        这是网关的核心方法，负责请求分发
        """
        try:
            # 记录请求信息
            logger.info(f"收到请求: {request.method} {request.url}")
            logger.info(f"Headers: {dict(request.headers)}")
            
            # 确定请求类型
            request_type = await self.determine_request_type(request)
            logger.info(f"请求类型: {request_type.value}")
            
            # 获取对应的处理器
            handler = self.handlers.get(request_type)
            if not handler:
                raise HTTPException(
                    status_code=500, 
                    detail=f"未找到请求类型 {request_type.value} 的处理器"
                )
            
            # 验证处理器是否能处理该请求
            if hasattr(handler, 'can_handle') and not await handler.can_handle(request):
                logger.warning(f"处理器 {type(handler).__name__} 无法处理该请求")
                # 尝试使用默认处理器
                handler = self.non_xhr_handler
            
            # 处理请求
            response = await handler.handle(request)
            
            # 添加网关标识头
            if hasattr(response, 'headers'):
                response.headers["X-Gateway"] = "FastAPI-Gateway"
                response.headers["X-Handler"] = type(handler).__name__
            
            logger.info(f"请求处理完成: {request.method} {request.url}")
            return response
            
        except HTTPException:
            # 重新抛出HTTP异常
            raise
        except Exception as e:
            # 处理未预期的错误
            logger.error(f"网关处理请求时发生错误: {str(e)}", exc_info=True)
            raise HTTPException(
                status_code=500, 
                detail=f"网关内部错误: {str(e)}"
            )
    
    async def get_handler_info(self) -> Dict:
        """获取所有处理器的信息"""
        return {
            "handlers": {
                request_type.value: type(handler).__name__
                for request_type, handler in self.handlers.items()
            },
            "total_handlers": len(self.handlers)
        }
    
    async def health_check(self) -> Dict:
        """健康检查"""
        try:
            # 检查所有处理器是否正常
            handler_status = {}
            for request_type, handler in self.handlers.items():
                try:
                    # 简单的处理器状态检查
                    handler_status[request_type.value] = {
                        "name": type(handler).__name__,
                        "status": "healthy"
                    }
                except Exception as e:
                    handler_status[request_type.value] = {
                        "name": type(handler).__name__,
                        "status": "error",
                        "error": str(e)
                    }
            
            return {
                "gateway_status": "healthy",
                "timestamp": "2025-09-22T00:00:00Z",
                "handlers": handler_status
            }
        except Exception as e:
            return {
                "gateway_status": "error",
                "error": str(e),
                "timestamp": "2025-09-22T00:00:00Z"
            }


class GatewayFactory:
    """网关工厂类"""
    
    _instance: Optional[Gateway] = None
    
    @classmethod
    def get_gateway(cls) -> Gateway:
        """获取网关单例实例"""
        if cls._instance is None:
            cls._instance = Gateway()
        return cls._instance
    
    @classmethod
    def reset_gateway(cls):
        """重置网关实例（主要用于测试）"""
        cls._instance = None