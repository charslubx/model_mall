#!/bin/bash

# 完整系统启动脚本

set -e

echo "========================================="
echo "图片分类系统 - 完整启动"
echo "========================================="

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 检查Docker
echo -e "\n${YELLOW}[1/6] 检查Docker环境...${NC}"
if ! command -v docker &> /dev/null; then
    echo -e "${RED}错误: Docker未安装${NC}"
    exit 1
fi

if ! command -v docker-compose &> /dev/null; then
    echo -e "${RED}错误: docker-compose未安装${NC}"
    exit 1
fi
echo -e "${GREEN}✓ Docker环境正常${NC}"

# 检查模型文件
echo -e "\n${YELLOW}[2/6] 检查模型文件...${NC}"
if [ ! -d "model_service/models" ]; then
    mkdir -p model_service/models
fi

model_count=$(find model_service/models -type f ! -name ".*" ! -name "*.txt" | wc -l)
if [ "$model_count" -eq 0 ]; then
    echo -e "${RED}警告: 未找到模型文件${NC}"
    echo "请将训练好的模型文件放到 model_service/models/ 目录"
    echo "例如: cp /path/to/your/model.mph model_service/models/"
    read -p "是否继续？(y/n) " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        exit 1
    fi
else
    echo -e "${GREEN}✓ 找到 $model_count 个模型文件${NC}"
fi

# 启动基础服务（数据库、Redis）
echo -e "\n${YELLOW}[3/6] 启动基础服务（PostgreSQL、Redis）...${NC}"
docker-compose up -d postgres redis

# 等待数据库启动
echo -e "${YELLOW}等待数据库启动...${NC}"
sleep 5

# 检查数据库连接
max_retries=30
retry_count=0
while ! docker exec postgres pg_isready -U postgres > /dev/null 2>&1; do
    retry_count=$((retry_count+1))
    if [ $retry_count -ge $max_retries ]; then
        echo -e "${RED}错误: 数据库启动超时${NC}"
        exit 1
    fi
    echo "等待数据库... ($retry_count/$max_retries)"
    sleep 2
done
echo -e "${GREEN}✓ 数据库已就绪${NC}"

# 运行数据库迁移
echo -e "\n${YELLOW}[4/6] 运行数据库迁移...${NC}"
if [ -f "migrations/run_migrations.sh" ]; then
    cd migrations
    chmod +x run_migrations.sh
    ./run_migrations.sh || echo -e "${YELLOW}警告: 数据库迁移可能失败，请检查${NC}"
    cd ..
else
    echo -e "${YELLOW}警告: 未找到迁移脚本${NC}"
fi
echo -e "${GREEN}✓ 数据库迁移完成${NC}"

# 启动模型服务
echo -e "\n${YELLOW}[5/6] 启动模型服务...${NC}"
docker-compose up -d model-service

# 等待模型服务启动
echo -e "${YELLOW}等待模型服务启动...${NC}"
max_retries=60
retry_count=0
while ! curl -f http://localhost:5000/health > /dev/null 2>&1; do
    retry_count=$((retry_count+1))
    if [ $retry_count -ge $max_retries ]; then
        echo -e "${RED}错误: 模型服务启动超时${NC}"
        echo "查看日志: docker-compose logs model-service"
        exit 1
    fi
    echo "等待模型服务... ($retry_count/$max_retries)"
    sleep 2
done
echo -e "${GREEN}✓ 模型服务已就绪${NC}"

# 启动Go后端（如果有Dockerfile）
echo -e "\n${YELLOW}[6/6] 启动Go后端...${NC}"
if [ -f "backend/Dockerfile" ]; then
    docker-compose up -d backend
    echo -e "${GREEN}✓ Go后端已启动${NC}"
else
    echo -e "${YELLOW}提示: 请手动启动Go后端${NC}"
    echo "cd backend && go run backend.go"
fi

# 显示服务状态
echo -e "\n${GREEN}=========================================${NC}"
echo -e "${GREEN}所有服务已启动！${NC}"
echo -e "${GREEN}=========================================${NC}"
echo ""
echo "服务地址："
echo "  - Go后端:      http://localhost:8888"
echo "  - 模型服务:    http://localhost:5000"
echo "  - PostgreSQL:  localhost:5432"
echo "  - Redis:       localhost:6379"
echo ""
echo "健康检查："
echo "  - 模型服务:    curl http://localhost:5000/health"
echo "  - 模型信息:    curl http://localhost:5000/info"
echo ""
echo "测试上传："
echo "  curl -X POST http://localhost:8888/api/images/upload \\"
echo "    -F \"image=@test_image.jpg\""
echo ""
echo "查看日志："
echo "  docker-compose logs -f"
echo ""
echo "停止服务："
echo "  docker-compose down"
echo ""
echo -e "${GREEN}=========================================${NC}"
