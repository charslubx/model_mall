import { Component } from '@angular/core';
import { GratitudeComponent } from './gratitude/gratitude.component';

@Component({
  selector: 'app-root',
  imports: [GratitudeComponent],
  templateUrl: './app.html',
  styleUrl: './app.css'
})
export class App {
}
