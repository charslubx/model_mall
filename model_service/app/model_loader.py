"""
模型加载器
支持多种模型格式的加载
"""
import os
import logging
from typing import List, Dict, Any
import numpy as np
from PIL import Image
import io

logger = logging.getLogger(__name__)


class ModelLoader:
    """模型加载器基类"""
    
    def __init__(self, model_path: str, model_name: str = "image-classifier"):
        self.model_path = model_path
        self.model_name = model_name
        self.model = None
        self.is_loaded = False
    
    def load_model(self):
        """加载模型"""
        raise NotImplementedError
    
    def predict(self, image_data: bytes) -> List[Dict[str, Any]]:
        """预测图片分类"""
        raise NotImplementedError
    
    def preprocess_image(self, image_data: bytes, target_size=(224, 224)) -> np.ndarray:
        """预处理图片"""
        try:
            # 将字节数据转换为PIL Image
            image = Image.open(io.BytesIO(image_data))
            
            # 转换为RGB（如果是RGBA或其他格式）
            if image.mode != 'RGB':
                image = image.convert('RGB')
            
            # 调整大小
            image = image.resize(target_size)
            
            # 转换为numpy数组
            img_array = np.array(image, dtype=np.float32)
            
            # 归一化到 [0, 1]
            img_array = img_array / 255.0
            
            # 添加batch维度
            img_array = np.expand_dims(img_array, axis=0)
            
            return img_array
        except Exception as e:
            logger.error(f"图片预处理失败: {str(e)}")
            raise


class KerasModelLoader(ModelLoader):
    """Keras/TensorFlow模型加载器（支持.h5格式）"""
    
    def load_model(self):
        """加载Keras模型"""
        try:
            import tensorflow as tf
            logger.info(f"正在加载Keras模型: {self.model_path}")
            self.model = tf.keras.models.load_model(self.model_path)
            self.is_loaded = True
            logger.info("模型加载成功")
        except Exception as e:
            logger.error(f"加载Keras模型失败: {str(e)}")
            raise
    
    def predict(self, image_data: bytes) -> List[Dict[str, Any]]:
        """预测图片分类"""
        if not self.is_loaded:
            raise RuntimeError("模型未加载")
        
        try:
            # 预处理图片
            img_array = self.preprocess_image(image_data)
            
            # 模型推理
            predictions = self.model.predict(img_array, verbose=0)
            
            # 解析预测结果
            results = self._parse_predictions(predictions[0])
            
            return results
        except Exception as e:
            logger.error(f"预测失败: {str(e)}")
            raise
    
    def _parse_predictions(self, predictions: np.ndarray, top_k: int = 5) -> List[Dict[str, Any]]:
        """解析预测结果"""
        # 获取top-k结果
        top_indices = np.argsort(predictions)[-top_k:][::-1]
        
        results = []
        for idx in top_indices:
            results.append({
                "label": self._get_label_name(int(idx)),
                "confidence": float(predictions[idx])
            })
        
        return results
    
    def _get_label_name(self, class_idx: int) -> str:
        """获取标签名称"""
        # 这里需要根据实际的类别映射来返回标签名
        # 可以从外部文件加载类别映射
        labels_file = os.path.join(os.path.dirname(self.model_path), 'labels.txt')
        
        if os.path.exists(labels_file):
            with open(labels_file, 'r', encoding='utf-8') as f:
                labels = [line.strip() for line in f.readlines()]
                if class_idx < len(labels):
                    return labels[class_idx]
        
        # 如果没有标签文件，返回类别索引
        return f"class_{class_idx}"


class ONNXModelLoader(ModelLoader):
    """ONNX模型加载器"""
    
    def load_model(self):
        """加载ONNX模型"""
        try:
            import onnxruntime as ort
            logger.info(f"正在加载ONNX模型: {self.model_path}")
            self.model = ort.InferenceSession(self.model_path)
            self.is_loaded = True
            logger.info("模型加载成功")
        except Exception as e:
            logger.error(f"加载ONNX模型失败: {str(e)}")
            raise
    
    def predict(self, image_data: bytes) -> List[Dict[str, Any]]:
        """预测图片分类"""
        if not self.is_loaded:
            raise RuntimeError("模型未加载")
        
        try:
            # 预处理图片
            img_array = self.preprocess_image(image_data)
            
            # 获取输入名称
            input_name = self.model.get_inputs()[0].name
            
            # 模型推理
            predictions = self.model.run(None, {input_name: img_array})
            
            # 解析预测结果
            results = self._parse_predictions(predictions[0][0])
            
            return results
        except Exception as e:
            logger.error(f"预测失败: {str(e)}")
            raise
    
    def _parse_predictions(self, predictions: np.ndarray, top_k: int = 5) -> List[Dict[str, Any]]:
        """解析预测结果"""
        # 获取top-k结果
        top_indices = np.argsort(predictions)[-top_k:][::-1]
        
        results = []
        for idx in top_indices:
            results.append({
                "label": self._get_label_name(int(idx)),
                "confidence": float(predictions[idx])
            })
        
        return results
    
    def _get_label_name(self, class_idx: int) -> str:
        """获取标签名称"""
        labels_file = os.path.join(os.path.dirname(self.model_path), 'labels.txt')
        
        if os.path.exists(labels_file):
            with open(labels_file, 'r', encoding='utf-8') as f:
                labels = [line.strip() for line in f.readlines()]
                if class_idx < len(labels):
                    return labels[class_idx]
        
        return f"class_{class_idx}"


class PyTorchModelLoader(ModelLoader):
    """PyTorch模型加载器（完整模型）"""
    
    def __init__(self, model_path: str, model_name: str = "image-classifier", device: str = None):
        super().__init__(model_path, model_name)
        # 自动选择设备：优先使用GPU
        if device is None:
            import torch
            self.device = torch.device('cuda' if torch.cuda.is_available() else 'cpu')
        else:
            import torch
            self.device = torch.device(device)
        logger.info(f"使用设备: {self.device}")
    
    def load_model(self):
        """加载PyTorch完整模型"""
        try:
            import torch
            logger.info(f"正在加载PyTorch模型: {self.model_path}")
            
            # 加载完整模型
            self.model = torch.load(self.model_path, map_location=self.device)
            
            # 检查是否是nn.Module实例
            if not isinstance(self.model, torch.nn.Module):
                raise RuntimeError(
                    f"加载的不是有效的PyTorch模型。\n"
                    f"类型: {type(self.model)}\n"
                    f"请确保使用 torch.save(model, 'model.pth') 保存完整模型"
                )
            
            logger.info(f"成功加载模型，类型: {type(self.model).__name__}")
            
            # 设置为评估模式
            self.model.eval()
            logger.info("已设置为评估模式")
            
            # 移动到指定设备
            self.model.to(self.device)
            
            self.is_loaded = True
            logger.info("PyTorch模型加载成功")
            
        except Exception as e:
            logger.error(f"加载PyTorch模型失败: {str(e)}")
            logger.error("提示: 请确保模型是使用 torch.save(model, path) 保存的完整模型")
            raise
    
    def predict(self, image_data: bytes) -> List[Dict[str, Any]]:
        """预测图片分类"""
        if not self.is_loaded:
            raise RuntimeError("模型未加载")
        
        try:
            import torch
            
            # 预处理图片
            img_array = self.preprocess_image(image_data)
            
            # 转换为PyTorch tensor
            img_tensor = torch.from_numpy(img_array).permute(0, 3, 1, 2)  # NHWC -> NCHW
            img_tensor = img_tensor.to(self.device)
            
            # 模型推理
            with torch.no_grad():
                outputs = self.model(img_tensor)
            
            # 如果输出是softmax前的logits，应用softmax
            if outputs.dim() == 2:
                # 检查是否需要softmax
                probabilities = torch.softmax(outputs, dim=1)
            else:
                probabilities = outputs
            
            # 转换为numpy数组
            predictions = probabilities.cpu().numpy()[0]
            
            # 解析预测结果
            results = self._parse_predictions(predictions)
            
            return results
            
        except Exception as e:
            logger.error(f"预测失败: {str(e)}")
            raise
    
    def _parse_predictions(self, predictions: np.ndarray, top_k: int = 5) -> List[Dict[str, Any]]:
        """解析预测结果"""
        # 获取top-k结果
        top_indices = np.argsort(predictions)[-top_k:][::-1]
        
        results = []
        for idx in top_indices:
            results.append({
                "label": self._get_label_name(int(idx)),
                "confidence": float(predictions[idx])
            })
        
        return results
    
    def _get_label_name(self, class_idx: int) -> str:
        """获取标签名称"""
        labels_file = os.path.join(os.path.dirname(self.model_path), 'labels.txt')
        
        if os.path.exists(labels_file):
            with open(labels_file, 'r', encoding='utf-8') as f:
                labels = [line.strip() for line in f.readlines()]
                if class_idx < len(labels):
                    return labels[class_idx]
        
        return f"class_{class_idx}"


class CustomModelLoader(ModelLoader):
    """自定义模型加载器（用于.mph或其他自定义格式）"""
    
    def load_model(self):
        """加载自定义格式模型"""
        try:
            import pickle
            logger.info(f"正在加载自定义模型: {self.model_path}")
            
            # 尝试使用pickle加载
            with open(self.model_path, 'rb') as f:
                self.model = pickle.load(f)
            
            self.is_loaded = True
            logger.info("模型加载成功")
        except Exception as e:
            logger.error(f"加载自定义模型失败: {str(e)}")
            logger.info("提示: 如果是特殊格式，请修改此方法的加载逻辑")
            raise
    
    def predict(self, image_data: bytes) -> List[Dict[str, Any]]:
        """预测图片分类"""
        if not self.is_loaded:
            raise RuntimeError("模型未加载")
        
        try:
            # 预处理图片
            img_array = self.preprocess_image(image_data)
            
            # 根据实际模型调整推理方式
            # 这里需要根据具体模型类型修改
            if hasattr(self.model, 'predict'):
                predictions = self.model.predict(img_array)
            elif hasattr(self.model, 'forward'):
                predictions = self.model.forward(img_array)
            else:
                raise RuntimeError("模型没有predict或forward方法")
            
            # 解析预测结果
            results = self._parse_predictions(predictions)
            
            return results
        except Exception as e:
            logger.error(f"预测失败: {str(e)}")
            raise
    
    def _parse_predictions(self, predictions, top_k: int = 5) -> List[Dict[str, Any]]:
        """解析预测结果"""
        # 转换为numpy数组
        if not isinstance(predictions, np.ndarray):
            predictions = np.array(predictions)
        
        # 如果是多维数组，取第一个
        if len(predictions.shape) > 1:
            predictions = predictions[0]
        
        # 获取top-k结果
        top_indices = np.argsort(predictions)[-top_k:][::-1]
        
        results = []
        for idx in top_indices:
            results.append({
                "label": self._get_label_name(int(idx)),
                "confidence": float(predictions[idx])
            })
        
        return results
    
    def _get_label_name(self, class_idx: int) -> str:
        """获取标签名称"""
        labels_file = os.path.join(os.path.dirname(self.model_path), 'labels.txt')
        
        if os.path.exists(labels_file):
            with open(labels_file, 'r', encoding='utf-8') as f:
                labels = [line.strip() for line in f.readlines()]
                if class_idx < len(labels):
                    return labels[class_idx]
        
        return f"class_{class_idx}"


def create_model_loader(model_path: str, model_name: str = "image-classifier", device: str = None) -> ModelLoader:
    """
    根据模型文件扩展名创建对应的模型加载器
    
    Args:
        model_path: 模型文件路径
        model_name: 模型名称
    
    Returns:
        ModelLoader实例
    """
    if not os.path.exists(model_path):
        raise FileNotFoundError(f"模型文件不存在: {model_path}")
    
    ext = os.path.splitext(model_path)[1].lower()
    
    if ext in ['.h5', '.keras']:
        logger.info("检测到Keras模型格式")
        return KerasModelLoader(model_path, model_name)
    elif ext == '.onnx':
        logger.info("检测到ONNX模型格式")
        return ONNXModelLoader(model_path, model_name)
    elif ext in ['.pt', '.pth', '.mph']:
        logger.info(f"检测到PyTorch模型格式: {ext}")
        return PyTorchModelLoader(model_path, model_name, device)
    elif ext in ['.pkl', '.pickle']:
        logger.info(f"检测到pickle模型格式: {ext}")
        logger.warning("需要根据实际模型格式修改CustomModelLoader的加载逻辑")
        return CustomModelLoader(model_path, model_name)
    else:
        logger.warning(f"未知模型格式: {ext}, 尝试使用PyTorch加载器")
        return PyTorchModelLoader(model_path, model_name, device)
