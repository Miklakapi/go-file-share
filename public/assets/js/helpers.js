export async function api(path, options = {}) {
    const res = await fetch('/api/v1' + path, {
        credentials: 'include',
        headers: { 'Content-Type': 'application/json' },
        ...options,
    })

    if (!res.ok) {
        let message = res.statusText
        try {
            const data = await res.json()
            message = data?.message || message
        } catch {
            const text = await res.text()
            if (text) message = text
        }
        throw new Error(message)
    }

    return res.status === 204 ? null : res.json()
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
