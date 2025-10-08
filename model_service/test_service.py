#!/usr/bin/env python3
"""
模型服务测试脚本
"""
import requests
import sys
import os


def test_health(base_url):
    """测试健康检查接口"""
    print("测试健康检查接口...")
    try:
        response = requests.get(f"{base_url}/health")
        print(f"状态码: {response.status_code}")
        print(f"响应: {response.json()}")
        return response.status_code == 200
    except Exception as e:
        print(f"错误: {str(e)}")
        return False


def test_info(base_url):
    """测试模型信息接口"""
    print("\n测试模型信息接口...")
    try:
        response = requests.get(f"{base_url}/info")
        print(f"状态码: {response.status_code}")
        print(f"响应: {response.json()}")
        return response.status_code == 200
    except Exception as e:
        print(f"错误: {str(e)}")
        return False


def test_classify(base_url, image_path):
    """测试分类接口"""
    print(f"\n测试分类接口...")
    print(f"图片路径: {image_path}")
    
    if not os.path.exists(image_path):
        print(f"错误: 图片文件不存在: {image_path}")
        return False
    
    try:
        with open(image_path, 'rb') as f:
            files = {'image': f}
            response = requests.post(f"{base_url}/classify", files=files)
        
        print(f"状态码: {response.status_code}")
        result = response.json()
        print(f"响应: {result}")
        
        if result.get('success'):
            print("\n分类结果:")
            for item in result.get('results', []):
                print(f"  - {item['label']}: {item['confidence']:.4f}")
        
        return response.status_code == 200
    except Exception as e:
        print(f"错误: {str(e)}")
        return False


def main():
    """主函数"""
    # 服务地址
    base_url = os.environ.get('SERVICE_URL', 'http://localhost:5000')
    
    print("========================================")
    print(f"模型服务测试")
    print(f"服务地址: {base_url}")
    print("========================================\n")
    
    # 测试健康检查
    if not test_health(base_url):
        print("\n健康检查失败！")
        sys.exit(1)
    
    # 测试模型信息
    if not test_info(base_url):
        print("\n获取模型信息失败！")
        sys.exit(1)
    
    # 测试分类接口（如果提供了测试图片）
    if len(sys.argv) > 1:
        image_path = sys.argv[1]
        if not test_classify(base_url, image_path):
            print("\n分类测试失败！")
            sys.exit(1)
    else:
        print("\n提示: 可以通过命令行参数提供测试图片路径")
        print(f"用法: python test_service.py <image_path>")
    
    print("\n========================================")
    print("所有测试通过！")
    print("========================================")


if __name__ == '__main__':
    main()
