#!/bin/bash

# 模型服务启动脚本

set -e

echo "========================================="
echo "图片分类模型服务启动脚本"
echo "========================================="

# 检查Python版本
echo "检查Python环境..."
python3 --version

# 检查是否存在虚拟环境
if [ ! -d "venv" ]; then
    echo "创建Python虚拟环境..."
    python3 -m venv venv
fi

# 激活虚拟环境
echo "激活虚拟环境..."
source venv/bin/activate

# 安装依赖
echo "安装Python依赖..."
pip install -r requirements.txt -i https://pypi.tuna.tsinghua.edu.cn/simple

# 检查模型文件
echo "检查模型文件..."
if [ ! -d "models" ]; then
    echo "创建models目录..."
    mkdir -p models
fi

# 设置环境变量
export MODEL_PATH=${MODEL_PATH:-"models/model.h5"}
export MODEL_NAME=${MODEL_NAME:-"image-classifier"}
export PORT=${PORT:-5000}
export HOST=${HOST:-"0.0.0.0"}

echo "========================================="
echo "服务配置:"
echo "模型路径: $MODEL_PATH"
echo "模型名称: $MODEL_NAME"
echo "服务地址: $HOST:$PORT"
echo "========================================="

# 检查模型文件是否存在
if [ ! -f "$MODEL_PATH" ]; then
    echo "警告: 模型文件不存在: $MODEL_PATH"
    echo "请将训练好的模型文件放在 models/ 目录下"
    read -p "是否继续启动服务？(y/n) " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        exit 1
    fi
fi

# 启动服务
echo "启动Flask服务..."
python -m app.api
