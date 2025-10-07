# 图片分类服务使用指南

## 概述

本项目实现了一个完整的图片分类中间层服务，支持图片上传、模型推理和分类结果存储。

## 功能特性

- ✅ 图片上传和存储
- ✅ 自动图片分类（支持本地和远程模型）
- ✅ 分类结果持久化存储
- ✅ 图片信息查询
- ✅ 图片列表查询

## 架构设计

### 核心组件

1. **Models（模型层）**
   - `Image`: 图片信息模型
   - `ImageClassification`: 图片分类标签模型

2. **Repository（仓储层）**
   - `ImageRepository`: 图片数据访问层，提供CRUD操作

3. **Service（服务层）**
   - `ModelService`: 模型服务接口
   - `LocalModelService`: 本地模型服务实现
   - `RemoteModelService`: 远程模型服务实现

4. **Logic（业务逻辑层）**
   - `UploadImageLogic`: 图片上传和分类业务逻辑
   - `GetImageLogic`: 获取图片信息业务逻辑
   - `ListImagesLogic`: 图片列表查询业务逻辑

5. **Handler（处理器层）**
   - `UploadImageHandler`: 处理图片上传请求
   - `GetImageHandler`: 处理图片查询请求
   - `ListImagesHandler`: 处理图片列表请求

## 数据库设计

### images 表

| 字段名 | 类型 | 说明 |
|--------|------|------|
| id | BIGSERIAL | 主键 |
| filename | VARCHAR(255) | 原始文件名 |
| file_path | VARCHAR(500) | 文件存储路径 |
| file_size | BIGINT | 文件大小(字节) |
| mime_type | VARCHAR(100) | 文件MIME类型 |
| width | INTEGER | 图片宽度 |
| height | INTEGER | 图片高度 |
| uploaded_by | BIGINT | 上传用户ID |
| status | SMALLINT | 状态：0-处理中 1-已分类 2-失败 |
| created_at | TIMESTAMP | 创建时间 |
| updated_at | TIMESTAMP | 更新时间 |

### image_classifications 表

| 字段名 | 类型 | 说明 |
|--------|------|------|
| id | BIGSERIAL | 主键 |
| image_id | BIGINT | 图片ID |
| label | VARCHAR(100) | 分类标签 |
| confidence | DECIMAL(5,4) | 置信度(0-1) |
| model_name | VARCHAR(100) | 使用的模型名称 |
| model_version | VARCHAR(50) | 模型版本 |
| created_at | TIMESTAMP | 创建时间 |

## 配置说明

在 `backend/etc/backend-api.yaml` 中添加模型配置：

```yaml
Model:
  Name: image-classifier        # 模型名称
  Version: v1.0                 # 模型版本
  Type: local                   # local 或 remote
  Path: /path/to/model         # 本地模型路径或远程服务端点
```

### 本地模型配置示例

```yaml
Model:
  Name: resnet50
  Version: v1.0
  Type: local
  Path: /models/resnet50.onnx
```

### 远程模型配置示例

```yaml
Model:
  Name: image-classifier
  Version: v1.0
  Type: remote
  Path: http://ml-service:5000/classify
```

## API 接口

### 1. 上传图片并分类

**接口**: `POST /api/images/upload`

**请求**:
- Content-Type: `multipart/form-data`
- 字段: `image` (文件)

**响应**:
```json
{
  "image_id": 1,
  "filename": "cat.jpg",
  "file_path": "./uploads/1699999999999.jpg",
  "file_size": 102400,
  "status": 1,
  "classifications": [
    {
      "label": "cat",
      "confidence": 0.8523
    },
    {
      "label": "dog",
      "confidence": 0.1234
    }
  ]
}
```

### 2. 获取图片信息

**接口**: `GET /api/images/:id`

**响应**:
```json
{
  "id": 1,
  "filename": "cat.jpg",
  "file_path": "./uploads/1699999999999.jpg",
  "file_size": 102400,
  "mime_type": "image/jpeg",
  "width": 1920,
  "height": 1080,
  "uploaded_by": 1,
  "status": 1,
  "classifications": [
    {
      "label": "cat",
      "confidence": 0.8523
    }
  ],
  "created_at": "2024-01-01 12:00:00",
  "updated_at": "2024-01-01 12:00:00"
}
```

### 3. 获取图片列表

**接口**: `GET /api/images?page=1&page_size=10`

**参数**:
- `page`: 页码（默认1）
- `page_size`: 每页数量（默认10）

**响应**:
```json
{
  "total": 100,
  "list": [
    {
      "id": 1,
      "filename": "cat.jpg",
      "file_path": "./uploads/1699999999999.jpg",
      "file_size": 102400,
      "mime_type": "image/jpeg",
      "status": 1,
      "classifications": [
        {
          "label": "cat",
          "confidence": 0.8523
        }
      ],
      "created_at": "2024-01-01 12:00:00",
      "updated_at": "2024-01-01 12:00:00"
    }
  ]
}
```

## 模型服务实现

### 本地模型服务

`LocalModelService` 提供了本地模型推理的框架，目前使用模拟数据。

**集成真实模型的方式**:

1. **使用 TensorFlow Go**
```go
import tf "github.com/tensorflow/tensorflow/tensorflow/go"
// 在 ClassifyImageFromBytes 中加载和运行模型
```

2. **使用 ONNX Runtime**
```go
import "github.com/yalue/onnxruntime_go"
// 在 ClassifyImageFromBytes 中加载和运行ONNX模型
```

3. **通过 CGO 调用 Python**
```go
// 使用 cgo 调用 Python 模型
```

### 远程模型服务

`RemoteModelService` 通过HTTP调用远程模型服务。

**远程服务需要实现的接口**:

```
POST /classify
Content-Type: multipart/form-data
Body: image=<binary>

Response:
{
  "results": [
    {
      "label": "cat",
      "confidence": 0.8523
    }
  ]
}
```

## 使用流程

### 1. 运行数据库迁移

```bash
cd migrations
chmod +x run_migrations.sh
./run_migrations.sh
```

### 2. 配置模型服务

修改 `backend/etc/backend-api.yaml` 中的 Model 配置。

### 3. 启动服务

```bash
cd backend
go run backend.go
```

### 4. 上传图片

```bash
curl -X POST http://localhost:8888/api/images/upload \
  -F "image=@/path/to/your/image.jpg"
```

## 工作流程

1. **图片上传**: 用户通过 `/api/images/upload` 接口上传图片
2. **文件存储**: 图片保存到 `./uploads` 目录
3. **数据库记录**: 创建图片记录，状态设为"处理中"
4. **模型推理**: 调用模型服务进行图片分类
5. **保存分类结果**: 将分类标签和置信度保存到数据库
6. **更新状态**: 更新图片状态为"已分类"
7. **返回结果**: 返回图片信息和分类结果

## 扩展性

### 添加新的分类模型

1. 实现 `ModelService` 接口
2. 在 `ServiceContext` 初始化时注册新服务
3. 通过配置文件选择使用的模型

### 支持多模型并行

修改业务逻辑，支持同时使用多个模型进行分类：

```go
type ServiceContext struct {
    ModelServices []ModelService
}
```

### 异步处理

对于大文件或慢模型，可以实现异步处理：

1. 上传后立即返回，状态为"处理中"
2. 使用消息队列异步处理分类
3. 用户通过查询接口获取最新状态

## 性能优化

1. **批量推理**: 积累多个图片后批量调用模型
2. **结果缓存**: 使用 Redis 缓存相似图片的分类结果
3. **模型预加载**: 服务启动时预加载模型到内存
4. **并发处理**: 使用 goroutine 池并发处理多个请求

## 注意事项

1. **文件大小限制**: 当前限制为32MB，可在 Handler 中调整
2. **存储空间**: 确保 `./uploads` 目录有足够空间
3. **模型性能**: 复杂模型可能需要较长推理时间
4. **错误处理**: 分类失败时图片状态会标记为"失败"
5. **安全性**: 建议添加文件类型验证和病毒扫描

## 故障排查

### 图片上传失败

- 检查文件大小是否超过限制
- 检查 `./uploads` 目录权限
- 检查磁盘空间

### 分类失败

- 检查模型配置是否正确
- 查看日志中的错误信息
- 验证模型服务是否可访问（远程模式）

### 数据库错误

- 检查数据库连接配置
- 确认已运行数据库迁移
- 检查数据库权限

## 后续开发建议

1. **集成真实ML模型**: 替换模拟数据为实际模型推理
2. **图片预处理**: 添加图片压缩、格式转换等功能
3. **分类结果过滤**: 根据置信度阈值过滤结果
4. **用户权限管理**: 与现有RBAC系统集成
5. **审计日志**: 记录图片上传和分类操作
6. **图片下载**: 提供图片下载接口
7. **批量操作**: 支持批量上传和分类
8. **统计分析**: 提供分类结果统计接口