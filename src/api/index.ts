export async function sendMessage(message: string) {
  const result = await fetch("http://localhost:8090/query", {
    method: "POST",
    body: JSON.stringify({
      message: message,
    }),
  });
  return result.json();
}
