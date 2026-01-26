import { CommonModule } from '@angular/common';
import { Component, ElementRef, ViewChild, inject } from '@angular/core';
import { FormBuilder, ReactiveFormsModule, Validators } from '@angular/forms';
import { Chart, type ChartConfiguration } from 'chart.js/auto';

export type GratitudeMessage = {
  sender: string;
  receiver: string;
  reason: string;
  timestamp: string;
};

@Component({
  selector: 'app-gratitude',
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule],
  templateUrl: './gratitude.component.html',
  styleUrl: './gratitude.component.css',
})
export class GratitudeComponent {
  private readonly fb = inject(FormBuilder);

  // 组件内 Demo 数据：后续可替换为从后端 API 拉取
  users: string[] = ['Alice', 'Bob', 'Charlie', 'David'];
  messages: GratitudeMessage[] = [];

  chartLabels: string[] = [];
  chartData: number[] = [];

  submitting = false;

  form = this.fb.group({
    sender: ['', Validators.required],
    receiver: ['', Validators.required],
    reason: ['', [Validators.required, Validators.maxLength(500)]],
  });

  private chart?: Chart;
  private chartCanvas?: ElementRef<HTMLCanvasElement>;

  @ViewChild('gratitudeChart')
  set gratitudeChartRef(ref: ElementRef<HTMLCanvasElement> | undefined) {
    this.chartCanvas = ref;
    queueMicrotask(() => this.renderChart());
  }

  submit(): void {
    if (this.submitting) return;
    if (this.form.invalid) {
      this.form.markAllAsTouched();
      return;
    }

    const { sender, receiver, reason } = this.form.getRawValue();
    if (!sender || !receiver || !reason) return;

    this.submitting = true;
    try {
      const timestamp = new Date().toLocaleString();
      this.messages = [{ sender, receiver, reason, timestamp }, ...this.messages];

      this.recomputeLeaderboard();

      this.form.reset({ sender: '', receiver: '', reason: '' });
    } finally {
      this.submitting = false;
    }
  }

  trackByMessage(_idx: number, msg: GratitudeMessage): string {
    return `${msg.timestamp}-${msg.sender}-${msg.receiver}-${msg.reason}`;
  }

  private recomputeLeaderboard(): void {
    const counts = new Map<string, number>();
    for (const msg of this.messages) {
      counts.set(msg.receiver, (counts.get(msg.receiver) ?? 0) + 1);
    }

    const sorted = [...counts.entries()].sort((a, b) => b[1] - a[1]);
    this.chartLabels = sorted.map(([name]) => name);
    this.chartData = sorted.map(([, count]) => count);

    this.renderChart();
  }

  private renderChart(): void {
    if (!this.chartCanvas) return;
    if (!this.chartLabels?.length) {
      this.destroyChart();
      return;
    }

    const ctx = this.chartCanvas.nativeElement.getContext('2d');
    if (!ctx) return;

    const config: ChartConfiguration<'bar'> = {
      type: 'bar',
      data: {
        labels: this.chartLabels,
        datasets: [
          {
            label: 'Appreciations Received',
            data: this.chartData,
            backgroundColor: 'rgba(54, 162, 235, 0.6)',
            borderColor: 'rgba(54, 162, 235, 1)',
            borderWidth: 1,
          },
        ],
      },
      options: {
        responsive: true,
        scales: {
          y: {
            beginAtZero: true,
            ticks: {
              stepSize: 1,
            },
          },
        },
      },
    };

    this.destroyChart();
    this.chart = new Chart(ctx, config);
  }

  private destroyChart(): void {
    if (!this.chart) return;
    this.chart.destroy();
    this.chart = undefined;
  }
}

