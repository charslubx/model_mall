# Public 文件夹使用指南

## 文件结构

```
your-angular-project/
├── public/
│   └── img/
│       ├── Christmas/
│       │   ├── bg.json              # 原始 Lottie JSON (20MB)
│       │   └── bg.optimized.json    # 优化后的文件 (2-5MB)
│       ├── img1.png
│       ├── img2.png
│       ├── img3.png
│       └── img4.png
└── src/
    └── app/
        └── your-component/
```

## 1. 优化命令（针对 public 目录）

```bash
# 优化 Lottie JSON
lottie-optimizer -i public/img/Christmas/bg.json -o public/img/Christmas/bg.optimized.json

# 如果需要更激进的压缩
lottie-optimizer -i public/img/Christmas/bg.json -o public/img/Christmas/bg.optimized.json --compress --precision 2
```

## 2. 图片优化配置

修改 `optimize-images.js` 配置：

```javascript
const CONFIG = {
  inputDir: './public/img',              // 输入目录改为 public
  outputDir: './public/img/optimized',   // 输出也在 public 下
  quality: 85,
  formats: ['.png', '.jpg', '.jpeg'],
  convertToWebP: true,
  webpQuality: 80
};
```

运行优化：
```bash
node optimize-images.js
```

## 3. Angular 代码中访问 public 文件

### 方式 1：直接路径（推荐）

```typescript
// Lottie JSON
this.http.get('img/Christmas/bg.json').subscribe(data => {
  // ...
});

// 图片
<img src="img/img1.png" alt="img1">
```

### 方式 2：绝对路径

```typescript
// Lottie JSON
this.http.get('/img/Christmas/bg.json').subscribe(data => {
  // ...
});

// 图片
<img src="/img/img1.png" alt="img1">
```

## 4. angular.json 配置

确保 `angular.json` 中配置了 public 目录：

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

这样配置后，`public` 目录下的文件会被复制到构建输出的根目录。

## 5. 完整的 Package.json 脚本

```json
{
  "scripts": {
    "optimize:lottie": "lottie-optimizer -i public/img/Christmas/bg.json -o public/img/Christmas/bg.optimized.json",
    "optimize:images": "node optimize-images.js",
    "optimize:all": "npm run optimize:lottie && npm run optimize:images",
    "build:prod": "npm run optimize:all && ng build --configuration production"
  }
}
```

## 6. 使用优化后的文件

如果优化效果满意，修改代码使用优化后的文件：

```typescript
// 使用优化后的 JSON
this.http.get('img/Christmas/bg.optimized.json').subscribe(data => {
  lottie.loadAnimation({
    container: this.container11.nativeElement,
    renderer: 'svg',
    loop: true,
    autoplay: true,
    animationData: data
  });
});
```

## 7. 验证文件是否正确加载

### 开发环境测试

```bash
ng serve
```

打开浏览器开发者工具（F12），查看 Network 选项卡：
- 应该能看到 `img/Christmas/bg.json` 的请求
- 检查文件大小和加载时间
- 确保状态码是 200

### 构建后测试

```bash
ng build
cd dist/your-app
npx http-server -p 8080
```

访问 `http://localhost:8080` 测试。

## 8. 常见问题

### Q: 404 错误，找不到文件

**A:** 检查以下几点：
1. 确认 `angular.json` 中正确配置了 public 目录
2. 路径不要以 `/` 开头（除非你确实需要绝对路径）
3. 文件名和路径大小写是否正确

### Q: 本地开发可以，打包后找不到文件

**A:** 检查构建配置：
```json
// angular.json
"assets": [
  {
    "glob": "**/*",
    "input": "public",
    "output": "/"  // 确保输出到根目录
  }
]
```

### Q: 如何确认文件被正确复制到 dist 目录？

**A:** 构建后检查：
```bash
ng build
ls -lh dist/your-app/img/Christmas/
```

应该能看到 `bg.json` 文件。

## 9. 性能优化检查清单

- [ ] 已运行 `lottie-optimizer` 优化 JSON 文件
- [ ] 已运行图片优化脚本
- [ ] 检查优化后的文件大小（应该减少 70% 以上）
- [ ] 在浏览器中测试加载时间
- [ ] 检查 Network 面板，确认文件正确加载
- [ ] 测试不同网络速度下的表现（Chrome DevTools → Network → Slow 3G）

## 10. 快速开始

```bash
# 1. 安装工具
npm install -g @lottiefiles/lottie-optimizer
npm install sharp

# 2. 优化所有文件
npm run optimize:all

# 3. 检查优化效果
ls -lh public/img/Christmas/

# 4. 启动开发服务器测试
ng serve

# 5. 构建生产版本
npm run build:prod
```

完成！现在您的 Lottie 动画应该加载速度快多了！🚀
