export async function api(path, options = {}, timeout = 10000) {
    const url = '/api/v1' + path
    const isFormData = options?.body instanceof FormData
    const wantsBlob = options?.responseType === 'blob'
    const headers = new Headers(options.headers || {})

    if (!isFormData && !headers.has('Content-Type')) {
        headers.set('Content-Type', 'application/json')
    }

    let res
    try {
        res = await fetch(url, {
            credentials: 'include',
            headers,
            ...options,
            signal: AbortSignal.timeout(timeout)
        })
    } catch {
        throw Error('Request timeout')
    }

    if (!res.ok) {
        let message = res.statusText
        try {
            const data = await res.json()
            message = data?.message || message
        } catch {
            try {
                const text = await res.text()
                if (text) message = text
            } catch { }
        }
        throw new Error(message)
    }

    if (res.status === 204) return null
    if (wantsBlob) {
        return {
            data: await res.blob(),
            headers: res.headers
        }
    }
    return res.json()
}

export function formatDate(iso) {
    try {
        return new Date(iso).toLocaleString()
    } catch {
        return iso
    }
}

export function shortId(id) {
    if (!id) return ''
    return id.length > 10 ? `${id.slice(0, 8)}…` : id
}

export async function sleep(ms) {
    return new Promise(resolve => setTimeout(resolve, ms))
} 
