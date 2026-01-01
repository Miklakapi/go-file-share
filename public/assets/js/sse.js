export function useSSE() {
    const es = new EventSource("http://localhost:8080/api/v1/sse")

    es.onmessage = event => {
        console.log("Received:", event.data)
    }

    es.addEventListener("time", event => {
        console.log("TIME:", event.data);
    })

    es.onerror = err => {
        console.error("EventSource failed:", err);
    }
}