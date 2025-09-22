"""
FastAPI网关应用程序入口
"""
from fastapi import FastAPI, Request, WebSocket, WebSocketDisconnect
from fastapi.middleware.cors import CORSMiddleware
from fastapi.responses import JSONResponse
import uvicorn
import logging
from gateway import GatewayFactory
from typing import Dict, Any
import json
import asyncio

# 配置日志
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s'
)
logger = logging.getLogger(__name__)

# 创建FastAPI应用
app = FastAPI(
    title="Gateway System",
    description="基于FastAPI的分层网关系统",
    version="1.0.0",
    docs_url="/docs",
    redoc_url="/redoc"
)

# 添加CORS中间件
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],  # 在生产环境中应该设置具体的域名
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

# 获取网关实例
gateway = GatewayFactory.get_gateway()


@app.middleware("http")
async def log_requests(request: Request, call_next):
    """请求日志中间件"""
    start_time = asyncio.get_event_loop().time()
    
    # 记录请求开始
    logger.info(f"请求开始: {request.method} {request.url}")
    
    # 处理请求
    response = await call_next(request)
    
    # 计算处理时间
    process_time = asyncio.get_event_loop().time() - start_time
    
    # 记录请求结束
    logger.info(f"请求完成: {request.method} {request.url} - 状态码: {response.status_code} - 耗时: {process_time:.4f}s")
    
    # 添加处理时间到响应头
    response.headers["X-Process-Time"] = str(process_time)
    
    return response


@app.exception_handler(Exception)
async def global_exception_handler(request: Request, exc: Exception):
    """全局异常处理器"""
    logger.error(f"未处理的异常: {str(exc)}", exc_info=True)
    return JSONResponse(
        status_code=500,
        content={
            "error": "内部服务器错误",
            "message": str(exc),
            "path": str(request.url)
        }
    )


# 主要的网关路由 - 捕获所有HTTP请求
@app.api_route("/{path:path}", methods=["GET", "POST", "PUT", "DELETE", "PATCH", "HEAD", "OPTIONS"])
async def gateway_handler(request: Request):
    """
    网关主处理器
    捕获所有HTTP请求并通过网关处理
    """
    return await gateway.process_request(request)


# WebSocket端点
@app.websocket("/ws")
async def websocket_endpoint(websocket: WebSocket):
    """WebSocket端点处理"""
    await websocket.accept()
    logger.info(f"WebSocket连接已建立: {websocket.client}")
    
    try:
        # 发送欢迎消息
        await websocket.send_text(json.dumps({
            "type": "welcome",
            "message": "WebSocket连接成功建立",
            "timestamp": "2025-09-22T00:00:00Z"
        }))
        
        # 消息循环
        message_count = 0
        while True:
            try:
                # 接收客户端消息
                data = await websocket.receive_text()
                logger.info(f"收到WebSocket消息: {data}")
                
                # 解析消息
                try:
                    message = json.loads(data)
                except json.JSONDecodeError:
                    message = {"raw_message": data}
                
                message_count += 1
                
                # 回复消息
                response = {
                    "type": "response",
                    "message": f"收到消息 #{message_count}",
                    "received_data": message,
                    "timestamp": "2025-09-22T00:00:00Z"
                }
                
                await websocket.send_text(json.dumps(response))
                
            except WebSocketDisconnect:
                logger.info(f"WebSocket客户端断开连接: {websocket.client}")
                break
            except Exception as e:
                logger.error(f"WebSocket处理消息时出错: {str(e)}")
                await websocket.send_text(json.dumps({
                    "type": "error",
                    "message": f"处理消息时出错: {str(e)}"
                }))
                
    except WebSocketDisconnect:
        logger.info(f"WebSocket连接断开: {websocket.client}")
    except Exception as e:
        logger.error(f"WebSocket连接出错: {str(e)}")


# 特殊端点：网关状态和管理
@app.get("/gateway/health")
async def gateway_health():
    """网关健康检查"""
    return await gateway.health_check()


@app.get("/gateway/info")
async def gateway_info():
    """获取网关信息"""
    return {
        "gateway": "FastAPI Gateway System",
        "version": "1.0.0",
        "handlers": await gateway.get_handler_info(),
        "endpoints": {
            "health": "/gateway/health",
            "info": "/gateway/info",
            "websocket": "/ws",
            "sse": "/sse/events",
            "docs": "/docs"
        }
    }


# SSE端点
@app.get("/sse/events")
async def sse_events(request: Request):
    """SSE事件流端点"""
    from handlers import SSEHandler
    sse_handler = SSEHandler()
    return await sse_handler.handle_sse_request(request)


# 根路径
@app.get("/")
async def root():
    """根路径处理"""
    return {
        "message": "欢迎使用FastAPI网关系统",
        "version": "1.0.0",
        "docs": "/docs",
        "health": "/gateway/health",
        "info": "/gateway/info"
    }


if __name__ == "__main__":
    # 开发环境启动配置
    uvicorn.run(
        "main:app",
        host="0.0.0.0",
        port=8000,
        reload=True,
        log_level="info"
    )