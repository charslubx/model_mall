# 🚀 针对你的Checkpoint模型的快速部署指南

你的模型是**checkpoint格式**（只包含权重，不包含架构），需要提供模型定义。

## ⚡ 3步完成部署

### 步骤1️⃣: 找到并复制你的模型定义

在你的训练代码中找到模型类定义，例如：

```python
# 你的 train.py 或 model.py
class MyModel(nn.Module):
    def __init__(self, num_classes=10):
        super().__init__()
        # ... 你的模型层定义
    
    def forward(self, x):
        # ... 前向传播逻辑
        return x
```

**将这个类复制到** `/workspace/model_service/app/model_architecture.py`

具体操作：

```bash
# 1. 编辑文件
nano /workspace/model_service/app/model_architecture.py

# 2. 在文件末尾添加你的模型类

# 3. 在 create_model 函数中添加选项：
def create_model(model_name='resnet18', num_classes=10):
    # ... 已有的代码 ...
    
    elif model_name == 'my_model':  # ← 添加这里
        return MyModel(num_classes)
```

### 步骤2️⃣: 配置服务

编辑 `/workspace/model_service/docker-compose.yml`：

```yaml
environment:
  - MODEL_PATH=/app/models/checkpoint_best.pth  # ← 你的文件名
  - MODEL_ARCH=my_model     # ← 你在步骤1中添加的模型名
  - NUM_CLASSES=10          # ← 你的类别数量
```

### 步骤3️⃣: 部署

```bash
# 1. 复制checkpoint文件
cp /path/to/your/checkpoint_best.pth /workspace/model_service/models/

# 2. 创建标签文件
cat > /workspace/model_service/models/labels.txt << EOF
类别1
类别2
类别3
...（共NUM_CLASSES个类别）
EOF

# 3. 启动服务
cd /workspace
./start_all_services.sh
```

## ✅ 验证

```bash
# 测试模型服务
curl http://localhost:5000/health

# 上传图片测试
curl -X POST http://localhost:8888/api/images/upload \
  -F "image=@test.jpg"
```

## 📋 常见情况

### 情况A: 使用的是ResNet/EfficientNet/MobileNet

如果你训练时用的是这些预训练模型，更简单：

```yaml
# docker-compose.yml
environment:
  - MODEL_PATH=/app/models/checkpoint_best.pth
  - MODEL_ARCH=resnet18        # 或 efficientnet, mobilenet
  - NUM_CLASSES=10
```

**不需要修改 model_architecture.py！** 这些模型已经内置了。

### 情况B: 自定义CNN模型

需要复制你的模型类到 `model_architecture.py`。

## 🆘 需要帮助？

查看详细文档：

- **[CHECKPOINT_GUIDE.md](model_service/CHECKPOINT_GUIDE.md)** - Checkpoint格式完整指南
- **[PYTORCH_MODEL_GUIDE.md](model_service/PYTORCH_MODEL_GUIDE.md)** - PyTorch模型详细说明

## 🎯 快速模板

### 模板：如果你不确定模型类名

在 `model_architecture.py` 中添加：

```python
# 将你的模型定义复制到这里
class YourModelName(nn.Module):  # ← 改成你的模型类名
    def __init__(self, num_classes=10):
        super().__init__()
        # 复制你的 __init__ 内容
        pass
    
    def forward(self, x):
        # 复制你的 forward 内容
        pass

# 添加到 create_model
def create_model(model_name='resnet18', num_classes=10):
    # ...
    elif model_name == 'your_model':  # ← 改成你想用的名字
        return YourModelName(num_classes)
```

然后在 `docker-compose.yml` 中：

```yaml
environment:
  - MODEL_ARCH=your_model  # ← 使用你定义的名字
```

就这么简单！🎉
