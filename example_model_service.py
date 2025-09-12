#!/usr/bin/env python3
"""
示例模型服务
这是一个简单的Flask应用，用于演示如何与Go后端集成
实际使用时，您应该替换为真实的机器学习模型
"""

from flask import Flask, request, jsonify
import time
import random
from PIL import Image
import io
import base64

app = Flask(__name__)

# 模拟的标签列表
SAMPLE_LABELS = [
    {"label": "猫", "code": "cat", "confidence": 0.95},
    {"label": "狗", "code": "dog", "confidence": 0.88},
    {"label": "汽车", "code": "car", "confidence": 0.76},
    {"label": "建筑物", "code": "building", "confidence": 0.82},
    {"label": "树", "code": "tree", "confidence": 0.91},
    {"label": "人", "code": "person", "confidence": 0.89},
    {"label": "天空", "code": "sky", "confidence": 0.73},
    {"label": "水", "code": "water", "confidence": 0.67},
]

@app.route('/health', methods=['GET'])
def health_check():
    """健康检查接口"""
    return jsonify({
        "status": "healthy",
        "timestamp": time.time()
    })

@app.route('/predict', methods=['POST'])
def predict():
    """图片分类预测接口"""
    start_time = time.time()
    
    try:
        # 检查是否有文件上传
        if 'image' not in request.files:
            return jsonify({
                "success": False,
                "message": "没有上传图片文件",
                "predictions": []
            }), 400
        
        file = request.files['image']
        if file.filename == '':
            return jsonify({
                "success": False,
                "message": "文件名为空",
                "predictions": []
            }), 400
        
        # 获取参数
        model_name = request.form.get('model_name', 'default')
        min_confidence = float(request.form.get('min_confidence', 0.5))
        
        # 验证图片
        try:
            image = Image.open(io.BytesIO(file.read()))
            # 这里可以添加图片预处理逻辑
            print(f"处理图片: {file.filename}, 尺寸: {image.size}, 模式: {image.mode}")
        except Exception as e:
            return jsonify({
                "success": False,
                "message": f"图片格式错误: {str(e)}",
                "predictions": []
            }), 400
        
        # 模拟模型推理过程
        time.sleep(random.uniform(0.1, 0.5))  # 模拟处理时间
        
        # 生成随机预测结果
        num_predictions = random.randint(1, 4)
        predictions = []
        
        selected_labels = random.sample(SAMPLE_LABELS, num_predictions)
        for label_info in selected_labels:
            # 添加一些随机性
            confidence = label_info["confidence"] + random.uniform(-0.1, 0.1)
            confidence = max(0.0, min(1.0, confidence))  # 确保在0-1范围内
            
            if confidence >= min_confidence:
                prediction = {
                    "label": label_info["label"],
                    "code": label_info["code"],
                    "confidence": round(confidence, 4)
                }
                
                # 随机添加边界框信息（可选）
                if random.random() > 0.5:
                    prediction["bounding_box"] = {
                        "x": random.uniform(0, 0.5),
                        "y": random.uniform(0, 0.5),
                        "width": random.uniform(0.2, 0.5),
                        "height": random.uniform(0.2, 0.5)
                    }
                
                predictions.append(prediction)
        
        # 按置信度排序
        predictions.sort(key=lambda x: x["confidence"], reverse=True)
        
        process_time = int((time.time() - start_time) * 1000)  # 转换为毫秒
        
        return jsonify({
            "success": True,
            "message": "分类成功",
            "process_time": process_time,
            "predictions": predictions
        })
        
    except Exception as e:
        return jsonify({
            "success": False,
            "message": f"服务器内部错误: {str(e)}",
            "predictions": []
        }), 500

@app.route('/models', methods=['GET'])
def get_models():
    """获取可用模型列表"""
    return jsonify({
        "models": [
            {
                "name": "default",
                "version": "1.0.0",
                "description": "默认图片分类模型",
                "supported_formats": ["jpeg", "jpg", "png", "gif", "webp", "bmp"]
            },
            {
                "name": "resnet50",
                "version": "1.0.0",
                "description": "ResNet-50图片分类模型",
                "supported_formats": ["jpeg", "jpg", "png"]
            },
            {
                "name": "efficientnet",
                "version": "1.0.0",
                "description": "EfficientNet图片分类模型",
                "supported_formats": ["jpeg", "jpg", "png"]
            }
        ]
    })

@app.route('/info', methods=['GET'])
def get_info():
    """获取服务信息"""
    return jsonify({
        "service": "Image Classification Model Service",
        "version": "1.0.0",
        "status": "running",
        "endpoints": {
            "health": "/health",
            "predict": "/predict",
            "models": "/models",
            "info": "/info"
        },
        "supported_formats": ["jpeg", "jpg", "png", "gif", "webp", "bmp"],
        "max_file_size": "10MB"
    })

if __name__ == '__main__':
    print("启动示例模型服务...")
    print("服务地址: http://localhost:8080")
    print("健康检查: http://localhost:8080/health")
    print("模型信息: http://localhost:8080/info")
    
    app.run(
        host='0.0.0.0',
        port=8080,
        debug=True
    )