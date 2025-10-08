# 🚀 最终使用指南

## 📋 项目说明

这是一个图片分类和存储系统，包括：
- Go后端服务（处理业务逻辑）
- Python模型服务（PyTorch图片分类）
- PostgreSQL数据库（存储数据）
- Redis缓存

## ⚡ 2步部署（超简单）

### 步骤1: 保存完整模型

在你的训练代码中：

```python
# 训练完成后，保存完整模型
model.eval()
torch.save(model, 'model.pth')
```

### 步骤2: 启动系统

```bash
# 1. 复制模型文件
cp model.pth /workspace/model_service/models/

# 2. 创建标签文件
cat > /workspace/model_service/models/labels.txt << EOF
类别1
类别2
类别3
EOF

# 3. 一键启动
cd /workspace
./start_all_services.sh
```

**完成！** 🎉

## 🧪 测试

```bash
# 上传图片
curl -X POST http://localhost:8888/api/images/upload \
  -F "image=@test.jpg"

# 响应示例
{
  "image_id": 1,
  "filename": "test.jpg",
  "classifications": [
    {"label": "类别1", "confidence": 0.85}
  ]
}
```

## 📁 关键文件位置

```
model_service/
├── models/
│   ├── model.pth          ← 放你的PyTorch完整模型
│   └── labels.txt         ← 分类标签，每行一个
└── docker-compose.yml     ← 配置文件（通常不需要改）
```

## 🔧 配置说明

### 如果模型文件名不是 model.pth

编辑 `model_service/docker-compose.yml`:

```yaml
environment:
  - MODEL_PATH=/app/models/your_model_name.pth  # 改成你的文件名
```

### 如果需要GPU加速

```yaml
environment:
  - DEVICE=cuda  # 使用GPU
```

需要安装NVIDIA Docker Runtime。

## 📊 API接口

### 1. 上传图片
```bash
POST http://localhost:8888/api/images/upload
```

### 2. 获取图片信息
```bash
GET http://localhost:8888/api/images/{id}
```

### 3. 获取图片列表
```bash
GET http://localhost:8888/api/images?page=1&page_size=10
```

## 🐛 问题排查

### 服务启动失败

```bash
# 查看日志
docker-compose logs -f model-service

# 常见原因：
# 1. 模型文件不存在 → 检查 model_service/models/
# 2. 模型格式错误 → 确保是 torch.save(model, ...) 保存的
```

### 模型加载失败

```bash
# 检查模型文件
ls -lh model_service/models/

# 测试模型
python model_service/test_pytorch_model.py model_service/models/model.pth
```

## 📚 详细文档

- **[SIMPLE_DEPLOYMENT.md](SIMPLE_DEPLOYMENT.md)** - 完整部署步骤
- **[SAVE_FULL_MODEL_GUIDE.md](SAVE_FULL_MODEL_GUIDE.md)** - 如何保存完整模型
- **[IMAGE_CLASSIFICATION_GUIDE.md](IMAGE_CLASSIFICATION_GUIDE.md)** - API使用文档

## 💡 重要提示

1. **必须使用完整模型** - `torch.save(model, 'model.pth')`
2. **标签顺序很重要** - 必须与模型训练时的类别索引对应
3. **模型要设为eval模式** - 保存前执行 `model.eval()`

## 🎯 服务地址

| 服务 | 地址 | 说明 |
|------|------|------|
| Go后端 | http://localhost:8888 | 业务API |
| 模型服务 | http://localhost:5000 | 图片分类 |
| 数据库 | localhost:5432 | PostgreSQL |
| Redis | localhost:6379 | 缓存 |

## ✅ 验证部署成功

```bash
# 1. 模型服务健康检查
curl http://localhost:5000/health
# 期望: {"status":"healthy","model_loaded":true}

# 2. 测试图片分类
curl -X POST http://localhost:5000/classify -F "image=@test.jpg"
# 期望: {"success":true,"results":[...]}

# 3. 测试完整流程
curl -X POST http://localhost:8888/api/images/upload -F "image=@test.jpg"
# 期望: {"image_id":1,"classifications":[...]}
```

---

**就这么简单！享受你的图片分类服务吧！** 🎊
