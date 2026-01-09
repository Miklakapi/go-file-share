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
        const signal = timeout ? AbortSignal.timeout(timeout) : undefined
        res = await fetch(url, {
            credentials: 'include',
            headers,
            signal,
            ...options,
        })
    } catch (error) {
        if (options.signal?.aborted) {
            throw new DOMException('Aborted', 'AbortError');
        }
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

export function generateNumericCode(length = 12) {
    const bytes = new Uint8Array(length)
    crypto.getRandomValues(bytes)

    let out = ''
    for (let i = 0; i < length; i++) {
        out += (bytes[i] % 10).toString()
    }

    return out
}

export function triggerBrowserDownload(blob, filename) {
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = filename || 'download'
    document.body.appendChild(a)
    a.click()
    a.remove()
    URL.revokeObjectURL(url)
}

export function filenameFromDisposition(disposition) {
    if (!disposition) return null

    const match = disposition.match(/filename\*?=(?:UTF-8'')?("?)([^";]+)\1/i)
    if (!match) return null

    try {
        return decodeURIComponent(match[2])
    } catch {
        return match[2]
    }
}
