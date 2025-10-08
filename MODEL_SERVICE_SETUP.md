# 模型服务集成 - 快速开始指南

## 功能概述

本次更新添加了与模型服务交互的完整功能：

1. ✅ **图片上传** - 支持多种格式图片上传
2. ✅ **模型识别** - 自动调用模型服务进行图像识别
3. ✅ **进度跟踪** - 实时查询任务处理进度
4. ✅ **标签管理** - 存储和查询分类标签
5. ✅ **回调机制** - 模型服务完成后自动回调
6. ✅ **数据持久化** - 所有数据存储在数据库，无需等待模型

## 文件结构

```
backend/
├── internal/
│   ├── models/                    # 数据模型
│   │   ├── image.go              # 图片模型
│   │   ├── recognition_task.go   # 识别任务模型
│   │   └── classification_label.go # 分类标签模型
│   ├── repository/               # 数据访问层
│   │   ├── image_repository.go
│   │   ├── recognition_task_repository.go
│   │   └── classification_label_repository.go
│   ├── handler/                  # HTTP处理器
│   │   ├── ImageHandler.go
│   │   ├── RecognitionTaskHandler.go
│   │   ├── ClassificationLabelHandler.go
│   │   └── ModelCallbackHandler.go
│   ├── logic/                    # 业务逻辑
│   │   ├── ImageLogic.go
│   │   ├── RecognitionTaskLogic.go
│   │   ├── ClassificationLabelLogic.go
│   │   └── ModelCallbackLogic.go
│   ├── svc/
│   │   └── ModelServiceClient.go # 模型服务客户端
│   └── types/
│       └── types.go              # 类型定义（已更新）
├── etc/
│   └── backend-api.yaml          # 配置文件（已更新）
migrations/
├── 005_create_images_table.sql
├── 006_create_recognition_tasks_table.sql
└── 007_create_classification_labels_table.sql
```

## 安装步骤

### 1. 数据库迁移

执行新的数据库迁移文件来创建所需的表：

```bash
cd /workspace/migrations
chmod +x run_migrations.sh
./run_migrations.sh
```

或者手动执行：

```bash
cd /workspace/migrations
go run migrate.go
```

### 2. 配置文件

编辑 `backend/etc/backend-api.yaml`，配置模型服务和上传设置：

```yaml
ModelService:
  BaseURL: http://localhost:8000  # 修改为实际的模型服务地址
  APIKey: ""                       # 如果需要，填入API密钥
  Timeout: 300                     # 超时时间（秒）

Upload:
  MaxSize: 10485760               # 最大文件大小 10MB
  AllowedTypes: ".jpg,.jpeg,.png,.gif,.bmp,.webp"
  StoragePath: "./uploads"        # 本地存储路径
  BaseURL: "http://localhost:8888/uploads"  # 访问基础URL
```

### 3. 创建上传目录

```bash
mkdir -p /workspace/uploads
chmod 755 /workspace/uploads
```

### 4. 安装依赖

如果有新的依赖需要安装：

```bash
cd /workspace
go mod tidy
```

### 5. 启动服务

```bash
cd /workspace/backend
go run backend.go -f etc/backend-api.yaml
```

## API使用示例

### 场景1：上传图片并等待识别完成

```bash
# 1. 上传图片
curl -X POST http://localhost:8888/api/images/upload \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -F "image=@/path/to/photo.jpg" \
  -F "model_name=resnet50"

# 响应：
# {
#   "code": 0,
#   "message": "success",
#   "data": {
#     "image_id": 123,
#     "task_id": "550e8400-e29b-41d4-a716-446655440000",
#     "file_url": "http://localhost:8888/uploads/2024/01/15/abc.jpg",
#     "filename": "abc.jpg"
#   }
# }

# 2. 查询任务状态（轮询）
curl -X GET http://localhost:8888/api/tasks/550e8400-e29b-41d4-a716-446655440000/status \
  -H "Authorization: Bearer YOUR_TOKEN"

# 3. 获取识别结果
curl -X GET http://localhost:8888/api/images/123/labels \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### 场景2：直接获取标签（无需轮询）

```bash
# 1. 上传图片
curl -X POST http://localhost:8888/api/images/upload \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -F "image=@photo.jpg"

# 获得 image_id: 123

# 2. 等待一段时间后，直接查询标签
# （模型完成后会自动保存到数据库）
curl -X GET http://localhost:8888/api/images/123/labels \
  -H "Authorization: Bearer YOUR_TOKEN"

# 如果返回空数组，说明还在处理中
# 如果有数据，说明识别已完成
```

### 场景3：批量查看所有任务

```bash
# 获取任务列表
curl -X GET "http://localhost:8888/api/tasks?page=1&page_size=20" \
  -H "Authorization: Bearer YOUR_TOKEN"

# 获取图片列表
curl -X GET "http://localhost:8888/api/images?page=1&page_size=20" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

## 模型服务对接

### 模型服务需要实现的接口

模型服务需要提供以下接口：

#### 1. 接收图片识别请求

**方式A：上传文件**
```
POST /api/v1/recognize/upload
Content-Type: multipart/form-data

参数：
- image: 图片文件
- task_id: 任务ID
- model_name: 模型名称（可选）
- callback: 回调URL
```

**方式B：提供URL**
```
POST /api/v1/recognize/url
Content-Type: application/json

{
  "task_id": "xxx",
  "image_url": "http://...",
  "model_name": "resnet50",
  "callback": "http://localhost:8888/api/model/callback"
}
```

#### 2. 处理完成后回调

模型服务在识别完成后，需要调用回调URL：

```bash
curl -X POST http://localhost:8888/api/model/callback \
  -H "Content-Type: application/json" \
  -d '{
    "task_id": "550e8400-e29b-41d4-a716-446655440000",
    "status": "completed",
    "progress": 100,
    "results": [
      {
        "name": "猫",
        "code": "cat",
        "confidence": 0.9856,
        "bbox": {
          "x": 100,
          "y": 150,
          "width": 200,
          "height": 180
        }
      }
    ]
  }'
```

### 回调状态说明

| 状态 | 说明 |
|------|------|
| pending | 待处理 |
| processing | 处理中（可多次回调更新进度） |
| completed | 已完成（需包含results） |
| failed | 失败（需包含error信息） |

## 数据流程图

```
用户端                    本服务                    模型服务
  │                        │                          │
  ├─► 上传图片              │                          │
  │   POST /images/upload  │                          │
  │                        ├─► 保存图片到数据库        │
  │                        ├─► 创建识别任务            │
  │                        ├─► 异步调用模型服务 ───────►
  │◄─ 返回 image_id,       │                          │
  │   task_id              │                          ├─► 开始处理
  │                        │                          │
  ├─► 查询任务状态          │                          │
  │   GET /tasks/:id/status│                          │
  │◄─ 返回进度信息          │                          │
  │                        │                          ├─► 处理中...
  │                        │                          │
  │                        │◄─ 回调更新进度 ──────────┤
  │                        ├─► 更新任务进度            │
  │                        │                          │
  │                        │◄─ 回调识别结果 ──────────┤
  │                        ├─► 保存标签到数据库        │
  │                        ├─► 更新任务状态为完成      │
  │                        │                          │
  ├─► 获取标签              │                          │
  │   GET /images/:id/labels                          │
  │◄─ 返回标签列表          │                          │
  │   （从数据库直接读取）  │                          │
```

## 测试清单

### 测试步骤

- [ ] 1. 执行数据库迁移
- [ ] 2. 配置模型服务地址
- [ ] 3. 创建上传目录
- [ ] 4. 启动服务
- [ ] 5. 测试图片上传
- [ ] 6. 测试任务查询
- [ ] 7. 模拟模型回调
- [ ] 8. 测试标签查询
- [ ] 9. 测试列表接口

### 手动测试回调（模拟模型服务）

```bash
# 模拟模型服务回调
curl -X POST http://localhost:8888/api/model/callback \
  -H "Content-Type: application/json" \
  -d '{
    "task_id": "YOUR_TASK_ID",
    "status": "completed",
    "progress": 100,
    "results": [
      {
        "name": "示例标签",
        "code": "example",
        "confidence": 0.95
      }
    ]
  }'
```

## 常见问题

### Q1: 上传的图片存储在哪里？
A: 默认存储在 `./uploads` 目录下，按日期组织：`uploads/2024/01/15/xxx.jpg`

### Q2: 如何处理重复上传？
A: 系统会自动计算文件MD5，相同文件不会重复存储。

### Q3: 模型服务调用失败怎么办？
A: 任务状态会更新为"失败"，错误信息会记录在数据库中。

### Q4: 可以不等待模型结果直接查询标签吗？
A: 可以。标签存储在数据库中，随时可以查询。如果标签为空，说明模型还在处理。

### Q5: 支持WebSocket实时推送吗？
A: 当前版本使用轮询或回调机制。WebSocket可以作为未来的扩展功能。

### Q6: 如何保证回调接口的安全性？
A: 建议：
1. 添加API密钥验证
2. 使用HTTPS
3. 验证请求签名
4. IP白名单限制

## 性能优化建议

1. **文件存储**：
   - 对于大规模应用，建议使用对象存储服务（如MinIO、S3）
   - 配置CDN加速图片访问

2. **数据库优化**：
   - 为常用查询字段添加索引
   - 考虑使用Redis缓存热点数据

3. **异步处理**：
   - 当前已实现异步调用模型服务
   - 可以考虑使用消息队列（如RabbitMQ、Kafka）解耦

4. **并发控制**：
   - 限制同时处理的任务数量
   - 实现任务队列管理

## 后续扩展方向

- [ ] WebSocket实时通知
- [ ] 批量上传支持
- [ ] 图片预处理（缩放、裁剪）
- [ ] 多模型选择和对比
- [ ] 识别历史统计
- [ ] 标签搜索和过滤
- [ ] 导出识别结果（CSV、JSON）
- [ ] 图片标注功能
- [ ] 模型训练反馈

## 联系支持

如有问题，请参考：
- 详细API文档：`MODEL_SERVICE_API.md`
- 数据库设计：`DATABASE_DESIGN.md`
- 项目README：`README.md`

---

**祝您使用愉快！** 🎉
