import { api, filenameFromDisposition, triggerBrowserDownload } from "./helpers.js"

export function useFiles() {
    async function get(roomId) {
        const files = (await api(`/rooms/${roomId}/files`)).data
        return files
    }

    async function download(roomId, fileId) {
        const res = await api(`/rooms/${roomId}/files/${fileId}/download`, {
            responseType: 'blob'
        })

        const blob = res.data
        const filename = filenameFromDisposition(res.headers?.['content-disposition']
            || res.headers?.get?.('content-disposition'))
            || `file-${fileId}`

        triggerBrowserDownload(blob, filename)

        return true
    }

    async function upload(roomId, file) {
        if (!file) throw new Error('No file provided')

        const form = new FormData()
        form.append('file', file, file.name)

        const res = await api(`/rooms/${roomId}/files`, {
            method: 'POST',
            body: form
        })

        return res.data.id
    }

    async function remove(roomId, fileId) {
        if (!confirm('Delete this file?')) return false
        await api(`/rooms/${roomId}/files/${fileId}`, { method: 'DELETE' })
        return true
    }

    return {
        get,
        download,
        upload,
        remove
    }
}
