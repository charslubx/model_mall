# 模型服务集成 - 变更日志

## 📅 更新日期
2024年（当前）

## 🎯 更新概述

添加了完整的模型服务交互功能，支持图片上传、模型识别、进度跟踪和标签管理。

## 📦 新增文件

### 数据库迁移文件
- `migrations/005_create_images_table.sql` - 图片表
- `migrations/006_create_recognition_tasks_table.sql` - 识别任务表
- `migrations/007_create_classification_labels_table.sql` - 分类标签表

### 数据模型 (Models)
- `backend/internal/models/image.go` - 图片模型
- `backend/internal/models/recognition_task.go` - 识别任务模型
- `backend/internal/models/classification_label.go` - 分类标签模型

### 数据访问层 (Repository)
- `backend/internal/repository/image_repository.go` - 图片数据操作
- `backend/internal/repository/recognition_task_repository.go` - 任务数据操作
- `backend/internal/repository/classification_label_repository.go` - 标签数据操作

### 服务层 (Service)
- `backend/internal/svc/ModelServiceClient.go` - 模型服务HTTP客户端

### 处理器 (Handler)
- `backend/internal/handler/ImageHandler.go` - 图片上传处理
- `backend/internal/handler/RecognitionTaskHandler.go` - 任务查询处理
- `backend/internal/handler/ClassificationLabelHandler.go` - 标签查询处理
- `backend/internal/handler/ModelCallbackHandler.go` - 模型回调处理
- `backend/internal/handler/StaticFileHandler.go` - 静态文件服务

### 业务逻辑 (Logic)
- `backend/internal/logic/ImageLogic.go` - 图片业务逻辑
- `backend/internal/logic/RecognitionTaskLogic.go` - 任务业务逻辑
- `backend/internal/logic/ClassificationLabelLogic.go` - 标签业务逻辑
- `backend/internal/logic/ModelCallbackLogic.go` - 回调业务逻辑

### 文档和测试
- `MODEL_SERVICE_API.md` - 详细的API文档
- `MODEL_SERVICE_SETUP.md` - 快速开始指南
- `CHANGELOG_MODEL_SERVICE.md` - 本文件
- `test_model_api.sh` - Shell测试脚本
- `test_model_api.py` - Python测试脚本
- `mock_model_service.py` - Mock模型服务（用于测试）

## 🔧 修改文件

### 配置文件
**文件**: `backend/internal/config/config.go`
- 新增 `ModelService` 配置项（BaseURL, APIKey, Timeout）
- 新增 `Upload` 配置项（MaxSize, AllowedTypes, StoragePath, BaseURL）

**文件**: `backend/etc/backend-api.yaml`
- 添加模型服务配置
- 添加文件上传配置

### 核心文件
**文件**: `backend/internal/types/types.go`
- 新增图片上传相关类型
- 新增识别任务相关类型
- 新增分类标签相关类型
- 新增模型回调相关类型

**文件**: `backend/internal/handler/routes.go`
- 添加图片上传路由
- 添加任务查询路由
- 添加标签查询路由
- 添加模型回调路由

**文件**: `backend/internal/repository/repository.go`
- 添加 `Image` 仓库
- 添加 `RecognitionTask` 仓库
- 添加 `ClassificationLabel` 仓库

**文件**: `backend/internal/svc/servicecontext.go`
- 添加 `ModelServiceClient` 字段
- 初始化模型服务客户端

**文件**: `backend/backend.go`
- 添加静态文件服务路由

**文件**: `migrations/migrate.go`
- 更新GORM自动迁移，包含新表

## 🌟 核心功能

### 1. 图片上传
- 支持多种图片格式（.jpg, .jpeg, .png, .gif, .bmp, .webp）
- 文件大小限制（默认10MB）
- MD5去重，避免重复存储
- 按日期组织存储目录
- 自动生成唯一文件名

### 2. 模型识别
- 异步调用模型服务
- 支持两种调用方式：上传文件 / 提供URL
- 自动创建识别任务
- 传递回调URL给模型服务

### 3. 进度跟踪
- 实时查询任务状态
- 支持进度更新（0-100%）
- 任务状态：待处理、处理中、已完成、失败
- 记录开始和完成时间

### 4. 标签管理
- 存储识别结果到数据库
- 支持多个标签
- 包含置信度信息
- 支持边界框（目标检测）
- 支持额外数据（JSON格式）
- 按置信度排序

### 5. 回调机制
- 模型服务完成后主动回调
- 更新任务状态
- 保存识别结果
- 错误处理

### 6. 数据持久化
- 所有数据存储在PostgreSQL
- 支持快速查询
- 无需等待模型处理
- 支持历史记录查询

## 📊 数据表结构

### images（图片表）
- 存储上传的图片信息
- 关联用户ID
- 包含文件元数据（大小、类型、尺寸等）
- MD5用于去重

### recognition_tasks（识别任务表）
- 任务唯一标识（UUID）
- 关联图片和用户
- 状态和进度跟踪
- 结果数量统计
- 时间记录

### classification_labels（分类标签表）
- 关联任务和图片
- 标签名称和代码
- 置信度
- 边界框信息（可选）
- 额外数据（JSONB）

## 🔌 API接口

### 需要认证的接口
1. `POST /api/images/upload` - 上传图片
2. `GET /api/images` - 获取图片列表
3. `GET /api/images/:image_id/labels` - 获取图片标签
4. `GET /api/tasks/:task_id/status` - 查询任务状态
5. `GET /api/tasks` - 获取任务列表

### 公开接口
1. `POST /api/model/callback` - 模型服务回调
2. `GET /uploads/*` - 访问上传的图片

## 🚀 使用流程

### 标准流程
```
1. 用户上传图片 → 返回 image_id 和 task_id
2. 系统异步调用模型服务
3. 用户轮询任务状态（可选）
4. 模型服务完成后回调
5. 系统保存识别结果到数据库
6. 用户查询标签（从数据库直接读取）
```

### 快速流程
```
1. 用户上传图片 → 获得 image_id
2. 等待一段时间
3. 直接查询标签 → 从数据库加载
   - 如果有数据：识别已完成
   - 如果为空：还在处理中
```

## 🔒 安全性

1. **认证**：除回调接口外，所有接口需要JWT认证
2. **文件验证**：检查文件类型、大小、格式
3. **用户隔离**：用户只能访问自己的数据
4. **路径安全**：防止目录遍历攻击
5. **错误处理**：统一的错误响应格式

## ⚡ 性能优化

1. **异步处理**：上传后立即返回，不阻塞用户
2. **MD5去重**：避免重复存储相同文件
3. **数据库索引**：为常用查询字段添加索引
4. **分页查询**：支持大量数据的分页加载
5. **静态文件服务**：直接提供图片访问

## 📝 配置说明

### 模型服务配置
```yaml
ModelService:
  BaseURL: http://localhost:8000  # 模型服务地址
  APIKey: ""                       # API密钥（可选）
  Timeout: 300                     # 超时时间（秒）
```

### 上传配置
```yaml
Upload:
  MaxSize: 10485760               # 最大文件大小（字节）
  AllowedTypes: ".jpg,.jpeg,.png,.gif,.bmp,.webp"
  StoragePath: "./uploads"        # 本地存储路径
  BaseURL: "http://localhost:8888/uploads"
```

## 🧪 测试方法

### 方法1：使用测试脚本
```bash
# Shell脚本
chmod +x test_model_api.sh
./test_model_api.sh

# Python脚本
python3 test_model_api.py
```

### 方法2：使用Mock模型服务
```bash
# 1. 启动Mock模型服务
python3 mock_model_service.py

# 2. 启动后端服务
cd backend
go run backend.go

# 3. 上传图片测试
curl -X POST http://localhost:8888/api/images/upload \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -F "image=@test.jpg"
```

### 方法3：手动测试
参考 `MODEL_SERVICE_API.md` 中的示例代码

## 🐛 已知问题

无

## 📋 待办事项

- [ ] WebSocket实时通知
- [ ] 批量上传支持
- [ ] 图片预处理（缩放、裁剪）
- [ ] 多模型选择和对比
- [ ] 识别历史统计
- [ ] 标签搜索和过滤
- [ ] 导出识别结果
- [ ] 图片标注功能
- [ ] 对象存储集成（MinIO/S3）

## 💡 使用建议

1. **开发环境**：使用Mock模型服务进行测试
2. **生产环境**：
   - 配置真实的模型服务地址
   - 使用对象存储服务
   - 配置CDN加速
   - 启用HTTPS
   - 添加监控和日志

## 📞 技术支持

如有问题，请参考：
- 详细API文档：`MODEL_SERVICE_API.md`
- 快速开始指南：`MODEL_SERVICE_SETUP.md`
- 数据库设计：`DATABASE_DESIGN.md`
- 项目README：`README.md`

## 🤝 贡献

欢迎提交Issue和Pull Request！

---

**版本**: v1.0.0  
**更新时间**: 2024年  
**状态**: ✅ 已完成
