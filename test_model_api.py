#!/usr/bin/env python3
"""
模型服务API测试脚本
"""

import requests
import json
import time
import sys
from pathlib import Path

# 配置
BASE_URL = "http://localhost:8888"
TOKEN = "YOUR_ACCESS_TOKEN"  # 替换为实际的token
TEST_IMAGE = "test.jpg"  # 替换为实际的测试图片


def print_section(title):
    """打印分节标题"""
    print("\n" + "=" * 60)
    print(f"  {title}")
    print("=" * 60)


def print_response(response):
    """格式化打印响应"""
    try:
        data = response.json()
        print(json.dumps(data, indent=2, ensure_ascii=False))
    except:
        print(response.text)
    print(f"状态码: {response.status_code}")


def check_config():
    """检查配置"""
    print_section("配置检查")
    
    if TOKEN == "YOUR_ACCESS_TOKEN":
        print("❌ 错误: 请先设置有效的access token")
        print("   修改脚本中的 TOKEN 变量")
        return False
    
    if not Path(TEST_IMAGE).exists():
        print(f"❌ 错误: 测试图片 {TEST_IMAGE} 不存在")
        print("   请准备一张测试图片并更新脚本中的 TEST_IMAGE 变量")
        return False
    
    print("✅ 配置检查通过")
    return True


def upload_image():
    """上传图片"""
    print_section("步骤1: 上传图片")
    
    headers = {"Authorization": f"Bearer {TOKEN}"}
    files = {"image": open(TEST_IMAGE, "rb")}
    data = {"model_name": "resnet50"}
    
    response = requests.post(
        f"{BASE_URL}/api/images/upload",
        headers=headers,
        files=files,
        data=data
    )
    
    print_response(response)
    
    if response.status_code == 200:
        result = response.json()
        if result.get("code") == 0:
            data = result.get("data", {})
            return data.get("image_id"), data.get("task_id")
    
    return None, None


def get_task_status(task_id):
    """获取任务状态"""
    print_section("步骤2: 查询任务状态")
    
    headers = {"Authorization": f"Bearer {TOKEN}"}
    response = requests.get(
        f"{BASE_URL}/api/tasks/{task_id}/status",
        headers=headers
    )
    
    print_response(response)
    return response


def get_task_list():
    """获取任务列表"""
    print_section("步骤3: 获取任务列表")
    
    headers = {"Authorization": f"Bearer {TOKEN}"}
    params = {"page": 1, "page_size": 10}
    
    response = requests.get(
        f"{BASE_URL}/api/tasks",
        headers=headers,
        params=params
    )
    
    print_response(response)
    return response


def simulate_model_callback(task_id):
    """模拟模型服务回调"""
    print_section("步骤4: 模拟模型服务回调")
    
    payload = {
        "task_id": task_id,
        "status": "completed",
        "progress": 100,
        "results": [
            {
                "name": "猫",
                "code": "cat",
                "confidence": 0.9523
            },
            {
                "name": "动物",
                "code": "animal",
                "confidence": 0.8712
            },
            {
                "name": "宠物",
                "code": "pet",
                "confidence": 0.8234,
                "bbox": {
                    "x": 100,
                    "y": 150,
                    "width": 200,
                    "height": 180
                }
            }
        ]
    }
    
    response = requests.post(
        f"{BASE_URL}/api/model/callback",
        json=payload
    )
    
    print_response(response)
    return response


def get_image_labels(image_id):
    """获取图片标签"""
    print_section("步骤5: 获取图片标签")
    
    headers = {"Authorization": f"Bearer {TOKEN}"}
    response = requests.get(
        f"{BASE_URL}/api/images/{image_id}/labels",
        headers=headers
    )
    
    print_response(response)
    return response


def get_image_list():
    """获取图片列表"""
    print_section("步骤6: 获取图片列表")
    
    headers = {"Authorization": f"Bearer {TOKEN}"}
    params = {"page": 1, "page_size": 10}
    
    response = requests.get(
        f"{BASE_URL}/api/images",
        headers=headers,
        params=params
    )
    
    print_response(response)
    return response


def main():
    """主函数"""
    print("\n" + "=" * 60)
    print("  模型服务API测试脚本")
    print("=" * 60)
    
    # 检查配置
    if not check_config():
        return 1
    
    # 1. 上传图片
    image_id, task_id = upload_image()
    if not image_id or not task_id:
        print("\n❌ 上传失败，测试终止")
        return 1
    
    print(f"\n✅ 上传成功!")
    print(f"   Image ID: {image_id}")
    print(f"   Task ID: {task_id}")
    
    # 2. 查询任务状态
    time.sleep(1)
    get_task_status(task_id)
    
    # 3. 获取任务列表
    get_task_list()
    
    # 4. 模拟模型服务回调
    simulate_model_callback(task_id)
    
    # 等待一下让数据写入数据库
    time.sleep(1)
    
    # 5. 获取图片标签
    labels_response = get_image_labels(image_id)
    
    # 6. 获取图片列表
    get_image_list()
    
    # 最后再查询一次任务状态
    print_section("步骤7: 最终任务状态查询")
    get_task_status(task_id)
    
    # 总结
    print_section("测试完成")
    print(f"✅ 所有测试步骤已完成!")
    print(f"\n总结:")
    print(f"  - Image ID: {image_id}")
    print(f"  - Task ID: {task_id}")
    
    if labels_response.status_code == 200:
        result = labels_response.json()
        if result.get("code") == 0:
            labels = result.get("data", {}).get("labels", [])
            print(f"  - 识别标签数量: {len(labels)}")
            if labels:
                print(f"\n  识别标签:")
                for label in labels:
                    print(f"    - {label['label_name']}: {label['confidence']:.2%}")
    
    print(f"\n你可以在浏览器中访问上传的图片:")
    print(f"  {BASE_URL}/uploads/[文件路径]")
    print()
    
    return 0


if __name__ == "__main__":
    try:
        sys.exit(main())
    except KeyboardInterrupt:
        print("\n\n⚠️  测试被用户中断")
        sys.exit(1)
    except Exception as e:
        print(f"\n\n❌ 发生错误: {e}")
        import traceback
        traceback.print_exc()
        sys.exit(1)
