import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';
import type { EChartsCoreOption } from 'echarts/core';
import { NzButtonModule } from 'ng-zorro-antd/button';
import { NzCardModule } from 'ng-zorro-antd/card';
import { NzFormModule } from 'ng-zorro-antd/form';
import { NzGridModule } from 'ng-zorro-antd/grid';
import { NzIconModule } from 'ng-zorro-antd/icon';
import { NzInputModule } from 'ng-zorro-antd/input';
import { NzPaginationModule } from 'ng-zorro-antd/pagination';
import { NzSelectModule } from 'ng-zorro-antd/select';
import { NzSpaceModule } from 'ng-zorro-antd/space';
import { NgxEchartsDirective } from 'ngx-echarts';

@Component({
  selector: 'app-gratitude-page',
  standalone: true,
  imports: [
    CommonModule,
    NzButtonModule,
    NzCardModule,
    NzFormModule,
    NzGridModule,
    NzIconModule,
    NzInputModule,
    NzPaginationModule,
    NzSelectModule,
    NzSpaceModule,
    NgxEchartsDirective
  ],
  templateUrl: './gratitude.page.html',
  styleUrl: './gratitude.page.scss'
})
export class GratitudePage {
  protected pageIndex = 1;
  protected pageSize = 5;

  protected readonly messages = [
    { title: 'AAA thanked CCC', body: '感谢感谢', time: '2026-01-20 08:56:39' },
    { title: 'AAA thanked BBB', body: '感谢感谢', time: '2026-01-20 08:56:39' },
    { title: 'BBB thanked AAA', body: '多谢支持', time: '2026-01-19 14:12:10' },
    { title: 'CCC thanked BBB', body: '辛苦啦', time: '2026-01-18 09:01:22' },
    { title: 'AAA thanked BBB', body: '感谢一起排查', time: '2026-01-17 18:40:05' },
    { title: 'BBB thanked CCC', body: 'Thanks!', time: '2026-01-16 10:20:33' },
    { title: 'CCC thanked AAA', body: '很棒的建议', time: '2026-01-15 16:08:54' }
  ];

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

