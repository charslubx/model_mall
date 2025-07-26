# Model Mall Backend

基于 Go-Zero 框架构建的模型商城后端服务

## 项目结构说明

```text
├── backend                          # 后端服务主目录
│   ├── backend.api                  # API 接口定义文件，使用 go-zero 的 api 语法
│   ├── backend.go                   # 服务入口文件，包含主函数
│   ├── etc                          # 配置文件目录
│   │   └── backend-api.yaml         # 服务配置文件，包含数据库、缓存等配置
│   └── internal                     # 内部代码目录
│       ├── config                   # 配置结构定义目录
│       │   └── config.go            # 配置结构体定义文件
│       ├── handler                  # HTTP 处理器目录
│       │   ├── backendhandler.go    # 请求处理器实现文件
│       │   └── routes.go            # 路由注册和中间件配置
│       ├── logic                    # 业务逻辑目录
│       │   └── backendlogic.go      # 具体业务逻辑实现
│       ├── svc                      # 服务上下文目录
│       │   └── servicecontext.go    # 服务上下文定义，用于依赖注入
│       └── types                    # 数据类型定义目录
│           └── types.go             # 请求响应等数据结构定义
├── go.mod                           # Go 模块依赖定义文件
└── go.sum                           # Go 模块依赖版本锁定文件
```

## 目录结构说明

- `backend/`: 包含所有后端服务相关代码
  - `backend.api`: 使用 go-zero 的 api 语法定义 HTTP 接口
  - `backend.go`: 服务启动入口，初始化配置和服务
  - `etc/`: 存放配置文件
  - `internal/`: 内部实现代码
    - `config/`: 配置相关代码
    - `handler/`: HTTP 请求处理层，负责解析请求和返回响应
    - `logic/`: 业务逻辑层，实现具体的业务功能
    - `svc/`: 服务上下文，管理服务依赖（如数据库连接）
    - `types/`: 定义请求、响应等数据结构