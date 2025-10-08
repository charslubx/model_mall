# Checkpoint格式模型使用指南

你的模型保存格式是这样的：

```python
state = {
    'epoch': epoch,
    'model_state': self.model.state_dict(),  # ← 模型权重
    'optimizer_state': self.optimizer.state_dict(),
    'scheduler_state': self.scheduler.state_dict(),
    'best_val_loss': self.best_val_loss
}
torch.save(state, 'checkpoint_best.pth')
```

这是一个**checkpoint格式**，包含训练状态但不包含模型架构，所以需要提供模型定义。

## 🚀 快速部署（3步）

### 步骤1: 准备模型架构定义

找到你训练时的模型类定义，复制到 `app/model_architecture.py`。

**示例：如果你的训练代码是这样的**

```python
# 你的训练代码 train.py
import torch.nn as nn

class MyImageClassifier(nn.Module):
    def __init__(self, num_classes=10):
        super().__init__()
        self.conv1 = nn.Conv2d(3, 64, 3, padding=1)
        self.conv2 = nn.Conv2d(64, 128, 3, padding=1)
        self.fc1 = nn.Linear(128 * 56 * 56, 256)
        self.fc2 = nn.Linear(256, num_classes)
    
    def forward(self, x):
        x = F.relu(self.conv1(x))
        x = F.max_pool2d(x, 2)
        x = F.relu(self.conv2(x))
        x = F.max_pool2d(x, 2)
        x = x.view(x.size(0), -1)
        x = F.relu(self.fc1(x))
        x = self.fc2(x)
        return x
```

**那么编辑 `app/model_architecture.py`：**

```python
# app/model_architecture.py
import torch.nn as nn
import torch.nn.functional as F

class MyImageClassifier(nn.Module):
    """复制你的模型定义到这里"""
    def __init__(self, num_classes=10):
        super().__init__()
        self.conv1 = nn.Conv2d(3, 64, 3, padding=1)
        self.conv2 = nn.Conv2d(64, 128, 3, padding=1)
        self.fc1 = nn.Linear(128 * 56 * 56, 256)
        self.fc2 = nn.Linear(256, num_classes)
    
    def forward(self, x):
        x = F.relu(self.conv1(x))
        x = F.max_pool2d(x, 2)
        x = F.relu(self.conv2(x))
        x = F.max_pool2d(x, 2)
        x = x.view(x.size(0), -1)
        x = F.relu(self.fc1(x))
        x = self.fc2(x)
        return x

# 添加创建函数
def create_my_model(num_classes=10):
    return MyImageClassifier(num_classes)

# 更新 create_model 函数
def create_model(model_name='resnet18', num_classes=10):
    if model_name == 'my_model':  # 添加你的模型
        return create_my_model(num_classes)
    # ... 其他模型
```

### 步骤2: 配置环境变量

编辑 `docker-compose.yml`：

```yaml
model-service:
  environment:
    - MODEL_PATH=/app/models/checkpoint_best.pth  # 你的checkpoint文件
    - MODEL_NAME=image-classifier
    - MODEL_ARCH=my_model        # ← 关键：指定模型架构
    - NUM_CLASSES=10             # ← 你的类别数量
```

### 步骤3: 部署

```bash
# 1. 复制checkpoint文件
cp checkpoint_best.pth /workspace/model_service/models/

# 2. 准备标签文件
cat > /workspace/model_service/models/labels.txt << EOF
类别1
类别2
类别3
EOF

# 3. 启动服务
cd /workspace
./start_all_services.sh
```

完成！🎉

## 📝 常见模型架构示例

### 使用ResNet

如果你训练时使用的是ResNet：

```python
# 训练时
import torchvision.models as models
model = models.resnet18(pretrained=True)
model.fc = nn.Linear(512, num_classes)
```

**配置：**
```yaml
environment:
  - MODEL_ARCH=resnet18
  - NUM_CLASSES=10
```

`app/model_architecture.py` 已经包含了这个定义，不需要修改！

### 使用EfficientNet

```yaml
environment:
  - MODEL_ARCH=efficientnet
  - NUM_CLASSES=10
```

### 使用MobileNet

```yaml
environment:
  - MODEL_ARCH=mobilenet
  - NUM_CLASSES=10
```

## 🔍 如何找到你的模型定义？

### 方法1: 查看训练代码

在你的训练脚本中搜索：
- `class xxx(nn.Module):`
- `def __init__`
- `def forward`

### 方法2: 查看模型配置文件

有些项目会把模型定义放在单独的文件：
- `models.py`
- `architecture.py`
- `network.py`

### 方法3: 检查checkpoint内容

```python
import torch

checkpoint = torch.load('checkpoint_best.pth')
print(checkpoint.keys())  # 查看包含的内容
# 输出: dict_keys(['epoch', 'model_state', 'optimizer_state', ...])

# 查看模型层的名称（可以推断架构）
state_dict = checkpoint['model_state']
for key in list(state_dict.keys())[:10]:
    print(key, state_dict[key].shape)
```

## 📊 完整示例：自定义CNN

### 1. 训练代码（你已有的）

```python
# train.py
class CustomImageClassifier(nn.Module):
    def __init__(self, num_classes=5):
        super().__init__()
        
        self.features = nn.Sequential(
            nn.Conv2d(3, 32, 3, padding=1),
            nn.ReLU(),
            nn.MaxPool2d(2),
            
            nn.Conv2d(32, 64, 3, padding=1),
            nn.ReLU(),
            nn.MaxPool2d(2),
            
            nn.Conv2d(64, 128, 3, padding=1),
            nn.ReLU(),
            nn.MaxPool2d(2),
        )
        
        self.classifier = nn.Sequential(
            nn.Linear(128 * 28 * 28, 512),
            nn.ReLU(),
            nn.Dropout(0.5),
            nn.Linear(512, num_classes)
        )
    
    def forward(self, x):
        x = self.features(x)
        x = x.view(x.size(0), -1)
        x = self.classifier(x)
        return x

# 训练和保存
model = CustomImageClassifier(num_classes=5)
# ... 训练过程 ...

# 保存checkpoint
state = {
    'epoch': epoch,
    'model_state': model.state_dict(),
    'optimizer_state': optimizer.state_dict(),
    'best_val_loss': best_loss
}
torch.save(state, 'checkpoint_best.pth')
```

### 2. 修改 model_architecture.py

```python
# app/model_architecture.py

# 将CustomImageClassifier类复制到这里
class CustomImageClassifier(nn.Module):
    def __init__(self, num_classes=5):
        super().__init__()
        
        self.features = nn.Sequential(
            nn.Conv2d(3, 32, 3, padding=1),
            nn.ReLU(),
            nn.MaxPool2d(2),
            
            nn.Conv2d(32, 64, 3, padding=1),
            nn.ReLU(),
            nn.MaxPool2d(2),
            
            nn.Conv2d(64, 128, 3, padding=1),
            nn.ReLU(),
            nn.MaxPool2d(2),
        )
        
        self.classifier = nn.Sequential(
            nn.Linear(128 * 28 * 28, 512),
            nn.ReLU(),
            nn.Dropout(0.5),
            nn.Linear(512, num_classes)
        )
    
    def forward(self, x):
        x = self.features(x)
        x = x.view(x.size(0), -1)
        x = self.classifier(x)
        return x

# 在 create_model 函数中添加
def create_model(model_name='resnet18', num_classes=10):
    # ... 其他模型 ...
    
    elif model_name == 'custom_image_classifier':
        return CustomImageClassifier(num_classes)
    
    # ...
```

### 3. 配置 docker-compose.yml

```yaml
environment:
  - MODEL_PATH=/app/models/checkpoint_best.pth
  - MODEL_ARCH=custom_image_classifier  # ← 使用你定义的名称
  - NUM_CLASSES=5                       # ← 你的类别数
```

### 4. 部署

```bash
./start_all_services.sh
```

## ⚠️ 注意事项

### 1. 模型架构必须完全一致

- 层的数量、类型、参数必须与训练时完全相同
- 变量名必须一致
- `__init__` 和 `forward` 的逻辑必须一致

### 2. 输入尺寸要匹配

如果你的模型期望特定输入尺寸（如 `224x224`），确保在 `model_loader.py` 中的 `preprocess_image()` 方法使用相同的尺寸。

### 3. 类别数量要正确

`NUM_CLASSES` 必须与训练时的类别数一致。

### 4. 标签顺序要对应

`labels.txt` 中的标签顺序必须与训练时的类别索引对应。

## 🧪 测试

部署后测试：

```bash
# 1. 健康检查
curl http://localhost:5000/health

# 2. 查看模型信息
curl http://localhost:5000/info

# 3. 测试分类
curl -X POST http://localhost:5000/classify \
  -F "image=@test.jpg"
```

## 🐛 故障排查

### 问题1: "RuntimeError: Error(s) in loading state_dict"

**原因：** 模型架构定义与checkpoint不匹配

**解决：**
1. 检查模型类定义是否与训练时完全一致
2. 检查类别数量是否正确
3. 运行下面的代码对比：

```python
import torch

# 加载checkpoint
checkpoint = torch.load('checkpoint_best.pth')
state_dict = checkpoint['model_state']

# 查看所有层的名称和形状
for name, param in state_dict.items():
    print(f"{name}: {param.shape}")

# 对比你定义的模型
from app.model_architecture import CustomImageClassifier
model = CustomImageClassifier(num_classes=5)
for name, param in model.state_dict().items():
    print(f"{name}: {param.shape}")
```

### 问题2: 模型加载失败

查看日志：
```bash
docker-compose logs model-service
```

### 问题3: 预测结果不准确

1. 检查标签文件顺序
2. 检查图片预处理是否与训练时一致
3. 确认使用了最佳checkpoint（`checkpoint_best.pth`）

## 📚 相关文档

- [PyTorch模型指南](PYTORCH_MODEL_GUIDE.md)
- [快速开始](../PYTORCH_QUICKSTART.md)
- [部署指南](../DEPLOYMENT.md)

## 💡 小贴士

1. **保留训练代码** - 方便随时查看模型定义
2. **使用有意义的模型名称** - 如 `my_resnet_classifier` 而不是 `custom_cnn`
3. **添加注释** - 在 `model_architecture.py` 中注释模型的用途和特点
4. **版本管理** - 如果有多个模型版本，使用不同的函数名

祝部署顺利！🎉
