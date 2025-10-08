# 部署指南

## 快速部署（适用于PyTorch Checkpoint模型）

### 前提条件
- Docker 和 Docker Compose
- PyTorch checkpoint文件（.pth/.pt/.mph）
- 训练代码中的模型类定义

### 部署步骤

#### 1. 准备模型定义

将训练代码中的模型类复制到：
```
/workspace/model_service/app/model_architecture.py
```

详细步骤参考: [CUSTOM_MODEL_INTEGRATION.md](CUSTOM_MODEL_INTEGRATION.md)

#### 2. 配置环境

编辑 `model_service/docker-compose.yml`:
```yaml
environment:
  - MODEL_PATH=/app/models/checkpoint_best.pth
  - MODEL_ARCH=my_model      # 你定义的模型名
  - NUM_CLASSES=10           # 类别数量
```

#### 3. 准备文件和启动

```bash
# 复制checkpoint
cp checkpoint_best.pth model_service/models/

# 创建标签文件
cat > model_service/models/labels.txt << EOF
类别1
类别2
...
EOF

# 一键启动
./start_all_services.sh
```

#### 4. 验证

```bash
curl http://localhost:5000/health
curl -X POST http://localhost:8888/api/images/upload -F "image=@test.jpg"
```

## 详细部署步骤

### 方式1: 使用Docker Compose（推荐）

**优点：** 简单、一致、易于管理

```bash
# 1. 准备模型文件
cp /path/to/your/model.mph model_service/models/

# 2. 启动所有服务
docker-compose up -d

# 3. 查看日志
docker-compose logs -f

# 4. 停止服务
docker-compose down
```

### 方式2: 分别部署各服务

#### 2.1 启动数据库

```bash
docker run -d \
  --name postgres \
  -e POSTGRES_PASSWORD=your-password \
  -e POSTGRES_DB=model_mall \
  -p 5432:5432 \
  postgres:15
```

#### 2.2 启动Redis

```bash
docker run -d \
  --name redis \
  -p 6379:6379 \
  redis:7-alpine redis-server --requirepass zqy5483201
```

#### 2.3 运行数据库迁移

```bash
cd migrations
./run_migrations.sh
cd ..
```

#### 2.4 启动模型服务

```bash
cd model_service
./start_with_docker.sh
cd ..
```

#### 2.5 启动Go后端

```bash
cd backend
go run backend.go
```

## 配置说明

### 模型服务配置

文件：`model_service/docker-compose.yml`

```yaml
environment:
  - MODEL_PATH=/app/models/model.mph  # 模型文件路径
  - MODEL_NAME=image-classifier        # 模型名称
```

### Go后端配置

文件：`backend/etc/backend-api.yaml`

```yaml
Model:
  Name: image-classifier
  Version: v1.0
  Type: remote
  Path: http://localhost:5000/classify
```

如果使用Docker网络：
```yaml
Path: http://model-service:5000/classify
```

### 数据库配置

文件：`backend/etc/backend-api.yaml`

```yaml
PostgreSQL:
  Host: localhost
  Port: 5432
  Username: postgres
  Password: your-password
  Database: model_mall
```

## 环境变量

创建`.env`文件（可选）：

```bash
# 数据库
POSTGRES_PASSWORD=your-password
POSTGRES_DB=model_mall

# Redis
REDIS_PASSWORD=zqy5483201

# 模型服务
MODEL_PATH=/app/models/model.mph
MODEL_NAME=image-classifier

# Go后端
PORT=8888
```

## 验证部署

### 检查服务状态

```bash
# 查看所有容器
docker-compose ps

# 检查模型服务
curl http://localhost:5000/health
# 预期: {"status":"healthy","model_loaded":true}

# 检查模型信息
curl http://localhost:5000/info

# 测试数据库连接
docker exec postgres psql -U postgres -d model_mall -c "SELECT 1;"
```

### 完整测试流程

```bash
# 1. 上传图片
curl -X POST http://localhost:8888/api/images/upload \
  -F "image=@test_image.jpg"

# 预期响应:
# {
#   "image_id": 1,
#   "filename": "test_image.jpg",
#   "classifications": [
#     {"label": "cat", "confidence": 0.85}
#   ]
# }

# 2. 获取图片信息
curl http://localhost:8888/api/images/1

# 3. 获取图片列表
curl http://localhost:8888/api/images?page=1&page_size=10
```

## 日志查看

```bash
# 查看所有服务日志
docker-compose logs -f

# 查看特定服务日志
docker-compose logs -f model-service
docker-compose logs -f postgres
docker-compose logs -f redis

# 查看Go后端日志（如果本地运行）
tail -f backend/logs/app.log
```

## 常见问题

### 1. 模型服务启动失败

**检查：**
```bash
docker-compose logs model-service
```

**常见原因：**
- 模型文件不存在
- 模型格式不支持
- 内存不足

**解决：**
```bash
# 检查模型文件
ls -lh model_service/models/

# 查看详细错误
docker-compose logs --tail=100 model-service
```

### 2. 数据库连接失败

**检查：**
```bash
docker exec postgres pg_isready -U postgres
```

**解决：**
```bash
# 重启数据库
docker-compose restart postgres

# 检查配置
cat backend/etc/backend-api.yaml | grep -A 5 PostgreSQL
```

### 3. Go后端无法连接模型服务

**检查：**
```bash
# 从Go后端容器测试连接
docker exec backend curl http://model-service:5000/health
```

**解决：**
- 确认服务在同一Docker网络
- 检查防火墙设置
- 验证配置文件中的URL

## 性能优化

### 1. 模型服务

```yaml
# docker-compose.yml
model-service:
  environment:
    - WORKERS=4        # 增加worker数量
    - THREADS=2        # 每个worker的线程数
  deploy:
    resources:
      limits:
        cpus: '2'      # CPU限制
        memory: 4G     # 内存限制
```

### 2. 数据库

```yaml
postgres:
  command: postgres -c max_connections=200 -c shared_buffers=256MB
```

### 3. Redis

```yaml
redis:
  command: redis-server --maxmemory 512mb --maxmemory-policy allkeys-lru
```

## 生产环境建议

### 1. 使用环境变量

不要在配置文件中硬编码密码，使用环境变量：

```yaml
environment:
  - POSTGRES_PASSWORD=${POSTGRES_PASSWORD}
  - REDIS_PASSWORD=${REDIS_PASSWORD}
```

### 2. 数据持久化

确保数据库使用volume：

```yaml
volumes:
  - postgres-data:/var/lib/postgresql/data
```

### 3. 健康检查

所有服务都配置健康检查：

```yaml
healthcheck:
  test: ["CMD", "curl", "-f", "http://localhost:5000/health"]
  interval: 30s
  timeout: 10s
  retries: 3
```

### 4. 自动重启

```yaml
restart: unless-stopped
```

### 5. 日志管理

```yaml
logging:
  driver: "json-file"
  options:
    max-size: "10m"
    max-file: "3"
```

## 更新部署

### 更新模型

```bash
# 1. 替换模型文件
cp /path/to/new/model.mph model_service/models/

# 2. 重启模型服务
docker-compose restart model-service
```

### 更新Go后端代码

```bash
# 1. 停止服务
docker-compose stop backend

# 2. 重新构建
docker-compose build backend

# 3. 启动服务
docker-compose up -d backend
```

### 数据库迁移

```bash
# 添加新的迁移文件
# migrations/006_xxx.sql

# 运行迁移
cd migrations
./run_migrations.sh
```

## 备份和恢复

### 备份数据库

```bash
docker exec postgres pg_dump -U postgres model_mall > backup_$(date +%Y%m%d).sql
```

### 恢复数据库

```bash
docker exec -i postgres psql -U postgres model_mall < backup_20240101.sql
```

### 备份上传的图片

```bash
tar -czf uploads_backup_$(date +%Y%m%d).tar.gz uploads/
```

## 监控

### 使用Prometheus + Grafana

创建`monitoring/docker-compose.yml`：

```yaml
version: '3.8'

services:
  prometheus:
    image: prom/prometheus
    ports:
      - "9090:9090"
    volumes:
      - ./prometheus.yml:/etc/prometheus/prometheus.yml

  grafana:
    image: grafana/grafana
    ports:
      - "3000:3000"
```

## 卸载

```bash
# 停止并删除所有容器
docker-compose down

# 删除数据卷（注意：会删除数据库数据）
docker-compose down -v

# 删除镜像
docker-compose down --rmi all
```

## 参考文档

- [模型集成指南](MODEL_INTEGRATION_GUIDE.md)
- [模型服务README](model_service/README.md)
- [图片分类使用指南](IMAGE_CLASSIFICATION_GUIDE.md)
- [数据库设计](DATABASE_DESIGN.md)
