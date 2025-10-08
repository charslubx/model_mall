# PyTorch模型快速开始指南

本指南专门针对使用PyTorch训练的模型（.mph, .pt, .pth格式）。

## 🚀 5分钟快速部署

### 步骤1: 测试你的模型文件

首先测试你的`.mph`文件是否可以正常加载：

```bash
cd /workspace/model_service

# 安装PyTorch（如果本地没有）
pip install torch torchvision

# 测试模型文件
python test_pytorch_model.py /path/to/your/model.mph
```

**输出示例：**
```
==============================================================
PyTorch模型测试工具
==============================================================

📁 模型路径: model.mph
📦 文件大小: 45.23 MB

🖥️  设备信息:
  ✓ CUDA可用
  ✓ GPU: NVIDIA GeForce RTX 3080
  ✓ PyTorch版本: 2.1.0

🔄 正在加载模型...
✅ 成功加载模型
   类型: ResNet

📊 模型信息:
  ✓ 总参数: 11,689,512
  ✓ 类别数量: 10

✅ 模型测试通过！
```

### 步骤2: 根据测试结果部署

#### 情况A: 测试通过（完整模型）✅

如果测试脚本显示"模型测试通过"，直接部署即可：

```bash
# 1. 复制模型文件
cp /path/to/your/model.mph /workspace/model_service/models/

# 2. 创建标签文件
cat > /workspace/model_service/models/labels.txt << EOF
类别1
类别2
类别3
...
EOF

# 3. 更新配置（如果文件名不是model.mph）
# 编辑 docker-compose.yml:
# MODEL_PATH=/app/models/your-model-name.mph

# 4. 启动服务
cd /workspace
./start_all_services.sh

# 5. 测试
curl -X POST http://localhost:8888/api/images/upload \
  -F "image=@test.jpg"
```

完成！🎉

#### 情况B: 需要模型架构（state_dict）⚠️

如果测试脚本提示"这是权重文件，需要先定义模型架构"：

**步骤1: 创建模型架构文件**

根据你的训练代码创建 `/workspace/model_service/app/model_architecture.py`：

```python
import torch
import torch.nn as nn
import torchvision.models as models

# 示例1: 如果你使用的是ResNet
class MyClassifier(nn.Module):
    def __init__(self, num_classes=10):
        super().__init__()
        self.model = models.resnet18(pretrained=False)
        self.model.fc = nn.Linear(512, num_classes)
    
    def forward(self, x):
        return self.model(x)

# 示例2: 如果是自定义CNN
class CustomCNN(nn.Module):
    def __init__(self, num_classes=10):
        super().__init__()
        # 复制你训练时的网络结构
        self.conv1 = nn.Conv2d(3, 64, 3, padding=1)
        self.conv2 = nn.Conv2d(64, 128, 3, padding=1)
        # ... 其他层
    
    def forward(self, x):
        # 复制你训练时的前向传播逻辑
        x = F.relu(self.conv1(x))
        # ...
        return x
```

**步骤2: 修改API加载代码**

编辑 `/workspace/model_service/app/api.py`，修改 `init_model()` 函数：

```python
def init_model():
    """初始化模型"""
    global model_loader
    
    model_path = os.environ.get('MODEL_PATH', '/app/models/model.mph')
    model_name = os.environ.get('MODEL_NAME', 'image-classifier')
    
    try:
        from app.model_loader import PyTorchModelLoader
        from app.model_architecture import MyClassifier  # 导入你定义的模型
        
        # 创建加载器
        model_loader = PyTorchModelLoader(model_path, model_name)
        
        # 使用架构加载权重
        model_class = MyClassifier(num_classes=10)  # 修改为实际类别数
        model_loader.load_model_with_architecture(model_class)
        
        logger.info("模型初始化成功")
    except Exception as e:
        logger.error(f"模型初始化失败: {str(e)}")
        raise
```

**步骤3: 部署**

```bash
cd /workspace
./start_all_services.sh
```

## 📋 标签文件格式

`models/labels.txt` 文件格式：

```
类别0的名称
类别1的名称
类别2的名称
...
```

**重要：**
- 每行一个标签
- 按照模型输出的类别索引顺序排列
- 标签数量必须与模型输出的类别数一致

**示例：**

```bash
# 10类图片分类
cat > models/labels.txt << EOF
飞机
汽车
鸟
猫
鹿
狗
青蛙
马
船
卡车
EOF
```

## 🐳 Docker部署配置

### 基础配置

`model_service/docker-compose.yml`:

```yaml
version: '3.8'

services:
  model-service:
    build: .
    ports:
      - "5000:5000"
    volumes:
      - ./models:/app/models
    environment:
      - MODEL_PATH=/app/models/model.mph
      - MODEL_NAME=image-classifier
      - DEVICE=cuda  # 或 cpu
```

### GPU加速配置

如果有GPU，修改 `docker-compose.yml`:

```yaml
model-service:
  build: .
  deploy:
    resources:
      reservations:
        devices:
          - driver: nvidia
            count: 1
            capabilities: [gpu]
  environment:
    - MODEL_PATH=/app/models/model.mph
    - DEVICE=cuda
```

并修改 `Dockerfile`:

```dockerfile
FROM pytorch/pytorch:2.1.0-cuda11.8-cudnn8-runtime

# ... 其他配置
```

## 🧪 测试和验证

### 1. 测试模型服务

```bash
# 健康检查
curl http://localhost:5000/health

# 预期响应
{
  "status": "healthy",
  "model_name": "image-classifier",
  "model_loaded": true
}
```

### 2. 测试图片分类

```bash
# 上传测试图片
curl -X POST http://localhost:5000/classify \
  -F "image=@test.jpg"

# 预期响应
{
  "success": true,
  "results": [
    {"label": "猫", "confidence": 0.8523},
    {"label": "狗", "confidence": 0.1234}
  ]
}
```

### 3. 测试完整流程

```bash
# 通过Go后端上传
curl -X POST http://localhost:8888/api/images/upload \
  -F "image=@test.jpg"

# 预期响应
{
  "image_id": 1,
  "filename": "test.jpg",
  "classifications": [
    {"label": "猫", "confidence": 0.8523}
  ]
}
```

## 🔧 常见问题

### Q1: 模型加载失败

**错误：** `RuntimeError: Error(s) in loading state_dict`

**解决：**
1. 检查模型架构定义是否与训练时一致
2. 确认类别数量正确
3. 查看详细错误日志

### Q2: CUDA out of memory

**错误：** GPU内存不足

**解决：**
```yaml
# 使用CPU
environment:
  - DEVICE=cpu

# 或减少worker数量
command: gunicorn --workers 2 ...
```

### Q3: 预测结果不准确

**可能原因：**
1. 标签顺序错误
2. 图片预处理不一致
3. 模型未处于eval模式

**解决：**
1. 检查 labels.txt 顺序
2. 确保预处理与训练时一致
3. 确认调用了 `model.eval()`

### Q4: 找不到模型文件

**错误：** `FileNotFoundError: 模型文件不存在`

**解决：**
```bash
# 检查文件是否存在
ls -lh model_service/models/

# 检查Docker volume挂载
docker-compose config | grep volumes -A 5

# 检查环境变量
docker-compose exec model-service env | grep MODEL
```

## 📊 性能优化

### 1. 批处理推理

适用于需要处理大量图片的场景。

### 2. 模型量化

减小模型大小，提高推理速度：

```python
# 动态量化
quantized_model = torch.quantization.quantize_dynamic(
    model, {torch.nn.Linear}, dtype=torch.qint8
)
```

### 3. TorchScript

将模型转换为TorchScript格式：

```python
# 训练后
scripted_model = torch.jit.script(model)
torch.jit.save(scripted_model, 'model_scripted.pt')

# 加载时
model = torch.jit.load('model_scripted.pt')
```

## 📚 相关文档

- [PyTorch模型详细指南](model_service/PYTORCH_MODEL_GUIDE.md)
- [模型集成完整指南](MODEL_INTEGRATION_GUIDE.md)
- [部署指南](DEPLOYMENT.md)
- [图片分类使用指南](IMAGE_CLASSIFICATION_GUIDE.md)

## 💡 下一步

部署成功后，你可以：

1. **集成到前端** - 创建上传界面
2. **批量处理** - 添加批量图片分类功能
3. **模型管理** - 支持多个模型切换
4. **性能监控** - 添加Prometheus监控
5. **结果缓存** - 使用Redis缓存结果

祝部署顺利！🚀
