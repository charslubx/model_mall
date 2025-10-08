#!/usr/bin/env python3
"""
PyTorch模型测试脚本
用于测试.mph/.pt/.pth模型文件是否可以正确加载
"""
import torch
import sys
import os


def test_load_model(model_path):
    """测试模型加载"""
    print("="*60)
    print("PyTorch模型测试工具")
    print("="*60)
    print(f"\n📁 模型路径: {model_path}")
    
    if not os.path.exists(model_path):
        print(f"❌ 错误: 文件不存在")
        return False
    
    file_size = os.path.getsize(model_path) / (1024 * 1024)
    print(f"📦 文件大小: {file_size:.2f} MB")
    
    # 检查CUDA
    print(f"\n🖥️  设备信息:")
    if torch.cuda.is_available():
        print(f"  ✓ CUDA可用")
        print(f"  ✓ GPU: {torch.cuda.get_device_name(0)}")
        print(f"  ✓ CUDA版本: {torch.version.cuda}")
        device = 'cuda'
    else:
        print(f"  ✗ CUDA不可用，使用CPU")
        device = 'cpu'
    
    print(f"  ✓ PyTorch版本: {torch.__version__}")
    
    try:
        # 尝试加载完整模型
        print(f"\n🔄 正在加载模型...")
        model = torch.load(model_path, map_location=device)
        print(f"✅ 成功加载模型")
        print(f"   类型: {type(model).__name__}")
        
        # 检查是否是nn.Module
        if isinstance(model, torch.nn.Module):
            print(f"\n📊 模型信息:")
            
            # 设置评估模式
            model.eval()
            print(f"  ✓ 已设置为评估模式")
            
            # 统计参数
            total_params = sum(p.numel() for p in model.parameters())
            trainable_params = sum(p.numel() for p in model.parameters() if p.requires_grad)
            print(f"  ✓ 总参数: {total_params:,}")
            print(f"  ✓ 可训练参数: {trainable_params:,}")
            
            # 打印模型结构（前几层）
            print(f"\n🏗️  模型结构（摘要）:")
            model_str = str(model)
            lines = model_str.split('\n')[:15]  # 只显示前15行
            for line in lines:
                print(f"  {line}")
            if len(model_str.split('\n')) > 15:
                print(f"  ... (省略 {len(model_str.split('\n')) - 15} 行)")
            
            # 测试前向传播
            print(f"\n🧪 测试前向传播:")
            try:
                dummy_input = torch.randn(1, 3, 224, 224).to(device)
                print(f"  ✓ 输入形状: {dummy_input.shape}")
                
                with torch.no_grad():
                    output = model(dummy_input)
                
                print(f"  ✓ 输出形状: {output.shape}")
                
                if len(output.shape) == 2:
                    num_classes = output.shape[1]
                    print(f"  ✓ 类别数量: {num_classes}")
                    
                    # 显示输出值范围
                    print(f"  ✓ 输出范围: [{output.min().item():.4f}, {output.max().item():.4f}]")
                    
                    # 应用softmax
                    probs = torch.softmax(output, dim=1)
                    top_prob, top_class = probs.max(dim=1)
                    print(f"  ✓ Softmax后最高概率: {top_prob.item():.4f}")
                
                print(f"\n✅ 模型测试通过！")
                print(f"\n💡 建议:")
                print(f"  1. 确保 labels.txt 文件包含 {num_classes if len(output.shape) == 2 else '正确数量的'} 个类别标签")
                print(f"  2. 将模型文件复制到: model_service/models/")
                print(f"  3. 启动服务: cd model_service && ./start_with_docker.sh")
                
                return True
                
            except Exception as e:
                print(f"  ❌ 前向传播失败: {str(e)}")
                print(f"\n💡 可能的问题:")
                print(f"  - 输入形状不匹配")
                print(f"  - 模型需要特定的输入格式")
                return False
        
        elif isinstance(model, dict):
            print(f"\n📦 这是一个字典对象")
            print(f"   键: {list(model.keys())}")
            
            # 检查常见的checkpoint格式
            if 'model_state_dict' in model:
                print(f"\n💡 检测到checkpoint格式")
                print(f"   包含: model_state_dict")
                if 'optimizer_state_dict' in model:
                    print(f"   包含: optimizer_state_dict")
                if 'epoch' in model:
                    print(f"   训练轮次: {model['epoch']}")
                
                state_dict = model['model_state_dict']
            else:
                print(f"\n💡 这可能是纯state_dict")
                state_dict = model
            
            # 分析state_dict
            print(f"\n📊 State Dict 信息:")
            print(f"  ✓ 参数层数: {len(state_dict)}")
            
            # 显示部分键
            keys = list(state_dict.keys())
            print(f"  ✓ 前几层:")
            for key in keys[:5]:
                if key in state_dict:
                    shape = state_dict[key].shape if hasattr(state_dict[key], 'shape') else 'N/A'
                    print(f"     - {key}: {shape}")
            
            if len(keys) > 5:
                print(f"     ... (省略 {len(keys) - 5} 层)")
            
            print(f"\n⚠️  这是权重文件，不是完整模型")
            print(f"\n💡 需要额外步骤:")
            print(f"  1. 创建 model_service/app/model_architecture.py")
            print(f"  2. 在其中定义与训练时相同的模型架构")
            print(f"  3. 使用 load_model_with_architecture() 方法加载")
            print(f"  4. 参考 PYTORCH_MODEL_GUIDE.md 中的示例")
            
            return False
        
        else:
            print(f"\n❌ 未知的模型格式: {type(model)}")
            return False
            
    except Exception as e:
        print(f"\n❌ 加载失败: {str(e)}")
        print(f"\n💡 可能的原因:")
        print(f"  - 模型文件损坏")
        print(f"  - PyTorch版本不兼容")
        print(f"  - 不是有效的PyTorch模型文件")
        return False


def main():
    """主函数"""
    if len(sys.argv) < 2:
        print("用法: python test_pytorch_model.py <model_path>")
        print("\n示例:")
        print("  python test_pytorch_model.py models/model.mph")
        print("  python test_pytorch_model.py models/model.pt")
        sys.exit(1)
    
    model_path = sys.argv[1]
    success = test_load_model(model_path)
    
    print("\n" + "="*60)
    if success:
        print("✅ 测试结果: 模型可以直接使用！")
    else:
        print("⚠️  测试结果: 需要额外配置")
    print("="*60)
    
    sys.exit(0 if success else 1)


if __name__ == '__main__':
    main()
