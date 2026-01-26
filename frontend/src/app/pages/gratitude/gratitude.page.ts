import { Component } from '@angular/core';
import type { EChartsCoreOption } from 'echarts/core';
import { NzButtonModule } from 'ng-zorro-antd/button';
import { NzCardModule } from 'ng-zorro-antd/card';
import { NzFormModule } from 'ng-zorro-antd/form';
import { NzGridModule } from 'ng-zorro-antd/grid';
import { NzIconModule } from 'ng-zorro-antd/icon';
import { NzInputModule } from 'ng-zorro-antd/input';
import { NzSelectModule } from 'ng-zorro-antd/select';
import { NzSpaceModule } from 'ng-zorro-antd/space';
import { NgxEchartsDirective } from 'ngx-echarts';

@Component({
  selector: 'app-gratitude-page',
  standalone: true,
  imports: [
    NzButtonModule,
    NzCardModule,
    NzFormModule,
    NzGridModule,
    NzIconModule,
    NzInputModule,
    NzSelectModule,
    NzSpaceModule,
    NgxEchartsDirective
  ],
  templateUrl: './gratitude.page.html',
  styleUrl: './gratitude.page.scss'
})
export class GratitudePage {
  protected readonly leaderboardOption: EChartsCoreOption = {
    tooltip: { trigger: 'axis', axisPointer: { type: 'shadow' } },
    legend: { data: ['Appreciations Received'], top: 0 },
    grid: { top: 40, left: 30, right: 20, bottom: 30, containLabel: true },
    xAxis: {
      type: 'category',
      data: ['BBB', 'CCC'],
      axisTick: { alignWithLabel: true }
    },
    yAxis: { type: 'value', minInterval: 1 },
    series: [
      {
        name: 'Appreciations Received',
        type: 'bar',
        data: [2, 1],
        itemStyle: {
          color: 'rgba(24, 144, 255, 0.45)',
          borderColor: 'rgba(24, 144, 255, 0.6)',
          borderWidth: 1
        }
      }
    ]
  };
}

