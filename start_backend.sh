#!/bin/bash

# 快速启动后端服务脚本

echo "正在启动模型商城后端服务..."

# 检查是否已编译
if [ ! -f "backend/backend" ]; then
    echo "编译后端服务..."
    cd /workspace && go build -o backend/backend backend/backend.go
    if [ $? -ne 0 ]; then
        echo "编译失败！"
        exit 1
    fi
fi

# 创建上传目录
mkdir -p uploads

echo "启动服务..."
cd backend && ./backend