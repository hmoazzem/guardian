import { Component, inject } from '@angular/core';
import { HwTempComponent } from './hw-temp.component';
import { CpuUtilzComponent } from './cpu-utilz.component';
import { CpuClockComponent } from './cpu-clock.component';
import { CommonModule } from '@angular/common';
import { rxResource } from '@angular/core/rxjs-interop';
import { HttpClient } from '@angular/common/http';

@Component({
  selector: 'e-system',
  imports: [CommonModule, CpuUtilzComponent, CpuClockComponent, HwTempComponent],
  template: `
    <e-cpu-utilz />
    <br>
    <e-cpu-clock />
    <br>
    <e-hw-temp />
    <br>
    <pre><code>{{ sysInfo.value() | json }}</code></pre>
  `,
  styles: ``
})
export class SystemComponent {
  private http = inject(HttpClient);
  
  sysInfo = rxResource({
    loader: () => this.http.get<any>(`http://${window.location.hostname}:8081/sysinfo`)
  })
}
