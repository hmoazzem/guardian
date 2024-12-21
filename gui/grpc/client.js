import { Empty } from './generated/metrics_pb.js';
import { MetricsClient } from './generated/metrics_grpc_web_pb.js';

// Create a client instance. replace with actual gRPC web server endpoint
// Maybe supply from Angular
const client = new MetricsClient('http://' + window.location.hostname + ':8080', null, null);

function startCPUClockStream() {
    const request = new Empty();
    const stream = client.streamCPUClock(request, {});

    stream.on('data', (response) => {
        const cpu_clocks = response.getCpuClockList();

        // Communicate with Service Worker
        if (navigator.serviceWorker && navigator.serviceWorker.controller) {
            navigator.serviceWorker.controller.postMessage({
                type: 'cpu_clock',
                payload: cpu_clocks,
            });
        }

        console.debug('grpc-web-client: cpu_clocks.length', cpu_clocks.length)
    });

    stream.on('error', (err) => {
        if (navigator.serviceWorker && navigator.serviceWorker.controller) {
            navigator.serviceWorker.controller.postMessage({
                type: 'clock_stream_error',
                payload: { code: err.code, message: err.message },
            });
        }

        console.log('CPU Clock Stream Error:', err)
    });

    stream.on('end', () => {
        if (navigator.serviceWorker && navigator.serviceWorker.controller) {
            navigator.serviceWorker.controller.postMessage({ type: 'clock_stream_end' });
        }

        console.log('CPU Clock Stream ended')
    });
}

function startCPUUtilizationStream() {
    const request = new Empty();
    const stream = client.streamCPUUtilization(request, {});

    stream.on('data', (response) => {
        const cpu_utilization = response.getCpuUtilizationList();

        // Communicate with Service Worker
        if (navigator.serviceWorker && navigator.serviceWorker.controller) {
            navigator.serviceWorker.controller.postMessage({
                type: 'cpu_utilization',
                payload: cpu_utilization,
            });
        }

        console.debug('grpc-web-client: cpu_utilization.length', cpu_utilization.length)
    });

    stream.on('error', (err) => {
        if (navigator.serviceWorker && navigator.serviceWorker.controller) {
            navigator.serviceWorker.controller.postMessage({
                type: 'utilization_stream_error',
                payload: { code: err.code, message: err.message },
            });
        }

        console.log('CPU Utilization Stream Error:', err)
    });

    stream.on('end', () => {
        if (navigator.serviceWorker && navigator.serviceWorker.controller) {
            navigator.serviceWorker.controller.postMessage({ type: 'utilization_stream_end' });
        }

        console.log('CPU Utilization Stream ended')
    });
}

function startHwmonStream() {
    const request = new Empty();
    const stream = client.streamHwmon(request, {});

    stream.on('data', (response) => {
        // Extract data from the Hwmon message
        const hwmonData = {
            name: response.getName(),
            composite: response.getComposite(),
            sensors: response.getSensorsList().map(sensor => ({
                name: sensor.getName(),
                temp: sensor.getTemp()
            }))
        };

        // Communicate with Service Worker
        if (navigator.serviceWorker && navigator.serviceWorker.controller) {
            navigator.serviceWorker.controller.postMessage({
                type: 'hwmon',
                payload: hwmonData
            });
        }

        console.log('grpc-web-client: hwmon length', Object.keys(hwmonData).length);
    });

    stream.on('error', (err) => {
        if (navigator.serviceWorker && navigator.serviceWorker.controller) {
            navigator.serviceWorker.controller.postMessage({
                type: 'hwmon_stream_error',
                payload: { code: err.code, message: err.message }
            });
        }

        console.log('Hwmon Stream Error:', err);
    });

    stream.on('end', () => {
        if (navigator.serviceWorker && navigator.serviceWorker.controller) {
            navigator.serviceWorker.controller.postMessage({
                type: 'hwmon_stream_end'
            });
        }

        console.log('Hwmon Stream ended');
    });
}

// Start streams when the script loads
// might consider trigger for start/stop from button clicks or other events
startCPUUtilizationStream();
startCPUClockStream();
startHwmonStream();


// Export functions if you want to use them in other modules
export default {
    startCPUUtilizationStream,
    startCPUClockStream,
    startHwmonStream,
};
