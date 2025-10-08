# Model Mall Backend

基于 Go-Zero 框架构建的模型商城后端服务，集成图片分类和存储功能

---

> 🎯 **推荐：[SAVE_FULL_MODEL_GUIDE.md](SAVE_FULL_MODEL_GUIDE.md)** - 保存完整模型，部署超简单！⭐
> 
> 💡 **或者：[QUICK_REFERENCE.md](QUICK_REFERENCE.md)** - 使用Checkpoint部署（需要提供模型定义）

---

## ✨ 功能特性

- 🔐 基于RBAC的用户权限管理系统
- 🖼️ **图片上传和智能分类**
- 🤖 **集成机器学习模型服务**
- 💾 分类结果持久化存储
- 🐳 Docker容器化部署
- 📊 RESTful API接口

## 🚀 快速开始

### 🎯 最快部署（PyTorch模型）

**如果你有PyTorch训练的模型文件，查看：**

👉 **[START_HERE_PYTORCH.md](START_HERE_PYTORCH.md)** - 3步完成部署

### 前置要求

- Docker 和 Docker Compose
- Go 1.21+
- PostgreSQL 15
- Python 3.10+ （用于模型服务）
- 训练好的模型文件（推荐：PyTorch .mph/.pt/.pth）

### 快速部署（3步）

```bash
# 1. 复制模型定义到 model_service/app/model_architecture.py
# 2. 配置 model_service/docker-compose.yml (MODEL_ARCH, NUM_CLASSES)
# 3. 启动服务

cp checkpoint_best.pth model_service/models/
./start_all_services.sh
```

详细步骤查看 **[CUSTOM_MODEL_INTEGRATION.md](CUSTOM_MODEL_INTEGRATION.md)**

## 📋 项目结构说明

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

```
├── backend/                          # Go后端服务
│   ├── backend.api                   # API定义文件
│   ├── backend.go                    # 服务入口
│   ├── etc/
│   │   └── backend-api.yaml          # 配置文件
│   └── internal/
│       ├── config/                   # 配置
│       ├── handler/                  # HTTP处理器
│       │   ├── ImageHandler.go       # 图片上传处理
│       │   └── LoginHandler.go       # 登录处理
│       ├── logic/                    # 业务逻辑
│       │   ├── ImageLogic.go         # 图片分类业务
│       │   └── LoginLogic.go         # 登录业务
│       ├── models/                   # 数据模型
│       │   ├── image.go              # 图片模型
│       │   └── user.go               # 用户模型
│       ├── repository/               # 数据访问层
│       │   ├── image_repository.go   # 图片仓储
│       │   └── user_repository.go    # 用户仓储
│       ├── svc/                      # 服务组件
│       │   ├── ModelHelper.go        # 模型服务
│       │   ├── PGHelper.go           # PostgreSQL
│       │   └── RedisHelper.go        # Redis
│       └── types/                    # 类型定义
│           └── types.go
├── model_service/                    # Python模型服务
│   ├── app/
│   │   ├── api.py                    # Flask API服务
│   │   └── model_loader.py           # 模型加载器
│   ├── models/                       # 模型文件目录
│   │   ├── model.mph                 # 你的模型文件
│   │   └── labels.txt                # 分类标签
│   ├── Dockerfile                    # Docker镜像
│   ├── docker-compose.yml            # Docker配置
│   ├── requirements.txt              # Python依赖
│   └── README.md

## 🖼️ 图片分类功能

### 架构设计

```
用户上传图片 → Go后端 → Python模型服务 → 返回分类结果 → 存储到数据库
```

### API接口

#### 1. 上传图片并分类

```bash
POST /api/images/upload
Content-Type: multipart/form-data

curl -X POST http://localhost:8888/api/images/upload \
  -F "image=@photo.jpg"

# 响应
{
  "image_id": 1,
  "filename": "photo.jpg",
  "classifications": [
    {"label": "cat", "confidence": 0.85},
    {"label": "dog", "confidence": 0.12}
  ]
}
```

#### 2. 获取图片信息

```bash
GET /api/images/:id

curl http://localhost:8888/api/images/1
```

#### 3. 获取图片列表

```bash
GET /api/images?page=1&page_size=10

curl http://localhost:8888/api/images?page=1&page_size=10
```

### 模型支持

- ✅ **PyTorch模型** (.pt, .pth, .mph) - 推荐
- ✅ Keras模型 (.h5)
- ✅ ONNX模型 (.onnx)
- ✅ GPU加速（CUDA）
- ✅ 本地模型和远程模型服务

## 📚 文档

### 核心文档
- **[自定义模型集成](CUSTOM_MODEL_INTEGRATION.md)** - ⭐ PyTorch自定义模型部署
- **[分步操作指南](model_service/STEP_BY_STEP.md)** - 详细操作步骤
- **[部署指南](DEPLOYMENT.md)** - 完整部署流程
- **[API使用指南](IMAGE_CLASSIFICATION_GUIDE.md)** - 接口文档和使用说明
- **[模型服务文档](model_service/README.md)** - Python服务详细说明
- **[数据库设计](DATABASE_DESIGN.md)** - 数据库表结构

## 🔧 配置

### Go后端配置

文件：`backend/etc/backend-api.yaml`

```yaml
Name: backend-api
Host: 0.0.0.0
Port: 8888

Auth:
  AccessSecret: your-secret-key
  AccessExpire: 86400

PostgreSQL:
  Host: localhost
  Port: 5432
  Username: postgres
  Password: your-password
  Database: model_mall

Redis:
  Host: localhost
  Port: 6379
  Password: your-password

Model:
  Name: image-classifier
  Version: v1.0
  Type: remote  # local 或 remote
  Path: http://localhost:5000/classify
```

### 模型服务配置

文件：`model_service/docker-compose.yml`

```yaml
environment:
  - MODEL_PATH=/app/models/model.mph
  - MODEL_NAME=image-classifier
  - PORT=5000
```

## 🧪 测试

### 测试模型服务

```bash
cd model_service
python test_service.py test_image.jpg
```

### 测试Go后端

```bash
# 上传图片
curl -X POST http://localhost:8888/api/images/upload \
  -F "image=@test.jpg"

# 查看结果
curl http://localhost:8888/api/images/1
```

## 🐛 故障排查

### 模型服务无法启动

```bash
# 查看日志
docker-compose logs model-service

# 常见问题：
# 1. 模型文件不存在 → 检查 model_service/models/
# 2. 模型格式不支持 → 修改 CustomModelLoader
# 3. 内存不足 → 增加Docker内存限制
```

### Go后端无法连接模型服务

```bash
# 测试连接
curl http://localhost:5000/health

# 检查配置
cat backend/etc/backend-api.yaml | grep -A 4 Model
```

### 数据库连接失败

```bash
# 检查数据库
docker exec postgres pg_isready -U postgres

# 运行迁移
cd migrations && ./run_migrations.sh
```

## 🚀 性能优化

### 1. 模型服务优化

- 使用GPU加速
- 增加worker数量
- 模型量化和剪枝
- 批处理推理

### 2. 后端优化

- Redis缓存分类结果
- 数据库连接池优化
- 异步处理大文件

### 3. 部署优化

- 使用Nginx负载均衡
- 多实例部署
- CDN加速静态资源

## 🔒 安全建议

1. **认证授权** - 使用JWT Token
2. **文件验证** - 验证上传文件类型和大小
3. **限流** - 防止API滥用
4. **HTTPS** - 生产环境使用SSL
5. **数据加密** - 敏感数据加密存储

## 📊 监控

建议添加：
- Prometheus + Grafana 监控
- ELK 日志收集
- APM 性能追踪
- 健康检查和告警

## 🤝 贡献

欢迎提交Issue和Pull Request！

## 📄 许可证

MIT License

## 📞 联系方式

如有问题，请查看文档或提交Issue。

---

**技术栈：**
- Go + Go-Zero
- Python + Flask
- PostgreSQL
- Redis
- Docker
- TensorFlow/PyTorch/ONNX                     # 模型服务文档
├── migrations/                       # 数据库迁移
│   ├── 001_create_permissions_table.sql
│   ├── 005_create_images_table.sql   # 图片表
│   └── run_migrations.sh
├── docker-compose.yml                # 总体Docker配置
├── start_all_services.sh             # 一键启动脚本
├── DEPLOYMENT.md                     # 部署指南
├── MODEL_INTEGRATION_GUIDE.md        # 模型集成指南
├── IMAGE_CLASSIFICATION_GUIDE.md     # 图片分类使用指南
└── DATABASE_DESIGN.md                # 数据库设计文档
```