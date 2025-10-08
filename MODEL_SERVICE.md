# 模型服务集成文档

## 快速开始

### 1. 数据库迁移
```bash
cd migrations
./run_migrations.sh
```

### 2. 配置文件
编辑 `backend/etc/backend-api.yaml`：
```yaml
ModelService:
  BaseURL: http://localhost:8000  # 模型服务地址
  APIKey: ""                       # API密钥（可选）
  Timeout: 300

Upload:
  MaxSize: 10485760               # 10MB
  AllowedTypes: ".jpg,.jpeg,.png,.gif,.bmp,.webp"
  StoragePath: "./uploads"
  BaseURL: "http://localhost:8888/uploads"
```

### 3. 创建上传目录
```bash
mkdir -p uploads
chmod 755 uploads
```

### 4. 启动服务
```bash
# 启动Mock模型服务（测试用）
python3 mock_model_service.py &

# 启动后端
cd backend
go run backend.go
```

### 5. 测试
```bash
python3 test_model_api.py
```

---

## API接口

### 1. 上传图片
```bash
POST /api/images/upload
Content-Type: multipart/form-data
Authorization: Bearer {token}

参数：
  - image: 图片文件（必填）
  - model_name: 模型名称（可选）

返回：
  - image_id: 图片ID
  - task_id: 任务ID
  - file_url: 访问URL
```

**示例**：
```bash
curl -X POST http://localhost:8888/api/images/upload \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -F "image=@photo.jpg" \
  -F "model_name=resnet50"
```

### 2. 查询任务状态
```bash
GET /api/tasks/:task_id/status
Authorization: Bearer {token}

返回：
  - status: 0-待处理 1-处理中 2-已完成 3-失败
  - progress: 进度 0-100
  - result_count: 识别结果数量
```

**示例**：
```bash
curl -X GET http://localhost:8888/api/tasks/YOUR_TASK_ID/status \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### 3. 获取图片标签
```bash
GET /api/images/:image_id/labels
Authorization: Bearer {token}

返回：
  - labels: 标签数组
    - label_name: 标签名称
    - confidence: 置信度
    - bbox: 边界框（可选）
```

**示例**：
```bash
curl -X GET http://localhost:8888/api/images/123/labels \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### 4. 获取任务列表
```bash
GET /api/tasks?page=1&page_size=10
Authorization: Bearer {token}
```

### 5. 获取图片列表
```bash
GET /api/images?page=1&page_size=10
Authorization: Bearer {token}
```

### 6. 模型回调（供模型服务调用）
```bash
POST /api/model/callback
Content-Type: application/json

{
  "task_id": "任务ID",
  "status": "completed",  // pending/processing/completed/failed
  "progress": 100,
  "results": [
    {
      "name": "标签名称",
      "code": "标签代码",
      "confidence": 0.95,
      "bbox": {"x": 100, "y": 150, "width": 200, "height": 180}
    }
  ]
}
```

---

## 使用流程

### 标准流程
```
1. 上传图片 → 获得 image_id 和 task_id
2. 系统异步调用模型服务
3. 轮询查询任务状态（可选）
4. 模型完成后回调本服务
5. 查询标签结果（从数据库直接读取）
```

### 快速流程（无需轮询）
```
1. 上传图片 → 获得 image_id
2. 等待一段时间
3. 直接查询标签
   - 有数据：识别完成
   - 为空：还在处理
```

---

## Python示例

```python
import requests

BASE_URL = "http://localhost:8888"
TOKEN = "your_token"
headers = {"Authorization": f"Bearer {TOKEN}"}

# 1. 上传图片
files = {"image": open("photo.jpg", "rb")}
resp = requests.post(f"{BASE_URL}/api/images/upload", 
                     headers=headers, files=files)
result = resp.json()
image_id = result["data"]["image_id"]
task_id = result["data"]["task_id"]

# 2. 查询状态
resp = requests.get(f"{BASE_URL}/api/tasks/{task_id}/status", headers=headers)
print(resp.json())

# 3. 获取标签
resp = requests.get(f"{BASE_URL}/api/images/{image_id}/labels", headers=headers)
labels = resp.json()["data"]["labels"]
for label in labels:
    print(f"{label['label_name']}: {label['confidence']:.2%}")
```

---

## 数据库表

### images（图片表）
- id, user_id, filename, file_path, file_url
- file_size, mime_type, width, height, md5
- status, created_at, updated_at

### recognition_tasks（任务表）
- id, task_id, image_id, user_id, model_name
- status, progress, result_count
- error_message, started_at, completed_at
- created_at, updated_at

### classification_labels（标签表）
- id, task_id, image_id
- label_name, label_code, confidence
- bbox_x, bbox_y, bbox_width, bbox_height
- extra_data, created_at, updated_at

---

## 模型服务对接

模型服务需要实现以下接口：

### 1. 接收识别请求
```
POST /api/v1/recognize/upload
参数：image文件, task_id, callback URL
```

### 2. 完成后回调
```
POST {callback_url}
Body: {"task_id": "xxx", "status": "completed", "results": [...]}
```

---

## 测试工具

### Mock模型服务
```bash
python3 mock_model_service.py
# 监听 8000 端口，模拟识别过程
```

### 测试脚本
```bash
python3 test_model_api.py
# 完整测试流程
```

---

## 核心特性

- ✅ 图片上传（多格式支持、MD5去重）
- ✅ 异步模型调用（不阻塞用户）
- ✅ 进度跟踪（0-100%实时查询）
- ✅ 结果持久化（数据库存储，快速加载）
- ✅ 回调机制（模型完成自动通知）
- ✅ 静态文件服务（直接访问图片）

---

## 文件结构

```
backend/internal/
├── models/                    # 数据模型
│   ├── image.go
│   ├── recognition_task.go
│   └── classification_label.go
├── repository/                # 数据访问
│   ├── image_repository.go
│   ├── recognition_task_repository.go
│   └── classification_label_repository.go
├── handler/                   # 请求处理
│   ├── ImageHandler.go
│   ├── RecognitionTaskHandler.go
│   ├── ClassificationLabelHandler.go
│   └── ModelCallbackHandler.go
├── logic/                     # 业务逻辑
│   ├── ImageLogic.go
│   ├── RecognitionTaskLogic.go
│   ├── ClassificationLabelLogic.go
│   └── ModelCallbackLogic.go
└── svc/
    └── ModelServiceClient.go  # 模型服务客户端

migrations/
├── 005_create_images_table.sql
├── 006_create_recognition_tasks_table.sql
└── 007_create_classification_labels_table.sql
```

---

## 常见问题

**Q: 图片存储在哪里？**  
A: 默认 `./uploads` 目录，按日期组织

**Q: 如何处理重复上传？**  
A: 自动计算MD5，相同文件不重复存储

**Q: 模型调用失败怎么办？**  
A: 任务状态更新为"失败"，错误信息记录到数据库

**Q: 可以不等待模型结果吗？**  
A: 可以，标签随时可查询，为空说明还在处理

---

## 注意事项

1. **生产环境建议**：
   - 使用对象存储（MinIO/S3）
   - 配置CDN加速
   - 启用HTTPS
   - 添加回调签名验证

2. **性能优化**：
   - Redis缓存热点数据
   - 数据库连接池
   - 消息队列解耦

3. **安全性**：
   - 所有接口（除回调）需要JWT认证
   - 文件类型和大小验证
   - 用户数据隔离
