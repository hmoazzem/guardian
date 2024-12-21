import { HttpClient } from '@angular/common/http';
import { Component, inject } from '@angular/core';
import { rxResource } from '@angular/core/rxjs-interop';
import { GPU } from '../../interfaces'
import { CommonModule } from '@angular/common';
import { EditorComponent, ExpandableTableComponent } from 'ng-essential';
import { MatIconModule } from '@angular/material/icon';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatButtonModule } from '@angular/material/button';

@Component({
  selector: 'e-gpus',
  imports: [CommonModule, ExpandableTableComponent, EditorComponent, MatIconModule, MatProgressSpinnerModule, MatButtonModule],
  templateUrl: './gpus.component.html',
  styles: ``
})
export class GpusComponent {
  private http = inject(HttpClient);

  columns = ['name', 'codename', 'vram'];
  cellDefs = ['board.manufacturer_name', 'board.product_name', 'vram.size | json'];
  currentExpandedRow?: GPU;

  isShowDetails = false;

  /** Updates the currently expanded row in the table. */
  handleRowChange(rowData: GPU) {
    this.currentExpandedRow = rowData;
  }

  /** Toggles the display of detailed network information. */
  toggleDetails() {
    this.isShowDetails = !this.isShowDetails;
  }

  gpus = rxResource({
    loader: () => this.http.get<GPU[]>(`http://${window.location.hostname}:8081/amdgpu`)
  })

  gpuMetric = rxResource({
    loader: () => this.http.get<GPU[]>(`http://${window.location.hostname}:8081/amdgpu-metrics`)
  })
}
