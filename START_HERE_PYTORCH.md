# 🚀 PyTorch模型快速部署（中文版）

你有一个PyTorch训练的`.mph`模型文件？这里是最简单的部署步骤！

## ⚡ 3步完成部署

### 步骤1️⃣：测试你的模型

```bash
cd /workspace/model_service

# 安装PyTorch（如果需要）
pip install torch torchvision

# 测试模型文件
python test_pytorch_model.py /path/to/your/model.mph
```

**如果看到 "✅ 测试结果: 模型可以直接使用！"** → 继续步骤2

**如果看到 "⚠️ 测试结果: 需要额外配置"** → 查看 [PYTORCH_MODEL_GUIDE.md](model_service/PYTORCH_MODEL_GUIDE.md)

### 步骤2️⃣：准备文件

```bash
# 复制模型到正确位置
cp /path/to/your/model.mph /workspace/model_service/models/

# 创建标签文件（替换为你的类别）
cat > /workspace/model_service/models/labels.txt << EOF
类别1
类别2
类别3
EOF
```

### 步骤3️⃣：一键启动

```bash
cd /workspace
./start_all_services.sh
```

等待30秒，所有服务自动启动！

## ✅ 测试

```bash
# 1. 测试模型服务
curl http://localhost:5000/health

# 2. 上传图片测试
curl -X POST http://localhost:8888/api/images/upload \
  -F "image=@test.jpg"
```

## 📊 查看结果

成功响应示例：
```json
{
  "image_id": 1,
  "filename": "test.jpg",
  "classifications": [
    {
      "label": "类别1",
      "confidence": 0.8523
    }
  ]
}
```

## 🎯 常见问题

### Q: 启动失败？

```bash
# 查看日志
docker-compose logs model-service

# 常见原因：
# 1. 模型文件不存在 → 检查 model_service/models/
# 2. 标签文件格式错误 → 每行一个标签，不要有空行
# 3. 端口被占用 → 修改 docker-compose.yml 中的端口
```

### Q: 分类结果不对？

检查 `model_service/models/labels.txt`：
- ✅ 标签顺序要与模型训练时一致
- ✅ 每行一个标签
- ✅ 类别数量要匹配

### Q: 需要GPU加速？

编辑 `model_service/docker-compose.yml`：

```yaml
model-service:
  image: pytorch/pytorch:2.1.0-cuda11.8-cudnn8-runtime
  deploy:
    resources:
      reservations:
        devices:
          - driver: nvidia
            count: 1
            capabilities: [gpu]
  environment:
    - DEVICE=cuda
```

## 📚 更多文档

- **详细教程**: [PYTORCH_QUICKSTART.md](PYTORCH_QUICKSTART.md)
- **完整指南**: [PYTORCH_MODEL_GUIDE.md](model_service/PYTORCH_MODEL_GUIDE.md)
- **故障排查**: [DEPLOYMENT.md](DEPLOYMENT.md)

## 💡 下一步

✅ 部署成功后，你可以：

1. **集成前端** - 创建图片上传界面
2. **查看图片列表** - `GET /api/images`
3. **查看单个图片** - `GET /api/images/{id}`
4. **添加更多模型** - 部署多个模型服务

---

需要帮助？查看 [完整文档](README.md) 或提交 Issue。

**祝你部署顺利！** 🎉
