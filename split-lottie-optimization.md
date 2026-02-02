# Lottie JSON 文件优化指南

## 方案1：优化 Lottie 文件大小

### 使用 lottie-optimizer
```bash
npm install -g @lottiefiles/lottie-optimizer

# 基础优化
lottie-optimizer -i bg.json -o bg.optimized.json

# 高级优化（可能会损失一些质量）
lottie-optimizer -i bg.json -o bg.optimized.json --compress --precision 2
```

### 使用在线工具
1. 访问 https://lottiefiles.com/
2. 上传您的 JSON 文件
3. 点击 "Optimize" 按钮
4. 可以减少 50-90% 的文件大小

### After Effects 导出优化
如果您有源文件（.aep），重新导出时：
1. 使用 Bodymovin 插件
2. 勾选 "Glyphs" 而不是 "Characters" （如果有文字）
3. 降低导出精度
4. 移除不必要的图层
5. 合并重复的资源

## 方案2：使用 CDN

### 1. 上传到 CDN（阿里云 OSS / 腾讯云 COS）
```typescript
// 从 CDN 加载，带缓存
const CDN_BASE = 'https://your-cdn.com/christmas/';

this.http.get(CDN_BASE + 'bg.json').subscribe(data => {
  // ...
});
```

### 2. 启用 CDN 压缩和缓存
- 在 CDN 控制台开启 Gzip/Brotli 压缩
- 设置缓存时间为 1 年
- 使用 HTTP/2 或 HTTP/3

## 方案3：渐进式加载

### 分割 Lottie 动画
将一个大的 Lottie 文件分割成多个小文件：
- bg-part1.json（前景）
- bg-part2.json（中景）
- bg-part3.json（背景）

先加载关键的前景，后台加载其他部分。

## 方案4：使用 WebP 替代 PNG

```bash
# 安装 cwebp
# Ubuntu/Debian
sudo apt-get install webp

# macOS
brew install webp

# 转换图片
cwebp -q 80 img1.png -o img1.webp
```

WebP 可以减少 25-35% 的文件大小。

## 方案5：图片压缩

### TinyPNG
访问 https://tinypng.com/ 批量压缩图片

### 使用命令行工具
```bash
npm install -g imagemin-cli

imagemin src/assets/img/*.png --out-dir=src/assets/img/compressed --plugin=pngquant
```

## 方案6：Service Worker 缓存

```typescript
// sw.js
self.addEventListener('install', (event) => {
  event.waitUntil(
    caches.open('christmas-v1').then((cache) => {
      return cache.addAll([
        '/assets/img/Christmas/bg.json',
        '/assets/img/img1.png',
        '/assets/img/img2.png',
        // ...
      ]);
    })
  );
});

self.addEventListener('fetch', (event) => {
  event.respondWith(
    caches.match(event.request).then((response) => {
      return response || fetch(event.request);
    })
  );
});
```

## 预期优化结果

| 优化项 | 原始大小 | 优化后 | 减少比例 |
|--------|---------|--------|---------|
| Lottie JSON | 20MB | 2-5MB | 75-90% |
| PNG 图片 | 10MB | 3-5MB | 50-70% |
| 总大小 | 30MB | 5-10MB | 67-83% |
| 加载时间 | 6.8分钟 | 10-30秒 | 95%+ |

## 推荐方案组合

1. **优先级 P0（必做）**：
   - 优化 Lottie JSON 文件（使用 lottie-optimizer）
   - 压缩所有图片（TinyPNG / WebP）
   - 启用服务器 Gzip/Brotli 压缩

2. **优先级 P1（强烈推荐）**：
   - 使用 CDN
   - 添加 Loading 进度条
   - 启用浏览器缓存

3. **优先级 P2（可选）**：
   - Service Worker 缓存
   - 渐进式加载
   - 预加载关键资源
