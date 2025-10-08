"""
Flask API服务
提供图片分类HTTP接口
"""
import os
import logging
from flask import Flask, request, jsonify
from werkzeug.exceptions import BadRequest
from app.model_loader import create_model_loader, ModelLoader

# 配置日志
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s'
)
logger = logging.getLogger(__name__)

# 创建Flask应用
app = Flask(__name__)
app.config['MAX_CONTENT_LENGTH'] = 32 * 1024 * 1024  # 最大32MB

# 全局模型加载器
model_loader: ModelLoader = None


def init_model():
    """初始化模型"""
    global model_loader
    
    # 从环境变量获取配置
    model_path = os.environ.get('MODEL_PATH', '/app/models/model.h5')
    model_name = os.environ.get('MODEL_NAME', 'image-classifier')
    
    logger.info(f"初始化模型服务...")
    logger.info(f"模型路径: {model_path}")
    logger.info(f"模型名称: {model_name}")
    
    try:
        # 创建模型加载器
        model_loader = create_model_loader(model_path, model_name)
        
        # 加载模型
        model_loader.load_model()
        
        logger.info("模型初始化成功")
    except Exception as e:
        logger.error(f"模型初始化失败: {str(e)}")
        raise


@app.route('/health', methods=['GET'])
def health_check():
    """健康检查接口"""
    if model_loader and model_loader.is_loaded:
        return jsonify({
            'status': 'healthy',
            'model_name': model_loader.model_name,
            'model_loaded': True
        }), 200
    else:
        return jsonify({
            'status': 'unhealthy',
            'model_loaded': False,
            'error': '模型未加载'
        }), 503


@app.route('/classify', methods=['POST'])
def classify_image():
    """
    图片分类接口
    
    请求:
        - Content-Type: multipart/form-data
        - 字段: image (文件)
    
    响应:
        {
            "success": true,
            "results": [
                {
                    "label": "cat",
                    "confidence": 0.8523
                }
            ]
        }
    """
    try:
        # 检查模型是否已加载
        if not model_loader or not model_loader.is_loaded:
            return jsonify({
                'success': False,
                'error': '模型未加载'
            }), 503
        
        # 检查是否有文件
        if 'image' not in request.files:
            return jsonify({
                'success': False,
                'error': '请上传图片文件（字段名：image）'
            }), 400
        
        file = request.files['image']
        
        # 检查文件名
        if file.filename == '':
            return jsonify({
                'success': False,
                'error': '文件名为空'
            }), 400
        
        # 检查文件类型
        allowed_extensions = {'png', 'jpg', 'jpeg', 'gif', 'bmp', 'webp'}
        ext = file.filename.rsplit('.', 1)[1].lower() if '.' in file.filename else ''
        
        if ext not in allowed_extensions:
            return jsonify({
                'success': False,
                'error': f'不支持的文件格式，支持的格式: {", ".join(allowed_extensions)}'
            }), 400
        
        # 读取图片数据
        image_data = file.read()
        
        # 检查文件大小
        if len(image_data) == 0:
            return jsonify({
                'success': False,
                'error': '文件内容为空'
            }), 400
        
        logger.info(f"收到分类请求: 文件名={file.filename}, 大小={len(image_data)} bytes")
        
        # 调用模型进行预测
        results = model_loader.predict(image_data)
        
        logger.info(f"分类完成: {len(results)} 个结果")
        
        return jsonify({
            'success': True,
            'results': results,
            'model_name': model_loader.model_name
        }), 200
        
    except BadRequest as e:
        logger.error(f"请求错误: {str(e)}")
        return jsonify({
            'success': False,
            'error': '请求格式错误'
        }), 400
    except Exception as e:
        logger.error(f"分类失败: {str(e)}", exc_info=True)
        return jsonify({
            'success': False,
            'error': f'分类失败: {str(e)}'
        }), 500


@app.route('/info', methods=['GET'])
def model_info():
    """获取模型信息"""
    if not model_loader:
        return jsonify({
            'error': '模型未初始化'
        }), 503
    
    return jsonify({
        'model_name': model_loader.model_name,
        'model_path': model_loader.model_path,
        'is_loaded': model_loader.is_loaded,
        'model_type': model_loader.__class__.__name__
    }), 200


@app.errorhandler(413)
def request_entity_too_large(error):
    """处理文件过大错误"""
    return jsonify({
        'success': False,
        'error': '文件过大，最大支持32MB'
    }), 413


@app.errorhandler(500)
def internal_server_error(error):
    """处理服务器内部错误"""
    logger.error(f"服务器错误: {str(error)}")
    return jsonify({
        'success': False,
        'error': '服务器内部错误'
    }), 500


if __name__ == '__main__':
    # 初始化模型
    init_model()
    
    # 启动服务
    port = int(os.environ.get('PORT', 5000))
    host = os.environ.get('HOST', '0.0.0.0')
    
    logger.info(f"启动Flask服务: {host}:{port}")
    app.run(host=host, port=port, debug=False)
