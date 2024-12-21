import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import {
  LegendComponent,
  TitleComponent,
  TooltipComponent,
  GridComponent,
  DataZoomComponent
} from 'echarts/components';
import { CanvasRenderer } from 'echarts/renderers';
import type { EChartsCoreOption } from 'echarts/core';
import { NgxEchartsDirective, provideEchartsCore } from 'ngx-echarts';
import * as echarts from 'echarts/core';
import { LineChart } from 'echarts/charts';

// Register components with echarts
echarts.use([
  LegendComponent,
  TooltipComponent,
  LineChart,
  TitleComponent,
  CanvasRenderer,
  GridComponent,
  DataZoomComponent
]);

@Component({
  selector: 'e-cpu-utilz',
  imports: [CommonModule, NgxEchartsDirective],
  template: `
    <div *ngIf="isChartReady" echarts [options]="chartOptions" [merge]="updateOptions" class="cpu_utilization"
      style="height: 16rem; width: 100%; margin: 0; background-color: antiquewhite;">
    </div>
  `,
  styles: ``,
  providers: [
    provideEchartsCore({ echarts }),
  ]
})
export class CpuUtilzComponent {
  // Flag to control chart rendering
  isChartReady = false;

  // Chart configuration options
  chartOptions!: EChartsCoreOption;

  // Update options for merging new data
  updateOptions: EChartsCoreOption | null = null;

  // Store CPU clock data
  private cpuClockData: number[][] = [];

  // Tracking data points for each core
  private coreDataSeries: number[][] = [];

  // Timestamps for x-axis
  private timestamps: string[] = [];

  // Number of cores
  private numberOfCores = 0;

  constructor() {}

  ngOnInit(): void {
    // Listen for CPU clock data from service worker
    if ('serviceWorker' in navigator) {
      navigator.serviceWorker.addEventListener('message', (event: MessageEvent) => {
        const { type, payload } = event.data;
        switch (type) {
          case 'cpu_utilization':
            // First time receiving data, initialize the chart
            if (!this.isChartReady) {
              this.initializeChart(payload as Number[]);
            } else {
              this.updateChartData(payload as Number[]);
            }
            break;
          case 'stream_error':
            console.error('Stream Error:', payload);
            break;
          case 'stream_end':
            console.log('Stream Ended');
            break;
        }
      });

    }
  }

  ngOnDestroy(): void {
    // Remove event listener if needed
    if ('serviceWorker' in navigator) {
      navigator.serviceWorker.removeEventListener('message', this.handleServiceWorkerMessage);
    }
  }

  private handleServiceWorkerMessage = (event: MessageEvent) => {
    const { type, payload } = event.data;

    if (type === 'cpu_utilization') {
      if (!this.isChartReady) {
        this.initializeChart(payload as Number[]);
      } else {
        this.updateChartData(payload as Number[]);
      }
    }
  }

  private initializeChart(initialData: Number[]): void {
    // Dynamically determine number of cores
    this.numberOfCores = initialData.length;

    // Initialize core data series
    this.coreDataSeries = Array.from({ length: this.numberOfCores }, () => []);

    // Prepare chart options with dynamic number of cores
    this.chartOptions = {
      title: {
        text: `CPU Utilization (${this.numberOfCores} vCPU)`,
        left: 'center'
      },
      tooltip: {
        trigger: 'axis',
        axisPointer: {
          type: 'line'
        }
      },
      grid: {
        left: '4%',
        right: '2%',
        top: '12%',
        bottom: '10%'
      },
      dataZoom: [
        {
          type: 'slider',
          start: 0,
          end: 100,
          bottom: 10,
        },
        {
          type: 'inside',
          start: 0,
          end: 100
        }
      ],
      xAxis: {
        type: 'category',
        boundaryGap: false,
        data: []
      },
      yAxis: {
        type: 'value',
        name: 'CPU Utilization (%)',
        min: 0,
        scale: true,
      },
      series: Array.from({ length: this.numberOfCores }, (_, i) => ({
        name: `CPU${i}`,
        type: 'line',
        showSymbol: false,
        smooth: true,
        data: []
      }))
    };

    // Mark chart as ready
    this.isChartReady = true;

    // Update with initial data
    this.updateChartData(initialData);
  }

  private updateChartData(newDataPoint: Number[]): void {
    // Ensure we have the correct number of cores
    if (newDataPoint.length !== this.numberOfCores) {
      console.warn(`Received data for ${newDataPoint.length} cores, expected ${this.numberOfCores}`);
      return;
    }

    // Convert Number[] to number[] and round to 2 decimal places
    const dataPoint = newDataPoint.map(val => Number(val.toFixed(2)));

    // Store the entire data point
    this.cpuClockData.push(dataPoint);

    // Store timestamp
    this.timestamps.push(new Date().toLocaleTimeString());

    // Update data for each core
    dataPoint.forEach((clockSpeed, coreIndex) => {
      this.coreDataSeries[coreIndex].push(clockSpeed);
    });

    // Calculate the min and max values from the updated data
    const allValues = this.coreDataSeries.flat();
    const dynamicMinY = Math.min(...allValues, 0);  // Ensure min is at least 0
    const dynamicMaxY = Math.max(...allValues);

    // Prepare update options for merging
    this.updateOptions = {
      series: this.coreDataSeries.map((coreData, index) => ({
        name: `CPU${index}`,
        data: coreData,
        smooth: true,
      })),
      xAxis: {
        data: this.timestamps
      },
      yAxis: {
        min: dynamicMinY,
        max: dynamicMaxY
      }
    };
  }
}
