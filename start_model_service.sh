#!/bin/bash

# 启动示例模型服务脚本

echo "正在安装Python依赖..."

# 检查是否安装了pip
if ! command -v pip3 &> /dev/null; then
    echo "错误: 未找到pip3，请先安装Python 3和pip"
    exit 1
fi

# 安装必要的依赖
pip3 install flask pillow

echo "启动模型服务..."
python3 example_model_service.py