"""
模型架构定义文件
请在此文件中定义你的模型架构，需要与训练时完全一致
"""
import torch
import torch.nn as nn
import torch.nn.functional as F
import torchvision.models as models


# ============================================
# 方式1: 使用预训练模型（如ResNet, EfficientNet等）
# ============================================

def create_resnet_model(num_classes=10):
    """
    创建ResNet模型
    
    示例用法:
        model = create_resnet_model(num_classes=10)
    """
    model = models.resnet18(pretrained=False)
    num_features = model.fc.in_features
    model.fc = nn.Linear(num_features, num_classes)
    return model


def create_efficientnet_model(num_classes=10):
    """创建EfficientNet模型"""
    model = models.efficientnet_b0(pretrained=False)
    num_features = model.classifier[1].in_features
    model.classifier[1] = nn.Linear(num_features, num_classes)
    return model


def create_mobilenet_model(num_classes=10):
    """创建MobileNet模型"""
    model = models.mobilenet_v2(pretrained=False)
    num_features = model.classifier[1].in_features
    model.classifier[1] = nn.Linear(num_features, num_classes)
    return model


# ============================================
# 方式2: 自定义CNN模型
# ============================================

class SimpleCNN(nn.Module):
    """
    简单的CNN分类器
    请根据你的训练代码修改这个类
    """
    def __init__(self, num_classes=10):
        super(SimpleCNN, self).__init__()
        
        # 卷积层
        self.conv1 = nn.Conv2d(3, 32, kernel_size=3, padding=1)
        self.bn1 = nn.BatchNorm2d(32)
        
        self.conv2 = nn.Conv2d(32, 64, kernel_size=3, padding=1)
        self.bn2 = nn.BatchNorm2d(64)
        
        self.conv3 = nn.Conv2d(64, 128, kernel_size=3, padding=1)
        self.bn3 = nn.BatchNorm2d(128)
        
        # 全连接层
        self.fc1 = nn.Linear(128 * 28 * 28, 512)
        self.dropout = nn.Dropout(0.5)
        self.fc2 = nn.Linear(512, num_classes)
    
    def forward(self, x):
        # 第一层
        x = self.conv1(x)
        x = self.bn1(x)
        x = F.relu(x)
        x = F.max_pool2d(x, 2)
        
        # 第二层
        x = self.conv2(x)
        x = self.bn2(x)
        x = F.relu(x)
        x = F.max_pool2d(x, 2)
        
        # 第三层
        x = self.conv3(x)
        x = self.bn3(x)
        x = F.relu(x)
        x = F.max_pool2d(x, 2)
        
        # 展平
        x = x.view(x.size(0), -1)
        
        # 全连接层
        x = self.fc1(x)
        x = F.relu(x)
        x = self.dropout(x)
        x = self.fc2(x)
        
        return x


class CustomCNN(nn.Module):
    """
    自定义CNN模型模板
    
    ⚠️ 重要：将下面的结构替换为你训练时使用的实际模型结构！
    """
    def __init__(self, num_classes=10):
        super(CustomCNN, self).__init__()
        
        # TODO: 在这里定义你的模型层
        # 示例:
        self.features = nn.Sequential(
            nn.Conv2d(3, 64, kernel_size=3, padding=1),
            nn.ReLU(inplace=True),
            nn.MaxPool2d(kernel_size=2, stride=2),
            
            nn.Conv2d(64, 128, kernel_size=3, padding=1),
            nn.ReLU(inplace=True),
            nn.MaxPool2d(kernel_size=2, stride=2),
        )
        
        self.classifier = nn.Sequential(
            nn.Linear(128 * 56 * 56, 256),
            nn.ReLU(inplace=True),
            nn.Dropout(0.5),
            nn.Linear(256, num_classes)
        )
    
    def forward(self, x):
        # TODO: 定义前向传播逻辑
        x = self.features(x)
        x = x.view(x.size(0), -1)
        x = self.classifier(x)
        return x


# ============================================
# 工厂函数：根据模型名称创建模型
# ============================================

def create_model(model_name='resnet18', num_classes=10):
    """
    根据模型名称创建模型实例
    
    Args:
        model_name: 模型名称 ('resnet18', 'efficientnet', 'mobilenet', 'simple_cnn', 'custom_cnn')
        num_classes: 分类类别数
    
    Returns:
        模型实例
    """
    model_name = model_name.lower()
    
    if model_name == 'resnet18':
        return create_resnet_model(num_classes)
    elif model_name == 'efficientnet':
        return create_efficientnet_model(num_classes)
    elif model_name == 'mobilenet':
        return create_mobilenet_model(num_classes)
    elif model_name == 'simple_cnn':
        return SimpleCNN(num_classes)
    elif model_name == 'custom_cnn':
        return CustomCNN(num_classes)
    else:
        raise ValueError(f"未知的模型名称: {model_name}")


# ============================================
# 使用说明
# ============================================

"""
📝 如何使用这个文件：

1. 找到你训练时的模型定义代码

2. 将模型类复制到这个文件（或修改上面的CustomCNN）

3. 确保类名和结构与训练时完全一致

4. 在 api.py 中使用：

   from app.model_architecture import create_model
   
   # 创建模型
   model = create_model('resnet18', num_classes=10)
   
   # 加载checkpoint
   model_loader.load_checkpoint_with_architecture(model)

示例：如果你的训练代码是这样的：

    class MyModel(nn.Module):
        def __init__(self):
            super().__init__()
            self.conv1 = nn.Conv2d(3, 64, 3)
            self.fc1 = nn.Linear(64*30*30, 10)
        
        def forward(self, x):
            x = F.relu(self.conv1(x))
            x = x.view(x.size(0), -1)
            x = self.fc1(x)
            return x
    
    model = MyModel()

那么你需要：
1. 将 MyModel 类复制到这个文件
2. 创建一个函数来实例化它：

    def create_my_model(num_classes=10):
        return MyModel()

3. 在 create_model() 中添加选项，或直接在 api.py 中使用：

    from app.model_architecture import MyModel
    model = MyModel()
"""
