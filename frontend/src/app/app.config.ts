import { ApplicationConfig, provideBrowserGlobalErrorListeners } from '@angular/core';
import { provideAnimations } from '@angular/platform-browser/animations';
import { provideRouter } from '@angular/router';
import { registerLocaleData } from '@angular/common';
import en from '@angular/common/locales/en';
import { provideNzI18n, en_US } from 'ng-zorro-antd/i18n';
import { provideEchartsCore } from 'ngx-echarts';

import { routes } from './app.routes';

registerLocaleData(en);

export const appConfig: ApplicationConfig = {
  providers: [
    provideBrowserGlobalErrorListeners(),
    provideAnimations(),
    provideRouter(routes),
    provideNzI18n(en_US),
    provideEchartsCore({
      echarts: () => import('echarts')
    })
  ]
};
