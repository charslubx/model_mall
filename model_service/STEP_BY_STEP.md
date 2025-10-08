# 📝 自定义模型部署 - 分步操作指南

## 你现在的情况

- ✅ 有checkpoint文件（checkpoint_best.pth）
- ✅ 模型是自己改过的，不是标准模型
- ✅ checkpoint格式：`{'epoch', 'model_state', 'optimizer_state', ...}`

## 🎯 需要做的事（具体操作）

### 操作1：找到模型定义

**在你的训练代码目录中：**

```bash
# 方法1：找Python文件
ls *.py

# 方法2：搜索模型类
grep -r "class.*nn.Module" .
grep -r "def __init__" .

# 常见文件名：
# - model.py
# - models.py  
# - network.py
# - train.py
```

**找到类似这样的代码：**

```python
class YourModel(nn.Module):
    def __init__(self, ...):
        ...
    def forward(self, x):
        ...
```

**👉 全选并复制这个类（包括所有相关的helper类）**

### 操作2：粘贴到项目

**打开这个文件：**
```
/workspace/model_service/app/model_architecture.py
```

**滚动到文件最底部，在最后添加：**

```python

# ============================================
# 我的自定义模型（从训练代码复制）
# ============================================

# 粘贴你复制的模型类到这里
class YourModel(nn.Module):
    """从训练代码完整复制"""
    def __init__(self, num_classes=10):
        # ... 粘贴的内容
        pass
    
    def forward(self, x):
        # ... 粘贴的内容
        pass
```

**然后找到 `create_model` 函数（在同一个文件中）：**

```python
def create_model(model_name='resnet18', num_classes=10):
    model_name = model_name.lower()
    
    # ... 已有的代码 ...
    
    # 👇 在这里添加
    elif model_name == 'my_model':  # ← 给个简单的名字
        return YourModel(num_classes)
    
    else:
        raise ValueError(f"未知的模型名称: {model_name}")
```

**保存文件！**

### 操作3：配置环境变量

**打开这个文件：**
```
/workspace/model_service/docker-compose.yml
```

**找到 `environment:` 部分，修改这3行：**

```yaml
environment:
  - MODEL_PATH=/app/models/checkpoint_best.pth  # ← 你的文件名
  - MODEL_ARCH=my_model  # ← 使用操作2中定义的名字
  - NUM_CLASSES=10       # ← 改成你的类别数量
```

**保存文件！**

### 操作4：准备文件

```bash
# 1. 复制checkpoint文件
cp /path/to/your/checkpoint_best.pth /workspace/model_service/models/

# 2. 创建标签文件（重要：类别数量要对！）
cat > /workspace/model_service/models/labels.txt << 'EOF'
类别1
类别2
类别3
类别4
类别5
类别6
类别7
类别8
类别9
类别10
EOF
```

### 操作5：启动服务

```bash
cd /workspace
./start_all_services.sh
```

**等待大约30-60秒...**

### 操作6：验证

```bash
# 1. 检查健康状态
curl http://localhost:5000/health

# 应该看到:
# {"status":"healthy","model_loaded":true}

# 2. 测试分类
curl -X POST http://localhost:5000/classify \
  -F "image=@/path/to/test/image.jpg"
```

## ✅ 成功的标志

**查看日志：**
```bash
docker-compose logs model-service | grep -A 10 "初始化模型"
```

**应该看到：**
```
✓ 初始化模型服务...
✓ 模型路径: /app/models/checkpoint_best.pth
✓ 模型架构: my_model
✓ 类别数量: 10
✓ 使用自定义架构加载checkpoint
✓ 从checkpoint加载权重 (epoch: XX)
✓ PyTorch模型加载成功
✓ 训练轮次: XX
✓ 最佳验证损失: X.XXXX
✓ 模型初始化成功
```

## 🐛 如果失败了

### 错误1：找不到模型文件

```
FileNotFoundError: 模型文件不存在
```

**检查：**
```bash
ls -lh /workspace/model_service/models/
# 应该看到 checkpoint_best.pth
```

### 错误2：未知的模型名称

```
ValueError: 未知的模型名称: my_model
```

**原因：** `docker-compose.yml` 中的 `MODEL_ARCH` 与 `create_model` 函数中的名字不匹配

**解决：** 确保两处名字完全一致

### 错误3：模型加载失败

```
RuntimeError: Error(s) in loading state_dict
```

**原因：** 模型定义与checkpoint不匹配

**解决：** 
1. 确认复制了完整的模型类
2. 确认 `NUM_CLASSES` 正确
3. 运行验证脚本（见下方）

### 错误4：import错误

```
ImportError: cannot import name 'XXX'
```

**原因：** 模型用了其他模块，但没有复制过来

**解决：** 把相关的类也复制到 `model_architecture.py`

## 🔍 验证脚本

如果遇到 `state_dict` 加载错误，运行这个：

```bash
# 进入容器
docker exec -it model-service python << 'EOF'
import torch
from app.model_architecture import YourModel  # 改成你的类名

# 加载checkpoint
checkpoint = torch.load('/app/models/checkpoint_best.pth')
saved_keys = set(checkpoint['model_state'].keys())

# 创建模型
model = YourModel(num_classes=10)  # 改成你的类别数
current_keys = set(model.state_dict().keys())

# 对比
missing = saved_keys - current_keys
extra = current_keys - saved_keys

if missing:
    print("❌ Checkpoint中有但模型中没有的层:")
    for k in missing:
        print(f"  - {k}")

if extra:
    print("❌ 模型中有但Checkpoint中没有的层:")
    for k in extra:
        print(f"  - {k}")

if not missing and not extra:
    print("✅ 完美匹配！")
else:
    print("\n请检查模型定义是否与训练时完全一致")
EOF
```

## 📋 完整检查清单

- [ ] 已找到训练代码中的模型类定义
- [ ] 已复制到 `model_architecture.py` 末尾
- [ ] 已在 `create_model()` 中添加选项
- [ ] `MODEL_ARCH` 名字与代码中一致
- [ ] `NUM_CLASSES` 正确
- [ ] checkpoint文件已复制到 `models/` 目录
- [ ] `labels.txt` 类别数量正确
- [ ] 运行了 `./start_all_services.sh`
- [ ] 服务健康检查通过

## 💡 快速参考

**主要文件：**
```
/workspace/model_service/app/model_architecture.py  ← 粘贴模型类
/workspace/model_service/docker-compose.yml        ← 配置环境变量  
/workspace/model_service/models/                   ← 放checkpoint和labels.txt
```

**关键环境变量：**
```yaml
MODEL_PATH=/app/models/checkpoint_best.pth
MODEL_ARCH=my_model        # ← 你定义的名字
NUM_CLASSES=10             # ← 类别数量
```

**查看日志：**
```bash
docker-compose logs -f model-service
```

**重启服务：**
```bash
docker-compose restart model-service
```

---

**按照这些步骤操作，应该就能成功！** 🚀

需要帮助？检查日志中的详细错误信息。
