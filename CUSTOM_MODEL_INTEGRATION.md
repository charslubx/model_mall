# 🔧 自定义模型集成指南

你的模型是**自己改过的**，不是标准的ResNet/EfficientNet等。这完全没问题！

## 🎯 你需要做的事情（只需3步）

### 步骤1️⃣：找到你的模型定义代码

在你的训练代码中找到模型类定义。通常在这些文件里：
- `model.py`
- `models.py`
- `network.py`
- `train.py`（有时模型直接定义在训练脚本里）

**找到类似这样的代码：**

```python
import torch
import torch.nn as nn
import torch.nn.functional as F

class YourCustomModel(nn.Module):  # ← 这个类名可能不同
    def __init__(self, num_classes=10):
        super(YourCustomModel, self).__init__()
        
        # 你的模型层定义
        self.conv1 = ...
        self.fc1 = ...
        # ... 等等
    
    def forward(self, x):
        # 你的前向传播逻辑
        x = ...
        return x
```

### 步骤2️⃣：复制到我们的项目

**打开文件：** `/workspace/model_service/app/model_architecture.py`

**在文件末尾添加你的模型类：**

```python
# ============================================
# 你的自定义模型（从训练代码复制）
# ============================================

class YourCustomModel(nn.Module):
    """
    你的自定义模型
    从训练代码复制而来
    """
    def __init__(self, num_classes=10):
        super(YourCustomModel, self).__init__()
        
        # 完整复制你的 __init__ 方法内容
        # ...你的层定义...
    
    def forward(self, x):
        # 完整复制你的 forward 方法内容
        # ...你的前向传播逻辑...
        return x


# 添加到 create_model 函数
def create_model(model_name='resnet18', num_classes=10):
    model_name = model_name.lower()
    
    # ... 已有的代码 ...
    
    # 添加你的模型
    elif model_name == 'my_custom_model':  # ← 给个名字
        return YourCustomModel(num_classes)
    
    else:
        raise ValueError(f"未知的模型名称: {model_name}")
```

### 步骤3️⃣：配置并启动

**编辑：** `/workspace/model_service/docker-compose.yml`

```yaml
environment:
  - MODEL_PATH=/app/models/checkpoint_best.pth
  - MODEL_ARCH=my_custom_model  # ← 使用你在步骤2中定义的名字
  - NUM_CLASSES=10              # ← 你的类别数量
```

**然后启动：**

```bash
# 1. 复制checkpoint
cp checkpoint_best.pth /workspace/model_service/models/

# 2. 创建标签文件
cat > /workspace/model_service/models/labels.txt << EOF
类别1
类别2
...
EOF

# 3. 启动服务
cd /workspace
./start_all_services.sh
```

完成！🎉

## 📋 完整示例

### 示例：你的训练代码

假设你的训练代码是这样的（train.py）：

```python
import torch
import torch.nn as nn
import torch.nn.functional as F

class MyImageClassifier(nn.Module):
    """我改过的分类器"""
    def __init__(self, num_classes=10):
        super(MyImageClassifier, self).__init__()
        
        # 第一个卷积块
        self.conv1 = nn.Conv2d(3, 64, kernel_size=3, padding=1)
        self.bn1 = nn.BatchNorm2d(64)
        
        # 第二个卷积块
        self.conv2 = nn.Conv2d(64, 128, kernel_size=3, padding=1)
        self.bn2 = nn.BatchNorm2d(128)
        
        # 第三个卷积块
        self.conv3 = nn.Conv2d(128, 256, kernel_size=3, padding=1)
        self.bn3 = nn.BatchNorm2d(256)
        
        # 第四个卷积块（你可能加了更多层）
        self.conv4 = nn.Conv2d(256, 512, kernel_size=3, padding=1)
        self.bn4 = nn.BatchNorm2d(512)
        
        # 全连接层
        self.fc1 = nn.Linear(512 * 14 * 14, 1024)
        self.dropout1 = nn.Dropout(0.5)
        self.fc2 = nn.Linear(1024, 512)
        self.dropout2 = nn.Dropout(0.3)
        self.fc3 = nn.Linear(512, num_classes)
    
    def forward(self, x):
        # 第一个卷积块
        x = self.conv1(x)
        x = self.bn1(x)
        x = F.relu(x)
        x = F.max_pool2d(x, 2)
        
        # 第二个卷积块
        x = self.conv2(x)
        x = self.bn2(x)
        x = F.relu(x)
        x = F.max_pool2d(x, 2)
        
        # 第三个卷积块
        x = self.conv3(x)
        x = self.bn3(x)
        x = F.relu(x)
        x = F.max_pool2d(x, 2)
        
        # 第四个卷积块
        x = self.conv4(x)
        x = self.bn4(x)
        x = F.relu(x)
        x = F.max_pool2d(x, 2)
        
        # 展平
        x = x.view(x.size(0), -1)
        
        # 全连接层
        x = self.fc1(x)
        x = F.relu(x)
        x = self.dropout1(x)
        
        x = self.fc2(x)
        x = F.relu(x)
        x = self.dropout2(x)
        
        x = self.fc3(x)
        
        return x

# 训练代码
model = MyImageClassifier(num_classes=10)
# ... 训练 ...

# 保存checkpoint
state = {
    'epoch': epoch,
    'model_state': model.state_dict(),
    'optimizer_state': optimizer.state_dict(),
    'best_val_loss': best_val_loss
}
torch.save(state, 'checkpoint_best.pth')
```

### 集成到项目

**编辑：** `/workspace/model_service/app/model_architecture.py`

在文件末尾添加：

```python
# ============================================
# MyImageClassifier - 我的自定义模型
# ============================================

class MyImageClassifier(nn.Module):
    """我改过的分类器 - 从训练代码完整复制"""
    def __init__(self, num_classes=10):
        super(MyImageClassifier, self).__init__()
        
        # 第一个卷积块
        self.conv1 = nn.Conv2d(3, 64, kernel_size=3, padding=1)
        self.bn1 = nn.BatchNorm2d(64)
        
        # 第二个卷积块
        self.conv2 = nn.Conv2d(64, 128, kernel_size=3, padding=1)
        self.bn2 = nn.BatchNorm2d(128)
        
        # 第三个卷积块
        self.conv3 = nn.Conv2d(128, 256, kernel_size=3, padding=1)
        self.bn3 = nn.BatchNorm2d(256)
        
        # 第四个卷积块
        self.conv4 = nn.Conv2d(256, 512, kernel_size=3, padding=1)
        self.bn4 = nn.BatchNorm2d(512)
        
        # 全连接层
        self.fc1 = nn.Linear(512 * 14 * 14, 1024)
        self.dropout1 = nn.Dropout(0.5)
        self.fc2 = nn.Linear(1024, 512)
        self.dropout2 = nn.Dropout(0.3)
        self.fc3 = nn.Linear(512, num_classes)
    
    def forward(self, x):
        # 第一个卷积块
        x = self.conv1(x)
        x = self.bn1(x)
        x = F.relu(x)
        x = F.max_pool2d(x, 2)
        
        # 第二个卷积块
        x = self.conv2(x)
        x = self.bn2(x)
        x = F.relu(x)
        x = F.max_pool2d(x, 2)
        
        # 第三个卷积块
        x = self.conv3(x)
        x = self.bn3(x)
        x = F.relu(x)
        x = F.max_pool2d(x, 2)
        
        # 第四个卷积块
        x = self.conv4(x)
        x = self.bn4(x)
        x = F.relu(x)
        x = F.max_pool2d(x, 2)
        
        # 展平
        x = x.view(x.size(0), -1)
        
        # 全连接层
        x = self.fc1(x)
        x = F.relu(x)
        x = self.dropout1(x)
        
        x = self.fc2(x)
        x = F.relu(x)
        x = self.dropout2(x)
        
        x = self.fc3(x)
        
        return x


# 在 create_model 函数中添加
def create_model(model_name='resnet18', num_classes=10):
    model_name = model_name.lower()
    
    if model_name == 'resnet18':
        return create_resnet_model(num_classes)
    elif model_name == 'efficientnet':
        return create_efficientnet_model(num_classes)
    elif model_name == 'mobilenet':
        return create_mobilenet_model(num_classes)
    elif model_name == 'simple_cnn':
        return SimpleCNN(num_classes)
    elif model_name == 'my_classifier':  # ← 添加这个
        return MyImageClassifier(num_classes)
    else:
        raise ValueError(f"未知的模型名称: {model_name}")
```

**配置文件：** `/workspace/model_service/docker-compose.yml`

```yaml
environment:
  - MODEL_PATH=/app/models/checkpoint_best.pth
  - MODEL_ARCH=my_classifier  # ← 使用上面定义的名字
  - NUM_CLASSES=10
```

## ⚠️ 重要注意事项

### 1. 必须完全一致

**关键：** 复制的模型定义必须与训练时**100%一致**

- ✅ 层的数量和类型
- ✅ 参数（kernel_size, padding等）
- ✅ 变量名（self.conv1, self.fc1等）
- ✅ forward方法的逻辑

### 2. import语句

确保文件开头有必要的导入：

```python
# model_architecture.py 文件开头
import torch
import torch.nn as nn
import torch.nn.functional as F
import torchvision.models as models
```

### 3. 类别数量

`num_classes` 参数必须与训练时一致。

### 4. 输入尺寸

如果你的模型期望特定的输入尺寸（如224x224），可能需要调整预处理。

## 🔍 如何验证复制正确

运行这个测试脚本：

```python
# test_model_match.py
import torch
from app.model_architecture import MyImageClassifier

# 1. 加载checkpoint
checkpoint = torch.load('checkpoint_best.pth')
saved_state = checkpoint['model_state']

# 2. 创建模型
model = MyImageClassifier(num_classes=10)
current_state = model.state_dict()

# 3. 对比所有层
print("对比模型结构：")
print(f"Checkpoint中的层数: {len(saved_state)}")
print(f"当前模型的层数: {len(current_state)}")

# 4. 逐层检查
all_match = True
for name in saved_state.keys():
    if name not in current_state:
        print(f"❌ 缺少层: {name}")
        all_match = False
    elif saved_state[name].shape != current_state[name].shape:
        print(f"❌ 形状不匹配: {name}")
        print(f"   Checkpoint: {saved_state[name].shape}")
        print(f"   当前模型: {current_state[name].shape}")
        all_match = False

if all_match:
    print("✅ 所有层都匹配！可以加载")
else:
    print("❌ 存在不匹配，请检查模型定义")
```

## 🐛 常见问题

### 问题1：找不到训练代码中的模型定义

**方法1：** 搜索文件

```bash
cd /path/to/your/training/code
grep -r "class.*nn.Module" .
grep -r "def __init__" .
```

**方法2：** 查看checkpoint内容

```python
import torch
checkpoint = torch.load('checkpoint_best.pth')

# 查看所有层的名称
for name in checkpoint['model_state'].keys():
    print(name)

# 这能帮你推断模型结构
```

### 问题2：模型用了其他文件的组件

如果你的模型用了其他文件定义的模块，也要一起复制：

```python
# 如果你的模型用了自定义的AttentionModule
from attention import AttentionModule  # 训练时

# 那么需要把AttentionModule也复制到model_architecture.py
class AttentionModule(nn.Module):
    # ... 复制完整定义
    pass

class MyCustomModel(nn.Module):
    def __init__(self, num_classes=10):
        super().__init__()
        self.attention = AttentionModule()  # 现在可以用了
        # ...
```

### 问题3：模型定义跨多个文件

将所有需要的类都复制到 `model_architecture.py` 中，确保依赖关系正确。

## 📝 检查清单

部署前检查：

- [ ] 已完整复制模型类定义到 `model_architecture.py`
- [ ] 已在 `create_model()` 函数中添加选项
- [ ] `docker-compose.yml` 中的 `MODEL_ARCH` 与函数中的名字匹配
- [ ] `NUM_CLASSES` 与训练时一致
- [ ] 所有必要的import语句都有
- [ ] 如果有自定义模块，都已复制
- [ ] （可选）运行了验证脚本确认结构匹配

## 🚀 完成后

```bash
cd /workspace
./start_all_services.sh
```

检查日志：
```bash
docker-compose logs -f model-service
```

应该看到：
```
✓ 使用自定义架构加载checkpoint
✓ 从checkpoint加载权重 (epoch: X)
✓ PyTorch模型加载成功
✓ 训练轮次: X
✓ 最佳验证损失: X.XXXX
```

## 💡 额外提示

### 技巧1：保持训练代码可访问

部署后仍然保留训练代码，方便以后参考。

### 技巧2：添加注释

在 `model_architecture.py` 中添加详细注释：

```python
class MyCustomModel(nn.Module):
    """
    自定义图片分类模型
    
    架构特点：
    - 4个卷积块，每块带BatchNorm
    - 3层全连接，使用Dropout防止过拟合
    - 输入尺寸: 224x224
    - 类别数: 10
    
    训练信息：
    - 数据集: XXX
    - 最佳准确率: XX%
    - 训练时间: XXXX
    """
    # ...
```

### 技巧3：版本管理

如果模型有多个版本：

```python
class MyCustomModelV1(nn.Module):
    # 版本1
    pass

class MyCustomModelV2(nn.Module):
    # 版本2
    pass

def create_model(model_name='resnet18', num_classes=10):
    # ...
    elif model_name == 'my_model_v1':
        return MyCustomModelV1(num_classes)
    elif model_name == 'my_model_v2':
        return MyCustomModelV2(num_classes)
```

---

**准备好了吗？** 找到你的模型定义，复制过来，配置一下，就能运行了！🎉

需要帮助？查看完整日志：`docker-compose logs -f model-service`
