# 模型服务交互接口文档

## 概述

本文档描述了与模型服务交互的所有接口，包括图片上传、模型识别、进度查询和标签获取等功能。

## 数据库设计

### 新增表

#### 1. images（图片表）
存储上传的图片信息。

| 字段 | 类型 | 说明 |
|------|------|------|
| id | BIGSERIAL | 主键 |
| user_id | BIGINT | 用户ID |
| filename | VARCHAR(255) | 文件名 |
| original_name | VARCHAR(255) | 原始文件名 |
| file_path | VARCHAR(500) | 文件路径 |
| file_url | VARCHAR(500) | 文件URL |
| file_size | BIGINT | 文件大小（字节） |
| mime_type | VARCHAR(100) | MIME类型 |
| width | INTEGER | 图片宽度 |
| height | INTEGER | 图片高度 |
| md5 | VARCHAR(32) | 文件MD5 |
| status | SMALLINT | 状态 0-已删除 1-正常 |
| created_at | TIMESTAMP | 创建时间 |
| updated_at | TIMESTAMP | 更新时间 |

#### 2. recognition_tasks（识别任务表）
存储模型识别任务及其进度。

| 字段 | 类型 | 说明 |
|------|------|------|
| id | BIGSERIAL | 主键 |
| task_id | VARCHAR(100) | 任务唯一标识 |
| image_id | BIGINT | 图片ID |
| user_id | BIGINT | 用户ID |
| model_name | VARCHAR(100) | 使用的模型名称 |
| status | SMALLINT | 状态 0-待处理 1-处理中 2-已完成 3-失败 |
| progress | INTEGER | 进度 0-100 |
| result_count | INTEGER | 识别结果数量 |
| error_message | TEXT | 错误信息 |
| started_at | TIMESTAMP | 开始时间 |
| completed_at | TIMESTAMP | 完成时间 |
| created_at | TIMESTAMP | 创建时间 |
| updated_at | TIMESTAMP | 更新时间 |

#### 3. classification_labels（分类标签表）
存储模型识别的分类标签结果。

| 字段 | 类型 | 说明 |
|------|------|------|
| id | BIGSERIAL | 主键 |
| task_id | BIGINT | 任务ID |
| image_id | BIGINT | 图片ID |
| label_name | VARCHAR(100) | 标签名称 |
| label_code | VARCHAR(100) | 标签代码 |
| confidence | DECIMAL(5,4) | 置信度 0.0000-1.0000 |
| bbox_x | INTEGER | 边界框X坐标 |
| bbox_y | INTEGER | 边界框Y坐标 |
| bbox_width | INTEGER | 边界框宽度 |
| bbox_height | INTEGER | 边界框高度 |
| extra_data | JSONB | 额外数据（JSON格式） |
| created_at | TIMESTAMP | 创建时间 |
| updated_at | TIMESTAMP | 更新时间 |

## API接口

### 1. 上传图片并创建识别任务

上传图片文件，系统会自动创建识别任务并异步调用模型服务。

**接口地址**：`POST /api/images/upload`

**请求头**：
```
Authorization: Bearer {access_token}
Content-Type: multipart/form-data
```

**请求参数**：
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| image | File | 是 | 图片文件 |
| model_name | String | 否 | 指定使用的模型名称 |

**响应示例**：
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "image_id": 123,
    "task_id": "550e8400-e29b-41d4-a716-446655440000",
    "file_url": "http://localhost:8888/uploads/2024/01/15/abc123.jpg",
    "filename": "abc123.jpg"
  }
}
```

**说明**：
- 上传成功后会立即返回，任务会在后台异步处理
- 如果相同MD5的文件已存在，会直接返回已有记录
- 支持的文件类型：.jpg, .jpeg, .png, .gif, .bmp, .webp
- 最大文件大小：10MB（可在配置中修改）

---

### 2. 查询任务状态

查询识别任务的当前状态和进度。

**接口地址**：`GET /api/tasks/:task_id/status`

**请求头**：
```
Authorization: Bearer {access_token}
```

**路径参数**：
| 参数 | 类型 | 说明 |
|------|------|------|
| task_id | String | 任务ID |

**响应示例**：
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "task_id": "550e8400-e29b-41d4-a716-446655440000",
    "image_id": 123,
    "status": 2,
    "progress": 100,
    "result_count": 5,
    "error_msg": "",
    "created_at": "2024-01-15 10:30:00",
    "completed_at": "2024-01-15 10:30:15"
  }
}
```

**状态说明**：
- 0: 待处理
- 1: 处理中
- 2: 已完成
- 3: 失败

---

### 3. 获取任务列表

获取当前用户的所有识别任务列表。

**接口地址**：`GET /api/tasks`

**请求头**：
```
Authorization: Bearer {access_token}
```

**查询参数**：
| 参数 | 类型 | 必填 | 默认值 | 说明 |
|------|------|------|--------|------|
| page | Integer | 否 | 1 | 页码 |
| page_size | Integer | 否 | 10 | 每页数量 |

**响应示例**：
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "list": [
      {
        "task_id": "550e8400-e29b-41d4-a716-446655440000",
        "image_id": 123,
        "model_name": "resnet50",
        "status": 2,
        "progress": 100,
        "result_count": 5,
        "created_at": "2024-01-15 10:30:00",
        "completed_at": "2024-01-15 10:30:15"
      }
    ],
    "total": 100,
    "page": 1,
    "page_size": 10
  }
}
```

---

### 4. 获取图片标签

获取图片的所有分类标签（从数据库加载，无需等待模型结果）。

**接口地址**：`GET /api/images/:image_id/labels`

**请求头**：
```
Authorization: Bearer {access_token}
```

**路径参数**：
| 参数 | 类型 | 说明 |
|------|------|------|
| image_id | Integer | 图片ID |

**响应示例**：
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "image_id": 123,
    "labels": [
      {
        "id": 1,
        "label_name": "猫",
        "label_code": "cat",
        "confidence": 0.9856,
        "bbox": {
          "x": 100,
          "y": 150,
          "width": 200,
          "height": 180
        },
        "extra_data": {
          "breed": "波斯猫"
        },
        "created_at": "2024-01-15 10:30:15"
      },
      {
        "id": 2,
        "label_name": "狗",
        "label_code": "dog",
        "confidence": 0.8234,
        "created_at": "2024-01-15 10:30:15"
      }
    ]
  }
}
```

**说明**：
- 标签按置信度从高到低排序
- bbox字段（边界框）为可选，仅在目标检测任务中存在
- extra_data为可选的额外数据，可包含任意JSON格式信息

---

### 5. 获取图片列表

获取当前用户上传的所有图片列表。

**接口地址**：`GET /api/images`

**请求头**：
```
Authorization: Bearer {access_token}
```

**查询参数**：
| 参数 | 类型 | 必填 | 默认值 | 说明 |
|------|------|------|--------|------|
| page | Integer | 否 | 1 | 页码 |
| page_size | Integer | 否 | 10 | 每页数量 |

**响应示例**：
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "list": [
      {
        "id": 123,
        "filename": "abc123.jpg",
        "original_name": "my_photo.jpg",
        "file_url": "http://localhost:8888/uploads/2024/01/15/abc123.jpg",
        "file_size": 1024000,
        "width": 1920,
        "height": 1080,
        "status": 1,
        "created_at": "2024-01-15 10:30:00"
      }
    ],
    "total": 50,
    "page": 1,
    "page_size": 10
  }
}
```

---

### 6. 模型服务回调接口（供模型服务调用）

此接口供模型服务在处理完成后回调，更新任务状态和保存识别结果。

**接口地址**：`POST /api/model/callback`

**请求头**：
```
Content-Type: application/json
```

**请求体**：
```json
{
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
      },
      "extra": {
        "breed": "波斯猫"
      }
    }
  ],
  "error": ""
}
```

**参数说明**：
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| task_id | String | 是 | 任务ID |
| status | String | 是 | 状态：pending/processing/completed/failed |
| progress | Integer | 否 | 进度 0-100 |
| results | Array | 否 | 识别结果数组（status为completed时提供） |
| error | String | 否 | 错误信息（status为failed时提供） |

**响应示例**：
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "success": true,
    "message": "success"
  }
}
```

---

## 工作流程

### 标准流程

1. **前端上传图片**
   ```
   POST /api/images/upload
   → 返回 image_id 和 task_id
   ```

2. **后端异步调用模型服务**
   ```
   系统自动调用模型服务
   → 传递图片和回调URL
   ```

3. **前端轮询任务状态**（可选）
   ```
   GET /api/tasks/:task_id/status
   → 获取实时进度
   ```

4. **模型服务处理完成后回调**
   ```
   模型服务调用: POST /api/model/callback
   → 保存识别结果到数据库
   ```

5. **前端获取标签结果**
   ```
   GET /api/images/:image_id/labels
   → 从数据库直接加载，无需等待
   ```

### 快速流程（无需轮询）

1. 上传图片获得 `image_id`
2. 等待一段时间（或通过WebSocket接收通知）
3. 直接查询标签：`GET /api/images/:image_id/labels`
4. 如果标签列表为空，说明模型还在处理中

---

## 配置说明

### 配置文件位置
`backend/etc/backend-api.yaml`

### 模型服务配置
```yaml
ModelService:
  BaseURL: http://localhost:8000  # 模型服务地址
  APIKey: ""  # 模型服务API密钥（如需要）
  Timeout: 300  # 超时时间（秒）
```

### 上传配置
```yaml
Upload:
  MaxSize: 10485760  # 最大文件大小（字节），默认10MB
  AllowedTypes: ".jpg,.jpeg,.png,.gif,.bmp,.webp"  # 允许的文件类型
  StoragePath: "./uploads"  # 本地存储路径
  BaseURL: "http://localhost:8888/uploads"  # 访问基础URL
```

---

## 数据库迁移

### 执行迁移
```bash
cd migrations
chmod +x run_migrations.sh
./run_migrations.sh
```

或使用Go直接执行：
```bash
cd migrations
go run migrate.go
```

### 迁移文件
- `005_create_images_table.sql` - 创建图片表
- `006_create_recognition_tasks_table.sql` - 创建识别任务表
- `007_create_classification_labels_table.sql` - 创建分类标签表

---

## 模型服务接口规范

模型服务需要实现以下接口供本系统调用：

### 1. 上传图片识别
**接口**：`POST /api/v1/recognize/upload`

**请求**：multipart/form-data
- `image`: 图片文件
- `task_id`: 任务ID
- `model_name`: 模型名称（可选）
- `callback`: 回调URL

**响应**：
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "task_id": "xxx",
    "status": "processing"
  }
}
```

### 2. URL识别
**接口**：`POST /api/v1/recognize/url`

**请求**：
```json
{
  "task_id": "xxx",
  "image_url": "http://...",
  "model_name": "resnet50",
  "callback": "http://localhost:8888/api/model/callback"
}
```

### 3. 查询任务状态
**接口**：`GET /api/v1/task/:task_id/status`

**响应**：
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "task_id": "xxx",
    "status": "completed",
    "progress": 100,
    "result": [...]
  }
}
```

---

## 错误处理

所有接口遵循统一的错误响应格式：

```json
{
  "code": 1001,
  "message": "错误描述",
  "data": null
}
```

**常见错误码**：
- `1001`: 参数错误
- `1002`: 未授权
- `1003`: 资源不存在
- `1004`: 操作失败
- `5000`: 服务器内部错误

---

## 性能优化建议

1. **图片去重**：通过MD5避免重复上传相同文件
2. **异步处理**：上传和识别异步进行，不阻塞用户
3. **进度查询**：支持轮询或WebSocket实时推送
4. **结果缓存**：识别结果存储在数据库，快速加载
5. **分页查询**：大量数据使用分页避免性能问题

---

## 安全性

1. **认证**：除回调接口外，所有接口需要JWT认证
2. **文件校验**：检查文件类型、大小、内容
3. **用户隔离**：用户只能访问自己的数据
4. **回调验证**：建议为回调接口添加签名验证（可选）

---

## 示例代码

### JavaScript/Fetch
```javascript
// 上传图片
const formData = new FormData();
formData.append('image', fileInput.files[0]);
formData.append('model_name', 'resnet50');

const response = await fetch('http://localhost:8888/api/images/upload', {
  method: 'POST',
  headers: {
    'Authorization': `Bearer ${accessToken}`
  },
  body: formData
});

const result = await response.json();
console.log('Task ID:', result.data.task_id);

// 查询任务状态
const taskResponse = await fetch(
  `http://localhost:8888/api/tasks/${result.data.task_id}/status`,
  {
    headers: {
      'Authorization': `Bearer ${accessToken}`
    }
  }
);

// 获取标签
const labelsResponse = await fetch(
  `http://localhost:8888/api/images/${result.data.image_id}/labels`,
  {
    headers: {
      'Authorization': `Bearer ${accessToken}`
    }
  }
);
```

### Python/Requests
```python
import requests

# 上传图片
files = {'image': open('photo.jpg', 'rb')}
data = {'model_name': 'resnet50'}
headers = {'Authorization': f'Bearer {access_token}'}

response = requests.post(
    'http://localhost:8888/api/images/upload',
    files=files,
    data=data,
    headers=headers
)

result = response.json()
task_id = result['data']['task_id']
image_id = result['data']['image_id']

# 查询任务状态
task_response = requests.get(
    f'http://localhost:8888/api/tasks/{task_id}/status',
    headers=headers
)

# 获取标签
labels_response = requests.get(
    f'http://localhost:8888/api/images/{image_id}/labels',
    headers=headers
)

labels = labels_response.json()['data']['labels']
for label in labels:
    print(f"{label['label_name']}: {label['confidence']:.2%}")
```

---

## 总结

本系统提供了完整的图片识别工作流：
- ✅ 图片上传与存储
- ✅ 异步模型调用
- ✅ 任务进度跟踪
- ✅ 结果持久化存储
- ✅ 快速标签查询
- ✅ 灵活的回调机制

系统支持两种使用模式：
1. **轮询模式**：上传后定期查询任务状态
2. **回调模式**：模型完成后主动通知，前端按需查询标签
