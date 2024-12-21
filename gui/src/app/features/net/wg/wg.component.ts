import { CommonModule } from '@angular/common';
import { HttpClient } from '@angular/common/http';
import { Component, inject, Inject } from '@angular/core';
import { rxResource } from '@angular/core/rxjs-interop'

interface Peer {
  public_key: string
  allowed_ips: string[]
  endpoint: string | null
  last_handshake_time: Date | null
  receive_bytes: number
  transmit_bytes: number
}

interface Device {
  name: string
  type: string
  public_key: string
  private_key?: string
  listen_port: number | null
  peers: Peer[]
} 

@Component({
  selector: 'e-wg',
  imports: [CommonModule],
  templateUrl: './wg.component.html',
  styles: ``
})
export class WgComponent {
  private http = inject(HttpClient);

  devices = rxResource({
    loader: () => this.http.get<Device[]>(`http://${window.location.hostname}:8081/net/wg-devices`)
  })
}
