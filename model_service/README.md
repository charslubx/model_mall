# 图片分类模型服务

这是一个基于Flask的图片分类模型服务，支持多种模型格式，提供HTTP API接口供Go后端调用。

## 📋 功能特性

- ✅ **支持PyTorch模型**：`.pt`、`.pth`、`.mph`（推荐）
- ✅ 支持多种模型格式：Keras (`.h5`)、ONNX (`.onnx`)
- ✅ 自动识别并加载对应格式的模型
- ✅ GPU加速支持（CUDA）
- ✅ RESTful API接口
- ✅ Docker容器化部署
- ✅ 健康检查和监控
- ✅ 生产级别的Gunicorn服务器

## 🗂️ 项目结构

```
model_service/
├── app/                        # 应用代码
│   ├── __init__.py
│   ├── api.py                 # Flask API服务
│   └── model_loader.py        # 模型加载器
├── models/                     # 模型文件目录
│   ├── .gitkeep
│   ├── labels.txt             # 分类标签文件
│   └── model.h5               # 你的模型文件（需要自行添加）
├── tests/                      # 测试文件
├── Dockerfile                  # Docker镜像定义
├── docker-compose.yml          # Docker Compose配置
├── requirements.txt            # Python依赖
├── .env.example               # 环境变量示例
├── start_service.sh           # 本地启动脚本
├── start_with_docker.sh       # Docker启动脚本
├── test_service.py            # 服务测试脚本
└── README.md                  # 本文件
```

## 🚀 快速开始

### PyTorch完整模型（推荐方式）

```bash
# 1. 复制模型文件
cp model.pth models/

# 2. 启动服务
./start_with_docker.sh
```

**就这么简单！**

### 配置文件

编辑`docker-compose.yml`，修改环境变量：

```yaml
environment:
  - MODEL_PATH=/app/models/model.mph  # 改为你的模型文件名
  - MODEL_NAME=image-classifier
```

3. **启动服务**

```bash
./start_with_docker.sh
```

或手动启动：

```bash
docker-compose up -d
```

4. **查看日志**

```bash
docker-compose logs -f
```

### 方式2：本地运行

1. **安装依赖**

```bash
# 创建虚拟环境
python3 -m venv venv
source venv/bin/activate  # Linux/Mac
# venv\Scripts\activate  # Windows

# 安装依赖
pip install -r requirements.txt
```

2. **准备模型文件**

```bash
cp /path/to/your/model.mph models/
```

3. **设置环境变量**

```bash
export MODEL_PATH=models/model.mph
export MODEL_NAME=image-classifier
export PORT=5000
```

4. **启动服务**

```bash
# 使用启动脚本
./start_service.sh

# 或直接运行
python -m app.api
```

## 🔧 PyTorch模型支持

### 完整模型格式

本服务支持PyTorch完整模型：

```python
# 训练时保存
torch.save(model, 'model.pth')

# 部署时直接加载，无需额外配置
```

**简单、快速、无需配置！**

### 常见模型格式转换

#### 如果是 Keras 模型（.h5）

直接使用即可，无需修改。

#### 如果是 PyTorch 模型（.pt/.pth）

1. 修改`requirements.txt`，取消注释：
```
torch==2.1.0
torchvision==0.16.0
```

2. 修改`app/model_loader.py`，添加PyTorch加载器。

#### 如果是自定义格式

修改`app/model_loader.py`中的`CustomModelLoader.load_model()`方法：

```python
def load_model(self):
    """加载自定义格式模型"""
    try:
        # 根据你的模型格式实现加载逻辑
        import your_model_library
        
        logger.info(f"正在加载自定义模型: {self.model_path}")
        self.model = your_model_library.load_model(self.model_path)
        self.is_loaded = True
        logger.info("模型加载成功")
    except Exception as e:
        logger.error(f"加载模型失败: {str(e)}")
        raise
```

## 📡 API 接口

### 1. 健康检查

```bash
GET /health
```

**响应：**
```json
{
  "status": "healthy",
  "model_name": "image-classifier",
  "model_loaded": true
}
```

### 2. 图片分类

```bash
POST /classify
Content-Type: multipart/form-data
```

**请求：**
- 字段名：`image`
- 类型：文件
- 支持格式：jpg, jpeg, png, gif, bmp, webp

**示例：**
```bash
curl -X POST http://localhost:5000/classify \
  -F "image=@/path/to/your/image.jpg"
```

**响应：**
```json
{
  "success": true,
  "results": [
    {
      "label": "cat",
      "confidence": 0.8523
    },
    {
      "label": "dog",
      "confidence": 0.1234
    },
    {
      "label": "bird",
      "confidence": 0.0243
    }
  ],
  "model_name": "image-classifier"
}
```

### 3. 模型信息

```bash
GET /info
```

**响应：**
```json
{
  "model_name": "image-classifier",
  "model_path": "/app/models/model.h5",
  "is_loaded": true,
  "model_type": "KerasModelLoader"
}
```

## 🧪 测试服务

使用测试脚本：

```bash
# 测试健康检查和模型信息
python test_service.py

# 测试分类功能
python test_service.py /path/to/test/image.jpg
```

或使用curl：

```bash
# 健康检查
curl http://localhost:5000/health

# 模型信息
curl http://localhost:5000/info

# 图片分类
curl -X POST http://localhost:5000/classify \
  -F "image=@test_image.jpg"
```

## 🏷️ 配置分类标签

编辑`models/labels.txt`文件，每行一个标签，按照模型输出的类别索引顺序排列：

```
cat
dog
bird
horse
cow
sheep
...
```

## ⚙️ 配置说明

### 环境变量

| 变量名 | 说明 | 默认值 |
|--------|------|--------|
| MODEL_PATH | 模型文件路径 | /app/models/model.h5 |
| MODEL_NAME | 模型名称 | image-classifier |
| PORT | 服务端口 | 5000 |
| HOST | 监听地址 | 0.0.0.0 |
| LOG_LEVEL | 日志级别 | INFO |

### 修改配置

**Docker方式：** 编辑`docker-compose.yml`

**本地方式：** 设置环境变量或编辑`.env`文件

## 🔗 与 Go 后端集成

### 1. 更新 Go 后端配置

编辑`backend/etc/backend-api.yaml`：

```yaml
Model:
  Name: image-classifier
  Version: v1.0
  Type: remote
  Path: http://localhost:5000/classify  # 或 http://model-service:5000/classify
```

### 2. 如果使用 Docker Compose 统一部署

创建一个总的`docker-compose.yml`将两个服务连接起来：

```yaml
version: '3.8'

services:
  model-service:
    build: ./model_service
    container_name: model-service
    volumes:
      - ./model_service/models:/app/models
    environment:
      - MODEL_PATH=/app/models/model.mph
      - MODEL_NAME=image-classifier
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
      - model-service
    networks:
      - app-network

networks:
  app-network:
    driver: bridge
```

## 🐛 故障排查

### 模型加载失败

**问题：** 服务启动后模型加载失败

**解决方案：**
1. 检查模型文件路径是否正确
2. 查看日志了解具体错误：`docker-compose logs -f`
3. 确认模型文件格式是否支持
4. 如果是自定义格式，需要修改`CustomModelLoader`

### 内存不足

**问题：** 服务运行时内存占用过高

**解决方案：**
1. 减少`docker-compose.yml`中的worker数量
2. 使用更小的模型
3. 增加服务器内存

### 推理速度慢

**问题：** 图片分类速度慢

**解决方案：**
1. 使用GPU加速（修改Dockerfile使用GPU版本TensorFlow）
2. 优化模型（量化、剪枝）
3. 增加worker数量
4. 使用批处理

## 📈 性能优化

### 1. 使用 GPU 加速

修改`Dockerfile`使用GPU版本的深度学习框架：

```dockerfile
FROM tensorflow/tensorflow:2.15.0-gpu
```

需要安装NVIDIA Docker运行时。

### 2. 批处理推理

修改API支持批量图片上传，减少网络开销。

### 3. 模型优化

- 模型量化（INT8量化）
- 模型剪枝
- 转换为ONNX格式使用ONNX Runtime

### 4. 添加缓存

使用Redis缓存常见图片的分类结果。

## 📝 开发建议

### 添加新的模型类型支持

1. 在`app/model_loader.py`中创建新的Loader类
2. 继承`ModelLoader`基类
3. 实现`load_model()`和`predict()`方法
4. 在`create_model_loader()`函数中添加格式识别逻辑

### 添加图片预处理

修改`ModelLoader.preprocess_image()`方法，添加自定义的预处理逻辑。

### 添加监控

使用Prometheus + Grafana监控服务：
- 请求数量
- 响应时间
- 错误率
- CPU/内存使用率

## 🔒 安全建议

1. **添加认证** - 使用API Key或JWT Token
2. **限流** - 防止API滥用
3. **输入验证** - 验证上传文件的安全性
4. **HTTPS** - 生产环境使用HTTPS
5. **网络隔离** - 使用Docker网络隔离

## 📚 相关文档

- [Flask文档](https://flask.palletsprojects.com/)
- [TensorFlow文档](https://www.tensorflow.org/)
- [Docker文档](https://docs.docker.com/)
- [Go后端集成文档](../IMAGE_CLASSIFICATION_GUIDE.md)

## 🤝 贡献

如果你有改进建议或发现bug，欢迎提交Issue或Pull Request。

## 📄 许可证

MIT License
