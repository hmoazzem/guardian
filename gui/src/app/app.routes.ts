import { Routes } from '@angular/router';

export const routes: Routes = [
  { path: 'system', loadComponent: () => import('./features/system/system.component').then((m) => m.SystemComponent) },
  { path: 'network', loadComponent: () => import('./features/net/wg/wg.component').then((m) => m.WgComponent) },
  { path: 'gpu', loadComponent: () => import('./features/gpus/gpus.component').then((m) => m.GpusComponent)},
  { path: '', redirectTo: 'system', pathMatch: 'full'}
];
