/**
 * 图片批量优化脚本
 * 使用方法：
 * 1. 安装依赖: npm install sharp
 * 2. 运行脚本: node optimize-images.js
 */

const fs = require('fs');
const path = require('path');
const sharp = require('sharp');

// 配置
const CONFIG = {
  inputDir: './public/img',               // 输入目录（修改为 public）
  outputDir: './public/img/optimized',    // 输出目录（修改为 public）
  quality: 85,                            // 压缩质量 (1-100)
  formats: ['.png', '.jpg', '.jpeg'],     // 支持的格式
  convertToWebP: true,                    // 是否转换为 WebP
  webpQuality: 80                         // WebP 质量
};

// 创建输出目录
if (!fs.existsSync(CONFIG.outputDir)) {
  fs.mkdirSync(CONFIG.outputDir, { recursive: true });
}

// 获取文件大小（MB）
function getFileSizeMB(filePath) {
  const stats = fs.statSync(filePath);
  return (stats.size / (1024 * 1024)).toFixed(2);
}

// 优化单个图片
async function optimizeImage(inputPath, filename) {
  const ext = path.extname(filename).toLowerCase();
  const basename = path.basename(filename, ext);
  
  try {
    const originalSize = getFileSizeMB(inputPath);
    console.log(`\n处理: ${filename} (${originalSize} MB)`);

    const image = sharp(inputPath);
    const metadata = await image.metadata();
    console.log(`  原始尺寸: ${metadata.width}x${metadata.height}`);

    // 压缩原格式
    let outputPath;
    if (ext === '.png') {
      outputPath = path.join(CONFIG.outputDir, filename);
      await image
        .png({ quality: CONFIG.quality, compressionLevel: 9 })
        .toFile(outputPath);
    } else if (ext === '.jpg' || ext === '.jpeg') {
      outputPath = path.join(CONFIG.outputDir, filename);
      await image
        .jpeg({ quality: CONFIG.quality, progressive: true })
        .toFile(outputPath);
    }

    const compressedSize = getFileSizeMB(outputPath);
    const saved = ((originalSize - compressedSize) / originalSize * 100).toFixed(1);
    console.log(`  压缩后: ${compressedSize} MB (节省 ${saved}%)`);

    // 转换为 WebP
    if (CONFIG.convertToWebP) {
      const webpPath = path.join(CONFIG.outputDir, `${basename}.webp`);
      await sharp(inputPath)
        .webp({ quality: CONFIG.webpQuality })
        .toFile(webpPath);
      
      const webpSize = getFileSizeMB(webpPath);
      const webpSaved = ((originalSize - webpSize) / originalSize * 100).toFixed(1);
      console.log(`  WebP: ${webpSize} MB (节省 ${webpSaved}%)`);
    }

    return {
      filename,
      originalSize: parseFloat(originalSize),
      compressedSize: parseFloat(compressedSize),
      webpSize: CONFIG.convertToWebP ? parseFloat(getFileSizeMB(
        path.join(CONFIG.outputDir, `${basename}.webp`)
      )) : 0
    };

  } catch (error) {
    console.error(`  错误: ${error.message}`);
    return null;
  }
}

// 批量处理
async function processDirectory() {
  console.log('='.repeat(60));
  console.log('图片批量优化工具');
  console.log('='.repeat(60));
  console.log(`输入目录: ${CONFIG.inputDir}`);
  console.log(`输出目录: ${CONFIG.outputDir}`);
  console.log(`压缩质量: ${CONFIG.quality}`);
  console.log('='.repeat(60));

  const files = fs.readdirSync(CONFIG.inputDir);
  const imageFiles = files.filter(file => 
    CONFIG.formats.includes(path.extname(file).toLowerCase())
  );

  if (imageFiles.length === 0) {
    console.log('未找到图片文件！');
    return;
  }

  console.log(`找到 ${imageFiles.length} 个图片文件\n`);

  const results = [];
  for (const file of imageFiles) {
    const inputPath = path.join(CONFIG.inputDir, file);
    const result = await optimizeImage(inputPath, file);
    if (result) {
      results.push(result);
    }
  }

  // 统计信息
  console.log('\n' + '='.repeat(60));
  console.log('优化完成！统计信息：');
  console.log('='.repeat(60));

  const totalOriginal = results.reduce((sum, r) => sum + r.originalSize, 0);
  const totalCompressed = results.reduce((sum, r) => sum + r.compressedSize, 0);
  const totalWebP = results.reduce((sum, r) => sum + r.webpSize, 0);

  console.log(`总原始大小: ${totalOriginal.toFixed(2)} MB`);
  console.log(`压缩后大小: ${totalCompressed.toFixed(2)} MB`);
  console.log(`节省空间: ${(totalOriginal - totalCompressed).toFixed(2)} MB (${((totalOriginal - totalCompressed) / totalOriginal * 100).toFixed(1)}%)`);
  
  if (CONFIG.convertToWebP) {
    console.log(`\nWebP 总大小: ${totalWebP.toFixed(2)} MB`);
    console.log(`WebP 节省: ${(totalOriginal - totalWebP).toFixed(2)} MB (${((totalOriginal - totalWebP) / totalOriginal * 100).toFixed(1)}%)`);
  }

  console.log('='.repeat(60));
}

// 运行
processDirectory().catch(console.error);
