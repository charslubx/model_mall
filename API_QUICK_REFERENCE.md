# API 快速参考

## 🔑 认证

所有需要认证的接口都需要在请求头中添加：
```
Authorization: Bearer {access_token}
```

---

## 📤 上传图片

```bash
POST /api/images/upload
Content-Type: multipart/form-data

参数:
  - image: File (必填) - 图片文件
  - model_name: String (可选) - 模型名称

返回:
  - image_id: 图片ID
  - task_id: 任务ID
  - file_url: 图片URL
  - filename: 文件名
```

**示例**:
```bash
curl -X POST http://localhost:8888/api/images/upload \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -F "image=@photo.jpg" \
  -F "model_name=resnet50"
```

---

## 🔍 查询任务状态

```bash
GET /api/tasks/:task_id/status

返回:
  - task_id: 任务ID
  - image_id: 图片ID
  - status: 状态 (0-待处理 1-处理中 2-已完成 3-失败)
  - progress: 进度 (0-100)
  - result_count: 识别结果数量
  - error_msg: 错误信息
  - created_at: 创建时间
  - completed_at: 完成时间
```

**示例**:
```bash
curl -X GET http://localhost:8888/api/tasks/YOUR_TASK_ID/status \
  -H "Authorization: Bearer YOUR_TOKEN"
```

---

## 📋 获取任务列表

```bash
GET /api/tasks?page=1&page_size=10

参数:
  - page: 页码 (默认: 1)
  - page_size: 每页数量 (默认: 10)

返回:
  - list: 任务列表
  - total: 总数
  - page: 当前页
  - page_size: 每页数量
```

**示例**:
```bash
curl -X GET "http://localhost:8888/api/tasks?page=1&page_size=20" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

---

## 🏷️ 获取图片标签

```bash
GET /api/images/:image_id/labels

返回:
  - image_id: 图片ID
  - labels: 标签列表
    - id: 标签ID
    - label_name: 标签名称
    - label_code: 标签代码
    - confidence: 置信度 (0.0-1.0)
    - bbox: 边界框 (可选)
      - x, y: 坐标
      - width, height: 宽高
    - extra_data: 额外数据 (可选)
    - created_at: 创建时间
```

**示例**:
```bash
curl -X GET http://localhost:8888/api/images/123/labels \
  -H "Authorization: Bearer YOUR_TOKEN"
```

---

## 🖼️ 获取图片列表

```bash
GET /api/images?page=1&page_size=10

参数:
  - page: 页码 (默认: 1)
  - page_size: 每页数量 (默认: 10)

返回:
  - list: 图片列表
  - total: 总数
  - page: 当前页
  - page_size: 每页数量
```

**示例**:
```bash
curl -X GET "http://localhost:8888/api/images?page=1&page_size=20" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

---

## 🔄 模型回调 (供模型服务调用)

```bash
POST /api/model/callback
Content-Type: application/json

{
  "task_id": "任务ID",
  "status": "状态 (pending/processing/completed/failed)",
  "progress": 进度 (0-100),
  "results": [  // status为completed时提供
    {
      "name": "标签名称",
      "code": "标签代码",
      "confidence": 置信度,
      "bbox": {  // 可选
        "x": x坐标,
        "y": y坐标,
        "width": 宽度,
        "height": 高度
      },
      "extra": {}  // 可选的额外数据
    }
  ],
  "error": "错误信息"  // status为failed时提供
}
```

**示例**:
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
        "confidence": 0.9856
      }
    ]
  }'
```

---

## 📁 访问图片

```bash
GET /uploads/{日期}/{文件名}

示例: http://localhost:8888/uploads/2024/01/15/abc123.jpg
```

---

## 📊 响应格式

### 成功响应
```json
{
  "code": 0,
  "message": "success",
  "data": { ... }
}
```

### 错误响应
```json
{
  "code": 1001,
  "message": "错误描述",
  "data": null
}
```

---

## 🎯 常用工作流

### 流程1: 上传并等待结果
```
1. POST /api/images/upload
   → 获得 image_id 和 task_id

2. 轮询 GET /api/tasks/:task_id/status
   → 查看进度和状态

3. 状态变为"已完成"后
   GET /api/images/:image_id/labels
   → 获取识别结果
```

### 流程2: 快速查询标签
```
1. POST /api/images/upload
   → 获得 image_id

2. 等待几秒（或通过其他方式得知已完成）

3. GET /api/images/:image_id/labels
   → 直接从数据库获取标签
   → 如果为空，说明还在处理
```

### 流程3: 查看历史记录
```
1. GET /api/images
   → 获取所有上传的图片

2. GET /api/tasks
   → 获取所有识别任务

3. GET /api/images/:image_id/labels
   → 查看特定图片的标签
```

---

## 🔧 配置示例

### backend/etc/backend-api.yaml
```yaml
ModelService:
  BaseURL: http://localhost:8000
  APIKey: ""
  Timeout: 300

Upload:
  MaxSize: 10485760  # 10MB
  AllowedTypes: ".jpg,.jpeg,.png,.gif,.bmp,.webp"
  StoragePath: "./uploads"
  BaseURL: "http://localhost:8888/uploads"
```

---

## 🐍 Python示例

```python
import requests

BASE_URL = "http://localhost:8888"
TOKEN = "your_token_here"
headers = {"Authorization": f"Bearer {TOKEN}"}

# 上传图片
files = {"image": open("photo.jpg", "rb")}
data = {"model_name": "resnet50"}
resp = requests.post(f"{BASE_URL}/api/images/upload", 
                     headers=headers, files=files, data=data)
result = resp.json()
image_id = result["data"]["image_id"]
task_id = result["data"]["task_id"]

# 查询状态
resp = requests.get(f"{BASE_URL}/api/tasks/{task_id}/status",
                    headers=headers)
print(resp.json())

# 获取标签
resp = requests.get(f"{BASE_URL}/api/images/{image_id}/labels",
                    headers=headers)
labels = resp.json()["data"]["labels"]
for label in labels:
    print(f"{label['label_name']}: {label['confidence']:.2%}")
```

---

## 🔍 JavaScript示例

```javascript
const BASE_URL = 'http://localhost:8888';
const TOKEN = 'your_token_here';

// 上传图片
const formData = new FormData();
formData.append('image', fileInput.files[0]);
formData.append('model_name', 'resnet50');

const uploadResp = await fetch(`${BASE_URL}/api/images/upload`, {
  method: 'POST',
  headers: { 'Authorization': `Bearer ${TOKEN}` },
  body: formData
});

const uploadResult = await uploadResp.json();
const { image_id, task_id } = uploadResult.data;

// 查询状态
const statusResp = await fetch(
  `${BASE_URL}/api/tasks/${task_id}/status`,
  { headers: { 'Authorization': `Bearer ${TOKEN}` }}
);

// 获取标签
const labelsResp = await fetch(
  `${BASE_URL}/api/images/${image_id}/labels`,
  { headers: { 'Authorization': `Bearer ${TOKEN}` }}
);

const labels = (await labelsResp.json()).data.labels;
labels.forEach(label => {
  console.log(`${label.label_name}: ${(label.confidence * 100).toFixed(2)}%`);
});
```

---

## 📞 获取帮助

- 详细文档: `MODEL_SERVICE_API.md`
- 快速开始: `MODEL_SERVICE_SETUP.md`
- 变更日志: `CHANGELOG_MODEL_SERVICE.md`
- 测试脚本: `test_model_api.py` 或 `test_model_api.sh`
- Mock服务: `mock_model_service.py`

---

**提示**: 使用 Mock 模型服务进行本地测试：
```bash
python3 mock_model_service.py
```
