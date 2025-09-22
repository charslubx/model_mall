"""
网关系统的具体处理器实现
"""
import json
import os
from typing import Any, Dict, Optional
from fastapi import Request, Response, HTTPException
from fastapi.responses import JSONResponse, FileResponse, StreamingResponse
from interfaces import (
    XHRHandlerInterface, NonXHRHandlerInterface, RestAPIInterface, 
    FileHandlerInterface, SSEHandlerInterface, WebSocketHandlerInterface,
    RequestType
)
import asyncio
import aiohttp
import logging

# 配置日志
logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)


class XHRHandler(XHRHandlerInterface):
    """XHR请求处理器实现"""
    
    def __init__(self):
        self.session: Optional[aiohttp.ClientSession] = None
    
    async def _get_session(self) -> aiohttp.ClientSession:
        """获取或创建HTTP会话"""
        if self.session is None or self.session.closed:
            self.session = aiohttp.ClientSession()
        return self.session
    
    def can_handle(self, request: Request) -> bool:
        """判断是否为XHR请求"""
        xhr_header = request.headers.get("X-Requested-With", "").lower()
        content_type = request.headers.get("Content-Type", "").lower()
        
        return (
            xhr_header == "xmlhttprequest" or 
            "application/json" in content_type or
            "application/xml" in content_type
        )
    
    async def handle(self, request: Request) -> Response:
        """处理XHR请求"""
        return await self.handle_xhr_request(request)
    
    async def handle_xhr_request(self, request: Request) -> Response:
        """处理XHR请求的具体逻辑"""
        try:
            logger.info(f"处理XHR请求: {request.method} {request.url}")
            
            # 获取请求数据
            body = await request.body()
            headers = dict(request.headers)
            
            # 移除可能导致问题的headers
            headers.pop('host', None)
            headers.pop('content-length', None)
            
            # 这里可以添加请求转发逻辑
            # 示例：转发到后端服务
            session = await self._get_session()
            
            # 模拟处理逻辑
            if request.method == "GET":
                return JSONResponse({
                    "status": "success",
                    "message": "XHR GET request processed",
                    "data": {"method": request.method, "url": str(request.url)}
                })
            elif request.method == "POST":
                try:
                    request_data = json.loads(body) if body else {}
                except json.JSONDecodeError:
                    request_data = {"raw_body": body.decode('utf-8')}
                
                return JSONResponse({
                    "status": "success",
                    "message": "XHR POST request processed",
                    "data": request_data
                })
            else:
                return JSONResponse({
                    "status": "success",
                    "message": f"XHR {request.method} request processed"
                })
                
        except Exception as e:
            logger.error(f"XHR请求处理错误: {str(e)}")
            raise HTTPException(status_code=500, detail=f"XHR处理错误: {str(e)}")


class NonXHRHandler(NonXHRHandlerInterface):
    """非XHR请求处理器实现"""
    
    def can_handle(self, request: Request) -> bool:
        """判断是否为非XHR请求"""
        xhr_header = request.headers.get("X-Requested-With", "").lower()
        content_type = request.headers.get("Content-Type", "").lower()
        
        return not (
            xhr_header == "xmlhttprequest" or 
            "application/json" in content_type or
            "application/xml" in content_type
        )
    
    async def handle(self, request: Request) -> Response:
        """处理非XHR请求"""
        return await self.handle_non_xhr_request(request)
    
    async def handle_non_xhr_request(self, request: Request) -> Response:
        """处理非XHR请求的具体逻辑"""
        try:
            logger.info(f"处理非XHR请求: {request.method} {request.url}")
            
            # 根据路径判断是否为前端资源加载
            path = str(request.url.path)
            
            if path.startswith('/api/'):
                # API请求但非XHR，可能是表单提交等
                return JSONResponse({
                    "status": "success",
                    "message": "Non-XHR API request processed",
                    "path": path
                })
            else:
                # 可能是页面请求或静态资源
                return Response(
                    content=f"<html><body><h1>Gateway Response</h1><p>Path: {path}</p></body></html>",
                    media_type="text/html"
                )
                
        except Exception as e:
            logger.error(f"非XHR请求处理错误: {str(e)}")
            raise HTTPException(status_code=500, detail=f"非XHR处理错误: {str(e)}")


class RestAPIHandler(RestAPIInterface):
    """REST API处理器实现"""
    
    def __init__(self):
        self.session: Optional[aiohttp.ClientSession] = None
    
    async def _get_session(self) -> aiohttp.ClientSession:
        """获取或创建HTTP会话"""
        if self.session is None or self.session.closed:
            self.session = aiohttp.ClientSession()
        return self.session
    
    async def process_http_request(self, request: Request) -> Response:
        """处理HTTP请求"""
        return await self.handle_rest_api(request)
    
    async def handle_rest_api(self, request: Request) -> Response:
        """处理REST API请求"""
        try:
            logger.info(f"处理REST API请求: {request.method} {request.url}")
            
            path = str(request.url.path)
            
            # 模拟不同的API端点
            if path.startswith('/api/users'):
                return await self._handle_users_api(request)
            elif path.startswith('/api/data'):
                return await self._handle_data_api(request)
            else:
                return JSONResponse({
                    "status": "success",
                    "message": "Generic REST API response",
                    "endpoint": path,
                    "method": request.method
                })
                
        except Exception as e:
            logger.error(f"REST API处理错误: {str(e)}")
            raise HTTPException(status_code=500, detail=f"REST API处理错误: {str(e)}")
    
    async def _handle_users_api(self, request: Request) -> Response:
        """处理用户相关API"""
        if request.method == "GET":
            return JSONResponse({
                "users": [
                    {"id": 1, "name": "用户1", "email": "user1@example.com"},
                    {"id": 2, "name": "用户2", "email": "user2@example.com"}
                ]
            })
        elif request.method == "POST":
            body = await request.body()
            try:
                user_data = json.loads(body) if body else {}
                return JSONResponse({
                    "status": "created",
                    "user": user_data,
                    "id": 123
                })
            except json.JSONDecodeError:
                raise HTTPException(status_code=400, detail="Invalid JSON")
    
    async def _handle_data_api(self, request: Request) -> Response:
        """处理数据相关API"""
        return JSONResponse({
            "data": {
                "timestamp": "2025-09-22T00:00:00Z",
                "values": [1, 2, 3, 4, 5],
                "status": "active"
            }
        })


class FileHandler(FileHandlerInterface):
    """文件处理器实现"""
    
    def __init__(self, static_dir: str = "/workspace/static"):
        self.static_dir = static_dir
        # 确保静态文件目录存在
        os.makedirs(static_dir, exist_ok=True)
    
    async def process_http_request(self, request: Request) -> Response:
        """处理HTTP请求"""
        return await self.handle_file_request(request)
    
    async def handle_file_request(self, request: Request) -> Response:
        """处理文件请求"""
        try:
            path = str(request.url.path)
            logger.info(f"处理文件请求: {path}")
            
            # 移除开头的斜杠并构建完整路径
            file_path = path.lstrip('/')
            full_path = os.path.join(self.static_dir, file_path)
            
            # 安全检查：确保文件在静态目录内
            if not os.path.abspath(full_path).startswith(os.path.abspath(self.static_dir)):
                raise HTTPException(status_code=403, detail="Access denied")
            
            # 检查文件是否存在
            if os.path.exists(full_path) and os.path.isfile(full_path):
                return FileResponse(full_path)
            else:
                # 如果文件不存在，创建一个示例文件
                if file_path.endswith('.txt'):
                    with open(full_path, 'w', encoding='utf-8') as f:
                        f.write(f"这是一个示例文件: {file_path}\n创建时间: 2025-09-22")
                    return FileResponse(full_path)
                else:
                    raise HTTPException(status_code=404, detail="File not found")
                    
        except HTTPException:
            raise
        except Exception as e:
            logger.error(f"文件处理错误: {str(e)}")
            raise HTTPException(status_code=500, detail=f"文件处理错误: {str(e)}")


class SSEHandler(SSEHandlerInterface):
    """SSE处理器实现"""
    
    async def handle_non_http_request(self, request: Request) -> Any:
        """处理非HTTP请求"""
        return await self.handle_sse_request(request)
    
    async def handle_sse_request(self, request: Request) -> Any:
        """处理SSE请求"""
        async def event_generator():
            """SSE事件生成器"""
            count = 0
            while True:
                # 检查客户端是否断开连接
                if await request.is_disconnected():
                    break
                
                # 发送事件数据
                yield f"data: {{\"message\": \"SSE事件 #{count}\", \"timestamp\": \"{count}\"}}\n\n"
                count += 1
                
                # 等待一段时间
                await asyncio.sleep(2)
        
        return StreamingResponse(
            event_generator(),
            media_type="text/event-stream",
            headers={
                "Cache-Control": "no-cache",
                "Connection": "keep-alive",
                "Access-Control-Allow-Origin": "*",
                "Access-Control-Allow-Headers": "Cache-Control"
            }
        )


class WebSocketHandler(WebSocketHandlerInterface):
    """WebSocket处理器实现"""
    
    async def handle_non_http_request(self, request: Request) -> Any:
        """处理非HTTP请求"""
        return await self.handle_websocket_request(request)
    
    async def handle_websocket_request(self, request: Request) -> Any:
        """处理WebSocket请求"""
        # 注意：在实际的FastAPI应用中，WebSocket需要特殊处理
        # 这里返回一个指示，实际的WebSocket端点应该在main.py中单独定义
        return {
            "type": "websocket",
            "message": "WebSocket连接需要在应用中单独处理",
            "endpoint": "/ws"
        }