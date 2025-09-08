#!/bin/bash

# 数据库迁移执行脚本
echo "开始执行数据库迁移..."

# 检查配置文件是否存在
if [ ! -f "../backend/etc/backend-api.yaml" ]; then
    echo "错误: 配置文件不存在，请先配置数据库连接信息"
    exit 1
fi

# 执行迁移
cd "$(dirname "$0")"
go run migrate.go

if [ $? -eq 0 ]; then
    echo "✓ 数据库迁移完成！"
    echo ""
    echo "默认用户账号:"
    echo "超级管理员 - 用户名: admin, 密码: admin123"
    echo "管理员     - 用户名: manager, 密码: admin123"  
    echo "普通用户   - 用户名: user, 密码: admin123"
    echo ""
    echo "请及时修改默认密码！"
else
    echo "✗ 数据库迁移失败！"
    exit 1
fi