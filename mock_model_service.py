#!/usr/bin/env python3
"""
Mock模型服务 - 用于测试模型服务集成

这是一个简单的模拟模型服务，用于测试后端与模型服务的集成。
它会接收图片，模拟识别过程，然后回调后端服务。

运行方式:
    python3 mock_model_service.py

然后启动后端服务，上传图片即可看到完整的流程。
"""

from flask import Flask, request, jsonify
import requests
import time
import threading
import random
import os
import statistics

app = Flask(__name__)

# 配置
CALLBACK_URL = "http://localhost:8888/api/model/callback"
PROCESSING_TIME = 3  # 模拟处理时间（秒）

# 模拟的识别结果
SAMPLE_LABELS = [
    {"name": "猫", "code": "cat", "confidence_range": (0.85, 0.98)},
    {"name": "狗", "code": "dog", "confidence_range": (0.80, 0.95)},
    {"name": "鸟", "code": "bird", "confidence_range": (0.75, 0.92)},
    {"name": "汽车", "code": "car", "confidence_range": (0.82, 0.96)},
    {"name": "建筑", "code": "building", "confidence_range": (0.78, 0.94)},
    {"name": "人", "code": "person", "confidence_range": (0.88, 0.99)},
    {"name": "风景", "code": "landscape", "confidence_range": (0.70, 0.90)},
    {"name": "食物", "code": "food", "confidence_range": (0.83, 0.97)},
]


def process_image(task_id, callback_url, image_filename):
    """
    模拟图片处理过程
    """
    try:
        print(f"\n[任务 {task_id}] 开始处理...")
        
        # 1. 更新状态为处理中
        callback_data = {
            "task_id": task_id,
            "status": "processing",
            "progress": 0
        }
        requests.post(callback_url, json=callback_data)
        print(f"[任务 {task_id}] 状态: 处理中 (0%)")
        
        # 2. 模拟处理过程，逐步更新进度
        for progress in [25, 50, 75]:
            time.sleep(PROCESSING_TIME / 4)
            callback_data = {
                "task_id": task_id,
                "status": "processing",
                "progress": progress
            }
            requests.post(callback_url, json=callback_data)
            print(f"[任务 {task_id}] 进度: {progress}%")
        
        # 3. 生成随机识别结果
        num_labels = random.randint(2, 5)
        selected_labels = random.sample(SAMPLE_LABELS, num_labels)
        
        results = []
        for label in selected_labels:
            confidence = random.uniform(*label["confidence_range"])
            result = {
                "name": label["name"],
                "code": label["code"],
                "confidence": round(confidence, 4)
            }
            
            # 50% 的概率添加边界框
            if random.random() > 0.5:
                result["bbox"] = {
                    "x": random.randint(10, 200),
                    "y": random.randint(10, 200),
                    "width": random.randint(100, 400),
                    "height": random.randint(100, 400)
                }
            
            results.append(result)
        
        # 按置信度排序
        results.sort(key=lambda x: x["confidence"], reverse=True)
        
        # 4. 完成处理，返回结果
        time.sleep(PROCESSING_TIME / 4)
        callback_data = {
            "task_id": task_id,
            "status": "completed",
            "progress": 100,
            "results": results
        }
        
        response = requests.post(callback_url, json=callback_data)
        print(f"[任务 {task_id}] 完成! 识别到 {len(results)} 个标签")
        print(f"[任务 {task_id}] 回调响应: {response.status_code}")
        
        for result in results:
            print(f"  - {result['name']}: {result['confidence']:.2%}")
        
    except Exception as e:
        print(f"[任务 {task_id}] 错误: {e}")
        # 发送失败回调
        callback_data = {
            "task_id": task_id,
            "status": "failed",
            "error": str(e)
        }
        try:
            requests.post(callback_url, json=callback_data)
        except:
            pass


@app.route('/api/v1/recognize/upload', methods=['POST'])
def recognize_upload():
    """
    接收上传的图片进行识别
    """
    try:
        # 获取参数
        task_id = request.form.get('task_id')
        model_name = request.form.get('model_name', 'default')
        callback = request.form.get('callback', CALLBACK_URL)
        
        if 'image' not in request.files:
            return jsonify({
                "code": 400,
                "message": "没有上传图片"
            }), 400
        
        image = request.files['image']
        if not task_id:
            return jsonify({
                "code": 400,
                "message": "缺少task_id参数"
            }), 400
        
        print(f"\n{'='*60}")
        print(f"收到识别请求:")
        print(f"  Task ID: {task_id}")
        print(f"  Model: {model_name}")
        print(f"  Callback: {callback}")
        print(f"  Image: {image.filename}")
        print(f"{'='*60}")
        
        # 异步处理图片
        thread = threading.Thread(
            target=process_image,
            args=(task_id, callback, image.filename)
        )
        thread.start()
        
        return jsonify({
            "code": 0,
            "message": "success",
            "data": {
                "task_id": task_id,
                "status": "pending"
            }
        })
        
    except Exception as e:
        print(f"错误: {e}")
        return jsonify({
            "code": 500,
            "message": str(e)
        }), 500


@app.route('/api/v1/recognize/url', methods=['POST'])
def recognize_url():
    """
    通过URL识别图片
    """
    try:
        data = request.json
        task_id = data.get('task_id')
        image_url = data.get('image_url')
        model_name = data.get('model_name', 'default')
        callback = data.get('callback', CALLBACK_URL)
        
        if not task_id or not image_url:
            return jsonify({
                "code": 400,
                "message": "缺少必要参数"
            }), 400
        
        print(f"\n{'='*60}")
        print(f"收到识别请求 (URL):")
        print(f"  Task ID: {task_id}")
        print(f"  Model: {model_name}")
        print(f"  Image URL: {image_url}")
        print(f"  Callback: {callback}")
        print(f"{'='*60}")
        
        # 异步处理
        thread = threading.Thread(
            target=process_image,
            args=(task_id, callback, os.path.basename(image_url))
        )
        thread.start()
        
        return jsonify({
            "code": 0,
            "message": "success",
            "data": {
                "task_id": task_id,
                "status": "pending"
            }
        })
        
    except Exception as e:
        print(f"错误: {e}")
        return jsonify({
            "code": 500,
            "message": str(e)
        }), 500


@app.route('/api/v1/task/<task_id>/status', methods=['GET'])
def get_task_status(task_id):
    """
    查询任务状态（简化版，实际应该从数据库查询）
    """
    # 这里简化处理，实际应该维护一个任务状态表
    return jsonify({
        "code": 0,
        "message": "success",
        "data": {
            "task_id": task_id,
            "status": "processing",
            "progress": 50
        }
    })

@app.route('/api/v1/anomaly/analyze', methods=['POST'])
def anomaly_analyze():
    """
    异常原因分析（mock）
    约定返回结构：
      { code:0, message:"success", data:{ analysis, causes, actions, confidence } }
    """
    try:
        data = request.json or {}
        metric = data.get("metric") or "指标"
        series_name = data.get("seriesName") or ""
        x_value = data.get("xValue") or ""
        y_value = data.get("yValue")
        series_values = data.get("seriesValues") or []
        confidence = 0.55

        # 简单统计（可用就用，不可用就返回通用建议）
        analysis_lines = []
        if series_name:
            analysis_lines.append(f"系列：{series_name}")
        if x_value:
            analysis_lines.append(f"位置：{x_value}")
        if y_value is not None:
            analysis_lines.append(f"{metric}：{y_value}")

        if isinstance(series_values, list) and len(series_values) >= 5 and y_value is not None:
            try:
                others = [v for v in series_values if isinstance(v, (int, float))]
                mean = statistics.mean(others) if others else 0
                stdev = statistics.pstdev(others) if len(others) > 1 else 0
                z = (float(y_value) - mean) / (stdev if stdev > 1e-9 else 1.0)
                analysis_lines.append(f"统计：均值≈{mean:.4g}，标准差≈{stdev:.4g}，z≈{z:.2f}")
                confidence = 0.75 if abs(z) >= 3 else (0.65 if abs(z) >= 2 else 0.55)
            except Exception:
                pass

        analysis = "；".join(analysis_lines) + "。建议结合业务变更与数据链路进行排查。"
        causes = [
            "业务侧：活动/投放/价格/库存变化",
            "系统侧：下单/支付/搜索等链路异常导致指标突变",
            "数据侧：延迟补数、重复上报、口径变更导致尖峰或塌陷"
        ]
        actions = [
            "拆解到渠道/品类/店铺/地区，定位贡献最大的子维度",
            "对比相关指标（访客→加购→下单→支付）确认是否全链路一致",
            "检查采集/ETL任务日志、去重与过滤规则是否变更或重跑"
        ]

        return jsonify({
            "code": 0,
            "message": "success",
            "data": {
                "analysis": analysis,
                "causes": causes,
                "actions": actions,
                "confidence": confidence
            }
        })
    except Exception as e:
        return jsonify({"code": 500, "message": str(e)}), 500


@app.route('/health', methods=['GET'])
def health():
    """健康检查"""
    return jsonify({"status": "ok"})


def print_banner():
    """打印欢迎信息"""
    print("\n" + "="*60)
    print("  Mock 模型服务")
    print("="*60)
    print(f"\n服务地址: http://0.0.0.0:8000")
    print(f"回调地址: {CALLBACK_URL}")
    print(f"处理时间: {PROCESSING_TIME}秒")
    print("\n支持的接口:")
    print("  POST /api/v1/recognize/upload  - 上传图片识别")
    print("  POST /api/v1/recognize/url     - URL图片识别")
    print("  GET  /api/v1/task/:id/status   - 查询任务状态")
    print("  POST /api/v1/anomaly/analyze   - 异常原因分析（mock）")
    print("  GET  /health                    - 健康检查")
    print("\n" + "="*60)
    print("服务已启动，等待请求...\n")


if __name__ == '__main__':
    print_banner()
    app.run(host='0.0.0.0', port=8000, debug=False)
