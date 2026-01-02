export function useSSE() {
    if (typeof EventSource === "undefined") {
        throw new Error("EventSource not supported in this browser")
    }

    const es = new EventSource("/api/v1/sse")

    function onMessage(handler) {
        es.onmessage = handler
    }

    function onEvent(eventName, handler) {
        es.addEventListener(eventName, handler)
        return () => es.removeEventListener(eventName, handler)
    }

    function onError(handler) {
        es.onerror = handler
    }

    function close() {
        es.close()
    }

    return {
        onMessage,
        onEvent,
        onError,
        close
    }
}