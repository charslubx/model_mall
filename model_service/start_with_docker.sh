#!/bin/bash

# 使用Docker启动模型服务

set -e

echo "========================================="
echo "使用Docker启动图片分类模型服务"
echo "========================================="

# 检查Docker是否安装
if ! command -v docker &> /dev/null; then
    echo "错误: Docker未安装，请先安装Docker"
    exit 1
fi

# 检查docker-compose是否安装
if ! command -v docker-compose &> /dev/null; then
    echo "错误: docker-compose未安装，请先安装docker-compose"
    exit 1
fi

# 检查模型文件
if [ ! -d "models" ]; then
    echo "创建models目录..."
    mkdir -p models
fi

echo "检查模型文件..."
model_count=$(find models -type f ! -name ".*" ! -name "*.txt" | wc -l)
if [ "$model_count" -eq 0 ]; then
    echo "警告: models目录下没有模型文件"
    echo "请将训练好的模型文件（如 model.h5, model.mph 等）放在 models/ 目录下"
    read -p "是否继续？(y/n) " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        exit 1
    fi
fi

# 构建并启动服务
echo "构建Docker镜像..."
docker-compose build

echo "启动服务..."
docker-compose up -d

echo "========================================="
echo "服务已启动！"
echo "服务地址: http://localhost:5000"
echo "健康检查: http://localhost:5000/health"
echo "模型信息: http://localhost:5000/info"
echo "========================================="
echo ""
echo "查看日志: docker-compose logs -f"
echo "停止服务: docker-compose down"
echo "========================================="
