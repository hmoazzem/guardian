import { MetricsClient } from './generated/MetricsServiceClientPb';
import { Empty, CPUClock } from './generated/metrics_pb';

async function main() {
  // Create a client instance
  // Replace with actual gRPC-web server endpoint
  const client = new MetricsClient('http://' + window.location.hostname + ':8080', null, null);

  // Create an empty request object
  const request = new Empty();

  try {
    // Call the StreamCPUClock server-streaming method
    const stream = client.streamCPUClock(request, {});

    // Handle the stream
    stream.on('data', (cpuClock: CPUClock) => {
      // Process each CPU clock data point
      const clockList = cpuClock.getCpuClockList();
      console.log('Received CPU Clock:', clockList);
    });

    stream.on('error', (err: Error) => {
      console.error('Stream error:', err);
    });

    stream.on('end', () => {
      console.log('Stream ended');
    });
  } catch (error) {
    console.error('Error calling StreamCPUClock:', error);
  }
}

// Run the main function
main().catch(console.error);
