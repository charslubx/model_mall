import { Component, ElementRef, ViewChild, AfterViewInit, OnDestroy } from '@angular/core';
import * as lottie from 'lottie-web';
import { LottieCacheService } from './lottie-cache.service';

interface ImageAnchor {
  id: string;
  src: string;
  x: number;
  y: number;
  w: number;
  h: number;
}

@Component({
  selector: 'app-lottie-optimized',
  template: `
    <!-- Loading 界面 -->
    <div class="loading-container" *ngIf="isLoading">
      <div class="loading-content">
        <div class="loading-spinner"></div>
        <div class="loading-text">{{ loadingMessage }}</div>
        <div class="loading-progress">
          <div class="progress-bar">
            <div class="progress-fill" [style.width.%]="loadingProgress"></div>
          </div>
          <div class="progress-text">{{ loadingProgress }}%</div>
        </div>
        <div class="loading-tip" *ngIf="isFirstLoad">
          首次加载较慢，第二次访问将秒开！
        </div>
      </div>
    </div>

    <!-- Lottie 容器 -->
    <div #container11 class="lottie-container" [class.visible]="!isLoading">
      <img *ngFor="let anchor of imageAnchors"
           [attr.data-anchor-id]="anchor.id"
           [src]="anchor.src" 
           class="overlay-image"
           [alt]="anchor.id"
           loading="lazy">
    </div>

    <!-- 缓存管理按钮（可选，仅用于测试） -->
    <div class="cache-controls" *ngIf="showCacheControls">
      <button (click)="clearCache()">清除缓存</button>
      <button (click)="showCacheInfo()">查看缓存大小</button>
    </div>
  `,
  styles: [`
    .loading-container {
      position: fixed;
      top: 0;
      left: 0;
      width: 100%;
      height: 100vh;
      background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
      display: flex;
      justify-content: center;
      align-items: center;
      z-index: 9999;
      animation: fadeIn 0.3s ease;
    }

    .loading-content {
      text-align: center;
      color: white;
    }

    .loading-spinner {
      width: 60px;
      height: 60px;
      border: 4px solid rgba(255, 255, 255, 0.3);
      border-top-color: white;
      border-radius: 50%;
      animation: spin 1s linear infinite;
      margin: 0 auto 24px;
    }

    @keyframes spin {
      to { transform: rotate(360deg); }
    }

    @keyframes fadeIn {
      from { opacity: 0; }
      to { opacity: 1; }
    }

    .loading-text {
      font-size: 20px;
      font-weight: 600;
      margin-bottom: 16px;
    }

    .loading-progress {
      width: 320px;
      margin: 0 auto;
    }

    .progress-bar {
      width: 100%;
      height: 8px;
      background: rgba(255, 255, 255, 0.3);
      border-radius: 4px;
      overflow: hidden;
      margin-bottom: 8px;
    }

    .progress-fill {
      height: 100%;
      background: white;
      transition: width 0.3s ease;
      border-radius: 4px;
    }

    .progress-text {
      font-size: 14px;
      color: rgba(255, 255, 255, 0.9);
    }

    .loading-tip {
      margin-top: 16px;
      font-size: 14px;
      color: rgba(255, 255, 255, 0.8);
      animation: pulse 2s ease-in-out infinite;
    }

    @keyframes pulse {
      0%, 100% { opacity: 0.8; }
      50% { opacity: 1; }
    }

    .lottie-container {
      position: relative;
      width: 100%;
      height: 100vh;
      overflow: hidden;
      opacity: 0;
      transition: opacity 0.5s ease;
    }

    .lottie-container.visible {
      opacity: 1;
    }

    .overlay-image {
      position: absolute;
      pointer-events: none;
    }

    .cache-controls {
      position: fixed;
      bottom: 20px;
      right: 20px;
      z-index: 1000;
      display: flex;
      gap: 10px;
    }

    .cache-controls button {
      padding: 8px 16px;
      background: white;
      border: none;
      border-radius: 4px;
      cursor: pointer;
      font-size: 12px;
      box-shadow: 0 2px 8px rgba(0,0,0,0.15);
    }

    .cache-controls button:hover {
      background: #f0f0f0;
    }
  `]
})
export class LottieOptimizedComponent implements AfterViewInit, OnDestroy {
  @ViewChild('container11', { static: false }) container11!: ElementRef;

  private resizeObserver?: ResizeObserver;
  private animation: any;

  private ORIGINAL_WIDTH = 2730;
  private ORIGINAL_HEIGHT = 1535;

  isLoading = true;
  isFirstLoad = false;
  loadingProgress = 0;
  loadingMessage = '正在加载...';
  showCacheControls = false; // 设置为 true 显示缓存控制按钮

  imageAnchors: ImageAnchor[] = [
    {
      id: 'img1',
      src: 'assets/img/img1.png',
      x: 300,
      y: 780,
      w: 150,
      h: 274
    },
    {
      id: 'img2',
      src: 'assets/img/img2.png',
      x: 500,
      y: 600,
      w: 200,
      h: 300
    },
    {
      id: 'img3',
      src: 'assets/img/img3.png',
      x: 800,
      y: 400,
      w: 180,
      h: 250
    },
    {
      id: 'img4',
      src: 'assets/img/img4.png',
      x: 1200,
      y: 900,
      w: 160,
      h: 220
    }
  ];

  constructor(private cacheService: LottieCacheService) {}

  ngAfterViewInit() {
    this.loadResourcesOptimized();
  }

  ngOnDestroy() {
    if (this.resizeObserver) {
      this.resizeObserver.disconnect();
    }
    if (this.animation) {
      this.animation.destroy();
    }
  }

  private async loadResourcesOptimized() {
    const startTime = performance.now();

    try {
      // 检查是否有缓存
      this.isFirstLoad = !(await this.checkCacheExists());
      
      if (this.isFirstLoad) {
        this.loadingMessage = '首次加载中，正在下载资源...';
      } else {
        this.loadingMessage = '从缓存加载中...';
      }

      // 步骤1: 加载 Lottie JSON（从缓存或网络）
      this.loadingProgress = 0;
      const lottieData = await this.loadLottieWithCache();
      
      this.loadingProgress = 60;
      this.loadingMessage = '正在初始化动画...';

      // 步骤2: 初始化动画
      await this.initLottieAnimation(lottieData);
      
      this.loadingProgress = 80;
      this.loadingMessage = '正在加载图片...';

      // 步骤3: 预加载图片
      await this.preloadImages();
      
      this.loadingProgress = 100;
      const loadTime = ((performance.now() - startTime) / 1000).toFixed(2);
      this.loadingMessage = `加载完成！耗时 ${loadTime} 秒`;

      console.log(`总加载时间: ${loadTime} 秒`);

      // 延迟隐藏 loading
      setTimeout(() => {
        this.isLoading = false;
        this.setupResizeObserver();
      }, 500);

    } catch (error) {
      console.error('加载失败:', error);
      this.loadingMessage = '加载失败，请刷新重试';
      this.isFirstLoad = false;
    }
  }

  private async checkCacheExists(): Promise<boolean> {
    try {
      const cacheSize = await this.cacheService.getCacheSize();
      return cacheSize > 0;
    } catch {
      return false;
    }
  }

  private loadLottieWithCache(): Promise<any> {
    return new Promise((resolve, reject) => {
      const url = 'assets/img/Christmas/bg.json';
      
      this.cacheService.loadLottieJSON(url).subscribe({
        next: (data) => {
          this.loadingProgress = 50;
          resolve(data);
        },
        error: (error) => reject(error)
      });
    });
  }

  private initLottieAnimation(data: any): Promise<void> {
    return new Promise((resolve) => {
      if (data.w && data.h) {
        this.ORIGINAL_WIDTH = data.w;
        this.ORIGINAL_HEIGHT = data.h;
      }

      this.animation = lottie.loadAnimation({
        container: this.container11.nativeElement,
        renderer: 'svg',
        loop: true,
        autoplay: true,
        animationData: data
      });

      this.animation.addEventListener('DOMLoaded', () => {
        setTimeout(() => {
          this.updateAllImagesPosition();
          resolve();
        }, 50);
      });
    });
  }

  private preloadImages(): Promise<void[]> {
    const imagePromises = this.imageAnchors.map((anchor, index) => {
      return new Promise<void>((resolve, reject) => {
        const img = new Image();
        img.onload = () => {
          const progress = 80 + ((index + 1) / this.imageAnchors.length) * 20;
          this.loadingProgress = Math.round(progress);
          resolve();
        };
        img.onerror = () => {
          console.warn(`图片加载失败: ${anchor.src}`);
          resolve(); // 即使失败也继续
        };
        img.src = anchor.src;
      });
    });

    return Promise.all(imagePromises);
  }

  private setupResizeObserver() {
    this.resizeObserver = new ResizeObserver(() => {
      requestAnimationFrame(() => {
        this.updateAllImagesPosition();
      });
    });

    this.resizeObserver.observe(this.container11.nativeElement);
  }

  private updateAllImagesPosition() {
    if (!this.container11) {
      return;
    }

    const container = this.container11.nativeElement;
    const containerRect = container.getBoundingClientRect();
    const containerWidth = containerRect.width;
    const containerHeight = containerRect.height;

    const originalRatio = this.ORIGINAL_WIDTH / this.ORIGINAL_HEIGHT;
    const containerRatio = containerWidth / containerHeight;

    let scaledWidth: number;
    let scaledHeight: number;
    let offsetX: number;
    let offsetY: number;

    if (containerRatio > originalRatio) {
      scaledHeight = containerHeight;
      scaledWidth = scaledHeight * originalRatio;
      offsetX = (containerWidth - scaledWidth) / 2;
      offsetY = 0;
    } else {
      scaledWidth = containerWidth;
      scaledHeight = scaledWidth / originalRatio;
      offsetX = 0;
      offsetY = (containerHeight - scaledHeight) / 2;
    }

    const scale = scaledWidth / this.ORIGINAL_WIDTH;

    const imgElements = container.querySelectorAll('.overlay-image');
    imgElements.forEach((imgElement: HTMLImageElement, index: number) => {
      if (index < this.imageAnchors.length) {
        const anchor = this.imageAnchors[index];
        this.updateSingleImagePosition(imgElement, anchor, scale, offsetX, offsetY);
      }
    });
  }

  private updateSingleImagePosition(
    img: HTMLImageElement,
    anchor: ImageAnchor,
    scale: number,
    offsetX: number,
    offsetY: number
  ) {
    const imgLeft = anchor.x * scale + offsetX;
    const imgTop = anchor.y * scale + offsetY;
    const imgWidth = anchor.w * scale;
    const imgHeight = anchor.h * scale;

    img.style.left = `${imgLeft}px`;
    img.style.top = `${imgTop}px`;
    img.style.width = `${imgWidth}px`;
    img.style.height = `${imgHeight}px`;
  }

  // 缓存管理方法
  async clearCache() {
    try {
      await this.cacheService.clearCache();
      alert('缓存已清除！刷新页面将重新下载资源。');
    } catch (error) {
      console.error('清除缓存失败:', error);
    }
  }

  async showCacheInfo() {
    try {
      const size = await this.cacheService.getCacheSize();
      const sizeMB = (size / 1024 / 1024).toFixed(2);
      alert(`当前缓存大小: ${sizeMB} MB`);
    } catch (error) {
      console.error('获取缓存信息失败:', error);
    }
  }
}
