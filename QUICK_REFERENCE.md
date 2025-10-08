# 快速参考手册

## 📋 你的模型情况

- ✅ PyTorch模型
- ✅ Checkpoint格式（只有权重）
- ✅ 自定义/改过的模型架构

## 🚀 3步部署

### 1. 复制模型定义

**编辑:** `/workspace/model_service/app/model_architecture.py`

```python
# 在文件末尾添加你的模型类
class YourModel(nn.Module):
    def __init__(self, num_classes=10):
        # ... 从训练代码复制
        pass
    
    def forward(self, x):
        # ... 从训练代码复制
        pass

# 在 create_model 函数中添加
def create_model(model_name='resnet18', num_classes=10):
    # ...
    elif model_name == 'my_model':
        return YourModel(num_classes)
```

### 2. 配置环境

**编辑:** `/workspace/model_service/docker-compose.yml`

```yaml
environment:
  - MODEL_PATH=/app/models/checkpoint_best.pth
  - MODEL_ARCH=my_model     # ← 步骤1中定义的名字
  - NUM_CLASSES=10          # ← 你的类别数
```

### 3. 启动

```bash
cp checkpoint_best.pth /workspace/model_service/models/
cat > /workspace/model_service/models/labels.txt << EOF
类别1
类别2
...
EOF
cd /workspace && ./start_all_services.sh
```

## 📁 关键文件位置

```
/workspace/
├── model_service/
│   ├── app/
│   │   └── model_architecture.py    ← 添加模型定义
│   ├── docker-compose.yml           ← 配置环境变量
│   └── models/
│       ├── checkpoint_best.pth      ← 你的checkpoint
│       └── labels.txt               ← 分类标签
└── start_all_services.sh            ← 启动脚本
```

## 🔍 验证命令

```bash
# 查看日志
docker-compose logs -f model-service

# 健康检查
curl http://localhost:5000/health

# 模型信息
curl http://localhost:5000/info

# 测试分类
curl -X POST http://localhost:5000/classify -F "image=@test.jpg"

# 完整流程测试
curl -X POST http://localhost:8888/api/images/upload -F "image=@test.jpg"
```

## 🐛 常见问题

### 问题1: state_dict加载错误

```bash
# 验证模型定义
docker exec -it model-service python << 'EOF'
import torch
from app.model_architecture import YourModel
checkpoint = torch.load('/app/models/checkpoint_best.pth')
model = YourModel(num_classes=10)

# 对比键
saved = set(checkpoint['model_state'].keys())
current = set(model.state_dict().keys())
print("缺少:", saved - current)
print("多余:", current - saved)
EOF
```

### 问题2: 模型名称不匹配

确保 `docker-compose.yml` 中的 `MODEL_ARCH` 与 `create_model()` 中的名字完全一致。

### 问题3: 服务启动失败

```bash
# 查看详细日志
docker-compose logs model-service | tail -50
```

## 📚 详细文档

- **操作步骤**: [model_service/STEP_BY_STEP.md](model_service/STEP_BY_STEP.md)
- **完整指南**: [CUSTOM_MODEL_INTEGRATION.md](CUSTOM_MODEL_INTEGRATION.md)
- **部署说明**: [DEPLOYMENT.md](DEPLOYMENT.md)

## 💡 关键配置项

| 配置项 | 说明 | 示例 |
|--------|------|------|
| MODEL_PATH | checkpoint文件路径 | `/app/models/checkpoint_best.pth` |
| MODEL_ARCH | 模型架构名称 | `my_model` |
| NUM_CLASSES | 分类类别数量 | `10` |
| DEVICE | 使用设备 | `cuda` 或 `cpu` |

## 🔄 重启服务

```bash
# 重启模型服务
docker-compose restart model-service

# 重启所有服务
docker-compose restart

# 停止并重新构建
docker-compose down
docker-compose build
docker-compose up -d
```

## 📊 服务端口

| 服务 | 端口 | 用途 |
|------|------|------|
| 模型服务 | 5000 | 图片分类API |
| Go后端 | 8888 | 完整业务API |
| PostgreSQL | 5432 | 数据库 |
| Redis | 6379 | 缓存 |

---

**需要帮助？** 查看 [STEP_BY_STEP.md](model_service/STEP_BY_STEP.md) 获取详细步骤！
