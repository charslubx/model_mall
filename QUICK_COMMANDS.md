# 快速命令参考

## 📁 文件位置：public 目录

```
public/
└── img/
    ├── Christmas/
    │   ├── bg.json              (原始文件 20MB)
    │   └── bg.optimized.json    (优化后 2-5MB)
    ├── img1.png
    ├── img2.png
    ├── img3.png
    └── img4.png
```

---

## 🚀 优化命令

### 1. 优化 Lottie JSON（基础）

```bash
lottie-optimizer -i public/img/Christmas/bg.json -o public/img/Christmas/bg.optimized.json
```

### 2. 优化 Lottie JSON（激进压缩，更小但可能损失质量）

```bash
lottie-optimizer -i public/img/Christmas/bg.json -o public/img/Christmas/bg.optimized.json --compress --precision 2
```

### 3. 优化图片

```bash
node optimize-images.js
```

### 4. 一键优化所有资源

```bash
npm run optimize:all
```

---

## 💻 代码中的路径

### TypeScript 代码

```typescript
// ✅ 正确 - public 目录下的文件
this.http.get('img/Christmas/bg.json').subscribe(data => { ... });

// ✅ 使用优化后的文件
this.http.get('img/Christmas/bg.optimized.json').subscribe(data => { ... });

// ❌ 错误 - assets 已经不用了
this.http.get('assets/img/Christmas/bg.json').subscribe(data => { ... });
```

### HTML 模板

```html
<!-- ✅ 正确 -->
<img src="img/img1.png" alt="img1">

<!-- ❌ 错误 -->
<img src="assets/img/img1.png" alt="img1">
```

---

## 📊 检查文件大小

```bash
# 查看原始文件大小
ls -lh public/img/Christmas/bg.json

# 查看优化后文件大小
ls -lh public/img/Christmas/bg.optimized.json

# 查看所有图片大小
ls -lh public/img/*.png

# 查看优化后的图片
ls -lh public/img/optimized/
```

---

## 🧪 测试

### 开发环境测试

```bash
ng serve
# 打开 http://localhost:4200
# 按 F12 打开开发者工具 → Network 选项卡
# 查看文件加载情况
```

### 构建并测试

```bash
# 构建
ng build

# 查看构建后的文件
ls -lh dist/your-app/img/Christmas/

# 使用简单服务器测试
cd dist/your-app
npx http-server -p 8080
# 打开 http://localhost:8080
```

### 测试慢速网络

1. 打开 Chrome DevTools (F12)
2. 切换到 Network 选项卡
3. 选择 "Slow 3G" 或 "Fast 3G"
4. 刷新页面，查看加载时间

---

## 🔧 angular.json 配置

确保您的 `angular.json` 包含以下配置：

```json
{
  "projects": {
    "your-app": {
      "architect": {
        "build": {
          "options": {
            "assets": [
              "src/favicon.ico",
              "src/assets",
              {
                "glob": "**/*",
                "input": "public",
                "output": "/"
              }
            ]
          }
        }
      }
    }
  }
}
```

---

## 📦 NPM 脚本

将这些添加到您的 `package.json`：

```json
{
  "scripts": {
    "optimize:lottie": "lottie-optimizer -i public/img/Christmas/bg.json -o public/img/Christmas/bg.optimized.json",
    "optimize:lottie:aggressive": "lottie-optimizer -i public/img/Christmas/bg.json -o public/img/Christmas/bg.optimized.json --compress --precision 2",
    "optimize:images": "node optimize-images.js",
    "optimize:all": "npm run optimize:lottie && npm run optimize:images",
    "build:prod": "npm run optimize:all && ng build --configuration production"
  }
}
```

使用：

```bash
npm run optimize:lottie          # 优化 Lottie
npm run optimize:lottie:aggressive  # 激进优化
npm run optimize:images          # 优化图片
npm run optimize:all             # 优化所有
npm run build:prod               # 优化并构建
```

---

## 🎯 完整流程（从头到尾）

```bash
# 步骤 1: 安装依赖
npm install -g @lottiefiles/lottie-optimizer
npm install sharp

# 步骤 2: 确认文件在正确位置
ls public/img/Christmas/bg.json

# 步骤 3: 优化所有资源
npm run optimize:all

# 步骤 4: 检查优化效果
ls -lh public/img/Christmas/bg.optimized.json

# 步骤 5: 更新代码使用优化后的文件
# 修改 .ts 文件：'img/Christmas/bg.optimized.json'

# 步骤 6: 测试开发环境
ng serve

# 步骤 7: 构建生产版本
npm run build:prod

# 步骤 8: 测试生产构建
cd dist/your-app && npx http-server -p 8080
```

---

## 📈 预期结果

| 项目 | 原始 | 优化后 | 改善 |
|------|------|--------|------|
| bg.json | 20 MB | 2-5 MB | 75-90% ↓ |
| 图片总计 | 10 MB | 3-5 MB | 50-70% ↓ |
| 首次加载 | 6.8 分钟 | 10-30 秒 | 95% ↓ |
| 二次加载 | 6.8 分钟 | 2-5 秒 | 99% ↓ |

---

## 🆘 常见问题快速解决

### 问题：404 找不到文件

```bash
# 检查 1: 文件是否存在
ls public/img/Christmas/bg.json

# 检查 2: angular.json 配置是否正确
cat angular.json | grep -A 5 "assets"

# 检查 3: 路径是否正确（不要用 assets/）
# ✅ 'img/Christmas/bg.json'
# ❌ 'assets/img/Christmas/bg.json'
```

### 问题：优化命令找不到

```bash
# 全局安装 lottie-optimizer
npm install -g @lottiefiles/lottie-optimizer

# 检查是否安装成功
lottie-optimizer --version
```

### 问题：图片优化脚本报错

```bash
# 安装依赖
npm install sharp

# 检查目录是否存在
mkdir -p public/img/optimized
```

---

## 🎉 完成！

现在您的 Lottie 动画应该快多了！

- ✅ 文件从 30MB 减少到 5-10MB
- ✅ 加载时间从 6.8 分钟减少到 10-30 秒
- ✅ 第二次访问只需 2-5 秒（有缓存）

享受飞速的加载体验吧！🚀
