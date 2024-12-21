self.addEventListener('message', (event) => {
  const { type, payload } = event.data;

  console.log("NGSW", type) // doesn't get printed on built asset or non-localhost addresses

  switch (type) {
    case 'cpu_utilization':
      broadcastMessage('cpu_utilization', payload);
      break;
    case 'cpu_clock':
      broadcastMessage('cpu_clock', payload);
      break;
    case 'hwmon':
      broadcastMessage('hwmon', payload);
      break;
    case 'stream_error':
      broadcastMessage('stream_error', payload);
      break;
    case 'stream_end':
      broadcastMessage('stream_end');
      break;
  }
});

// Broadcast data to all clients (Angular app)
function broadcastMessage(type, payload = null) {
  self.clients.matchAll({ includeUncontrolled: true }).then((clients) => {
    clients.forEach((client) => {
      client.postMessage({ type, payload });
    });
  });
}
