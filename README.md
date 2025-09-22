# FastAPI Gateway System

基于FastAPI架构的分层网关系统，采用抽象类+实现类的设计模式。

## 架构设计

```
<<interface>>
    |
    ├── XHR处理类 ──┐
    │              ├── HTTP请求 ──┬── rest api
    │              │             └── file
    └── 非XHR处理类 ─┘
    
    非HTTP请求(sse/websocket)
```

## 系统特性

- **分层架构**: 采用抽象接口+具体实现的分层设计
- **请求分发**: 智能识别请求类型并路由到对应处理器
- **多协议支持**: 支持HTTP、WebSocket、SSE等多种协议
- **可扩展性**: 基于接口的设计，易于扩展新的处理器
- **统一处理**: 所有请求通过统一的网关入口处理

## 项目结构

```
/workspace/
├── interfaces.py      # 抽象接口定义
├── handlers.py        # 具体处理器实现
├── gateway.py         # 网关核心逻辑
├── main.py           # FastAPI应用入口
├── config.py         # 配置管理
├── requirements.txt  # 依赖包
└── static/          # 静态文件目录
```

## 核心组件

### 1. 抽象接口层 (interfaces.py)
- `GatewayInterface`: 网关顶层接口
- `RequestHandlerInterface`: 请求处理器基础接口
- `XHRHandlerInterface`: XHR处理器接口
- `NonXHRHandlerInterface`: 非XHR处理器接口
- 其他专用接口...

### 2. 处理器实现层 (handlers.py)
- `XHRHandler`: XHR请求处理器
- `NonXHRHandler`: 非XHR请求处理器
- `RestAPIHandler`: REST API处理器
- `FileHandler`: 文件处理器
- `SSEHandler`: SSE处理器
- `WebSocketHandler`: WebSocket处理器

### 3. 网关核心层 (gateway.py)
- `Gateway`: 主网关实现类
- `GatewayFactory`: 网关工厂类

### 4. 应用入口层 (main.py)
- FastAPI应用配置
- 路由定义
- 中间件配置
- 异常处理

## 快速开始

### 1. 安装依赖
```bash
pip install -r requirements.txt
```

### 2. 启动服务
```bash
python main.py
```

### 3. 访问服务
- 主页: http://localhost:8000/
- API文档: http://localhost:8000/docs
- 网关信息: http://localhost:8000/gateway/info
- 健康检查: http://localhost:8000/gateway/health
- SSE事件: http://localhost:8000/sse/events
- WebSocket: ws://localhost:8000/ws

## API端点

### 网关管理端点
- `GET /gateway/health` - 健康检查
- `GET /gateway/info` - 网关信息
- `GET /` - 根路径

### 特殊协议端点
- `GET /sse/events` - SSE事件流
- `WebSocket /ws` - WebSocket连接

### 通用端点
- `/{path:path}` - 所有其他请求通过网关处理

## 请求类型识别

网关会根据以下规则自动识别请求类型：

1. **SSE请求**: `Accept: text/event-stream` 或路径以 `/sse` 开头
2. **WebSocket请求**: `Upgrade: websocket` 头或路径以 `/ws` 开头
3. **文件请求**: 路径包含文件扩展名或静态文件路径
4. **REST API请求**: 路径以 `/api/` 开头
5. **XHR请求**: `X-Requested-With: XMLHttpRequest` 或 `Content-Type: application/json`
6. **其他**: 默认为非XHR请求

## 扩展指南

### 添加新的处理器

1. 在 `interfaces.py` 中定义抽象接口:
```python
class NewHandlerInterface(RequestHandlerInterface):
    @abstractmethod
    async def handle_new_request(self, request: Request) -> Response:
        pass
```

2. 在 `handlers.py` 中实现具体处理器:
```python
class NewHandler(NewHandlerInterface):
    async def handle(self, request: Request) -> Response:
        return await self.handle_new_request(request)
    
    async def handle_new_request(self, request: Request) -> Response:
        # 具体实现逻辑
        pass
```

3. 在 `gateway.py` 中注册处理器:
```python
def _init_handlers(self):
    # ... 现有处理器
    self.new_handler = NewHandler()
    self.handlers[RequestType.NEW_TYPE] = self.new_handler
```

## 配置说明

通过环境变量配置系统参数：
- `GATEWAY_HOST`: 服务器地址 (默认: 0.0.0.0)
- `GATEWAY_PORT`: 服务器端口 (默认: 8000)
- `GATEWAY_DEBUG`: 调试模式 (默认: true)
- `GATEWAY_STATIC_DIR`: 静态文件目录 (默认: /workspace/static)
- `GATEWAY_LOG_LEVEL`: 日志级别 (默认: INFO)
- `GATEWAY_REQUEST_TIMEOUT`: 请求超时时间 (默认: 30秒)

## 日志说明

系统包含完整的日志记录：
- 请求/响应日志
- 处理器选择日志
- 错误异常日志
- 性能统计日志

## 开发建议

1. **遵循接口契约**: 新增处理器必须实现对应的抽象接口
2. **错误处理**: 在处理器中适当处理异常并记录日志
3. **性能考虑**: 对于高频请求，考虑添加缓存机制
4. **安全性**: 文件处理器已包含路径安全检查，其他处理器也应考虑安全性
5. **测试**: 为新增的处理器编写单元测试

## 许可证

MIT License