#!/usr/bin/env python3
"""
简单的PyTorch模型训练示例
演示如何训练并保存一个可用于本服务的模型
"""
import torch
import torch.nn as nn
import torch.optim as optim
import torchvision
import torchvision.transforms as transforms
from torch.utils.data import DataLoader


class SimpleCNN(nn.Module):
    """简单的CNN分类器"""
    def __init__(self, num_classes=10):
        super(SimpleCNN, self).__init__()
        
        self.features = nn.Sequential(
            nn.Conv2d(3, 32, kernel_size=3, padding=1),
            nn.ReLU(inplace=True),
            nn.MaxPool2d(kernel_size=2, stride=2),
            
            nn.Conv2d(32, 64, kernel_size=3, padding=1),
            nn.ReLU(inplace=True),
            nn.MaxPool2d(kernel_size=2, stride=2),
            
            nn.Conv2d(64, 128, kernel_size=3, padding=1),
            nn.ReLU(inplace=True),
            nn.MaxPool2d(kernel_size=2, stride=2),
        )
        
        self.classifier = nn.Sequential(
            nn.Linear(128 * 28 * 28, 256),
            nn.ReLU(inplace=True),
            nn.Dropout(0.5),
            nn.Linear(256, num_classes)
        )
    
    def forward(self, x):
        x = self.features(x)
        x = x.view(x.size(0), -1)
        x = self.classifier(x)
        return x


def train_model():
    """训练模型示例"""
    print("="*60)
    print("PyTorch模型训练示例")
    print("="*60)
    
    # 设备
    device = torch.device('cuda' if torch.cuda.is_available() else 'cpu')
    print(f"\n使用设备: {device}")
    
    # 数据预处理
    transform = transforms.Compose([
        transforms.Resize((224, 224)),
        transforms.ToTensor(),
        transforms.Normalize((0.5, 0.5, 0.5), (0.5, 0.5, 0.5))
    ])
    
    # 加载CIFAR-10数据集作为示例
    print("\n加载数据集...")
    trainset = torchvision.datasets.CIFAR10(
        root='./data', 
        train=True,
        download=True, 
        transform=transform
    )
    
    trainloader = DataLoader(
        trainset, 
        batch_size=32,
        shuffle=True, 
        num_workers=2
    )
    
    # 类别名称
    classes = ('飞机', '汽车', '鸟', '猫', '鹿', '狗', '青蛙', '马', '船', '卡车')
    
    # 创建模型
    print("\n创建模型...")
    model = SimpleCNN(num_classes=10).to(device)
    
    # 损失函数和优化器
    criterion = nn.CrossEntropyLoss()
    optimizer = optim.Adam(model.parameters(), lr=0.001)
    
    # 训练几个epoch作为示例
    print("\n开始训练...")
    num_epochs = 2
    
    for epoch in range(num_epochs):
        running_loss = 0.0
        for i, data in enumerate(trainloader, 0):
            inputs, labels = data[0].to(device), data[1].to(device)
            
            # 前向传播
            optimizer.zero_grad()
            outputs = model(inputs)
            loss = criterion(outputs, labels)
            
            # 反向传播
            loss.backward()
            optimizer.step()
            
            # 统计
            running_loss += loss.item()
            if i % 100 == 99:
                print(f'  Epoch [{epoch + 1}/{num_epochs}], '
                      f'Batch [{i + 1}], '
                      f'Loss: {running_loss / 100:.4f}')
                running_loss = 0.0
    
    print('\n训练完成！')
    
    # 保存模型（方式1：完整模型 - 推荐）
    model_path = '../models/example_model.mph'
    print(f"\n保存完整模型到: {model_path}")
    model.eval()  # 设置为评估模式
    torch.save(model, model_path)
    print("✓ 完整模型已保存")
    
    # 保存模型（方式2：仅权重）
    weights_path = '../models/example_model_weights.pth'
    print(f"\n保存模型权重到: {weights_path}")
    torch.save(model.state_dict(), weights_path)
    print("✓ 模型权重已保存")
    
    # 保存标签文件
    labels_path = '../models/labels.txt'
    print(f"\n保存标签文件到: {labels_path}")
    with open(labels_path, 'w', encoding='utf-8') as f:
        for class_name in classes:
            f.write(f"{class_name}\n")
    print("✓ 标签文件已保存")
    
    # 测试加载
    print("\n测试模型加载...")
    loaded_model = torch.load(model_path, map_location=device)
    loaded_model.eval()
    
    # 测试前向传播
    dummy_input = torch.randn(1, 3, 224, 224).to(device)
    with torch.no_grad():
        output = loaded_model(dummy_input)
    print(f"✓ 模型加载成功，输出形状: {output.shape}")
    
    print("\n" + "="*60)
    print("模型已准备就绪！")
    print("="*60)
    print(f"\n下一步:")
    print(f"1. 使用完整模型: mv {model_path} ../models/model.mph")
    print(f"2. 启动服务: ./start_with_docker.sh")
    print(f"3. 测试: curl -X POST http://localhost:5000/classify -F 'image=@test.jpg'")


if __name__ == '__main__':
    train_model()
