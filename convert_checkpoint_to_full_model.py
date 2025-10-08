#!/usr/bin/env python3
"""
将checkpoint格式转换为完整模型

使用方法:
1. 修改下面的模型定义（从你的训练代码复制）
2. 运行: python convert_checkpoint_to_full_model.py
"""

import torch
import torch.nn as nn
import torch.nn.functional as F


# ============================================
# 第1步：复制你的模型定义到这里
# ============================================

class YourModel(nn.Module):
    """
    将你训练代码中的模型类完整复制到这里
    """
    def __init__(self, num_classes=10):
        super(YourModel, self).__init__()
        
        # TODO: 复制你的模型层定义
        # 示例:
        # self.conv1 = nn.Conv2d(3, 64, 3, padding=1)
        # self.fc1 = nn.Linear(512, num_classes)
        pass
    
    def forward(self, x):
        # TODO: 复制你的前向传播逻辑
        # 示例:
        # x = F.relu(self.conv1(x))
        # x = self.fc1(x)
        pass


# ============================================
# 第2步：运行转换
# ============================================

def convert_checkpoint_to_full_model(
    checkpoint_path='checkpoint_best.pth',
    output_path='model_full.pth',
    num_classes=10
):
    """
    转换checkpoint为完整模型
    
    Args:
        checkpoint_path: checkpoint文件路径
        output_path: 输出的完整模型路径
        num_classes: 分类类别数量
    """
    print("="*60)
    print("Checkpoint转完整模型工具")
    print("="*60)
    
    # 1. 加载checkpoint
    print(f"\n[1/4] 加载checkpoint: {checkpoint_path}")
    checkpoint = torch.load(checkpoint_path, map_location='cpu')
    
    # 显示checkpoint信息
    print(f"  ✓ Checkpoint键: {list(checkpoint.keys())}")
    if 'epoch' in checkpoint:
        print(f"  ✓ 训练轮次: {checkpoint['epoch']}")
    if 'best_val_loss' in checkpoint:
        print(f"  ✓ 最佳验证损失: {checkpoint['best_val_loss']:.4f}")
    
    # 2. 创建模型
    print(f"\n[2/4] 创建模型实例 (num_classes={num_classes})")
    model = YourModel(num_classes=num_classes)
    print(f"  ✓ 模型类型: {type(model).__name__}")
    
    # 3. 加载权重
    print(f"\n[3/4] 加载权重")
    if 'model_state' in checkpoint:
        model_state = checkpoint['model_state']
    elif 'model_state_dict' in checkpoint:
        model_state = checkpoint['model_state_dict']
    else:
        model_state = checkpoint
    
    model.load_state_dict(model_state)
    print(f"  ✓ 权重加载成功")
    
    # 设置为评估模式
    model.eval()
    print(f"  ✓ 设置为评估模式")
    
    # 4. 保存完整模型
    print(f"\n[4/4] 保存完整模型: {output_path}")
    torch.save(model, output_path)
    
    # 验证
    print(f"\n验证转换结果:")
    loaded_model = torch.load(output_path, map_location='cpu')
    print(f"  ✓ 模型可以正常加载")
    print(f"  ✓ 模型类型: {type(loaded_model).__name__}")
    
    # 测试前向传播
    dummy_input = torch.randn(1, 3, 224, 224)
    with torch.no_grad():
        output = loaded_model(dummy_input)
    print(f"  ✓ 前向传播正常")
    print(f"  ✓ 输出形状: {output.shape}")
    
    print(f"\n{'='*60}")
    print(f"✅ 转换成功！")
    print(f"{'='*60}")
    print(f"\n完整模型已保存到: {output_path}")
    print(f"\n🎉 现在你可以直接使用这个模型，无需提供模型架构！")
    print(f"\n部署步骤:")
    print(f"  1. cp {output_path} /workspace/model_service/models/")
    print(f"  2. 编辑 docker-compose.yml，移除 MODEL_ARCH 配置")
    print(f"  3. ./start_all_services.sh")


if __name__ == '__main__':
    import sys
    
    # 配置参数
    checkpoint_path = 'checkpoint_best.pth'  # 你的checkpoint文件
    output_path = 'model_full.pth'           # 输出的完整模型
    num_classes = 10                         # 你的类别数量
    
    # 从命令行参数读取（可选）
    if len(sys.argv) > 1:
        checkpoint_path = sys.argv[1]
    if len(sys.argv) > 2:
        output_path = sys.argv[2]
    if len(sys.argv) > 3:
        num_classes = int(sys.argv[3])
    
    try:
        convert_checkpoint_to_full_model(checkpoint_path, output_path, num_classes)
    except Exception as e:
        print(f"\n❌ 转换失败: {str(e)}")
        print(f"\n请确保:")
        print(f"  1. 已修改脚本中的 YourModel 类定义")
        print(f"  2. checkpoint文件存在")
        print(f"  3. num_classes 参数正确")
        sys.exit(1)
