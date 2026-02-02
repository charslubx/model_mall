import { Component, ElementRef, ViewChild, AfterViewInit, OnDestroy } from '@angular/core';
import { HttpClient, HttpEventType } from '@angular/common/http';
import * as lottie from 'lottie-web';

interface ImageAnchor {
  id: string;
  src: string;
  x: number;
  y: number;
  w: number;
  h: number;
}

@Component({
  selector: 'app-lottie-lazy-load',
  template: `
    <div class="loading-container" *ngIf="isLoading">
      <div class="loading-spinner"></div>
      <div class="loading-text">
        加载中... {{ loadingProgress }}%
      </div>
      <div class="loading-bar">
        <div class="loading-bar-fill" [style.width.%]="loadingProgress"></div>
      </div>
      <div class="loading-details">
        {{ loadingMessage }}
      </div>
    </div>

    <div #container11 class="lottie-container" [class.hidden]="isLoading">
      <img *ngFor="let anchor of imageAnchors"
           [attr.data-anchor-id]="anchor.id"
           [src]="anchor.src" 
           class="overlay-image"
           [alt]="anchor.id"
           loading="lazy">
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
      flex-direction: column;
      justify-content: center;
      align-items: center;
      z-index: 9999;
    }

    .loading-spinner {
      width: 60px;
      height: 60px;
      border: 4px solid rgba(255, 255, 255, 0.3);
      border-top-color: white;
      border-radius: 50%;
      animation: spin 1s linear infinite;
      margin-bottom: 20px;
    }

    @keyframes spin {
      to { transform: rotate(360deg); }
    }

    .loading-text {
      color: white;
      font-size: 24px;
      font-weight: bold;
      margin-bottom: 20px;
    }

    .loading-bar {
      width: 300px;
      height: 8px;
      background: rgba(255, 255, 255, 0.3);
      border-radius: 4px;
      overflow: hidden;
      margin-bottom: 10px;
    }

    .loading-bar-fill {
      height: 100%;
      background: white;
      transition: width 0.3s ease;
    }

    .loading-details {
      color: rgba(255, 255, 255, 0.8);
      font-size: 14px;
      margin-top: 10px;
    }

    .lottie-container {
      position: relative;
      width: 100%;
      height: 100vh;
      overflow: hidden;
      transition: opacity 0.5s ease;
    }

    .lottie-container.hidden {
      opacity: 0;
      pointer-events: none;
    }

    .overlay-image {
      position: absolute;
      pointer-events: none;
    }
  `]
})
export class LottieLazyLoadComponent implements AfterViewInit, OnDestroy {
  @ViewChild('container11', { static: false }) container11!: ElementRef;

  private resizeObserver?: ResizeObserver;
  private animation: any;

  private ORIGINAL_WIDTH = 2730;
  private ORIGINAL_HEIGHT = 1535;

  isLoading = true;
  loadingProgress = 0;
  loadingMessage = '准备加载资源...';

  imageAnchors: ImageAnchor[] = [
    {
      id: 'img1',
      src: 'img/img1.png',  // 使用 public 目录下的路径
      x: 300,
      y: 780,
      w: 150,
      h: 274
    },
    {
      id: 'img2',
      src: 'img/img2.png',
      x: 500,
      y: 600,
      w: 200,
      h: 300
    },
    {
      id: 'img3',
      src: 'img/img3.png',
      x: 800,
      y: 400,
      w: 180,
      h: 250
    },
    {
      id: 'img4',
      src: 'img/img4.png',
      x: 1200,
      y: 900,
      w: 160,
      h: 220
    }
  ];

  constructor(private http: HttpClient) {}

  ngAfterViewInit() {
    this.loadResourcesWithProgress();
  }

  ngOnDestroy() {
    if (this.resizeObserver) {
      this.resizeObserver.disconnect();
    }
    if (this.animation) {
      this.animation.destroy();
    }
  }

  private async loadResourcesWithProgress() {
    try {
      // 加载 Lottie JSON（带进度）
      this.loadingMessage = '正在加载动画文件...';
      const lottieData = await this.loadLottieWithProgress();
      
      this.loadingProgress = 50;
      this.loadingMessage = '正在初始化动画...';
      
      // 初始化 Lottie 动画
      await this.initLottieAnimation(lottieData);
      
      this.loadingProgress = 80;
      this.loadingMessage = '正在加载图片资源...';
      
      // 预加载图片
      await this.preloadImages();
      
      this.loadingProgress = 100;
      this.loadingMessage = '加载完成！';
      
      // 延迟一点再隐藏loading，让用户看到100%
      setTimeout(() => {
        this.isLoading = false;
        this.setupResizeObserver();
      }, 500);
      
    } catch (error) {
      console.error('资源加载失败:', error);
      this.loadingMessage = '加载失败，请刷新页面重试';
    }
  }

  private loadLottieWithProgress(): Promise<any> {
    return new Promise((resolve, reject) => {
      // 使用 public 目录下的路径
      this.http.get('img/Christmas/bg.json', {
        reportProgress: true,
        observe: 'events',
        responseType: 'json'
      }).subscribe({
        next: (event) => {
          if (event.type === HttpEventType.DownloadProgress) {
            if (event.total) {
              const progress = Math.round((event.loaded / event.total) * 50);
              this.loadingProgress = progress;
              
              // 显示文件大小信息
              const loadedMB = (event.loaded / 1024 / 1024).toFixed(2);
              const totalMB = (event.total / 1024 / 1024).toFixed(2);
              this.loadingMessage = `正在下载动画文件... ${loadedMB}MB / ${totalMB}MB`;
            }
          } else if (event.type === HttpEventType.Response) {
            resolve(event.body);
          }
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
          this.loadingProgress = 80 + (index + 1) / this.imageAnchors.length * 20;
          resolve();
        };
        img.onerror = reject;
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
}
