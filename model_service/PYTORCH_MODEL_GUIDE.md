# PyTorch模型使用指南

本文档专门介绍如何在模型服务中使用PyTorch训练的模型（.pt, .pth, .mph格式）。

## 📋 PyTorch模型保存格式

PyTorch模型通常有两种保存方式：

### 1. 完整模型保存（推荐）

```python
# 训练代码中
import torch

# 保存完整模型（包含架构和权重）
torch.save(model, 'model.pt')
# 或
torch.save(model, 'model.mph')
```

**优点：** 加载简单，不需要额外定义模型架构

**缺点：** 文件较大，可能有兼容性问题

### 2. 仅保存权重（state_dict）

```python
# 训练代码中
import torch

# 仅保存模型权重
torch.save(model.state_dict(), 'model_weights.pt')
```

**优点：** 文件较小，兼容性好

**缺点：** 加载时需要先定义模型架构

## 🚀 快速开始

### 方式1: 使用完整模型（推荐）

如果你的`.mph`文件是通过`torch.save(model, ...)`保存的完整模型：

```bash
# 1. 复制模型文件
cp /path/to/your/model.mph model_service/models/

# 2. 准备标签文件
cat > model_service/models/labels.txt << EOF
类别1
类别2
类别3
EOF

# 3. 启动服务
cd model_service
./start_with_docker.sh

# 4. 测试
curl -X POST http://localhost:5000/classify \
  -F "image=@test.jpg"
```

模型会自动被识别为PyTorch格式并加载！

### 方式2: 使用state_dict

如果你的模型是权重文件，需要先定义模型架构。

#### 步骤1: 创建模型定义文件

创建 `model_service/app/model_architecture.py`：

```python
import torch
import torch.nn as nn
import torchvision.models as models

class MyImageClassifier(nn.Module):
    """
    你的模型架构定义
    这里需要与训练时使用的架构完全一致
    """
    def __init__(self, num_classes=10):
        super(MyImageClassifier, self).__init__()
        
        # 示例：使用ResNet18作为backbone
        self.backbone = models.resnet18(pretrained=False)
        
        # 替换最后的全连接层
        num_features = self.backbone.fc.in_features
        self.backbone.fc = nn.Linear(num_features, num_classes)
    
    def forward(self, x):
        return self.backbone(x)


# 或者如果你的模型是自定义卷积网络
class CustomCNN(nn.Module):
    def __init__(self, num_classes=10):
        super(CustomCNN, self).__init__()
        
        self.features = nn.Sequential(
            nn.Conv2d(3, 64, kernel_size=3, padding=1),
            nn.ReLU(inplace=True),
            nn.MaxPool2d(kernel_size=2, stride=2),
            
            nn.Conv2d(64, 128, kernel_size=3, padding=1),
            nn.ReLU(inplace=True),
            nn.MaxPool2d(kernel_size=2, stride=2),
            
            nn.Conv2d(128, 256, kernel_size=3, padding=1),
            nn.ReLU(inplace=True),
            nn.MaxPool2d(kernel_size=2, stride=2),
        )
        
        self.classifier = nn.Sequential(
            nn.Linear(256 * 28 * 28, 512),
            nn.ReLU(inplace=True),
            nn.Dropout(0.5),
            nn.Linear(512, num_classes)
        )
    
    def forward(self, x):
        x = self.features(x)
        x = x.view(x.size(0), -1)
        x = self.classifier(x)
        return x
```

#### 步骤2: 修改API加载逻辑

修改 `model_service/app/api.py` 中的 `init_model()` 函数：

```python
def init_model():
    """初始化模型"""
    global model_loader
    
    model_path = os.environ.get('MODEL_PATH', '/app/models/model.mph')
    model_name = os.environ.get('MODEL_NAME', 'image-classifier')
    
    logger.info(f"初始化模型服务...")
    logger.info(f"模型路径: {model_path}")
    
    try:
        from app.model_loader import PyTorchModelLoader
        from app.model_architecture import MyImageClassifier  # 导入你的模型定义
        
        # 创建PyTorch加载器
        model_loader = PyTorchModelLoader(model_path, model_name)
        
        # 使用自定义架构加载state_dict
        model_class = MyImageClassifier(num_classes=10)  # 根据实际类别数修改
        model_loader.load_model_with_architecture(model_class)
        
        logger.info("模型初始化成功")
    except Exception as e:
        logger.error(f"模型初始化失败: {str(e)}")
        raise
```

## 🔧 常见模型架构

### ResNet系列

```python
import torchvision.models as models

# ResNet18
model = models.resnet18(pretrained=False)
model.fc = nn.Linear(512, num_classes)

# ResNet50
model = models.resnet50(pretrained=False)
model.fc = nn.Linear(2048, num_classes)
```

### EfficientNet

```python
import torchvision.models as models

model = models.efficientnet_b0(pretrained=False)
model.classifier[1] = nn.Linear(1280, num_classes)
```

### VGG系列

```python
import torchvision.models as models

model = models.vgg16(pretrained=False)
model.classifier[6] = nn.Linear(4096, num_classes)
```

### MobileNet

```python
import torchvision.models as models

model = models.mobilenet_v2(pretrained=False)
model.classifier[1] = nn.Linear(1280, num_classes)
```

## 🖼️ 图片预处理

PyTorch模型通常需要特定的图片预处理。如果需要自定义预处理，修改 `PyTorchModelLoader.preprocess_image()` 方法：

```python
def preprocess_image(self, image_data: bytes, target_size=(224, 224)) -> np.ndarray:
    """自定义预处理"""
    from PIL import Image
    import torchvision.transforms as transforms
    
    # 加载图片
    image = Image.open(io.BytesIO(image_data)).convert('RGB')
    
    # 定义预处理流程（与训练时保持一致）
    transform = transforms.Compose([
        transforms.Resize(256),
        transforms.CenterCrop(224),
        transforms.ToTensor(),
        transforms.Normalize(
            mean=[0.485, 0.456, 0.406],  # ImageNet标准化参数
            std=[0.229, 0.224, 0.225]
        )
    ])
    
    # 应用转换
    img_tensor = transform(image)
    img_tensor = img_tensor.unsqueeze(0)  # 添加batch维度
    
    return img_tensor
```

## 🎮 GPU加速

### 自动检测GPU

模型服务会自动检测并使用可用的GPU：

```python
# PyTorchModelLoader会自动检测
device = torch.device('cuda' if torch.cuda.is_available() else 'cpu')
```

### 强制使用CPU或GPU

通过环境变量设置：

```bash
# docker-compose.yml
environment:
  - DEVICE=cuda  # 或 cpu
```

### Docker GPU支持

修改 `Dockerfile`：

```dockerfile
# 使用CUDA基础镜像
FROM pytorch/pytorch:2.1.0-cuda11.8-cudnn8-runtime

# ... 其他配置
```

修改 `docker-compose.yml`：

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
```

需要先安装 NVIDIA Docker Runtime。

## 🧪 测试模型加载

创建测试脚本 `test_model.py`：

```python
#!/usr/bin/env python3
import torch
import sys

def test_load_model(model_path):
    """测试模型加载"""
    print(f"测试加载模型: {model_path}")
    
    # 检查CUDA
    if torch.cuda.is_available():
        print(f"✓ CUDA可用: {torch.cuda.get_device_name(0)}")
        device = 'cuda'
    else:
        print("✗ CUDA不可用，使用CPU")
        device = 'cpu'
    
    try:
        # 尝试加载完整模型
        print("\n尝试加载完整模型...")
        model = torch.load(model_path, map_location=device)
        print("✓ 成功加载完整模型")
        print(f"  模型类型: {type(model)}")
        
        # 检查模型结构
        if hasattr(model, 'eval'):
            model.eval()
            print("✓ 模型已设置为评估模式")
        
        # 测试前向传播
        print("\n测试前向传播...")
        dummy_input = torch.randn(1, 3, 224, 224).to(device)
        with torch.no_grad():
            output = model(dummy_input)
        print(f"✓ 前向传播成功")
        print(f"  输出形状: {output.shape}")
        print(f"  类别数量: {output.shape[1]}")
        
        return True
        
    except Exception as e:
        print(f"✗ 加载完整模型失败: {str(e)}")
        
        # 尝试加载state_dict
        try:
            print("\n尝试加载state_dict...")
            state_dict = torch.load(model_path, map_location=device)
            
            if isinstance(state_dict, dict):
                print("✓ 成功加载state_dict")
                print(f"  参数数量: {len(state_dict)}")
                print(f"  部分键: {list(state_dict.keys())[:5]}")
                print("\n⚠ 这是权重文件，需要先定义模型架构")
                return False
            else:
                print(f"✗ 未知格式: {type(state_dict)}")
                return False
                
        except Exception as e2:
            print(f"✗ 加载state_dict也失败: {str(e2)}")
            return False

if __name__ == '__main__':
    if len(sys.argv) < 2:
        print("用法: python test_model.py <model_path>")
        sys.exit(1)
    
    model_path = sys.argv[1]
    success = test_load_model(model_path)
    
    if success:
        print("\n✅ 模型可以直接使用！")
    else:
        print("\n❌ 需要先定义模型架构")
```

运行测试：

```bash
python test_model.py models/model.mph
```

## 🔍 调试技巧

### 1. 检查模型结构

```python
import torch

model = torch.load('model.mph')
print(model)  # 打印模型结构
```

### 2. 查看state_dict

```python
import torch

checkpoint = torch.load('model.mph')
if isinstance(checkpoint, dict):
    print("Keys:", checkpoint.keys())
    if 'model_state_dict' in checkpoint:
        print("这是checkpoint格式")
    else:
        print("这是纯state_dict")
```

### 3. 检查输入输出形状

```python
import torch

model = torch.load('model.mph')
model.eval()

# 测试输入
dummy_input = torch.randn(1, 3, 224, 224)
output = model(dummy_input)
print(f"输入形状: {dummy_input.shape}")
print(f"输出形状: {output.shape}")
print(f"类别数: {output.shape[1]}")
```

## 📝 完整示例

### 示例1: ResNet18图片分类

```python
# model_architecture.py
import torch
import torch.nn as nn
import torchvision.models as models

def create_resnet18_classifier(num_classes):
    model = models.resnet18(pretrained=False)
    model.fc = nn.Linear(512, num_classes)
    return model
```

### 示例2: 自定义CNN

```python
# model_architecture.py
import torch
import torch.nn as nn

class SimpleCNN(nn.Module):
    def __init__(self, num_classes=10):
        super(SimpleCNN, self).__init__()
        
        self.conv_layers = nn.Sequential(
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
        
        self.fc_layers = nn.Sequential(
            nn.Linear(128 * 28 * 28, 256),
            nn.ReLU(),
            nn.Dropout(0.5),
            nn.Linear(256, num_classes)
        )
    
    def forward(self, x):
        x = self.conv_layers(x)
        x = x.view(x.size(0), -1)
        x = self.fc_layers(x)
        return x
```

## 🚨 常见问题

### 问题1: RuntimeError: Expected all tensors to be on the same device

**原因：** 模型和输入数据不在同一设备上

**解决：**
```python
# 确保模型和数据都在同一设备
device = torch.device('cuda' if torch.cuda.is_available() else 'cpu')
model = model.to(device)
input_tensor = input_tensor.to(device)
```

### 问题2: 模型输出维度错误

**原因：** 类别数量不匹配

**解决：**
```python
# 检查并修改输出层
if hasattr(model, 'fc'):
    num_features = model.fc.in_features
    model.fc = nn.Linear(num_features, correct_num_classes)
```

### 问题3: 预处理不一致导致准确率低

**原因：** 推理时的预处理与训练时不一致

**解决：** 确保使用与训练时完全相同的预处理流程

## 📚 参考资源

- [PyTorch官方文档](https://pytorch.org/docs/stable/index.html)
- [torchvision模型](https://pytorch.org/vision/stable/models.html)
- [PyTorch模型保存和加载](https://pytorch.org/tutorials/beginner/saving_loading_models.html)
