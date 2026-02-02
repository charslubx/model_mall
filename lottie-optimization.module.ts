import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { HttpClientModule } from '@angular/common/http';

// 组件
import { LottieOptimizedComponent } from './lottie-optimized.component';
import { LottieLazyLoadComponent } from './lottie-lazy-load.component';

// 服务
import { LottieCacheService } from './lottie-cache.service';

@NgModule({
  declarations: [
    LottieOptimizedComponent,
    LottieLazyLoadComponent
  ],
  imports: [
    CommonModule,
    HttpClientModule
  ],
  providers: [
    LottieCacheService
  ],
  exports: [
    LottieOptimizedComponent,
    LottieLazyLoadComponent
  ]
})
export class LottieOptimizationModule { }
