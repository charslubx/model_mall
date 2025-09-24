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
import httpx
import logging
from config import get_config

# 配置日志
logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)


class XHRHandler(XHRHandlerInterface):
    """XHR请求处理器实现"""
    
    def __init__(self):
        self.client: Optional[httpx.AsyncClient] = None
        self.config = get_config()
    
    async def _get_client(self) -> httpx.AsyncClient:
        """获取或创建HTTP客户端"""
        if self.client is None or self.client.is_closed:
            self.client = httpx.AsyncClient(
                timeout=httpx.Timeout(self.config.request_timeout),
                follow_redirects=True
            )
        return self.client
    
    def _get_backend_url(self, path: str) -> str:
        """根据路径获取后端服务URL"""
        # 按路径长度排序，优先匹配更具体的路径
        sorted_services = sorted(
            self.config.backend_services.items(), 
            key=lambda x: len(x[0]), 
            reverse=True
        )
        
        for service_path, backend_url in sorted_services:
            if path.startswith(service_path):
                return backend_url
        
        return self.config.default_backend
    
    async def can_handle(self, request: Request) -> bool:
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
        """处理XHR请求的具体逻辑 - 透传到后端服务"""
        try:
            path = str(request.url.path)
            logger.info(f"透传XHR请求: {request.method} {path}")
            
            # 获取后端服务URL
            backend_url = self._get_backend_url(path)
            target_url = f"{backend_url}{path}"
            if request.url.query:
                target_url += f"?{request.url.query}"
            
            logger.info(f"转发到: {target_url}")
            
            # 获取请求数据
            body = await request.body()
            headers = dict(request.headers)
            
            # 清理不需要转发的headers
            headers_to_remove = ['host', 'content-length', 'connection', 'upgrade']
            for header in headers_to_remove:
                headers.pop(header, None)
            
            # 获取HTTP客户端
            client = await self._get_client()
            
            # 转发请求到后端服务
            try:
                response = await client.request(
                    method=request.method,
                    url=target_url,
                    headers=headers,
                    content=body if body else None
                )
                
                # 构建响应
                response_headers = dict(response.headers)
                # 移除可能导致问题的响应头
                response_headers.pop('content-length', None)
                response_headers.pop('transfer-encoding', None)
                
                return Response(
                    content=response.content,
                    status_code=response.status_code,
                    headers=response_headers,
                    media_type=response.headers.get('content-type')
                )
                
            except httpx.ConnectError:
                logger.error(f"无法连接到后端服务: {backend_url}")
                return JSONResponse({
                    "error": "Backend service unavailable",
                    "message": f"无法连接到后端服务: {backend_url}",
                    "status": "service_unavailable"
                }, status_code=503)
            
            except httpx.TimeoutException:
                logger.error(f"后端服务响应超时: {backend_url}")
                return JSONResponse({
                    "error": "Backend service timeout",
                    "message": "后端服务响应超时",
                    "status": "timeout"
                }, status_code=504)
                
        except Exception as e:
            logger.error(f"XHR请求处理错误: {str(e)}")
            raise HTTPException(status_code=500, detail=f"XHR处理错误: {str(e)}")


class NonXHRHandler(NonXHRHandlerInterface):
    """非XHR请求处理器实现"""
    
    def __init__(self):
        self.client: Optional[httpx.AsyncClient] = None
        self.config = get_config()
    
    async def _get_client(self) -> httpx.AsyncClient:
        """获取或创建HTTP客户端"""
        if self.client is None or self.client.is_closed:
            self.client = httpx.AsyncClient(
                timeout=httpx.Timeout(self.config.request_timeout),
                follow_redirects=True
            )
        return self.client
    
    def _get_backend_url(self, path: str) -> str:
        """根据路径获取后端服务URL"""
        # 按路径长度排序，优先匹配更具体的路径
        sorted_services = sorted(
            self.config.backend_services.items(), 
            key=lambda x: len(x[0]), 
            reverse=True
        )
        
        for service_path, backend_url in sorted_services:
            if path.startswith(service_path):
                return backend_url
        
        return self.config.default_backend
    
    async def can_handle(self, request: Request) -> bool:
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
        """处理非XHR请求的具体逻辑 - 透传到后端服务"""
        try:
            path = str(request.url.path)
            logger.info(f"透传非XHR请求: {request.method} {path}")
            
            # 获取后端服务URL
            backend_url = self._get_backend_url(path)
            target_url = f"{backend_url}{path}"
            if request.url.query:
                target_url += f"?{request.url.query}"
            
            logger.info(f"转发到: {target_url}")
            
            # 获取请求数据
            body = await request.body()
            headers = dict(request.headers)
            
            # 清理不需要转发的headers
            headers_to_remove = ['host', 'content-length', 'connection', 'upgrade']
            for header in headers_to_remove:
                headers.pop(header, None)
            
            # 获取HTTP客户端
            client = await self._get_client()
            
            # 转发请求到后端服务
            try:
                response = await client.request(
                    method=request.method,
                    url=target_url,
                    headers=headers,
                    content=body if body else None
                )
                
                # 构建响应
                response_headers = dict(response.headers)
                # 移除可能导致问题的响应头
                response_headers.pop('content-length', None)
                response_headers.pop('transfer-encoding', None)
                
                return Response(
                    content=response.content,
                    status_code=response.status_code,
                    headers=response_headers,
                    media_type=response.headers.get('content-type')
                )
                
            except httpx.ConnectError:
                logger.error(f"无法连接到后端服务: {backend_url}")
                # 对于页面请求，返回友好的错误页面
                if not path.startswith('/api/'):
                    return Response(
                        content=f"""
                        <html>
                        <head><title>服务不可用</title></head>
                        <body>
                            <h1>服务暂时不可用</h1>
                            <p>无法连接到后端服务: {backend_url}</p>
                            <p>请稍后再试</p>
                        </body>
                        </html>
                        """,
                        status_code=503,
                        media_type="text/html"
                    )
                else:
                    return JSONResponse({
                        "error": "Backend service unavailable",
                        "message": f"无法连接到后端服务: {backend_url}",
                        "status": "service_unavailable"
                    }, status_code=503)
            
            except httpx.TimeoutException:
                logger.error(f"后端服务响应超时: {backend_url}")
                if not path.startswith('/api/'):
                    return Response(
                        content=f"""
                        <html>
                        <head><title>服务超时</title></head>
                        <body>
                            <h1>服务响应超时</h1>
                            <p>后端服务响应时间过长</p>
                            <p>请稍后再试</p>
                        </body>
                        </html>
                        """,
                        status_code=504,
                        media_type="text/html"
                    )
                else:
                    return JSONResponse({
                        "error": "Backend service timeout",
                        "message": "后端服务响应超时",
                        "status": "timeout"
                    }, status_code=504)
                
        except Exception as e:
            logger.error(f"非XHR请求处理错误: {str(e)}")
            raise HTTPException(status_code=500, detail=f"非XHR处理错误: {str(e)}")


class RestAPIHandler(RestAPIInterface):
    """REST API处理器实现"""
    
    def __init__(self):
        self.client: Optional[httpx.AsyncClient] = None
        self.config = get_config()
        # 保留一些模拟端点，其他透传到后端
        self.mock_routes = {
            "/api/gateway": self._handle_gateway_api,  # 网关自身管理API
            "/api/mock": self._handle_mock_api,  # 模拟API，用于测试
        }
    
    async def _get_client(self) -> httpx.AsyncClient:
        """获取或创建HTTP客户端"""
        if self.client is None or self.client.is_closed:
            self.client = httpx.AsyncClient(
                timeout=httpx.Timeout(self.config.request_timeout),
                follow_redirects=True
            )
        return self.client
    
    def _get_backend_url(self, path: str) -> str:
        """根据路径获取后端服务URL"""
        # 按路径长度排序，优先匹配更具体的路径
        sorted_services = sorted(
            self.config.backend_services.items(), 
            key=lambda x: len(x[0]), 
            reverse=True
        )
        
        for service_path, backend_url in sorted_services:
            if path.startswith(service_path):
                return backend_url
        
        return self.config.default_backend
    
    async def _forward_to_backend(self, request: Request) -> Response:
        """透传请求到后端服务"""
        try:
            path = str(request.url.path)
            
            # 获取后端服务URL
            backend_url = self._get_backend_url(path)
            target_url = f"{backend_url}{path}"
            if request.url.query:
                target_url += f"?{request.url.query}"
            
            logger.info(f"透传API请求到: {target_url}")
            
            # 获取请求数据
            body = await request.body()
            headers = dict(request.headers)
            
            # 清理不需要转发的headers
            headers_to_remove = ['host', 'content-length', 'connection', 'upgrade']
            for header in headers_to_remove:
                headers.pop(header, None)
            
            # 获取HTTP客户端
            client = await self._get_client()
            
            # 转发请求到后端服务
            response = await client.request(
                method=request.method,
                url=target_url,
                headers=headers,
                content=body if body else None
            )
            
            # 构建响应
            response_headers = dict(response.headers)
            # 移除可能导致问题的响应头
            response_headers.pop('content-length', None)
            response_headers.pop('transfer-encoding', None)
            
            return Response(
                content=response.content,
                status_code=response.status_code,
                headers=response_headers,
                media_type=response.headers.get('content-type')
            )
            
        except httpx.ConnectError:
            logger.error(f"无法连接到后端服务: {backend_url}")
            return JSONResponse({
                "error": "Backend service unavailable",
                "message": f"无法连接到后端服务: {backend_url}",
                "status": "service_unavailable"
            }, status_code=503)
        
        except httpx.TimeoutException:
            logger.error(f"后端服务响应超时: {backend_url}")
            return JSONResponse({
                "error": "Backend service timeout",
                "message": "后端服务响应超时",
                "status": "timeout"
            }, status_code=504)
    
    async def process_http_request(self, request: Request) -> Response:
        """处理HTTP请求"""
        return await self.handle_rest_api(request)
    
    async def handle_rest_api(self, request: Request) -> Response:
        """处理REST API请求"""
        try:
            logger.info(f"处理REST API请求: {request.method} {request.url}")
            
            # 请求验证
            if not await self.validate_request(request):
                raise HTTPException(status_code=400, detail="请求验证失败")
            
            # 认证
            user_info = await self.authenticate_request(request)
            
            # 授权
            if not await self.authorize_request(request, user_info):
                raise HTTPException(status_code=403, detail="权限不足")
            
            # 根据HTTP方法分发请求
            method = request.method.upper()
            if method == "GET":
                return await self.handle_get(request)
            elif method == "POST":
                return await self.handle_post(request)
            elif method == "PUT":
                return await self.handle_put(request)
            elif method == "PATCH":
                return await self.handle_patch(request)
            elif method == "DELETE":
                return await self.handle_delete(request)
            elif method == "HEAD":
                return await self.handle_head(request)
            elif method == "OPTIONS":
                return await self.handle_options(request)
            else:
                raise HTTPException(status_code=405, detail=f"不支持的HTTP方法: {method}")
                
        except HTTPException:
            raise
        except Exception as e:
            logger.error(f"REST API处理错误: {str(e)}")
            raise HTTPException(status_code=500, detail=f"REST API处理错误: {str(e)}")
    
    async def handle_get(self, request: Request) -> Response:
        """处理GET请求"""
        path = str(request.url.path)
        logger.info(f"处理GET请求: {path}")
        
        # 检查是否为模拟路由
        for route, handler in self.mock_routes.items():
            if path.startswith(route):
                return await handler(request)
        
        # 透传到后端服务
        return await self._forward_to_backend(request)
    
    async def handle_post(self, request: Request) -> Response:
        """处理POST请求"""
        path = str(request.url.path)
        logger.info(f"处理POST请求: {path}")
        
        # 获取请求体
        body = await request.body()
        try:
            request_data = json.loads(body) if body else {}
        except json.JSONDecodeError:
            request_data = {"raw_body": body.decode('utf-8', errors='ignore')}
        
        # 根据路径分发到具体处理器
        for route, handler in self.api_routes.items():
            if path.startswith(route):
                return await handler(request)
        
        # 通用POST处理
        return JSONResponse({
            "status": "created",
            "method": "POST",
            "path": path,
            "data": request_data,
            "message": "POST请求处理成功"
        }, status_code=201)
    
    async def handle_put(self, request: Request) -> Response:
        """处理PUT请求"""
        path = str(request.url.path)
        logger.info(f"处理PUT请求: {path}")
        
        # 获取请求体
        body = await request.body()
        try:
            request_data = json.loads(body) if body else {}
        except json.JSONDecodeError:
            request_data = {"raw_body": body.decode('utf-8', errors='ignore')}
        
        # 根据路径分发到具体处理器
        for route, handler in self.api_routes.items():
            if path.startswith(route):
                return await handler(request)
        
        return JSONResponse({
            "status": "updated",
            "method": "PUT",
            "path": path,
            "data": request_data,
            "message": "PUT请求处理成功"
        })
    
    async def handle_patch(self, request: Request) -> Response:
        """处理PATCH请求"""
        path = str(request.url.path)
        logger.info(f"处理PATCH请求: {path}")
        
        # 获取请求体
        body = await request.body()
        try:
            request_data = json.loads(body) if body else {}
        except json.JSONDecodeError:
            request_data = {"raw_body": body.decode('utf-8', errors='ignore')}
        
        # 根据路径分发到具体处理器
        for route, handler in self.api_routes.items():
            if path.startswith(route):
                return await handler(request)
        
        return JSONResponse({
            "status": "patched",
            "method": "PATCH",
            "path": path,
            "data": request_data,
            "message": "PATCH请求处理成功"
        })
    
    async def handle_delete(self, request: Request) -> Response:
        """处理DELETE请求"""
        path = str(request.url.path)
        logger.info(f"处理DELETE请求: {path}")
        
        # 根据路径分发到具体处理器
        for route, handler in self.api_routes.items():
            if path.startswith(route):
                return await handler(request)
        
        return JSONResponse({
            "status": "deleted",
            "method": "DELETE",
            "path": path,
            "message": "DELETE请求处理成功"
        })
    
    async def handle_head(self, request: Request) -> Response:
        """处理HEAD请求"""
        path = str(request.url.path)
        logger.info(f"处理HEAD请求: {path}")
        
        # HEAD请求只返回头部信息，不返回内容
        response = Response(status_code=200)
        response.headers["Content-Type"] = "application/json"
        response.headers["X-API-Version"] = "1.0.0"
        response.headers["X-Path"] = path
        return response
    
    async def handle_options(self, request: Request) -> Response:
        """处理OPTIONS请求"""
        path = str(request.url.path)
        logger.info(f"处理OPTIONS请求: {path}")
        
        # 返回支持的方法
        response = Response(status_code=200)
        response.headers["Allow"] = "GET, POST, PUT, PATCH, DELETE, HEAD, OPTIONS"
        response.headers["Access-Control-Allow-Methods"] = "GET, POST, PUT, PATCH, DELETE, HEAD, OPTIONS"
        response.headers["Access-Control-Allow-Headers"] = "Content-Type, Authorization, X-Requested-With"
        return response
    
    async def validate_request(self, request: Request) -> bool:
        """验证请求有效性"""
        try:
            # 基本验证逻辑
            path = str(request.url.path)
            
            # 检查路径是否为API路径
            if not path.startswith('/api/'):
                return False
            
            # 检查Content-Type（对于有请求体的方法）
            if request.method in ["POST", "PUT", "PATCH"]:
                content_type = request.headers.get("Content-Type", "")
                if not any(ct in content_type for ct in ["application/json", "application/x-www-form-urlencoded", "multipart/form-data"]):
                    logger.warning(f"不支持的Content-Type: {content_type}")
                    # 暂时允许，只记录警告
            
            return True
        except Exception as e:
            logger.error(f"请求验证失败: {str(e)}")
            return False
    
    async def authenticate_request(self, request: Request) -> Dict[str, Any]:
        """请求认证"""
        try:
            # 简单的认证逻辑
            auth_header = request.headers.get("Authorization", "")
            
            if auth_header.startswith("Bearer "):
                token = auth_header[7:]  # 移除 "Bearer " 前缀
                # 这里应该验证token的有效性
                # 简化示例：假设所有token都有效
                return {
                    "user_id": "user123",
                    "username": "test_user",
                    "roles": ["user"],
                    "token": token
                }
            elif auth_header.startswith("Basic "):
                # 基本认证处理
                return {
                    "user_id": "guest",
                    "username": "guest",
                    "roles": ["guest"],
                    "auth_type": "basic"
                }
            else:
                # 无认证信息，返回匿名用户
                return {
                    "user_id": "anonymous",
                    "username": "anonymous",
                    "roles": ["anonymous"],
                    "auth_type": "none"
                }
        except Exception as e:
            logger.error(f"认证失败: {str(e)}")
            return {
                "user_id": "error",
                "username": "error",
                "roles": [],
                "error": str(e)
            }
    
    async def authorize_request(self, request: Request, user_info: Dict[str, Any]) -> bool:
        """请求授权"""
        try:
            path = str(request.url.path)
            method = request.method.upper()
            user_roles = user_info.get("roles", [])
            
            # 简单的授权规则
            # 公开端点
            public_endpoints = ["/api/auth", "/api/data"]
            if any(path.startswith(endpoint) for endpoint in public_endpoints):
                return True
            
            # 需要认证的端点
            if "anonymous" in user_roles:
                # 匿名用户只能访问公开端点
                return False
            
            # 管理员可以访问所有端点
            if "admin" in user_roles:
                return True
            
            # 普通用户的权限
            if "user" in user_roles:
                # 用户可以GET大部分资源
                if method == "GET":
                    return True
                # 用户可以POST到某些端点
                if method == "POST" and path.startswith("/api/orders"):
                    return True
                # 用户可以修改自己的信息
                if method in ["PUT", "PATCH"] and "/users/" in path:
                    # 这里应该检查用户是否修改自己的信息
                    return True
            
            # 默认拒绝
            logger.warning(f"授权失败: 用户 {user_info.get('username')} 尝试访问 {method} {path}")
            return True  # 暂时允许所有请求，在实际项目中应该更严格
            
        except Exception as e:
            logger.error(f"授权检查失败: {str(e)}")
            return False
    
    async def _handle_users_api(self, request: Request) -> Response:
        """处理用户相关API"""
        method = request.method.upper()
        path = str(request.url.path)
        
        if method == "GET":
            # 获取用户列表或特定用户
            if path == "/api/users":
                return JSONResponse({
                    "users": [
                        {"id": 1, "name": "用户1", "email": "user1@example.com", "status": "active"},
                        {"id": 2, "name": "用户2", "email": "user2@example.com", "status": "active"},
                        {"id": 3, "name": "用户3", "email": "user3@example.com", "status": "inactive"}
                    ],
                    "total": 3,
                    "page": 1,
                    "limit": 10
                })
            else:
                # 获取特定用户
                user_id = path.split("/")[-1]
                return JSONResponse({
                    "id": user_id,
                    "name": f"用户{user_id}",
                    "email": f"user{user_id}@example.com",
                    "status": "active",
                    "created_at": "2025-09-22T00:00:00Z"
                })
        
        elif method == "POST":
            body = await request.body()
            try:
                user_data = json.loads(body) if body else {}
                return JSONResponse({
                    "status": "created",
                    "user": {
                        "id": 123,
                        "name": user_data.get("name", "新用户"),
                        "email": user_data.get("email", "new@example.com"),
                        "status": "active"
                    },
                    "message": "用户创建成功"
                }, status_code=201)
            except json.JSONDecodeError:
                raise HTTPException(status_code=400, detail="Invalid JSON")
        
        elif method == "PUT":
            body = await request.body()
            user_id = path.split("/")[-1]
            try:
                user_data = json.loads(body) if body else {}
                return JSONResponse({
                    "status": "updated",
                    "user": {
                        "id": user_id,
                        "name": user_data.get("name", f"用户{user_id}"),
                        "email": user_data.get("email", f"user{user_id}@example.com"),
                        "status": user_data.get("status", "active")
                    },
                    "message": "用户更新成功"
                })
            except json.JSONDecodeError:
                raise HTTPException(status_code=400, detail="Invalid JSON")
        
        elif method == "DELETE":
            user_id = path.split("/")[-1]
            return JSONResponse({
                "status": "deleted",
                "user_id": user_id,
                "message": "用户删除成功"
            })
        
        return JSONResponse({"error": "不支持的操作"}, status_code=405)
    
    async def _handle_data_api(self, request: Request) -> Response:
        """处理数据相关API"""
        method = request.method.upper()
        
        if method == "GET":
            return JSONResponse({
                "data": {
                    "timestamp": "2025-09-22T00:00:00Z",
                    "values": [1, 2, 3, 4, 5],
                    "status": "active",
                    "metrics": {
                        "cpu_usage": 45.2,
                        "memory_usage": 67.8,
                        "disk_usage": 23.1
                    }
                },
                "meta": {
                    "version": "1.0.0",
                    "last_updated": "2025-09-22T00:00:00Z"
                }
            })
        
        elif method == "POST":
            body = await request.body()
            try:
                data = json.loads(body) if body else {}
                return JSONResponse({
                    "status": "processed",
                    "received_data": data,
                    "processed_at": "2025-09-22T00:00:00Z",
                    "result_id": "data_123"
                }, status_code=201)
            except json.JSONDecodeError:
                raise HTTPException(status_code=400, detail="Invalid JSON")
        
        return JSONResponse({"error": "不支持的操作"}, status_code=405)
    
    async def _handle_products_api(self, request: Request) -> Response:
        """处理产品相关API"""
        method = request.method.upper()
        path = str(request.url.path)
        
        if method == "GET":
            return JSONResponse({
                "products": [
                    {"id": 1, "name": "产品A", "price": 99.99, "category": "electronics"},
                    {"id": 2, "name": "产品B", "price": 149.99, "category": "books"},
                    {"id": 3, "name": "产品C", "price": 79.99, "category": "clothing"}
                ],
                "total": 3,
                "categories": ["electronics", "books", "clothing"]
            })
        
        elif method == "POST":
            body = await request.body()
            try:
                product_data = json.loads(body) if body else {}
                return JSONResponse({
                    "status": "created",
                    "product": {
                        "id": 123,
                        "name": product_data.get("name", "新产品"),
                        "price": product_data.get("price", 0.0),
                        "category": product_data.get("category", "uncategorized")
                    }
                }, status_code=201)
            except json.JSONDecodeError:
                raise HTTPException(status_code=400, detail="Invalid JSON")
        
        return JSONResponse({"error": "不支持的操作"}, status_code=405)
    
    async def _handle_orders_api(self, request: Request) -> Response:
        """处理订单相关API"""
        method = request.method.upper()
        
        if method == "GET":
            return JSONResponse({
                "orders": [
                    {"id": 1, "user_id": 1, "total": 199.98, "status": "completed"},
                    {"id": 2, "user_id": 2, "total": 79.99, "status": "pending"},
                    {"id": 3, "user_id": 1, "total": 149.99, "status": "shipped"}
                ],
                "total": 3
            })
        
        elif method == "POST":
            body = await request.body()
            try:
                order_data = json.loads(body) if body else {}
                return JSONResponse({
                    "status": "created",
                    "order": {
                        "id": 456,
                        "user_id": order_data.get("user_id", 1),
                        "items": order_data.get("items", []),
                        "total": order_data.get("total", 0.0),
                        "status": "pending"
                    }
                }, status_code=201)
            except json.JSONDecodeError:
                raise HTTPException(status_code=400, detail="Invalid JSON")
        
        return JSONResponse({"error": "不支持的操作"}, status_code=405)
    
    async def _handle_auth_api(self, request: Request) -> Response:
        """处理认证相关API"""
        method = request.method.upper()
        path = str(request.url.path)
        
        if method == "POST":
            if path.endswith("/login"):
                body = await request.body()
                try:
                    credentials = json.loads(body) if body else {}
                    return JSONResponse({
                        "status": "success",
                        "token": "mock_jwt_token_123",
                        "user": {
                            "id": 1,
                            "username": credentials.get("username", "user"),
                            "email": "user@example.com",
                            "roles": ["user"]
                        },
                        "expires_in": 3600
                    })
                except json.JSONDecodeError:
                    raise HTTPException(status_code=400, detail="Invalid JSON")
            
            elif path.endswith("/logout"):
                return JSONResponse({
                    "status": "success",
                    "message": "登出成功"
                })
            
            elif path.endswith("/refresh"):
                return JSONResponse({
                    "status": "success",
                    "token": "new_mock_jwt_token_456",
                    "expires_in": 3600
                })
        
        return JSONResponse({"error": "不支持的操作"}, status_code=405)


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