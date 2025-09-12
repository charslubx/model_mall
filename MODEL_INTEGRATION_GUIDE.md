# 模型对接中间层使用指南

## 概述

本项目已成功集成了图片分类模型对接功能，包括：
- 图片上传和预处理
- 模型服务调用
- 多标签分类结果存储
- 完整的API接口

## 功能特性

### 1. 图片分类功能
- 支持多种图片格式：JPEG, PNG, GIF, WebP, BMP
- 最大文件大小：10MB
- 支持多标签分类结果
- 置信度过滤
- 边界框信息存储

### 2. 数据存储
- 图片分类记录表 (`image_classifications`)
- 分类标签表 (`image_classification_labels`)
- 支持用户关联
- 处理状态跟踪

### 3. API接口
- POST `/api/classify/image` - 图片分类
- GET `/api/classify/result` - 获取分类结果
- GET `/api/classify/my` - 获取我的分类记录
- GET `/api/classify/statistics` - 获取统计信息
- DELETE `/api/classify/result` - 删除分类记录
- GET `/api/classify/search` - 搜索分类记录

## 部署步骤

### 1. 数据库迁移

```bash
cd migrations
chmod +x run_migrations.sh
./run_migrations.sh
```

或者手动执行：

```bash
cd migrations
go run migrate.go
```

### 2. 启动模型服务（示例）

```bash
# 使用提供的示例服务
./start_model_service.sh
```

或者手动启动：

```bash
pip3 install flask pillow
python3 example_model_service.py
```

### 3. 配置后端服务

修改 `backend/etc/backend-api.yaml` 中的配置：

```yaml
ModelService:
  Endpoint: http://localhost:8080  # 模型服务地址
  Timeout: 30  # 超时时间（秒）
Upload:
  Path: ./uploads  # 文件上传路径
  MaxSize: 10485760  # 最大文件大小（10MB）
  AllowedTypes: ["jpeg", "jpg", "png", "gif", "webp", "bmp"]
```

### 4. 启动后端服务

```bash
cd backend
go run backend.go
```

## API使用示例

### 1. 图片分类

```bash
curl -X POST \
  http://localhost:8888/api/classify/image \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -F "image=@/path/to/your/image.jpg" \
  -F "model_name=default" \
  -F "save_image=true" \
  -F "min_confidence=0.5"
```

响应示例：
```json
{
  "code": 200,
  "message": "分类成功",
  "data": {
    "id": 1,
    "image_path": "2024/01/15/20240115143022_abc12345.jpg",
    "image_name": "test_image.jpg",
    "model_name": "default",
    "process_time": 234,
    "confidence": 0.95,
    "status": 1,
    "labels": [
      {
        "id": 1,
        "classification_id": 1,
        "label_name": "猫",
        "label_code": "cat",
        "confidence": 0.95,
        "bounding_box": "{\"x\":0.1,\"y\":0.2,\"width\":0.5,\"height\":0.6}",
        "created_at": "2024-01-15T14:30:22Z"
      }
    ],
    "created_at": "2024-01-15T14:30:22Z"
  }
}
```

### 2. 获取分类结果

```bash
curl -X GET \
  "http://localhost:8888/api/classify/result?id=1" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### 3. 获取我的分类记录

```bash
curl -X GET \
  "http://localhost:8888/api/classify/my?page=1&page_size=20" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### 4. 获取统计信息

```bash
curl -X GET \
  "http://localhost:8888/api/classify/statistics" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### 5. 搜索分类记录

```bash
curl -X GET \
  "http://localhost:8888/api/classify/search?model_name=default&status=1&page=1&page_size=20" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

## 数据库表结构

### image_classifications 表
| 字段名 | 类型 | 说明 |
|--------|------|------|
| id | BIGSERIAL | 主键 |
| image_path | VARCHAR(500) | 图片路径 |
| image_name | VARCHAR(255) | 图片名称 |
| image_size | BIGINT | 图片大小 |
| image_format | VARCHAR(20) | 图片格式 |
| model_name | VARCHAR(100) | 模型名称 |
| model_version | VARCHAR(50) | 模型版本 |
| process_time | BIGINT | 处理耗时(毫秒) |
| confidence | DECIMAL(5,4) | 总体置信度 |
| status | SMALLINT | 状态 0-失败 1-成功 2-处理中 |
| user_id | BIGINT | 用户ID |
| created_at | TIMESTAMP | 创建时间 |
| updated_at | TIMESTAMP | 更新时间 |

### image_classification_labels 表
| 字段名 | 类型 | 说明 |
|--------|------|------|
| id | BIGSERIAL | 主键 |
| classification_id | BIGINT | 分类记录ID |
| label_name | VARCHAR(100) | 标签名称 |
| label_code | VARCHAR(100) | 标签代码 |
| confidence | DECIMAL(5,4) | 置信度 |
| bounding_box | TEXT | 边界框信息(JSON格式) |
| created_at | TIMESTAMP | 创建时间 |

## 模型服务接口规范

### 健康检查
- 接口：GET `/health`
- 响应：`{"status": "healthy", "timestamp": 1642234567}`

### 图片分类
- 接口：POST `/predict`
- 请求：multipart/form-data
  - `image`: 图片文件
  - `model_name`: 模型名称
  - `min_confidence`: 最小置信度（可选）
- 响应：
```json
{
  "success": true,
  "message": "分类成功",
  "process_time": 234,
  "predictions": [
    {
      "label": "猫",
      "code": "cat",
      "confidence": 0.95,
      "bounding_box": {
        "x": 0.1,
        "y": 0.2,
        "width": 0.5,
        "height": 0.6
      }
    }
  ]
}
```

## 错误处理

系统包含完整的错误处理机制：
- 图片格式验证
- 文件大小限制
- 模型服务超时处理
- 数据库事务回滚
- 详细的错误日志

## 扩展性

### 1. 添加新的模型
- 更新模型服务以支持新模型
- 修改配置文件添加模型信息
- 无需修改数据库结构

### 2. 支持更多图片格式
- 更新 `ModelService.GetImageFormat()` 方法
- 修改配置文件中的 `AllowedTypes`

### 3. 添加批量处理
- 扩展API接口支持多文件上传
- 修改业务逻辑支持批量处理

## 性能优化

### 1. 数据库优化
- 已创建必要的索引
- 支持分页查询
- 使用事务确保数据一致性

### 2. 文件存储优化
- 按日期组织文件目录
- 使用MD5避免重复存储
- 可扩展支持云存储

### 3. 缓存优化
- 可集成Redis缓存热点数据
- 模型结果缓存
- 统计信息缓存

## 监控和日志

系统已集成完整的日志记录：
- 请求日志
- 错误日志
- 性能指标
- 模型调用统计

## 安全考虑

- JWT认证保护API
- 文件类型验证
- 文件大小限制
- SQL注入防护
- 用户权限控制

## 故障排除

### 常见问题

1. **模型服务连接失败**
   - 检查模型服务是否启动
   - 验证配置文件中的endpoint地址
   - 检查网络连接

2. **图片上传失败**
   - 检查文件大小是否超限
   - 验证图片格式是否支持
   - 检查上传目录权限

3. **数据库连接错误**
   - 检查数据库服务状态
   - 验证连接配置
   - 确认数据库表是否已创建

## 总结

本模型对接中间层提供了完整的图片分类功能，包括：
- ✅ 图片上传和预处理
- ✅ 模型服务调用
- ✅ 多标签分类结果存储
- ✅ 完整的RESTful API
- ✅ 用户权限控制
- ✅ 错误处理和日志记录
- ✅ 数据库表和索引优化
- ✅ 示例模型服务

系统已经可以立即投入使用，您只需要：
1. 运行数据库迁移
2. 启动模型服务
3. 配置后端服务
4. 开始使用API接口

如需自定义模型，请替换示例模型服务为您的实际模型实现。