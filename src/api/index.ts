export async function sendMessage(message: string) {
  const result = await fetch("http://localhost:8090/query", {
    method: "POST",
    body: JSON.stringify({
      message: message,
    }),
  });
  return result.json();
}

export function sendMessageStream(
  message: string,
  onMessage: (data: string) => void,
  onError: (error: any) => void
) {
  const eventSource = new EventSource(
    `http://localhost:8090/chat?message=${encodeURIComponent(message)}`
  );

  eventSource.onmessage = (event) => {
    const data = JSON.parse(event.data);
    onMessage(data.content);
  };

  eventSource.onerror = (error) => {
    eventSource.close();
    onError(error);
  };

  return eventSource;
}

export function sendSSEMessage(message: string): EventSource {
  const eventSource = new EventSource(
    `http://localhost:8090/chat?message=${encodeURIComponent(message)}`
  );

  return eventSource;
}
