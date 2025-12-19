import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable, from, of } from 'rxjs';
import { switchMap, tap, catchError } from 'rxjs/operators';

@Injectable({
  providedIn: 'root'
})
export class LottieCacheService {
  private dbName = 'LottieCache';
  private storeName = 'animations';
  private db: IDBDatabase | null = null;

  constructor(private http: HttpClient) {
    this.initDB();
  }

  private async initDB(): Promise<void> {
    return new Promise((resolve, reject) => {
      const request = indexedDB.open(this.dbName, 1);

      request.onerror = () => reject(request.error);
      request.onsuccess = () => {
        this.db = request.result;
        resolve();
      };

      request.onupgradeneeded = (event: any) => {
        const db = event.target.result;
        if (!db.objectStoreNames.contains(this.storeName)) {
          db.createObjectStore(this.storeName, { keyPath: 'url' });
        }
      };
    });
  }

  /**
   * 加载 Lottie JSON，优先从缓存读取
   */
  loadLottieJSON(url: string): Observable<any> {
    return from(this.initDB()).pipe(
      switchMap(() => from(this.getFromCache(url))),
      switchMap((cached) => {
        if (cached) {
          console.log('从缓存加载 Lottie:', url);
          return of(cached.data);
        } else {
          console.log('从网络加载 Lottie:', url);
          return this.http.get(url).pipe(
            tap((data) => this.saveToCache(url, data))
          );
        }
      }),
      catchError((error) => {
        console.error('加载失败，尝试从网络重新加载:', error);
        return this.http.get(url).pipe(
          tap((data) => this.saveToCache(url, data))
        );
      })
    );
  }

  /**
   * 从 IndexedDB 缓存读取
   */
  private async getFromCache(url: string): Promise<any> {
    if (!this.db) {
      await this.initDB();
    }

    return new Promise((resolve, reject) => {
      const transaction = this.db!.transaction([this.storeName], 'readonly');
      const store = transaction.objectStore(this.storeName);
      const request = store.get(url);

      request.onsuccess = () => {
        const result = request.result;
        if (result && this.isCacheValid(result)) {
          resolve(result);
        } else {
          resolve(null);
        }
      };

      request.onerror = () => reject(request.error);
    });
  }

  /**
   * 保存到 IndexedDB 缓存
   */
  private async saveToCache(url: string, data: any): Promise<void> {
    if (!this.db) {
      await this.initDB();
    }

    return new Promise((resolve, reject) => {
      const transaction = this.db!.transaction([this.storeName], 'readwrite');
      const store = transaction.objectStore(this.storeName);
      const cacheData = {
        url: url,
        data: data,
        timestamp: Date.now(),
        version: '1.0'
      };

      const request = store.put(cacheData);

      request.onsuccess = () => {
        console.log('缓存保存成功:', url);
        resolve();
      };

      request.onerror = () => reject(request.error);
    });
  }

  /**
   * 检查缓存是否有效（7天有效期）
   */
  private isCacheValid(cached: any): boolean {
    const maxAge = 7 * 24 * 60 * 60 * 1000; // 7 天
    return (Date.now() - cached.timestamp) < maxAge;
  }

  /**
   * 清除缓存
   */
  async clearCache(): Promise<void> {
    if (!this.db) {
      await this.initDB();
    }

    return new Promise((resolve, reject) => {
      const transaction = this.db!.transaction([this.storeName], 'readwrite');
      const store = transaction.objectStore(this.storeName);
      const request = store.clear();

      request.onsuccess = () => {
        console.log('缓存已清除');
        resolve();
      };

      request.onerror = () => reject(request.error);
    });
  }

  /**
   * 获取缓存大小
   */
  async getCacheSize(): Promise<number> {
    if (!this.db) {
      await this.initDB();
    }

    return new Promise((resolve, reject) => {
      const transaction = this.db!.transaction([this.storeName], 'readonly');
      const store = transaction.objectStore(this.storeName);
      const request = store.getAll();

      request.onsuccess = () => {
        const items = request.result;
        const size = items.reduce((total: number, item: any) => {
          return total + JSON.stringify(item).length;
        }, 0);
        resolve(size);
      };

      request.onerror = () => reject(request.error);
    });
  }
}
