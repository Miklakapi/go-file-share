import { api } from "./helpers.js"

export function useFiles() {
    async function get(roomId) {
        const files = (await api(`/rooms/${roomId}/files`)).data
        return files
    }

    async function download(roomId, fileId) {
        return (await api(`/rooms/${roomId}/files/${fileId}/download`)).data
    }

    async function upload(roomId, file) {
        const response = await api(`/rooms/${roomId}/file`, {
            method: 'POST'
        })
        return response.data.ID
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
