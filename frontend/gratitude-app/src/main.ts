import { registerLocaleData } from '@angular/common';
import zh from '@angular/common/locales/zh';
import { bootstrapApplication } from '@angular/platform-browser';
import { appConfig } from './app/app.config';
import { App } from './app/app';

registerLocaleData(zh);

bootstrapApplication(App, appConfig)
  .catch((err) => console.error(err));
