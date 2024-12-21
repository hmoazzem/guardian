import { Component, OnDestroy, OnInit } from '@angular/core';
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

export interface Hwmon {
    name: string;
    composite: number;
    sensors: HwmonSensor[];
}

interface HwmonSensor {
    name: string;
    temp: number;
}

interface TimeSeriesData {
    [key: string]: {
        times: Date[];
        temps: number[];
    };
}

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
    selector: 'e-hw-temp',
    standalone: true,
    imports: [CommonModule, NgxEchartsDirective],
    template: `
      <div echarts [options]="chartOption" class="hwmon"
          style="height: 16rem; width: 100%; background-color: antiquewhite;">
      </div>
    `,
    providers: [provideEchartsCore({ echarts })]
})
export class HwTempComponent implements OnInit, OnDestroy {
    chartOption: EChartsCoreOption = {};
    private timeSeriesData: TimeSeriesData = {};
    private readonly MAX_POINTS = 100; // Maximum number of points to show

    constructor() {
        this.initializeChartOptions();
    }

    ngOnInit(): void {
        if ('serviceWorker' in navigator) {
            navigator.serviceWorker.addEventListener('message', this.handleServiceWorkerMessage);
        }
    }

    ngOnDestroy(): void {
        if ('serviceWorker' in navigator) {
            navigator.serviceWorker.removeEventListener('message', this.handleServiceWorkerMessage);
        }
    }

    private handleServiceWorkerMessage = (event: MessageEvent) => {
        const { type, payload } = event.data;
        switch (type) {
            case 'hwmon':
                this.updateChartData(payload as Hwmon);
                break;
            case 'stream_error':
                console.error('Stream Error:', payload);
                break;
            case 'stream_end':
                console.log('Stream Ended');
                break;
        }
    };

    private initializeChartOptions(): void {
        this.chartOption = {
            title: {
                text: 'HW Temp',
                left: 'left'
            },
            tooltip: {
                trigger: 'axis',
                formatter: (params: any[]) => {
                    let result = `${new Date(params[0].value[0]).toLocaleTimeString()}<br/>`;
                    params.forEach(param => {
                        result += `${param.seriesName}: ${param.value[1].toFixed(1)}°C<br/>`;
                    });
                    return result;
                }
            },
            legend: {
                data: [],
                type: 'scroll',
                top: 25
            },
            grid: {
                left: '3%',
                right: '4%',
                bottom: '10%',
                containLabel: true
            },
            xAxis: {
                type: 'time',
                splitLine: {
                    show: false
                }
            },
            yAxis: {
                type: 'value',
                name: 'Temperature (°C)',
                splitLine: {
                    show: true
                }
            },
            dataZoom: [
                {
                    type: 'slider',
                    show: true,
                    bottom: 5,
                    start: 0,
                    end: 100
                }
            ],
            series: []
        };
    }

    private updateChartData(hwmon: Hwmon): void {
        const currentTime = new Date();

        // Update time series data for each sensor
        hwmon.sensors.forEach(sensor => {
            const sensorKey = `${hwmon.name}-${sensor.name}`;
            if (!this.timeSeriesData[sensorKey]) {
                this.timeSeriesData[sensorKey] = {
                    times: [],
                    temps: []
                };
            }

            const data = this.timeSeriesData[sensorKey];
            data.times.push(currentTime);
            data.temps.push(sensor.temp);

            // Maintain maximum number of points
            if (data.times.length > this.MAX_POINTS) {
                data.times.shift();
                data.temps.shift();
            }
        });

        // Update chart series
        const series = Object.keys(this.timeSeriesData).map(sensorKey => ({
            name: sensorKey,
            type: 'line',
            showSymbol: false,
            data: this.timeSeriesData[sensorKey].times.map((time, index) => [
                time,
                this.timeSeriesData[sensorKey].temps[index]
            ]),
            emphasis: {
                focus: 'series'
            }
        }));

        // Update chart options
        this.chartOption = {
            ...this.chartOption,
            legend: {
                // ...this.chartOption.legend,
                data: Object.keys(this.timeSeriesData)
            },
            series
        };
    }
}
