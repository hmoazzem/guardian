import { bootstrapApplication } from '@angular/platform-browser';
import { appConfig } from './app/app.config';
import { AppComponent } from './app/app.component';

if ('serviceWorker' in navigator) {
  navigator.serviceWorker
    .register('/grpc-service-worker.js')
    .then((registration) => {
      console.log('Service Worker registered:', registration);
      navigator.serviceWorker.ready.then((sw) => {
        sw.active?.postMessage({ type: 'INIT', ready: true }); // Initial message to the service worker
      });
    })
    .catch((error) => console.error('Service Worker registration failed:', error));
}

bootstrapApplication(AppComponent, appConfig)
  .catch((err) => console.error(err));
