# 模型集成完整指南

本文档提供了将PyTorch完整模型集成到项目中的完整步骤。

> 💡 **推荐使用完整模型：** 使用 `torch.save(model, 'model.pth')` 保存，部署超简单！

## 📋 目录

1. [系统架构](#系统架构)
2. [快速开始](#快速开始)
3. [模型服务部署](#模型服务部署)
4. [Go后端配置](#go后端配置)
5. [完整部署流程](#完整部署流程)
6. [测试和验证](#测试和验证)
7. [故障排查](#故障排查)

## 🏗️ 系统架构

```
┌─────────────┐         ┌──────────────┐         ┌─────────────┐
│   用户端     │ HTTP    │   Go后端      │  HTTP   │  Python     │
│  (前端)     │────────▶│  (backend)   │────────▶│  模型服务   │
│             │         │              │         │             │
└─────────────┘         └──────────────┘         └─────────────┘
                               │                        │
                               │                        │
                               ▼                        ▼
                        ┌──────────────┐         ┌─────────────┐
                        │  PostgreSQL  │         │  模型文件   │
                        │   数据库     │         │  (.mph)     │
                        └──────────────┘         └─────────────┘
```

**工作流程：**

1. 用户上传图片到Go后端（`POST /api/images/upload`）
2. Go后端保存图片文件到磁盘
3. Go后端调用Python模型服务进行分类
4. Python模型服务加载模型，执行推理
5. 返回分类结果给Go后端
6. Go后端将结果保存到数据库
7. 返回分类结果给用户

## 🚀 快速开始

### 前置要求

- Docker 和 Docker Compose
- Go 1.21+
- PostgreSQL 数据库
- 训练好的`.mph`模型文件

### 5分钟快速部署

```bash
# 1. 克隆或进入项目目录
cd /workspace

# 2. 准备模型文件
cp /path/to/your/model.mph model_service/models/

# 3. 准备分类标签文件
# 编辑 model_service/models/labels.txt，每行一个标签

# 4. 启动模型服务
cd model_service
./start_with_docker.sh

# 5. 验证模型服务
curl http://localhost:5000/health

# 6. 运行数据库迁移
cd ../migrations
./run_migrations.sh

# 7. 启动Go后端
cd ../backend
go run backend.go
```

## 🐳 模型服务部署

### 步骤1: 准备模型文件

将你的`.mph`模型文件复制到模型目录：

```bash
cd /workspace/model_service
cp /path/to/your/model.mph models/
```

### 步骤2: 配置标签文件

编辑`models/labels.txt`，添加你的分类标签：

```bash
cat
dog
bird
horse
cow
# ... 更多标签
```

每行一个标签，顺序要与模型输出的类别索引对应。

### 步骤3: 修改模型加载器（如果需要）

如果`.mph`是特殊格式，需要修改`app/model_loader.py`中的`CustomModelLoader`类：

```python
def load_model(self):
    """加载.mph格式模型"""
    try:
        # 根据实际情况修改加载逻辑
        import your_model_library
        
        logger.info(f"正在加载.mph模型: {self.model_path}")
        self.model = your_model_library.load(self.model_path)
        self.is_loaded = True
        logger.info("模型加载成功")
    except Exception as e:
        logger.error(f"加载模型失败: {str(e)}")
        raise
```

### 步骤4: 更新配置

编辑`docker-compose.yml`：

```yaml
environment:
  - MODEL_PATH=/app/models/model.mph  # 改为你的模型文件名
  - MODEL_NAME=image-classifier
```

### 步骤5: 启动服务

**使用Docker（推荐）：**

```bash
cd /workspace/model_service
./start_with_docker.sh
```

**本地运行：**

```bash
cd /workspace/model_service
./start_service.sh
```

### 步骤6: 验证服务

```bash
# 检查健康状态
curl http://localhost:5000/health

# 查看模型信息
curl http://localhost:5000/info

# 测试分类（如果有测试图片）
curl -X POST http://localhost:5000/classify \
  -F "image=@test.jpg"
```

## ⚙️ Go后端配置

### 步骤1: 更新配置文件

编辑`backend/etc/backend-api.yaml`：

```yaml
Model:
  Name: image-classifier
  Version: v1.0
  Type: remote  # 使用远程模型服务
  Path: http://localhost:5000/classify  # 模型服务地址
```

如果使用Docker部署，地址改为：
```yaml
Path: http://model-service:5000/classify
```

### 步骤2: 运行数据库迁移

```bash
cd /workspace/migrations
chmod +x run_migrations.sh
./run_migrations.sh
```

这会创建以下表：
- `images` - 存储图片信息
- `image_classifications` - 存储分类标签

### 步骤3: 启动Go后端

```bash
cd /workspace/backend
go run backend.go
```

服务将在端口8888启动。

## 🌐 完整部署流程

### 方案1: Docker Compose 统一部署（推荐）

创建`/workspace/docker-compose.yml`：

```yaml
version: '3.8'

services:
  postgres:
    image: postgres:15
    container_name: postgres
    environment:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: your-password
      POSTGRES_DB: model_mall
    ports:
      - "5432:5432"
    volumes:
      - postgres-data:/var/lib/postgresql/data
      - ./migrations:/docker-entrypoint-initdb.d
    networks:
      - app-network

  redis:
    image: redis:7
    container_name: redis
    ports:
      - "6379:6379"
    command: redis-server --requirepass zqy5483201
    networks:
      - app-network

  model-service:
    build: ./model_service
    container_name: model-service
    volumes:
      - ./model_service/models:/app/models
    environment:
      - MODEL_PATH=/app/models/model.mph
      - MODEL_NAME=image-classifier
      - PORT=5000
    ports:
      - "5000:5000"
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:5000/health"]
      interval: 30s
      timeout: 10s
      retries: 3
    networks:
      - app-network

  backend:
    build: ./backend
    container_name: backend
    ports:
      - "8888:8888"
    environment:
      - MODEL_SERVICE_URL=http://model-service:5000/classify
    depends_on:
      - postgres
      - redis
      - model-service
    volumes:
      - ./backend/etc:/app/etc
      - ./uploads:/app/uploads
    networks:
      - app-network

volumes:
  postgres-data:

networks:
  app-network:
    driver: bridge
```

启动所有服务：

```bash
cd /workspace
docker-compose up -d
```

### 方案2: 分别部署

**1. 启动数据库（如果未启动）**

```bash
docker run -d \
  --name postgres \
  -e POSTGRES_PASSWORD=your-password \
  -e POSTGRES_DB=model_mall \
  -p 5432:5432 \
  postgres:15
```

**2. 启动模型服务**

```bash
cd /workspace/model_service
./start_with_docker.sh
```

**3. 运行数据库迁移**

```bash
cd /workspace/migrations
./run_migrations.sh
```

**4. 启动Go后端**

```bash
cd /workspace/backend
go run backend.go
```

## 🧪 测试和验证

### 1. 测试模型服务

```bash
cd /workspace/model_service
python test_service.py test_image.jpg
```

### 2. 测试完整流程

**上传图片：**

```bash
curl -X POST http://localhost:8888/api/images/upload \
  -H "Content-Type: multipart/form-data" \
  -F "image=@test_image.jpg"
```

**预期响应：**

```json
{
  "image_id": 1,
  "filename": "test_image.jpg",
  "file_path": "./uploads/1699999999.jpg",
  "file_size": 102400,
  "status": 1,
  "classifications": [
    {
      "label": "cat",
      "confidence": 0.8523
    },
    {
      "label": "dog",
      "confidence": 0.1234
    }
  ]
}
```

**获取图片信息：**

```bash
curl http://localhost:8888/api/images/1
```

**获取图片列表：**

```bash
curl http://localhost:8888/api/images?page=1&page_size=10
```

### 3. 查看日志

**模型服务日志：**
```bash
docker-compose logs -f model-service
```

**Go后端日志：**
```bash
# 如果使用Docker
docker-compose logs -f backend

# 如果本地运行，查看终端输出
```

## 🐛 故障排查

### 问题1: 模型服务启动失败

**症状：** `docker-compose up`后model-service容器退出

**检查：**
```bash
docker-compose logs model-service
```

**常见原因：**
1. 模型文件不存在
2. 模型文件格式不支持
3. 依赖包安装失败

**解决方案：**
1. 确认模型文件在`model_service/models/`目录
2. 修改`CustomModelLoader`支持`.mph`格式
3. 查看日志中的具体错误信息

### 问题2: Go后端无法连接模型服务

**症状：** 上传图片后返回"模型服务返回错误"

**检查：**
```bash
# 测试模型服务是否可访问
curl http://localhost:5000/health

# 如果使用Docker网络
docker exec backend curl http://model-service:5000/health
```

**解决方案：**
1. 检查`backend-api.yaml`中的模型服务地址配置
2. 确保模型服务已启动
3. 检查网络连接和防火墙设置

### 问题3: 分类结果为空或错误

**症状：** 能上传图片但分类结果不对

**检查：**
1. 查看模型服务日志
2. 验证labels.txt文件是否正确
3. 测试模型服务单独运行

**解决方案：**
1. 确认模型加载正确
2. 检查图片预处理逻辑
3. 验证标签映射关系

### 问题4: 数据库连接失败

**症状：** Go后端启动时报数据库连接错误

**检查：**
```bash
# 测试数据库连接
psql -h localhost -U postgres -d model_mall
```

**解决方案：**
1. 确认PostgreSQL已启动
2. 检查`backend-api.yaml`中的数据库配置
3. 运行数据库迁移脚本

## 📊 监控和维护

### 健康检查

定期检查服务状态：

```bash
# 模型服务
curl http://localhost:5000/health

# Go后端（需要添加健康检查端点）
curl http://localhost:8888/health
```

### 日志管理

```bash
# 查看模型服务日志
docker-compose logs -f model-service

# 查看Go后端日志
docker-compose logs -f backend

# 查看所有服务日志
docker-compose logs -f
```

### 性能监控

建议添加：
- Prometheus + Grafana 监控
- ELK 日志收集
- APM 性能追踪

## 🔒 生产环境建议

### 安全性

1. **使用HTTPS** - 配置SSL证书
2. **添加认证** - API Key或JWT
3. **限流** - 防止API滥用
4. **输入验证** - 验证上传文件

### 高可用

1. **负载均衡** - 使用Nginx
2. **服务冗余** - 多实例部署
3. **数据备份** - 定期备份数据库
4. **容器编排** - 使用Kubernetes

### 性能优化

1. **模型优化** - 量化、剪枝
2. **缓存** - Redis缓存结果
3. **GPU加速** - 使用GPU推理
4. **批处理** - 批量处理请求

## 📚 相关文档

- [图片分类服务使用指南](IMAGE_CLASSIFICATION_GUIDE.md)
- [模型服务README](model_service/README.md)
- [数据库设计文档](DATABASE_DESIGN.md)

## 💡 常见问题

**Q: .mph文件是什么格式？**

A: 请确认模型的训练框架。如果是Keras模型，通常是`.h5`格式。如需支持特殊格式，需要修改`CustomModelLoader`。

**Q: 如何更换模型？**

A: 将新模型文件放到`model_service/models/`，更新配置文件中的`MODEL_PATH`，重启服务即可。

**Q: 能否使用GPU加速？**

A: 可以。修改Dockerfile使用GPU版本的深度学习框架，并安装NVIDIA Docker运行时。

**Q: 如何添加多个模型？**

A: 可以修改代码支持多模型，或者部署多个模型服务实例，每个服务负责不同的模型。

## 📞 技术支持

如有问题，请查看：
1. 服务日志
2. 相关文档
3. 提交Issue
